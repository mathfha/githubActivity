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
	"CreateEvent":       "Created a new branch/tag in %s",
	"DeleteEvent":       "Deleted a branch/tag in %s",
	"IssueCommentEvent": "Commented on an issue in %s",
	"WatchEvent":        "Starred %s",
	"ForkEvent":         "Forked %s",
	"ReleaseEvent":      "Published a new release in %s",
}

func formatAction(action, eventName, repoName string) string {
	if action == "" {
		action = "activity"
	}
	return fmt.Sprintf("Performed '%s' on %s in %s", actionStyle.Render(action), eventName, repoName)
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
			return fmt.Sprintf("Pushed to %s", repoName)
		}
		c := countStyle.Render(fmt.Sprintf("%d", commitCount))
		return fmt.Sprintf("Pushed %s commits to %s", c, repoName)

	case "PullRequestEvent":
		return formatAction(event.Payload.Action, "Pull Request", repoName)

	case "IssuesEvent":
		return formatAction(event.Payload.Action, "Issue", repoName)

	default:
		return fmt.Sprintf("Performed '%s' in %s", actionStyle.Render(event.Type), repoName)
	}
}
