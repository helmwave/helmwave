all: build
	@echo "Finish"

build: test
	go build github.com/helmwave/helmwave/cmd/helmwave
	cp -f helmwave /usr/local/bin/helmwave

test:
	go test github.com/helmwave/helmwave/pkg/yml

