package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/FournyP/hardhat-abigen/tui"
)

type HardhatOutput struct {
	ABI      json.RawMessage `json:"abi"`
	ByteCode string          `json:"bytecode"`
}

func main() {
	abiFile := flag.String("abi", "", "Path to the ABI file")
	outFile := flag.String("out", "", "Path to the output file")
	typeName := flag.String("type", "", "Type name")
	pkgName := flag.String("pkg", "", "Package name")

	// Parse flags
	flag.Parse()

	// Prompt for missing values
	if *abiFile == "" {
		*abiFile = tui.PromptInput("Enter the ABI file path:")
	}

	if *outFile == "" {
		*outFile = tui.PromptInput("Enter the output file:")
	}

	if *typeName == "" {
		*typeName = tui.PromptInput("Enter the type name:")
	}

	if *pkgName == "" {
		*pkgName = tui.PromptInput("Enter the package name:")
	}

	log.Println("ABI file:", *abiFile)
	log.Println("Output file:", *outFile)
	log.Println("Type name:", *typeName)
	log.Println("Package name:", *pkgName)

	// Read and extract ABI
	file, err := os.Open(*abiFile)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}

	var output HardhatOutput
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	if err := json.Unmarshal(fileBytes, &output); err != nil {
		log.Fatalf("Failed to unmarshal ABI: %v", err)
	}

	// Write ABI to temporary file
	abiTempFile, err := os.CreateTemp("", "abi_temp_abi_*.json")
	if err != nil {
		log.Fatalf("Failed to create temp file: %v", err)
	}
	if _, err := abiTempFile.Write(output.ABI); err != nil {
		log.Fatalf("Failed to write to temp file: %v", err)
	}
	abiTempPath := abiTempFile.Name()
	defer os.Remove(abiTempPath)
	abiTempFile.Close()

	binTempFile, err := os.CreateTemp("", "bin_temp_bin_*.json")
	if err != nil {
		log.Fatalf("Failed to create temp file: %v", err)
	}
	if _, err := binTempFile.Write([]byte(output.ByteCode)); err != nil {
		log.Fatalf("Failed to write to temp file: %v", err)
	}
	binTempPath := binTempFile.Name()
	defer os.Remove(binTempPath)
	binTempFile.Close()

	log.Printf("ABI temp file: %s", abiTempPath)
	log.Printf("Bytecode temp file: %s", binTempPath)

	log.Printf("Generating Go bindings...")

	// Run abigen
	cmd := exec.Command("abigen", "--abi="+abiTempPath, "--pkg="+*pkgName, "--type="+*typeName, "--out="+*outFile, "--bin="+binTempPath)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to run abigen: %v", err)
	}
}
