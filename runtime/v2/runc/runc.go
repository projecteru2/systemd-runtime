package runc

import (
	"context"
	"io"
	"net"
	"os"
	"os/exec"
	"time"

	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// Container hold information for a runc container
type Container struct {
	ID          string            `json:"id"`
	Pid         int               `json:"pid"`
	Status      string            `json:"status"`
	Bundle      string            `json:"bundle"`
	Rootfs      string            `json:"rootfs"`
	Created     time.Time         `json:"created"`
	Annotations map[string]string `json:"annotations"`
}

// Socket is a unix socket that accepts the pty master created by runc
type Socket struct {
	rmdir bool
	l     *net.UnixListener
}

type IO interface {
	io.Closer
	Stdin() io.WriteCloser
	Stdout() io.ReadCloser
	Stderr() io.ReadCloser
	Set(*exec.Cmd)
}

type ConsoleSocket interface {
	Path() string
}

type CreateOpts struct {
	IO
	// PidFile is a path to where a pid file should be created
	PidFile       string
	ConsoleSocket ConsoleSocket
	Detach        bool
	NoPivot       bool
	NoNewKeyring  bool
	ExtraFiles    []*os.File
}

type Exit struct {
	Timestamp time.Time
	Pid       int
	Status    int
}

type ExecOpts struct {
	IO
	PidFile       string
	ConsoleSocket ConsoleSocket
	Detach        bool
}

type DeleteOpts struct {
	Force bool
}

type KillOpts struct {
	All bool
}

type Stats struct {
	Cpu     Cpu                `json:"cpu"`
	Memory  Memory             `json:"memory"`
	Pids    Pids               `json:"pids"`
	Blkio   Blkio              `json:"blkio"`
	Hugetlb map[string]Hugetlb `json:"hugetlb"`
}

type Event struct {
	// Type are the event type generated by runc
	// If the type is "error" then check the Err field on the event for
	// the actual error
	Type  string `json:"type"`
	ID    string `json:"id"`
	Stats *Stats `json:"data,omitempty"`
	// Err has a read error if we were unable to decode the event from runc
	Err error `json:"-"`
}

type Hugetlb struct {
	Usage   uint64 `json:"usage,omitempty"`
	Max     uint64 `json:"max,omitempty"`
	Failcnt uint64 `json:"failcnt"`
}

type BlkioEntry struct {
	Major uint64 `json:"major,omitempty"`
	Minor uint64 `json:"minor,omitempty"`
	Op    string `json:"op,omitempty"`
	Value uint64 `json:"value,omitempty"`
}

type Blkio struct {
	IoServiceBytesRecursive []BlkioEntry `json:"ioServiceBytesRecursive,omitempty"`
	IoServicedRecursive     []BlkioEntry `json:"ioServicedRecursive,omitempty"`
	IoQueuedRecursive       []BlkioEntry `json:"ioQueueRecursive,omitempty"`
	IoServiceTimeRecursive  []BlkioEntry `json:"ioServiceTimeRecursive,omitempty"`
	IoWaitTimeRecursive     []BlkioEntry `json:"ioWaitTimeRecursive,omitempty"`
	IoMergedRecursive       []BlkioEntry `json:"ioMergedRecursive,omitempty"`
	IoTimeRecursive         []BlkioEntry `json:"ioTimeRecursive,omitempty"`
	SectorsRecursive        []BlkioEntry `json:"sectorsRecursive,omitempty"`
}

type Pids struct {
	Current uint64 `json:"current,omitempty"`
	Limit   uint64 `json:"limit,omitempty"`
}

type Throttling struct {
	Periods          uint64 `json:"periods,omitempty"`
	ThrottledPeriods uint64 `json:"throttledPeriods,omitempty"`
	ThrottledTime    uint64 `json:"throttledTime,omitempty"`
}

type CpuUsage struct {
	// Units: nanoseconds.
	Total  uint64   `json:"total,omitempty"`
	Percpu []uint64 `json:"percpu,omitempty"`
	Kernel uint64   `json:"kernel"`
	User   uint64   `json:"user"`
}

type Cpu struct {
	Usage      CpuUsage   `json:"usage,omitempty"`
	Throttling Throttling `json:"throttling,omitempty"`
}

type MemoryEntry struct {
	Limit   uint64 `json:"limit"`
	Usage   uint64 `json:"usage,omitempty"`
	Max     uint64 `json:"max,omitempty"`
	Failcnt uint64 `json:"failcnt"`
}

type Memory struct {
	Cache     uint64            `json:"cache,omitempty"`
	Usage     MemoryEntry       `json:"usage,omitempty"`
	Swap      MemoryEntry       `json:"swap,omitempty"`
	Kernel    MemoryEntry       `json:"kernel,omitempty"`
	KernelTCP MemoryEntry       `json:"kernelTCP,omitempty"`
	Raw       map[string]uint64 `json:"raw,omitempty"`
}

type TopResults struct {
	// Processes running in the container, where each is process is an array of values corresponding to the headers
	Processes [][]string `json:"Processes"`

	// Headers are the names of the columns
	Headers []string `json:"Headers"`
}

type CheckpointOpts struct {
	// ImagePath is the path for saving the criu image file
	ImagePath string
	// WorkDir is the working directory for criu
	WorkDir string
	// ParentPath is the path for previous image files from a pre-dump
	ParentPath string
	// AllowOpenTCP allows open tcp connections to be checkpointed
	AllowOpenTCP bool
	// AllowExternalUnixSockets allows external unix sockets to be checkpointed
	AllowExternalUnixSockets bool
	// AllowTerminal allows the terminal(pty) to be checkpointed with a container
	AllowTerminal bool
	// CriuPageServer is the address:port for the criu page server
	CriuPageServer string
	// FileLocks handle file locks held by the container
	FileLocks bool
	// Cgroups is the cgroup mode for how to handle the checkpoint of a container's cgroups
	Cgroups CgroupMode
	// EmptyNamespaces creates a namespace for the container but does not save its properties
	// Provide the namespaces you wish to be checkpointed without their settings on restore
	EmptyNamespaces []string
}

type CgroupMode string

type CheckpointAction func([]string) []string

type RestoreOpts struct {
	CheckpointOpts
	IO

	Detach        bool
	PidFile       string
	NoSubreaper   bool
	NoPivot       bool
	ConsoleSocket ConsoleSocket
}

type Version struct {
	Runc   string
	Commit string
	Spec   string
}

// Runc .
type Runc interface {
	List(context context.Context) ([]*Container, error)
	State(context context.Context, id string) (*Container, error)
	Create(context context.Context, id, bundle string, opts *CreateOpts) error
	Start(context context.Context, id string) error
	Exec(context context.Context, id string, spec specs.Process, opts *ExecOpts) error
	Run(context context.Context, id, bundle string, opts *CreateOpts) (int, error)
	Delete(context context.Context, id string, opts *DeleteOpts) error
	Kill(context context.Context, id string, sig int, opts *KillOpts) error
	Stats(context context.Context, id string) (*Stats, error)
	Events(context context.Context, id string, interval time.Duration) (chan *Event, error)
	Pause(context context.Context, id string) error
	Resume(context context.Context, id string) error
	Ps(context context.Context, id string) ([]int, error)
	Top(context context.Context, id string, psOptions string) (*TopResults, error)
	Checkpoint(context context.Context, id string, opts *CheckpointOpts, actions ...CheckpointAction) error
	Restore(context context.Context, id, bundle string, opts *RestoreOpts) (int, error)
	Update(context context.Context, id string, resources *specs.LinuxResources) error
	Version(context context.Context) (Version, error)
}