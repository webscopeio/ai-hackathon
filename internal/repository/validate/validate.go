package validate

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func Validate(tempDir string) error {
	// Get the absolute path to the src directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return err
	}
	
	templatePath := filepath.Join(currentDir, "nodeTemplate")
	
	// Files to copy
	filesToCopy := []string{"tsconfig.json", "pnpm-lock.yaml", "package.json", "playwright.config.ts"}
	
	// Copy each file from template to temp directory
	for _, file := range filesToCopy {
		srcFile := filepath.Join(templatePath, file)
		dstFile := filepath.Join(tempDir, file)
		
		src, err := os.Open(srcFile)
		if err != nil {
			fmt.Printf("Error opening source file %s: %v\n", file, err)
			return err
		}
		defer src.Close()
		
		dst, err := os.Create(dstFile)
		if err != nil {
			fmt.Printf("Error creating destination file %s: %v\n", file, err)
			return err
		}
		defer dst.Close()
		
		if _, err = io.Copy(dst, src); err != nil {
			fmt.Printf("Error copying file %s: %v\n", file, err)
			return err
		}
	}
	
	// Run pnpm install
	installCmd := exec.Command("pnpm", "i")
	installCmd.Dir = tempDir
	fmt.Printf("Running pnpm install in %s...\n", tempDir)
	output, err := installCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ Installation failed!\n")
		fmt.Printf("Error executing pnpm install: %v\n", err)
		fmt.Printf("Installation output: %s\n", output)
		return err
	}
	fmt.Printf("✅ Installation completed successfully!\n")
	
	// Run pnpm test
	testCmd := exec.Command("pnpm", "test")
	testCmd.Dir = tempDir
	fmt.Printf("Running tests in %s...\n", tempDir)
	output, err = testCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ Tests failed!\n")
		fmt.Printf("Error executing pnpm test: %v\n", err)
		fmt.Printf("Test output: %s\n", output)
		return err
	}
	
	fmt.Printf("✅ Tests passed successfully!\n")
	fmt.Printf("Test output: %s\n", output)
	return nil
}