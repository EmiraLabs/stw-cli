# Quick Start

Get started with stw-cli in under 5 minutes. This guide will walk you through creating your first static site.

## 1. Create a New Site
 
 Run the following command to create a new site with the default structure and configuration:
 
 ```bash
 stw init my-site
 # Or to also configure Wrangler for deployment:
 # stw init my-site --wrangler
 cd my-site
 ```
 
 This command clones the [official template](https://github.com/EmiraLabs/stw) and sets up your new site.
 
 ## 2. Build and Serve

Build the site:

```bash
/path/to/stw build
```

Serve the site:

```bash
/path/to/stw serve
```

Open http://localhost:8080 in your browser.

## 8. Add Auto-Reload (Optional)

For development, enable auto-reload:

```bash
/path/to/stw serve --watch
```

Now edit your files and see changes instantly in the browser!

## What's Next?

- [Configure SEO metadata](seo-meta.md)
- [Customize templates](templates.md)
- [Deploy your site](deployment-cloudflare.md)

## Example Project

For a complete example, check out the [example site](https://github.com/EmiraLabs/stw-cli/tree/main/example) in the repository.