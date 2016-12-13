cd paddle/scripts/docker/
docker build -f Dockerfile.cpu -t paddle_ssh .
docker rmi paddle_ssh
