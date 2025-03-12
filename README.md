# Hardhat abigen

This is a simple tool to generate Go structures from Hardhat json ABI. It uses the `abigen` tool from the `go-ethereum` repository.

## Usage

```bash
hardhat-abigen --abi <path-to-abi> --type <output-type> --pkg <output-package-name> --out <output-file>
```

To build the tool, run:

```bash
go build .
```

Move the binary to the golang bin directory:

```bash
mv hardhat-abigen $GOPATH/bin
```
