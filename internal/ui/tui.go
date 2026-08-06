package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github-activity/internal/formatter"
	"github-activity/internal/github"
)

type state int

const (
	stateIdle state = iota
	stateLoading
	stateLoaded
	stateError
)

type Model struct {
	client   github.ClientAPI
	username string
	days     int

	state        state
	spinner      spinner.Model
	viewport     viewport.Model
	repoViewport viewport.Model
	searchBar    textinput.Model
	
	events        []github.Event
	repos         []github.Repository
	
	err          error
	ready        bool
	
	eventsLoaded bool
	reposLoaded  bool

	// Navigation
	activeTab int
	tabs      []string
}

func parseSearchQuery(query string) (string, int) {
	parts := strings.Fields(query)
	if len(parts) >= 3 && strings.ToLower(parts[0]) == "days" {
		var days int
		if _, err := fmt.Sscanf(parts[1], "%d", &days); err == nil {
			return strings.Join(parts[2:], " "), days
		}
	}
	return query, -1
}

func NewModel(client github.ClientAPI) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	ti := textinput.New()
	ti.Placeholder = "Enter GitHub username or command (e.g., days 0 torvalds)..."
	ti.CharLimit = 100
	ti.Width = 60
	ti.Focus()

	return Model{
		client:    client,
		username:  "",
		days:      -1,
		state:     stateIdle,
		spinner:   s,
		searchBar: ti,
		tabs:      []string{"Activity", "Repositories"},
		activeTab: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

type dataFetchedMsg struct {
	events []github.Event
}

type reposFetchedMsg struct {
	repos []github.Repository
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

func (m Model) fetchReposCmd() tea.Msg {
	ctx := context.Background()
	repos, err := m.client.FetchUserRepositories(ctx, m.username)
	if err != nil {
		return errMsg{err: err}
	}
	return reposFetchedMsg{repos: repos}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch k := msg.String(); k {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab":
			m.activeTab = (m.activeTab + 1) % len(m.tabs)
			return m, nil
		case "shift+tab":
			m.activeTab = (m.activeTab - 1 + len(m.tabs)) % len(m.tabs)
			return m, nil
		case "enter":
			val := strings.TrimSpace(m.searchBar.Value())
			if val != "" {
				parsedUsername, parsedDays := parseSearchQuery(val)
				m.username = parsedUsername
				m.days = parsedDays
				
				m.state = stateLoading
				m.eventsLoaded = false
				m.reposLoaded = false
				return m, tea.Batch(
					m.spinner.Tick,
					m.fetchDataCmd,
					m.fetchReposCmd,
				)
			}
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-(headerHeight+footerHeight))
			m.viewport.YPosition = headerHeight

			m.repoViewport = viewport.New(msg.Width, msg.Height-(headerHeight+footerHeight))
			m.repoViewport.YPosition = headerHeight

			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - (headerHeight + footerHeight)
			
			m.repoViewport.Width = msg.Width
			m.repoViewport.Height = msg.Height - (headerHeight + footerHeight)
		}

	case dataFetchedMsg:
		m.events = msg.events
		m.eventsLoaded = true
		if m.reposLoaded {
			m.state = stateLoaded
		}
		m.updateEventsViewport()
		return m, nil

	case reposFetchedMsg:
		m.repos = msg.repos
		m.reposLoaded = true
		if m.eventsLoaded {
			m.state = stateLoaded
		}
		m.updateReposViewport()
		return m, nil

	case errMsg:
		m.err = msg.err
		m.state = stateError
		return m, nil
	}

	// Update search bar textinput globally
	m.searchBar, cmd = m.searchBar.Update(msg)
	cmds = append(cmds, cmd)

	if m.state == stateLoading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.state == stateLoaded {
		if m.activeTab == 0 {
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.activeTab == 1 {
			m.repoViewport, cmd = m.repoViewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) updateEventsViewport() {
	var b strings.Builder
	if len(m.events) == 0 {
		b.WriteString("No public activity found.\n")
	} else {
		for _, e := range m.events {
			b.WriteString(fmt.Sprintf("• %s\n", formatter.FormatEvent(e)))
		}
	}
	m.viewport.SetContent(b.String())
}

func (m *Model) updateReposViewport() {
	var b strings.Builder
	if len(m.repos) == 0 {
		b.WriteString("No public repositories found.\n")
	} else {
		titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("36")).Bold(true)
		descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
		statStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("220"))

		for _, r := range m.repos {
			stats := fmt.Sprintf("⭐ %d", r.StargazersCount)
			if r.Language != "" {
				stats += fmt.Sprintf(" • %s", r.Language)
			}
			if r.Fork {
				stats += " • (Fork)"
			}
			b.WriteString(fmt.Sprintf(" %s | %s\n", titleStyle.Render(r.Name), statStyle.Render(stats)))
			if r.Description != "" {
				b.WriteString("  " + descStyle.Render(r.Description) + "\n")
			}
			b.WriteString("\n")
		}
	}
	m.repoViewport.SetContent(b.String())
}

func (m Model) headerView() string {
	// Search Bar (Global)
	searchView := m.searchBar.View()
	
	// Jika belum ada pencarian, cukup tampilkan search bar saja
	if m.state == stateIdle {
		return "\n" + searchView + "\n"
	}

	// Title
	title := fmt.Sprintf(" Showing results for: %s ", m.username)
	if m.days == 0 {
		title += "(Today) "
	} else if m.days > 0 {
		title += fmt.Sprintf("(Last %d Days) ", m.days)
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		MarginBottom(1).
		MarginTop(1)

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

	tabsRendered := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	return "\n" + searchView + "\n" + titleRendered + "\n" + tabsRendered + "\n"
}

func (m Model) footerView() string {
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)
	if m.state == stateLoaded {
		return helpStyle.Render("Up/Down: Scroll List • Tab/Shift+Tab: Switch Tab • Enter: Search New User • Esc: Quit")
	}
	return helpStyle.Render("Enter: Search • Esc: Quit")
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing...\n"
	}

	if m.state == stateIdle {
		guideStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).MarginTop(1)
		
		var guide strings.Builder
		guide.WriteString(guideStyle.Render("Welcome to GitHub Activity CLI!") + "\n\n")
		guide.WriteString("Search bar usage:\n")
		guide.WriteString("  torvalds          (Fetch the latest 30 public activities)\n")
		guide.WriteString("  days 0 torvalds   (Fetch activities from today only)\n")
		guide.WriteString("  days 3 torvalds   (Fetch activities from the last 3 days)\n")
		
		return fmt.Sprintf("%s\n%s\n\n%s", m.headerView(), guide.String(), m.footerView())
	}

	if m.state == stateLoading {
		return fmt.Sprintf("%s\n %s Fetching data for %s...\n\n%s", m.headerView(), m.spinner.View(), m.username, m.footerView())
	}

	if m.state == stateError {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
		return fmt.Sprintf("%s\n%s\n%v\n\n%s", m.headerView(), errorStyle.Render("An Error Occurred:"), m.err, m.footerView())
	}

	var content string
	if m.activeTab == 0 {
		content = m.viewport.View()
	} else if m.activeTab == 1 {
		content = m.repoViewport.View()
	}

	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), content, m.footerView())
}
