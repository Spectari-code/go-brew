# Contributing to Go Brew CLI

Thank you for your interest in contributing to Go Brew CLI! This document provides guidelines and information to help you get started with contributing to this project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)
- [Pull Request Process](#pull-request-process)
- [Release Process](#release-process)

## Code of Conduct

This project follows the [Go Community Code of Conduct](https://golang.org/conduct). Please read and follow these guidelines to ensure a welcoming environment for all contributors.

## Development Setup

### Prerequisites

- Go 1.23 or later
- Git
- Make (optional, for build automation)

### Getting Started

1. **Fork the Repository**
   ```bash
   # Fork the repository on GitHub, then clone your fork
   git clone https://github.com/YOUR_USERNAME/go-brew-cli.git
   cd go-brew-cli
   ```

2. **Add Upstream Remote**
   ```bash
   git remote add upstream https://github.com/Spectari-code/go-brew-cli.git
   ```

3. **Install Dependencies**
   ```bash
   go mod tidy
   ```

4. **Verify Setup**
   ```bash
   go run .
   go test
   ```

### Development Workflow

```bash
# Create a new branch for your feature
git checkout -b feature/your-feature-name

# Make your changes
# ... edit files ...

# Run tests
go test -v

# Format code
go fmt ./...

# Run linter (optional)
go vet ./...

# Commit your changes
git add .
git commit -m "feat: add your feature description"

# Push to your fork
git push origin feature/your-feature-name
```

## How to Contribute

### Types of Contributions

We welcome the following types of contributions:

1. **Bug Fixes** - Found and fixed issues
2. **New Features** - Enhancements to existing functionality
3. **Documentation** - Improved documentation and examples
4. **Performance Improvements** - Optimizations and efficiency gains
5. **Test Coverage** - Additional tests for existing code
6. **Refactoring** - Code cleanup and architectural improvements

### Reporting Issues

When filing bug reports, please include:

- **Clear Description**: What the bug is and how to reproduce it
- **Steps to Reproduce**: Detailed steps to trigger the bug
- **Expected Behavior**: What you expected to happen
- **Actual Behavior**: What actually happened
- **Environment**: Go version, OS, terminal information
- **Additional Context**: Any relevant screenshots or logs

### Feature Requests

For feature requests:

- **Use Case**: Describe the problem you're trying to solve
- **Proposed Solution**: How you envision the feature working
- **Alternatives**: Other approaches you've considered
- **Additional Context**: Any relevant information

## Coding Standards

### Go Guidelines

Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go.html).

### Project-Specific Standards

1. **Package Comments**: Every package should have a godoc comment explaining its purpose
2. **Exported Functions**: All exported functions must have godoc comments
3. **Error Handling**: Use explicit error handling with proper context
4. **Naming**: Use clear, descriptive names following Go conventions
5. **Formatting**: Use `gofmt` for consistent formatting

### Code Style Examples

```go
// Good: Clear function with proper documentation
// StartTimer begins the countdown timer with the specified duration.
// It returns any error encountered during timer initialization.
func StartTimer(duration time.Duration) error {
    if duration <= 0 {
        return errors.New("duration must be positive")
    }
    // implementation...
    return nil
}

// Bad: Unclear naming, missing documentation
func start(d int) {
    // unclear implementation
}
```

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```
feat(ui): add progress bar animation
fix(audio): handle missing mp3 file gracefully
docs(readme): update installation instructions
```

## Testing Guidelines

### Test Coverage

- Aim for >80% test coverage for new code
- Test both happy paths and error cases
- Test edge cases and boundary conditions

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        want     OutputType
        wantErr  bool
    }{
        {
            name:    "valid input",
            input:   validInput,
            want:    expectedOutput,
            wantErr: false,
        },
        {
            name:    "invalid input",
            input:   invalidInput,
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("FunctionName() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Running Tests

```bash
# Run all tests
go test

# Run with verbose output
go test -v

# Run with coverage
go test -cover

# Run with coverage report
go test -coverprofile=coverage.out && go tool cover -html=coverage.out

# Run specific test
go test -run TestSpecificFunction

# Run benchmarks
go test -bench=.
```

## Documentation

### Code Documentation

- All exported functions, types, and methods need godoc comments
- Include parameter descriptions, return values, and usage examples
- Explain design decisions and trade-offs when relevant

### Example Documentation

```go
// TeaPreset represents a pre-configured tea brewing setting with all necessary
// information for proper tea preparation. Each preset includes brew time,
// recommended temperature, and helpful notes for the best results.
//
// Example:
//     preset := TeaPreset{
//         Name:     "Green Tea",
//         Duration: 2 * time.Minute,
//         Temp:     "80°C",
//         Notes:    "Don't overbrew to avoid bitterness",
//     }
type TeaPreset struct {
    Name     string        // Human-readable name of the tea type
    Duration time.Duration // Recommended brewing time
    Temp     string        // Recommended water temperature
    Notes    string        // Additional brewing notes or tips
}
```

## Pull Request Process

### Before Submitting

1. **Test Your Changes**
   ```bash
   go test -v
   go build .
   ```

2. **Format Code**
   ```bash
   go fmt ./...
   ```

3. **Check for Issues**
   ```bash
   go vet ./...
   ```

4. **Update Documentation**
   - Update README.md if needed
   - Add godoc comments for new functions
   - Update or add tests

### Pull Request Template

```markdown
## Description
Brief description of changes made.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Tests pass locally
- [ ] New tests added for new functionality
- [ ] Manual testing completed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No merge conflicts
```

### Review Process

1. **Automated Checks**: CI will run tests and linters
2. **Code Review**: Maintainers will review your changes
3. **Approval**: At least one maintainer approval required
4. **Merge**: Changes will be merged after approval

## Release Process

### Versioning

This project follows [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Checklist

1. **Update Version**
   - Update version constants
   - Update CHANGELOG.md

2. **Create Release**
   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```

3. **Build Binaries**
   ```bash
   goreleaser build --clean
   ```

## Getting Help

- **GitHub Issues**: For bug reports and feature requests
- **GitHub Discussions**: For questions and community discussion
- **Maintainers**: Tag maintainers in relevant issues or discussions

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Bubbletea Documentation](https://github.com/charmbracelet/bubbletea)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

---

Thank you for contributing to Go Brew CLI! Your contributions help make this project better for everyone.