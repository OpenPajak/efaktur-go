# Efaktur-Go Code Wiki

Welcome to the **Efaktur-Go** Code Wiki! This documentation provides a comprehensive overview of the repository, including its architecture, modules, key components, dependencies, and instructions on how to use and run the project.

---

## 1. Overall Project Architecture

**Efaktur-Go** (`github.com/OpenPajak/efaktur-go`) is a library and command-line interface (CLI) tool designed to interact with the Indonesian Web E-Faktur service (`web-efaktur.pajak.go.id`). 

The architecture is divided into two primary parts:
1. **The SDK Library (`pkg/provider/web`)**: A Go package that wraps the E-Faktur Web API. It handles PKCS#12 certificate-based TLS authentication, session management (login/logout via cookies), and defines structured payloads for the API endpoints.
2. **The CLI Tool (`cmd/efakturtool`)**: A terminal application built on top of the SDK. It allows users to dump E-Faktur records, retrieve prepopulated data, and output the results in various formats (JSON, YAML, CSV, Wide).

---

## 2. Responsibilities of Major Modules

### `cmd/efakturtool` (CLI Application)
The entry point for the end-user CLI application.
- **[main.go](file:///workspace/cmd/efakturtool/main.go)**: The application entry point. It parses flags, initializes configuration via `go-binder`, sets up the `web.Client` with TLS certificates, and manages the pre-run (login/profile fetch) and post-run (logout) hooks.
- **[cmd_dump.go](file:///workspace/cmd/efakturtool/cmd_dump.go)**: Implements the `dump` command. It concurrently fetches SPT Headers and their Lampiran (Details) using goroutines, and aggregates the results.
- **[cmd_prepopulated.go](file:///workspace/cmd/efakturtool/cmd_prepopulated.go)**: Implements the `prepopulated` command for downloading prepopulated CSV files from DJP.
- **[formatter.go](file:///workspace/cmd/efakturtool/formatter.go)**: Provides dynamic output formatting for the CLI, allowing data to be printed as wide-column text, CSV, JSON, or YAML.

### `pkg/provider/web` (Core SDK)
The HTTP client library for the DJP Web E-Faktur service.
- **[client.go](file:///workspace/cmd/pkg/provider/web/client.go)**: Contains the core `Client` struct, handles HTTP connections, CookieJars, certificate-based TLS configuration, and generic REST invocations (`restInvoke`).
- **[endpoint.go](file:///workspace/cmd/pkg/provider/web/endpoint.go)**: Centralized registry for all Web E-Faktur API endpoints (e.g., Login, SptHeaderList, PrepopulatedDownload).
- **[spt.go](file:///workspace/cmd/pkg/provider/web/spt.go)**: Defines enums and constants such as `KodeFormSpt` (A1, A2, B1, B2, B3) and `PrepopulatedJenisDokumen` (FPM, PIB, PEB, etc.) generated using `enumer`.
- **`ext_*.go`**: Domain-specific API extensions attached to the main `Client`. Examples include:
  - **[ext_spt_header.go](file:///workspace/pkg/provider/web/ext_spt_header.go)**: Manages SPT headers (Posting, Delete, Cetak).
  - **[ext_spt_detail.go](file:///workspace/pkg/provider/web/ext_spt_detail.go)**: Retrieves individual invoice details (Lampiran).
  - **[ext_profile.go](file:///workspace/pkg/provider/web/ext_profile.go)**: Retrieves the PKP (Pengusaha Kena Pajak) profile data.

---

## 3. Key Classes and Functions

### `web.Client`
The primary HTTP client structure for all communications.
- **`web.NewClient(opts ClientOptions)`**: Initializes the client. Requires `ClientOptions` which must contain the parsed PKCS#12 TLS certificate and Certificate Authorities.
- **`Client.Login(ctx, password)` / `Client.Logout(ctx)`**: Handles authentication. It sends the service password, checks redirection rules, and manages the session cookies.
- **`web.restInvoke[Req, Resp]()`**: A generic utility function used internally by extension modules to make REST API requests. It handles JSON marshaling, execution, and unmarshaling.

### `web.PKCS12ToTLSCertificateFromFile`
Reads and decodes the PKCS#12 (`.p12` / `.pfx`) certificate required by the DJP server for MTLS (Mutual TLS) authentication. Returns a `tls.Certificate` used in `ClientOptions`.

### `cmd[T]` & `Commands`
CLI command routing interfaces in `efakturtool`. Each command (like `cmdDump`) defines its own flags, arguments (`data`), and a `callback` function executed when the command is called.

---

## 4. Dependency Relationships

The project relies on a few critical external libraries to function:
- **`github.com/ii64/go-binder`**: Used in `efakturtool` to bind configuration files (TOML, YAML, JSON) to Go structs automatically.
- **`software.sslmate.com/src/go-pkcs12`** (Replaced by `github.com/ii64/go-pkcs12` in `go.mod`): Crucial for decoding the PKCS#12 digital certificates issued by DJP. Standard Go library does not fully support all PKCS#12 encryption schemes used by these certificates.
- **`golang.org/x/net`**: Provides `publicsuffix` which is used to initialize the HTTP `cookiejar`, maintaining the session properly across API calls.
- **`github.com/dmarkham/enumer`**: Used as a `go:generate` tool to automatically generate JSON, YAML, and SQL un/marshalers for enumeration types like `KodeFormSpt`.

---

## 5. Instructions for Running the Project

### Prerequisites
- Go 1.20 or higher.
- A valid DJP E-Faktur digital certificate (`.p12` file) and its passphrase.
- Your E-Faktur / E-Nova service password.

### Installation
You can install the CLI tool directly using `go install`:
```bash
go install github.com/OpenPajak/efaktur-go/cmd/efakturtool@latest
```
Alternatively, clone the repository and build it manually:
```bash
git clone https://github.com/OpenPajak/efaktur-go.git
cd efaktur-go
go build -o efakturtool ./cmd/efakturtool
```

### Configuration
Create an `efaktur.toml` file in the same directory where you run the tool (or set the `CONFIG` environment variable). The configuration requires your certificate details:

```toml
[conf]
  certPath = "./path/to/your_certificate.p12"
  certPasswd = "YOUR_CERT_PASSPHRASE"
  svcPasswd = "YOUR_SERVICE_PASSWORD"
```
*(Note: You can also pass these via CLI arguments `--certPath`, `--certPasswd`, and `--svcPasswd` or via environment variables `CERT_PATH`, `CERT_PASS`, `SVC_PASSWD`)*

### Usage Examples

1. **Dump Lampiran Data**
Dump SPT detail data (Lampiran A1, A2, B1, etc.) for month `1` (January) and year `2023`, outputting the result in JSON format:
```bash
efakturtool dump -m 1 -y 2023 -t A1,B1,B2 -o json
```

2. **Download Prepopulated Data**
Retrieve prepopulated CSV data for Faktur Pajak Masukan (FPM) for December 2022:
```bash
efakturtool prepopulated -m 12 -y 2022 -t FPM -o csv
```

By default, the tool outputs in a `wide` tabular format if `-o` is omitted. Available formats are `wide`, `json`, `yaml`, and `csv`.
