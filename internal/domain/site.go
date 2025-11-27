package domain

const (
	IndexFile          = "index.html"
	BaseTemplate       = "base.html"
	HeaderTemplateFile = "components/header.html"
	FooterTemplateFile = "components/footer.html"
	HeadTemplateFile   = "partials/head.html"
)

// Site represents the static site configuration
type Site struct {
	PagesDir         string
	TemplatesDir     string
	AssetsDir        string
	DistDir          string
	EnableAutoReload bool
	Config           map[string]interface{}
	ConfigPath       string
}
