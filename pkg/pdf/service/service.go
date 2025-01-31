// Package service provides PDF generation and upload functionality
package service

import (
	"bytes"
	"context"
	"fmt"

	"github.com/josephmojoo/pdfgen/pkg/pdf/generator/dynamic"
)

// Service defines the PDF service interface
type Service interface {
	GenerateAndUpload(ctx context.Context, data map[string]interface{}, config UploadConfig) (*UploadResponse, error)
	GenerateOnly(ctx context.Context, data map[string]interface{}) ([]byte, error)
}

type service struct {
	generator *dynamic.Generator
	uploader  Uploader
}

// New creates a new PDF service
func New(config Config) Service {
	return &service{
		generator: dynamic.NewGenerator(),
		uploader:  newUploader(config),
	}
}

// GenerateAndUpload generates a PDF and uploads it
func (s *service) GenerateAndUpload(ctx context.Context, data map[string]interface{}, config UploadConfig) (*UploadResponse, error) {
	// Generate PDF
	pdfData, err := s.GenerateOnly(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Upload PDF
	response, err := s.uploader.Upload(ctx, pdfData, config)
	if err != nil {
		return nil, fmt.Errorf("failed to upload PDF: %w", err)
	}

	return response, nil
}

// GenerateOnly generates a PDF without uploading
func (s *service) GenerateOnly(ctx context.Context, data map[string]interface{}) ([]byte, error) {
	// Create template from data
	template := s.generator.GenerateTemplate(data)

	// Create a new generator instance
	gen := dynamic.NewGenerator()

	// Generate PDF
	var buf bytes.Buffer
	if err := gen.Generate(ctx, &buf, template); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}
