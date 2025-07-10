package resource

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
)

func Check() {
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