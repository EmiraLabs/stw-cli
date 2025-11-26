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
				PagesDir:         "pages",
				TemplatesDir:     "templates",
				AssetsDir:        "assets",
				DistDir:          "dist",
				EnableAutoReload: false,
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
			watch, _ := cmd.Flags().GetBool("watch")

			site := &domain.Site{
				PagesDir:         "pages",
				TemplatesDir:     "templates",
				AssetsDir:        "assets",
				DistDir:          "dist",
				EnableAutoReload: watch,
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

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize the project with optional features",
		Run: func(cmd *cobra.Command, args []string) {
			tailwind, _ := cmd.Flags().GetBool("tailwind")

			fs := &infrastructure.OSFileSystem{}

			initializer := application.NewSiteInitializer(fs)

			if tailwind {
				if err := initializer.InitTailwind(); err != nil {
					log.Fatal(err)
				}
			}
		},
	}

	initCmd.Flags().Bool("tailwind", false, "Initialize with Tailwind CSS setup")

	serveCmd.Flags().StringP("port", "p", "8080", "Port to serve on")
	serveCmd.Flags().BoolP("watch", "w", true, "Enable auto-reload on file changes")

	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(initCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
