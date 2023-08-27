package app

import (
	"github.com/pietjan/tm/app/command"
	"github.com/pietjan/tm/app/query"
)

type Application struct {
	Command Commands
	Query   Queries
}

type Queries struct {
	AllTasks    query.AllTasksHandler
	AllProjects query.AllProjectsHandler
	AllTimers   query.AllTimersHandler
}

type Commands struct {
	AddTask       command.AddTaskHandler
	AddProject    command.AddProjectHandler
	SetTaskStatus command.SetTaskStatusHandler
	StartTimer    command.StartTimerHandler
	StopTimer     command.StopTimerHandler
}
