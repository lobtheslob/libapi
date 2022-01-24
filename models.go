package libapi

// Book struct (Model)
type Book struct {
	ID     int    `json:"id"`
	Name   string `json:"name,omitempty"`
	Author string `json:"author,omitempty"`
	ISBN   string `json:"isbn,omitempty"`
}
