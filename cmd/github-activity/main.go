package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github-activity/internal/formatter"
	"github-activity/internal/github"
)

func main() {
	flag.Usage = func() {
		fmt.Println("GitHub Activity CLI")
		fmt.Println("Menampilkan daftar aktivitas publik terbaru dari pengguna GitHub.")
		fmt.Println("\nCara penggunaan:")
		fmt.Println("  github-activity [flags] <username>")
		fmt.Println("\nFlags yang tersedia:")
		flag.PrintDefaults()
		fmt.Println("\nContoh:")
		fmt.Println("  github-activity mathfha           (Menampilkan default 30 aktivitas terbaru)")
		fmt.Println("  github-activity -days 0 mathfha   (Menampilkan aktivitas hari ini saja)")
		fmt.Println("  github-activity -days 3 mathfha   (Menampilkan aktivitas 3 hari terakhir)")
	}

	daysPtr := flag.Int("days", -1, "Jumlah hari ke belakang (0 = hari ini, 1 = kemarin, dst)")
	flag.Parse()

	username := flag.Arg(0)

	if username == "" {
		flag.Usage()
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := github.NewClient()

	events, err := client.FetchUserActivity(ctx, username)
	if err != nil {
		fmt.Printf("Gagal mengambil data:\n%v\n", err)
		os.Exit(1)
	}

	days := *daysPtr
	if days >= 0 {
		var filteredEvents []github.Event
		
		now := time.Now()
		threshold := time.Date(now.Year(), now.Month(), now.Day()-days, 0, 0, 0, 0, now.Location())
		
		for _, event := range events {
			if event.CreatedAt.After(threshold) {
				filteredEvents = append(filteredEvents, event)
			}
		}
		events = filteredEvents
	}

	formatter.PrintActivity(username, events, days)
}
