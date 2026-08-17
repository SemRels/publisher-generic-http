## v0.3.1 (2026-08-17)

### Bug Fixes

* remove commit_changelog:false with no replacement plugin - CHANGELOG.md was never updated
* resolve OCI digests via jq instead of oras tag lookup (Docker rewrites bare tags to :latest)

### Other Changes

* **schema:** add JSON Schema for editor validation (#6)
* sync from plugin-template@e9e7179 -- .github/workflows/ci.yml .github/workflows/release.yml .github/workflows/security.yml Dockerfile GOVERNANCE.md Makefile SECURITY.md (#15)

