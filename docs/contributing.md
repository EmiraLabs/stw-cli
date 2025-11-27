# Contributing

We welcome contributions to stw-cli! This document outlines the process for contributing to the project.

## Getting Started

### Development Setup

1. **Fork and clone the repository:**
   ```bash
   git clone https://github.com/EmiraLabs/stw-cli.git
   cd stw-cli
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Build the project:**
   ```bash
   go build -o stw ./cmd/stw
   ```

4. **Run tests:**
   ```bash
   go test ./...
   ```

### Development Workflow

1. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**

3. **Run tests and linting:**
   ```bash
   go test ./...
   go vet ./...
   ```

4. **Build and test manually:**
   ```bash
   go build -o stw ./cmd/stw
   ./stw --help
   ```

5. **Commit your changes:**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

6. **Push and create PR:**
   ```bash
   git push origin feature/your-feature-name
   ```

## Code Style

### Go Code

- Follow standard Go formatting: `go fmt`
- Use `go vet` for static analysis
- Follow Go naming conventions
- Write comprehensive tests
- Add documentation comments for exported functions

### Commit Messages

Use conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Testing
- `chore`: Maintenance

Examples:
```
feat: add SEO meta support
fix: correct template parsing error
docs: update installation guide
```

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/meta

# Run with coverage
go test -cover ./...

# Run with verbose output
go test -v ./...
```

### Writing Tests

- Use table-driven tests for multiple test cases
- Test both success and error cases
- Mock external dependencies
- Use descriptive test names

Example:

```go
func TestParseFrontMatter(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        wantMeta meta.Meta
        wantBody string
        wantErr  bool
    }{
        {
            name:  "valid YAML front matter",
            input: "---\ntitle: Test\n---\nContent",
            wantMeta: meta.Meta{Title: "Test"},
            wantBody: "Content",
            wantErr:  false,
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            gotMeta, gotBody, err := meta.ParseFrontMatter(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseFrontMatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(gotMeta, tt.wantMeta) {
                t.Errorf("ParseFrontMatter() gotMeta = %v, want %v", gotMeta, tt.wantMeta)
            }
            if gotBody != tt.wantBody {
                t.Errorf("ParseFrontMatter() gotBody = %v, want %v", gotBody, tt.wantBody)
            }
        })
    }
}
```

## Documentation

### Code Documentation

- Add package comments
- Document exported functions, types, and methods
- Use examples in documentation

### User Documentation

- Update docs in `docs/` directory
- Keep README.md current
- Add examples for new features

## Pull Request Process

### Before Submitting

1. **Update tests** for any changed functionality
2. **Update documentation** if needed
3. **Ensure CI passes** locally
4. **Squash commits** if needed
5. **Write a clear PR description**

### PR Template

```
## Description
Brief description of the changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Manual testing performed
- [ ] CI passes

## Checklist
- [ ] Code follows Go conventions
- [ ] Tests pass
- [ ] Documentation updated
- [ ] Commit messages follow conventions
```

### Review Process

1. **Automated checks** run on PR
2. **Code review** by maintainers
3. **Approval** and merge
4. **Release** if needed

## Architecture

### Package Structure

```
cmd/stw/           # CLI application
internal/
  application/     # Application services
  domain/          # Domain models
  infrastructure/  # Infrastructure implementations
  meta/            # Metadata handling
docs/              # Documentation
```

### Design Principles

- **Clean Architecture**: Separation of concerns
- **Dependency Injection**: Interfaces for testability
- **SOLID Principles**: Single responsibility, open/closed, etc.
- **Error Handling**: Explicit error returns
- **Testing**: Comprehensive test coverage

## Issue Reporting

### Bug Reports

Use the bug report template:

```
## Description
Clear description of the bug

## Steps to Reproduce
1. Step 1
2. Step 2
3. Step 3

## Expected Behavior
What should happen

## Actual Behavior
What actually happens

## Environment
- OS: [e.g., macOS 12.0]
- Go version: [e.g., 1.19]
- stw-cli version: [e.g., v1.0.0]
```

### Feature Requests

Use the feature request template:

```
## Problem
What's the problem this feature would solve?

## Solution
Describe the proposed solution

## Alternatives
Any alternative solutions considered?

## Additional Context
Any other context or screenshots
```

## Release Process

### Versioning

Follows [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features
- **PATCH**: Bug fixes

### Release Checklist

- [ ] Update version in code
- [ ] Update CHANGELOG.md
- [ ] Create git tag
- [ ] Create GitHub release
- [ ] Update documentation
- [ ] Announce release

## Community

### Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help newcomers learn

### Getting Help

- Check existing issues and documentation
- Ask questions in discussions
- Join our community chat

## Recognition

Contributors are recognized in:
- CHANGELOG.md for releases
- GitHub contributors list
- Release notes

Thank you for contributing to stw-cli! 🎉