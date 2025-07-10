# Concourse Replicated Channel Resource

A Concourse CI resource for managing releases on the Replicated Platform with channel support.

## Source Configuration

- `api_token`: *Required.* Your Replicated API token
- `app_id`: *Required.* The Replicated application ID
- `channel`: *Required.* The channel name (e.g., "stable", "beta", "unstable")

## Behavior

### `check`: Check for new releases

Returns the latest release sequence number from the specified channel.

### `in`: Fetch a release

Downloads release information and optionally unpacks the release configuration.

#### Parameters

- `unpack`: *Optional.* If true, writes the release YAML to `release.yaml` and release notes to `release-notes.md`

### `out`: Create and promote a release

Creates a new release and promotes it to the specified channel.

#### Parameters

- `release_yaml`: *Required.* Path to the release YAML file
- `release_notes`: *Optional.* Release notes (string or path to file)
- `version`: *Required.* Version string for the release

## Example Configuration

```yaml
resource_types:
  - name: replicated-channel
    type: registry-image
    source:
      repository: ghcr.io/replicatedhq/concourse-replicated-channel-resource
      tag: latest

resources:
  - name: replicated-release
    type: replicated-channel
    source:
      api_token: ((replicated-api-token))
      app_id: ((replicated-app-id))
      channel: stable

jobs:
  - name: deploy-to-replicated
    plan:
      - get: source-code
        trigger: true
      - put: replicated-release
        params:
          release_yaml: source-code/release.yaml
          release_notes: source-code/CHANGELOG.md
          version: "1.0.0"
```

## Development

### Building

```bash
go build -o concourse-replicated-channel-resource .
```

### Testing

```bash
go test -v ./...
```

### Building Container Image

```bash
docker build -t concourse-replicated-channel-resource .
```

## License

Apache 2.0