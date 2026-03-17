package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
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
		"err_numeric":        "Must be a valid number",
		"err_negative":       "Cannot be negative",
		"err_parents":        "Parents must be 0, 1, or 2",
		"err_ssb":            "Maximum SSB is 360,000",
		"res_income":         "📊 Income Details",
		"res_reliefs":        "🛡️  Tax Reliefs",
		"res_total_income":   "Total Taxable Income",
		"res_total_reliefs":  "Total Reliefs",
		"res_final_tax":      "💎 Final Tax",
		"export_prompt":      "Choose Export Format",
		"success_copy":       "📋 Copied to clipboard!",
		"success_export":     "📁 Exported to PIT_Report.",
		"help_footer":        "c: Copy to clipboard • e: Export file • q: Quit",
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
		"err_numeric":        "ကိန်းဂဏန်းသာ ဖြစ်ရမည်",
		"err_negative":       "အနုတ်မရပါ",
		"err_parents":        "မိဘ ယောက်ရေ ၀, ၁, သို့မဟုတ် ၂ သာ ထည့်ပါ",
		"err_ssb":            "အများဆုံး ထည့်ဝင်ငွေ ၃၆၀,၀၀၀ ဖြစ်သည်",
		"res_income":         "📊 ဝင်ငွေ အသေးစိတ်",
		"res_reliefs":        "🛡️  အခွန်သက်သာခွင့်များ",
		"res_total_income":   "အခွန်စည်းကြပ်ရန် ဝင်ငွေ",
		"res_total_reliefs":  "သက်သာခွင့် စုစုပေါင်း",
		"res_final_tax":      "💎 ကျသင့် အခွန်ငွေ",
		"export_prompt":      "ပို့ဆောင်မည့် ပုံစံရွေးပါ",
		"success_copy":       "📋 ကူးယူပြီးပါပြီ!",
		"success_export":     "📁 PIT_Report သို့ မှတ်တမ်းတင်ပြီးပါပြီ။",
		"help_footer":        "c: ကူးယူမည် • e: ဖိုင်ထုတ်မည် • q: ထွက်မည်",
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

func validateNumeric(l langKey) func(string) error {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return nil
		}
		val, err := parseNumericInput(s)
		if err != nil {
			return errors.New(t(l, "err_numeric"))
		}
		if val != nil && *val < 0 {
			return errors.New(t(l, "err_negative"))
		}
		return nil
	}
}

func validateParents(l langKey) func(string) error {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return nil
		}
		val, err := parseNumericInput(s)
		if err != nil {
			return errors.New(t(l, "err_numeric"))
		}
		if val != nil && (*val < 0 || *val > 2) {
			return errors.New(t(l, "err_parents"))
		}
		return nil
	}
}

func validateSSB(l langKey) func(string) error {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return nil
		}
		val, err := parseNumericInput(s)
		if err != nil {
			return errors.New(t(l, "err_numeric"))
		}
		if val != nil && *val < 0 {
			return errors.New(t(l, "err_negative"))
		}
		if val != nil && *val > 360000 {
			return errors.New(t(l, "err_ssb"))
		}
		return nil
	}
}

func t(lang langKey, id string) string {
	return trans[lang][id]
}

func currencyFormat(amount float64) string {
	return message.NewPrinter(language.English).Sprintf("%.2f MMK", amount)
}

type state int

const (
	stateLang state = iota
	stateForm
	stateResult
	stateExport
)

type model struct {
	state        state
	langForm     *huh.Form
	taxForm      *huh.Form
	exportForm   *huh.Form
	selectedLang langKey
	errMessage   string
	actionAlert  string
	calcResult   *pitcalc.CalculatePITOutput
	table        table.Model

	valSalary   string
	valBonus    string
	valSpouse   bool
	valChildren string
	valParents  string
	valSSB      string
	
	valExportFormat string
}

func initialModel() model {
	m := model{
		state:        stateLang,
		selectedLang: langEN,
		valSSB:       "72000",
		valChildren:  "0",
		valParents:   "0",
		valExportFormat: "txt",
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

func (m *model) initExportForm() {
	m.exportForm = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(t(m.selectedLang, "export_prompt")).
				Options(
					huh.NewOption("TXT Document", "txt"),
					huh.NewOption("JSON Data", "json"),
					huh.NewOption("CSV Spreadsheet", "csv"),
				).
				Value(&m.valExportFormat),
		),
	).WithTheme(huh.ThemeDracula())
	m.exportForm.Init()
}

func (m *model) initTaxForm() {
	l := m.selectedLang
	m.taxForm = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(t(l, "salary_prompt")).
				Placeholder("500000").
				Validate(validateNumeric(l)).
				Value(&m.valSalary),
			huh.NewInput().
				Title(t(l, "bonus_prompt")).
				Placeholder("0").
				Validate(validateNumeric(l)).
				Value(&m.valBonus),
		).Title(t(l, "income_group")),
		
		huh.NewGroup(
			huh.NewConfirm().
				Title(t(l, "spouse_prompt")).
				Description(t(l, "spouse_desc")).
				Value(&m.valSpouse),
			huh.NewInput().
				Title(t(l, "children_prompt")).
				Validate(validateNumeric(l)).
				Value(&m.valChildren),
			huh.NewInput().
				Title(t(l, "parents_prompt")).
				Validate(validateParents(l)).
				Value(&m.valParents),
		).Title(t(l, "reliefs_group")),
		
		huh.NewGroup(
			huh.NewInput().
				Title(t(l, "ssb_prompt")).
				Placeholder("72000").
				Validate(validateSSB(l)).
				Value(&m.valSSB),
		).Title(t(l, "other_group")),
	).WithTheme(huh.ThemeDracula())

	m.taxForm.Init()
}

func buildResultTable(c *pitcalc.CalculatePITOutput) table.Model {
	columns := []table.Column{
		{Title: "From", Width: 16},
		{Title: "To", Width: 16},
		{Title: "Tax Amount", Width: 18},
	}
	
	breakdown := c.TaxBreakdown
	sort.Slice(breakdown, func(i, j int) bool {
		return breakdown[i].Start < breakdown[j].Start
	})
	
	var rows []table.Row
	for _, v := range breakdown {
		var limitStr string
		if v.Limit == math.Inf(1) {
			limitStr = "And above"
		} else {
			limitStr = currencyFormat(v.Limit)
		}
		rows = append(rows, table.Row{
			currencyFormat(v.Start),
			limitStr,
			currencyFormat(v.Amount),
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)+1),
	)
	
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(themeBorder).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(themeText).
		Bold(false)
	t.SetStyles(s)
	
	return t
}

func buildResultView(m model) string {
	if m.calcResult == nil { return "" }
	c := m.calcResult
	l := m.selectedLang

	// Income Box
	incomeText := fmt.Sprintf("%s\n%s: %s\n",
		successStyle.Render(t(l, "res_income")),
		t(l, "res_total_income"),
		currencyFormat(c.TotalTexable))
	
	incomeBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(themeBorder).
		Padding(1, 2).
		Width(38).
		Render(incomeText)

	// Reliefs Box
	reliefsText := fmt.Sprintf("%s\n%s: %s\n",
		successStyle.Render(t(l, "res_reliefs")),
		t(l, "res_total_reliefs"),
		currencyFormat(c.TotalRelief))
		
	reliefsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(themeBorder).
		Padding(1, 2).
		Width(38).
		Render(reliefsText)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, incomeBox, "  ", reliefsBox)
	
	// Final Result Box
	finalBox := lipgloss.NewStyle().
		MarginTop(1).
		Padding(1, 2).
		Background(lipgloss.Color("#1E293B")).
		Foreground(lipgloss.Color("#F8FAFC")).
		Render(fmt.Sprintf("%s: %s", t(l, "res_final_tax"), successStyle.Render(currencyFormat(c.TotalTax))))

	tableRender := "\n" + m.table.View() + "\n"
	
	footer := lipgloss.NewStyle().Foreground(themeBorder).Render(t(l, "help_footer"))
	if m.actionAlert != "" {
		footer = successStyle.Render(m.actionAlert) + "\n" + footer
	}

	return topRow + "\n" + finalBox + "\n" + tableRender + "\n" + footer
}

func generatePlainTextReport(c *pitcalc.CalculatePITOutput) string {
	var b strings.Builder
	b.WriteString("Myanmar PIT Calculator Report\n==============================\n")
	b.WriteString(fmt.Sprintf("Total Taxable Income: %s\n", currencyFormat(c.TotalTexable)))
	b.WriteString(fmt.Sprintf("Total Reliefs: %s\n", currencyFormat(c.TotalRelief)))
	b.WriteString(fmt.Sprintf("TOTAL TAX: %s\n\n", currencyFormat(c.TotalTax)))
	
	b.WriteString("Tax Breakdown:\n")
	for _, v := range c.TaxBreakdown {
		limitStr := "And above"
		if v.Limit != math.Inf(1) {
			limitStr = currencyFormat(v.Limit)
		}
		b.WriteString(fmt.Sprintf("  %s to %s -> %s\n", currencyFormat(v.Start), limitStr, currencyFormat(v.Amount)))
	}
	return b.String()
}

// --- T020: Export Writers ---
func exportToFile(format string, c *pitcalc.CalculatePITOutput) error {
	filename := "PIT_Report." + format
	
	switch format {
	case "json":
		data, err := json.MarshalIndent(c, "", "  ")
		if err != nil { return err }
		return os.WriteFile(filename, data, 0644)
	
	case "csv":
		f, err := os.Create(filename)
		if err != nil { return err }
		defer f.Close()
		
		w := csv.NewWriter(f)
		w.Write([]string{"Metric", "Value (MMK)"})
		w.Write([]string{"Total Taxable Income", fmt.Sprintf("%.2f", c.TotalTexable)})
		w.Write([]string{"Total Reliefs", fmt.Sprintf("%.2f", c.TotalRelief)})
		w.Write([]string{"Total Tax", fmt.Sprintf("%.2f", c.TotalTax)})
		w.Write([]string{"", ""})
		
		w.Write([]string{"Breakdown From", "Breakdown To", "Tax Amount"})
		for _, tb := range c.TaxBreakdown {
			limit := fmt.Sprintf("%.2f", tb.Limit)
			if tb.Limit == math.Inf(1) { limit = "And above" }
			w.Write([]string{fmt.Sprintf("%.2f", tb.Start), limit, fmt.Sprintf("%.2f", tb.Amount)})
		}
		w.Flush()
		return w.Error()
		
	default: // txt
		return os.WriteFile(filename, []byte(generatePlainTextReport(c)), 0644)
	}
}

func (m model) Init() tea.Cmd {
	return m.langForm.Init()
}

// --- T005: State Transition Logic ---
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		
		if m.state == stateResult {
			if msg.String() == "q" {
				return m, tea.Quit
			}
			if msg.String() == "c" {
				err := clipboard.WriteAll(generatePlainTextReport(m.calcResult))
				if err != nil {
					m.actionAlert = "Failed to copy"
				} else {
					m.actionAlert = t(m.selectedLang, "success_copy")
				}
				return m, nil
			}
			if msg.String() == "e" {
				m.state = stateExport
				m.initExportForm()
				return m, nil
			}
			
			// Table navigation
			var cmd tea.Cmd
			m.table, cmd = m.table.Update(msg)
			return m, cmd
		}
	
	case tea.WindowSizeMsg:
		m.table.SetWidth(msg.Width - 4)
		return m, nil
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
			
			output, err := pitcalc.CalculatePIT(pitcalc.CalculatePITInput{
				MonthlyIncome:    rawSalary + (rawBonus / 12),
				StartingMonth:    4,
				DependentParents: int64(rawParents),
				DependentSpouse:  func() int64 { if m.valSpouse { return 1 }; return 0 }(),
				Childrens:        int64(rawChildren),
				SSB:              rawSSB,
			})
			if err != nil {
				m.errMessage = err.Error()
			} else {
				m.calcResult = output
				m.table = buildResultTable(output)
			}
			return m, nil
		}
		return m, cmd

	case stateExport:
		form, cmd := m.exportForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.exportForm = f
		}
		if m.exportForm.State == huh.StateCompleted {
			m.state = stateResult
			err := exportToFile(m.valExportFormat, m.calcResult)
			if err != nil {
				m.actionAlert = "Export failed: " + err.Error()
			} else {
				m.actionAlert = t(m.selectedLang, "success_export")
			}
			return m, nil
		}
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	banner := headerStyle.Render(t(m.selectedLang, "title"))
	
	switch m.state {
	case stateLang:
		return "\n" + banner + "\n\n" + m.langForm.View()
		
	case stateForm:
		return "\n" + banner + "\n\n" + m.taxForm.View()

	case stateExport:
		return "\n" + banner + "\n\n" + m.exportForm.View()

	case stateResult:
		if m.errMessage != "" {
			return "\n" + banner + "\n\n" + errorStyle.Render("Error: " + m.errMessage)
		}
		if m.calcResult != nil {
			return "\n" + banner + "\n\n" + buildResultView(m)
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
