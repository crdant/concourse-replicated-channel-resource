package resource

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
)

func In() {
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