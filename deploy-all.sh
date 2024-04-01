#!/bin/bash

# Initialize default values
DOCKER_USERNAME="kaiohenricunha"
DOCKER_PASSWORD="default_password"
DB_CLEANUP="false" # Default to not cleaning up the database

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

# Configuration variables
IMAGE_NAME="go-music-k8s"
VERSION="latest"
MUSICAPI_DEPLOYMENT_NAME="musicapi"
MYSQL_DEPLOYMENT_NAME="mysql"
MUSICAPI_DEPLOYMENT_YAML="deploy/k8s/backend/musicapi"
MYSQL_DEPLOYMENT_YAML="deploy/k8s/backend/mysql"

# Change into the backend directory
echo "Changing into backend directory..."
cd backend || exit

# Run Go tests
echo "Running Go tests..."
if ! go test ./...; then
  echo "Tests failed. Exiting..."
  exit 1
fi

# Run golangci-lint
echo "Running golangci-lint..."
if ! golangci-lint run ./...; then
  echo "Linting errors detected. Exiting..."
  exit 1
fi

echo "Tests and Linting passed. Proceeding with the build."

# Build Docker image
echo "Building Docker image..."
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
  kubectl delete pvc --all -n db-ns
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

# Deploy Music API to Minikube
echo "Deploying Music API to Minikube..."
kubectl apply -f "${MUSICAPI_DEPLOYMENT_YAML}/api-music-ns.yaml"

## create Spotify and JWT secrets
./create-secrets.sh

kubectl apply -f "${MUSICAPI_DEPLOYMENT_YAML}/"

# Wait for Music API deployment to complete
echo "Waiting for Music API deployment to complete..."
kubectl rollout status deployment/${MUSICAPI_DEPLOYMENT_NAME} -n music-ns

echo "Deployment completed successfully."

# Display logs
echo "Displaying Music API logs..."
kubectl logs -f -n music-ns -l app=musicapi --max-log-requests=1
