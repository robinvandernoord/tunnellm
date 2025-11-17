# Makefile
.SILENT:

IMAGE := tunellm
VERSION := $(shell cat VERSION)

# Shared bump logic
define do_bump
	NEW_VERSION=$$(semantic-release version $(1)); \
	echo $$NEW_VERSION > VERSION; \
	semantic-release changelog; \
	semantic-release publish; \
	echo "Bumped to $$NEW_VERSION"
endef

version:
	echo $(VERSION)

build:
	go build \
		-ldflags "-X main.version=$(VERSION)" \
		-o target/tunellm \
		.

run: build
	./target/tunellm

clean:
	rm -rf target

bump:
	$(call do_bump,)

bump-major:
	$(call do_bump,--major)

bump-minor:
	$(call do_bump,--minor)

bump-patch:
	$(call do_bump,--patch)

docker:
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--push \
		-t $(IMAGE):$(VERSION) \
		-t $(IMAGE):latest \
		.

publish: bump build docker
	echo "Published version $(shell cat VERSION)"
