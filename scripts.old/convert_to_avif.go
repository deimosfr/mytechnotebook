// Script: convert_to_avif.go
//
// Purpose: This script converts all JPG, PNG, and GIF images in the static/images folder
// to the AVIF format, updates references in markdown files in the content folder,
// and deletes the original images after successful conversion.
//
// CLAUDE AI Prompt used to generate this code:
// "Create a go script that will convert all images (jpg, png, gif)
// to the new avif format and update image extension references in markdown files."

package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

const (
	imagesDir  = "static/images"
	contentDir = "content"
)

func main() {
	fmt.Println("Starting image conversion to AVIF and markdown updates...")

	// Get the current working directory (project root)
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// Convert images to AVIF
	imageFiles, err := convertImagesToAVIF(filepath.Join(cwd, imagesDir))
	if err != nil {
		fmt.Printf("Error converting images: %v\n", err)
		return
	}

	fmt.Printf("Successfully converted %d images to AVIF format\n", len(imageFiles))

	// Update markdown files
	updatedFiles, err := updateMarkdownFiles(filepath.Join(cwd, contentDir), imageFiles)
	if err != nil {
		fmt.Printf("Error updating markdown files: %v\n", err)
		return
	}

	fmt.Printf("Successfully updated %d markdown files\n", updatedFiles)
}

// convertImagesToAVIF converts all JPG, PNG, and GIF files to AVIF format
// and deletes the original files after successful conversion
func convertImagesToAVIF(imagesPath string) (map[string]string, error) {
	convertedFiles := make(map[string]string)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var conversionErrors []string

	err := filepath.WalkDir(imagesPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// Check if file is JPG, PNG, or GIF
		lowerName := strings.ToLower(d.Name())
		if !strings.HasSuffix(lowerName, ".jpg") &&
			!strings.HasSuffix(lowerName, ".jpeg") &&
			!strings.HasSuffix(lowerName, ".png") &&
			!strings.HasSuffix(lowerName, ".gif") {
			return nil
		}

		wg.Add(1)
		go func(filePath string, fileName string) {
			defer wg.Done()

			fmt.Printf("Converting %s to AVIF...\n", filePath)

			// Output file name
			baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
			outputPath := filepath.Join(filepath.Dir(filePath), baseName+".avif")

			// Use the convert command from ImageMagick to convert to AVIF
			// You may need to install ImageMagick: brew install imagemagick
			cmd := exec.Command("convert", filePath, outputPath)
			if output, err := cmd.CombinedOutput(); err != nil {
				errorMsg := fmt.Sprintf("Error converting %s: %v - %s", filePath, err, output)
				mu.Lock()
				conversionErrors = append(conversionErrors, errorMsg)
				mu.Unlock()
				return
			}

			// Store original -> new mapping
			relativePath, _ := filepath.Rel(imagesPath, filePath)
			relativeOutputPath, _ := filepath.Rel(imagesPath, outputPath)

			mu.Lock()
			convertedFiles[relativePath] = relativeOutputPath
			mu.Unlock()

			// Delete the original file
			if err := os.Remove(filePath); err != nil {
				fmt.Printf("Error deleting original file %s: %v\n", filePath, err)
			}
		}(path, d.Name())

		return nil
	})

	wg.Wait()

	if len(conversionErrors) > 0 {
		fmt.Println("Some conversions failed:")
		for _, err := range conversionErrors {
			fmt.Println(" -", err)
		}
		fmt.Println("Please make sure ImageMagick is installed: brew install imagemagick")
	}

	return convertedFiles, err
}

// updateMarkdownFiles updates image references in markdown files
func updateMarkdownFiles(contentPath string, imageMap map[string]string) (int, error) {
	updatedFileCount := 0

	err := filepath.WalkDir(contentPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// Only process markdown files
		if !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}

		// Read the markdown file
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		originalContent := string(content)
		updatedContent := originalContent

		// Replace all image references
		for oldPath, newPath := range imageMap {
			// Handle various markdown image syntax
			// ![alt text](/images/old.jpg) -> ![alt text](/images/old.avif)

			// Replace references with the full path
			oldRef := "/images/" + oldPath
			newRef := "/images/" + newPath
			updatedContent = strings.Replace(updatedContent, oldRef, newRef, -1)

			// Handle just the filename references
			oldFileName := filepath.Base(oldPath)
			newFileName := filepath.Base(newPath)
			updatedContent = strings.Replace(updatedContent, oldFileName, newFileName, -1)

			// Handle Hugo shortcodes: {{< figure src="/images/file.jpg" >}}
			updatedContent = strings.Replace(updatedContent,
				fmt.Sprintf(`src="/images/%s"`, oldPath),
				fmt.Sprintf(`src="/images/%s"`, newPath), -1)
		}

		// Only write the file if changes were made
		if updatedContent != originalContent {
			err = os.WriteFile(path, []byte(updatedContent), 0644)
			if err != nil {
				return err
			}
			updatedFileCount++
			fmt.Printf("Updated references in %s\n", path)
		}

		return nil
	})

	return updatedFileCount, err
}
