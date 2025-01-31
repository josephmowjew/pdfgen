package service

import (
	"context"
	"testing"
)

func TestService_GenerateOnly(t *testing.T) {
	// Initialize service
	svc := New(Config{
		UploadBaseURL: "https://example.com",
		BearerToken:   "test-token",
	})

	// Test data
	data := map[string]interface{}{
		"title": "Test Document",
		"items": []interface{}{
			map[string]interface{}{
				"name":  "Item 1",
				"value": 100,
			},
		},
	}

	// Test PDF generation
	ctx := context.Background()
	pdfData, err := svc.GenerateOnly(ctx, data)
	if err != nil {
		t.Errorf("GenerateOnly() error = %v", err)
		return
	}

	if len(pdfData) == 0 {
		t.Error("GenerateOnly() returned empty PDF data")
	}
}

func TestService_GenerateAndUpload(t *testing.T) {
	// Initialize service
	svc := New(Config{
		UploadBaseURL: "https://example.com",
		BearerToken:   "test-token",
	})

	// Test data
	data := map[string]interface{}{
		"title": "Test Document",
		"items": []interface{}{
			map[string]interface{}{
				"name":  "Item 1",
				"value": 100,
			},
		},
	}

	// Test config
	config := UploadConfig{
		OrganizationalID: "test-org",
		BranchID:         "test-branch",
		CreatedBy:        "test-user",
		FileName:         "test.pdf",
	}

	// Test PDF generation and upload
	ctx := context.Background()
	response, err := svc.GenerateAndUpload(ctx, data, config)
	if err != nil {
		t.Errorf("GenerateAndUpload() error = %v", err)
		return
	}

	// Verify response
	if response.FileName != "test.pdf" {
		t.Errorf("GenerateAndUpload() filename = %v, want %v", response.FileName, "test.pdf")
	}
}

// Test configuration validation
func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				UploadBaseURL: "https://example.com",
				BearerToken:   "test-token",
			},
			wantErr: false,
		},
		{
			name: "missing URL",
			config: Config{
				BearerToken: "test-token",
			},
			wantErr: true,
		},
		{
			name: "missing token",
			config: Config{
				UploadBaseURL: "https://example.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.config.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
