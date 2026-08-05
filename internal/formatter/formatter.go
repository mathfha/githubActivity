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

func FormatEvent(event github.Event) string {
	repoName := repoStyle.Render(event.Repo.Name)

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
	case "CreateEvent":
		return fmt.Sprintf("Membuat branch/tag baru di %s", repoName)
	case "DeleteEvent":
		return fmt.Sprintf("Menghapus branch/tag di %s", repoName)
	case "PullRequestEvent":
		action := event.Payload.Action
		if action == "" {
			action = "aktivitas"
		}
		act := actionStyle.Render(action)
		return fmt.Sprintf("Melakukan '%s' Pull Request pada %s", act, repoName)
	case "IssuesEvent":
		action := event.Payload.Action
		if action == "" {
			action = "aktivitas"
		}
		act := actionStyle.Render(action)
		return fmt.Sprintf("Melakukan '%s' Issue pada %s", act, repoName)
	case "IssueCommentEvent":
		return fmt.Sprintf("Memberikan komentar pada Issue di %s", repoName)
	case "WatchEvent":
		return fmt.Sprintf("Memberikan bintang (star) pada %s", repoName)
	case "ForkEvent":
		return fmt.Sprintf("Melakukan fork repo %s", repoName)
	case "ReleaseEvent":
		return fmt.Sprintf("Mempublikasikan release baru di %s", repoName)
	default:
		return fmt.Sprintf("Melakukan aktivitas jenis '%s' di %s", actionStyle.Render(event.Type), repoName)
	}
}
