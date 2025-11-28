# Commands

stw-cli provides three main commands: `build`, `serve`, and `init`. This reference covers all available commands and their options.

## Global Options

All commands support these global options:

- `--help`, `-h`: Show help information
- `--version`, `-v`: Show version information

## build

Builds the static site from source files.

```bash
stw build
```

**Description:** Parses pages, applies templates, processes metadata, and copies assets to the `dist` directory.

**What it does:**
- Parses all HTML files in `pages/`
- Applies templates from `templates/`
- Processes SEO metadata from `config.yaml` and page front matter
- Copies all files from `assets/` to `dist/assets/`
- Generates the complete static site in `dist/`

**Output:** Static files in the `dist/` directory ready for deployment.

## serve

Builds and serves the site locally with optional auto-reload.

```bash
stw serve [options]
```

**Options:**
- `--port`, `-p` (string): Port to serve on (default: "8080")
- `--watch`, `-w` (bool): Enable auto-reload on file changes (default: true)

**Examples:**
```bash
# Serve on default port 8080 with auto-reload
stw serve

# Serve on custom port
stw serve --port 3000

# Serve without auto-reload
stw serve --watch=false

# Short form
stw serve -p 3000 -w
```

**What it does:**
- Builds the site initially
- Starts a local HTTP server
- Serves files from `dist/`
- If `--watch` is enabled:
  - Watches for changes in `pages/`, `templates/`, `assets/`, and `config.yaml`
  - Automatically rebuilds when files change
  - Notifies connected browsers to reload

**Auto-reload:** When enabled, the server injects JavaScript that connects to a Server-Sent Events endpoint. Changes trigger a browser reload.

## init

Initialize a new project from the official template, with optional Cloudflare Pages configuration.

```bash
stw init <project-name> [flags]
```

**Description:** Scaffolds a new static site project by cloning the [official stw template](https://github.com/EmiraLabs/stw) and optionally configures it for Cloudflare Pages deployment via Wrangler.

**Arguments:**
- `<project-name>` (required): Name of the project directory to create

**Flags:**
- `--wrangler`, `-w` (bool): Also configure Wrangler for Cloudflare Pages deployment (default: false)

**What it does:**
1. **Scaffolding:**
   - Clones the official template from `https://github.com/EmiraLabs/stw`
   - Creates a new directory with your project name
   - Sets up the complete project structure (pages/, templates/, assets/, config.yaml)
   - Removes Git history and initializes a fresh repository

2. **Wrangler Configuration (if `--wrangler` flag used):**
   - Prompts for project name (default: from project directory name)
   - Prompts for custom domain (e.g., yoursite.com)
   - Updates `wrangler.json` with the provided configuration
   - Prepares the project for immediate Cloudflare Pages deployment

**Examples:**

```bash
# Create a basic site
stw init my-website

# Create site and configure for Cloudflare Pages
stw init my-website --wrangler

# Short form
stw init my-website -w
```

**Example interaction with --wrangler:**
```
Initializing project: my-website
✓ Cloned template
✓ Created project directory

Configure Cloudflare Pages deployment:
Enter project name (default: my-website): my-awesome-site
Enter custom domain (e.g., yoursite.com): myawesome.com
✓ wrangler.json configured

Next steps:
  cd my-website
  stw serve              # Start development
  git push origin main   # Deploy to Cloudflare Pages
```

## Command Structure

```
stw [command] [flags]

Available Commands:
  build       Build the static site
  init        Initialize Wrangler configuration for deployment
  serve       Build and serve the static site

Flags:
  -h, --help   help for stw
  -v, --version   version for stw

Use "stw [command] --help" for more information about a command.
```

## Exit Codes

- `0`: Success
- `1`: Error (build failure, missing files, etc.)

## Examples

### Development Workflow

```bash
# Initial build
stw build

# Start development server with auto-reload
stw serve

# In another terminal, edit files and see changes automatically
# The server will rebuild and reload the browser
```

### Production Build

```bash
# Build for production
stw build

# Deploy the dist/ directory to your hosting provider
```

### Deployment Setup

```bash
# Initialize for Cloudflare Pages
stw init

# Then deploy
wrangler pages deploy dist
```

## Troubleshooting

### Build fails

- Check that `config.yaml` exists and is valid YAML
- Ensure `pages/`, `templates/`, and `assets/` directories exist
- Verify template files are valid HTML with Go template syntax

### Serve doesn't start

- Check if port 8080 is already in use: `lsof -i :8080`
- Try a different port: `stw serve --port 3000`

### Auto-reload not working

- Ensure your browser supports Server-Sent Events
- Check browser console for JavaScript errors
- Verify files are being saved (some editors have auto-save disabled)

### Init fails

- Ensure `wrangler.json` exists in the project root
- Check that the file is writable
- Verify the template syntax in `wrangler.json`