package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type step int64

const (
	stepMonthlyIncome step = iota
	stepStartingMonth
	stepConfirm
)

type model struct {
	step               step
	monthlyIncomeInput textinput.Model
	startingMonthInput list.Model

	startingMonthList []string

	monthlyIncome float64
	startingMonth int64

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
	fmt.Println("Monthly Income:	", m.monthlyIncome)
	fmt.Println("Starting Month:	", m.startingMonth)
	fmt.Println("=============================")
}

func initialModel() model {

	startingMonthList := []string{
		"January",
		"February",
		"March",
		"April",
		"May",
		"June",
		"July",
		"August",
		"September",
		"October",
		"November",
		"December",
	}

	monthlyIncomeInput := textinput.New()
	monthlyIncomeInput.Placeholder = "500000"
	monthlyIncomeInput.Width = 20
	monthlyIncomeInput.Focus()

	startingMonthItemList := make([]list.Item, len(startingMonthList))
	for i, v := range startingMonthList {

		startingMonthItemList[i] = item(v)
	}
	startingMonthInput := list.New(
		startingMonthItemList, itemDelegate{}, 30, 15)
	startingMonthInput.SetShowTitle(false)
	startingMonthInput.SetShowStatusBar(false)
	startingMonthInput.SetFilteringEnabled(false)
	startingMonthInput.Select(3)

	return model{

		step:               stepMonthlyIncome,
		monthlyIncomeInput: monthlyIncomeInput,
		startingMonthInput: startingMonthInput,

		startingMonthList: startingMonthList,
	}
}

type item string

func (i item) FilterValue() string { return string(i) }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(
	w io.Writer, m list.Model, index int, listItem list.Item) {

	i, ok := listItem.(item)
	if !ok {

		return
	}

	fn := lipgloss.NewStyle().PaddingLeft(4).Render
	if index == m.Index() {

		fn = func(s ...string) string {

			return lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("170")).
				Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(fmt.Sprintf("%2d. %s", index+1, string(i))))
}

func currencyFormat(amount float64) string {

	return message.NewPrinter(language.English).Sprintf("%.2f MMK", amount)
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

				m.step = stepStartingMonth
				m.monthlyIncome = *v
				m.errMessage = nil
			case stepStartingMonth:

				m.startingMonth = int64(m.startingMonthInput.Index() + 1)
				m.step = stepConfirm
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
	case stepStartingMonth:

		m.startingMonthInput, cmd = m.startingMonthInput.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	var b strings.Builder

	switch m.step {
	case stepMonthlyIncome:
		b.WriteString("Step 1: Enter monthly income (MMK):\n")
		b.WriteString(m.monthlyIncomeInput.View())
	case stepStartingMonth:
		b.WriteString("Step 2: Select starting month:\n\n")
		b.WriteString(m.startingMonthInput.View())
	case stepConfirm:
		b.WriteString("Final Step: Confirm your information\n\n")
		b.WriteString(
			fmt.Sprintf(
				"Monthly Income: %s\n", currencyFormat(m.monthlyIncome)))
		b.WriteString(
			fmt.Sprintf(
				"Starting Month: %s\n",
				m.startingMonthList[m.startingMonth-1]))
		b.WriteString("\nPress ENTER to confirm or ESC to quit.")
	}

	if m.errMessage != nil {
		b.WriteString("\n\n" + *m.errMessage)
	}

	return b.String()
}
