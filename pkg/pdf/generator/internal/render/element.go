package render

import (
	"fmt"
	"strings"

	"github.com/josephmojoo/pdfgen/pkg/pdf/model"
	"github.com/jung-kurt/gofpdf"
)

// Context holds the rendering context for PDF elements
type Context struct {
	PDF      *gofpdf.Fpdf
	PageSize model.Size
	Margins  model.Padding
}

// ElementRenderer defines the interface for rendering PDF elements
type ElementRenderer interface {
	Render(ctx *Context, element model.Element) error
}

// TextRenderer handles rendering of text elements
type TextRenderer struct{}

func (r *TextRenderer) Render(ctx *Context, element model.Element) error {
	content, ok := element.Content.(string)
	if !ok {
		return fmt.Errorf("invalid content type for text element")
	}

	pdf := ctx.PDF
	style := element.Style
	if style != nil {
		if style.FontFamily != "" {
			pdf.SetFont(style.FontFamily, "", style.FontSize)
		}
		if style.FontColor != "" {
			pdf.SetTextColor(0, 0, 0) // TODO: Parse color string
		}
	}

	// Handle multiline text
	lines := strings.Split(content, "\n")
	lineHeight := pdf.PointToUnitConvert(style.FontSize)

	for i, line := range lines {
		// Calculate text width for positioning
		textWidth := pdf.GetStringWidth(line)
		textX := element.Bounds.X

		// Handle text alignment
		if style != nil && style.Alignment != "" {
			switch style.Alignment {
			case model.AlignLeft:
				textX = element.Bounds.X
			case model.AlignCenter:
				textX = element.Bounds.X + (element.Bounds.Width-textWidth)/2
			case model.AlignRight:
				textX = element.Bounds.X + element.Bounds.Width - textWidth
			case model.AlignJustify:
				// TODO: Implement text justification
				textX = element.Bounds.X
			}
		} else {
			// Default to left alignment
			textX = element.Bounds.X
		}

		// Calculate Y position for each line
		textY := element.Bounds.Y + lineHeight + float64(i)*lineHeight*1.5

		pdf.Text(textX, textY, line)
	}

	return nil
}

// TableRenderer handles rendering of table elements
type TableRenderer struct{}

func (r *TableRenderer) Render(ctx *Context, element model.Element) error {
	// First, try to convert the content to []interface{}
	rawContent, ok := element.Content.([]interface{})
	if !ok {
		return fmt.Errorf("invalid content type for table element: expected []interface{}, got %T", element.Content)
	}

	// Convert the raw content to [][]string
	content := make([][]string, len(rawContent))
	for i, row := range rawContent {
		rowArray, ok := row.([]interface{})
		if !ok {
			return fmt.Errorf("invalid row type at index %d: expected []interface{}, got %T", i, row)
		}

		content[i] = make([]string, len(rowArray))
		for j, cell := range rowArray {
			// Convert each cell to string
			switch v := cell.(type) {
			case string:
				content[i][j] = v
			case float64:
				content[i][j] = fmt.Sprintf("%.2f", v)
			case int:
				content[i][j] = fmt.Sprintf("%d", v)
			case bool:
				content[i][j] = fmt.Sprintf("%v", v)
			default:
				content[i][j] = fmt.Sprintf("%v", v)
			}
		}
	}

	pdf := ctx.PDF
	x, y := element.Bounds.X, element.Bounds.Y

	// Apply table styles
	fontSize := 10.0 // Default font size
	if element.Style != nil {
		if element.Style.FontFamily != "" {
			pdf.SetFont(element.Style.FontFamily, "", element.Style.FontSize)
			fontSize = element.Style.FontSize
		}
	}

	// Calculate cell dimensions
	colCount := len(content[0])
	cellWidth := element.Bounds.Width / float64(colCount)
	cellHeight := pdf.PointToUnitConvert(fontSize) * 2

	// Draw table
	for rowIndex, row := range content {
		for colIndex, cell := range row {
			currentX := x + float64(colIndex)*cellWidth
			currentY := y + float64(rowIndex)*cellHeight

			// Draw cell border
			if element.Style != nil && element.Style.Border != nil {
				pdf.SetLineWidth(element.Style.Border.Width)
				pdf.Rect(currentX, currentY, cellWidth, cellHeight, "D")
			}

			// Calculate text position based on alignment
			textWidth := pdf.GetStringWidth(cell)
			textX := currentX

			// Handle cell alignment
			if element.Style != nil && element.Style.Alignment != "" {
				switch element.Style.Alignment {
				case model.AlignLeft:
					textX = currentX + 2 // Small padding for left alignment
				case model.AlignCenter:
					textX = currentX + (cellWidth-textWidth)/2
				case model.AlignRight:
					textX = currentX + cellWidth - textWidth - 2 // Small padding for right alignment
				}
			} else {
				// Default to center alignment for tables
				textX = currentX + (cellWidth-textWidth)/2
			}

			// Adjust Y position to center text vertically
			fontHeight := pdf.PointToUnitConvert(fontSize)
			textY := currentY + (cellHeight-fontHeight)/2 + fontHeight

			// Draw cell content
			pdf.Text(textX, textY, cell)
		}
	}

	return nil
}

// ImageRenderer handles rendering of image elements
type ImageRenderer struct{}

func (r *ImageRenderer) Render(ctx *Context, element model.Element) error {
	content, ok := element.Content.(string) // Assuming content is image path
	if !ok {
		return fmt.Errorf("invalid content type for image element")
	}

	pdf := ctx.PDF
	pdf.Image(content, element.Bounds.X, element.Bounds.Y, element.Bounds.Width, element.Bounds.Height, false, "", 0, "")
	return nil
}

// Registry maps element types to their renderers
type Registry struct {
	renderers map[model.ElementType]ElementRenderer
}

// NewRegistry creates a new renderer registry with default renderers
func NewRegistry() *Registry {
	r := &Registry{
		renderers: make(map[model.ElementType]ElementRenderer),
	}

	// Register default renderers
	r.renderers[model.ElementTypeText] = &TextRenderer{}
	r.renderers[model.ElementTypeTable] = &TableRenderer{}
	r.renderers[model.ElementTypeImage] = &ImageRenderer{}

	return r
}

// GetRenderer returns the appropriate renderer for an element type
func (r *Registry) GetRenderer(elementType model.ElementType) (ElementRenderer, error) {
	renderer, ok := r.renderers[elementType]
	if !ok {
		return nil, fmt.Errorf("no renderer found for element type: %s", elementType)
	}
	return renderer, nil
}

// RegisterRenderer adds a custom renderer for an element type
func (r *Registry) RegisterRenderer(elementType model.ElementType, renderer ElementRenderer) {
	r.renderers[elementType] = renderer
}
