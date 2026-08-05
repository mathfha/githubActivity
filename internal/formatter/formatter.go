package formatter

import (
	"fmt"

	"github-activity/internal/github"
)

func FormatEvent(event github.Event) string {
	repoName := event.Repo.Name

	switch event.Type {
	case "PushEvent":
		commitCount := event.Payload.Size
		if commitCount == 0 {
			commitCount = len(event.Payload.Commits)
		}
		
		if commitCount == 0 {
			return fmt.Sprintf("Melakukan push ke %s", repoName)
		}
		return fmt.Sprintf("Melakukan push %d commit ke %s", commitCount, repoName)
	case "CreateEvent":
		return fmt.Sprintf("Membuat branch/tag baru di %s", repoName)
	case "DeleteEvent":
		return fmt.Sprintf("Menghapus branch/tag di %s", repoName)
	case "PullRequestEvent":
		action := event.Payload.Action
		if action == "" {
			action = "aktivitas"
		}
		return fmt.Sprintf("Melakukan '%s' Pull Request pada %s", action, repoName)
	case "IssuesEvent":
		action := event.Payload.Action
		if action == "" {
			action = "aktivitas"
		}
		return fmt.Sprintf("Melakukan '%s' Issue pada %s", action, repoName)
	case "IssueCommentEvent":
		return fmt.Sprintf("Memberikan komentar pada Issue di %s", repoName)
	case "WatchEvent":
		return fmt.Sprintf("Memberikan bintang (star) pada %s", repoName)
	case "ForkEvent":
		return fmt.Sprintf("Melakukan fork repo %s", repoName)
	case "ReleaseEvent":
		return fmt.Sprintf("Mempublikasikan release baru di %s", repoName)
	default:
		return fmt.Sprintf("Melakukan aktivitas jenis '%s' di %s", event.Type, repoName)
	}
}

func PrintActivity(username string, events []github.Event, days int) {
	if len(events) == 0 {
		if days == 0 {
			fmt.Printf("Tidak ada aktivitas publik terbaru untuk pengguna '%s' hari ini\n", username)
		} else if days > 0 {
			fmt.Printf("Tidak ada aktivitas publik terbaru untuk pengguna '%s' dalam %d hari terakhir\n", username, days)
		} else {
			fmt.Printf("Tidak ada aktivitas publik terbaru untuk pengguna '%s'\n", username)
		}
		return
	}

	if days == 0 {
		fmt.Printf("Aktivitas %s hari ini:\n\n", username)
	} else if days > 0 {
		fmt.Printf("Aktivitas %s dalam %d hari terakhir:\n\n", username, days)
	} else {
		fmt.Printf("Aktivitas terbaru %s:\n\n", username)
	}
	
	for _, event := range events {
		fmt.Printf("• %s\n", FormatEvent(event))
	}
}
