# Quick Start

Get started with stw-cli in under 5 minutes. This guide will walk you through creating your first static site.

## 1. Create a New Site

Create a directory for your site:

```bash
mkdir my-site
cd my-site
```

## 2. Initialize Site Structure

Create the basic directory structure:

```bash
mkdir -p pages assets/css assets/js templates/components templates/partials
```

## 3. Create Configuration

Create `config.yaml`:

```yaml
navigations:
  - title: Home
    url: /
  - title: About
    url: /about/

home:
  contents:
    - title: "<h1>Welcome to My Site</h1>"
      content: "This is my first static site built with stw-cli!"

meta:
  title: "My Site"
  description: "A static site built with stw-cli"
```

## 4. Create Templates

Create `templates/base.html`:

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Meta.Title}}</title>
    <meta name="description" content="{{.Meta.Description}}">
    <link rel="stylesheet" href="/assets/css/styles.css">
</head>
<body>
    <header>
        <nav>
            {{range .Config.navigations}}
            <a href="{{.url}}">{{.title}}</a>
            {{end}}
        </nav>
    </header>
    <main>
        {{.Content}}
    </main>
</body>
</html>
```

## 5. Create Pages

Create `pages/index.html`:

```html
<h1>Welcome</h1>
<p>This is the home page of my site.</p>

{{range .Config.home.contents}}
<div>
    {{.title}}
    <p>{{.content}}</p>
</div>
{{end}}
```

Create `pages/about/index.html`:

```html
<h1>About</h1>
<p>Learn more about this site.</p>
```

## 6. Add Styles

Create `assets/css/styles.css`:

```css
body {
    font-family: Arial, sans-serif;
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
}

nav {
    margin-bottom: 20px;
}

nav a {
    margin-right: 15px;
    text-decoration: none;
    color: #007acc;
}
```

## 7. Build and Serve

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