package drawing

import (
	"github.com/kvnxiao/pictorio/model"
)

type History interface {
	Append(line model.Line) bool
	GetAll() []model.Line
	Redo() bool
	Undo() bool
	Clear() bool
}

type Drawing struct {
	lines     []model.Line
	redoStack []model.Line
}

func NewDrawingHistory() History {
	return &Drawing{
		lines:     nil,
		redoStack: nil,
	}
}

func (d *Drawing) Append(line model.Line) bool {
	d.lines = append(d.lines, line)
	return true
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

	return true
}
