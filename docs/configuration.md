# Configuration

stw-cli uses a `config.yaml` file to configure site-wide settings, navigation, content, and metadata defaults.

## File Location

The configuration file should be named `config.yaml` and placed in the root of your project.

## Basic Structure

```yaml
# Navigation menu items
navigations:
  - title: Home
    url: /
  - title: About
    url: /about/
  - title: Contact
    url: /contact/

# Site content data
home:
  contents:
    - title: "<h1>Welcome</h1>"
      content: "Welcome to our site!"

# SEO metadata defaults
meta:
  title: "Default Site Title"
  description: "Default site description"
  robots: "index,follow"
```

## Configuration Sections

### Navigation (`navigations`)

Defines the site navigation menu. Each item has:

- `title`: Display text for the menu item
- `url`: Relative URL path

```yaml
navigations:
  - title: Home
    url: /
  - title: About
    url: /about/
  - title: Services
    url: /services/
  - title: Contact
    url: /contact/
```

### Content Data

You can define arbitrary content sections that can be accessed in templates via `{{.Config.sectionName}}`.

```yaml
# Home page content
home:
  hero_title: "Welcome to Our Site"
  hero_subtitle: "We build amazing things"
  features:
    - title: "Fast"
      description: "Lightning fast performance"
    - title: "Secure"
      description: "Built with security in mind"

# Footer content
footer:
  copyright: "© 2024 My Company"
  links:
    - title: "Privacy Policy"
      url: "/privacy/"
    - title: "Terms of Service"
      url: "/terms/"
```

### SEO Metadata (`meta`)

Site-wide defaults for SEO metadata. These can be overridden per page using front matter.

```yaml
meta:
  title: "My Website"
  description: "A description of my website under 160 characters"
  canonical: ""
  robots: "index,follow"
  keywords: "keyword1,keyword2,keyword3"
  og_title: ""
  og_description: ""
  og_image: ""
  twitter_title: ""
  twitter_description: ""
  twitter_image: ""
  jsonld: {}
```

See [SEO Meta System](seo-meta.md) for complete details on metadata configuration.

## Template Usage

Access configuration data in templates using `{{.Config.key}}`:

```html
<!-- Navigation -->
<nav>
    {{range .Config.navigations}}
    <a href="{{.url}}">{{.title}}</a>
    {{end}}
</nav>

<!-- Content -->
<h1>{{.Config.home.hero_title}}</h1>
<p>{{.Config.home.hero_subtitle}}</p>

<!-- Features -->
{{range .Config.home.features}}
<div class="feature">
    <h3>{{.title}}</h3>
    <p>{{.description}}</p>
</div>
{{end}}
```

## Data Types

Configuration supports various data types:

```yaml
# Strings
site_name: "My Site"

# Numbers
version: 1.0

# Booleans
debug: false

# Arrays
tags: ["tag1", "tag2", "tag3"]

# Objects
author:
  name: "John Doe"
  email: "john@example.com"

# Nested structures
menu:
  main:
    - title: "Home"
      url: "/"
    - title: "About"
      url: "/about/"
  footer:
    - title: "Privacy"
      url: "/privacy/"
```

## HTML Content

For HTML content in configuration, use quotes or the `|` multiline syntax:

```yaml
# Using quotes
content: "<p>This is <strong>HTML</strong> content.</p>"

# Using multiline
content: |
  <div class="hero">
    <h1>Welcome</h1>
    <p>This is multiline HTML content.</p>
  </div>
```

## Validation

Configuration is validated during build. Common issues:

- Invalid YAML syntax
- Missing required fields (none currently required)
- Incorrect data types

## Examples

### Blog Configuration

```yaml
navigations:
  - title: Home
    url: /
  - title: Blog
    url: /blog/
  - title: About
    url: /about/

blog:
  title: "My Blog"
  description: "Thoughts on technology and life"
  posts_per_page: 10

meta:
  title: "My Blog"
  description: "A personal blog about technology"
  og_image: "/assets/images/blog-og.jpg"
```

### Business Site Configuration

```yaml
navigations:
  - title: Home
    url: /
  - title: Services
    url: /services/
  - title: About
    url: /about/
  - title: Contact
    url: /contact/

company:
  name: "Acme Corp"
  tagline: "Building the future"
  founded: 2020

services:
  - title: "Web Development"
    description: "Custom web applications"
    icon: "💻"
  - title: "Consulting"
    description: "Technical consulting"
    icon: "🤝"

meta:
  title: "Acme Corp - Building the Future"
  description: "Professional web development and consulting services"
```

## Environment-Specific Configuration

For different environments, you can use multiple config files and rename them as needed:

- `config.dev.yaml` - Development configuration
- `config.prod.yaml` - Production configuration

Copy the appropriate file to `config.yaml` before building.