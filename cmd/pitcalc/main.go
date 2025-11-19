package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type step int64

const (
	stepMonthlyIncome step = iota
	stepConfirm
)

type model struct {
	step               step
	monthlyIncomeInput textinput.Model

	monthlyIncome float64

	errMessage *string
}

func main() {

	p := tea.NewProgram(initialModel())

	final, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	m := final.(model)

	fmt.Println("========== RESULTS ==========")
	fmt.Println("Monthly Income:     ", m.monthlyIncome)
	fmt.Println("=============================")
}

func initialModel() model {

	monthlyIncomeInput := textinput.New()
	monthlyIncomeInput.Placeholder = "500000"
	monthlyIncomeInput.Width = 20
	monthlyIncomeInput.Focus()

	return model{

		step:               stepMonthlyIncome,
		monthlyIncomeInput: monthlyIncomeInput,
	}
}

func parseNumericInput(input string) (*float64, error) {

	clean := strings.ReplaceAll(strings.TrimSpace(input), ",", "")

	value := 0.0
	if clean == "" {
		return &value, nil
	}

	value, err := strconv.ParseFloat(clean, 64)
	if err != nil {

		return nil, errors.New("invalid numeric format")
	}

	return &value, nil
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.Type {

		case tea.KeyEnter:

			switch m.step {

			case stepMonthlyIncome:

				var errMessage string
				v, err := parseNumericInput(m.monthlyIncomeInput.Value())
				if err != nil {

					errMessage = "❌ Invalid input, try again."
				} else if v == nil || *v <= 0 {

					errMessage = "❌ Monthly income must be greater than 0."
				}
				if errMessage != "" {

					m.errMessage = &errMessage
					return m, nil
				}

				m.step = stepConfirm
				m.monthlyIncome = *v
				m.errMessage = nil
			case stepConfirm:

				return m, tea.Quit
			}
		case tea.KeyCtrlC, tea.KeyEsc:

			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	switch m.step {

	case stepMonthlyIncome:

		m.monthlyIncomeInput, cmd = m.monthlyIncomeInput.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	var b strings.Builder

	switch m.step {
	case stepMonthlyIncome:
		b.WriteString("Step 1: Enter monthly income (MMK):\n")
		b.WriteString(m.monthlyIncomeInput.View())
	case stepConfirm:
		b.WriteString("Final Step: Confirm your information\n\n")
		b.WriteString(fmt.Sprintf("Monthly Income: %.2f MMK\n", m.monthlyIncome))
		b.WriteString("\nPress ENTER to confirm or ESC to quit.")
	}

	if m.errMessage != nil {
		b.WriteString("\n\n" + *m.errMessage)
	}

	return b.String()
}
