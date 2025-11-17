SHELL := /bin/bash
.SILENT:
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c

REGISTRY ?= robinvandernoord
IMAGE := $(REGISTRY)/tunnellm
DOCKER_CONFIG = .docker

VERSION_FILE = app/VERSION


# Helper macro: safe git-cliff bump wrapper
define do_bump
	echo "🔍 Checking for new commits since last tag..."
	# Determine the next version
	NEW_VERSION=$$(git-cliff --bump $(1) --bumped-version 2>/dev/null || true); \
	if [ -z "$$NEW_VERSION" ]; then \
		echo "⚠️  Nothing to bump."; \
		exit 0; \
	fi; \
	# Guard against repeated 'v' prefixes
	CLEAN_VERSION=$$(echo "$$NEW_VERSION" | sed 's/^v*//'); \
	echo "$$CLEAN_VERSION" > $(VERSION_FILE); \
	# Generate changelog
	if git-cliff --bump $(1) --output CHANGELOG.md 2>&1 | grep -q "There is nothing to bump"; then \
		echo "⚠️  Nothing new to release."; \
		exit 0; \
	fi; \
	git add ./$(VERSION_FILE) CHANGELOG.md; \
	if git diff --cached --quiet; then \
		echo "⚠️  Nothing changed — skipping tag."; \
	else \
		git commit -m "chore(release): v$$CLEAN_VERSION"; \
		git tag "v$$CLEAN_VERSION"; \
		echo "✅ Released v$$CLEAN_VERSION ($(1))"; \
	fi
endef

version:
	version=$(shell cat $(VERSION_FILE))
	echo $$version

install:
	echo "Installing local release dependencies..."
	cargo install git-cliff
	if [ ! -f .gitcliff.toml ]; then \
		git-cliff --init; \
	fi
	echo "Dependencies installed."

build:
	# release build
	# -h: halt on error
	# -s: disable symbol table
	# -w: disable DWARF generation
	# -trimpath strips computers file paths from error messages
	version=$(shell cat $(VERSION_FILE))

	go build \
		-C app \
		-trimpath \
		-ldflags "-s -w -h -X main.version=$$version" \
		-o ../target/tunnellm

build-dev:
	version=$(shell cat $(VERSION_FILE))

	go build \
		-C app \
		-ldflags "-X main.version=$$version" \
		-o ../target/tunnellm

run: build-dev
	./target/tunnellm

clean:
	rm -rf target

# Auto bump (from commit history)
bump:
	$(call do_bump,auto)

# Forced bumps
bump-major:
	$(call do_bump,major)

bump-minor:
	$(call do_bump,minor)

bump-patch:
	$(call do_bump,patch)

# Preview / dry-run
bump-dry:
	NEXT_VERSION=$$(git-cliff --bump auto --bumped-version 2>/dev/null || echo "unknown")
	echo "🔍 Would bump to: $$NEXT_VERSION"

docker:
	version=$(shell cat $(VERSION_FILE))

	BUILDER_NAME="tunnellm-builder-$$RANDOM"
	trap 'echo "🧹 Removing builder $(BUILDER_NAME)..."; \
	      docker --config $(DOCKER_CONFIG) buildx rm -f "$$BUILDER_NAME" >/dev/null 2>&1 || true' EXIT

	mkdir -p $(DOCKER_CONFIG)
	echo "🔐 Logging into Docker registry..."
	docker --config $(DOCKER_CONFIG) login

	echo "🔧 Creating temporary builder $(BUILDER_NAME)..."
	docker --config $(DOCKER_CONFIG) buildx create \
		--name "$$BUILDER_NAME" --driver docker-container --use >/dev/null

	echo "🚀 Building and pushing multi-arch images (linux/amd64, linux/arm64)..."
	docker --config $(DOCKER_CONFIG) buildx bake \
		--file docker-compose.yml \
		--set *.tags+="$(IMAGE):latest" \
		--set *.tags+="$(IMAGE):$$version" \
		--push || { echo "❌ Docker build or push failed"; exit 1; }

	echo "✅ Images pushed: $(IMAGE):latest and $(IMAGE):$$version"

git-push:
	version=$(shell cat $(VERSION_FILE))
	git push
	git push --tags
	echo "📦 Published version $$version"

publish: bump docker git-push
