# Quick Start

Get started with stw-cli in under 5 minutes. This guide will walk you through creating your first static site and deploying it to Cloudflare Pages.

## Prerequisites

- Go 1.19+ installed ([download](https://golang.org/dl/))
- Git installed
- (Optional) Cloudflare account for deployment

## 1. Install stw-cli

```bash
# Using Go install (recommended)
go install github.com/EmiraLabs/stw-cli/cmd/stw@latest

# Or build from source
git clone https://github.com/EmiraLabs/stw-cli.git
cd stw-cli
go build -o stw ./cmd/stw
```

## 2. Create a New Site
 
Run the following command to create a new site from the official template:
 
```bash
# Basic site
stw init my-site

# Or with Cloudflare Pages configuration
stw init my-site --wrangler

cd my-site
```
 
This command clones the [official stw template](https://github.com/EmiraLabs/stw) and sets up your new project with all necessary files.

> **Live Example:** See [os.emiralabs.com](https://os.emiralabs.com) - a real site built with stw-cli. ([View source](https://github.com/EmiraLabs/os))
 
## 3. Build and Serve

Build the site:

```bash
stw build
```

Serve the site with live reload:

```bash
stw serve
# Or with custom port
stw serve --port 3000
```

Open http://localhost:8080 in your browser.

## 4. Enable Auto-Reload

Auto-reload is enabled by default with `stw serve`. Edit your files and see changes instantly!

To disable auto-reload:

```bash
stw serve --watch=false
```

## 5. Deploy to Cloudflare Pages

### Option A: GitHub Integration (Recommended)

1. **Push to GitHub:**
   ```bash
   git add .
   git commit -m "Initial commit"
   git remote add origin https://github.com/yourusername/my-site.git
   git push -u origin main
   ```

2. **Connect to Cloudflare Pages:**
   - Go to [Cloudflare Dashboard](https://dash.cloudflare.com/)
   - Navigate to Workers & Pages → Create application → Pages
   - Connect your GitHub repository
   - Set build command: `stw build`
   - Set output directory: `dist`
   - Deploy!

3. **Automatic Deployments:**
   - Every push to `main` triggers a new deployment
   - Pull requests get preview deployments

### Option B: Manual Deployment

```bash
# Install Wrangler
npm install -g wrangler

# Authenticate
wrangler auth login

# Build and deploy
stw build
wrangler pages deploy dist
```

## What's Next?

- [Configure SEO metadata](seo-meta.md)
- [Customize templates](templates.md)
- [Learn about site structure](site-structure.md)
- [Read the full documentation](../README.md)

## Live Examples

Check out these real sites built with stw-cli:

- **[os.emiralabs.com](https://os.emiralabs.com)** - Open source tools showcase
  - [Source code](https://github.com/EmiraLabs/os)
  - Features: Multiple tools, clean design, deployed on Cloudflare Pages

## Troubleshooting

### Command not found

Make sure `$(go env GOPATH)/bin` is in your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Template clone fails

Ensure you have Git installed and internet connection. The template is cloned from:
```
https://github.com/EmiraLabs/stw
```

### Build fails

Check that you're in the project directory with valid `config.yaml`, `pages/`, `templates/`, and `assets/` directories.

For more help, see the [Troubleshooting Guide](troubleshooting.md).