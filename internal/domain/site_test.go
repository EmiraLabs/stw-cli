package domain

import "testing"

func TestSiteConstants(t *testing.T) {
	if IndexFile != "index.html" {
		t.Errorf("Expected IndexFile 'index.html', got %s", IndexFile)
	}

	if BaseTemplate != "base.html" {
		t.Errorf("Expected BaseTemplate 'base.html', got %s", BaseTemplate)
	}

	if HeaderTemplateFile != "components/header.html" {
		t.Errorf("Expected HeaderTemplateFile 'components/header.html', got %s", HeaderTemplateFile)
	}

	if FooterTemplateFile != "components/footer.html" {
		t.Errorf("Expected FooterTemplateFile 'components/footer.html', got %s", FooterTemplateFile)
	}

	if HeadTemplateFile != "partials/head.html" {
		t.Errorf("Expected HeadTemplateFile 'partials/head.html', got %s", HeadTemplateFile)
	}
}

func TestSiteStruct(t *testing.T) {
	s := Site{
		PagesDir:         "pages",
		TemplatesDir:     "templates",
		AssetsDir:        "assets",
		DistDir:          "dist",
		EnableAutoReload: true,
		Config:           map[string]interface{}{},
		ConfigPath:       "config.yaml",
	}

	if s.PagesDir != "pages" {
		t.Errorf("Expected PagesDir 'pages', got %s", s.PagesDir)
	}

	if !s.EnableAutoReload {
		t.Error("Expected EnableAutoReload true")
	}
}
