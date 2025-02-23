package util

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type SelectionModel struct {
	initialMessage string
	Choices        []string
	Selected       int
}

func NewModel(initialMessage string, choices []string) SelectionModel {
	return SelectionModel{
		initialMessage: initialMessage,
		Choices:        choices,
		Selected:       0,
	}
}

func (m SelectionModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (model SelectionModel) View() string {
	content := fmt.Sprintf("%s\n\n", model.initialMessage)

	for i, choice := range model.Choices {
		cursor := " "
		if model.Selected == i {
			cursor = ">"
		}
		content += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	content += "\nPress q to quit.\n"

	return content
}

func (model SelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			model.Selected = len(model.Choices) - 1
			return model, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if model.Selected > 0 {
				model.Selected--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if model.Selected < len(model.Choices)-1 {
				model.Selected++
			}

		case "enter", " ":
			return model, tea.Quit
		}
	}

	return model, nil
}
