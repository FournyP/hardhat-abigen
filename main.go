package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/FournyP/hardhat-abigen/tui"
	tea "github.com/charmbracelet/bubbletea"
)

type HardhatOutput struct {
	ABI      json.RawMessage `json:"abi"`
	ByteCode string          `json:"bytecode"`
}

func main() {
	requiredParams := []string{"abi", "out", "type", "pkg"}
	params := make(map[string]string)

	for _, key := range requiredParams {
		if val, exists := os.LookupEnv(strings.ToUpper(key)); exists {
			params[key] = val
		}
	}

	missingKeys := []string{}
	for _, key := range requiredParams {
		if _, exists := params[key]; !exists {
			missingKeys = append(missingKeys, key)
		}
	}

	m := tui.NewModel(params, missingKeys)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	fmt.Println("Collected parameters:")
	for k, v := range m.Params {
		fmt.Printf("%s: %s\n", k, v)
	}

	// Read and extract ABI
	file, err := os.Open(m.Params["abi"])
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

	// Run abigen
	cmd := exec.Command("abigen", "--abi="+abiTempPath, "--pkg="+m.Params["pkg"], "--type="+m.Params["type"], "--out="+m.Params["out"], "--bin="+binTempPath)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to run abigen: %v", err)
	}
}
