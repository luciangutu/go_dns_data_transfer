IMAGE_NAME = dns-data-transfer
RESPONSE_IP = 192.168.100.1
DATA_DIR = ./data

all: build run

build: build-server build-client

build-server:
	docker build -t $(IMAGE_NAME) .

build-client:
	go build -o client/client ./client

run:
	docker run --rm -v $(DATA_DIR):/app/data --name $(IMAGE_NAME) $(IMAGE_NAME) $(RESPONSE_IP)
