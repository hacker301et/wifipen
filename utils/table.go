package utils

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Row table.Row

type Model struct {
	err   error
	rows  []table.Row
	table table.Model
	Sub   chan Row
	exist bool
}

func waitForActivity(sub chan Row) tea.Cmd {

	return func() tea.Msg {
		return Row(<-sub)
	}
}
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		waitForActivity(m.Sub),
	)
}
func (m Model) GetNewTable(row Row) table.Row {
	tempRow := table.Row{}
	for _, rowValues := range row {
		tempRow = append(tempRow, rowValues)
	}

	return tempRow

}
func NewView(columns []table.Column) *Model {
	m := &Model{}
	m.Sub = make(chan Row, 0)
	tempTable := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(15),
	)
	style := table.DefaultStyles()
	style.Header = style.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("34")).
		BorderBottom(true).
		Bold(false)
	style.Selected = style.Selected.
		Foreground(lipgloss.Color("234")).
		Background(lipgloss.Color("34")).
		Bold(false)
	tempTable.SetStyles(style)
	m.table = tempTable
	return m
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch teaMsg := msg.(type) {
	case tea.KeyMsg:
		switch teaMsg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up":
			m.table, cmd = m.table.Update(msg)
			return m, cmd
		case "down":
			m.table, cmd = m.table.Update(msg)
			return m, cmd
		}
	case Row:

		row := m.GetNewTable(teaMsg)
		tempRow := []table.Row{}

		if len(m.rows) > 0 {
			for _, r := range m.rows {
				if r[0] == row[0] {
					tempRow = append(tempRow, row)
					m.exist = true
					continue
				}
				tempRow = append(tempRow, r)
			}
		}
		if !m.exist || len(m.rows) == 0 {
			tempRow = append(m.rows, row)
		}
		m.exist = false
		m.rows = tempRow
		m.table.SetRows(m.rows)
		return m, waitForActivity(m.Sub)

	}

	return m, cmd
}

func (m Model) View() string {
	view := baseStyle.Render(" ⬆  To Move Up  \n ⬇  To Move down \n"+m.table.View()) + "\n"
	return view

}
