package application

import (
	"bytes"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/EmiraLabs/stw-cli/internal/domain"
	"github.com/EmiraLabs/stw-cli/internal/infrastructure"
	"github.com/EmiraLabs/stw-cli/internal/meta"
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

	// Load site meta
	siteMeta := meta.LoadSiteMeta(sb.site.Config)

	// Parse templates
	tmpl, err := sb.renderer.ParseFiles(
		filepath.Join(sb.site.TemplatesDir, domain.BaseTemplate),
		filepath.Join(sb.site.TemplatesDir, domain.HeaderTemplateFile),
		filepath.Join(sb.site.TemplatesDir, domain.FooterTemplateFile),
		filepath.Join(sb.site.TemplatesDir, domain.HeadTemplateFile),
	)
	if err != nil {
		return err
	}
	// Set the template in renderer if possible, but since interface, perhaps cast or change.

	// For simplicity, use the tmpl directly
	sb.buildPages(tmpl, siteMeta)
	sb.copyAssets()

	return nil
}

func (sb *SiteBuilder) buildPages(tmpl *template.Template, siteMeta meta.Meta) error {
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

			// Parse front matter
			pageMeta, body, err := meta.ParseFrontMatter(string(content))
			if err != nil {
				return err
			}

			// Merge meta
			mergedMeta := meta.Merge(siteMeta, pageMeta)

			// Validate meta
			if err := mergedMeta.Validate(sb.site.AssetsDir); err != nil {
				return err
			}

			// Parse page content as template
			pageTmpl, err := template.New("page").Parse(body)
			if err != nil {
				return err
			}

			// Create page data without Content
			pageData := domain.Page{
				Title:  title,
				Path:   rel,
				IsDev:  sb.site.EnableAutoReload,
				Config: sb.site.Config,
				Meta:   mergedMeta,
			}

			// Execute page template
			var buf bytes.Buffer
			if err := pageTmpl.Execute(&buf, pageData); err != nil {
				return err
			}

			page := domain.Page{
				Title:   title,
				Content: template.HTML(buf.String()),
				Path:    rel,
				IsDev:   sb.site.EnableAutoReload,
				Config:  sb.site.Config,
				Meta:    mergedMeta,
			}

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
