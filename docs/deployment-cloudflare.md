# Cloudflare Pages Deployment

Deploy your stw-cli site to Cloudflare Pages for fast, global content delivery with automatic HTTPS and CDN.

## Prerequisites

- Cloudflare account
- GitHub repository (recommended)
- Wrangler CLI (optional, for manual deployment)

## Quick Setup

### 1. Initialize Configuration

Run the init command to set up deployment:

```bash
stw init
```

This will prompt for:
- **Project name** (default: stw-site)
- **Custom domain** (your domain name)

### 2. Connect to Cloudflare Pages

#### Option A: GitHub Integration (Recommended)

1. **Authorize Cloudflare in GitHub:**
   - Go to your GitHub repository settings
   - Navigate to "Integrations" > "Applications"
   - Find "Cloudflare Pages" and click "Configure"
   - Select your repository and allow access

2. **Set up Pages in Cloudflare Dashboard:**
   - Go to [Cloudflare Dashboard](https://dash.cloudflare.com/)
   - Navigate to "Workers & Pages" > "Pages"
   - Click "Create application" > "Pages" tab
   - Select "Connect to Git"
   - Choose your GitHub repository
   - Click "Begin setup"

3. **Configure build settings:**
   - **Build command:** `stw build`
   - **Build output directory:** `dist`
   - **Root directory:** `/` (leave empty)

4. **Environment variables** (optional):
   - Add any environment variables your build needs

5. **Deploy:**
   - Cloudflare will build and deploy on every push to main
   - You can also trigger manual deployments

#### Option B: Manual Deployment

1. **Install Wrangler:**
   ```bash
   npm install -g wrangler
   ```

2. **Authenticate:**
   ```bash
   wrangler auth login
   ```

3. **Deploy:**
   ```bash
   stw build
   wrangler pages deploy dist
   ```

## Custom Domain

### Using Your Own Domain

1. **Add domain to Cloudflare:**
   - In Cloudflare dashboard, go to "Websites"
   - Click "Add site"
   - Enter your domain and follow the setup

2. **Configure Pages custom domain:**
   - Go to Pages project settings
   - Under "Custom domains", click "Add custom domain"
   - Enter your domain and click "Add"

3. **Update DNS:**
   - Cloudflare will provide DNS instructions
   - Update your domain registrar's DNS settings

### Cloudflare-provided Domain

Your site will be available at `https://your-project.pages.dev` automatically.

## Build Configuration

### Build Commands

For GitHub integration, use:
- **Build command:** `stw build`
- **Build output directory:** `dist`

### Environment Variables

Set these in your Pages project settings if needed:

- `NODE_VERSION`: If using Node.js tools
- `GO_VERSION`: If your build process needs Go
- Custom variables for your site configuration

## Deployment Process

### Automatic Deployment

With GitHub integration:
1. Push changes to your main branch
2. Cloudflare automatically detects the push
3. Runs `stw build` in a clean environment
4. Deploys the `dist/` directory
5. Site is live within minutes

### Manual Deployment

For testing or one-off deployments:

```bash
# Build locally
stw build

# Deploy
wrangler pages deploy dist --project-name your-project
```

## Preview Deployments

Cloudflare Pages creates preview deployments for:
- Pull requests
- Branch pushes (if configured)

Access previews at:
- `https://branch-name.your-project.pages.dev`
- `https://pr-number.your-project.pages.dev`

## Custom Build Scripts

If you need custom build steps, create a build script:

```bash
#!/bin/bash
# build.sh
stw build
# Additional build steps here
```

Then use `build.sh` as your build command.

## Environment-specific Builds

For different environments:

```yaml
# In GitHub Actions or build script
if [ "$CF_PAGES_BRANCH" = "main" ]; then
    # Production build
    cp config.prod.yaml config.yaml
else
    # Preview build
    cp config.dev.yaml config.yaml
fi
stw build
```

## Troubleshooting

### Build Fails

**Check build logs:**
- In Cloudflare dashboard, go to Pages > your project > Builds
- Click on a build to see detailed logs

**Common issues:**
- Missing `config.yaml`
- Invalid template syntax
- Missing template files
- Permission issues

### Site Not Loading

**Check deployment status:**
- Ensure build completed successfully
- Verify custom domain configuration
- Check DNS propagation (can take up to 24 hours)

**Clear cache:**
```bash
# Clear Cloudflare cache
wrangler pages deployment tail
```

### Custom Domain Issues

**DNS not updating:**
- Wait for DNS propagation
- Check DNS settings in domain registrar
- Verify domain is added to Cloudflare

**SSL certificate:**
- Cloudflare provides automatic HTTPS
- May take a few minutes for certificate issuance

## Performance Optimization

### Cloudflare Features

- **CDN:** Global content delivery
- **Caching:** Automatic caching rules
- **Compression:** Gzip/Brotli compression
- **Minification:** Automatic minification

### stw-cli Optimizations

- Optimize images in `assets/`
- Minify CSS and JavaScript
- Use appropriate caching headers
- Enable compression

## Cost

Cloudflare Pages is free for:
- Unlimited static sites
- 100 GB bandwidth/month
- 30,000 requests/month
- Custom domains

Paid plans available for higher limits.

## Migration from Other Platforms

### From Netlify

1. Export your site content
2. Set up stw-cli project structure
3. Configure build settings as above
4. Deploy to Cloudflare Pages

### From Vercel

1. Export your static files
2. Create stw-cli configuration
3. Use manual deployment or GitHub integration
4. Update DNS to point to Cloudflare

## Advanced Configuration

### Wrangler Configuration

The `wrangler.json` file controls deployment:

```json
{
  "name": "my-site",
  "compatibility_date": "2024-01-01",
  "assets": {
    "directory": "./dist",
    "binding": "ASSETS"
  },
  "routes": [
    {
      "pattern": "mysite.com",
      "custom_domain": true
    }
  ]
}
```

### Build Hooks

Use Cloudflare's build hooks for external notifications:

1. Create a build hook in Pages settings
2. Use the webhook URL in external services
3. Trigger builds programmatically

## Security

- All sites get automatic HTTPS
- DDoS protection included
- Web Application Firewall (WAF) available
- Access controls for team collaboration