SHELL := /bin/bash
.SILENT:
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c

IMAGE := tunellm
VERSION := $(shell cat VERSION)

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
	echo "$$CLEAN_VERSION" > VERSION; \
	# Generate changelog
	if git-cliff --bump $(1) --output CHANGELOG.md 2>&1 | grep -q "There is nothing to bump"; then \
		echo "⚠️  Nothing new to release."; \
		exit 0; \
	fi; \
	git add VERSION CHANGELOG.md; \
	if git diff --cached --quiet; then \
		echo "⚠️  Nothing changed — skipping tag."; \
	else \
		git commit -m "chore(release): v$$CLEAN_VERSION"; \
		git tag "v$$CLEAN_VERSION"; \
		echo "✅ Released v$$CLEAN_VERSION ($(1))"; \
	fi
endef

version:
	echo $(VERSION)

install:
	echo "Installing local release dependencies..."
	cargo install git-cliff
	if [ ! -f .gitcliff.toml ]; then \
		git-cliff --init; \
	fi
	echo "Dependencies installed."

build:
	go build \
		-ldflags "-X main.version=$(VERSION)" \
		-o target/tunellm \
		.

run: build
	./target/tunellm

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
	docker buildx build -t $(IMAGE):$(VERSION) .

publish: bump build docker
	echo "📦 Published version $(shell cat VERSION)"