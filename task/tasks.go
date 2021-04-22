/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package task

import (
	"context"
	"sync"

	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/runtime"
	"github.com/pkg/errors"
)

type Tasks interface {
	Init(context.Context) error
	Delete(context.Context, runtime.Task, func() error) error
}

// NewtaskList returns a new taskList
func NewTaskList() *TaskList {
	return &TaskList{
		tasks: make(map[string]map[string]runtime.Task),
	}
}

// taskList holds and provides locking around tasks
type TaskList struct {
	mu    sync.Mutex
	tasks map[string]map[string]runtime.Task
}

// Get a task
func (l *TaskList) Get(ctx context.Context, id string) (runtime.Task, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil {
		return nil, err
	}
	tasks, ok := l.tasks[namespace]
	if !ok {
		return nil, runtime.ErrTaskNotExists
	}
	t, ok := tasks[id]
	if !ok {
		return nil, runtime.ErrTaskNotExists
	}
	return t, nil
}

// GetAll tasks under a namespace
func (l *TaskList) GetAll(ctx context.Context, noNS bool) ([]runtime.Task, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	var o []runtime.Task
	if noNS {
		for ns := range l.tasks {
			for _, t := range l.tasks[ns] {
				o = append(o, t)
			}
		}
		return o, nil
	}
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil {
		return nil, err
	}
	tasks, ok := l.tasks[namespace]
	if !ok {
		return o, nil
	}
	for _, t := range tasks {
		o = append(o, t)
	}
	return o, nil
}

// Add a task
func (l *TaskList) Add(ctx context.Context, t runtime.Task) error {
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil {
		return err
	}
	return l.AddWithNamespace(namespace, t)
}

// AddWithNamespace adds a task with the provided namespace
func (l *TaskList) AddWithNamespace(namespace string, t runtime.Task) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	id := t.ID()
	if _, ok := l.tasks[namespace]; !ok {
		l.tasks[namespace] = make(map[string]runtime.Task)
	}
	if _, ok := l.tasks[namespace][id]; ok {
		return errors.Wrap(runtime.ErrTaskAlreadyExists, id)
	}
	l.tasks[namespace][id] = t
	return nil
}

// Delete a task
func (l *TaskList) Delete(ctx context.Context, t runtime.Task) {
	l.mu.Lock()
	defer l.mu.Unlock()
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil {
		return
	}
	tasks, ok := l.tasks[namespace]
	if ok {
		id := t.ID()
		task, exists := tasks[id]
		if exists && task == t {
			delete(tasks, id)
		}
	}
}