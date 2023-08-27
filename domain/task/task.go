package task

import (
	"context"
)

type UpdateFn = func(ctx context.Context, t *Task) error

type ReadModel interface {
	GetTask(ctx context.Context, uid string) (task *Task, err error)
	// AllTasks(ctx context.Context, filter *filter.Filter) (tasks []*Task, err error)
}

type WriteModel interface {
	AddTask(ctx context.Context, task *Task) error
	UpdateTask(ctx context.Context, timerUID string, updateFn UpdateFn) error
}

type Repository interface {
	ReadModel
	WriteModel
}

type Option = func(*Task)

type Task struct {
	uid     string
	ref     string
	name    string
	status  int
	project string
}

func New(options ...Option) (*Task, error) {
	t := &Task{}

	err := apply(t, options...)

	return t, err
}

func (t *Task) Update(options ...Option) error {
	return apply(t, options...)
}

func apply(t *Task, options ...Option) error {
	for _, fn := range options {
		fn(t)
	}

	return nil
}

func UID(uid string) Option {
	return func(t *Task) {
		t.uid = uid
	}
}

func (t Task) UID() string {
	return t.uid
}

func Ref(ref string) Option {
	return func(t *Task) {
		t.ref = ref
	}
}

func (t Task) Ref() string {
	return t.ref
}

func Name(name string) Option {
	return func(t *Task) {
		t.name = name
	}
}

func (t Task) Name() string {
	return t.name
}

func Status(status int) Option {
	return func(t *Task) {
		t.status = status
	}
}

func (t Task) Status() int {
	return t.status
}

func Project(project string) Option {
	return func(t *Task) {
		t.project = project
	}
}

func (t Task) Project() string {
	return t.project
}
