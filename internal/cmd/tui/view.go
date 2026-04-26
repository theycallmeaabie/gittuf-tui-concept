// Copyright The gittuf Authors
// SPDX-License-Identifier: Apache-2.0

package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

const (
	colorRegularText = "#FFFFFF"
	colorFocus       = "#007AFF"
	colorBlur        = "#A0A0A0"
	colorFooter      = "#11ff00"
	colorSubtext     = "#555555"
	colorErrorMsg    = "#FF0000"

	colorValid       = "#10B981"
	colorExpired     = "#EF4444"
	colorWarning     = "#F59E0B"
	colorAccent      = "#38BDF8"
	colorBorder      = "#374151"
	colorBarBg       = "#1F2937"
	colorBadgeEditBg = "#1D4ED8"
	colorBadgeRoBg   = "#7F1D1D"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorRegularText)).
			Padding(0, 2).
			MarginTop(1).
			Bold(true)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(4).
			Foreground(lipgloss.Color(colorRegularText))

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(4).
				Foreground(lipgloss.Color(colorRegularText)).
				Background(lipgloss.Color(colorFocus))

	focusedStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	blurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorBlur))

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorRegularText))

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colorBorder)).
			Padding(1, 2).
			MarginBottom(1)

	rootCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colorAccent)).
			Padding(1, 2).
			MarginBottom(1)

	keyValueLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorBlur)).
				Width(14)

	fingerprintStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorAccent))

	validBadgeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorValid)).
			Bold(true)

	expiredBadgeStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorExpired)).
				Bold(true)

	warningBadgeStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorWarning)).
				Bold(true)

	branchPatternStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorAccent)).
				Bold(true)

	statusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(colorBarBg)).
			Foreground(lipgloss.Color(colorRegularText)).
			Padding(0, 1)

	statusBadgeEditStyle = lipgloss.NewStyle().
				Background(lipgloss.Color(colorBadgeEditBg)).
				Foreground(lipgloss.Color(colorRegularText)).
				Bold(true).
				Padding(0, 1)

	statusBadgeRoStyle = lipgloss.NewStyle().
				Background(lipgloss.Color(colorBadgeRoBg)).
				Foreground(lipgloss.Color(colorRegularText)).
				Bold(true).
				Padding(0, 1)
)

// renderWithMargin wraps content in the standard margin used by all screens.
func renderWithMargin(content string) string {
	return lipgloss.NewStyle().Margin(1, 2).Render(content)
}

// renderFooter returns the footer text styled in the footer color.
func renderFooter(text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorFooter)).Render(text)
}

// renderErrorMsg returns error messages styled in the error color.
func renderErrorMsg(text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorErrorMsg)).Render(text)
}

// renderFormScreen renders a form screen with a title, input fields, help text, and footer.
func (m model) renderFormScreen(title string) string {
	var b strings.Builder
	b.WriteString(titleStyle.Render(title) + "\n\n")
	for _, input := range m.inputs {
		b.WriteString(input.View() + "\n")
	}
	b.WriteString("\n" + "Press Tab to advance, Enter to advance/submit, and Esc to go back." + "\n")
	b.WriteString(renderFooter(m.footer))
	b.WriteString(renderErrorMsg(m.errorMsg))
	return renderWithMargin(b.String())
}

// renderListScreen renders a list with help text and footer.
func (m model) renderListScreen(l list.Model, helpText string, emptyMsg string, isEmpty bool) string {
	listView := l.View()
	if isEmpty {
		emptyMsgStyled := lipgloss.NewStyle().Foreground(lipgloss.Color(colorSubtext)).Render(emptyMsg)
		listView = l.Title + "\n\n" + emptyMsgStyled
	}

	return renderWithMargin(
		listView + "\n\n" +
			renderFooter(m.footer) +
			renderErrorMsg(m.errorMsg) +
			"\n" + helpText,
	)
}

// screenPolicyRulesHelp returns the help bar for the policy rules view screen.
func screenPolicyRulesHelp(readOnly bool) string {
	if readOnly {
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colorBlur)).Render(
			"esc: back  q: quit",
		)
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorBlur)).Render(
		"a: add  e: edit  d: delete  k: move-up  j: move-down  esc: back  q: quit",
	)
}

// screenTrustGlobalRulesHelp returns the help bar for the global rules view screen.
func screenTrustGlobalRulesHelp(readOnly bool) string {
	if readOnly {
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colorBlur)).Render(
			"esc: back  q: quit",
		)
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorBlur)).Render(
		"a: add  e: edit  d: delete  esc: back  q: quit",
	)
}

// renderDeleteOverlay renders the delete confirmation prompt.
func renderDeleteOverlay(target string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true).
		Render(fmt.Sprintf("Delete rule %q? [y/n]", target))
}

func breadcrumb(s screen, keyIdx int) string {
	switch s {
	case screenMockTrust:
		return "Home > Trust"
	case screenMockPolicy:
		return "Home > Policy"
	case screenKeyDetail:
		if keyIdx < len(mockTrustedKeys) {
			return "Home > Trust > " + mockTrustedKeys[keyIdx].name
		}
		return "Home > Trust > Key"
	case screenMockVerify:
		return "Home > Verify Ref"
	default:
		return "Home"
	}
}

func (m model) renderStatusBar() string {
	crumb := breadcrumb(m.screen, m.trustKeyIdx)

	var badge string
	if m.readOnly {
		badge = statusBadgeRoStyle.Render(" READ-ONLY ")
	} else {
		badge = statusBadgeEditStyle.Render(" EDIT ")
	}

	width := m.width
	if width <= 0 {
		width = 80
	}

	badgeW := lipgloss.Width(badge)
	leftW := max(width-badgeW, 0)

	left := statusBarStyle.Width(leftW).Render(crumb)
	return lipgloss.JoinHorizontal(lipgloss.Top, left, badge)
}

func renderRootKeyCard() string {
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorAccent)).
		Bold(true).
		Render("◆ " + mockRootKey.name)

	body := title + "\n\n" +
		keyValueLabelStyle.Render("Type:") + mockRootKey.keyType + "\n" +
		keyValueLabelStyle.Render("Fingerprint:") + fingerprintStyle.Render(mockRootKey.fingerprint) + "\n" +
		keyValueLabelStyle.Render("Added:") + mockRootKey.addedDate + "\n" +
		keyValueLabelStyle.Render("Expires:") + validBadgeStyle.Render(mockRootKey.expiresDate)

	return rootCardStyle.Render(body)
}

func renderTrustedKeyRow(k mockKey, selected bool) string {
	var statusBadge string
	if k.expired {
		statusBadge = expiredBadgeStyle.Render("● EXPIRED")
	} else {
		statusBadge = validBadgeStyle.Render("● VALID")
	}

	nameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorRegularText)).
		Bold(true)

	prefix := "  "
	if selected {
		prefix = lipgloss.NewStyle().Foreground(lipgloss.Color(colorFocus)).Bold(true).Render("▶ ")
		nameStyle = nameStyle.Foreground(lipgloss.Color(colorFocus))
	}

	line1 := prefix + nameStyle.Render(k.name) + "  " + statusBadge
	line2 := "    " + keyValueLabelStyle.Render("Type:") + k.keyType +
		"  " + keyValueLabelStyle.Render("Expires:") + k.expiresDate
	line3 := "    " + fingerprintStyle.Render(k.fingerprint)

	return line1 + "\n" + line2 + "\n" + line3
}

func (m model) renderMockTrust() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Root of Trust") + "\n\n")
	b.WriteString(renderRootKeyCard() + "\n")

	section := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorRegularText)).
		Render(fmt.Sprintf("Trusted Keys  (%d)", len(mockTrustedKeys)))
	b.WriteString(section + "\n\n")

	for i, k := range mockTrustedKeys {
		b.WriteString(renderTrustedKeyRow(k, i == m.trustKeyIdx) + "\n\n")
	}

	help := lipgloss.NewStyle().Foreground(lipgloss.Color(colorBlur)).Render(
		"↑/↓: navigate    enter: view details    esc: back    q: quit",
	)
	b.WriteString(help)

	content := renderWithMargin(b.String())
	return content + "\n" + m.renderStatusBar()
}

func renderBranchCard(br mockBranch) string {
	pattern := branchPatternStyle.Render(br.pattern)
	threshold := fmt.Sprintf("%d-of-%d signatures required", br.threshold, br.totalApprovers)

	var fpStatus string
	if br.forcePushBlocked {
		fpStatus = validBadgeStyle.Render("● Force push blocked")
	} else {
		fpStatus = warningBadgeStyle.Render("● Force push allowed")
	}

	approvers := strings.Join(br.requiredApprovers, ", ")

	body := pattern + "\n\n" +
		keyValueLabelStyle.Render("Threshold:") + threshold + "\n" +
		keyValueLabelStyle.Render("Approvers:") + approvers + "\n" +
		keyValueLabelStyle.Render("Status:") + fpStatus

	return cardStyle.Render(body)
}

func (m model) renderMockPolicy() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Branch Protection Policy") + "\n\n")

	summary := lipgloss.NewStyle().Foreground(lipgloss.Color(colorBlur)).Render(
		fmt.Sprintf("%d branches under protection.", len(mockProtectedBranches)),
	)
	b.WriteString(summary + "\n\n")

	for _, br := range mockProtectedBranches {
		b.WriteString(renderBranchCard(br) + "\n")
	}

	help := lipgloss.NewStyle().Foreground(lipgloss.Color(colorBlur)).Render(
		"esc: back    q: quit",
	)
	b.WriteString(help)

	content := renderWithMargin(b.String())
	return content + "\n" + m.renderStatusBar()
}

func renderVerifyRefRow(c mockVerifyCase, selected bool) string {
	var statusBadge string
	if c.passed {
		statusBadge = validBadgeStyle.Render("● PASS")
	} else {
		statusBadge = expiredBadgeStyle.Render("● FAIL")
	}

	refStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorRegularText)).
		Bold(true)

	prefix := "  "
	if selected {
		prefix = lipgloss.NewStyle().Foreground(lipgloss.Color(colorFocus)).Bold(true).Render("▶ ")
		refStyle = refStyle.Foreground(lipgloss.Color(colorFocus))
	}

	return prefix + refStyle.Render(c.refName) + "  " + statusBadge
}

func renderVerifyResultCard(c mockVerifyCase) string {
	var headline string
	if c.passed {
		headline = validBadgeStyle.Render("✓ POLICY COMPLIANT")
	} else {
		headline = expiredBadgeStyle.Render("✗ POLICY VIOLATION")
	}

	signers := strings.Join(c.signers, ", ")
	if signers == "" {
		signers = "(none)"
	}

	body := headline + "\n\n" +
		keyValueLabelStyle.Render("Ref:") + branchPatternStyle.Render(c.refName) + "\n" +
		keyValueLabelStyle.Render("Commit:") + fingerprintStyle.Render(c.commit) + "\n" +
		keyValueLabelStyle.Render("Signers:") + signers + "\n" +
		keyValueLabelStyle.Render("Threshold:") + fmt.Sprintf("%d-of-%d required", len(c.signers), c.threshold) + "\n\n"

	checksHeader := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorRegularText)).Render("Checks")
	body += checksHeader + "\n"

	for _, ch := range c.checks {
		var icon string
		if ch.passed {
			icon = validBadgeStyle.Render("✓")
		} else {
			icon = expiredBadgeStyle.Render("✗")
		}
		detail := lipgloss.NewStyle().Foreground(lipgloss.Color(colorBlur)).Render(ch.detail)
		body += "  " + icon + "  " + ch.name + "  " + detail + "\n"
	}

	if !c.passed && c.failureReason != "" {
		body += "\n" + warningBadgeStyle.Render("!  "+c.failureReason)
	}

	if c.passed {
		return rootCardStyle.Render(body)
	}
	return cardStyle.Render(body)
}

func (m model) renderMockVerify() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Verify Ref") + "\n\n")

	intro := lipgloss.NewStyle().Foreground(lipgloss.Color(colorBlur)).Render(
		"Pick a ref to check it against the current policy.",
	)
	b.WriteString(intro + "\n\n")

	for i, c := range mockVerifyCases {
		b.WriteString(renderVerifyRefRow(c, i == m.verifyIdx) + "\n")
	}

	b.WriteString("\n")
	if m.verifyIdx < len(mockVerifyCases) {
		b.WriteString(renderVerifyResultCard(mockVerifyCases[m.verifyIdx]) + "\n")
	}

	help := lipgloss.NewStyle().Foreground(lipgloss.Color(colorBlur)).Render(
		"↑/↓: select ref    esc: back    q: quit",
	)
	b.WriteString(help)

	return renderWithMargin(b.String()) + "\n" + m.renderStatusBar()
}

func (m model) renderKeyDetail() string {
	if m.trustKeyIdx >= len(mockTrustedKeys) {
		return renderWithMargin("No key selected.") + "\n" + m.renderStatusBar()
	}
	k := mockTrustedKeys[m.trustKeyIdx]

	var statusBadge string
	if k.expired {
		statusBadge = expiredBadgeStyle.Render("● EXPIRED")
	} else {
		statusBadge = validBadgeStyle.Render("● VALID")
	}

	name := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorAccent)).
		Bold(true).
		Render(k.name)

	body := name + "  " + statusBadge + "\n\n" +
		keyValueLabelStyle.Render("Type:") + k.keyType + "\n" +
		keyValueLabelStyle.Render("Fingerprint:") + fingerprintStyle.Render(k.fingerprint) + "\n" +
		keyValueLabelStyle.Render("Added:") + k.addedDate + "\n" +
		keyValueLabelStyle.Render("Expires:") + k.expiresDate

	if k.expired {
		body += "\n\n" + warningBadgeStyle.Render(
			"!  Key has expired. Signatures from it won't be accepted.",
		)
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("Key Details") + "\n\n")
	b.WriteString(cardStyle.Render(body) + "\n\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(colorBlur)).Render(
		"esc: back    q: quit",
	))

	return renderWithMargin(b.String()) + "\n" + m.renderStatusBar()
}

func (m model) View() string {
	switch m.screen {
	case screenLoading:
		if m.errorMsg != "" {
			return renderWithMargin(
				titleStyle.Render("gittuf TUI") + "\n\n" +
					renderErrorMsg(m.errorMsg) + "\n\n" +
					lipgloss.NewStyle().Foreground(lipgloss.Color(colorBlur)).Render("Press Q or Ctrl+C to quit."),
			)
		}
		return renderWithMargin(
			titleStyle.Render("gittuf TUI") + "\n\n" +
				m.spinner.View() + " Loading, please wait...\n",
		)
	case screenChoice:
		body := renderWithMargin(m.choiceList.View() + "\n" + renderFooter(m.footer) + renderErrorMsg(m.errorMsg))
		return body + "\n" + m.renderStatusBar()
	case screenMockTrust:
		return m.renderMockTrust()
	case screenMockPolicy:
		return m.renderMockPolicy()
	case screenKeyDetail:
		return m.renderKeyDetail()
	case screenMockVerify:
		return m.renderMockVerify()
	case screenPolicy:
		return renderWithMargin(m.policyScreenList.View() + "\n" + renderFooter(m.footer) + renderErrorMsg(m.errorMsg))
	case screenTrust:
		return renderWithMargin(m.trustScreenList.View() + "\n" + renderFooter(m.footer) + renderErrorMsg(m.errorMsg))
	case screenPolicyRules:
		overlay := ""
		if m.confirmDelete {
			overlay = "\n" + renderDeleteOverlay(m.deleteTarget) + "\n"
		}
		hint := ""
		if !m.readOnly {
			hint = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color(colorSubtext)).Render(
				"Run `gittuf policy apply` to apply staged changes to the selected policy file.",
			)
		}

		emptyMsg := "No rules configured. Press 'A' to add one."
		return m.renderListScreen(m.ruleList, overlay+screenPolicyRulesHelp(m.readOnly)+hint, emptyMsg, len(m.rules) == 0)
	case screenTrustGlobalRules:
		overlay := ""
		if m.confirmDelete {
			overlay = "\n" + renderDeleteOverlay(m.deleteTarget) + "\n"
		}
		hint := ""
		if !m.readOnly {
			hint = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color(colorSubtext)).Render(
				"Run `gittuf trust apply` to apply staged changes to the selected policy file.",
			)
		}

		emptyMsg := "No rules configured. Press 'A' to add one."
		return m.renderListScreen(m.globalRuleList, overlay+screenTrustGlobalRulesHelp(m.readOnly)+hint, emptyMsg, len(m.globalRules) == 0)
	case screenPolicyAddRule:
		return m.renderFormScreen("Add Rule")
	case screenPolicyEditRule:
		return m.renderFormScreen("Edit Rule")
	case screenTrustAddGlobalRule:
		return m.renderFormScreen("Add Global Rule")
	case screenTrustEditGlobalRule:
		return m.renderFormScreen("Edit Global Rule")
	default:
		return "Unknown screen"
	}
}
