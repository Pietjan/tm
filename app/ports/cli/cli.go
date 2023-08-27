package cli

import (
	"fmt"
	"strings"

	"github.com/pietjan/tm/app"
	"github.com/pietjan/tm/app/adapters/sqlite"
	"github.com/pietjan/tm/app/command"
	"github.com/pietjan/tm/app/ports/tui"
	"github.com/pietjan/tm/app/query"
	"github.com/pietjan/tm/domain/filter"

	"github.com/segmentio/ksuid"
	"github.com/urfave/cli/v2"
)

type App struct {
	*cli.App
	app app.Application
}

func New(repo sqlite.Repository) App {
	a := App{
		App: &cli.App{
			Name:  "task-cli",
			Usage: "manage tasks from the cli",
		},
		app: app.Application{
			Command: app.Commands{
				AddTask:       command.NewAddTaskHandler(repo),
				AddProject:    command.NewAddProjectHandler(repo),
				StartTimer:    command.NewStartTimerHandler(repo),
				StopTimer:     command.NewStopTimerHandler(repo),
				SetTaskStatus: command.NewSetTaskStatusHandler(repo),
			},
			Query: app.Queries{
				AllTasks:    query.NewAllTasksHandler(repo),
				AllProjects: query.NewAllProjectsHandler(repo),
				AllTimers:   query.NewAllTimersHandler(repo),
			},
		},
	}

	a.Action = addTask(a.app)
	a.Flags = append(a.Flags, &cli.StringFlag{
		Name:    `project`,
		Aliases: []string{`p`},
	})

	a.Commands = []*cli.Command{
		{Name: `start`, Action: startTask(a.app)},
		{Name: `stop`, Action: stopTask(a.app)},
		{Name: `list`, Aliases: []string{`ls`}, Action: listTasks(a.app)},
	}

	return a
}

func addTask(app app.Application) cli.ActionFunc {
	return func(c *cli.Context) error {
		var projectUID string

		if c.NumFlags() > 0 && len(c.String(`project`)) > 0 {
			projectName := c.String(`project`)
			projects, err := app.Query.AllProjects.Handle(c.Context, query.AllProjects{
				Filter: filter.New(
					filter.Eq(`name`, projectName),
				),
			})
			if err != nil {
				return err
			}

			if len(projects) == 1 {
				projectUID = projects[0].UID()
			}

			if len(projects) == 0 {
				projectUID = ksuid.New().String()
				app.Command.AddProject.Handle(c.Context, command.AddProject{
					ProjectUID: projectUID,
					Name:       projectName,
				})
			}
		}

		return app.Command.AddTask.Handle(c.Context, command.AddTask{
			TaskUID: ksuid.New().String(),
			Project: projectUID,
			Name:    strings.Join(c.Args().Slice(), ` `),
		})
	}
}

func startTask(app app.Application) cli.ActionFunc {
	return func(c *cli.Context) error {
		taskName := strings.Join(c.Args().Slice(), ` `)

		tasks, err := app.Query.AllTasks.Handle(c.Context, query.AllTasks{
			Filter: filter.New(
				filter.Eq(`name`, taskName),
			),
		})
		if err != nil {
			return err
		}

		if len(tasks) != 1 {
			return fmt.Errorf(`task %q not found`, taskName)
		}

		err = stopTask(app)(c)
		if err != nil {
			return err
		}

		err = app.Command.SetTaskStatus.Handle(c.Context, command.SetTaskStatus{
			TaskUID: tasks[0].UID,
			Status:  1, // in progress
		})
		if err != nil {
			return err
		}

		return app.Command.StartTimer.Handle(c.Context, command.StartTimer{
			TimerUID: ksuid.New().String(),
			TaskUID:  tasks[0].UID,
		})
	}
}

func stopTask(app app.Application) cli.ActionFunc {
	return func(c *cli.Context) error {
		timers, err := app.Query.AllTimers.Handle(c.Context, query.AllTimers{
			Filter: filter.New(filter.Eq(`end`, nil)),
		})
		if err != nil {
			return err
		}

		for _, timer := range timers {
			err := app.Command.StopTimer.Handle(c.Context, command.StopTimer{
				TimerUID: timer.UID,
			})
			if err != nil {
				return err
			}

			err = app.Command.SetTaskStatus.Handle(c.Context, command.SetTaskStatus{
				TaskUID: timer.Task,
				Status:  2, // done
			})
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func listTasks(app app.Application) cli.ActionFunc {
	return func(c *cli.Context) error {
		tasks, err := app.Query.AllTasks.Handle(c.Context, query.AllTasks{
			Filter: nil,
		})
		if err != nil {
			return err
		}

		fmt.Println(tui.TasksTable(tasks).View())

		return nil
	}
}
