# orion-cli

Doğal dilde verilen geliştirme görevlerini analiz edip GitHub Issue'larına dönüştüren CLI.

## Kurulum (adım adım)

```bash
mkdir orion-cli
cd orion-cli

go mod init orion-cli
go get github.com/spf13/cobra@latest
go get github.com/google/go-github/v66/github@latest
go get golang.org/x/oauth2@latest
```

## Klasör yapısı

```text
.
├── cmd/
│   ├── root.go
│   └── task.go
├── internal/
│   ├── ai/
│   │   └── gemini_adapter.go
│   ├── config/
│   │   └── env.go
│   ├── core/
│   │   ├── facade.go
│   │   ├── ports.go
│   │   └── types.go
│   ├── github/
│   │   └── issues_adapter.go
│   └── gitrepo/
│       └── detector.go
├── main.go
├── go.mod
└── README.md
```

## Kullanım

Önce ortam değişkenlerini tanımlayın:

```bash
export GEMINI_API_KEY="..."
export GITHUB_TOKEN="..."
```

Çalıştırma:

```bash
go run . task "Kullanıcı kayıt olduğunda RabbitMQ'ya mesaj atacak bir servis yazacağım, bunu planla"
```
