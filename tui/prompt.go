package tui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PromptInput asks the user for input
func PromptInput(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question + " ")
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
