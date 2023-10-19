GOFILES = $(wildcard **/*.go)

.PHONY: run/wip
run/wip: build/wip
	@WIPFILE=file:tmp/.wip $< $(WIP_ARGS)

build/wip: build/wip-$(shell go env GOOS)-$(shell go env GOARCH)
	cp -f $< $@

build/wip-%: $(GOFILES) go.mod
	GOOS=$(word 1,$(subst -, ,$(basename $*))) \
	GOARCH=$(word 2,$(subst -, ,$(basename $*))) \
	go build -o $@ ./cmd/...

.PHONY: release
release: build/wip-linux-amd64 build/wip-linux-arm64 build/wip-darwin-amd64 build/wip-darwin-arm64
	gh release create v$(shell build/wip-$(shell go env GOOS)-$(shell go env GOARCH) version) --draft --generate-notes $^

.PHONY: clean
clean:
	rm -rf build