NAME=packer-post-processor-tfvars
GIT_COMMIT=$(git rev-parse HEAD)
GIT_DIRTY=$(test -n "$(git status --porcelain)" && echo "+CHANGES" || true)
REVISION=$(GIT_COMMIT)$(GIT_DIRTY)

default: test build

build:
	gox -ldflags "-X github.com/sj26/$(NAME)/VersionRevision=$(REVISION)" .

test:
	go test .

.PHONY: build test
