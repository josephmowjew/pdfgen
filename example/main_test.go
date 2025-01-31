package main

import (
	"context"
	"os"
	"testing"

	"github.com/josephmojoo/pdfgen/pkg/pdf/service"
)

func TestPDFGeneration(t *testing.T) {
	// Test data
	data := map[string]interface{}{
		"customer": map[string]interface{}{
			"name":    "John Doe",
			"email":   "john@example.com",
			"address": "123 Main St, City, Country",
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

	// Initialize service
	svc := service.New(service.Config{})

	// Test PDF generation
	t.Run("Generate PDF", func(t *testing.T) {
		pdfData, err := svc.GenerateOnly(context.Background(), data)
		if err != nil {
			t.Fatalf("Failed to generate PDF: %v", err)
		}
		if len(pdfData) == 0 {
			t.Error("Generated PDF is empty")
		}

		// Save to test file
		if err := os.WriteFile("test_output.pdf", pdfData, 0644); err != nil {
			t.Fatalf("Failed to save PDF: %v", err)
		}
		defer os.Remove("test_output.pdf")
	})
}
