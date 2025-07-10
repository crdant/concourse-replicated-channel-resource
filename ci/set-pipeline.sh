#!/bin/bash

set -e

if [ -z "$CONCOURSE_URL" ]; then
    echo "CONCOURSE_URL environment variable is required"
    exit 1
fi

if [ -z "$CONCOURSE_USERNAME" ]; then
    echo "CONCOURSE_USERNAME environment variable is required"
    exit 1
fi

if [ -z "$CONCOURSE_PASSWORD" ]; then
    echo "CONCOURSE_PASSWORD environment variable is required"
    exit 1
fi

if [ -z "$CONCOURSE_TEAM" ]; then
    CONCOURSE_TEAM="main"
fi

# Login to Concourse
fly -t ci login \
    -c "$CONCOURSE_URL" \
    -u "$CONCOURSE_USERNAME" \
    -p "$CONCOURSE_PASSWORD" \
    -n "$CONCOURSE_TEAM"

# Set the pipeline
fly -t ci set-pipeline \
    -p concourse-replicated-channel-resource \
    -c ci/pipeline.yml \
    --var github-username="$GITHUB_USERNAME" \
    --var github-token="$GITHUB_TOKEN"

# Unpause the pipeline
fly -t ci unpause-pipeline -p concourse-replicated-channel-resource

echo "Pipeline set successfully!"
echo "View it at: $CONCOURSE_URL/teams/$CONCOURSE_TEAM/pipelines/concourse-replicated-channel-resource"