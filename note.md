
### steps

1. run the cmd to build image, here *pi* should be my image name but was a mistake.

```bash

# debug with this 
docker build -t test-website .
docker run -it --entrypoint /bin/sh test-website
docker run -p 5000:5000 test-website

docker build --target debug -t pi/raspberry-go-app:debug .

docker build --target production -t pi/raspberry-go-app:production .

docker build -t pi/raspberry-go-app:1.0 .

docker images
```
2. test the container locally
```bash
docker run -p 5000:5000 pi/raspberry-go-app:1.0
```
3. ensure the container always on 
```bash
docker run -d -p 5000:5000 --name my-go-app --restart unless-stopped pi/raspberry-go-app:1.0
```
4. verify the container is running 
```bash
docker ps
```
5. view log
`docker logs my-go-app`

6. stop and start container
```bash
docker stop my-go-app

docker start my-go-app
```
### debug

start an interactive shell in the container

```bash
docker run -it --entrypoint /bin/sh pi/raspberry-go-app:1.0
```