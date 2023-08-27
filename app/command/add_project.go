package command

import (
	"context"

	"github.com/pietjan/tm/domain/project"
)

type AddProject struct {
	ProjectUID string
	Name       string
}

type AddProjectHandler struct {
	repo project.WriteModel
}

func NewAddProjectHandler(repo project.Repository) AddProjectHandler {
	if repo == nil {
		panic("nil repo")
	}

	return AddProjectHandler{repo: repo}
}

func (h AddProjectHandler) Handle(ctx context.Context, cmd AddProject) (err error) {
	t, err := project.New(
		project.UID(cmd.ProjectUID),
		project.Name(cmd.Name),
	)
	if err != nil {
		return err
	}

	return h.repo.AddProject(ctx, t)
}
