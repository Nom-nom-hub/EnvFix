# envfix Integration with CI/CD

This guide shows how to integrate envfix into your CI/CD pipelines.

## GitHub Actions

### Basic CI Check

```yaml
name: Check Environment

on: [push, pull_request]

jobs:
  envcheck:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Build envfix
        run: go build -o envfix ./cmd
      
      - name: Check environment
        run: ./envfix scan
      
      - name: Validate versions
        run: ./envfix repair --dry-run
```

### With Artifact Upload

```yaml
- name: Generate scan report
  run: ./envfix scan --json > scan-report.json

- name: Upload report
  uses: actions/upload-artifact@v3
  with:
    name: environment-report
    path: scan-report.json
```

## GitLab CI

```yaml
check-environment:
  stage: setup
  image: golang:1.21
  script:
    - go build -o envfix ./cmd
    - ./envfix scan
    - ./envfix repair --dry-run
  artifacts:
    reports:
      dotenv: env-report.txt
```

## Jenkins

```groovy
pipeline {
    agent any
    
    stages {
        stage('Setup') {
            steps {
                sh 'go build -o envfix ./cmd'
            }
        }
        
        stage('Check Environment') {
            steps {
                sh './envfix scan'
                sh './envfix repair --dry-run'
            }
        }
        
        stage('Generate Report') {
            steps {
                sh './envfix scan --json > env-report.json'
                archiveArtifacts artifacts: 'env-report.json'
            }
        }
    }
}
```

## Docker

### Dockerfile

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app
COPY . .

RUN go build -o envfix ./cmd

# Run environment check
RUN ./envfix scan

# Continue with application build
FROM alpine:latest
COPY --from=0 /app/envfix /usr/local/bin/
```

### Docker Compose

```yaml
services:
  app:
    build: .
    environment:
      - PATH=/usr/local/go/bin:$PATH
    command: |
      sh -c "
        envfix scan &&
        envfix repair &&
        envfix clean &&
        npm test
      "
```

## CircleCI

```yaml
version: 2.1

jobs:
  setup:
    docker:
      - image: cimg/go:1.21
    steps:
      - checkout
      - run:
          name: Build envfix
          command: go build -o envfix ./cmd
      - run:
          name: Check environment
          command: ./envfix scan
      - run:
          name: Validate versions
          command: ./envfix repair --dry-run
      - store_artifacts:
          path: scan-report.json
```

## Travis CI

```yaml
language: go
go:
  - "1.21"

before_install:
  - go build -o envfix ./cmd

script:
  - ./envfix scan
  - ./envfix repair --dry-run
  - npm test
```

## Script Examples

### Check and Fix Environment

```bash
#!/bin/bash
set -e

echo "Checking environment..."
envfix scan

if [ $? -ne 0 ]; then
    echo "Issues found, attempting to fix..."
    envfix repair
    envfix clean
fi

echo "Verifying fixes..."
envfix scan

echo "Environment check passed!"
```

### Generate Report

```bash
#!/bin/bash

echo "Generating environment report..."
envfix scan --json > env-report.json

echo "Environment Report"
echo "=================="
cat env-report.json | jq .

# Archive for later
gzip env-report.json
mv env-report.json.gz artifacts/
```

### Conditional Repair

```bash
#!/bin/bash

echo "Checking for issues..."
ISSUES=$(envfix scan --json | jq '.issues | length')

if [ "$ISSUES" -gt 0 ]; then
    echo "Found $ISSUES issues"
    envfix repair --yes
    envfix clean
else
    echo "No issues found!"
fi
```

## Docker Compose Example

```yaml
version: '3.8'

services:
  check-env:
    image: golang:1.21
    working_dir: /app
    volumes:
      - .:/app
    command: |
      sh -c "
        echo 'Building envfix...' &&
        go build -o envfix ./cmd &&
        echo 'Scanning environment...' &&
        ./envfix scan &&
        echo 'Validating repairs...' &&
        ./envfix repair --dry-run &&
        echo 'Environment check passed!'
      "

  app:
    depends_on:
      - check-env
    image: node:18-alpine
    working_dir: /app
    volumes:
      - .:/app
    command: npm test
```

## Best Practices

### 1. Use Dry-Run for CI
```bash
# Preview changes without applying
envfix repair --dry-run
```

### 2. Save Reports
```bash
# Generate JSON for analysis
envfix scan --json > env-report.json
```

### 3. Fail Fast
```bash
# Fail if critical issues found
envfix scan --json | jq '.issues[] | select(.severity == "critical")'
```

### 4. Cache Dependencies
```yaml
- uses: actions/cache@v3
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
```

### 5. Multi-Platform Testing
```yaml
strategy:
  matrix:
    os: [ubuntu-latest, windows-latest, macos-latest]
    go-version: ['1.21', '1.25']
```

## Troubleshooting

### Build Fails on First Run
```bash
# Clean cache and retry
envfix clean
envfix repair
```

### Version Mismatch in CI
```bash
# Ensure consistent versions
envfix scan --json | jq '.python.version'
```

### Dependency Issues
```bash
# Detailed dependency check
envfix explain nodejs_missing
```

## Integration Patterns

### Pattern 1: Pre-Build Check
```bash
envfix scan || exit 1  # Fail if issues
go build ...
```

### Pattern 2: Pre-Test Prep
```bash
envfix repair --yes  # Fix everything
envfix clean         # Clean caches
npm test            # Run tests
```

### Pattern 3: Reporting
```bash
envfix scan --json > report.json
curl -X POST -d @report.json https://monitoring.example.com
```

### Pattern 4: Version Lock
```bash
envfix lock --output env.yaml
git add env.yaml
```

---

For more examples, check the main documentation.
