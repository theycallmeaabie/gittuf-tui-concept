// Copyright The gittuf Authors
// SPDX-License-Identifier: Apache-2.0

package tui

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gittuf/gittuf/experimental/gittuf"
	"github.com/gittuf/gittuf/internal/tuf"
	"github.com/secure-systems-lab/go-securesystemslib/dsse"
)

type screen int

const (
	screenLoading             screen = iota // Loading screen shown on startup
	screenChoice                            // Initial menu
	screenPolicy                            // Menu for Policy operations
	screenPolicyRules                       // Rule management screen
	screenPolicyAddRule                     // Form: add a new policy rule
	screenPolicyEditRule                    // Form: edit selected rule (prefilled)
	screenTrust                             // Menu for Trust operations
	screenTrustGlobalRules                  // Global rule management screen
	screenTrustAddGlobalRule                // Form: add a new global rule
	screenTrustEditGlobalRule               // Form: edit selected global rule (prefilled)
	screenMockTrust
	screenMockPolicy
	screenKeyDetail
	screenMockVerify
)

// item is what we hand to bubbles/list. The list package needs these three
// methods; nothing else is going on here.
type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	ctx              context.Context
	screen           screen
	spinner          spinner.Model
	choiceList       list.Model
	policyScreenList list.Model
	trustScreenList  list.Model
	rules            []rule
	ruleList         list.Model
	globalRules      []globalRule
	globalRuleList   list.Model
	inputs           []textinput.Model
	focusIndex       int
	cursorMode       cursor.Mode
	repo             *gittuf.Repository
	signer           dsse.SignerVerifier
	policyName       string
	options          *options
	footer           string
	errorMsg         string
	readOnly         bool
	confirmDelete    bool
	deleteTarget     string

	trustKeyIdx int
	verifyIdx   int

	width  int
	height int
}

// initDoneMsg is sent back from loadRepoCmd once the repo, signer and rules
// have finished loading (or failed). The Update loop swaps the loading screen
// for the real menu when this lands.
type initDoneMsg struct {
	repo        *gittuf.Repository
	signer      dsse.SignerVerifier
	rules       []rule
	globalRules []globalRule
	readOnly    bool
	footer      string
	err         error
}

type inputField struct {
	placeholder string
	prompt      string
}

func newDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = selectedItemStyle
	d.Styles.SelectedDesc = selectedItemStyle
	d.Styles.NormalTitle = itemStyle
	d.Styles.NormalDesc = itemStyle
	return d
}

func newMenuList(title string, items []list.Item, delegate list.DefaultDelegate) list.Model {
	l := list.New(items, delegate, 0, 0)
	l.Title = title
	l.Styles.Title = titleStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	return l
}

// initInputs builds the text inputs for a form. First one gets focus, the
// rest start blurred so tab cycling has somewhere to go.
func initInputs(fields []inputField) []textinput.Model {
	inputs := make([]textinput.Model, len(fields))
	for i, f := range fields {
		t := textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 64
		t.Placeholder = f.placeholder
		t.Prompt = f.prompt
		if i == 0 {
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		} else {
			t.Blur()
			t.PromptStyle = blurredStyle
			t.TextStyle = blurredStyle
		}
		inputs[i] = t
	}
	return inputs
}

// initialModel just spins up the lists and a spinner so we can paint the
// loading screen immediately. Anything that touches disk or git happens later
// in loadRepoCmd — otherwise the first frame is noticeably slow on big repos.
func initialModel(ctx context.Context, o *options) model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	delegate := newDelegate()

	m := model{
		ctx:        ctx,
		screen:     screenLoading,
		spinner:    s,
		cursorMode: cursor.CursorBlink,
		policyName: o.policyName,
		options:    o,

		choiceList: newMenuList("gittuf TUI", []list.Item{
			item{title: "Policy", desc: "View and manage gittuf Policy"},
			item{title: "Trust", desc: "View and manage gittuf Root of Trust"},
			item{title: "Verify Ref", desc: "Check whether a Git ref complies with policy"},
		}, delegate),
		policyScreenList: newMenuList("gittuf Policy Operations", []list.Item{
			item{title: "View Rules", desc: "View and manage policy rules"},
			item{title: "View Branch Protection", desc: "Demo: visual summary of protected branches"},
		}, delegate),
		trustScreenList: newMenuList("gittuf Trust Operations", []list.Item{
			item{title: "View Global Rules", desc: "View and manage global rules"},
			item{title: "View Trusted Keys", desc: "Demo: root of trust + trusted signers"},
		}, delegate),
		ruleList:       newMenuList("Policy Rules", []list.Item{}, delegate),
		globalRuleList: newMenuList("Global Rules", []list.Item{}, delegate),
	}

	return m
}

// loadRepoCmd does the slow startup work off the UI thread: opening the repo,
// loading the signing key, pulling the current rules. Result comes back as an
// initDoneMsg.
func loadRepoCmd(ctx context.Context, o *options) tea.Cmd {
	return func() tea.Msg {
		repo, err := gittuf.LoadRepository(".")
		if err != nil {
			// No gittuf repo in cwd — fall back to demo mode so reviewers can
			// still poke around the mock screens without setting one up.
			return initDoneMsg{
				readOnly: true,
				footer:   "Demo mode: no gittuf repository detected.",
			}
		}

		readOnly := o.readOnly
		var signer dsse.SignerVerifier
		var footer string

		if !readOnly {
			signer, err = gittuf.LoadSigner(repo, o.p.SigningKey)
			if err != nil {
				if !errors.Is(err, gittuf.ErrSigningKeyNotSpecified) {
					return initDoneMsg{err: fmt.Errorf("failed to load signing key from Git config: %w", err)}
				}
				readOnly = true
				footer = "No signing key found in Git config, running in read-only mode."
			}
		}

		return initDoneMsg{
			repo:        repo,
			signer:      signer,
			rules:       getCurrRules(ctx, o),
			globalRules: getGlobalRules(ctx, o),
			readOnly:    readOnly,
			footer:      footer,
		}
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick, loadRepoCmd(m.ctx, m.options))
}

func (m *model) initRuleInputs() {
	m.inputs = initInputs([]inputField{
		{"Enter Rule Name Here", "Rule Name:"},
		{"Enter Rule Pattern Here", " Rule Pattern:"},
		{"Enter Principal IDs Here (comma-separated)", "Authorized Principals:"},
		{"Enter Threshold", "Threshold:"},
	})
	m.focusIndex = 0
}

func (m *model) initRuleInputsPrefilled(r rule) {
	m.initRuleInputs()
	m.inputs[0].SetValue(r.name)
	m.inputs[1].SetValue(r.pattern)
	m.inputs[2].SetValue(r.key)
	m.inputs[3].SetValue(fmt.Sprintf("%d", r.threshold))
}

func (m *model) initGlobalRuleInputs() {
	m.inputs = initInputs([]inputField{
		{"Enter Global Rule Name Here", "Rule Name:"},
		{"Enter Global Rule Type (threshold|block-force-pushes)", "Type:"},
		{"Enter Namespaces (comma-separated)", "Namespaces:"},
		{"Enter Threshold (if threshold type)", "Threshold:"},
	})
	m.focusIndex = 0
}

func (m *model) initGlobalRuleInputsPrefilled(gr globalRule) {
	m.initGlobalRuleInputs()
	m.inputs[0].SetValue(gr.ruleName)
	m.inputs[1].SetValue(gr.ruleType)
	m.inputs[2].SetValue(strings.Join(gr.rulePatterns, ", "))
	if gr.ruleType == tuf.GlobalRuleThresholdType {
		m.inputs[3].SetValue(fmt.Sprintf("%d", gr.threshold))
	}
}

// refreshRules / refreshGlobalRules re-read from the repo after a write.
// We always go back to disk rather than mutating the in-memory slice so the
// list reflects whatever the policy actually is, including anything signing
// or validation rejected.
func (m *model) refreshRules() {
	m.rules = getCurrRules(m.ctx, m.options)
	m.updateRuleList()
}

func (m *model) refreshGlobalRules() {
	m.globalRules = getGlobalRules(m.ctx, m.options)
	m.updateGlobalRuleList()
}

func (m *model) updateRuleList() {
	items := make([]list.Item, len(m.rules))
	for i, rule := range m.rules {
		items[i] = item{title: rule.name, desc: fmt.Sprintf("Pattern: %s, Key: %s, Threshold: %d", rule.pattern, rule.key, rule.threshold)}
	}
	m.ruleList.SetItems(items)
}

func (m *model) updateGlobalRuleList() {
	items := make([]list.Item, len(m.globalRules))
	for i, gr := range m.globalRules {
		desc := fmt.Sprintf(
			"Type: %s\nNamespaces: %s",
			gr.ruleType,
			strings.Join(gr.rulePatterns, ", "),
		)
		if gr.ruleType == tuf.GlobalRuleThresholdType {
			desc += fmt.Sprintf("\nThreshold: %d", gr.threshold)
		}
		items[i] = item{title: gr.ruleName, desc: desc}
	}
	m.globalRuleList.SetItems(items)
}
