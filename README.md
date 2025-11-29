# stw-cli

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![codecov](https://codecov.io/gh/EmiraLabs/stw-cli/branch/main/graph/badge.svg)](https://codecov.io/gh/EmiraLabs/stw-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/EmiraLabs/stw-cli)](https://goreportcard.com/report/github.com/EmiraLabs/stw-cli)

> **The static site generator designed for Cloudflare Pages**

Build clean, fast static sites and deploy to Cloudflare Pages with zero friction. Perfect for company profiles, landing pages, and portfolios.

**See it in action:** [os.emiralabs.com](https://os.emiralabs.com) → Built and deployed with stw-cli

---

## Why stw-cli?

**Stop wrestling with complex tools.** Most static sites don't need 300+ features—they need simplicity and fast Cloudflare deployment.

### Comparison with Popular Static Site Generators

| Feature | Jekyll | Hugo | 11ty | **stw-cli** |
|---------|--------|------|------|-------------|
| **Platform** | GitHub Pages | Generic | Generic | **Cloudflare Pages** ✅ |
| **Language** | Ruby | Go | JavaScript | **Go** |
| **Setup** | `gem install` | `brew install` | `npm install` | **`go install`** |
| **Config** | _config.yml | config.toml/yaml | .eleventy.js | **config.yaml** ✅ |
| **Cloudflare Deploy** | Manual | Manual | Manual | **`stw init --wrangler`** ✅ |
| **Setup Time** | 15 min | 30 min | 20 min | **5 min** ✅ |
| **Learning Curve** | Moderate | Steep | Moderate | **Gentle** ✅ |
| **Best For** | GitHub Pages blogs | Complex sites | JS developers | **Cloudflare Pages sites** ✅ |

**Key Insight:** If you're moving from **Jekyll on GitHub Pages** to **Cloudflare Pages**, stw-cli is the natural choice. Same simplicity, better platform.


### Perfect For

- 🚀 **Company profiles** - Professional sites deployed globally
- 💼 **Portfolio sites** - Showcase your work on Cloudflare's edge
- 📄 **Documentation** - Simple docs with fast global delivery
- ⚡ **Cloudflare-first projects** - Native Wrangler integration

### Not For

- Complex blogs with 1000+ posts → Use Hugo
- Multi-language sites with taxonomies → Use Hugo
- Need 300+ themes → Use Hugo

---

## Quick Start

### 1. Install stw-cli

**Quick Install (Linux/macOS):**
```bash
curl -sSL https://raw.githubusercontent.com/EmiraLabs/stw-cli/main/install.sh | sh
```

**Using Go:**
```bash
go install github.com/EmiraLabs/stw-cli/cmd/stw@latest
```

**Pre-built Binaries:**
Download for your platform from [GitHub Releases](https://github.com/EmiraLabs/stw-cli/releases/latest)

**Build from Source:**
```bash
git clone https://github.com/EmiraLabs/stw-cli.git
cd stw-cli
go build -o stw ./cmd/stw
```

### 2. Create Your Site

```bash
# Scaffold a new site from template
stw init my-site

# Or scaffold + configure for Cloudflare Pages
stw init my-site --wrangler

cd my-site
```

### 3. Build & Preview

```bash
# Start dev server with live reload
stw serve

# Or build for production
stw build
```

### 4. Deploy to Cloudflare Pages

**Option A: GitHub Integration (Recommended)**

1. Push your code to GitHub
2. Connect repository to Cloudflare Pages
3. Set build command: `stw build`
4. Set output directory: `dist`
5. Deploy! ✨

**Option B: Manual Deployment**

```bash
# Build locally
stw build

# Deploy with Wrangler
wrangler pages deploy dist
```

**That's it.** Your site is live on Cloudflare's global CDN.

---

## Features

✅ **Cloudflare-First** - Native Wrangler integration built-in  
✅ **Fast Builds** - Go-powered static generation in milliseconds  
✅ **Simple Config** - One YAML file, no complexity  
✅ **SEO-Ready** - Built-in meta tags, OpenGraph, and JSON-LD  
✅ **Live Reload** - Auto-refresh browser during development  
✅ **Free Hosting** - Deploy to Cloudflare Pages at no cost  
✅ **Clean Architecture** - Testable, maintainable codebase  

---

## Documentation

Complete documentation is available in the [docs/](docs/) directory:

- [Installation](docs/installation.md) - Setup and requirements
- [Quick Start](docs/quick-start.md) - 5-minute guide
- [Commands](docs/commands.md) - CLI reference
- [Configuration](docs/configuration.md) - YAML config guide
- [Templates](docs/templates.md) - Go template system
- [SEO Meta](docs/seo-meta.md) - Metadata configuration
- [Site Structure](docs/site-structure.md) - Project organization
- [Deployment](docs/deployment-cloudflare.md) - Cloudflare Pages guide
- [Troubleshooting](docs/troubleshooting.md) - Common issues

---

## Live Examples

**Real sites built with stw-cli:**

- [os.emiralabs.com](https://os.emiralabs.com) - EmiraLabs open source tools showcase ([source](https://github.com/EmiraLabs/os))

*Built a site with stw-cli? [Submit a PR](https://github.com/EmiraLabs/stw-cli/pulls) to add it here!*

---

## Project Structure

```
my-site/
├── config.yaml          # Site configuration
├── wrangler.json        # Cloudflare settings (if using --wrangler)
├── pages/               # Your content (HTML)
│   ├── index.html
│   └── about/
│       └── index.html
├── templates/           # Go templates
│   ├── base.html
│   ├── components/
│   └── partials/
├── assets/              # Static files
│   ├── css/
│   ├── js/
│   └── images/
└── dist/                # Generated output (git ignored)
```

---

## When to Use stw-cli vs Hugo

### Use stw-cli for:

- ✅ Company profiles and landing pages
- ✅ Portfolio and showcase sites
- ✅ Simple documentation sites
- ✅ Cloudflare Pages deployments
- ✅ Projects needing simplicity over features

### Use Hugo for:

- ✅ Blogs with 100+ posts
- ✅ Multi-language sites
- ✅ Complex taxonomy and categorization
- ✅ Need for extensive theme ecosystem

**Philosophy:** We built stw-cli for the 80% use case—simple sites that deserve simple tooling.

---

## Migrating from Jekyll?

If you're currently using Jekyll on GitHub Pages and want to move to Cloudflare Pages, stw-cli is your smoothest path:

**What's Familiar:**
- ✅ Simple YAML configuration (like Jekyll's `_config.yml`)
- ✅ Front matter for page metadata
- ✅ Template-based rendering
- ✅ Static file generation

**What's Better:**
- ⚡ Cloudflare's global CDN (faster than GitHub Pages)
- 🔒 Automatic SSL with custom domains (no CNAME gymnastics)
- 🚀 Workers integration for dynamic features
- 💰 More generous free tier

**Migration Steps:**
1. `stw init my-site --wrangler` - scaffold new project
2. Copy content from Jekyll's `_posts/` to stw's `pages/`
3. Adapt Liquid templates to Go templates (similar syntax)
4. Update `config.yaml` with your site settings
5. `git push` - auto-deploy to Cloudflare Pages

**Not ready to switch?** That's fine. Use Hugo for complex blogs. Use stw-cli when you need simplicity + Cloudflare.

---

## Development

### Running Tests

```bash
go test ./...
```

### Test Coverage

```bash
go test ./... -cover
```

### Contributing

We welcome contributions! See [CONTRIBUTING.md](docs/contributing.md) for guidelines.

---

## Why We Built This

After using Jekyll on GitHub Pages, we discovered Cloudflare Pages offered more power—better CDN, Workers integration, and easier SSL management. But Hugo felt like overkill for simple sites, and Jekyll wasn't designed for Cloudflare.

**stw-cli** fills that gap: a simple, Cloudflare-native static site generator for the majority of use cases.

---

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

## Links

- **Template Repository:** [github.com/EmiraLabs/stw](https://github.com/EmiraLabs/stw)
- **Documentation:** [docs/](docs/)
- **Live Example:** [os.emiralabs.com](https://os.emiralabs.com)
- **Issues:** [GitHub Issues](https://github.com/EmiraLabs/stw-cli/issues)

---

**Built with ❤️ by [EmiraLabs](https://emiralabs.com)**

*Making static site deployment to Cloudflare Pages stupidly simple.*