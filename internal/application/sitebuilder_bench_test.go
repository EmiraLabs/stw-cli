package application

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/EmiraLabs/stw-cli/internal/domain"
	"github.com/EmiraLabs/stw-cli/internal/meta"
)

// BenchmarkBuild measures the performance of the site build process
func BenchmarkBuild(b *testing.B) {
	// Create temp dir with test site
	tmpDir := b.TempDir()

	// Create pages/index.html
	pagesDir := filepath.Join(tmpDir, "pages")
	os.MkdirAll(pagesDir, 0755)
	indexPath := filepath.Join(pagesDir, "index.html")
	content := `---
title: Benchmark Test
description: Testing build performance
---
<h1>{{.Title}}</h1>
<p>Some content for benchmarking</p>`
	os.WriteFile(indexPath, []byte(content), 0644)

	// Create templates
	templatesDir := filepath.Join(tmpDir, "templates")
	os.MkdirAll(filepath.Join(templatesDir, "components"), 0755)
	os.MkdirAll(filepath.Join(templatesDir, "partials"), 0755)

	basePath := filepath.Join(templatesDir, "base.html")
	baseContent := `<!DOCTYPE html>
<html>
<head><title>{{.Title}}</title></head>
<body>
<header><h1>{{.Title}}</h1></header>
{{.Content}}
<footer>Footer</footer>
</body>
</html>`
	os.WriteFile(basePath, []byte(baseContent), 0644)

	headerPath := filepath.Join(templatesDir, "components", "header.html")
	os.WriteFile(headerPath, []byte(""), 0644)

	footerPath := filepath.Join(templatesDir, "components", "footer.html")
	os.WriteFile(footerPath, []byte(""), 0644)

	headPath := filepath.Join(templatesDir, "partials", "head.html")
	os.WriteFile(headPath, []byte(""), 0644)

	// Create assets
	assetsDir := filepath.Join(tmpDir, "assets")
	os.MkdirAll(assetsDir, 0755)
	cssDir := filepath.Join(assetsDir, "css")
	os.MkdirAll(cssDir, 0755)
	cssPath := filepath.Join(cssDir, "styles.css")
	os.WriteFile(cssPath, []byte("body { margin: 0; }"), 0644)

	// Change to temp dir
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	site := &domain.Site{
		PagesDir:         "pages",
		TemplatesDir:     "templates",
		AssetsDir:        "assets",
		DistDir:          "dist",
		EnableAutoReload: false,
		Config:           map[string]interface{}{"test": "value"},
	}

	fs := NewMockFileSystem()
	fs.hasPages = true
	fs.hasAssets = true
	renderer := NewMockTemplateRenderer()

	builder := NewSiteBuilder(site, fs, renderer)

	// Reset timer after setup
	b.ResetTimer()

	// Run the build b.N times
	for i := 0; i < b.N; i++ {
		builder.Build()
	}
}

// BenchmarkBuildPages measures just the page building performance
func BenchmarkBuildPages(b *testing.B) {
	site := &domain.Site{
		DistDir:  "dist",
		PagesDir: "pages",
		Config:   map[string]interface{}{"test": "value"},
	}

	fs := NewMockFileSystem()
	fs.hasPages = true
	renderer := NewMockTemplateRenderer()
	renderer.ParseFiles(filepath.Join("templates", domain.BaseTemplate))

	builder := &SiteBuilder{site: site, fs: fs, renderer: renderer}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		builder.buildPages(meta.Meta{})
	}
}

// BenchmarkCopyAssets measures asset copying performance
func BenchmarkCopyAssets(b *testing.B) {
	site := &domain.Site{DistDir: "dist", AssetsDir: "assets"}
	fs := NewMockFileSystem()
	fs.hasAssets = true

	builder := &SiteBuilder{site: site, fs: fs}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		builder.copyAssets()
	}
}
