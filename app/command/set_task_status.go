package command

import (
	"context"

	"github.com/pietjan/tm/domain/task"
)

type SetTaskStatus struct {
	TaskUID string
	Status  int
}

type SetTaskStatusHandler struct {
	repo task.WriteModel
}

func NewSetTaskStatusHandler(repo task.Repository) SetTaskStatusHandler {
	if repo == nil {
		panic("nil repo")
	}

	return SetTaskStatusHandler{repo: repo}
}

func (h SetTaskStatusHandler) Handle(ctx context.Context, cmd SetTaskStatus) (err error) {
	return h.repo.UpdateTask(ctx, cmd.TaskUID,
		func(ctx context.Context, t *task.Task) error {
			return t.Update(
				task.Status(cmd.Status),
			)
		})
}
