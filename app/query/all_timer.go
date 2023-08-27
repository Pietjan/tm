package query

import (
	"context"
	"time"

	"github.com/pietjan/tm/domain/filter"
)

type Timer struct {
	UID   string
	Task  string
	Start time.Time
	End   time.Time
}

type AllTimers struct {
	Filter *filter.Filter
}

type AllTimersHandler struct {
	readModel AllTimersReadModel
}

type AllTimersReadModel interface {
	AllTimers(ctx context.Context, filter *filter.Filter) (timers []Timer, err error)
}

func NewAllTimersHandler(readModel AllTimersReadModel) AllTimersHandler {
	if readModel == nil {
		panic(`nil readModel`)
	}

	return AllTimersHandler{readModel: readModel}
}

func (h AllTimersHandler) Handle(ctx context.Context, query AllTimers) (t []Timer, err error) {
	t, err = h.readModel.AllTimers(ctx, query.Filter)
	return
}
