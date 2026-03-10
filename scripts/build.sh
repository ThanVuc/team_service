#!/bin/bash
set -e 
# Test before setup cicd pipeline
DOCKER_USERNAME="sinhnguyen417"
DOCKER_REPO="team_service"
DOCKER_TAG="latest"
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

echo "🐳 Building Docker image..."
docker build -t $DOCKER_USERNAME/$DOCKER_REPO:$DOCKER_TAG $PROJECT_ROOT

echo "🔐 Logging in to Docker Hub..."
docker login -u "$DOCKER_USERNAME"

echo "📤 Pushing image to Docker Hub..."
docker push $DOCKER_USERNAME/$DOCKER_REPO:$DOCKER_TAG

echo "✅ Done! Image pushed to https://hub.docker.com/r/$DOCKER_USERNAME/$DOCKER_REPO"
