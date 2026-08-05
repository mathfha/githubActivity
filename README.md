# GitHub User Activity CLI

Aplikasi **GitHub User Activity CLI** adalah program berbasis Command Line Interface (CLI) yang dibuat menggunakan bahasa **Go**. Program ini mengambil aktivitas terbaru dari seorang pengguna GitHub melalui **GitHub REST API**, kemudian menampilkannya dalam format yang mudah dibaca di terminal.

Project ini dibuat sebagai sarana belajar mengenai:

- Pemrograman Go dasar
- HTTP Client (`net/http`)
- Parsing JSON (`encoding/json`)
- Command Line Interface (CLI) & Flag Parsing
- Penanganan error (Error Handling)
- Konsumsi REST API
- Struktur project Go yang ringkas dan *Lean*

---

## Fitur

- Mengambil aktivitas terbaru pengguna GitHub
- Menampilkan aktivitas dalam format yang mudah dibaca
- **Filter Berdasarkan Waktu**: Menampilkan aktivitas untuk hari ini, kemarin, atau beberapa hari ke belakang (menggunakan flag `-days`)
- Menerima username GitHub melalui command line
- Menangani error dengan pesan bantuan bawaan (`--help`)
- Menggunakan **library bawaan Go** tanpa dependensi pihak ketiga

---

## Teknologi

- Go
- GitHub REST API
- HTTP Client (`net/http`)
- JSON (`encoding/json`)
- CLI Flags (`flag`)

---

## Struktur Project

```text
github-user-activity/
│
├── cmd/
│   └── github-activity/
│       └── main.go          # Entry point aplikasi & logika filter
│
├── internal/
│   ├── github/
│   │   ├── client.go        # Komunikasi jaringan (HTTP Client)
│   │   └── models.go        # Definisi struktur data (JSON Struct)
│   │
│   └── formatter/
│       └── formatter.go     # Pencetakan teks terminal (CLI Output)
│
├── go.mod
└── README.md
```

---

## Cara Kerja

1. Pengguna menjalankan program melalui terminal.
2. Program menguraikan parameter CLI (`username` dan filter `-days`).
3. Program mengirim permintaan HTTP GET ke GitHub API.
4. GitHub mengembalikan JSON data aktivitas (maks 30 event).
5. Program *parse* JSON ke objek Go dan menyaring event jika diperlukan.
6. Hasilnya ditampilkan ke layar terminal secara dinamis.

---

## Endpoint API

Program menggunakan endpoint berikut:

```
GET https://api.github.com/users/{username}/events
```

---

## Cara Menjalankan

Clone repository:

```bash
git clone https://github.com/username/github-user-activity.git
```

Masuk ke folder project:

```bash
cd github-user-activity
```

Jalankan program (default 30 aktivitas terbaru tanpa filter waktu):

```bash
go run ./cmd/github-activity/main.go torvalds
```

Jalankan dengan memfilter jumlah hari ke belakang:

```bash
go run ./cmd/github-activity/main.go -days 0 torvalds   # Hanya aktivitas hari ini
go run ./cmd/github-activity/main.go -days 3 torvalds   # Aktivitas dalam 3 hari terakhir
```

Tampilkan bantuan (Help):

```bash
go run ./cmd/github-activity/main.go --help
```

Atau setelah di-build:

```bash
go build -o github-activity ./cmd/github-activity/main.go
./github-activity -days 1 torvalds
```

---

## Contoh Output

```text
Aktivitas torvalds dalam 3 hari terakhir:

• Melakukan push ke torvalds/linux
• Melakukan push 3 commit ke torvalds/linux
• Melakukan 'opened' Pull Request pada torvalds/linux
• Memberikan bintang (star) pada golang/go
• Membuat branch/tag baru di torvalds/linux
```

---

## Penanganan Error

Program akan memberikan pesan yang jelas apabila terjadi kesalahan, misalnya:

### Username tidak ditemukan

```text
Gagal mengambil data:
pengguna 'abcdef123456' tidak ditemukan (404)
```

### Username tidak diberikan

```text
GitHub Activity CLI
Menampilkan daftar aktivitas publik terbaru dari pengguna GitHub.

Cara penggunaan:
  github-activity [flags] <username>
...
```

---

## Event yang Didukung

Beberapa event GitHub yang dikenali:

| Event | Deskripsi |
|---------|----------|
| PushEvent | Push commit |
| CreateEvent | Membuat branch atau tag |
| DeleteEvent | Menghapus branch atau tag |
| PullRequestEvent | Pull Request |
| IssuesEvent | Issue |
| IssueCommentEvent | Komentar pada Issue |
| WatchEvent | Memberikan Star |
| ForkEvent | Fork repository |
| ReleaseEvent | Membuat Release |

Event lain yang belum didukung akan tetap ditangani dan ditampilkan secara aman sebagai format umum.

---

## Tujuan Pembelajaran

Melalui project ini diharapkan memahami:

- Menguasai `net/http` dan pembacaan *stream* IO
- Parsing JSON
- Clean Code (menghindari abstraksi yang terlalu *over-engineered* pada project kecil)
- Parsing CLI Flag dengan standard library
- Formatting Terminal

---

## Referensi

- Dokumentasi Go: https://go.dev/doc/
- GitHub REST API: https://docs.github.com/en/rest
- Roadmap.sh Project: https://roadmap.sh/projects/github-user-activity

---
## Lisensi

Project ini dibuat sebagai media pembelajaran dan dapat digunakan, dimodifikasi, maupun dikembangkan sesuai kebutuhan.