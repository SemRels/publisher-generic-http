# publisher-generic-http

[![Latest Release](https://img.shields.io/github/v/release/SemRels/publisher-generic-http?label=version&color=blue)](https://github.com/SemRels/publisher-generic-http/releases/latest)
[![CI](https://github.com/SemRels/publisher-generic-http/actions/workflows/ci.yml/badge.svg)](https://github.com/SemRels/publisher-generic-http/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/SemRels/publisher-generic-http)](LICENSE)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/SemRels/publisher-generic-http/badge)](https://scorecard.dev/viewer/?uri=github.com/SemRels/publisher-generic-http)
[![REUSE status](https://api.reuse.software/badge/github.com/SemRels/publisher-generic-http)](https://api.reuse.software/info/github.com/SemRels/publisher-generic-http)

Generic HTTP publisher plugin for semrel. It uploads release artifacts to custom endpoints via HTTP PUT/POST.

## Repository Layout

```text
cmd/plugin/              Plugin entry point
internal/plugin/         Business logic scaffold
internal/grpc/           gRPC transport scaffold
proto/v1                 Symlink to the SemRel protobuf contract
.github/workflows/       CI, release, and security automation
```

## Installation

Published binaries are distributed through releases and synchronized to `registry.semrel.io`.

## Development

```bash
go build ./cmd/plugin
go test ./...
```

## Configuration

Configure in `.semrel.yaml`:

```yaml
plugins:
  - uses: "@semrel/publisher-generic-http"
    phase: release
    args:
      url: https://uploads.example.com/releases/{version}/{artifact}
      method: PUT
      token: ${UPLOAD_TOKEN}
      artifacts: dist/myapp-linux-amd64,dist/myapp-linux-arm64
```

## Options

Every option can be set as an entry under `args` in `.semrel.yaml`. semrel passes each entry to the plugin as the environment variable shown in the table, so you can also set that variable directly (for example in your CI environment) instead of using `args`.

| Argument | Environment variable | Required | Description | Default |
| --- | --- | --- | --- | --- |
| `url` | `SEMREL_PLUGIN_URL` | Required | Upload URL. Supports the `{version}` and `{artifact}` tokens. | None |
| `method` | `SEMREL_PLUGIN_METHOD` | Optional | HTTP method used for the upload. | `PUT` |
| `token` | `SEMREL_PLUGIN_TOKEN` | Optional | Bearer token sent with each request. | None |
| `headers_json` | `SEMREL_PLUGIN_HEADERS_JSON` | Optional | JSON object with custom request headers. | None |
| `artifacts` | `SEMREL_PLUGIN_ARTIFACTS` | Optional | Comma-separated list of artifacts to upload. | None |
| `artifacts_json` | `SEMREL_PLUGIN_ARTIFACTS_JSON` | Optional | JSON array of artifacts, used as a fallback for `artifacts`. | None |
| `artifact` | `SEMREL_PLUGIN_ARTIFACT` | Optional | Single artifact, used as a fallback. | None |

Release context:

| Variable | Description |
| --- | --- |
| `SEMREL_VERSION` / `SEMREL_NEXT_VERSION` | Required release version. |
| `SEMREL_DRY_RUN` | Whether semrel is running in dry-run mode. |
