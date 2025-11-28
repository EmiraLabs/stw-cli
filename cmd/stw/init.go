package main

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [site-name]",
	Short: "Initialize a new static site",
	Long:  `Initialize a new static site by cloning the official template.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		siteName := args[0]
		withWrangler, _ := cmd.Flags().GetBool("wrangler")
		
		// Check if git is installed
		if _, err := exec.LookPath("git"); err != nil {
			return fmt.Errorf("git is required to create a new site: %w", err)
		}

		fmt.Printf("Creating new site '%s' from template...\n", siteName)

		// Clone the repository
		cloneCmd := exec.Command("git", "clone", "https://github.com/EmiraLabs/stw.git", siteName)
		cloneCmd.Stdout = os.Stdout
		cloneCmd.Stderr = os.Stderr
		if err := cloneCmd.Run(); err != nil {
			return fmt.Errorf("failed to clone template: %w", err)
		}

		// Remove .git directory to start fresh
		if err := os.RemoveAll(filepath.Join(siteName, ".git")); err != nil {
			return fmt.Errorf("failed to remove .git directory: %w", err)
		}

		fmt.Printf("\nSuccessfully created site '%s'\n", siteName)

		if withWrangler {
			fmt.Println("\nInitializing Wrangler configuration...")
			
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Enter custom domain (e.g., yoursite.com): ")
			domain, _ := reader.ReadString('\n')
			domain = strings.TrimSpace(domain)
			if domain == "" {
				return fmt.Errorf("domain is required")
			}

			wranglerPath := filepath.Join(siteName, "wrangler.json")
			// Read the template from wrangler.json in the new site
			templateData, err := os.ReadFile(wranglerPath)
			if err != nil {
				return fmt.Errorf("failed to read wrangler.json template: %w", err)
			}

			// Parse and execute template
			tmpl, err := template.New("wrangler").Parse(string(templateData))
			if err != nil {
				return fmt.Errorf("failed to parse template: %w", err)
			}

			var buf strings.Builder
			data := map[string]interface{}{
				"project_name": siteName,
				"domain":       domain,
			}

			if err := tmpl.Execute(&buf, data); err != nil {
				return fmt.Errorf("failed to execute template: %w", err)
			}

			// Write back to file
			if err := os.WriteFile(wranglerPath, []byte(buf.String()), 0644); err != nil {
				return fmt.Errorf("failed to update wrangler.json: %w", err)
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
		}

		fmt.Printf("To get started:\n  cd %s\n  stw serve\n", siteName)

		return nil
	},
}

func init() {
	initCmd.Flags().Bool("wrangler", false, "Initialize Wrangler configuration for Cloudflare deployment")
}
