# GitHub User Activity CLI

The **GitHub User Activity CLI** is a terminal-based interactive application built with **Go** and the **Charmbracelet** ecosystem (`bubbletea`, `lipgloss`). It fetches the latest public activity of a specified GitHub user via the **GitHub REST API** and presents it in a beautiful, colorful, and interactive Terminal User Interface (TUI).

This project serves as a practical learning medium for:

- Fundamental Go programming
- HTTP Clients (`net/http`) & JSON Parsing (`encoding/json`)
- Building Terminal User Interfaces (TUI) using the Elm architecture (*Model-Update-View*)
- Terminal layout and styling with `lipgloss`
- Utilizing interactive components such as `spinner`, `viewport`, and Tab navigation
- Consuming RESTful APIs

---

## Screenshots

### 1. Activity List View (Scrollable)
![Activity View](.assets/activity.jpg)

### 2. Repositories Tab Navigation
![Repositories Tab](.assets/repo.jpg)

### 3. Help Menu
![Help Menu](.assets/help.jpg)

---

## Key Features

- **Interactive UI**: Built on top of Bubble Tea, providing a seamless and highly responsive terminal experience.
- **Tabbed Navigation**: Easily switch contexts (e.g., *Activities* and *Repositories*) using the `Left/Right` arrow keys or the `Tab` key.
- **Loading Animations**: Displays a smooth `spinner` while asynchronously fetching data from the GitHub API.
- **Scrollable Viewport**: Handles extensive activity logs effortlessly by allowing users to scroll up and down if the content exceeds the terminal height.
- **Time-based Filtering**: Filter activities from today, yesterday, or several days ago using the `-days` flag.
- **Dynamic Styling**: Beautifully highlights repository names and commit counts using customizable `lipgloss` color palettes.

---

## Technologies Used

- **Language**: Go 1.26
- **API**: GitHub REST API
- **TUI Libraries**:
  - `github.com/charmbracelet/bubbletea` (Core framework)
  - `github.com/charmbracelet/bubbles` (Spinner & Viewport components)
  - `github.com/charmbracelet/lipgloss` (Styling and layout engine)

---

## Project Structure

```text
github-user-activity/
│
├── cmd/
│   └── github-activity/
│       └── main.go          # Entry point & Bubble Tea program initialization
│
├── internal/
│   ├── github/
│   │   ├── client.go        # Network communication (HTTP Client)
│   │   └── models.go        # Data structure definitions (JSON Structs)
│   │
│   ├── ui/
│   │   └── tui.go           # Interactive TUI logic (Model, Update, View)
│   │
│   └── formatter/
│       └── formatter.go     # Event formatting integrated with lipgloss
│
├── .assets/                 # Screenshot assets for documentation
├── go.mod
├── go.sum
└── README.md
```

---

## How It Works

1. The user executes the application from the terminal by providing a `username` parameter.
2. The program launches the `bubbletea` UI in full-screen mode (AltScreen) and displays a loading animation.
3. A background command (`tea.Cmd`) sends an asynchronous HTTP GET request to the GitHub API.
4. Once the JSON response is parsed and filtered, the program dispatches a message to the TUI.
5. The interface hides the loading state and updates the View to render the activity list inside a scrollable Viewport.
6. The user can safely exit the application at any time by pressing `q`, `Esc`, or `Ctrl+C`.

---

## Getting Started

Clone the repository:

```bash
git clone https://github.com/username/github-user-activity.git
cd github-user-activity
```

Run the program (automatically downloads Go module dependencies):

```bash
go run ./cmd/github-activity <username>
```

Run with time-based filters:

```bash
go run ./cmd/github-activity -days 0 <username>   # Only show today's activities
go run ./cmd/github-activity -days 3 <username>   # Show activities from the last 3 days
```

Display help manual:

```bash
go run ./cmd/github-activity --help
```

Alternatively, build it as a standalone binary:

```bash
go build -o github-activity ./cmd/github-activity
./github-activity -days 1 <username>
```

---

## Supported GitHub Events

| Event | UI Output Translation (ID) |
|---------|----------|
| PushEvent | Melakukan push X commit ke repo |
| CreateEvent | Membuat branch atau tag |
| DeleteEvent | Menghapus branch atau tag |
| PullRequestEvent | Melakukan pull request |
| IssuesEvent | Melakukan issue |
| IssueCommentEvent | Memberikan komentar pada Issue |
| WatchEvent | Memberikan bintang (star) |
| ForkEvent | Melakukan fork repo |
| ReleaseEvent | Mempublikasikan release baru |

Unsupported events will fallback to a generic, safe action format indicating the exact event type performed by the user.

---

## License

This project was built for educational purposes and is open for use, modification, and further development based on your needs.