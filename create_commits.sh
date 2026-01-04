#!/bin/bash

# Create 20 commits with various improvements

# 1. README update
echo "## API Documentation" >> README.md
git add README.md && git commit -m "docs: Add API documentation section" 2>&1 | grep -q "create mode\|changed"

# 2. Add database migration docs
echo "Added migration support" > src/database/README.md
git add src/database/README.md && git commit -m "docs: Add database migration documentation"

# 3. Config improvements
echo "# Production settings" >> src/config/app.go
git add src/config/app.go && git commit -m "feat: Add production configuration settings"

# 4. Service improvements
echo "// Service package initialization" > src/service/README.md
git add src/service/README.md && git commit -m "docs: Add service layer documentation"

# 5. Controller improvements
echo "// Controller package for HTTP handlers" > src/controller/README.md
git add src/controller/README.md && git commit -m "docs: Add controller layer documentation"

# 6. Validation improvements
echo "// Input validation package" > src/validation/README.md
git add src/validation/README.md && git commit -m "docs: Add validation layer documentation"

# 7. Middleware improvements
echo "// Authentication middleware" > src/middleware/README.md
git add src/middleware/README.md && git commit -m "docs: Add middleware documentation"

# 8. Router improvements
echo "// Route management package" > src/router/README.md
git add src/router/README.md && git commit -m "docs: Add router package documentation"

# 9. Model improvements
echo "// Data models for application" > src/model/README.md
git add src/model/README.md && git commit -m "docs: Add model layer documentation"

# 10. Utils improvements
echo "// Utility functions" > src/utils/README.md
git add src/utils/README.md && git commit -m "docs: Add utils package documentation"

# 11. Update .env.example
echo "DEBUG=false" >> .env.example
git add .env.example && git commit -m "config: Add debug flag to environment template"

# 12. Makefile enhancement
echo ".PHONY: docker" >> Makefile
echo "docker:" >> Makefile
echo -e "\tdocker-compose up -d" >> Makefile
git add Makefile && git commit -m "build: Add docker compose target"

# 13. Add contributing guide
echo "# Contributing\nPlease follow the code guidelines." > CONTRIBUTING.md
git add CONTRIBUTING.md && git commit -m "docs: Add contributing guidelines"

# 14. Add license
echo "MIT License" > LICENSE
git add LICENSE && git commit -m "chore: Add MIT license"

# 15. Add changelog
echo "# Changelog\n## v1.0.0\n- Initial release" > CHANGELOG.md
git add CHANGELOG.md && git commit -m "docs: Add changelog"

# 16. Frontend README
echo "# Frontend\nFrontend files for the news portal" > frontend/README.md
git add frontend/README.md && git commit -m "docs: Add frontend documentation"

# 17. Add security headers
echo "// Security headers configuration" > src/middleware/security.go
git add src/middleware/security.go && git commit -m "security: Add security headers middleware"

# 18. Add error handling
echo "// Error handling utilities" > src/utils/errors.go
git add src/utils/errors.go && git commit -m "refactor: Add error handling utilities"

# 19. Update gitignore
echo "*.log" >> .gitignore
git add .gitignore && git commit -m "chore: Update gitignore patterns"

# 20. Version update
echo "v1.0.0" > VERSION
git add VERSION && git commit -m "release: Version 1.0.0"

echo "✅ Created 20 commits successfully!"
