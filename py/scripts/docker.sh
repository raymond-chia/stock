docker rm stock
docker run --name stock -it -v $(pwd):/app stock
