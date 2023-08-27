package timer

import (
	"context"
	"time"
)

type UpdateFn = func(ctx context.Context, s *Timer) error

type ReadModel interface {
	GetTimer(ctx context.Context, timerUID string) (timer *Timer, err error)
	// AllTimers(ctx context.Context, filter *filter.Filter) (timers []*Timer, err error)
}

type WriteModel interface {
	AddTimer(ctx context.Context, timer *Timer) error
	UpdateTimer(ctx context.Context, timerUID string, updateFn UpdateFn) error
}

type Repository interface {
	ReadModel
	WriteModel
}

type Option = func(*Timer)

type Timer struct {
	uid   string
	task  string
	start time.Time
	end   time.Time
}

func New(options ...Option) (*Timer, error) {
	t := &Timer{}

	if err := apply(t, options...); err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Timer) Update(options ...Option) error {
	return apply(t, options...)
}

func apply(t *Timer, options ...Option) error {
	for _, fn := range options {
		fn(t)
	}

	return nil
}

func UID(uid string) Option {
	return func(t *Timer) {
		t.uid = uid
	}
}

func (t Timer) UID() string {
	return t.uid
}

func Task(task string) Option {
	return func(t *Timer) {
		t.task = task
	}
}

func (t Timer) Task() string {
	return t.task
}

func Start(start time.Time) Option {
	return func(t *Timer) {
		t.start = start
	}
}

func (t Timer) Start() time.Time {
	return t.start
}

func End(end time.Time) Option {
	return func(t *Timer) {
		t.end = end
	}
}

func (t Timer) End() time.Time {
	return t.end
}

func (t Timer) Duration() time.Duration {
	end := time.Now()
	if !t.end.IsZero() {
		end = t.end
	}

	return end.Sub(t.start)
}
