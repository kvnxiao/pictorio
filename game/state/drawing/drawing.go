package drawing

import (
	"github.com/kvnxiao/pictorio/model"
)

type History interface {
	Append(line model.Line) bool
	AppendFromTempLine(tempLine model.Line)
	PromoteLine() bool
	SetTempColour(colourIdx int)
	SetTempThickness(thicknessIdx int)
	GetAll() []model.Line
	Redo() bool
	Undo() bool
	Clear() bool
}

type Drawing struct {
	tempColour    int
	tempThickness int
	tempPoints    []model.Point
	lines         []model.Line
	redoStack     []model.Line
}

func NewDrawingHistory() History {
	return &Drawing{
		tempPoints: nil,
		lines:      nil,
		redoStack:  nil,
	}
}

func (d *Drawing) Append(line model.Line) bool {
	d.lines = append(d.lines, line)
	return true
}

func (d *Drawing) AppendFromTempLine(tempLine model.Line) {
	d.tempColour = tempLine.ColourIdx
	d.tempThickness = tempLine.ThicknessIdx
	d.tempPoints = append(d.tempPoints, tempLine.Points...)
}

func (d *Drawing) PromoteLine() bool {
	d.lines = append(d.lines, model.Line{
		Points:       d.tempPoints,
		ColourIdx:    d.tempColour,
		ThicknessIdx: d.tempThickness,
	})
	d.tempPoints = nil
	return true
}

func (d *Drawing) SetTempColour(colourIdx int) {
	d.tempColour = colourIdx
}

func (d *Drawing) SetTempThickness(thicknessIdx int) {
	d.tempThickness = thicknessIdx
}

func (d *Drawing) GetAll() []model.Line {
	allLines := make([]model.Line, len(d.lines))
	copy(allLines, d.lines)
	return allLines
}

func (d *Drawing) Redo() bool {
	// No-op if no lines to redo
	if len(d.redoStack) <= 0 {
		return false
	}

	var line model.Line
	// Pop most recent line added to redo stack
	line, d.redoStack = d.redoStack[len(d.redoStack)-1], d.redoStack[:len(d.redoStack)-1]
	// Push into drawing history
	d.lines = append(d.lines, line)

	return true
}

func (d *Drawing) Undo() bool {
	// No-op if no lines to undo
	if len(d.lines) <= 0 {
		return false
	}

	var line model.Line
	// Pop most recent line added to drawing
	line, d.lines = d.lines[len(d.lines)-1], d.lines[:len(d.lines)-1]
	// Push into redo stack
	d.redoStack = append(d.redoStack, line)

	return true
}

func (d *Drawing) Clear() bool {
	d.lines = nil
	d.redoStack = nil
	d.tempPoints = nil
	d.tempColour = 0
	d.tempThickness = 0

	return true
}
