package layout

import (
	"fmt"

	"github.com/josephmojoo/pdfgen/pkg/pdf/model"
)

// Manager handles the positioning and layout of PDF elements
type Manager struct {
	pageSize     model.Size
	margins      model.Padding
	currentPage  int
	currentY     float64
	elements     []model.Element
	pageElements map[int][]model.Element
}

// NewManager creates a new layout manager
func NewManager(pageSize model.Size, margins model.Padding) *Manager {
	return &Manager{
		pageSize:     pageSize,
		margins:      margins,
		currentPage:  1,
		currentY:     margins.Top,
		pageElements: make(map[int][]model.Element),
	}
}

// CalculateLayout positions all elements on pages
func (m *Manager) CalculateLayout(elements []model.Element) error {
	m.elements = elements

	for _, element := range elements {
		if err := m.positionElement(&element); err != nil {
			return fmt.Errorf("failed to position element: %w", err)
		}
	}

	return nil
}

// positionElement calculates the position for a single element
func (m *Manager) positionElement(element *model.Element) error {
	// Calculate available space
	availableHeight := m.pageSize.Height - m.currentY - m.margins.Bottom

	// Check if element fits on current page
	if element.Bounds.Height > availableHeight {
		m.startNewPage()
	}

	// Set element position
	element.Bounds.X = m.margins.Left
	element.Bounds.Y = m.currentY

	// Update current Y position
	m.currentY += element.Bounds.Height

	// Add element to current page
	m.pageElements[m.currentPage] = append(m.pageElements[m.currentPage], *element)

	return nil
}

// startNewPage begins a new page for element positioning
func (m *Manager) startNewPage() {
	m.currentPage++
	m.currentY = m.margins.Top
}

// GetPageElements returns all elements for a specific page
func (m *Manager) GetPageElements(page int) []model.Element {
	return m.pageElements[page]
}

// TotalPages returns the total number of pages
func (m *Manager) TotalPages() int {
	return m.currentPage
}
