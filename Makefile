GOCC = go1.16

.PHONY: multi install test multi-cross
.PHONY: multi-darwin multi-darwin-386 multi-darwin-amd64
.PHONY: multi-linux multi-linux-386 multi-linux-amd64
.PHONY: multi-windows multi-windows-386 multi-windows-amd64

all: multi

multi:
	$(GOCC) build -o ./bin/multi ./cmd/multi

install:
	$(GOCC) install ./cmd/multi

install-systemd: install
	mkdir -p $(HOME)/.config/systemd/user/
	cp misc/multiverse-user.service $(HOME)/.config/systemd/user

test:
	$(GOCC) test ./... -cover

multi-cross: multi-darwin multi-linux multi-windows
	@echo "Full cross compilation done."

multi-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(GOCC) build -o ./bin/multi-darwin-amd64 ./cmd/multi
	@echo "Darwin amd64 cross compilation done."

multi-darwin: multi-darwin-amd64
	@echo "Darwin cross compilation done."

multi-linux-386:
	GOOS=linux GOARCH=386 $(GOCC) build -o ./bin/multi-linux-386 ./cmd/multi
	@echo "Linux 386 cross compilation done."

multi-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GOCC) build -o ./bin/multi-linux-amd64 ./cmd/multi
	@echo "Linux amd64 cross compilation done."

multi-linux: multi-linux-386 multi-linux-amd64
	@echo "Linux cross compilation done."

multi-windows-386:
	GOOS=windows GOARCH=386 $(GOCC) build -o ./bin/multi-windows-386 ./cmd/multi
	@echo "Windows 386 cross compilation done."

multi-windows-amd64:
	GOOS=windows GOARCH=amd64 $(GOCC) build -o ./bin/multi-windows-amd64 ./cmd/multi
	@echo "Windows amd64 cross compilation done."

multi-windows: multi-windows-386 multi-windows-amd64
	@echo "Windows cross compilation done."
