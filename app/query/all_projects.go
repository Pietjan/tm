package query

import (
	"context"

	"github.com/pietjan/tm/domain/filter"
	"github.com/pietjan/tm/domain/project"
)

type AllProjects struct {
	Filter *filter.Filter
}

type AllProjectsHandler struct {
	readModel project.ReadModel
}

func NewAllProjectsHandler(readModel project.ReadModel) AllProjectsHandler {
	if readModel == nil {
		panic(`nil readModel`)
	}

	return AllProjectsHandler{readModel: readModel}
}

func (h AllProjectsHandler) Handle(ctx context.Context, query AllProjects) (p []*project.Project, err error) {
	p, err = h.readModel.AllProjects(ctx, query.Filter)
	return
}
