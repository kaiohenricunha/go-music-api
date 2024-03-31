#!/bin/bash

# Initialize default values
DOCKER_USERNAME="kaiohenricunha"
DOCKER_PASSWORD="default_password"

# Parse command-line arguments for username and password
while getopts u:p: flag
do
    case "${flag}" in
        u) DOCKER_USERNAME=${OPTARG};;
        p) DOCKER_PASSWORD=${OPTARG};;
    esac
done

echo "Using Docker Username: $DOCKER_USERNAME"

# Configuration variables
IMAGE_NAME="go-music-k8s"
VERSION="latest"
MUSICAPI_DEPLOYMENT_NAME="musicapi"
MYSQL_DEPLOYMENT_NAME="mysql"
MUSICAPI_DEPLOYMENT_YAML="deploy/k8s/backend/musicapi"
MYSQL_DEPLOYMENT_YAML="deploy/k8s/backend/mysql"

# Build Docker image
echo "Building Docker image..."
cd backend || exit
docker build -t "${DOCKER_USERNAME}/${IMAGE_NAME}:${VERSION}" .
cd .. || exit

# Push Docker image
echo "Pushing Docker image to Docker Hub..."
docker login -u "${DOCKER_USERNAME}" -p "${DOCKER_PASSWORD}"
docker push "${DOCKER_USERNAME}/${IMAGE_NAME}:${VERSION}"

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

# Check if the MySQL deployment already exists
kubectl get deployment ${MYSQL_DEPLOYMENT_NAME} -n db-ns &> /dev/null
if [ $? -eq 0 ]; then
  echo "MySQL Deployment already exists. Deleting the existing deployment..."
  kubectl delete deployment ${MYSQL_DEPLOYMENT_NAME} -n db-ns
fi

# Deploy MySQL
echo "Deploying MySQL to Minikube..."
kubectl apply -f "${MYSQL_DEPLOYMENT_YAML}/db-ns.yaml"
kubectl apply -f "${MYSQL_DEPLOYMENT_YAML}/"

# Wait for MySQL deployment to complete
echo "Waiting for MySQL deployment to complete..."
kubectl rollout status deployment/${MYSQL_DEPLOYMENT_NAME} -n db-ns

# Check if the Music API deployment already exists
kubectl get deployment ${MUSICAPI_DEPLOYMENT_NAME} -n music-ns &> /dev/null
if [ $? -eq 0 ]; then
  echo "Music API Deployment already exists. Deleting the existing deployment..."
  kubectl delete deployment ${MUSICAPI_DEPLOYMENT_NAME} -n music-ns
fi

# call create-spotify-secret.sh
./create-spotify-secret.sh

# Deploy Music API to Minikube
echo "Deploying Music API to Minikube..."
kubectl apply -f "${MUSICAPI_DEPLOYMENT_YAML}/api-music-ns.yaml"
kubectl apply -f "${MUSICAPI_DEPLOYMENT_YAML}/"

# Wait for Music API deployment to complete
echo "Waiting for Music API deployment to complete..."
kubectl rollout status deployment/${MUSICAPI_DEPLOYMENT_NAME} -n music-ns

echo "Deployment completed successfully."
