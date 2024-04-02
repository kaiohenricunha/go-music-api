#!/bin/bash

# Initialize default values
DOCKER_USERNAME="kaiohenricunha"
DOCKER_PASSWORD="default_password"
DB_CLEANUP="false" # Default to not cleaning up the database
IMAGE_NAME_BACKEND="go-music-k8s"
IMAGE_NAME_FRONTEND="react-music-app"
VERSION="latest"
MUSICAPI_DEPLOYMENT_NAME="musicapi"
MYSQL_DEPLOYMENT_NAME="mysql"
FRONTEND_DEPLOYMENT_NAME="react-frontend"
MUSICAPI_DEPLOYMENT_YAML="deploy/k8s/backend/musicapi"
MYSQL_DEPLOYMENT_YAML="deploy/k8s/backend/mysql"
FRONTEND_DEPLOYMENT_YAML="deploy/k8s/frontend"

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

# Make create-secrets.sh executable
chmod +x create-secrets.sh

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
cd backend || exit
docker build -t "${DOCKER_USERNAME}/${IMAGE_NAME_BACKEND}:${VERSION}" .
echo "Pushing Backend Docker image to Docker Hub..."
docker push "${DOCKER_USERNAME}/${IMAGE_NAME_BACKEND}:${VERSION}"
cd .. || exit

## Frontend
echo "Building Frontend Docker image..."
cd frontend || exit
docker build -t "${DOCKER_USERNAME}/${IMAGE_NAME_FRONTEND}:${VERSION}" .
echo "Pushing Frontend Docker image to Docker Hub..."
docker push "${DOCKER_USERNAME}/${IMAGE_NAME_FRONTEND}:${VERSION}"
cd .. || exit

# Deploy to Minikube

## MySQL
echo "Deploying MySQL to Minikube..."
kubectl apply -f "${MYSQL_DEPLOYMENT_YAML}/db-ns.yaml"
kubectl apply -f "${MYSQL_DEPLOYMENT_YAML}/"
echo "Waiting for MySQL deployment to complete..."
kubectl rollout status deployment/${MYSQL_DEPLOYMENT_NAME} -n db-ns

## Backend API
echo "Deploying Music API to Minikube..."
kubectl delete deployment ${MUSICAPI_DEPLOYMENT_NAME} -n music-ns
kubectl apply -f "${MUSICAPI_DEPLOYMENT_YAML}/api-music-ns.yaml"
./create-secrets.sh
kubectl apply -f "${MUSICAPI_DEPLOYMENT_YAML}/"
echo "Waiting for Music API deployment to complete..."
kubectl rollout status deployment/${MUSICAPI_DEPLOYMENT_NAME} -n music-ns

## Frontend
echo "Deploying React Frontend to Minikube..."
kubectl delete deployment ${FRONTEND_DEPLOYMENT_NAME} -n music-ns
kubectl apply -f "${FRONTEND_DEPLOYMENT_YAML}"
echo "Waiting for React Frontend deployment to complete..."
kubectl rollout status deployment/${FRONTEND_DEPLOYMENT_NAME} -n music-ns

echo "Deployment completed successfully."

MINIKUBE_IP=$(minikube ip)
echo "To update /etc/hosts, please run the following command with root permissions:"
echo "echo '$MINIKUBE_IP musicapi.local' | sudo tee -a /etc/hosts"

kubectl port-forward service/react-frontend-service 3000:3000 -n music-ns
