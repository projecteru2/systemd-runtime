package task

import (
	"context"
	"sync"
	"time"

	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/runtime"
	taskapi "github.com/containerd/containerd/runtime/v2/task"
	"github.com/containerd/ttrpc"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/projecteru2/systemd-runtime/common"
	"github.com/projecteru2/systemd-runtime/shim"
)

const (
	maxRetry     = 10
	waitInterval = time.Duration(120) * time.Second
	connTimeout  = time.Duration(10) * time.Second
)

// ServiceSubscribe .
type ServiceSubscribe struct {
	service *Service
	exit    *runtime.Exit
}

type ConnMng struct {
	sync.Mutex
	cancelConn    CancelConn
	disabled      bool
	logger        *logrus.Entry
	lastSubscribe ServiceSubscribe
	bundle        common.Bundle
	subscribers   []chan<- ServiceSubscribe
}

func (s *ConnMng) getService() *Service {
	s.Lock()
	defer s.Unlock()

	return s.lastSubscribe.service
}

// return error either task is paused or disabled
func (s *ConnMng) getRunningService(_ context.Context) (*Service, error) {
	s.Lock()
	defer s.Unlock()

	if s.lastSubscribe.service != nil {
		return s.lastSubscribe.service, nil
	}

	if s.lastSubscribe.exit != nil {
		return nil, ErrTaskIsKilled
	}
	return nil, ErrTaskIsPaused
}

// wait service connection, if service is exited, return exit status
func (s *ConnMng) subscribeService(ctx context.Context) (sub ServiceSubscribe, err error) {
	ch := s.subscribe()

	select {
	case <-ctx.Done():
		return sub, ctx.Err()
	case <-time.After(connTimeout):
		return sub, ErrTaskIsPaused
	case service := <-ch:
		return service, nil
	}
}

func (s *ConnMng) waitServiceConnected(ctx context.Context) (*Service, error) {
	serviceSubscribe, err := s.subscribeService(ctx)
	if err != nil {
		return nil, err
	}
	if serviceSubscribe.exit != nil {
		return nil, ErrTaskIsKilled
	}
	return serviceSubscribe.service, nil
}

func (s *ConnMng) Disable() {
	s.Lock()
	defer s.Unlock()

	s.disabled = true
}

func (s *ConnMng) Disabled(exit *runtime.Exit) {
	subscribers := func() []chan<- ServiceSubscribe {
		s.Lock()
		defer s.Unlock()

		s.cancelConn.Cancel()

		s.disabled = true
		s.lastSubscribe = ServiceSubscribe{
			exit: exit,
		}
		subscribers := s.subscribers
		s.subscribers = nil
		return subscribers
	}()

	for _, sub := range subscribers {
		sub <- ServiceSubscribe{
			exit: exit,
		}
		close(sub)
	}
}

func (s *ConnMng) subscribe() <-chan ServiceSubscribe {
	ch := make(chan ServiceSubscribe, 1)
	s.Lock()
	defer s.Unlock()

	if s.lastSubscribe.service != nil {
		ch <- ServiceSubscribe{
			service: s.lastSubscribe.service,
		}
		close(ch)
		return ch
	}

	if s.lastSubscribe.exit != nil {
		ch <- ServiceSubscribe{
			exit: s.lastSubscribe.exit,
		}
		close(ch)
		return ch
	}

	s.subscribers = append(s.subscribers, ch)
	return ch
}

func (s *ConnMng) reset() {
	s.Lock()
	defer s.Unlock()

	s.lastSubscribe = ServiceSubscribe{}
}

func (s *ConnMng) setService(service *Service) {
	subscribers := func() []chan<- ServiceSubscribe {
		s.Lock()
		defer s.Unlock()

		s.lastSubscribe = ServiceSubscribe{
			service: service,
		}
		subscribers := s.subscribers
		s.subscribers = nil

		return subscribers
	}()

	for _, sub := range subscribers {
		sub <- ServiceSubscribe{
			service: service,
		}
		close(sub)
	}
}

func (s *ConnMng) connect() {
	retryCount := 0
	for {
		if func() bool {
			s.Lock()
			defer s.Unlock()

			return s.disabled
		}() {
			return
		}

		// if containerd crushes and failed to catch a pause event, we can make up here
		// also we can clean up crushed shim process here
		// there will be a delay with docker ps for paused status, because docker will not handle the event on startup
		// docker will receive pause event on next connect loop
		err := s.bundle.Cleanup(context.Background())
		if err != nil {
			s.logger.WithError(err).Error("cleanup bundle before connect error")
		}

		status, running, err := s.bundle.ShimStatus(context.Background())
		if err != nil {
			logrus.WithError(err).Error("check shim status error")
			retryCount++
		}
		if retryCount >= maxRetry {
			if err := s.bundle.Delete(context.Background()); err != nil {
				logrus.WithError(err).Error("delete bundle error")
			}
			return
		}
		if status.Disabled && !running {
			s.Disabled(status.Exit)
			return
		}
		if s.doConnect(waitInterval, connTimeout) {
			return
		}
	}
}

func (s *ConnMng) getAddr(timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	s.cancelConn.Set(cancel)
	defer s.cancelConn.Cancel()

	return common.ReceiveAddressOverFifo(ctx, s.bundle.Path())
}

func (s *ConnMng) connOnAddr(timeout time.Duration, addr string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.G(ctx).Debug("create shim connection")
	conn, err := common.Connect(addr, shim.AnonReconnectDialer)
	if err != nil {
		return errors.WithMessage(err, "make shim connection failed")
	}

	client := ttrpc.NewClient(conn, ttrpc.WithOnClose(s.reconnect))
	taskClient := taskapi.NewTaskClient(client)

	log.G(ctx).Debug("connect shim service")
	taskPid, err := connect(ctx, taskClient, s.bundle.ID())
	if err != nil {
		return errors.WithMessage(err, "connect task service failed")
	}

	se := &Service{
		taskPid:     taskPid,
		id:          s.bundle.ID(),
		namespace:   s.bundle.Namespace(),
		taskService: taskClient,
	}
	s.setService(se)
	return nil
}

func (s *ConnMng) doConnect(waitInterval time.Duration, connTimeout time.Duration) bool {

	addr, err := s.getAddr(waitInterval)
	if err != nil {
		logrus.WithError(err).Error("receive address over fifo error")
		return false
	}

	if err := s.connOnAddr(connTimeout, addr); err != nil {
		s.logger.WithField(
			"address", addr,
		).WithError(
			err,
		).Error("connect on address error")
		return false
	}
	return true
}

func (s *ConnMng) reconnect() {
	s.reset()
	err := s.bundle.Cleanup(context.Background())
	if err != nil {
		s.logger.WithError(err).Error("cleanup bundle error")
	}

	s.logger.Debug("ttrpc conn disconnected, reconnect")
	go s.connect()
}

type CancelConn struct {
	sync.Mutex
	cancel func()
}

func (c *CancelConn) Set(cancel func()) {
	c.Lock()
	defer c.Unlock()
	c.cancel = cancel
}

func (c *CancelConn) Cancel() {
	c.Lock()
	defer c.Unlock()

	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
	}
}
