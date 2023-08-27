package query

import (
	"context"
	"time"

	"github.com/pietjan/tm/domain/filter"
)

type Task struct {
	UID      string
	Ref      string
	Name     string
	Project  string
	Status   int
	Duration time.Duration
}

type AllTasks struct {
	Filter *filter.Filter
}

type AllTasksHandler struct {
	readModel AllTasksReadModel
}

type AllTasksReadModel interface {
	AllTasks(ctx context.Context, filter *filter.Filter) (tasks []Task, err error)
}

func NewAllTasksHandler(readModel AllTasksReadModel) AllTasksHandler {
	if readModel == nil {
		panic(`nil readModel`)
	}

	return AllTasksHandler{readModel: readModel}
}

func (h AllTasksHandler) Handle(ctx context.Context, query AllTasks) (t []Task, err error) {
	t, err = h.readModel.AllTasks(ctx, query.Filter)
	return
}
