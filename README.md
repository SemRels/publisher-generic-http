# publisher-generic-http

[![Latest Release](https://img.shields.io/github/v/release/SemRels/publisher-generic-http?label=version\&color=blue)](https://github.com/SemRels/publisher-generic-http/releases/latest)

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
	- uses: publisher-generic-http
		args:
			url: https://uploads.example.com/releases/{version}/{artifact}
			method: PUT
			token: ${UPLOAD_TOKEN}
			artifacts: dist/myapp-linux-amd64,dist/myapp-linux-arm64
```

Runtime inputs:

- `SEMREL_VERSION` / `SEMREL_NEXT_VERSION` (required)
- `SEMREL_DRY_RUN`
- `SEMREL_PLUGIN_URL` (required; supports `{version}` and `{artifact}` tokens)
- `SEMREL_PLUGIN_METHOD` (default: `PUT`)
- `SEMREL_PLUGIN_TOKEN` (optional bearer token)
- `SEMREL_PLUGIN_HEADERS_JSON` (optional JSON object for custom headers)
- `SEMREL_PLUGIN_ARTIFACTS` / `SEMREL_PLUGIN_ARTIFACTS_JSON` / `SEMREL_PLUGIN_ARTIFACT`
