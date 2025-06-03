# ?? Contributing to GOLINE Looking Glass

Thank you for your interest in contributing to GOLINE Looking Glass! This document provides guidelines and information for contributors.

## ?? Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contributing Guidelines](#contributing-guidelines)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)

## Code of Conduct

This project adheres to a code of conduct that we expect all contributors to follow:

- **Be respectful** and inclusive in your communications
- **Be collaborative** and help others learn and grow
- **Be constructive** in your feedback and criticism
- **Be patient** with newcomers and those learning
- **Focus on what's best** for the project and community

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Basic understanding of networking concepts
- Familiarity with SSH and router command-line interfaces

### Areas for Contribution

We welcome contributions in several areas:

1. **?? Bug Fixes** - Help us identify and fix issues
2. **? New Features** - Add support for new router vendors or commands
3. **?? Documentation** - Improve guides, examples, and API docs
4. **?? Testing** - Add test cases and improve test coverage
5. **?? UI/UX** - Enhance the web interface and user experience
6. **?? Security** - Identify and fix security vulnerabilities
7. **? Performance** - Optimize performance and resource usage

## Development Setup

### 1. Fork and Clone

```bash
# Fork the repository on GitHub, then clone your fork
git clone https://github.com/yourusername/goline-looking-glass.git
cd goline-looking-glass

# Add upstream remote
git remote add upstream https://github.com/paolokappa/goline-looking-glass.git
```

### 2. Setup Development Environment

```bash
# Install dependencies
make deps

# Setup development configuration
make setup

# Build the application
make build

# Run tests
make test
```

### 3. Development Workflow

```bash
# Create a feature branch
git checkout -b feature/your-feature-name

# Make your changes
# ... edit files ...

# Test your changes
make test
make lint

# Build and test locally
make build
./build/looking-glass

# Commit your changes
git add .
git commit -m "Add: descriptive commit message"

# Push to your fork
git push origin feature/your-feature-name
```

## Contributing Guidelines

### Types of Contributions

#### ?? Bug Reports

When filing a bug report, please include:

- **Clear description** of the issue
- **Steps to reproduce** the problem
- **Expected vs actual behavior**
- **Environment details** (OS, Go version, router type)
- **Log output** (if applicable)
- **Screenshots** (for UI issues)

#### ? Feature Requests

For feature requests, please provide:

- **Use case description** - Why is this feature needed?
- **Proposed solution** - How should it work?
- **Alternative solutions** - Other approaches considered
- **Impact assessment** - Who would benefit from this feature?

#### ?? Code Contributions

Before starting work on a significant feature:

1. **Check existing issues** - Make sure it's not already being worked on
2. **Create an issue** - Discuss the feature with maintainers
3. **Get approval** - Wait for feedback before starting implementation
4. **Follow standards** - Use our coding conventions and patterns

### Router Vendor Support

When adding support for a new router vendor:

1. **Research commands** - Document the CLI commands needed
2. **Test thoroughly** - Verify on actual hardware/simulation
3. **Add documentation** - Update router configuration guide
4. **Provide examples** - Include sample configurations
5. **Add tests** - Create unit tests for the new vendor

### Security Considerations

When contributing security-related changes:

- **Private disclosure** - Report vulnerabilities privately first
- **Minimal exposure** - Don't include sensitive data in commits
- **Testing** - Verify security fixes don't break functionality
- **Documentation** - Update security guides as needed

## Pull Request Process

### 1. Pre-submission Checklist

Before submitting a pull request:

- [ ] Code follows project style guidelines
- [ ] Tests pass (`make test`)
- [ ] Linting passes (`make lint`)
- [ ] Documentation is updated
- [ ] Commit messages are descriptive
- [ ] Changes are in a feature branch
- [ ] No sensitive data is included

### 2. Pull Request Template

Use this template for your pull request description:

```markdown
## Description
Brief description of changes made.

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Router Testing
If applicable:
- [ ] Tested on Juniper devices
- [ ] Tested on Huawei devices  
- [ ] Tested on Cisco devices

## Documentation
- [ ] Code comments updated
- [ ] README updated
- [ ] API documentation updated
- [ ] Configuration examples updated

## Screenshots
If applicable, add screenshots of UI changes.

## Additional Notes
Any additional information or context.
```

### 3. Review Process

1. **Automated checks** - CI/CD pipeline runs tests and checks
2. **Code review** - Maintainers review the changes
3. **Feedback** - Address any comments or requested changes
4. **Approval** - Get approval from at least one maintainer
5. **Merge** - Changes are merged into main branch

## Coding Standards

### Go Style Guide

Follow standard Go conventions:

- Use `go fmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use meaningful variable and function names
- Write clear, concise comments
- Handle errors appropriately

### Commit Messages

Use clear, descriptive commit messages:

```
Type: Brief description (50 chars or less)

More detailed explanation if needed. Wrap at 72 characters.
Explain what and why, not how.

- Use bullet points for multiple changes
- Reference issues: Fixes #123
- Include breaking change notes if applicable
```

Types:
- `Add:` New feature or functionality
- `Fix:` Bug fix
- `Update:` Modification to existing feature
- `Remove:` Deletion of feature or code
- `Docs:` Documentation changes
- `Test:` Adding or updating tests
- `Refactor:` Code restructuring without behavior change

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test package
go test ./internal/router/...

# Run tests verbosely
go test -v ./...
```

### Writing Tests

- **Unit tests** - Test individual functions and methods
- **Integration tests** - Test component interactions
- **Router tests** - Test actual router communications (when possible)
- **API tests** - Test HTTP endpoints and responses

## Documentation

### Types of Documentation

1. **Code comments** - Explain complex logic and decisions
2. **API documentation** - Document REST endpoints and parameters
3. **User guides** - Help users configure and use the system
4. **Developer guides** - Help contributors understand the codebase
5. **Configuration examples** - Provide working configuration samples

### Documentation Standards

- **Clear and concise** - Easy to understand
- **Up to date** - Keep documentation current with code changes
- **Examples** - Include practical examples
- **Screenshots** - Use images for UI-related documentation
- **Links** - Reference related documentation sections

## Getting Help

If you need help with contributing:

- ?? **Email**: [noc@goline.ch](mailto:noc@goline.ch)
- ?? **Discussions**: [GitHub Discussions](https://github.com/paolokappa/goline-looking-glass/discussions)
- ?? **Issues**: [GitHub Issues](https://github.com/paolokappa/goline-looking-glass/issues)

## Recognition

Contributors are recognized in several ways:

- **GitHub contributors list** - Automatic recognition
- **CHANGELOG.md** - Major contributions mentioned
- **Documentation credits** - Contributor acknowledgments
- **Social media** - Appreciation posts for significant contributions

## Thank You!

Your contributions help make GOLINE Looking Glass better for everyone. Whether you're fixing a typo, adding a feature, or helping with documentation, every contribution is valuable and appreciated.

---

**Happy Contributing!** ??

*The GOLINE Looking Glass Team*  
*Paolo Caparrelli - GOLINE SA*
