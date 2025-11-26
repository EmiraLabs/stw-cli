package application

import (
	"os/exec"

	"github.com/EmiraLabs/stw-cli/internal/infrastructure"
)

// SiteInitializer handles initializing the project with features
type SiteInitializer struct {
	fs infrastructure.FileSystem
}

// NewSiteInitializer creates a new SiteInitializer
func NewSiteInitializer(fs infrastructure.FileSystem) *SiteInitializer {
	return &SiteInitializer{
		fs: fs,
	}
}

// InitTailwind initializes the project with Tailwind CSS, PostCSS, and Autoprefixer
func (si *SiteInitializer) InitTailwind() error {
	// Create package.json
	packageJSON := `{
  "name": "stw-site",
  "version": "1.0.0",
  "scripts": {
    "build:css": "postcss assets/css/styles.css -o dist/assets/css/styles.css"
  },
  "devDependencies": {
    "tailwindcss": "^3.4.0",
    "postcss": "^8.4.0",
    "autoprefixer": "^10.4.0",
		"postcss-cli": "^10.1.0"
  }
}`
	if err := si.fs.WriteFile("package.json", []byte(packageJSON), 0644); err != nil {
		return err
	}

	// Create postcss.config.js
	postcssConfig := `module.exports = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  }
}`
	if err := si.fs.WriteFile("postcss.config.js", []byte(postcssConfig), 0644); err != nil {
		return err
	}

	// Create tailwind.config.js
	tailwindConfig := `module.exports = {
  content: ["./pages/**/*.{html}", "./templates/**/*.{html}", "./dist/**/*.html"],
  theme: {
    extend: {},
  },
  plugins: [],
}`
	if err := si.fs.WriteFile("tailwind.config.js", []byte(tailwindConfig), 0644); err != nil {
		return err
	}

	// Update styles.css
	stylesCSS := `@tailwind base;
@tailwind components;
@tailwind utilities;`
	if err := si.fs.WriteFile("assets/css/styles.css", []byte(stylesCSS), 0644); err != nil {
		return err
	}

	// Run npm install
	cmd := exec.Command("npm", "install")
	return cmd.Run()
}
