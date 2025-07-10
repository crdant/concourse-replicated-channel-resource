package resource

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
)

func Out() {
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