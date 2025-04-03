package validate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Validate() {
// Get the absolute path to the src directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}
	
	srcPath := filepath.Join(currentDir, "tests/src")
	
	// Create command to run pnpm test in the src directory
	cmd := exec.Command("pnpm", "test")
	cmd.Dir = srcPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing pnpm test: %v\n", err)
		fmt.Printf("Test output: %s\n", output)
		return
	}
	
	fmt.Printf("Test output: %s\n", output)
}