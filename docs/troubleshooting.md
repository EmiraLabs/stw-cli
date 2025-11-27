# Troubleshooting

Common issues and solutions when using stw-cli.

## Installation Issues

### "command not found" after installation

**Problem:** `stw` command not found in PATH.

**Solutions:**
1. **Check Go installation:**
   ```bash
   go version
   ```

2. **Add GOPATH/bin to PATH:**
   ```bash
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

3. **Or use full path:**
   ```bash
   /path/to/stw build
   ```

4. **For local builds, ensure executable:**
   ```bash
   chmod +x stw
   ```

### Build fails with "module not found"

**Problem:** Go modules not properly initialized.

**Solutions:**
```bash
# Clean module cache
go clean -modcache

# Download dependencies
go mod download

# Tidy modules
go mod tidy
```

## Build Issues

### "config.yaml not found"

**Problem:** Missing configuration file.

**Solutions:**
1. **Create config.yaml:**
   ```yaml
   meta:
     title: "My Site"
   ```

2. **Check current directory:**
   ```bash
   pwd
   ls -la config.yaml
   ```

### Template parsing errors

**Problem:** Invalid template syntax.

**Common issues:**
- Unclosed template tags: `{{if .Field}` without `{{end}}`
- Invalid variable references
- Missing template files

**Solutions:**
1. **Check template syntax:**
   ```bash
   # The build will show specific error lines
   stw build
   ```

2. **Validate template files:**
   - Ensure all `{{template}}` references exist
   - Check for typos in variable names
   - Verify template structure

### "assets directory not found"

**Problem:** Missing assets directory.

**Solutions:**
1. **Create directory structure:**
   ```bash
   mkdir -p assets/css assets/js assets/images
   ```

2. **Or disable assets copying** (not recommended)

## Serve Issues

### Port already in use

**Problem:** Port 8080 is occupied.

**Solutions:**
1. **Use different port:**
   ```bash
   stw serve --port 3000
   ```

2. **Find process using port:**
   ```bash
   lsof -i :8080
   kill -9 <PID>
   ```

### Auto-reload not working

**Problem:** Browser doesn't refresh on file changes.

**Solutions:**
1. **Check browser console** for JavaScript errors

2. **Verify file watching:**
   - Ensure files are saved (not just modified)
   - Check if editor has auto-save disabled

3. **Test manually:**
   ```bash
   # Modify a file and check if build triggers
   touch pages/index.html
   ```

4. **Check file permissions:**
   ```bash
   ls -la pages/
   ```

### Server not accessible

**Problem:** Cannot access http://localhost:8080

**Solutions:**
1. **Check if server started:**
   ```bash
   # Look for "Serving dist on http://localhost:8080"
   stw serve
   ```

2. **Firewall blocking:**
   - Check firewall settings
   - Try different port

3. **Network issues:**
   - Try 127.0.0.1:8080 instead of localhost
   - Check VPN/proxy settings

## Template Issues

### Variables not rendering

**Problem:** Template variables show as empty or undefined.

**Solutions:**
1. **Check variable names:**
   - Use `{{.Config.variable}}` for config
   - Use `{{.Meta.variable}}` for metadata
   - Use `{{.Title}}` for page title

2. **Debug template data:**
   ```html
   <!-- Temporary debug -->
   <pre>{{printf "%+v" .}}</pre>
   ```

3. **Check config.yaml syntax:**
   ```bash
   # Validate YAML
   python3 -c "import yaml; yaml.safe_load(open('config.yaml'))"
   ```

### Template inheritance not working

**Problem:** Included templates not rendering.

**Solutions:**
1. **Check template paths:**
   - `{{template "head.html" .}}` looks in `templates/`
   - `{{template "partials/head.html" .}}` looks in `templates/partials/`

2. **Verify file exists:**
   ```bash
   ls -la templates/
   ```

3. **Check define blocks:**
   ```html
   {{define "component.html"}}
   <!-- content -->
   {{end}}
   ```

## Metadata Issues

### Front matter not parsed

**Problem:** Page metadata not applied.

**Solutions:**
1. **Check front matter format:**
   ```yaml
   ---
   title: "Page Title"
   description: "Page description"
   ---
   <!-- page content -->
   ```

2. **Validate YAML syntax:**
   - Use online YAML validator
   - Check for tabs vs spaces

3. **Test parsing:**
   ```go
   // In Go code, test the ParseFrontMatter function
   ```

### SEO validation errors

**Problem:** Build fails with metadata validation.

**Common validations:**
- Title > 60 characters
- Description > 160 characters
- og_image not under /assets/

**Solutions:**
1. **Check lengths:**
   ```bash
   # Count characters
   echo "Your title" | wc -c
   ```

2. **Fix image paths:**
   ```yaml
   og_image: "/assets/images/og.jpg"  # Correct
   og_image: "og.jpg"                  # Wrong
   ```

## Asset Issues

### Assets not copied

**Problem:** CSS/JS/images not in dist.

**Solutions:**
1. **Check directory structure:**
   ```
   assets/
     css/
     js/
     images/
   ```

2. **Verify file permissions:**
   ```bash
   ls -la assets/
   ```

3. **Check build output:**
   ```bash
   stw build
   ls -la dist/assets/
   ```

### Broken asset links

**Problem:** Assets not loading in browser.

**Solutions:**
1. **Check paths in templates:**
   ```html
   <link rel="stylesheet" href="/assets/css/styles.css">
   ```

2. **Verify file exists in dist:**
   ```bash
   ls -la dist/assets/css/styles.css
   ```

## Deployment Issues

### Cloudflare Pages build fails

**Problem:** Deployment build fails.

**Solutions:**
1. **Check build logs** in Cloudflare dashboard

2. **Test build locally:**
   ```bash
   stw build
   ```

3. **Verify wrangler.json:**
   ```json
   {
     "name": "your-project",
     "assets": {
       "directory": "./dist"
     }
   }
   ```

### Custom domain not working

**Problem:** Site not accessible on custom domain.

**Solutions:**
1. **Check DNS propagation:**
   - May take up to 24 hours
   - Use `dig yourdomain.com`

2. **Verify Cloudflare setup:**
   - Domain added to Cloudflare
   - DNS records correct
   - SSL certificate issued

## Performance Issues

### Slow builds

**Problem:** Build takes too long.

**Solutions:**
1. **Check file sizes:**
   ```bash
   du -sh assets/
   ```

2. **Optimize images:**
   ```bash
   # Use image optimization tools
   ```

3. **Reduce template complexity**

### Large dist directory

**Problem:** Generated site too large.

**Solutions:**
1. **Clean dist before build:**
   ```bash
   rm -rf dist/
   stw build
   ```

2. **Exclude unnecessary files:**
   - Don't include source files
   - Minify assets

## Error Messages

### "no such file or directory"

**Context:** File system operations.

**Solutions:**
- Check file paths
- Verify files exist
- Check permissions

### "invalid YAML"

**Context:** Configuration parsing.

**Solutions:**
- Validate YAML syntax
- Check indentation
- Use spaces, not tabs

### "template: pattern matches no files"

**Context:** Template parsing.

**Solutions:**
- Check template file paths
- Verify files exist in templates/
- Check file permissions

### "port already in use"

**Context:** Server startup.

**Solutions:**
- Use different port: `stw serve --port 3000`
- Kill existing process
- Check for zombie processes

## Getting Help

### Debug Mode

Enable verbose logging:

```bash
# Set log level
export STW_LOG_LEVEL=debug
stw build
```

### System Information

Provide when reporting issues:

```bash
# System info
uname -a
go version
stw --version

# Project info
ls -la
head config.yaml
```

### Community Support

- Check existing GitHub issues
- Search documentation
- Ask in discussions

### Reporting Bugs

When reporting bugs, include:

1. **Steps to reproduce**
2. **Expected behavior**
3. **Actual behavior**
4. **Environment details**
5. **Error messages/logs**
6. **Minimal example** to reproduce

## Common Workflows

### Reset Project

```bash
# Clean everything
rm -rf dist/
rm -rf assets/ templates/ pages/
mkdir -p assets/css assets/js assets/images
mkdir -p templates/components templates/partials
mkdir -p pages

# Recreate basic files
# ... add your files
```

### Debug Template

```html
<!-- Add to template for debugging -->
<pre>{{printf "%+v" .}}</pre>
<pre>{{printf "%+v" .Config}}</pre>
<pre>{{printf "%+v" .Meta}}</pre>
```

### Test Configuration

```bash
# Validate config.yaml
go run -c 'import "gopkg.in/yaml.v3"; fmt.Println("Valid YAML")' < config.yaml
```