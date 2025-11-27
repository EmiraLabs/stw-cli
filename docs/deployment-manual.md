# Manual Deployment

Deploy your stw-cli site to any static hosting provider manually. This guide covers deployment to various platforms.

## General Process

1. **Build the site:**
   ```bash
   stw build
   ```

2. **Upload the `dist/` directory** to your hosting provider

3. **Configure your domain** (if needed)

## Hosting Providers

### Netlify

1. **Build locally:**
   ```bash
   stw build
   ```

2. **Drag and drop deployment:**
   - Go to [Netlify](https://netlify.com)
   - Drag the `dist/` folder to the deployment area
   - Or use Netlify CLI: `netlify deploy --dir dist`

3. **Set up continuous deployment:**
   - Connect your GitHub repository
   - Set build command: `stw build`
   - Set publish directory: `dist`

### Vercel

1. **Build locally:**
   ```bash
   stw build
   ```

2. **Deploy:**
   ```bash
   npx vercel --prod dist
   ```

3. **Or use Vercel CLI:**
   ```bash
   vercel dist
   ```

### GitHub Pages

1. **Build the site:**
   ```bash
   stw build
   ```

2. **Create `docs/` deployment:**
   - Rename `dist/` to `docs/`
   - Or configure GitHub Pages to use `dist/` branch

3. **Enable GitHub Pages:**
   - Go to repository Settings > Pages
   - Select "Deploy from a branch"
   - Choose `main` branch and `/docs` folder

### AWS S3 + CloudFront

1. **Build the site:**
   ```bash
   stw build
   ```

2. **Upload to S3:**
   ```bash
   aws s3 sync dist/ s3://your-bucket-name --delete
   ```

3. **Configure CloudFront:**
   - Create distribution pointing to S3 bucket
   - Set default root object to `index.html`
   - Configure error pages for SPA routing

### Firebase Hosting

1. **Install Firebase CLI:**
   ```bash
   npm install -g firebase-tools
   ```

2. **Initialize:**
   ```bash
   firebase init hosting
   ```

3. **Configure `firebase.json`:**
   ```json
   {
     "hosting": {
       "public": "dist",
       "ignore": [
         "firebase.json",
         "**/.*",
         "**/node_modules/**"
       ],
       "rewrites": [
         {
           "source": "**",
           "destination": "/index.html"
         }
       ]
     }
   }
   ```

4. **Deploy:**
   ```bash
   firebase deploy
   ```

### Surge

1. **Install Surge:**
   ```bash
   npm install -g surge
   ```

2. **Deploy:**
   ```bash
   stw build
   surge dist
   ```

### Render

1. **Connect repository** to Render
2. **Set build settings:**
   - Build command: `stw build`
   - Publish directory: `dist`

### Railway

1. **Connect repository** to Railway
2. **Configure as static site:**
   - Set build command: `stw build`
   - Set output directory: `dist`

## FTP/SFTP Upload

For traditional hosting:

1. **Build the site:**
   ```bash
   stw build
   ```

2. **Upload via FTP:**
   ```bash
   # Using lftp
   lftp -c "open ftp://user:pass@host; cd public_html; mirror -R dist/ ."
   ```

3. **Or use FileZilla/rclone** to upload the `dist/` directory

## Docker Deployment

Create a Dockerfile for containerized deployment:

```dockerfile
FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY dist/ /usr/share/nginx/html/
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

Build and run:
```bash
docker build -t my-site .
docker run -p 8080:80 my-site
```

## CI/CD Deployment

### GitHub Actions

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy
on:
  push:
    branches: [main]
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
      - name: Deploy to Netlify
        run: npx netlify-cli deploy --dir=dist --prod
        env:
          NETLIFY_AUTH_TOKEN: ${{ secrets.NETLIFY_AUTH_TOKEN }}
```

### GitLab CI

Create `.gitlab-ci.yml`:

```yaml
stages:
  - build
  - deploy

build:
  stage: build
  image: golang:1.21
  script:
    - go build -o stw ./cmd/stw
    - ./stw build
  artifacts:
    paths:
      - dist/

deploy:
  stage: deploy
  script:
    - echo "Deploy to your hosting provider"
  dependencies:
    - build
```

## Custom Domain Configuration

### DNS Setup

1. **Point domain to hosting provider:**
   - Update nameservers or DNS records
   - Configure CNAME or A records as required

2. **SSL certificates:**
   - Most providers offer automatic HTTPS
   - For manual setup, use Let's Encrypt

### Subdomain Setup

For `www` subdomain:
- Create CNAME record: `www.yourdomain.com` → `yourdomain.com`
- Or configure redirect in hosting provider

## Performance Optimization

### CDN Integration

- Use CDN for global distribution
- Configure caching headers
- Enable compression

### Build Optimizations

- Minify HTML, CSS, JS
- Optimize images
- Use WebP format
- Enable gzip compression

## Backup and Rollback

### Version Control

- Keep deployment history in Git
- Tag releases: `git tag v1.0.0`
- Rollback: `git checkout v1.0.0`

### Provider-specific Rollback

- Most platforms keep deployment history
- Use provider's rollback features
- Keep backup of `dist/` directory

## Monitoring

### Uptime Monitoring

- Use services like UptimeRobot, Pingdom
- Monitor response times
- Set up alerts for downtime

### Analytics

- Add Google Analytics or similar
- Monitor traffic and performance
- Track conversion goals

## Troubleshooting

### Common Issues

**404 errors:**
- Ensure `index.html` is in root
- Check server configuration for SPA routing
- Verify file permissions

**HTTPS issues:**
- Check SSL certificate status
- Update DNS configuration
- Clear browser cache

**Slow loading:**
- Optimize images
- Enable compression
- Use CDN
- Check server response times

**Build failures:**
- Test build locally first
- Check build logs
- Verify dependencies
- Ensure correct file paths

### Debugging

**Local testing:**
```bash
stw serve --port 3000
# Test locally before deploying
```

**Validate build:**
```bash
stw build
# Check dist/ directory contents
ls -la dist/
```

**Test deployment:**
- Deploy to staging environment first
- Use provider's preview features
- Test all pages and functionality