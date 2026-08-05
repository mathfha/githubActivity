package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github-activity/internal/formatter"
	"github-activity/internal/github"
)

type state int

const (
	stateLoading state = iota
	stateLoaded
	stateError
)

type Model struct {
	client   github.ClientAPI
	username string
	days     int

	state    state
	spinner  spinner.Model
	viewport viewport.Model
	events   []github.Event
	err      error

	ready bool
	
	// Navigation
	activeTab int
	tabs      []string
}

func NewModel(client github.ClientAPI, username string, days int) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		client:   client,
		username: username,
		days:     days,
		state:    stateLoading,
		spinner:  s,
		tabs:     []string{"Aktivitas", "Repositories"},
		activeTab: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.fetchDataCmd,
	)
}

type dataFetchedMsg struct {
	events []github.Event
}

type errMsg struct {
	err error
}

func (m Model) fetchDataCmd() tea.Msg {
	ctx := context.Background()
	events, err := m.client.FetchUserActivity(ctx, m.username)
	if err != nil {
		return errMsg{err: err}
	}

	if m.days >= 0 {
		var filtered []github.Event
		now := time.Now()
		threshold := time.Date(now.Year(), now.Month(), now.Day()-m.days, 0, 0, 0, 0, now.Location())

		for _, e := range events {
			if e.CreatedAt.After(threshold) {
				filtered = append(filtered, e)
			}
		}
		events = filtered
	}

	return dataFetchedMsg{events: events}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch k := msg.String(); k {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "tab", "right", "l":
			m.activeTab = (m.activeTab + 1) % len(m.tabs)
		case "shift+tab", "left", "h":
			m.activeTab = (m.activeTab - 1 + len(m.tabs)) % len(m.tabs)
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

	case dataFetchedMsg:
		m.events = msg.events
		m.state = stateLoaded

		var b strings.Builder
		if len(m.events) == 0 {
			b.WriteString("Tidak ada aktivitas publik.\n")
		} else {
			for _, e := range m.events {
				b.WriteString(fmt.Sprintf("• %s\n", formatter.FormatEvent(e)))
			}
		}

		m.viewport.SetContent(b.String())
		return m, nil

	case errMsg:
		m.err = msg.err
		m.state = stateError
		return m, nil
	}

	if m.state == stateLoading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.state == stateLoaded {
		// Hanya viewport aktivitas yang butuh update scrolling
		if m.activeTab == 0 {
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) headerView() string {
	// Judul
	title := fmt.Sprintf(" GitHub Activity: %s ", m.username)
	if m.days == 0 {
		title += "(Hari Ini) "
	} else if m.days > 0 {
		title += fmt.Sprintf("(%d Hari Terakhir) ", m.days)
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		MarginBottom(1)
		
	titleRendered := titleStyle.Render(title)

	// Tabs
	activeTabStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("205")).
		Padding(0, 1)

	inactiveTabStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	var renderedTabs []string
	for i, t := range m.tabs {
		if i == m.activeTab {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(t))
		} else {
			renderedTabs = append(renderedTabs, inactiveTabStyle.Render(t))
		}
	}

	// Gabungkan tab secara horizontal
	tabsRendered := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	return titleRendered + "\n" + tabsRendered + "\n"
}

func (m Model) footerView() string {
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)
	if m.state == stateLoaded {
		if m.activeTab == 0 {
			return helpStyle.Render("Tab/Kanan: Ganti Tab • Atas/Bawah: Scroll List • Tekan 'q' untuk keluar.")
		}
		return helpStyle.Render("Tab/Kiri: Ganti Tab • Tekan 'q' untuk keluar.")
	}
	return helpStyle.Render("Tekan 'q' untuk keluar.")
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Inisialisasi...\n"
	}

	if m.state == stateLoading {
		return fmt.Sprintf("\n %s Mengambil data %s...\n\n%s", m.spinner.View(), m.username, m.footerView())
	}

	if m.state == stateError {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
		return fmt.Sprintf("\n%s\n%v\n\n%s", errorStyle.Render("Terjadi Kesalahan:"), m.err, m.footerView())
	}

	// Konten berdasarkan tab aktif
	var content string
	if m.activeTab == 0 {
		content = m.viewport.View()
	} else if m.activeTab == 1 {
		placeholderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Italic(true).Padding(1, 2)
		content = placeholderStyle.Render(fmt.Sprintf("Fitur list Repositories untuk %s belum diimplementasikan.\n\nSegera datang di update berikutnya!", m.username))
	}

	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), content, m.footerView())
}
