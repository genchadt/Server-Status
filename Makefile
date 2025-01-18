.PHONY: build install uninstall clean

build:
	@./scripts/build.sh

install: build
	@sudo ./scripts/install.sh

uninstall:
	@sudo ./scripts/uninstall.sh

clean:
	@rm -f serverstatus
	@rm -rf logs
