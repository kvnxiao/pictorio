package model

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Line struct {
	Points       []Point `json:"points"`
	ColourIdx    int     `json:"colourIdx"`
	ThicknessIdx int     `json:"thicknessIdx"`
}
