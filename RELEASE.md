# Release Guide

This guide explains how to release new versions of the Ghion Go SDK.

## Versioning

We follow [Semantic Versioning](https://semver.org/):
- **MAJOR** (X.0.0): Breaking changes
- **MINOR** (0.X.0): New features, backward compatible
- **PATCH** (0.0.X): Bug fixes, backward compatible

Example: `v1.2.3` (Major 1, Minor 2, Patch 3)

## Pre-Release Checklist

Before releasing a new version:

1. **Update Tests**
   ```bash
   go test ./pkg/... -v
   go test ./tests/ -v -tags=integration
   ```

2. **Check Coverage**
   ```bash
   go test ./pkg/... -coverprofile=coverage.out -covermode=atomic
   go tool cover -func=coverage.out
   ```

3. **Update Documentation**
   - Update `README.md` with new features
   - Update `CHANGELOG.md` with release notes
   - Update examples if needed

4. **Run Linter**
   ```bash
   go fmt ./...
   golangci-lint run
   ```

5. **Update Dependencies**
   ```bash
   go mod tidy
   go mod verify
   ```

## Creating a Release

### Step 1: Update Version in Files

Update version references in:
- `README.md` (if version is mentioned)
- `CHANGELOG.md` (add new release notes)
- Any other documentation files

### Step 2: Commit Changes

```bash
git add .
git commit -m "Release v1.0.1: Add new feature"
```

### Step 3: Create Git Tag

```bash
# Annotated tag with release notes
git tag -a v1.0.1 -m "Release v1.0.1: Add new feature

- Added support for new payment channel
- Fixed bug in webhook signature verification
- Updated documentation"
```

### Step 4: Push to GitHub

```bash
# Push code
git push origin main

# Push tag
git push origin v1.0.1
```

### Step 5: Create GitHub Release (Optional)

1. Go to GitHub repository
2. Click "Releases" → "Create a new release"
3. Select the tag you just pushed
4. Add release title and description
5. Publish the release

## Post-Release Tasks

1. **Verify Installation**
   ```bash
   # Test that users can install the new version
   go get github.com/yayawallet/ghion-go-sdk@v1.0.1
   ```

2. **Update Go Module Proxy**
   - The Go module proxy (proxy.golang.org) will automatically fetch your release
   - No manual action needed

3. **Announce Release**
   - Update documentation if needed
   - Notify users of breaking changes (if any)

## Changelog Format

Maintain `CHANGELOG.md` with this format:

```markdown
## [1.0.1] - 2024-01-15

### Added
- Support for new payment channel
- New webhook event type

### Fixed
- Bug in webhook signature verification
- Memory leak in retry logic

### Changed
- Improved error messages for validation failures
- Updated dependencies

### Deprecated
- Old payment method (will be removed in v2.0.0)
```

## Version Bumping Examples

### Patch Release (Bug Fix)
```bash
git tag -a v1.0.1 -m "Release v1.0.1: Fix webhook signature bug"
git push origin v1.0.1
```

### Minor Release (New Feature)
```bash
git tag -a v1.1.0 -m "Release v1.1.0: Add Telebirr support"
git push origin v1.1.0
```

### Major Release (Breaking Change)
```bash
git tag -a v2.0.0 -m "Release v2.0.0: Breaking API changes"
git push origin v2.0.0
```

## Pre-Release Versions

For beta or alpha releases:

```bash
# Pre-release version
git tag -a v1.1.0-beta.1 -m "Beta release v1.1.0-beta.1"
git push origin v1.1.0-beta.1

# Users can install pre-release
go get github.com/yayawallet/ghion-go-sdk@v1.1.0-beta.1
```

## Rolling Back a Release

If you need to revert a release:

```bash
# Delete the tag locally
git tag -d v1.0.1

# Delete the tag from GitHub
git push origin :refs/tags/v1.0.1

# Create a new tag for the previous version
git tag -a v1.0.0 -m "Revert to v1.0.0"
git push origin v1.0.0
```

## Common Issues

### Tag Already Exists
```bash
# Delete existing tag first
git tag -d v1.0.1
git push origin :refs/tags/v1.0.1
# Then create new tag
```

### Module Not Found
- Wait a few minutes for Go module proxy to sync
- Verify tag is pushed: `git ls-remote --tags origin`
- Check GitHub repository settings (must be public)

## Summary

1. Run tests and check coverage
2. Update documentation and changelog
3. Commit changes
4. Create annotated git tag
5. Push code and tags to GitHub
6. Verify installation works
7. Create GitHub release (optional)
