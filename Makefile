NAME=disposable
IMAGE=0x19/disposable
GOPATH?=/Users/0x19/src/go
BRANCH=master

all: build-docker

proto-go:
	protoc -I/usr/local/include -I$(GOPATH)/src  -I. --python_out=. --grpc_out=. --plugin=protoc-gen-grpc=`which grpc_python_plugin` --go_out=plugins=grpc:. protos/*.proto
	ls protos/*.pb.go | xargs -n1 -IX bash -c "sed -e '/bool/ s/,omitempty//' X > X.tmp && mv X{.tmp,}"

submodules:
	git submodule init
	git submodule update
	git submodule foreach git checkout $(BRANCH)
	git submodule foreach git pull origin $(BRANCH)

update: submodules

build-docker:
	docker build -t $(NAME) .

run: build-docker
	docker run  -it --name=$(NAME) $(NAME)

sh: build-docker
	docker run -it $(NAME) /bin/sh

push: build-docker
	docker tag $(NAME) $(IMAGE)
	docker push $(IMAGE)
