# SEO Meta System

stw-cli includes a comprehensive SEO meta system that allows you to define search engine optimization metadata at both the site and page levels.

## Overview

The SEO meta system provides:
- Site-wide default meta configuration
- Per-page overrides via front matter
- Automatic validation of meta fields
- Generation of all standard SEO tags, Open Graph, Twitter Cards, and JSON-LD

## Site-Wide Defaults

Configure default meta values in `config.yaml`:

```yaml
meta:
  title: "Default Site Title"
  description: "Default site description for SEO."
  canonical: ""
  robots: "index,follow"
  keywords: "default,keywords"
  og_title: ""
  og_description: ""
  og_image: ""
  twitter_title: ""
  twitter_description: ""
  twitter_image: ""
  jsonld: {}
```

## Per-Page Overrides

Add front matter to any page in `pages/` to override site defaults:

```yaml
---
title: "Custom Page Title"
description: "Page-specific description under 160 characters."
canonical: "https://example.com/page"
robots: "index,follow"
keywords: "page,specific,keywords"
og_title: "Open Graph Title"
og_description: "Open Graph description."
og_image: "/assets/images/og-image.jpg"
twitter_title: "Twitter Card Title"
twitter_description: "Twitter description."
twitter_image: "/assets/images/twitter-image.jpg"
jsonld:
  "@context": "https://schema.org"
  "@type": "WebPage"
  name: "Page Name"
  description: "Structured data description"
---
<!-- Your page content here -->
```

## Supported Meta Fields

| Field | Type | Max Length | Description |
|-------|------|------------|-------------|
| `title` | string | 60 chars | Page title |
| `description` | string | 160 chars | Meta description |
| `canonical` | string | - | Canonical URL |
| `robots` | string | - | Robots directive (e.g., "index,follow") |
| `keywords` | string | - | Comma-separated keywords |
| `og_title` | string | - | Open Graph title |
| `og_description` | string | - | Open Graph description |
| `og_image` | string | - | Open Graph image (must be under `/assets/`) |
| `twitter_title` | string | - | Twitter Card title |
| `twitter_description` | string | - | Twitter Card description |
| `twitter_image` | string | - | Twitter Card image |
| `jsonld` | object | - | Raw JSON-LD structured data object |

## Validation

The system validates meta fields during build:
- **Title length**: ≤ 60 characters
- **Description length**: ≤ 160 characters
- **Open Graph image**: Must be under `/assets/` path

## Front Matter Format

Currently supports YAML front matter delimited by `---` markers:

```yaml
---
# meta fields here
---
<!-- page content -->
```

## Generated HTML

The build process injects all meta tags into the `<head>` section of each page, including:

- Standard meta tags (`<meta name="description">`, `<meta name="keywords">`, etc.)
- Open Graph tags (`<meta property="og:*">`) for social sharing
- Twitter Card tags (`<meta name="twitter:*">`)
- JSON-LD structured data (`<script type="application/ld+json">`)

## Backward Compatibility

Pages without front matter automatically use site defaults from `config.yaml`. Existing sites will continue to work without modification.

## Examples

### Basic Page Meta

```yaml
---
title: "About Us"
description: "Learn more about our company and mission."
keywords: "about,company,mission"
---
```

### Social Media Optimized Page

```yaml
---
title: "Product Launch"
description: "We're excited to announce our new product!"
og_title: "🚀 New Product Launch!"
og_description: "Check out our latest innovation."
og_image: "/assets/images/product-launch.jpg"
twitter_title: "🚀 New Product Launch!"
twitter_description: "Check out our latest innovation."
twitter_image: "/assets/images/product-launch-twitter.jpg"
---
```

### Page with Structured Data

```yaml
---
title: "Contact Us"
description: "Get in touch with our team."
jsonld:
  "@context": "https://schema.org"
  "@type": "Organization"
  name: "Example Company"
  contactPoint:
    "@type": "ContactPoint"
    telephone: "+1-555-123-4567"
    contactType: "customer service"
---
```