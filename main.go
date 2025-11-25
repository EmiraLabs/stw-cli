package main

import (
	"flag"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Page struct {
	Title   string
	Content template.HTML
}

var buildFlag = flag.Bool("build", false, "Build the site")
var serveFlag = flag.Bool("serve", false, "Build and serve the site")

func main() {
	flag.Parse()
	if *buildFlag {
		build()
	} else if *serveFlag {
		serve()
	} else {
		flag.Usage()
	}
}

func build() {
	os.RemoveAll("dist")
	os.MkdirAll("dist", 0755)

	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/components/header.html", "templates/components/footer.html"))

	filepath.WalkDir("pages", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == "index.html" {
			rel, _ := filepath.Rel("pages", path)
			dst := filepath.Join("dist", rel)
			os.MkdirAll(filepath.Dir(dst), 0755)

			var title string
			if rel == "index.html" {
				title = "Home"
			} else {
				dir := filepath.Dir(rel)
				title = strings.Title(filepath.Base(dir))
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			page := Page{Title: title, Content: template.HTML(content)}
			f, err := os.Create(dst)
			if err != nil {
				return err
			}
			defer f.Close()
			return tmpl.ExecuteTemplate(f, "base.html", page)
		}
		return nil
	})

	// Copy static files
	copyDir("assets", filepath.Join("dist", "assets"))
}

func serve() {
	build()
	fs := http.FileServer(http.Dir("./dist"))
	log.Println("Serving ./dist on http://localhost:8001")
	log.Fatal(http.ListenAndServe(":8001", fs))
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	os.MkdirAll(filepath.Dir(dst), 0755)
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	return err
}
