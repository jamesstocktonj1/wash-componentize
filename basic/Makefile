.PHONY: generate
generate:
	componentize-go --world basic bindings --pkg-name componentize-basic/gen --output gen --format

.PHONY: build
build: generate
	componentize-go --world basic componentize --output build/basic.wasm

.PHONY: clean
clean:
	rm -rf gen/wasi*
	rm -rf gen/wit*

	rm -rf build/