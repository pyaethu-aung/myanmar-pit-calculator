package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/myanmar-pit-calculator/pkg/pitcalc"
)

type step int64

const (
	stepMonthlyIncome step = iota
	stepStartingMonth
	stepDependentParents
	stepDependentSpouse
	stepDependentChildrens
	stepTotalSSBContribution
	stepConfirm
)

type model struct {
	step                      step
	monthlyIncomeInput        textinput.Model
	startingMonthInput        list.Model
	dependentParentsInput     list.Model
	dependentSpouseInput      list.Model
	dependentChildrensInput   textinput.Model
	totalSSBContributionInput textinput.Model

	startingMonthList    []string
	dependentParentsList []string
	dependentSpouseList  []string

	monthlyIncome        float64
	startingMonth        int64
	dependentParents     int64
	dependentSpouse      int64
	dependentChildrens   int64
	totalSSBContribution float64

	errMessage *string
}

func main() {

	p := tea.NewProgram(initialModel())

	final, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	m := final.(model)

	output, err := pitcalc.CalculatePIT(
		pitcalc.CalculatePITInput{
			MonthlyIncome:    float64(m.monthlyIncome),
			StartingMonth:    m.startingMonth,
			DependentParents: m.dependentParents,
			DependentSpouse:  m.dependentSpouse,
			Childrens:        m.dependentChildrens,
			SSB:              float64(m.totalSSBContribution),
		})
	if err != nil {

		fmt.Printf("Error in calculating PIT: %v\n", err)
	}
	fmt.Println("=====================================")
	fmt.Printf(
		"Total Taxable Income: %s\n", currencyFormat(output.TotalTexable))
	fmt.Printf("Total Reliefs: %s\n", currencyFormat(output.TotalRelief))
	fmt.Printf("Total Personal Income Tax: %s\n", currencyFormat(output.TotalTax))
	sort.Slice(output.TaxBreakdown, func(i, j int) bool {

		return output.TaxBreakdown[i].Start < output.TaxBreakdown[j].Start
	})
	for _, v := range output.TaxBreakdown {

		if v.Limit == math.Inf(1) {

			fmt.Printf(
				"  Above from %s: %s\n",
				currencyFormat(v.Start),
				currencyFormat(v.Amount))
		} else {

			fmt.Printf(
				"  Up to %s: %s\n",
				currencyFormat(v.Limit),
				currencyFormat(v.Amount))
		}
	}
	fmt.Println("=====================================")
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
	dependentParentsList := []string{
		"No dependent parents",
		"Only one dependent parent",
		"Two dependent parents",
	}
	dependentSpouseList := []string{
		"No dependent spouse",
		"Has dependent spouse",
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

	dependentParentsItemList := make([]list.Item, len(dependentParentsList))
	for i, v := range dependentParentsList {
		dependentParentsItemList[i] = item(v)
	}
	dependentParentsInput := list.New(
		dependentParentsItemList, itemDelegate{}, 10, 8)
	dependentParentsInput.SetShowTitle(false)
	dependentParentsInput.SetShowStatusBar(false)
	dependentParentsInput.SetFilteringEnabled(false)
	dependentParentsInput.Select(0)

	dependentSpouseItemList := make([]list.Item, len(dependentSpouseList))
	for i, v := range dependentSpouseList {
		dependentSpouseItemList[i] = item(v)
	}
	dependentSpouseInput := list.New(
		dependentSpouseItemList, itemDelegate{}, 10, 8)
	dependentSpouseInput.SetShowTitle(false)
	dependentSpouseInput.SetShowStatusBar(false)
	dependentSpouseInput.SetFilteringEnabled(false)
	dependentSpouseInput.Select(0)

	dependentChildrensInput := textinput.New()
	dependentChildrensInput.Placeholder = "0"
	dependentChildrensInput.Width = 20
	dependentChildrensInput.Focus()

	totalSSBContributionInput := textinput.New()
	totalSSBContributionInput.Placeholder = "72000"
	totalSSBContributionInput.Width = 20
	totalSSBContributionInput.Focus()

	return model{

		step:                      stepMonthlyIncome,
		monthlyIncomeInput:        monthlyIncomeInput,
		startingMonthInput:        startingMonthInput,
		dependentParentsInput:     dependentParentsInput,
		dependentSpouseInput:      dependentSpouseInput,
		dependentChildrensInput:   dependentChildrensInput,
		totalSSBContributionInput: totalSSBContributionInput,

		startingMonthList:    startingMonthList,
		dependentParentsList: dependentParentsList,
		dependentSpouseList:  dependentSpouseList,
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
				m.step = stepDependentParents
			case stepDependentParents:

				m.dependentParents = int64(m.dependentParentsInput.Index())
				m.step = stepDependentSpouse
			case stepDependentSpouse:

				m.dependentSpouse = int64(m.dependentSpouseInput.Index())
				m.step = stepDependentChildrens
			case stepDependentChildrens:

				var errMessage string
				v, err := parseNumericInput(m.dependentChildrensInput.Value())
				if err != nil {

					errMessage = "❌ Invalid input, try again."
				} else if v == nil || *v < 0 {

					errMessage = "❌ Number of dependent childrens cannot be negative."
				}
				if errMessage != "" {

					m.errMessage = &errMessage
					return m, nil
				}

				m.step = stepTotalSSBContribution
				m.dependentChildrens = int64(*v)
				m.errMessage = nil
			case stepTotalSSBContribution:

				var errMessage string
				v, err := parseNumericInput(m.totalSSBContributionInput.Value())
				if err != nil {

					errMessage = "❌ Invalid input, try again."
				} else if v == nil || *v < 0 {

					errMessage = "❌ Total SSB contribution cannot be negative."
				}
				if errMessage != "" {

					m.errMessage = &errMessage
					return m, nil
				}

				m.step = stepConfirm
				m.totalSSBContribution = *v
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
	case stepStartingMonth:

		m.startingMonthInput, cmd = m.startingMonthInput.Update(msg)
	case stepDependentParents:

		m.dependentParentsInput, cmd = m.dependentParentsInput.Update(msg)
	case stepDependentSpouse:

		m.dependentSpouseInput, cmd = m.dependentSpouseInput.Update(msg)
	case stepDependentChildrens:

		m.dependentChildrensInput, cmd = m.dependentChildrensInput.Update(msg)
	case stepTotalSSBContribution:

		m.totalSSBContributionInput, cmd = m.totalSSBContributionInput.Update(msg)
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
	case stepDependentParents:
		b.WriteString("Step 3: Select number of dependent parents:\n\n")
		b.WriteString(m.dependentParentsInput.View())
	case stepDependentSpouse:
		b.WriteString("Step 4: Select dependent spouse status:\n\n")
		b.WriteString(m.dependentSpouseInput.View())
	case stepDependentChildrens:
		b.WriteString("Step 5: Enter number of children (500,000 MMK for each):\n")
		b.WriteString(m.dependentChildrensInput.View())
	case stepTotalSSBContribution:
		b.WriteString("Step 6: Enter total SSB contribution yearly (MMK):\n")
		b.WriteString(m.totalSSBContributionInput.View())
	case stepConfirm:
		b.WriteString("Final Step: Confirm your information\n\n")
		b.WriteString(
			fmt.Sprintf(
				"Monthly Income: %s\n", currencyFormat(m.monthlyIncome)))
		b.WriteString(
			fmt.Sprintf(
				"Starting Month: %s\n",
				m.startingMonthList[m.startingMonth-1]))
		b.WriteString(
			fmt.Sprintf(
				"Dependent Parents: %s\n",
				m.dependentParentsList[m.dependentParents]))
		b.WriteString(
			fmt.Sprintf(
				"Dependent Spouse: %s\n",
				m.dependentSpouseList[m.dependentSpouse]))
		b.WriteString(
			fmt.Sprintf(
				"Number of Children: %d\n",
				m.dependentChildrens))
		b.WriteString(
			fmt.Sprintf(
				"Total SSB Contribution: %s\n",
				currencyFormat(m.totalSSBContribution)))
		b.WriteString("\nPress ENTER to confirm or ESC to quit.")
	}

	if m.errMessage != nil {
		b.WriteString("\n\n" + *m.errMessage)
	}

	return b.String()
}
