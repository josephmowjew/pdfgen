# PDFGen

A powerful and flexible PDF generation library for Go that enables dynamic PDF creation from structured data with optional cloud storage integration.

## Features

- 📄 Dynamic PDF generation from structured data
- 🎨 Template-based PDF generation with customizable layouts
- ☁️ Built-in upload functionality with configurable endpoints
- 🧩 Clean and modular architecture
- 🔧 Customizable page layouts and styling
- 🎯 Simple and intuitive API

## Installation

Requires Go 1.22 or higher.

```bash
go get github.com/josephmojoo/pdfgen
```

## Quick Start

Here's a simple example of generating a PDF:

```go
package main

import (
    "context"
    "log"
    "os"
    "github.com/josephmojoo/pdfgen/pkg/pdf/service"
)

func main() {
    // Initialize the PDF service
    svc := service.New(service.Config{})

    // Sample data
    data := map[string]interface{}{
        "customer": map[string]interface{}{
            "name":    "John Doe",
            "email":   "john@example.com",
            "address": "123 Main St",
        },
        "order": map[string]interface{}{
            "id":     "ORD-12345",
            "date":   "2024-03-15",
            "status": "Completed",
            "total":  109.97,
        },
    }

    // Generate PDF
    pdfData, err := svc.GenerateOnly(context.Background(), data)
    if err != nil {
        log.Fatalf("Failed to generate PDF: %v", err)
    }

    // Save to file
    if err := os.WriteFile("output.pdf", pdfData, 0644); err != nil {
        log.Fatalf("Failed to save PDF: %v", err)
    }
}
```

## Advanced Usage

### Generate and Upload PDF

```go
response, err := svc.GenerateAndUpload(context.Background(), data, service.UploadConfig{
    OrganizationalID: "your-org-id",
    BranchID:         "your-branch-id",
    CreatedBy:        "user-name",
    FileName:         "report.pdf",
})
if err != nil {
    log.Fatalf("Failed to generate and upload PDF: %v", err)
}

fmt.Printf("PDF uploaded successfully: %s\n", response.FileDownloadUri)
```

### Service Configuration

```go
svc := service.New(service.Config{
    UploadBaseURL: "https://your-storage-service.com/files",
    BearerToken:   "your-auth-token",
})
```

## Project Structure

```
.
├── pkg/
│   └── pdf/
│       ├── generator/     # PDF generation logic
│       ├── service/       # Main service interface
│       ├── model/         # Data models
│       └── errors/        # Error definitions
├── example/              # Usage examples
└── cmd/                  # Command line tools
```

## API Reference

### Service Interface

```go
type Service interface {
    GenerateAndUpload(ctx context.Context, data map[string]interface{}, config UploadConfig) (*UploadResponse, error)
    GenerateOnly(ctx context.Context, data map[string]interface{}) ([]byte, error)
}
```

### Configuration Types

```go
type Config struct {
    UploadBaseURL string
    BearerToken   string
}

type UploadConfig struct {
    OrganizationalID string
    BranchID         string
    CreatedBy        string
    FileName         string
}
```

### Upload Response

```go
type UploadResponse struct {
    FileName        string
    FileDownloadUri string
    FileType        string
    Size            int64
}
```

## Data Structure

The library accepts any structured data as a map[string]interface{}. Here's an example of supported data structure:

```go
data := map[string]interface{}{
    "customer": map[string]interface{}{
        "name":    "John Doe",
        "email":   "john@example.com",
        "address": "123 Main St",
    },
    "order": map[string]interface{}{
        "id":     "ORD-12345",
        "date":   "2024-03-15",
        "status": "Completed",
        "items": []interface{}{
            map[string]interface{}{
                "name":     "Product A",
                "quantity": 2,
                "price":    29.99,
            },
            map[string]interface{}{
                "name":     "Product B",
                "quantity": 1,
                "price":    49.99,
            },
        },
        "total": 109.97,
    },
}
```

## Error Handling

The library provides custom error types for better error handling:

```go
type PDFError struct {
    Code    ErrorCode
    Message string
    Cause   error
}
```

Common error codes:
- `ErrInvalidTemplate`
- `ErrInvalidData`
- `ErrRenderFailed`
- `ErrLayoutFailed`
- `ErrGenerationFailed`

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Testing

Run the tests:

```bash
go test ./...
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [jung-kurt/gofpdf](https://github.com/jung-kurt/gofpdf) - PDF generation library 