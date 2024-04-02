#!/bin/bash

# Initialize default values
DOCKER_USERNAME="kaiohenricunha"
DOCKER_PASSWORD="default_password"
DB_CLEANUP="false" # Default to not cleaning up the database
IMAGE_NAME_BACKEND="go-music-k8s"
VERSION="latest"
MUSICAPI_DEPLOYMENT_NAME="musicapi"
MYSQL_DEPLOYMENT_NAME="mysql"
MUSICAPI_DEPLOYMENT_YAML="../deploy/k8s/backend/musicapi"
MYSQL_DEPLOYMENT_YAML="../deploy/k8s/backend/mysql"

# Parse command-line arguments for username, password, and cleanup flag
while getopts u:p: flag
do
    case "${flag}" in
        u) DOCKER_USERNAME=${OPTARG};;
        p) DOCKER_PASSWORD=${OPTARG};;
    esac
done

echo "Using Docker Username: $DOCKER_USERNAME"
if [ "$DB_CLEANUP" == "true" ]; then
    echo "Database cleanup will be performed."
fi

# Pass the DB_CLEANUP variable to your Go application
export DB_CLEANUP

# Make create-secrets.sh executable
chmod +x ../create-secrets.sh

# Check if Minikube is running
minikube status &> /dev/null
if [ $? -ne 0 ]; then
  echo "Minikube is not running, starting Minikube..."
  minikube start
else
  echo "Minikube is already running."
fi

# Set Docker environment to Minikube
eval $(minikube docker-env)

# Build and Push Docker Images

## Backend
echo "Building Backend Docker image..."
docker build -t "${DOCKER_USERNAME}/${IMAGE_NAME_BACKEND}:${VERSION}" .
echo "Pushing Backend Docker image to Docker Hub..."
docker push "${DOCKER_USERNAME}/${IMAGE_NAME_BACKEND}:${VERSION}"

# Deploy to Minikube

## MySQL
echo "Deploying MySQL to Minikube..."
kubectl apply -f "${MYSQL_DEPLOYMENT_YAML}/db-ns.yaml"
kubectl apply -f "${MYSQL_DEPLOYMENT_YAML}/"
echo "Waiting for MySQL deployment to complete..."
kubectl rollout status deployment/${MYSQL_DEPLOYMENT_NAME} -n db-ns

## Backend API
echo "Deploying Music API to Minikube..."
kubectl apply -f "${MUSICAPI_DEPLOYMENT_YAML}/api-music-ns.yaml"
../create-secrets.sh
kubectl apply -f "${MUSICAPI_DEPLOYMENT_YAML}/"
echo "Waiting for Music API deployment to complete..."
kubectl rollout status deployment/${MUSICAPI_DEPLOYMENT_NAME} -n music-ns

# Run tests
echo "Running tests..."
go test -v ./...

# Linting
echo "Running linter..."
golangci-lint run

echo "Deployment completed successfully."

kubectl port-forward -n music-ns svc/musicapi 8081
