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
