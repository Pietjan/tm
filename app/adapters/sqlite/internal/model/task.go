package model

import (
	"database/sql"
	"time"

	"github.com/pietjan/tm/app/query"

	"git.ultraware.nl/NiseVoid/qb"
	"github.com/pietjan/model"
)

type TaskModel struct {
	model.Model
}

func (m TaskModel) UID() qb.Field {
	return m.Field(uidField)
}

func (m TaskModel) Ref() qb.Field {
	return m.Field(refField)
}

func (m TaskModel) Name() qb.Field {
	return m.Field(nameField)
}

func (m TaskModel) Status() qb.Field {
	return m.Field(statusField)
}

func (m TaskModel) Project() qb.Field {
	return m.Field(projectField)
}

func (m TaskModel) CreatedAt() qb.Field {
	return m.Field(createdAtField)
}

func (m TaskModel) UpdatedAt() qb.Field {
	return m.Field(updatedAtField)
}

func (m TaskModel) Fields() ModelFields {
	return map[string]qb.Field{
		uidField:       m.UID(),
		refField:       m.Ref(),
		nameField:      m.Name(),
		projectField:   m.Project(),
		createdAtField: m.CreatedAt(),
		updatedAtField: m.UpdatedAt(),
	}
}

func (m TaskModel) Scan(rows *sql.Rows) (t query.Task, err error) {
	var (
		ref      *string
		project  *string
		duration *time.Duration
	)

	if err := rows.Scan(&t.UID, &ref, &t.Name, &t.Status, &project, &duration); err != nil {
		return t, err
	}

	if ref != nil {
		t.Ref = *ref
	}

	if project != nil {
		t.Project = *project
	}

	if duration != nil {
		t.Duration = *duration
	}

	return
}
