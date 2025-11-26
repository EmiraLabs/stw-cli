package application

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/EmiraLabs/stw-cli/internal/domain"
	"github.com/EmiraLabs/stw-cli/internal/infrastructure"
)

// SiteBuilder handles building the static site
type SiteBuilder struct {
	site     *domain.Site
	fs       infrastructure.FileSystem
	renderer infrastructure.TemplateRenderer
}

// NewSiteBuilder creates a new SiteBuilder
func NewSiteBuilder(site *domain.Site, fs infrastructure.FileSystem, renderer infrastructure.TemplateRenderer) *SiteBuilder {
	return &SiteBuilder{
		site:     site,
		fs:       fs,
		renderer: renderer,
	}
}

// Build builds the site
func (sb *SiteBuilder) Build() error {
	// Remove and recreate dist dir
	if err := sb.fs.RemoveAll(sb.site.DistDir); err != nil {
		return err
	}
	if err := sb.fs.MkdirAll(sb.site.DistDir, 0755); err != nil {
		return err
	}

	// Parse templates
	tmpl, err := sb.renderer.ParseFiles(
		filepath.Join(sb.site.TemplatesDir, domain.BaseTemplate),
		filepath.Join(sb.site.TemplatesDir, domain.HeaderTemplateFile),
		filepath.Join(sb.site.TemplatesDir, domain.FooterTemplateFile),
	)
	if err != nil {
		return err
	}
	// Set the template in renderer if possible, but since interface, perhaps cast or change.

	// For simplicity, use the tmpl directly
	sb.buildPages(tmpl)
	sb.copyAssets()

	return nil
}

func (sb *SiteBuilder) buildPages(tmpl *template.Template) error {
	return sb.fs.WalkDir(sb.site.PagesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == domain.IndexFile {
			rel, _ := filepath.Rel(sb.site.PagesDir, path)
			dst := filepath.Join(sb.site.DistDir, rel)
			if err := sb.fs.MkdirAll(filepath.Dir(dst), 0755); err != nil {
				return err
			}

			var title string
			if rel == domain.IndexFile {
				title = "Home"
			} else {
				dir := filepath.Dir(rel)
				title = strings.Title(filepath.Base(dir))
			}

			content, err := sb.fs.ReadFile(path)
			if err != nil {
				return err
			}
			page := domain.Page{Title: title, Content: template.HTML(content), Path: rel, IsDev: sb.site.EnableAutoReload}

			f, err := sb.fs.Create(dst)
			if err != nil {
				return err
			}
			defer f.Close()
			return tmpl.ExecuteTemplate(f, domain.BaseTemplate, page)
		}
		return nil
	})
}

func (sb *SiteBuilder) copyAssets() error {
	src := sb.site.AssetsDir
	dst := filepath.Join(sb.site.DistDir, "assets")
	return sb.fs.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return sb.fs.MkdirAll(target, 0755)
		}
		if strings.HasSuffix(path, "styles.css") {
			return sb.processCSS(path, target)
		}
		return sb.copyFile(path, target)
	})
}

func (sb *SiteBuilder) copyFile(src, dst string) error {
	content, err := sb.fs.ReadFile(src)
	if err != nil {
		return err
	}
	f, err := sb.fs.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(content)
	return err
}

func (sb *SiteBuilder) processCSS(src, dst string) error {
	// Check if postcss.config.js exists
	if _, err := sb.fs.Stat("postcss.config.js"); err != nil {
		// If not, just copy
		return sb.copyFile(src, dst)
	}
	// Ensure target dir exists
	if err := sb.fs.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	// Run postcss
	fmt.Println("Processing CSS with PostCSS...")
	cmd := exec.Command("./node_modules/.bin/postcss", src, "-o", dst)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
