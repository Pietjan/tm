package project

import (
	"context"

	"github.com/pietjan/tm/domain/filter"
)

type ReadModel interface {
	AllProjects(ctx context.Context, filter *filter.Filter) (projects []*Project, err error)
}

type WriteModel interface {
	AddProject(ctx context.Context, project *Project) error
}

type Repository interface {
	ReadModel
	WriteModel
}

type Option = func(*Project)

type Project struct {
	uid  string
	ref  string
	name string
}

func New(options ...Option) (*Project, error) {
	p := &Project{}

	for _, fn := range options {
		fn(p)
	}

	return p, nil
}

func UID(uid string) Option {
	return func(p *Project) {
		p.uid = uid
	}
}

func (p Project) UID() string {
	return p.uid
}

func Ref(ref string) Option {
	return func(p *Project) {
		p.ref = ref
	}
}

func (p Project) Ref() string {
	return p.ref
}

func Name(name string) Option {
	return func(p *Project) {
		p.name = name
	}
}

func (p Project) Name() string {
	return p.name
}
