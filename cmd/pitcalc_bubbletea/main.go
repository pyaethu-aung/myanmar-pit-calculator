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
		"title":              "🇲🇲 Myanmar PIT Calculator",
		"lang_prompt":        "Select Language",
		"income_group":       "Income Details",
		"salary_prompt":      "Monthly Salary (MMK)",
		"bonus_prompt":       "Yearly Bonus (MMK) [Optional]",
		"reliefs_group":      "Tax Reliefs",
		"spouse_prompt":      "Dependent Spouse?",
		"spouse_desc":        "Is your spouse currently unemployed or not earning?",
		"children_prompt":    "Number of Dependent Children",
		"parents_prompt":     "Number of Dependent Parents",
		"other_group":        "Other Allowances",
		"ssb_prompt":         "Total SSB Contribution (MMK)",
		"calculating":        "Calculating...",
		"err_validation":     "❌ Invalid input, please fix errors.",
	},
	langMY: {
		"title":              "🇲🇲 မြန်မာ ဝင်ငွေခွန် တွက်စက်",
		"lang_prompt":        "ဘာသာစကား ရွေးချယ်ပါ",
		"income_group":       "ဝင်ငွေ အသေးစိတ်",
		"salary_prompt":      "လစဉ်လစာ (ကျပ်)",
		"bonus_prompt":       "နှစ်စဉ် ဆုကြေး (ကျပ်) [ရွေးချယ်ရန်]",
		"reliefs_group":      "အခွန်သက်သာခွင့်များ",
		"spouse_prompt":      "မှီခို ဇနီး/ခင်ပွန်း ရှိပါသလား?",
		"spouse_desc":        "အလုပ်လုပ်ကိုင်ခြင်းမရှိသော အိမ်ထောင်ဖက်",
		"children_prompt":    "မှီခို ကလေး အရေအတွက်",
		"parents_prompt":     "မှီခို မိဘ အရေအတွက်",
		"other_group":        "အခြားသော ခွင့်ပြုချက်များ",
		"ssb_prompt":         "လူမှုဖူလုံရေး ထည့်ဝင်ငွေ စုစုပေါင်း (ကျပ်)",
		"calculating":        "တွက်ချက်နေပါသည်...",
		"err_validation":     "❌ ထည့်သွင်းထားသော အချက်အလက်များ မှားယွင်းနေပါသည်။",
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

func t(lang langKey, id string) string {
	return trans[lang][id]
}

// --- T006: Currency Formatter ---
func currencyFormat(amount float64) string {
	return message.NewPrinter(language.English).Sprintf("%.2f MMK", amount)
}

// --- Main State Model ---
type state int

const (
	stateLang state = iota
	stateForm
	stateResult
)

type model struct {
	state        state
	langForm     *huh.Form
	taxForm      *huh.Form
	selectedLang langKey
	errMessage   string
	calcResult   *pitcalc.CalculatePITOutput
	table        table.Model

	// Form Data Pointers
	valSalary   string
	valBonus    string
	valSpouse   bool
	valChildren string
	valParents  string
	valSSB      string
}

func initialModel() model {
	m := model{
		state:        stateLang,
		selectedLang: langEN,
		valSSB:       "72000",
		valChildren:  "0",
		valParents:   "0",
	}

	m.initLangForm()
	return m
}

func (m *model) initLangForm() {
	m.langForm = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[langKey]().
				Title("Select Language / ဘာသာစကား ရွေးချယ်ပါ").
				Options(
					huh.NewOption("English", langEN),
					huh.NewOption("မြန်မာ", langMY),
				).
				Value(&m.selectedLang),
		),
	).WithTheme(huh.ThemeDracula())
	m.langForm.Init()
}

// --- T004 & US1: Define huh.Form ---
func (m *model) initTaxForm() {
	l := m.selectedLang
	m.taxForm = huh.NewForm(
		// T008: Income Group
		huh.NewGroup(
			huh.NewInput().
				Title(t(l, "salary_prompt")).
				Placeholder("500000").
				Value(&m.valSalary),
			huh.NewInput().
				Title(t(l, "bonus_prompt")).
				Placeholder("0").
				Value(&m.valBonus),
		).Title(t(l, "income_group")),
		
		// T009: Reliefs Group
		huh.NewGroup(
			huh.NewConfirm().
				Title(t(l, "spouse_prompt")).
				Description(t(l, "spouse_desc")).
				Value(&m.valSpouse),
			huh.NewInput().
				Title(t(l, "children_prompt")).
				Value(&m.valChildren),
			huh.NewInput().
				Title(t(l, "parents_prompt")).
				Value(&m.valParents),
		).Title(t(l, "reliefs_group")),
		
		// T010: Other Group
		huh.NewGroup(
			huh.NewInput().
				Title(t(l, "ssb_prompt")).
				Placeholder("72000").
				Value(&m.valSSB),
		).Title(t(l, "other_group")),
	).WithTheme(huh.ThemeDracula())

	m.taxForm.Init()
}

func (m model) Init() tea.Cmd {
	return m.langForm.Init()
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
	case stateLang:
		form, cmd := m.langForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.langForm = f
		}
		if m.langForm.State == huh.StateCompleted {
			m.state = stateForm
			m.initTaxForm()
			return m, nil
		}
		return m, cmd

	case stateForm:
		form, cmd := m.taxForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.taxForm = f
		}

		if m.taxForm.State == huh.StateCompleted {
			m.state = stateResult

			// Parse values
			salary, _ := parseNumericInput(m.valSalary)
			bonus, _ := parseNumericInput(m.valBonus)
			ssb, _ := parseNumericInput(m.valSSB)
			children, _ := parseNumericInput(m.valChildren)
			parents, _ := parseNumericInput(m.valParents)
			
			rawSalary := 0.0; if salary != nil { rawSalary = *salary }
			rawBonus := 0.0; if bonus != nil { rawBonus = *bonus }
			rawSSB := 0.0; if ssb != nil { rawSSB = *ssb }
			rawChildren := 0.0; if children != nil { rawChildren = *children }
			rawParents := 0.0; if parents != nil { rawParents = *parents }
			
			// Mock calculation for structural transition testing
			output, err := pitcalc.CalculatePIT(pitcalc.CalculatePITInput{
				MonthlyIncome:    rawSalary + (rawBonus / 12),
				StartingMonth:    4, // April default for now
				DependentParents: int64(rawParents),
				DependentSpouse:  func() int64 { if m.valSpouse { return 1 }; return 0 }(),
				Childrens:        int64(rawChildren),
				SSB:              rawSSB,
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
	// T011 [US1] Wrap form in a lipgloss styled header banner
	banner := headerStyle.Render(t(m.selectedLang, "title"))
	
	switch m.state {
	case stateLang:
		return "\n" + banner + "\n\n" + m.langForm.View()
		
	case stateForm:
		return "\n" + banner + "\n\n" + m.taxForm.View()

	case stateResult:
		if m.errMessage != "" {
			return errorStyle.Render("Error: " + m.errMessage)
		}
		if m.calcResult != nil {
			return successStyle.Render(fmt.Sprintf("\nCalculation Complete!\nTotal Tax: %s", currencyFormat(m.calcResult.TotalTax)))
		}
		return t(m.selectedLang, "calculating")
	}
	return ""
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
