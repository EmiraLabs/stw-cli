package main

import (
	"flag"
	"log"

	"github.com/EmiraLabs/stw-cli/internal/application"
	"github.com/EmiraLabs/stw-cli/internal/domain"
	"github.com/EmiraLabs/stw-cli/internal/infrastructure"
)

var buildFlag = flag.Bool("build", false, "Build the site")
var serveFlag = flag.Bool("serve", false, "Build and serve the site")

func main() {
	flag.Parse()

	site := &domain.Site{
		PagesDir:     "pages",
		TemplatesDir: "templates",
		AssetsDir:    "assets",
		DistDir:      "dist",
	}

	fs := &infrastructure.OSFileSystem{}
	renderer := &infrastructure.GoTemplateRenderer{}

	builder := application.NewSiteBuilder(site, fs, renderer)

	if *buildFlag {
		if err := builder.Build(); err != nil {
			log.Fatal(err)
		}
	} else if *serveFlag {
		server := application.NewSiteServer(site, builder)
		if err := server.Serve(); err != nil {
			log.Fatal(err)
		}
	} else {
		flag.Usage()
	}
}
