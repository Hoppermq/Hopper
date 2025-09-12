# Contributing to Hopper 🐇

Thank you for your interest in contributing to Hopper! This guide will help you get started with the development workflow and contribution standards.

---

## 🚀 Quick Start

### Prerequisites
- **Go 1.23.4+** (check with `go version`)
- **Git** for version control
- **Make/Task** (we use [Taskfile](https://taskfile.dev/))

### Development Setup

```bash
# Clone the repository
git clone https://github.com/hoppermq/hopper.git
cd hopper

# Install dependencies
go mod download

# Run tests to verify setup
go test ./...

# Start development server
go run main.go
```

---

## 🏗️ Project Structure

```
hopper/
├── cmd/                    # Application entry points
├── internal/               # Private application code
│   ├── application/        # Application orchestration
│   ├── config/            # Configuration management  
│   ├── events/            # Event bus implementation
│   ├── http/              # HTTP server and APIs
│   ├── mq/                # Message queue core logic
│   │   ├── core/          # Broker and protocol implementation
│   │   └── transport/     # Transport layer (TCP)
│   ├── storage/           # Data persistence
│   └── ui/                # Web dashboard
├── pkg/                   # Public API packages
│   ├── client/            # Go client SDK
│   └── domain/            # Domain models and interfaces
├── scripts/               # Build and deployment scripts
└── taskfile.yml          # Task definitions
```

---

## 📋 Development Workflow

### 1. Pick an Issue
- Browse [open issues](https://github.com/hoppermq/hopper/issues) 
- Comment on issues you'd like to work on
- For new features, open an issue first to discuss the approach
- *Note: GitHub bot integration coming soon for Linear ↔ GitHub sync*

### 2. Branch Naming Convention
**Use Linear issue IDs for branch names:**

```bash
# Feature branches (Linear ID format)
git checkout -b feat/HOP-034/create-ingestor-service
git checkout -b feat/HOP-042/implement-producer-api

# Bug fixes
git checkout -b fix/HOP-051/resolve-connection-leak
git checkout -b fix/HOP-063/fix-memory-usage

# Documentation
git checkout -b docs/HOP-028/update-client-api-docs
git checkout -b docs/HOP-071/add-deployment-guide
```

**Branch naming format**: `{type}/{LINEAR-ID}/{short-description}`
- **Linear ID**: Use exact Linear issue ID (e.g., `HOP-034`)
- **Type**: `feat|fix|docs|test|refactor|chore`
- **Description**: Brief kebab-case description

### 3. Commit Messages
Follow [Conventional Commits](https://www.conventionalcommits.org/):

```bash
feat(client): add message persistence API
fix(broker): resolve memory leak in connection pool
docs(readme): update quickstart examples
test(core): add broker integration tests
```

### 4. Testing Standards

**Package-level Testing Strategy:**
- **Internal packages**: Use same-package tests (`package mq`) for implementation details
- **Public packages**: Use separate test packages (`package client_test`) for API contracts  
- **Integration tests**: Always include for user-facing features

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run integration tests
go test -tags=integration ./...
```

**Test Structure:**
```go
// Internal package test (implementation details)
package mq

func TestBroker_HandleConnection(t *testing.T) {
    // Test internal broker logic
}

// Public package test (API contract)  
package client_test

import "github.com/hoppermq/hopper/pkg/client"

func TestClient_Connect(t *testing.T) {
    // Test public client API
}
```

### 5. Code Quality Standards

**Code Formatting:**
```bash
# Format code
go fmt ./...

# Run linters
go vet ./...
golangci-lint run
```

**Documentation:**
- All public functions must have godoc comments
- Include examples in godoc when helpful
- Update README.md if adding user-facing features

---

## 🔧 Common Development Tasks

### Running Services
```bash
# Start Hopper server
go run main.go

# Run with development config  
APP_ENV=dev go run main.go

# Build binary
go build -o hopper main.go
```

### Working with Client SDK
```bash
# Test client examples
cd pkg/client/examples
go run main.go

# Run client tests
go test ./pkg/client/...
```

### Database/Storage Development
```bash
# Run integration tests with storage
go test -tags=integration ./internal/storage/...
```

---

## 🎯 Contribution Areas

### 🔥 High Priority
- **Client SDK completion** → Producer/Consumer APIs
- **Message persistence** → Durable message storage
- **Performance optimization** → Benchmarking and tuning
- **Integration tests** → End-to-end testing scenarios

### 🚀 New Features
- **Multi-language clients** → Python, Node.js, Java SDKs
- **Advanced routing** → Topic patterns and filtering
- **Streamly integration** → Monitoring dashboard connectivity
- **Docker distribution** → Official container images

### 🐛 Bug Reports
- **Create in Linear** for internal project tracking
  - Include Go version, OS, and steps to reproduce
  - Add relevant logs and error messages
- Check existing Linear issues before creating duplicates
- *GitHub Issues may also be used for community reports*

---

## 📝 Pull Request Process

### Before Submitting
- [ ] All tests pass (`go test ./...`)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] No linting errors (`go vet ./...`)
- [ ] Documentation updated (if applicable)
- [ ] Integration tests added (for user-facing features)

### PR Template
```markdown
## Description
Brief description of changes and motivation.

## Type of Change
- [ ] Bug fix
- [ ] New feature  
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added (if applicable)
- [ ] Manual testing performed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests pass locally
```

### Review Process
1. **Automated checks** → CI pipeline runs tests and linting
2. **Maintainer review** → Code quality and architectural alignment
3. **Community feedback** → Additional input from contributors  
4. **Merge approval** → Final approval from maintainers

---

## 🤝 Community Guidelines

### Code of Conduct
- Be respectful and inclusive in all interactions
- Focus on constructive feedback and collaboration
- Help newcomers get started with the project

### Communication Channels
- **Linear** → Primary issue tracking and project management
- **GitHub Issues** → Community bug reports and feature requests
- **GitHub Discussions** → Architecture questions and ideas
- **Pull Requests** → Code review and collaboration

### Getting Help
- Check existing documentation and Linear/GitHub issues first
- Ask questions in GitHub Discussions for community support
- Tag maintainers for urgent issues: `@hoppermq/maintainers`
- *Linear access for core contributors*

---

## 📊 Project Governance

### Maintainers
- Review and approve pull requests
- Guide architectural decisions  
- Manage releases and roadmap

### Contributors
- Anyone who submits pull requests
- Recognition in project contributors list
- Opportunity to become maintainers based on contributions

---

## 🎉 Recognition

Contributors are recognized in:
- GitHub contributors graph
- Release notes acknowledgments  
- Project documentation credits

**Thank you for contributing to Hopper! Every contribution helps make message brokers more transparent and developer-friendly.** 🚀
