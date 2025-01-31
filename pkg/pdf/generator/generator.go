package generator

import (
	"bytes"
	"context"
	"fmt"

	"github.com/josephmojoo/pdfgen/pkg/pdf/generator/internal/layout"
	"github.com/josephmojoo/pdfgen/pkg/pdf/generator/internal/render"
	"github.com/josephmojoo/pdfgen/pkg/pdf/model"
	"github.com/jung-kurt/gofpdf"
)

// Generator handles PDF generation from templates
type Generator struct {
	template *model.Template
	layout   *layout.Manager
	registry *render.Registry
	margins  model.Padding
}

// New creates a new PDF generator
func New(template *model.Template) *Generator {
	margins := model.Padding{
		Top:    10,
		Right:  10,
		Bottom: 10,
		Left:   10,
	}

	return &Generator{
		template: template,
		layout:   layout.NewManager(template.Size, margins),
		registry: render.NewRegistry(),
		margins:  margins,
	}
}

// Generate creates a PDF document from the template and data
func (g *Generator) Generate(ctx context.Context, data interface{}) (*bytes.Buffer, error) {
	// Validate template and data
	if err := g.template.Validate(); err != nil {
		return nil, fmt.Errorf("invalid template: %w", err)
	}
	if err := g.template.ValidateData(data); err != nil {
		return nil, fmt.Errorf("invalid data: %w", err)
	}

	// Calculate layout
	if err := g.layout.CalculateLayout(g.template.Elements); err != nil {
		return nil, fmt.Errorf("layout calculation failed: %w", err)
	}

	// Create PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	renderCtx := &render.Context{
		PDF: pdf,
		PageSize: model.Size{
			Width:  g.template.Size.Width,
			Height: g.template.Size.Height,
		},
		Margins: g.margins,
	}

	// Render each page
	totalPages := g.layout.TotalPages()
	for page := 1; page <= totalPages; page++ {
		pdf.AddPage()

		// Render elements for current page
		elements := g.layout.GetPageElements(page)
		for _, element := range elements {
			renderer, err := g.registry.GetRenderer(element.Type)
			if err != nil {
				return nil, fmt.Errorf("failed to get renderer: %w", err)
			}

			if err := renderer.Render(renderCtx, element); err != nil {
				return nil, fmt.Errorf("failed to render element: %w", err)
			}
		}
	}

	// Write to buffer
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to write PDF: %w", err)
	}

	return &buf, nil
}

// RegisterRenderer registers a custom renderer for an element type
func (g *Generator) RegisterRenderer(elementType model.ElementType, renderer render.ElementRenderer) {
	g.registry.RegisterRenderer(elementType, renderer)
}

// SetMargins sets the page margins
func (g *Generator) SetMargins(margins model.Padding) {
	g.margins = margins
	g.layout = layout.NewManager(g.template.Size, margins)
}
