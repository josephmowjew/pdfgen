package dynamic

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/josephmojoo/pdfgen/pkg/pdf/model"
	"github.com/jung-kurt/gofpdf"
)

// Generator handles dynamic PDF generation based on data
type Generator struct {
	lineHeight float64
	margins    model.Padding
	fontSize   float64
}

// NewGenerator creates a new dynamic generator
func NewGenerator() *Generator {
	return &Generator{
		lineHeight: 8.0,
		margins: model.Padding{
			Top:    20,
			Right:  20,
			Bottom: 20,
			Left:   20,
		},
		fontSize: 12,
	}
}

// GenerateTemplate creates a template from the provided data
func (g *Generator) GenerateTemplate(data interface{}) *model.Template {
	elements := []model.Element{}
	currentY := g.margins.Top

	// Add title
	elements = append(elements, model.Element{
		ID:   "title",
		Type: "text",
		Bounds: model.Bounds{
			Position: model.Position{
				X: g.margins.Left,
				Y: currentY,
			},
			Size: model.Size{
				Width:  210 - g.margins.Left - g.margins.Right,
				Height: g.lineHeight * 2,
			},
		},
		Content: "Data Report",
		Style: &model.Style{
			FontFamily: "Arial",
			FontSize:   24,
			FontColor:  "#000000",
			Alignment:  "center",
		},
	})
	currentY += g.lineHeight * 3

	// Process data recursively
	elements = append(elements, g.processData(data, "", &currentY)...)

	// Create template
	return &model.Template{
		Name:    "Dynamic Template",
		Version: "1.0",
		Size: model.Size{
			Width:  210, // A4 width in mm
			Height: 297, // A4 height in mm
		},
		Elements: elements,
	}
}

// processData handles different types of data and creates appropriate elements
func (g *Generator) processData(data interface{}, prefix string, currentY *float64) []model.Element {
	elements := []model.Element{}

	switch v := data.(type) {
	case map[string]interface{}:
		if prefix != "" {
			elements = append(elements, g.createSectionHeader(prefix, currentY))
			*currentY += g.lineHeight * 1.5
		}

		for key, value := range v {
			fieldName := key
			if prefix != "" {
				fieldName = fmt.Sprintf("%s.%s", prefix, key)
			}
			elements = append(elements, g.processData(value, fieldName, currentY)...)
		}

	case []interface{}:
		elements = append(elements, g.createArrayElement(prefix, v, currentY))
		*currentY += g.lineHeight*float64(len(v)) + g.lineHeight

	default:
		elements = append(elements, g.createValueElement(prefix, v, currentY))
		*currentY += g.lineHeight
	}

	return elements
}

// createSectionHeader creates a header element for sections
func (g *Generator) createSectionHeader(title string, currentY *float64) model.Element {
	return model.Element{
		ID:   fmt.Sprintf("header-%s", title),
		Type: "text",
		Bounds: model.Bounds{
			Position: model.Position{
				X: g.margins.Left,
				Y: *currentY,
			},
			Size: model.Size{
				Width:  210 - g.margins.Left - g.margins.Right,
				Height: g.lineHeight * 1.5,
			},
		},
		Content: strings.ToUpper(title),
		Style: &model.Style{
			FontFamily: "Arial",
			FontSize:   14,
			FontColor:  "#333333",
			Background: "#f5f5f5",
			Padding:    &model.Padding{Left: 5, Top: 2, Bottom: 2, Right: 5},
		},
	}
}

// createValueElement creates an element for a key-value pair
func (g *Generator) createValueElement(key string, value interface{}, currentY *float64) model.Element {
	formattedValue := g.formatValue(value)
	return model.Element{
		ID:   fmt.Sprintf("field-%s", key),
		Type: "text",
		Bounds: model.Bounds{
			Position: model.Position{
				X: g.margins.Left,
				Y: *currentY,
			},
			Size: model.Size{
				Width:  210 - g.margins.Left - g.margins.Right,
				Height: g.lineHeight,
			},
		},
		Content: fmt.Sprintf("%s: %s", key, formattedValue),
		Style: &model.Style{
			FontFamily: "Arial",
			FontSize:   g.fontSize,
			FontColor:  "#000000",
		},
	}
}

// createArrayElement creates an element for an array
func (g *Generator) createArrayElement(key string, items []interface{}, currentY *float64) model.Element {
	var content strings.Builder
	content.WriteString(fmt.Sprintf("%s:\n", key))

	for _, item := range items {
		switch v := item.(type) {
		case map[string]interface{}:
			// Format map items in a readable way
			content.WriteString("-")
			var parts []string
			for k, val := range v {
				parts = append(parts, fmt.Sprintf("%s: %v", k, val))
			}
			// Join all parts with commas
			content.WriteString(" " + strings.Join(parts, ", ") + "\n")
		default:
			content.WriteString(fmt.Sprintf("- %s\n", g.formatValue(item)))
		}
	}

	return model.Element{
		ID:   fmt.Sprintf("array-%s", key),
		Type: "text",
		Bounds: model.Bounds{
			Position: model.Position{
				X: g.margins.Left,
				Y: *currentY,
			},
			Size: model.Size{
				Width:  210 - g.margins.Left - g.margins.Right,
				Height: g.lineHeight * float64(len(items)+1),
			},
		},
		Content: content.String(),
		Style: &model.Style{
			FontFamily: "Arial",
			FontSize:   g.fontSize,
			FontColor:  "#000000",
		},
	}
}

// formatValue converts a value to a formatted string
func (g *Generator) formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%.2f", v)
	case int:
		return fmt.Sprintf("%d", v)
	case bool:
		return fmt.Sprintf("%v", v)
	case nil:
		return "N/A"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Generate creates a PDF document from the template and writes it to the provided writer
func (g *Generator) Generate(ctx context.Context, w io.Writer, template *model.Template) error {
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Add page
	pdf.AddPage()

	// Render elements
	for _, element := range template.Elements {
		switch element.Type {
		case "text":
			if err := g.renderText(pdf, element); err != nil {
				return fmt.Errorf("failed to render text element: %w", err)
			}
		// Add other element types as needed
		default:
			return fmt.Errorf("unsupported element type: %s", element.Type)
		}
	}

	// Write to output
	return pdf.Output(w)
}

func (g *Generator) renderText(pdf *gofpdf.Fpdf, element model.Element) error {
	style := element.Style
	if style == nil {
		style = &model.Style{
			FontFamily: "Arial",
			FontSize:   12,
			FontColor:  "#000000",
		}
	}

	pdf.SetFont(style.FontFamily, "", style.FontSize)
	pdf.SetTextColor(0, 0, 0) // Parse style.FontColor and set RGB values

	content, ok := element.Content.(string)
	if !ok {
		return fmt.Errorf("invalid content type: expected string, got %T", element.Content)
	}

	x := element.Bounds.Position.X
	y := element.Bounds.Position.Y
	lineHeight := pdf.PointToUnitConvert(style.FontSize) * 1.2 // 1.2 for comfortable line spacing

	// Handle multiline text
	for _, line := range strings.Split(content, "\n") {
		pdf.Text(x, y, line)
		y += lineHeight
	}

	return nil
}
