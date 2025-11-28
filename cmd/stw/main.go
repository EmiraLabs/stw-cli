package main

import (
	"html/template"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/EmiraLabs/stw-cli/internal/application"
	"github.com/EmiraLabs/stw-cli/internal/domain"
	"github.com/EmiraLabs/stw-cli/internal/infrastructure"
)

func convertToHTML(data interface{}) interface{} {
	switch v := data.(type) {
	case string:
		return template.HTML(v)
	case map[string]interface{}:
		for key, val := range v {
			v[key] = convertToHTML(val)
		}
		return v
	case []interface{}:
		for i, val := range v {
			v[i] = convertToHTML(val)
		}
		return v
	default:
		return v
	}
}

func loadConfig() (map[string]interface{}, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]interface{}{}, nil // default empty config
		}
		return nil, err
	}
	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return convertToHTML(config).(map[string]interface{}), nil
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "stw",
		Short: "Static Web Generator",
	}

	var buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Build the static site",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := loadConfig()
			if err != nil {
				log.Fatal(err)
			}

			site := &domain.Site{
				PagesDir:         "pages",
				TemplatesDir:     "templates",
				AssetsDir:        "assets",
				DistDir:          "dist",
				EnableAutoReload: false,
				Config:           config,
				ConfigPath:       "config.yaml",
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

			config, err := loadConfig()
			if err != nil {
				log.Fatal(err)
			}

			site := &domain.Site{
				PagesDir:         "pages",
				TemplatesDir:     "templates",
				AssetsDir:        "assets",
				DistDir:          "dist",
				EnableAutoReload: watch,
				Config:           config,
				ConfigPath:       "config.yaml",
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
	serveCmd.Flags().BoolP("watch", "w", true, "Enable auto-reload on file changes")

	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(initCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
