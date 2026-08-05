package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github-activity/internal/github"
	"github-activity/internal/ui"
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

	client := github.NewClient()
	
	m := ui.NewModel(client, username, *daysPtr)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Terjadi kesalahan saat menjalankan antarmuka: %v\n", err)
		os.Exit(1)
	}
}
