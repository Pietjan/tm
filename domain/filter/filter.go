package filter

type Operator string

const (
	ValuesIn       Operator = `in`
	Equal          Operator = `eq`
	NotEqual       Operator = `ne`
	GreaterThan    Operator = `gt`
	GreaterOrEqual Operator = `ge`
	LessThan       Operator = `lt`
	LessOrEqual    Operator = `le`
)

type Condition struct {
	column   string
	operator Operator
	values   []any
}

type Filter struct {
	conditions []Condition
	limit      int
	offset     int
}

type Option = func(*Filter)

func New(options ...Option) *Filter {
	f := &Filter{}

	for _, fn := range options {
		fn(f)
	}

	return f
}

func (f Filter) Limit() int {
	return f.limit
}

func (f Filter) Offset() int {
	return f.offset
}

func (f Filter) Conditions() []Condition {
	return f.conditions
}

func (c Condition) Column() string {
	return c.column
}

func (c Condition) Operator() Operator {
	return c.operator
}

func (c Condition) Value() any {
	if len(c.values) > 0 {
		return c.values[0]
	}

	return nil
}

func (c Condition) Values() []any {
	return c.values
}

func condition(col string, op Operator, val ...any) Option {
	return func(f *Filter) {
		f.conditions = append(f.conditions, Condition{
			column:   col,
			operator: op,
			values:   val,
		})
	}
}

func Limit(limit int) Option {
	return func(f *Filter) {
		f.limit = limit
	}
}

func Offset(offset int) Option {
	return func(f *Filter) {
		f.offset = offset
	}
}

func In(col string, val ...any) Option {
	return condition(col, Equal, val...)
}

func Eq(col string, val ...any) Option {
	return condition(col, Equal, val...)
}

func Ne(col string, val ...any) Option {
	return condition(col, NotEqual, val...)
}

func Gt(col string, val ...any) Option {
	return condition(col, GreaterThan, val...)
}

func Ge(col string, val ...any) Option {
	return condition(col, GreaterOrEqual, val...)
}

func Lt(col string, val ...any) Option {
	return condition(col, LessThan, val...)
}

func Le(col string, val ...any) Option {
	return condition(col, LessOrEqual, val...)
}
