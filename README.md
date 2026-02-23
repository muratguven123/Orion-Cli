# orion-cli

Doğal dilde verilen geliştirme görevlerini analiz edip otomatik olarak GitHub Issue'larına dönüştüren yapay zeka destekli CLI aracı.

## Nedir?

`orion-cli`, yazılım geliştirme süreçlerini hızlandırmak için tasarlanmış bir komut satırı aracıdır. Bir geliştirme görevini doğal dilde yazmanız yeterlidir; araç bu metni Google Gemini yapay zekası aracılığıyla analiz ederek alt görevlere böler ve her birini otomatik olarak bulunduğunuz Git reposundaki GitHub Issues'a açar.

### Örnek Akış

```
Kullanıcı → "Kullanıcı kayıt olduğunda RabbitMQ'ya mesaj atacak bir servis yaz"
                         │
                         ▼
               Google Gemini AI (analiz)
                         │
                         ▼
          ┌──────────────────────────────┐
          │ Issue #1: RabbitMQ bağlantı  │
          │ Issue #2: Mesaj modeli       │
          │ Issue #3: Kayıt eventi       │
          └──────────────────────────────┘
                         │
                         ▼
              GitHub Issues (otomatik oluşturulur)
```

## Gereksinimler

- [Go](https://go.dev/) 1.22 veya üzeri
- [Google Gemini API anahtarı](https://aistudio.google.com/app/apikey)
- [GitHub Personal Access Token](https://github.com/settings/tokens) (`repo` kapsamı gerekli)
- Bir GitHub remote'u olan yerel bir Git reposu

## Kurulum

### 1. Repoyu klonlayın

```bash
git clone https://github.com/muratguven123/Orion-Cli.git
cd Orion-Cli
```

### 2. Bağımlılıkları yükleyin

```bash
go mod download
```

### 3. Derleyin

```bash
go build -o orion-cli .
```

İsterseniz binary'yi PATH'inize ekleyin:

```bash
mv orion-cli /usr/local/bin/
```

## Yapılandırma

Aşağıdaki ortam değişkenlerini tanımlamanız gerekir:

| Değişken         | Açıklama                                      | Zorunlu | Varsayılan           |
|------------------|-----------------------------------------------|---------|----------------------|
| `GEMINI_API_KEY` | Google Gemini API anahtarı                    | ✅ Evet  | —                    |
| `GITHUB_TOKEN`   | GitHub Personal Access Token (`repo` kapsamı) | ✅ Evet  | —                    |
| `GEMINI_MODEL`   | Kullanılacak Gemini model adı                 | ❌ Hayır | `gemini-2.5-flash`   |

### Linux / macOS

```bash
export GEMINI_API_KEY="AIza..."
export GITHUB_TOKEN="ghp_..."
```

### Windows (PowerShell)

```powershell
$env:GEMINI_API_KEY="AIza..."
$env:GITHUB_TOKEN="ghp_..."
```

## Kullanım

Komutun çalıştırılacağı dizin, GitHub'a bağlı bir Git reposu olmalıdır (`.git/config` içinde `origin` remote tanımlı olmalı).

```bash
orion-cli task "GÖREV_METNİ"
```

Veya doğrudan Go ile çalıştırmak isterseniz:

```bash
go run . task "GÖREV_METNİ"
```

### Örnekler

```bash
# Basit bir özellik planlaması
orion-cli task "Kullanıcı kayıt olduğunda RabbitMQ'ya mesaj atacak bir servis yaz"

# Servis mimarisi
orion-cli task "Ürün kataloğu için REST API yaz, CRUD işlemlerini desteklesin"

# Hata takibi
orion-cli task "Ödeme akışındaki timeout hatalarını araştır ve düzelt"
```

### Örnek Çıktı

```
✅ Issue'lar başarıyla oluşturuldu:
- #12 RabbitMQ bağlantı konfigürasyonu oluştur -> https://github.com/kullanici/repo/issues/12
- #13 Kullanıcı kayıt eventi için mesaj modeli tanımla -> https://github.com/kullanici/repo/issues/13
- #14 Kayıt servisine RabbitMQ publish entegrasyonu ekle -> https://github.com/kullanici/repo/issues/14
```

## Proje Mimarisi

Proje, Hexagonal Architecture (Ports & Adapters) prensibine göre tasarlanmıştır.

```text
.
├── main.go                        # Uygulama giriş noktası (gizli bilgi redaction içerir)
├── cmd/
│   ├── root.go                    # Cobra kök komutu
│   └── task.go                    # `task` alt komutu
├── internal/
│   ├── core/
│   │   ├── types.go               # Alan modelleri (IssueDraft, CreatedIssue)
│   │   ├── ports.go               # Arayüzler (TaskAnalyzer, IssueCreator, RepoLocator)
│   │   └── facade.go              # İş mantığını birleştiren facade
│   ├── ai/
│   │   └── gemini_adapter.go      # Google Gemini API adaptörü
│   ├── github/
│   │   └── issues_adapter.go      # GitHub Issues API adaptörü
│   ├── config/
│   │   └── env.go                 # Ortam değişkeni yükleme
│   └── gitrepo/
│       └── detector.go            # .git/config'den owner/repo tespiti
├── go.mod
└── README.md
```

### Katmanlar

| Katman     | Paket              | Sorumluluk                                      |
|------------|--------------------|-------------------------------------------------|
| CLI        | `cmd`              | Kullanıcı girişini alır, bağımlılıkları oluşturur |
| Core       | `internal/core`    | İş mantığı ve arayüz tanımları                  |
| AI Adapter | `internal/ai`      | Gemini API ile iletişim                         |
| GH Adapter | `internal/github`  | GitHub Issues API ile iletişim                  |
| Config     | `internal/config`  | Ortam değişkenlerini okur ve doğrular           |
| Git Detect | `internal/gitrepo` | Yerel repo'dan GitHub owner/repo bilgisini çeker |

## Testleri Çalıştırma

```bash
go test ./...
```

## Güvenlik

- `GEMINI_API_KEY` ve `GITHUB_TOKEN` değerleri asla kaynak koduna eklenmemelidir.
- Hata mesajlarında API anahtarı veya token bilgisi görünmemesi için `main.go` içinde otomatik redaction uygulanır.

## Bağımlılıklar

| Paket                             | Amaç                        |
|-----------------------------------|-----------------------------|
| `github.com/spf13/cobra`          | CLI çerçevesi               |
| `github.com/google/go-github/v66` | GitHub REST API istemcisi   |
| `golang.org/x/oauth2`             | GitHub token kimlik doğrulama |

## Lisans

Bu proje MIT lisansı ile lisanslanmıştır.
