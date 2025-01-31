package model

import (
	"encoding/json"

	"github.com/josephmojoo/pdfgen/pkg/pdf/errors"
)

// Position represents x,y coordinates in the PDF
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Size represents width and height
type Size struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// Bounds combines position and size
type Bounds struct {
	Position
	Size
}

// ElementType defines the type of PDF element
type ElementType string

const (
	ElementTypeText    ElementType = "text"
	ElementTypeTable   ElementType = "table"
	ElementTypeImage   ElementType = "image"
	ElementTypeBarcode ElementType = "barcode"
	ElementTypeForm    ElementType = "form"
)

// TextAlignment defines text alignment options
type TextAlignment string

const (
	AlignLeft    TextAlignment = "left"
	AlignCenter  TextAlignment = "center"
	AlignRight   TextAlignment = "right"
	AlignJustify TextAlignment = "justify"
)

// Element represents a PDF element configuration
type Element struct {
	ID       string          `json:"id"`
	Type     ElementType     `json:"type"`
	Bounds   Bounds          `json:"bounds"`
	Content  interface{}     `json:"content"`
	Style    *Style          `json:"style,omitempty"`
	Metadata json.RawMessage `json:"metadata,omitempty"`
}

// Style defines the visual properties of an element
type Style struct {
	FontFamily string        `json:"fontFamily,omitempty"`
	FontSize   float64       `json:"fontSize,omitempty"`
	FontColor  string        `json:"fontColor,omitempty"`
	Background string        `json:"background,omitempty"`
	Border     *Border       `json:"border,omitempty"`
	Padding    *Padding      `json:"padding,omitempty"`
	Alignment  TextAlignment `json:"alignment,omitempty"`
}

// Border defines border properties
type Border struct {
	Width float64 `json:"width"`
	Color string  `json:"color"`
	Style string  `json:"style"`
}

// Padding defines padding properties
type Padding struct {
	Top    float64 `json:"top"`
	Right  float64 `json:"right"`
	Bottom float64 `json:"bottom"`
	Left   float64 `json:"left"`
}

// Template defines the structure of a PDF template
type Template struct {
	Name     string                 `json:"name"`
	Version  string                 `json:"version"`
	Size     Size                   `json:"size"`
	Elements []Element              `json:"elements"`
	Schema   map[string]interface{} `json:"schema"`
}

// Validate ensures the template configuration is valid
func (t *Template) Validate() error {
	if t.Name == "" {
		return errors.NewPDFError(errors.ErrInvalidTemplate, "template name is required", nil)
	}
	if t.Size.Width <= 0 || t.Size.Height <= 0 {
		return errors.NewPDFError(errors.ErrInvalidTemplate, "invalid template size", nil)
	}
	if len(t.Elements) == 0 {
		return errors.NewPDFError(errors.ErrInvalidTemplate, "template must contain at least one element", nil)
	}
	return nil
}

// ValidateData ensures the provided data matches the template schema
func (t *Template) ValidateData(data interface{}) error {
	if t.Schema == nil {
		return nil
	}
	// TODO: Implement schema validation
	return nil
}
