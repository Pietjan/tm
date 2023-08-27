package model

import (
	"git.ultraware.nl/NiseVoid/qb"
	"github.com/pietjan/model"
)

type ModelRow interface {
	Scan(dest ...any) error
}

type ModelFields = map[string]qb.Field

const (
	timerTable   = `timer`
	projectTable = `project`
	taskTable    = `task`
)

const (
	uidField       = `uid`
	refField       = `ref`
	nameField      = `name`
	createdAtField = `created_at`
	updatedAtField = `updated_at`
	projectField   = `project`
	taskField      = `task`
	startField     = `start`
	endField       = `end`
	statusField    = `status`
)

func Project() ProjectModel {
	return ProjectModel{model.New(projectTable,
		model.Columns(uidField, refField, nameField, createdAtField, updatedAtField),
	)}
}

func Task() TaskModel {
	return TaskModel{model.New(taskTable,
		model.Columns(
			uidField, refField, nameField, statusField, projectField, createdAtField, updatedAtField,
		),
	)}
}

func Timer() TimerModel {
	return TimerModel{model.New(timerTable,
		model.Columns(
			uidField, taskField, startField, endField,
		),
	)}
}
