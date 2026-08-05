package formatter

import (
	"fmt"

	"github-activity/internal/github"
	"github.com/charmbracelet/lipgloss"
)

var (
	repoStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	actionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("36"))
	countStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true)
)

var simpleEvents = map[string]string{
	"CreateEvent":       "Membuat branch/tag baru di %s",
	"DeleteEvent":       "Menghapus branch/tag di %s",
	"IssueCommentEvent": "Memberikan komentar pada Issue di %s",
	"WatchEvent":        "Memberikan bintang (star) pada %s",
	"ForkEvent":         "Melakukan fork repo %s",
	"ReleaseEvent":      "Mempublikasikan release baru di %s",
}

func formatAction(action, eventName, repoName string) string {
	if action == "" {
		action = "aktivitas"
	}
	return fmt.Sprintf("Melakukan '%s' %s pada %s", actionStyle.Render(action), eventName, repoName)
}

func FormatEvent(event github.Event) string {
	repoName := repoStyle.Render(event.Repo.Name)

	if formatStr, ok := simpleEvents[event.Type]; ok {
		return fmt.Sprintf(formatStr, repoName)
	}

	switch event.Type {
	case "PushEvent":
		commitCount := event.Payload.Size
		if commitCount == 0 {
			commitCount = len(event.Payload.Commits)
		}

		if commitCount == 0 {
			return fmt.Sprintf("Melakukan push ke %s", repoName)
		}
		c := countStyle.Render(fmt.Sprintf("%d", commitCount))
		return fmt.Sprintf("Melakukan push %s commit ke %s", c, repoName)

	case "PullRequestEvent":
		return formatAction(event.Payload.Action, "Pull Request", repoName)

	case "IssuesEvent":
		return formatAction(event.Payload.Action, "Issue", repoName)

	default:
		return fmt.Sprintf("Melakukan aktivitas jenis '%s' di %s", actionStyle.Render(event.Type), repoName)
	}
}
