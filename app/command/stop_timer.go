package command

import (
	"context"
	"time"

	"github.com/pietjan/tm/domain/timer"
)

type StopTimer struct {
	TimerUID string
}

type StopTimerHandler struct {
	repo timer.WriteModel
}

func NewStopTimerHandler(repo timer.Repository) StopTimerHandler {
	if repo == nil {
		panic("nil repo")
	}

	return StopTimerHandler{repo: repo}
}

func (h StopTimerHandler) Handle(ctx context.Context, cmd StopTimer) (err error) {
	return h.repo.UpdateTimer(ctx, cmd.TimerUID,
		func(ctx context.Context, t *timer.Timer) error {
			return t.Update(
				timer.End(time.Now()),
			)
		})
}
