# Copyright 2020 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

ARCHS = amd64 arm64
COMMONENVVAR=GOOS=linux
BUILDENVVAR=CGO_ENABLED=0

RELEASE_REGISTRY?=harbor.devops.qdb.com/star
RELEASE_VERSION?=$(shell git rev-parse --short HEAD)
RELEASE_IMAGE:=$(RELEASE_REGISTRY)/star_llm_backend:$(RELEASE_VERSION)


.PHONY: all
all: build

.PHONY: build
build: build-star_llm_backend

.PHONY: build.amd64
build.amd64: build-star_llm_backend.amd64

.PHONY: build.arm64v8
build.arm64v8: build-star_llm_backend.arm64v8

.PHONY: build-star_llm_backend
build-star_llm_backend:
	$(COMMONENVVAR) $(BUILDENVVAR) go build -ldflags '-w' -o bin/star_llm_backend main.go

.PHONY: build-star_llm_backend.amd64
build-star_llm_backend.amd64:
	$(COMMONENVVAR) $(BUILDENVVAR) GOARCH=amd64 go build -ldflags '-w' -o bin/star_llm_backend main.go

.PHONY: build-star_llm_backend.amd64.dlv
build-star_llm_backend.amd64.dlv:
	$(COMMONENVVAR) $(BUILDENVVAR) GOARCH=amd64 go build -gcflags "all=-N -l" -o bin/star_llm_backend main.go

.PHONY: build-star_llm_backend.arm64v8
build-star_llm_backend.arm64v8:
	GOOS=linux $(BUILDENVVAR) GOARCH=arm64 go build -ldflags '-w' -o bin/star_llm_backend main.go

.PHONY: release-image.amd64
release-image.amd64: clean
	nerdctl build  --build-arg ARCH="amd64" --build-arg RELEASE_VERSION="$(RELEASE_VERSION)" --build-arg RELEASE_REGISTRY="$(RELEASE_REGISTRY)" --build-arg RELEASE_IMAGE="$(RELEASE_IMAGE)" -t $(RELEASE_IMAGE)-amd64 .

.PHONY: release-image.arm64v8
release-image.arm64v8: clean
	nerdctl build --build-arg ARCH="arm64v8" --build-arg RELEASE_VERSION="$(RELEASE_VERSION)" --build-arg RELEASE_REGISTRY="$(RELEASE_REGISTRY)" --build-arg RELEASE_IMAGE="$(RELEASE_IMAGE)" -t $(RELEASE_IMAGE)-arm64 .

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

