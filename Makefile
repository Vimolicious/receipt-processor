.PHONY: all run prod image dev dev-container dev-image clean

IMAGE_NAME := receipt-processor
DEV_IMAGE_NAME := $(IMAGE_NAME):dev
EXTERNAL_PORT := 8080
INTERNAL_PORT := 8080
CONTAINER_NAME := receipt-app
DEV_CONTAINER_NAME := $(CONTAINER_NAME)-dev
APP_ENTRYPOINT := main.go

N_DEV_CONTAINERS := $(shell echo $$((`docker ps --filter "ancestor=$(DEV_IMAGE_NAME)" | wc -l` - 1)))

all: run
run: prod

prod: image
	docker run -d --rm \
		-p $(EXTERNAL_PORT):$(INTERNAL_PORT) \
		--name $(CONTAINER_NAME) $(IMAGE_NAME)

image:
	docker build -t $(IMAGE_NAME) .

dev: dev-container
	docker exec -it $(DEV_CONTAINER_NAME) go run $(APP_ENTRYPOINT)

dev-container:
ifeq ($(N_DEV_CONTAINERS),0)
	docker run -d --rm \
		-p $(EXTERNAL_PORT):$(INTERNAL_PORT) \
		-v $(PWD):/usr/src/app \
		--name $(DEV_CONTAINER_NAME) \
		$(DEV_IMAGE_NAME)
endif

dev-image:
	docker build -t $(DEV_IMAGE_NAME) -f Dockerfile.dev .

clean:
	docker rmi -f $(IMAGE_NAME) $(DEV_IMAGE_NAME)
