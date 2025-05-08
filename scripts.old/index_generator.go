/*

CLAUDE AI Prompt to generate this code:

Create a Go script that automatically updates Hugo documentation _index.md files with links to subdirectories. The script should:

1. Recursively scan all directories starting from "content/docs"
2. For each directory, identify all its direct subdirectories
3. Update the "description" field in each _index.md file to include Markdown links to all subdirectories
4. Format the links as: [Subdirectory name](./parent-directory/subdirectory-name)

Requirements:
- The first letter of each subdirectory name should be capitalized in the displayed text
- All links should be lowercase in the URL part
- Spaces in directory names should be replaced with hyphens in URLs
- Links should be separated by the "•" symbol
- Links should be relative (starting with "./")
- Only update files that have subdirectories; don't modify descriptions in directories without subdirectories
- Include the parent directory name in the link path
- If a file doesn't exist, create it with a default template

Example of desired output in a BSD _index.md file:
description: "[Kernel](./bsd/kernel) • [Misc](./bsd/misc) • [Network](./bsd/network) • [Security](./bsd/security)"

All code comments should be in English.

*/

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Structure to store directory information
type DirInfo struct {
	Path       string
	Name       string
	SubDirs    []string
	ParentPath string
}

func main() {
	// Starting point for the scan
	startDir := "content/docs"

	// Check if the directory exists
	if _, err := os.Stat(startDir); os.IsNotExist(err) {
		fmt.Printf("Directory %s does not exist\n", startDir)
		return
	}

	// Get all directories and subdirectories
	dirs, err := scanDirectories(startDir)
	if err != nil {
		fmt.Printf("Error scanning directories: %v\n", err)
		return
	}

	// For each directory, update its _index.md file
	for _, dir := range dirs {
		updateIndexFile(dir)
	}

	fmt.Println("Successfully updated all _index.md files!")
}

// Recursively scan directories from the starting point
func scanDirectories(startDir string) ([]DirInfo, error) {
	var result []DirInfo

	err := filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// If it's a directory (and not the root directory)
		if info.IsDir() && path != startDir {
			// Get the relative path from content/docs
			relPath, err := filepath.Rel(startDir, path)
			if err != nil {
				return err
			}

			// Get the current directory name and its parent
			dirName := filepath.Base(path)
			parentPath := filepath.Dir(path)
			parentRelPath, _ := filepath.Rel(startDir, parentPath)

			// Get direct subdirectories
			subDirs, err := getSubDirectories(path)
			if err != nil {
				return err
			}

			// Create DirInfo object and add it to the result
			dirInfo := DirInfo{
				Path:       relPath,
				Name:       dirName,
				SubDirs:    subDirs,
				ParentPath: parentRelPath,
			}
			result = append(result, dirInfo)
		}
		return nil
	})

	return result, err
}

// Get direct subdirectories of a directory
func getSubDirectories(dirPath string) ([]string, error) {
	var subDirs []string

	entries, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			subDirs = append(subDirs, entry.Name())
		}
	}

	// Sort subdirectories alphabetically
	sort.Strings(subDirs)
	return subDirs, nil
}

// Update the _index.md file in the specified directory
func updateIndexFile(dir DirInfo) {
	// Full path to the _index.md file
	indexPath := filepath.Join("content/docs", dir.Path, "_index.md")

	// Check if the file exists
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		fmt.Printf("File %s not found, creating...\n", indexPath)
		// If the file does not exist, create it with a default template
		createDefaultIndexFile(indexPath, dir)
		return
	}

	// Skip update if no subdirectories exist
	if len(dir.SubDirs) == 0 {
		fmt.Printf("Skipping %s - no subdirectories\n", indexPath)
		return
	}

	// Open the file for reading
	file, err := os.Open(indexPath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", indexPath, err)
		return
	}
	defer file.Close()

	// Read the file content
	var lines []string
	scanner := bufio.NewScanner(file)
	descriptionFound := false
	descPattern := regexp.MustCompile(`^description:\s*(.*)$`)

	for scanner.Scan() {
		line := scanner.Text()

		// If it's the description line
		if descPattern.MatchString(line) {
			descriptionFound = true
			// Replace the line with the new description
			line = formatDescriptionLine(dir)
		}

		lines = append(lines, line)
	}

	// If no description line was found, add it
	if !descriptionFound {
		// Find the appropriate place to add the description (after title or at the end of frontmatter)
		titleIndex := -1
		frontmatterEndIndex := -1

		for i, line := range lines {
			if strings.HasPrefix(line, "title:") {
				titleIndex = i
			} else if line == "---" && i > 0 {
				frontmatterEndIndex = i
			}
		}

		descLine := formatDescriptionLine(dir)

		// Insert description after the title or before the end of frontmatter
		if titleIndex >= 0 {
			newLines := append(lines[:titleIndex+1], append([]string{descLine}, lines[titleIndex+1:]...)...)
			lines = newLines
		} else if frontmatterEndIndex >= 0 {
			newLines := append(lines[:frontmatterEndIndex], append([]string{descLine}, lines[frontmatterEndIndex:]...)...)
			lines = newLines
		} else {
			// Otherwise, add it at the end (rare case)
			lines = append(lines, descLine)
		}
	}

	// Write the updated content to the file
	newContent := strings.Join(lines, "\n")
	err = ioutil.WriteFile(indexPath, []byte(newContent), 0644)
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", indexPath, err)
		return
	}

	fmt.Printf("Successfully updated file %s\n", indexPath)
}

// Create a default _index.md file
func createDefaultIndexFile(path string, dir DirInfo) {
	// Default content for a new _index.md file
	title := strings.Title(dir.Name)

	var content string
	if len(dir.SubDirs) > 0 {
		content = fmt.Sprintf("---\ntitle: %s\n%s\n---\n\nContent for %s",
			title, formatDescriptionLine(dir), title)
	} else {
		content = fmt.Sprintf("---\ntitle: %s\ndescription: \"Documentation for %s\"\n---\n\nContent for %s",
			title, title, title)
	}

	// Create the file
	err := ioutil.WriteFile(path, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", path, err)
		return
	}

	fmt.Printf("Successfully created file %s\n", path)
}

// Format the description line with links to subdirectories
func formatDescriptionLine(dir DirInfo) string {
	// If no subdirectories, return a default description
	if len(dir.SubDirs) == 0 {
		return fmt.Sprintf("description: \"Documentation for %s\"", strings.Title(dir.Name))
	}

	// Build links for each subdirectory
	var links []string

	// Get the current directory name for prefixing links
	dirName := strings.ToLower(strings.Replace(dir.Name, " ", "-", -1))

	for _, subDir := range dir.SubDirs {
		// Format name for display (first letter capitalized)
		displayName := strings.Title(subDir)

		// Process subdirectory name: lowercase and replace spaces with dashes
		formattedSubDir := strings.ToLower(strings.Replace(subDir, " ", "-", -1))

		// Create link path with "./" prefix and include the current directory name
		linkPath := fmt.Sprintf("./%s/%s", dirName, formattedSubDir)

		// Format the Markdown link
		link := fmt.Sprintf("[%s](%s)", displayName, linkPath)
		links = append(links, link)
	}

	// Combine links with the • separator
	description := strings.Join(links, " • ")
	return fmt.Sprintf("description: \"%s\"", description)
}
