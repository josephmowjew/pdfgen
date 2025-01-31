package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/josephmojoo/pdfgen/pkg/pdf/service"
)

func main() {
	// Initialize the PDF service with actual credentials
	svc := service.New(service.Config{
		UploadBaseURL: "https://staging-storage.lyvepulse.com/storage/files",
		BearerToken:   "eyJhbGciOiJIUzM4NCJ9.eyJzdWIiOiJuby1yZXBseUBvcHNzaWZ5LmNvbSIsInVzZXJuYW1lIjoibm8tcmVwbHlAb3Bzc2lmeS5jb20iLCJlbXBsb3llZUlkIjoibm8tcmVwbHlAb3Bzc2lmeS5jb20iLCJmaXJzdE5hbWUiOiJTWVNURU0iLCJsYXN0TmFtZSI6IlNZU1RFTSIsInBob25lTnVtYmVyIjoiODc2NzUyMzQyIiwiZW5hYmxlZCI6dHJ1ZSwicGVuZGluZ1Jlc2V0IjpmYWxzZSwicm9sZXMiOlt7InJvbGVJZCI6IlNZU19BRE1JTiIsImJyYW5jaElkIjoiQlItMTAwMiIsIm9yZ2FuaXNhdGlvbmFsSWQiOiI1NDMyMSJ9XSwiaWF0IjoxNzM4MzA0OTQ0LCJleHAiOjE3MzgzMzM3NDR9.zK35mp8s8wbHFVOf7M6qpewdQvtfbiF1yo9B6hvfabnAgztsC4WbBZ7j4mY_dPyX",
	})

	// Sample data to be included in the PDF
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

	// Generate and save PDF to file
	pdfData, err := svc.GenerateOnly(context.Background(), data)
	if err != nil {
		log.Fatalf("Failed to generate PDF: %v", err)
	}

	// Save to file
	if err := os.WriteFile("output.pdf", pdfData, 0644); err != nil {
		log.Fatalf("Failed to save PDF: %v", err)
	}
	fmt.Println("PDF saved to output.pdf")

	// Generate and upload PDF with correct organizational details
	response, err := svc.GenerateAndUpload(context.Background(), data, service.UploadConfig{
		OrganizationalID: "54321",
		BranchID:         "BR-1002",
		CreatedBy:        "Joseph",
		FileName:         "order-report.pdf",
	})
	if err != nil {
		log.Fatalf("Failed to generate and upload PDF: %v", err)
	}

	fmt.Printf("PDF uploaded successfully:\n")
	fmt.Printf("File Name: %s\n", response.FileName)
	fmt.Printf("Download URL: %s\n", response.FileDownloadUri)
	fmt.Printf("File Type: %s\n", response.FileType)
	fmt.Printf("Size: %d bytes\n", response.Size)
}
