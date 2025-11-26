# stw-cli

A simple static website generator and server CLI tool written in Go.

## Features

- Build static websites from HTML pages and templates
- Serve the built site locally
- Copy static assets automatically
- Auto-reload browser when files change during development
- Deploy to Cloudflare Pages using Wrangler

## Installation

Ensure you have Go installed on your system.

Clone the repository and build the binary:

```bash
git clone https://github.com/EmiraLabs/stw-cli.git
cd stw-cli
go build -o stw ./cmd/stw
```

## Usage

### Build the site

```bash
./stw build
```

This will generate the site in the `dist` directory.

### Serve the site

```bash
./stw serve
```

This will build the site and serve it on `http://localhost:8080`. The site will automatically rebuild and reload in the browser when HTML, template, or asset files are changed.

To disable auto-reload:

```bash
./stw serve --watch=false
```

### Initialize Deployment Configuration

To set up deployment configuration:

```bash
./stw init
```

This will prompt you for:
- **Project name** (default: stw-site)
- **Custom domain** (required, e.g., yoursite.com)

It creates a `wrangler.json` file with your specified values, configured for Cloudflare Workers deployment with static assets.

### Deploy with Cloudflare Pages

After running `init`, to deploy your site to Cloudflare Pages:

1. **Authorize Cloudflare in GitHub:**
   - Go to your GitHub repository settings
   - Navigate to "Integrations" > "Applications"
   - Find "Cloudflare Pages" and click "Configure"
   - Select your repository and allow access

2. **Set up Pages in Cloudflare Dashboard:**
   - In the Cloudflare dashboard, go to the Workers & Pages page
   - Click "Create application"
   - Select the "Pages" tab
   - Select "Connect to Git"
   - Choose your GitHub repository and click "Begin setup"

3. **Configure build settings:**
   - **Build command:** `./stw build`
   - **Build output directory:** `dist`
   - **Root directory:** `/` (leave empty)

4. **Deploy:**
   - Cloudflare will automatically build and deploy your site on every push to the main branch
   - You can also trigger manual deployments from the dashboard

### Deploy the Site

#### Option 1: Manual Deployment with Wrangler

After building the site:

```bash
./stw build
wrangler pages deploy dist
```

#### Option 2: Automatic Deployment with GitHub and Cloudflare Workers

For seamless deployments on every push:

1. **Connect your repository to Cloudflare Workers:**
   - Use Wrangler to deploy initially: `wrangler deploy`
   - Or set up GitHub Actions with Wrangler

2. **Configure build settings:**
   - **Build command:** `./stw build`
   - **Deploy command:** `wrangler deploy`

3. **Environment variables:**
   - Set `CLOUDFLARE_API_TOKEN` and `CLOUDFLARE_ACCOUNT_ID` in your GitHub repo secrets

4. **GitHub Actions example:**
   ```yaml
   name: Deploy to Cloudflare Workers
   on:
     push:
       branches:
         - main
   jobs:
     deploy:
       runs-on: ubuntu-latest
       steps:
         - uses: actions/checkout@v3
         - uses: actions/setup-go@v3
           with:
             go-version: '1.21'
         - run: go build -o stw ./cmd/stw
         - run: ./stw build
         - run: npm install -g wrangler
         - run: wrangler deploy
           env:
             CLOUDFLARE_API_TOKEN: ${{ secrets.CLOUDFLARE_API_TOKEN }}
   ```

This setup ensures your static site is automatically updated whenever you push changes to GitHub.

## Project Structure

- `pages/`: Directory containing HTML pages
- `templates/`: Directory containing base template and components
- `assets/`: Directory containing static assets (CSS, JS, images, etc.)
- `dist/`: Output directory for the built site (generated)