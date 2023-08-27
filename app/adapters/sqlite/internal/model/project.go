package model

import (
	"database/sql"

	"github.com/pietjan/tm/domain/project"

	"git.ultraware.nl/NiseVoid/qb"
	"github.com/pietjan/model"
)

type ProjectModel struct {
	model.Model
}

func (m ProjectModel) UID() qb.Field {
	return m.Field(uidField)
}

func (m ProjectModel) Ref() qb.Field {
	return m.Field(refField)
}

func (m ProjectModel) Name() qb.Field {
	return m.Field(nameField)
}

func (m ProjectModel) CreatedAt() qb.Field {
	return m.Field(createdAtField)
}

func (m ProjectModel) UpdatedAt() qb.Field {
	return m.Field(updatedAtField)
}

func (m ProjectModel) Fields() ModelFields {
	return map[string]qb.Field{
		uidField:       m.UID(),
		refField:       m.Ref(),
		nameField:      m.Name(),
		createdAtField: m.CreatedAt(),
		updatedAtField: m.UpdatedAt(),
	}
}

func (m ProjectModel) Scan(rows *sql.Rows) (*project.Project, error) {
	var (
		uid  string
		ref  *string
		name string
	)

	if err := rows.Scan(&uid, &ref, &name); err != nil {
		return nil, err
	}

	options := []project.Option{project.UID(uid), project.Name(name)}

	if ref != nil {
		options = append(options, project.Ref(*ref))
	}

	return project.New(options...)
}
