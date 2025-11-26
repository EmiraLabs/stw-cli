package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

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

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize Wrangler configuration for deployment",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Enter project name (default: stw-site): ")
			projectName, _ := reader.ReadString('\n')
			projectName = strings.TrimSpace(projectName)
			if projectName == "" {
				projectName = "stw-site"
			}

			fmt.Print("Enter custom domain (e.g., yoursite.com): ")
			domain, _ := reader.ReadString('\n')
			domain = strings.TrimSpace(domain)
			if domain == "" {
				log.Fatal("Domain is required")
			}

			// Read the template from wrangler.json
			templateData, err := os.ReadFile("wrangler.json")
			if err != nil {
				log.Fatal("Failed to read wrangler.json template:", err)
			}

			// Parse and execute template
			tmpl, err := template.New("wrangler").Parse(string(templateData))
			if err != nil {
				log.Fatal("Failed to parse template:", err)
			}

			var buf strings.Builder
			data := map[string]interface{}{
				"project_name": projectName,
				"domain":       domain,
			}

			if err := tmpl.Execute(&buf, data); err != nil {
				log.Fatal("Failed to execute template:", err)
			}

			// Write back to file
			if err := os.WriteFile("wrangler.json", []byte(buf.String()), 0644); err != nil {
				log.Fatal("Failed to update wrangler.json:", err)
			}

			fmt.Println("wrangler.json updated successfully.")
			fmt.Println("")
			fmt.Println("Next steps for deployment:")
			fmt.Println("")
			fmt.Println("1. Authorize Cloudflare in GitHub:")
			fmt.Println("   - Go to your GitHub repository settings")
			fmt.Println("   - Navigate to 'Integrations' > 'Applications'")
			fmt.Println("   - Find 'Cloudflare Pages' and click 'Configure'")
			fmt.Println("   - Select your repository and allow access")
			fmt.Println("")
			fmt.Println("2. Set up Pages in Cloudflare Dashboard:")
			fmt.Println("   - Go to Workers & Pages page")
			fmt.Println("   - Click 'Create application'")
			fmt.Println("   - Select 'Pages' tab")
			fmt.Println("   - Select 'Connect to Git'")
			fmt.Println("   - Choose your repository and click 'Begin setup'")
			fmt.Println("")
			fmt.Println("3. Configure build settings:")
			fmt.Println("   - Build command: ./stw build")
			fmt.Println("   - Build output directory: dist")
			fmt.Println("   - Root directory: / (leave empty)")
			fmt.Println("")
			fmt.Println("4. Deploy automatically on every push to main branch!")
		},
	}

	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(initCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
