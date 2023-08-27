package command

import (
	"context"
	"time"

	"github.com/pietjan/tm/domain/timer"
)

type StartTimer struct {
	TimerUID string
	TaskUID  string
}

type StartTimerHandler struct {
	repo timer.WriteModel
}

func NewStartTimerHandler(repo timer.Repository) StartTimerHandler {
	if repo == nil {
		panic("nil repo")
	}

	return StartTimerHandler{repo: repo}
}

func (h StartTimerHandler) Handle(ctx context.Context, cmd StartTimer) (err error) {
	t, err := timer.New(
		timer.UID(cmd.TimerUID),
		timer.Task(cmd.TaskUID),
		timer.Start(time.Now()),
	)
	if err != nil {
		return err
	}

	return h.repo.AddTimer(ctx, t)
}
