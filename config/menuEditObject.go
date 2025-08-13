package config

type MenuEditPageData struct {
	Title      string
	Action     string
	Type       string
	Item       Item
	Category   Category
	Categories []Category
	Error      string
}
