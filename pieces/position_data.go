package pieces

// PositionData is used by modules to define and store the coordinates
// and other positional data that TView needs to render a view to screen
type PositionData struct {
	Row       int
	Col       int
	RowSpan   int
	ColSpan   int
	MinHeight int
	MinWidth  int
}

/* -------------------- Exported Functions -------------------- */

func (p *PositionData) GetRow() int {
	return p.Row
}

func (p *PositionData) GetCol() int {
	return p.Col
}

func (p *PositionData) GetRowSpan() int {
	return p.RowSpan
}

func (p *PositionData) GetColSpan() int {
	return p.ColSpan
}

func (p *PositionData) GetMinWidth() int {
	return p.MinWidth
}

func (p *PositionData) GetMinHeight() int {
	return p.MinHeight
}
