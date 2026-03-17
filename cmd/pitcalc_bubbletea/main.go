package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/myanmar-pit-calculator/pkg/pitcalc"
)

// To satisfy the compiler for now
var _ = clipboard.WriteAll
var _ = table.New

// --- T002: Lipgloss Styles & Tokens ---
var (
	themePrimary   = lipgloss.Color("#10B981") // Emerald Green
	themeSecondary = lipgloss.Color("#3B82F6") // Blue
	themeText      = lipgloss.Color("#E2E8F0") // Slate 200
	themeBorder    = lipgloss.Color("#475569") // Slate 600
	themeError     = lipgloss.Color("#EF4444") // Red

	headerStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(themePrimary).
			Padding(0, 2).
			Bold(true).
			Foreground(themeText)

	errorStyle = lipgloss.NewStyle().
			Foreground(themeError).
			Bold(true)
			
	successStyle = lipgloss.NewStyle().
			Foreground(themePrimary).
			Bold(true)
)

// --- T003: Translation Map ---
type langKey string
const (
	langEN langKey = "EN"
	langMY langKey = "MY"
)

var trans = map[langKey]map[string]string{
	langEN: {
		"title":           "🇲🇲 Myanmar PIT Calculator",
		"lang_prompt":     "Select Language",
		"income_group":    "Income",
		"salary_prompt":   "Monthly Salary (MMK)",
		"bonus_prompt":    "Yearly Bonus (MMK) [Optional]",
	},
	langMY: {
		"title":           "🇲🇲 မြန်မာ ဝင်ငွေခွန် တွက်စက်",
		"lang_prompt":     "ဘာသာစကား ရွေးချယ်ပါ",
		"income_group":    "ဝင်ငွေ",
		"salary_prompt":   "လစဉ်လစာ (ကျပ်)",
		"bonus_prompt":    "နှစ်စဉ် ဆုကြေး (ကျပ်) [ရွေးချယ်ရန်]",
	},
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

// --- T006: Currency Formatter ---
func currencyFormat(amount float64) string {
	return message.NewPrinter(language.English).Sprintf("%.2f MMK", amount)
}

// --- Main State Model ---
type state int

const (
	stateForm state = iota
	stateResult
)

type model struct {
	state          state
	form           *huh.Form
	selectedLang   langKey
	errMessage     string
	calcResult     *pitcalc.CalculatePITOutput
	table          table.Model
	
	// Form Data Pointers
	valSalary      string
	valBonus       string
}

func initialModel() model {
	m := model{
		state:        stateForm,
		selectedLang: langEN,
	}
	
	m.initForm()
	return m
}

// --- T004: Define huh.Form ---
func (m *model) initForm() {
	// Simple scaffold for now - will be expanded in Phase 3
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[langKey]().
				Title("Language / ဘာသာစကား").
				Options(
					huh.NewOption("English", langEN),
					huh.NewOption("မြန်မာ", langMY),
				).
				Value(&m.selectedLang),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Monthly Salary").
				Value(&m.valSalary),
		),
	).WithTheme(huh.ThemeDracula()) // Use a default dark theme for now
	
	m.form.Init()
}

func (m model) Init() tea.Cmd {
	return m.form.Init()
}

// --- T005: State Transition Logic ---
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}

	switch m.state {
	case stateForm:
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}
		
		if m.form.State == huh.StateCompleted {
			m.state = stateResult
			
			// Mock calculation for structural transition testing
			salary, _ := parseNumericInput(m.valSalary)
			safeSalary := 0.0
			if salary != nil {
				safeSalary = *salary
			}
			
			output, err := pitcalc.CalculatePIT(pitcalc.CalculatePITInput{
				MonthlyIncome: safeSalary,
				StartingMonth: 4, // April default for now
			})
			if err != nil {
				m.errMessage = err.Error()
			} else {
				m.calcResult = output
			}
			return m, nil
		}
		return m, cmd
		
	case stateResult:
		// Result state interactions will be handled in Phase 5
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	switch m.state {
	case stateForm:
		return "\n" + headerStyle.Render(trans[m.selectedLang]["title"]) + "\n\n" + m.form.View()
		
	case stateResult:
		if m.errMessage != "" {
			return errorStyle.Render("Error: " + m.errMessage)
		}
		if m.calcResult != nil {
			return successStyle.Render(fmt.Sprintf("\nCalculation Complete!\nTotal Tax: %s", currencyFormat(m.calcResult.TotalTax)))
		}
		return "Calculating..."
	}
	return ""
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
