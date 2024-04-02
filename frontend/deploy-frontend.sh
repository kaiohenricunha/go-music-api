#!/bin/bash

# Initialize default values
DOCKER_USERNAME="kaiohenricunha"
DOCKER_PASSWORD="default_password"
IMAGE_NAME_FRONTEND="react-music-app"
VERSION="latest"
FRONTEND_DEPLOYMENT_NAME="react-frontend"
FRONTEND_DEPLOYMENT_YAML="deploy/k8s/frontend"

# Parse command-line arguments for username, password flags
while getopts u:p: flag
do
    case "${flag}" in
        u) DOCKER_USERNAME=${OPTARG};;
        p) DOCKER_PASSWORD=${OPTARG};;
    esac
done

echo "Using Docker Username: $DOCKER_USERNAME"

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

## Frontend
echo "Building Frontend Docker image..."
docker build -t "${DOCKER_USERNAME}/${IMAGE_NAME_FRONTEND}:${VERSION}" .
echo "Pushing Frontend Docker image to Docker Hub..."
docker push "${DOCKER_USERNAME}/${IMAGE_NAME_FRONTEND}:${VERSION}"

# Deploy to Minikube

## Frontend
echo "Deploying React Frontend to Minikube..."

cd ..

kubectl apply -f ${FRONTEND_DEPLOYMENT_YAML} -n music-ns
echo "Waiting for React Frontend deployment to complete..."
kubectl delete pods --selector=app=react-frontend -n music-ns
kubectl rollout status deployment/${FRONTEND_DEPLOYMENT_NAME} -n music-ns

echo "Deployment completed successfully."

kubectl port-forward service/react-frontend-service 3000:3000 -n music-ns
