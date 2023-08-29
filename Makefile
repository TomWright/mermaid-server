DOCKER_IMAGE=tomwright/mermaid-server:latest
CONTAINER_NAME=mermaid-server

docker-image:
	docker build --platform linux/x86_64 -t ${DOCKER_IMAGE} .

docker-run:
	docker run -d --platform linux/x86_64 --name ${CONTAINER_NAME} -p 80:80 ${DOCKER_IMAGE}

docker-stop:
	docker stop ${CONTAINER_NAME} || true

docker-rm:
	make docker-stop
	docker rm ${CONTAINER_NAME} || true

docker-logs:
	docker logs -f ${CONTAINER_NAME}

docker-push:
	docker push ${DOCKER_IMAGE}
