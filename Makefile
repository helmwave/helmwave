all: build
	@echo "Finish"

build: test
	go build github.com/zhilyaev/helmwave/cmd/helmwave
	cp -f helmwave /usr/local/bin/helmwave

test:
	go test github.com/zhilyaev/helmwave/cmd/helmwave

