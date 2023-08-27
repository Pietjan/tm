package tui

import (
	"log"
	"os"
	"time"

	"github.com/pietjan/tm/app/query"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type status int

const (
	todo status = iota
	inProgress
	done
)

func (s status) String() string {
	return []string{`todo`, `in progress`, `done`}[s]
}

func calculateWidth(min, width int) int {
	p := width / 10
	switch min {
	case XS:
		if p < XS {
			return XS
		}
		return p / 2

	case SM:
		if p < SM {
			return SM
		}
		return p / 2
	case MD:
		if p < MD {
			return MD
		}
		return p * 2
	case LG:
		if p < LG {
			return LG
		}
		return p * 3
	default:
		return p
	}
}

const (
	XS int = 1
	SM int = 3
	MD int = 5
	LG int = 10
)

func TasksTable(tasks []query.Task) table.Model {
	// get term size
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// we don't really want to fail it...
		log.Println(`unable to calculate height and width of terminal`)
	}

	columns := []table.Column{
		{Title: `UID`, Width: calculateWidth(MD, w)},
		{Title: `Name`, Width: calculateWidth(LG, w)},
		{Title: `Project`, Width: calculateWidth(MD, w)},
		{Title: `Status`, Width: calculateWidth(SM, w)},
		{Title: `Duration`, Width: calculateWidth(MD, w)},
	}
	var rows []table.Row
	for _, task := range tasks {
		rows = append(rows, table.Row{
			task.UID,
			task.Name,
			task.Project,
			status(task.Status).String(),
			formatDuration(task.Duration),
		})
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(len(tasks)),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(`240`)).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(nil).
		Background(nil).
		Bold(false)

	t.SetStyles(s)
	return t
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return ``
	}

	return d.Round(time.Second).String()
}
