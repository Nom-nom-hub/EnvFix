# Documentation Index

Welcome to the envfix documentation. This directory contains comprehensive guides for users and developers.

## Quick Navigation

### 📖 User Guides
- [Getting Started](../QUICK_START.md) - Start here if you're new
- [Command Reference](../QUICK_REFERENCE.md) - All commands at a glance
- [README](../README.md) - Project overview

### 🚀 Feature Guides
- [Orphaned Environments](../ORPHANED_ENVIRONMENTS.md) - Detect and clean unused venvs
- [Dependency Management](../DEPENDENCIES.md) - Manage dependencies across languages
- [Version Management](../VERSION_MANAGEMENT.md) - Handle version conflicts and EOL

### 🔧 Developer Guides
- [Architecture](../ARCHITECTURE.md) - System design and overview
- [Development Guide](../DEVELOPMENT.md) - Contributing code
- [Testing Guide](../TESTING.md) - Running and writing tests

### 📝 Project Documentation
- [Changelog](../CHANGELOG.md) - Complete change history
- [Contributing](../CONTRIBUTING.md) - Contribution guidelines
- [License](../LICENSE) - MIT License

### 🎯 Release & Status
- [Release Notes](../V1_0_1_RELEASE_NOTES.md) - v1.0.1 release information
- [Release Checklist](../RELEASE_CHECKLIST.md) - QA verification checklist
- [Master Checklist](../MASTER_CHECKLIST.md) - Feature completion tracking
- [Deliverables](../DELIVERABLES.md) - What's included in the release

### 📚 Examples & Integrations
- [Basic Usage Examples](../examples/basic_usage.sh) - Example commands
- [CI/CD Integration](../examples/ci_cd_integration.md) - Pipeline integration
- [Docker Examples](../examples/ci_cd_integration.md#docker) - Container usage

---

## Documentation Structure

```
envfix/
├── docs/                          # Documentation directory
│   ├── README.md                  # This file
│   └── ... (guides listed above)
├── examples/                      # Example scripts and configs
│   ├── basic_usage.sh
│   └── ci_cd_integration.md
├── .github/
│   ├── workflows/                 # CI/CD pipelines
│   │   ├── ci.yml
│   │   └── release.yml
│   └── ISSUE_TEMPLATE/            # Issue templates
│       ├── bug_report.md
│       └── feature_request.md
├── README.md                      # Main overview
├── QUICK_START.md                 # Getting started
├── CONTRIBUTING.md                # How to contribute
├── CHANGELOG.md                   # Version history
└── [other docs]
```

---

## For Different Users

### I'm a New User
1. Start with [README.md](../README.md)
2. Follow [QUICK_START.md](../QUICK_START.md)
3. Check [QUICK_REFERENCE.md](../QUICK_REFERENCE.md)
4. Read feature guides as needed

### I'm a Developer
1. Read [ARCHITECTURE.md](../ARCHITECTURE.md)
2. Follow [DEVELOPMENT.md](../DEVELOPMENT.md)
3. Check [TESTING.md](../TESTING.md)
4. Review [CONTRIBUTING.md](../CONTRIBUTING.md)

### I'm Setting Up CI/CD
1. See [CI/CD Integration Guide](../examples/ci_cd_integration.md)
2. Review [.github/workflows/](.github/workflows/)
3. Check [examples/](../examples/)

### I Want to Report an Issue
1. Search existing issues first
2. Use [bug report template](.github/ISSUE_TEMPLATE/bug_report.md)
3. Provide reproduction steps
4. Include environment information

### I Want to Suggest a Feature
1. Use [feature request template](.github/ISSUE_TEMPLATE/feature_request.md)
2. Describe the problem it solves
3. Provide use case examples
4. Suggest implementation approach

---

## Key Documentation Files

| File | Purpose | Audience |
|------|---------|----------|
| README.md | Project overview | Everyone |
| QUICK_START.md | Getting started guide | New users |
| QUICK_REFERENCE.md | Command reference | All users |
| ARCHITECTURE.md | System design | Developers |
| DEVELOPMENT.md | Development setup | Contributors |
| CONTRIBUTING.md | Contribution guidelines | Contributors |
| TESTING.md | Testing approach | QA, Developers |
| VERSION_MANAGEMENT.md | Version handling | All users |
| DEPENDENCIES.md | Dependency management | All users |
| ORPHANED_ENVIRONMENTS.md | Orphaned cleanup | All users |
| CHANGELOG.md | Change history | All users |

---

## Getting Help

### Common Questions

**Q: How do I get started?**
A: Start with [QUICK_START.md](../QUICK_START.md)

**Q: What commands are available?**
A: Check [QUICK_REFERENCE.md](../QUICK_REFERENCE.md)

**Q: How do I set up development?**
A: See [DEVELOPMENT.md](../DEVELOPMENT.md)

**Q: How do I contribute?**
A: Read [CONTRIBUTING.md](../CONTRIBUTING.md)

**Q: How do I report a bug?**
A: Use the [bug report template](.github/ISSUE_TEMPLATE/bug_report.md)

**Q: How do I integrate with CI/CD?**
A: See [CI/CD Integration Guide](../examples/ci_cd_integration.md)

### Need More Help?

- **Search existing issues**: [GitHub Issues](https://github.com/aicode-dev/envfix/issues)
- **Read the guides**: Links above
- **Check examples**: [examples/](../examples/)
- **Review code**: [Source code](../cmd/, ../internal/)

---

## Documentation Maintenance

This documentation is maintained alongside the code. When reporting issues or suggesting improvements, please let us know!

### Keeping Documentation Fresh

- Documentation is reviewed with every code change
- Examples are tested and verified
- Links are checked regularly
- Version information is kept current

---

## Related Resources

- [GitHub Repository](https://github.com/aicode-dev/envfix)
- [Project Homepage](https://github.com/aicode-dev/envfix)
- [Release Page](https://github.com/aicode-dev/envfix/releases)
- [Issue Tracker](https://github.com/aicode-dev/envfix/issues)

---

**Last Updated**: December 6, 2025  
**Version**: 1.0.1
