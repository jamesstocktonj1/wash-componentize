.PHONY: generate
generate:
	componentize-go --world redis-echo bindings --pkg-name redis-echo/gen --output gen --format

.PHONY: build
build: generate
	componentize-go --world redis-echo componentize --output build/redis-echo.wasm

.PHONY: clean
clean:
	rm -rf gen/wasi*
	rm -rf gen/wit*

	rm -rf build/