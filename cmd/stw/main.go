package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/EmiraLabs/stw-cli/internal/application"
	"github.com/EmiraLabs/stw-cli/internal/domain"
	"github.com/EmiraLabs/stw-cli/internal/infrastructure"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "stw",
		Short: "Static Web Generator",
	}

	var buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Build the static site",
		Run: func(cmd *cobra.Command, args []string) {
			site := &domain.Site{
				PagesDir:     "pages",
				TemplatesDir: "templates",
				AssetsDir:    "assets",
				DistDir:      "dist",
			}

			fs := &infrastructure.OSFileSystem{}
			renderer := &infrastructure.GoTemplateRenderer{}

			builder := application.NewSiteBuilder(site, fs, renderer)

			if err := builder.Build(); err != nil {
				log.Fatal(err)
			}
		},
	}

	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Build and serve the static site",
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetString("port")

			site := &domain.Site{
				PagesDir:     "pages",
				TemplatesDir: "templates",
				AssetsDir:    "assets",
				DistDir:      "dist",
			}

			fs := &infrastructure.OSFileSystem{}
			renderer := &infrastructure.GoTemplateRenderer{}

			builder := application.NewSiteBuilder(site, fs, renderer)

			server := application.NewSiteServer(site, builder, port)
			if err := server.Serve(); err != nil {
				log.Fatal(err)
			}
		},
	}

	serveCmd.Flags().StringP("port", "p", "8080", "Port to serve on")

	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(serveCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
