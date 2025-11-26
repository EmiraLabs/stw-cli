package domain

// Site represents the static site configuration
type Site struct {
	PagesDir     string
	TemplatesDir string
	AssetsDir    string
	DistDir      string
}
