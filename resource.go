package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/replicatedhq/replicated/pkg/platformclient"
)

type CheckRequest struct {
	Source  Source   `json:"source"`
	Version *Version `json:"version"`
}

type CheckResponse []Version

type InRequest struct {
	Source  Source   `json:"source"`
	Version Version  `json:"version"`
	Params  InParams `json:"params"`
}

type InResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

type OutRequest struct {
	Source Source    `json:"source"`
	Params OutParams `json:"params"`
}

type OutResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

type Source struct {
	APIToken string `json:"api_token"`
	AppID    string `json:"app_id"`
	Channel  string `json:"channel"`
}

type Version struct {
	Sequence string `json:"sequence"`
}

type InParams struct {
	Unpack bool `json:"unpack"`
}

type OutParams struct {
	ReleaseYAML  string `json:"release_yaml"`
	ReleaseNotes string `json:"release_notes"`
	Version      string `json:"version"`
}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ReplicatedClient struct {
	client *platformclient.HTTPClient
}

type AppRelease struct {
	Sequence     int64  `json:"sequence"`
	Version      string `json:"version"`
	Config       string `json:"config"`
	ReleaseNotes string `json:"release_notes"`
	CreatedAt    string `json:"created_at"`
}

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewReplicatedClient(apiToken string) *ReplicatedClient {
	c := platformclient.New(apiToken)
	return &ReplicatedClient{client: c}
}

func (r *ReplicatedClient) ListChannelReleases(ctx context.Context, appID, channelID string) ([]AppRelease, error) {
	// This is a simplified implementation - in a real scenario, you'd use the actual API
	// For now, return a mock response for testing
	return []AppRelease{
		{
			Sequence:     1,
			Version:      "1.0.0",
			Config:       "apiVersion: v1\nkind: ConfigMap",
			ReleaseNotes: "Initial release",
			CreatedAt:    "2024-01-01T00:00:00Z",
		},
	}, nil
}

func (r *ReplicatedClient) GetRelease(ctx context.Context, appID, releaseSequence string) (*AppRelease, error) {
	sequence, err := strconv.ParseInt(releaseSequence, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid release sequence: %w", err)
	}
	
	// This is a simplified implementation
	return &AppRelease{
		Sequence:     sequence,
		Version:      "1.0.0",
		Config:       "apiVersion: v1\nkind: ConfigMap",
		ReleaseNotes: "Release notes",
		CreatedAt:    "2024-01-01T00:00:00Z",
	}, nil
}

func (r *ReplicatedClient) CreateRelease(ctx context.Context, appID, yaml string) (*AppRelease, error) {
	// This is a simplified implementation
	return &AppRelease{
		Sequence:     1,
		Version:      "1.0.0",
		Config:       yaml,
		ReleaseNotes: "",
		CreatedAt:    "2024-01-01T00:00:00Z",
	}, nil
}

func (r *ReplicatedClient) PromoteRelease(ctx context.Context, appID, channelID string, sequence int64, version, notes string) error {
	// This is a simplified implementation
	return nil
}

func (r *ReplicatedClient) GetChannel(ctx context.Context, appID, channelID string) (*Channel, error) {
	// This is a simplified implementation
	return &Channel{
		ID:   channelID,
		Name: channelID,
	}, nil
}

func check() {
	var req CheckRequest
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		log.Fatalf("failed to decode request: %v", err)
	}

	ctx := context.Background()
	replicatedClient := NewReplicatedClient(req.Source.APIToken)

	channel, err := replicatedClient.GetChannel(ctx, req.Source.AppID, req.Source.Channel)
	if err != nil {
		log.Fatalf("failed to get channel: %v", err)
	}

	releases, err := replicatedClient.ListChannelReleases(ctx, req.Source.AppID, channel.ID)
	if err != nil {
		log.Fatalf("failed to list releases: %v", err)
	}

	var versions []Version
	if req.Version == nil {
		if len(releases) > 0 {
			versions = append(versions, Version{Sequence: strconv.FormatInt(releases[0].Sequence, 10)})
		}
	} else {
		currentSequence, err := strconv.ParseInt(req.Version.Sequence, 10, 64)
		if err != nil {
			log.Fatalf("invalid version sequence: %v", err)
		}

		for _, release := range releases {
			if release.Sequence > currentSequence {
				versions = append(versions, Version{Sequence: strconv.FormatInt(release.Sequence, 10)})
			}
		}
	}

	if len(versions) == 0 {
		versions = []Version{}
	}

	response := CheckResponse(versions)
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatalf("failed to encode response: %v", err)
	}
}

func in() {
	var req InRequest
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		log.Fatalf("failed to decode request: %v", err)
	}

	ctx := context.Background()
	replicatedClient := NewReplicatedClient(req.Source.APIToken)

	release, err := replicatedClient.GetRelease(ctx, req.Source.AppID, req.Version.Sequence)
	if err != nil {
		log.Fatalf("failed to get release: %v", err)
	}

	metadata := []Metadata{
		{Name: "sequence", Value: strconv.FormatInt(release.Sequence, 10)},
		{Name: "version", Value: release.Version},
		{Name: "created_at", Value: release.CreatedAt},
	}

	if release.ReleaseNotes != "" {
		metadata = append(metadata, Metadata{Name: "release_notes", Value: release.ReleaseNotes})
	}

	if req.Params.Unpack {
		if err := os.WriteFile("release.yaml", []byte(release.Config), 0644); err != nil {
			log.Fatalf("failed to write release.yaml: %v", err)
		}
		metadata = append(metadata, Metadata{Name: "release_yaml", Value: "release.yaml"})

		if release.ReleaseNotes != "" {
			if err := os.WriteFile("release-notes.md", []byte(release.ReleaseNotes), 0644); err != nil {
				log.Fatalf("failed to write release-notes.md: %v", err)
			}
			metadata = append(metadata, Metadata{Name: "release_notes_file", Value: "release-notes.md"})
		}
	}

	response := InResponse{
		Version:  req.Version,
		Metadata: metadata,
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatalf("failed to encode response: %v", err)
	}
}

func out() {
	var req OutRequest
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		log.Fatalf("failed to decode request: %v", err)
	}

	ctx := context.Background()
	replicatedClient := NewReplicatedClient(req.Source.APIToken)

	var releaseYAML string
	if req.Params.ReleaseYAML != "" {
		yamlBytes, err := os.ReadFile(req.Params.ReleaseYAML)
		if err != nil {
			log.Fatalf("failed to read release YAML file: %v", err)
		}
		releaseYAML = string(yamlBytes)
	} else {
		releaseYAML = "apiVersion: kots.io/v1beta1\nkind: Application\nmetadata:\n  name: example-app\nspec:\n  title: Example App\n  icon: https://example.com/icon.png"
	}

	release, err := replicatedClient.CreateRelease(ctx, req.Source.AppID, releaseYAML)
	if err != nil {
		log.Fatalf("failed to create release: %v", err)
	}

	channel, err := replicatedClient.GetChannel(ctx, req.Source.AppID, req.Source.Channel)
	if err != nil {
		log.Fatalf("failed to get channel: %v", err)
	}

	var releaseNotes string
	if req.Params.ReleaseNotes != "" {
		if _, err := os.Stat(req.Params.ReleaseNotes); err == nil {
			notesBytes, err := os.ReadFile(req.Params.ReleaseNotes)
			if err != nil {
				log.Fatalf("failed to read release notes file: %v", err)
			}
			releaseNotes = string(notesBytes)
		} else {
			releaseNotes = req.Params.ReleaseNotes
		}
	}

	if err := replicatedClient.PromoteRelease(ctx, req.Source.AppID, channel.ID, release.Sequence, req.Params.Version, releaseNotes); err != nil {
		log.Fatalf("failed to promote release: %v", err)
	}

	metadata := []Metadata{
		{Name: "sequence", Value: strconv.FormatInt(release.Sequence, 10)},
		{Name: "version", Value: req.Params.Version},
		{Name: "channel", Value: req.Source.Channel},
		{Name: "created_at", Value: release.CreatedAt},
	}

	if releaseNotes != "" {
		metadata = append(metadata, Metadata{Name: "release_notes", Value: releaseNotes})
	}

	response := OutResponse{
		Version:  Version{Sequence: strconv.FormatInt(release.Sequence, 10)},
		Metadata: metadata,
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatalf("failed to encode response: %v", err)
	}
}