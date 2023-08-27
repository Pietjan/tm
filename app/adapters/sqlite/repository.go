package sqlite

import (
	"context"
	"time"

	"github.com/pietjan/tm/app/adapters/sqlite/internal/database"
	"github.com/pietjan/tm/app/adapters/sqlite/internal/migrate"
	"github.com/pietjan/tm/app/adapters/sqlite/internal/model"
	"github.com/pietjan/tm/app/query"
	"github.com/pietjan/tm/domain/filter"
	"github.com/pietjan/tm/domain/project"
	"github.com/pietjan/tm/domain/task"
	"github.com/pietjan/tm/domain/timer"

	"git.ultraware.nl/NiseVoid/qb"
	"git.ultraware.nl/NiseVoid/qb/qc"
	"git.ultraware.nl/NiseVoid/qb/qf"
)

type Repository interface {
	project.Repository
	task.Repository
	query.AllTasksReadModel
	timer.Repository
	query.AllTimersReadModel
}

// `file:”?cache=shared&mode=memory`
func New(dsn string) Repository {
	r := &repository{
		database.New(dsn),
	}

	if err := migrate.Run(r.DB()); err != nil {
		panic(err)
	}

	return r
}

type repository struct {
	database.Database
}

// UpdateTask implements Repository.
func (r repository) UpdateTask(ctx context.Context, taskUID string, updateFn func(ctx context.Context, t *task.Task) error) error {
	task, err := r.GetTask(ctx, taskUID)
	if err != nil {
		return err
	}

	if err := updateFn(ctx, task); err != nil {
		return err
	}

	m := model.Task()
	q := m.Update().Where(qc.Eq(m.UID(), task.UID())).
		Set(m.Ref(), task.Ref()).
		Set(m.Name(), task.Name()).
		Set(m.Status(), task.Status()).
		Set(m.Project(), task.Project()).
		Set(m.UpdatedAt(), time.Now())

	_, err = r.ExecContext(ctx, q)
	return err
}

// GetTask implements Repository.
func (r repository) GetTask(ctx context.Context, taskUID string) (t *task.Task, err error) {
	m := model.Task()

	q := m.Select(m.UID(), m.Ref(), m.Name(), m.Status(), m.Project()).
		Where(
			qc.Eq(m.UID(), taskUID),
		)

	var (
		uid     string
		ref     *string
		name    string
		status  int
		project *string
	)

	if err := r.QueryRowContext(ctx, q).Scan(&uid, &ref, &name, &status, &project); err != nil {
		return nil, err
	}

	options := []task.Option{
		task.UID(uid),
		task.Name(name),
		task.Status(status),
	}

	if ref != nil {
		options = append(options, task.Ref(*ref))
	}

	if project != nil {
		options = append(options, task.Project(*project))
	}

	return task.New(
		options...,
	)
}

// AddProject implements Repository.
func (r repository) AddProject(ctx context.Context, project *project.Project) error {
	m := model.Project()
	q := m.Insert(
		m.UID(), m.Ref(), m.Name(), m.CreatedAt(),
	).Values(
		project.UID(), nilString(project.Ref()), project.Name(), time.Now().UnixNano(),
	)

	_, err := r.ExecContext(ctx, q)

	return err
}

// AddTask implements Repository.
func (r repository) AddTask(ctx context.Context, task *task.Task) error {
	m := model.Task()
	q := m.Insert(
		m.UID(), m.Ref(), m.Name(),
		m.Status(), m.Project(), m.CreatedAt(),
	).Values(
		task.UID(), task.Ref(), task.Name(),
		task.Status(), task.Project(), time.Now().UnixNano(),
	)

	_, err := r.ExecContext(ctx, q)

	return err
}

// AddTimer implements Repository.
func (r repository) AddTimer(ctx context.Context, timer *timer.Timer) error {
	m := model.Timer()
	q := m.Insert(
		m.UID(), m.Task(), m.Start(),
	).Values(
		timer.UID(), timer.Task(), timer.Start().UnixNano(),
	)

	_, err := r.ExecContext(ctx, q)

	return err
}

// AllProjects implements Repository.
func (r repository) AllProjects(ctx context.Context, filter *filter.Filter) (projects []*project.Project, err error) {
	m := model.Project()
	q := m.Select(m.UID(), m.Ref(), m.Name())

	applyFilter(m.Fields(), q, filter)

	rows, err := r.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		project, err := m.Scan(rows)
		if err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}

	return
}

// AllTasks implements Repository.
func (r repository) AllTasks(ctx context.Context, filter *filter.Filter) (tasks []query.Task, err error) {
	t, p, d := model.Task(), model.Project(), model.Timer()

	q := t.Select(
		t.UID(), t.Ref(), t.Name(), t.Status(), p.Name(), qf.Sum(qf.NewCalculatedField(qf.Coalesce(d.End(), time.Now().UnixNano()), `-`, d.Start())),
	).
		LeftJoin(p.UID(), t.Project()).
		LeftJoin(d.Task(), t.UID()).GroupBy(t.UID()).
		OrderBy(
			qb.Asc(t.Status()),
			qb.Desc(t.CreatedAt()),
		)

	applyFilter(t.Fields(), q, filter)

	rows, err := r.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task, err := t.Scan(rows)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return
}

// AllTimers implements Repository.
func (r repository) AllTimers(ctx context.Context, filter *filter.Filter) (timers []query.Timer, err error) {
	t := model.Timer()
	q := t.Select(t.UID(), t.Task(), t.Start(), t.End())

	applyFilter(t.Fields(), q, filter)

	rows, err := r.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		timer, err := t.Scan(rows)
		if err != nil {
			return nil, err
		}

		timers = append(timers, timer)
	}

	return
}

func (r repository) GetTimer(ctx context.Context, timerUID string) (*timer.Timer, error) {
	m := model.Timer()

	q := m.Select(m.Task(), m.Start(), m.End()).
		Where(
			qc.Eq(m.UID(), timerUID),
		)

	var (
		task  string
		start int64
		end   *int64
	)

	if err := r.QueryRowContext(ctx, q).Scan(&task, &start, &end); err != nil {
		return nil, err
	}

	options := []timer.Option{timer.UID(timerUID), timer.Task(task), timer.Start(time.Unix(0, start))}

	if end != nil {
		options = append(options, timer.End(time.Unix(0, *end)))
	}

	return timer.New(options...)
}

// UpdateTimer implements Repository.
func (r repository) UpdateTimer(ctx context.Context, timerUID string, updateFn timer.UpdateFn) error {
	timer, err := r.GetTimer(ctx, timerUID)
	if err != nil {
		return err
	}

	if err := updateFn(ctx, timer); err != nil {
		return err
	}

	m := model.Timer()
	q := m.Update().Where(qc.Eq(m.UID(), timerUID)).
		Set(m.Task(), timer.Task()).
		Set(m.Start(), timer.Start().UnixNano()).
		Set(m.End(), timer.End().UnixNano())

	_, err = r.ExecContext(ctx, q)
	return err
}

func applyFilter(m model.ModelFields, q *qb.SelectBuilder, f *filter.Filter) {
	if f == nil {
		return
	}

	q.Where(toQueryConditions(m, f)...)

	if f.Offset() > 0 {
		q.Offset(f.Offset())
	}

	if f.Limit() > 0 {
		q.Limit(f.Limit())
	}
}

func toQueryConditions(m model.ModelFields, f *filter.Filter) []qb.Condition {
	var conditions []qb.Condition
	if f == nil {
		return conditions
	}

	for _, condition := range f.Conditions() {
		if field, found := m[condition.Column()]; found {
			switch condition.Operator() {
			case filter.ValuesIn:
				conditions = append(conditions, qc.In(field, condition.Values()))
			case filter.Equal:
				if condition.Value() == nil {
					conditions = append(conditions, qc.IsNull(field))
					break
				}
				conditions = append(conditions, qc.Eq(field, condition.Value()))
			case filter.NotEqual:
				if condition.Value() == nil {
					conditions = append(conditions, qc.NotNull(field))
					break
				}
				conditions = append(conditions, qc.Ne(field, condition.Value()))
			case filter.GreaterOrEqual:
				conditions = append(conditions, qc.Gte(field, condition.Value()))
			case filter.GreaterThan:
				conditions = append(conditions, qc.Gt(field, condition.Value()))
			case filter.LessOrEqual:
				conditions = append(conditions, qc.Lte(field, condition.Value()))
			case filter.LessThan:
				conditions = append(conditions, qc.Lt(field, condition.Value()))
			}
		}
	}

	return conditions
}

func nilString(s string) any {
	if len(s) > 0 {
		return s
	}

	return nil
}
