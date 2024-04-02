#!/bin/bash

# Initialize default values
DOCKER_USERNAME="kaiohenricunha"
DOCKER_PASSWORD="default_password"
DB_CLEANUP="false" # Default to not cleaning up the database
IMAGE_NAME_BACKEND="go-music-k8s"
VERSION="latest"
MUSICAPI_DEPLOYMENT_NAME="musicapi"
MYSQL_DEPLOYMENT_NAME="mysql"
MUSICAPI_DEPLOYMENT_YAML="deploy/k8s/backend/musicapi"
MYSQL_DEPLOYMENT_YAML="deploy/k8s/backend/mysql"

# Parse command-line arguments for username, password, and cleanup flag
while getopts u:p:c flag
do
    case "${flag}" in
        u) DOCKER_USERNAME=${OPTARG};;
        p) DOCKER_PASSWORD=${OPTARG};;
        c) DB_CLEANUP="true";;
    esac
done

echo "Using Docker Username: $DOCKER_USERNAME"
if [ "$DB_CLEANUP" == "true" ]; then
    echo "Database cleanup will be performed."
fi

# Pass the DB_CLEANUP variable to your Go application
export DB_CLEANUP

# Check if Minikube is running
minikube status &> /dev/null
if [ $? -ne 0 ]; then
  echo "Minikube is not running, starting Minikube..."
  minikube start
  minikube addons enable ingress && minikube addons enable ingress-dns && minikube addons enable metrics-server
else
  echo "Minikube is already running."
fi

# Set Docker environment to Minikube
eval $(minikube docker-env)

# Build and Push Docker Images
docker login -u "${DOCKER_USERNAME}" -p "${DOCKER_PASSWORD}"

## Backend
echo "Building Backend Docker image..."
docker build -t "${DOCKER_USERNAME}/${IMAGE_NAME_BACKEND}:${VERSION}" .
echo "Pushing Backend Docker image to Docker Hub..."
docker push "${DOCKER_USERNAME}/${IMAGE_NAME_BACKEND}:${VERSION}"

# Deploy to Minikube

cd ..

## Make create-secrets.sh executable
chmod +x ./create-secrets.sh

## MySQL
echo "Deploying MySQL to Minikube..."
kubectl apply -f "${MYSQL_DEPLOYMENT_YAML}/db-ns.yaml"
kubectl apply -f "${MYSQL_DEPLOYMENT_YAML}/"

## Backend API
echo "Deploying Music API to Minikube..."
kubectl apply -f "${MUSICAPI_DEPLOYMENT_YAML}/api-music-ns.yaml"
./create-secrets.sh
kubectl apply -f "${MUSICAPI_DEPLOYMENT_YAML}/"
kubectl delete pods --selector=app=${MUSICAPI_DEPLOYMENT_NAME} -n music-ns

## Rollout deployments
kubectl rollout status deployment/${MYSQL_DEPLOYMENT_NAME} -n db-ns
kubectl rollout status deployment/${MUSICAPI_DEPLOYMENT_NAME} -n music-ns

echo "Deployment completed successfully."

kubectl port-forward -n music-ns svc/musicapi 8081 
