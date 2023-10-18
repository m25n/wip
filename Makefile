build/wip-%: main.go go.mod
	GOOS=$(word 1,$(subst -, ,$(basename $*))) \
	GOARCH=$(word 2,$(subst -, ,$(basename $*))) \
	go build -o $@ ./...

.PHONY: release
release: build/wip-linux-amd64 build/wip-linux-arm64 build/wip-darwin-amd64 build/wip-darwin-arm64
	gh release create v$(shell build/wip-$(shell go env GOOS)-$(shell go env GOARCH) version) --draft --generate-notes $^