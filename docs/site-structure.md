# Site Structure

Understanding the directory structure of your stw-cli site is key to working effectively with stw-cli. This guide explains the purpose of each directory and file in your site project.

## Root Directory

```
my-site/
├── config.yaml          # Site configuration
├── wrangler.json        # Cloudflare Pages deployment config
├── pages/               # HTML pages
├── templates/           # Template files
├── assets/              # Static assets
└── dist/                # Generated site (created by build)
```

## Core Directories

### pages/

Contains your site's HTML pages. Each `.html` file becomes a page in the built site.

```
pages/
├── index.html              # Home page (/)
├── about/
│   ├── index.html          # About page (/about/)
│   └── contact/
│       └── index.html      # Contact page (/about/contact/)
├── blog/
│   ├── index.html          # Blog listing (/blog/)
│   └── my-post.html        # Blog post (/blog/my-post.html)
└── 404.html                # 404 error page
```

**Rules:**
- Files named `index.html` create directory routes
- Other `.html` files create file routes
- Supports nested directories
- Each page can have YAML or JSON front matter for metadata

### templates/

Contains Go HTML templates that define the site's layout and structure.

```
templates/
├── base.html              # Main template (required)
├── components/
│   ├── header.html        # Header component
│   ├── footer.html        # Footer component
│   └── sidebar.html       # Sidebar component
└── partials/
    ├── head.html          # Head section with meta tags
    └── scripts.html       # JavaScript includes
```

**Required templates:**
- `base.html`: Main template that includes other templates

**Common patterns:**
- `components/`: Reusable UI components
- `partials/`: Template fragments included in base.html

### assets/

Static files that are copied unchanged to the built site.

```
assets/
├── css/
│   ├── styles.css         # Main stylesheet
│   └── theme.css          # Theme styles
├── js/
│   ├── app.js             # Main JavaScript
│   └── utils.js           # Utility functions
├── images/
│   ├── logo.png           # Site logo
│   └── hero.jpg           # Hero image
└── fonts/
    └── custom-font.woff   # Custom fonts
```

**Rules:**
- All files are copied to `dist/assets/`
- Maintain the same directory structure
- Referenced in templates with `/assets/` prefix

### dist/

Generated directory containing the built static site. Created by `stw build`.

```
dist/
├── index.html             # Built home page
├── about/
│   └── index.html         # Built about page
├── assets/
│   ├── css/
│   ├── js/
│   └── images/
└── 404.html               # Built 404 page
```

**Note:** This directory is recreated on each build. Don't edit files here directly.

## Configuration Files

### config.yaml

Site-wide configuration including navigation, content data, and SEO defaults.

```yaml
navigations:
  - title: Home
    url: /

meta:
  title: "My Site"
  description: "Site description"
```

### wrangler.json

Cloudflare Pages deployment configuration. Modified by `stw init`.

```json
{
  "name": "my-site",
  "assets": {
    "directory": "./dist"
  },
  "routes": [
    {
      "pattern": "mysite.com",
      "custom_domain": true
    }
  ]
}
```

## File Naming Conventions

### Pages
- Use `index.html` for directory routes
- Use kebab-case for filenames: `my-blog-post.html`
- Use lowercase for directories: `blog/`, `about/`

### Templates
- Use `base.html` for the main template
- Use descriptive names: `header.html`, `footer.html`
- Group related templates in subdirectories

### Assets
- Use organized subdirectories: `css/`, `js/`, `images/`
- Use descriptive filenames with extensions
- Use kebab-case or camelCase consistently

## Example Project

```
my-blog/
├── config.yaml
├── wrangler.json
├── pages/
│   ├── index.html
│   ├── about/
│   │   └── index.html
│   └── blog/
│       ├── index.html
│       └── hello-world.html
├── templates/
│   ├── base.html
│   ├── components/
│   │   ├── header.html
│   │   └── footer.html
│   └── partials/
│       └── head.html
├── assets/
│   ├── css/
│   │   └── styles.css
│   ├── js/
│   │   └── app.js
│   └── images/
│       └── logo.png
└── dist/
    ├── index.html
    ├── about/
    │   └── index.html
    ├── blog/
    │   ├── index.html
    │   └── hello-world.html
    └── assets/
        ├── css/
        ├── js/
        └── images/
```

## Best Practices

### Organization
- Keep pages organized in logical directory structures
- Group related templates in subdirectories
- Use consistent naming conventions

### Maintenance
- Don't edit files in `dist/` - they get overwritten
- Keep configuration in `config.yaml`
- Use front matter in pages for page-specific settings

### Performance
- Optimize images in `assets/images/`
- Minify CSS and JS files
- Use appropriate file formats

### SEO
- Use descriptive filenames for pages
- Include metadata in `config.yaml` and page front matter
- Ensure proper heading hierarchy in templates