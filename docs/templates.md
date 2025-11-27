# Templates

stw-cli uses Go HTML templates to render your pages. Templates define the structure and layout of your site, allowing you to create reusable components and maintain consistent design.

## Template Engine

stw-cli uses Go's `html/template` package, which provides:
- Safe HTML escaping
- Template inheritance
- Custom functions
- Conditional rendering
- Loops and variables

## Template Structure

Templates are organized in the `templates/` directory:

```
templates/
├── base.html              # Main template
├── components/
│   ├── header.html        # Header component
│   ├── footer.html        # Footer component
│   └── navigation.html    # Navigation component
└── partials/
    ├── head.html          # Head section
    └── scripts.html       # JavaScript includes
```

## Base Template

The `base.html` template is the main template that defines the HTML structure. It includes other templates using `{{template}}` directives.

```html
<!DOCTYPE html>
<html lang="en">
<head>
    {{template "head.html" .}}
</head>
<body>
    {{template "header.html" .}}
    <main>
        {{.Content}}
    </main>
    {{template "footer.html" .}}
</body>
</html>
```

**Key points:**
- `{{.Content}}` is replaced with the rendered page content
- `{{template "name.html" .}}` includes other templates
- The dot `.` passes the current context (page data) to included templates

## Template Data

Each template receives a `Page` struct with these fields:

- `.Title`: Page title
- `.Content`: Rendered HTML content
- `.Path`: Relative path
- `.IsDev`: Boolean indicating development mode
- `.Config`: Site configuration from `config.yaml`
- `.Meta`: SEO metadata

## Head Template

The head template (`partials/head.html`) handles meta tags and includes:

```html
{{define "head.html"}}
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Meta.Title}}</title>
    {{if .Meta.Description}}
    <meta name="description" content="{{.Meta.Description}}">
    {{end}}
    {{if .Meta.Canonical}}
    <link rel="canonical" href="{{.Meta.Canonical}}">
    {{end}}
    <link rel="stylesheet" href="/assets/css/styles.css">
{{end}}
```

## Component Templates

Components are reusable template fragments. Define them with `{{define "name"}}`:

```html
<!-- components/header.html -->
{{define "header.html"}}
<header class="header">
    <div class="container">
        <h1>{{.Config.site_name}}</h1>
        <nav>
            {{range .Config.navigations}}
            <a href="{{.url}}">{{.title}}</a>
            {{end}}
        </nav>
    </div>
</header>
{{end}}
```

## Template Functions

### Built-in Functions

- `{{if condition}}...{{end}}`: Conditional rendering
- `{{range .Items}}...{{end}}`: Loop over arrays
- `{{.Field}}`: Access struct fields
- `{{len .Array}}`: Get array length

### Custom Functions

stw-cli provides a `toJson` function for JSON-LD structured data:

```html
<script type="application/ld+json">
{{.Meta.JsonLd | toJson}}
</script>
```

## Accessing Configuration

Use `{{.Config.key}}` to access data from `config.yaml`:

```html
<!-- Navigation -->
<nav>
    {{range .Config.navigations}}
    <a href="{{.url}}">{{.title}}</a>
    {{end}}
</nav>

<!-- Site info -->
<footer>
    <p>{{.Config.footer.copyright}}</p>
</footer>
```

## Conditional Rendering

Use `{{if}}` for conditional content:

```html
{{if .IsDev}}
<script>
    console.log('Development mode');
</script>
{{end}}
```

## Loops

Use `{{range}}` to iterate over arrays:

```html
<ul>
    {{range .Config.features}}
    <li>
        <h3>{{.title}}</h3>
        <p>{{.description}}</p>
    </li>
    {{end}}
</ul>
```

## Page Content

Page content is available as `{{.Content}}`. Pages can use template syntax too:

```html
<!-- pages/index.html -->
<h1>Welcome</h1>
<p>Site tagline: {{.Config.tagline}}</p>

{{range .Config.features}}
<div class="feature">
    <h2>{{.title}}</h2>
    <p>{{.description}}</p>
</div>
{{end}}
```

## Template Inheritance

Templates can include other templates. The base template includes components, which can include partials.

## Best Practices

### Organization
- Keep templates modular and reusable
- Use `components/` for UI components
- Use `partials/` for template fragments
- Name templates descriptively

### Performance
- Minimize template complexity
- Avoid deep nesting
- Use conditional rendering sparingly

### Maintainability
- Document complex template logic
- Keep business logic out of templates
- Use consistent naming conventions

## Examples

### Complete Base Template

```html
<!DOCTYPE html>
<html lang="en">
<head>
    {{template "partials/head.html" .}}
</head>
<body class="{{if .IsDev}}dev{{end}}">
    {{template "components/header.html" .}}

    <main class="main">
        <div class="container">
            {{.Content}}
        </div>
    </main>

    {{template "components/footer.html" .}}

    {{template "partials/scripts.html" .}}
</body>
</html>
```

### Navigation Component

```html
{{define "components/navigation.html"}}
<nav class="nav">
    <ul>
        {{range .Config.navigations}}
        <li class="{{if eq $.Path .url}}active{{end}}">
            <a href="{{.url}}">{{.title}}</a>
        </li>
        {{end}}
    </ul>
</nav>
{{end}}
```

### Blog Post Template

```html
{{define "blog-post.html"}}
<article class="post">
    <header>
        <h1>{{.Title}}</h1>
        <time>{{.Date}}</time>
    </header>

    <div class="content">
        {{.Content}}
    </div>

    <footer>
        <div class="tags">
            {{range .Tags}}
            <span class="tag">{{.}}</span>
            {{end}}
        </div>
    </footer>
</article>
{{end}}
```

## Troubleshooting

### Template not found
- Ensure template files exist in `templates/`
- Check spelling in `{{template}}` directives
- Verify file permissions

### Variable undefined
- Check if the variable is in the Page struct
- Use `{{if .Field}}` to check existence
- Debug with `{{printf "%+v" .}}`

### HTML escaping issues
- Use `template.HTML` for trusted HTML in config
- Be careful with user-generated content
- Use `{{.Field | safeHTML}}` if needed (not recommended)