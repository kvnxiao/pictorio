package drawing

import (
	"github.com/kvnxiao/pictorio/model"
)

type History interface {
	Append(line model.Line)
	GetAll() []model.Line
	Redo()
	Undo()
	Clear()
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

func (d *Drawing) Append(line model.Line) {
	d.lines = append(d.lines, line)
}

func (d *Drawing) GetAll() []model.Line {
	allLines := make([]model.Line, len(d.lines))
	copy(allLines, d.lines)
	return allLines
}

func (d *Drawing) Redo() {
	var line model.Line
	// Pop most recent line added to redo stack
	line, d.redoStack = d.redoStack[len(d.redoStack)-1], d.redoStack[:len(d.redoStack)-1]
	// Push into drawing history
	d.lines = append(d.lines, line)
}

func (d *Drawing) Undo() {
	var line model.Line
	// Pop most recent line added to drawing
	line, d.lines = d.lines[len(d.lines)-1], d.lines[:len(d.lines)-1]
	// Push into redo stack
	d.redoStack = append(d.redoStack, line)
}

func (d *Drawing) Clear() {
	d.lines = nil
	d.redoStack = nil
}
