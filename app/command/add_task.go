package command

import (
	"context"

	"github.com/pietjan/tm/domain/task"
)

type AddTask struct {
	TaskUID string
	Project string
	Name    string
}

type AddTaskHandler struct {
	repo task.WriteModel
}

func NewAddTaskHandler(repo task.Repository) AddTaskHandler {
	if repo == nil {
		panic("nil repo")
	}

	return AddTaskHandler{repo: repo}
}

func (h AddTaskHandler) Handle(ctx context.Context, cmd AddTask) (err error) {
	t, err := task.New(
		task.UID(cmd.TaskUID),
		task.Project(cmd.Project),
		task.Name(cmd.Name),
	)
	if err != nil {
		return err
	}

	return h.repo.AddTask(ctx, t)
}
