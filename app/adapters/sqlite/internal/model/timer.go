package model

import (
	"database/sql"
	"time"

	"github.com/pietjan/tm/app/query"

	"git.ultraware.nl/NiseVoid/qb"
	"github.com/pietjan/model"
)

type TimerModel struct {
	model.Model
}

func (m TimerModel) UID() qb.Field {
	return m.Field(uidField)
}

func (m TimerModel) Task() qb.Field {
	return m.Field(taskField)
}

func (m TimerModel) Start() qb.Field {
	return m.Field(startField)
}

func (m TimerModel) End() qb.Field {
	return m.Field(endField)
}

func (m TimerModel) Fields() ModelFields {
	return map[string]qb.Field{
		uidField:   m.UID(),
		taskField:  m.Task(),
		startField: m.Start(),
		endField:   m.End(),
	}
}

func (m TimerModel) Scan(rows *sql.Rows) (t query.Timer, err error) {
	var (
		start int64
		end   *int64
	)

	if err := rows.Scan(&t.UID, &t.Task, &start, &end); err != nil {
		return t, err
	}

	t.Start = time.Unix(0, start)

	if end != nil {
		t.End = time.Unix(0, *end)
	}

	return
}
