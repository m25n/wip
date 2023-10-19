.PHONY: run/wip
run/wip: build/wip
	@WIPFILE=tmp/.wip $< $(WIP_ARGS)

build/wip: build/wip-$(shell go env GOOS)-$(shell go env GOARCH)
	cp -f $< $@

build/wip-%: main.go go.mod
	GOOS=$(word 1,$(subst -, ,$(basename $*))) \
	GOARCH=$(word 2,$(subst -, ,$(basename $*))) \
	go build -o $@ ./main.go

.PHONY: release
release: build/wip-linux-amd64 build/wip-linux-arm64 build/wip-darwin-amd64 build/wip-darwin-arm64
	gh release create v$(shell build/wip-$(shell go env GOOS)-$(shell go env GOARCH) version) --draft --generate-notes $^