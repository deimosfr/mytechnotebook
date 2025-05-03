/*
CLAUDE AI Prompt to generate this code:

I want you to:

* create a go script to manage navigations in nav.md files present in the content/docs directory
* Store this script in scripts/update_navigation.go
* You have to update nav.md ONLY betwwen <!--# start navigation here --> and <!--# stop navigation here -->. You must not touch anything else in those files.

# nav.md instruction

* Look at the first level sub-directories name and their respective name to create a card containing:
 1. The link to the the folder, ex: /docs/linux
 2. The icon you'll be able to find in the _index.md file of each sub-directory and update it inside <i class="material-icons align-middle">. Ex: <i class="material-icons align-middle">article</i>
 3. Create links on a single line like this way: <a href="/docs/linux/applications/">Applications</a> • <a href="/docs/linux/coding/">Coding</a> • <a href="/docs/linux/desktop/">Desktop</a></li> • <a href="/docs/linux/filesystems-and-storage/">Filesystems & Storage</a>

Here is a model I want you to use:
```
<div class="row flex-xl-wrap pb-4">

<div id="list-item" class="col-md-6 col-12 py-2">

	<a class="text-decoration-none text-reset" href="/docs/linux/">
	<div class="card h-100 features feature-full-bg rounded p-4 position-relative overflow-hidden border-1">
	    <span class="h1 icon-color">
	      <i class="material-icons align-middle">article</i>
	    </span>
	    <div class="card-body p-0 content">
	      <p class="fs-5 fw-semibold card-title mb-1">Linux</p>
	      <p class="para card-text mb-0">
	      <ul>
	        <a href="/docs/linux/applications/">Applications</a> • <a href="/docs/linux/coding/">Coding</a> • <a href="/docs/linux/desktop/">Desktop</a></li> • <a href="/docs/linux/filesystems-and-storage/">Filesystems & Storage</a> • <a href="/docs/linux/firewalls/">Firewalls</a> • <a href="/docs/linux/kernel/">Kernel</a> • <a href="/docs/linux/misc/">Misc</a> • <a href="/docs/linux/multimedia/">Multimedia</a> • <a href="/docs/linux/network/">Network</a> • <a href="/docs/linux/packages/">Packages</a> • <a href="/docs/linux/security/">Security</a>
	      </ul>
	      </p>
	    </div>
	  </div>
	</a>

</div>

<div id="list-item" class="col-md-6 col-12 py-2">

	<a class="text-decoration-none text-reset" href="../guides/landing-page/overview/">
	  <div class="card h-100 features feature-full-bg rounded p-4 position-relative overflow-hidden border-1">
	    <span class="h1 icon-color">
	      <i class="material-icons align-middle">flight_land</i>
	    </span>
	    <div class="card-body p-0 content">
	      <p class="fs-5 fw-semibold card-title mb-1">Landing Page</p>
	      <p class="para card-text mb-0">Customizable landing page and templates</p>
	    </div>
	  </div>
	</a>

</div>
</div>
```
 4. Sort the directories in a case-insensitive way
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

// Navigation markers
const (
	StartNavMarker = "<!--# start navigation here -->"
	StopNavMarker  = "<!--# stop navigation here -->"
)

// Content from _index.md file
type IndexContent struct {
	Title    string
	Icon     string
	IconType string
}

// For case-insensitive sorting
type caseInsensitiveSort []string

func (s caseInsensitiveSort) Len() int {
	return len(s)
}

func (s caseInsensitiveSort) Less(i, j int) bool {
	return strings.ToLower(s[i]) < strings.ToLower(s[j])
}

func (s caseInsensitiveSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// readIndexFile reads the _index.md file and extracts the title and icon
func readIndexFile(path string) (IndexContent, error) {
	content := IndexContent{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return content, err
	}

	// Extract title and icon using regex
	titleRegex := regexp.MustCompile(`title: ["'](.+?)["']`)
	iconRegex := regexp.MustCompile(`icon: ["'](.+?)["']`)
	iconTypeRegex := regexp.MustCompile(`icontype: ["'](.+?)["']`)

	titleMatches := titleRegex.FindSubmatch(data)
	if len(titleMatches) > 1 {
		content.Title = string(titleMatches[1])
	}

	iconMatches := iconRegex.FindSubmatch(data)
	if len(iconMatches) > 1 {
		content.Icon = string(iconMatches[1])
	}

	iconTypeMatches := iconTypeRegex.FindSubmatch(data)
	if len(iconTypeMatches) > 1 {
		content.IconType = string(iconTypeMatches[1])
	}

	return content, nil
}

// generateMainNavHTML generates the HTML for the main navigation
func generateMainNavHTML(docsDir string) (string, error) {
	// Get all first-level subdirectories
	entries, err := os.ReadDir(docsDir)
	if err != nil {
		return "", err
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), "_") && !strings.HasPrefix(entry.Name(), ".") {
			dirs = append(dirs, entry.Name())
		}
	}
	// Use case-insensitive sort
	sort.Sort(caseInsensitiveSort(dirs))

	var htmlBuilder strings.Builder
	htmlBuilder.WriteString("<div class=\"row flex-xl-wrap pb-4\">\n\n")

	// Process each directory to create a card
	count := 0
	for _, dir := range dirs {
		dirPath := filepath.Join(docsDir, dir)
		indexPath := filepath.Join(dirPath, "_index.md")

		// Skip if _index.md doesn't exist
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			continue
		}

		// Read index content
		idxContent, err := readIndexFile(indexPath)
		if err != nil {
			fmt.Printf("Warning: Could not read %s: %v\n", indexPath, err)
			continue
		}

		// Default icon if not found
		if idxContent.Icon == "" {
			idxContent.Icon = "article"
		}

		// Use directory name as fallback title
		title := idxContent.Title
		if title == "" {
			title = dir
		}

		// URL-friendly directory name (lowercase, replace spaces with dashes)
		urlDir := strings.ToLower(strings.ReplaceAll(dir, " ", "-"))

		// Start card HTML
		htmlBuilder.WriteString("<div id=\"list-item\" class=\"col-md-6 col-12 py-2\">\n")
		htmlBuilder.WriteString(fmt.Sprintf("  <a class=\"text-decoration-none text-reset\" href=\"/docs/%s/\">\n", strings.ToLower(urlDir)))
		htmlBuilder.WriteString("  <div class=\"card h-100 features feature-full-bg rounded p-4 position-relative overflow-hidden border-1\">\n")
		htmlBuilder.WriteString("      <span class=\"h1 icon-color\">\n")
		if idxContent.IconType == "simple" {
			htmlBuilder.WriteString(fmt.Sprintf("        <i class=\"si si-%s align-middle\"></i>\n", idxContent.Icon))
		} else {
			htmlBuilder.WriteString(fmt.Sprintf("        <i class=\"material-icons align-middle\">%s</i>\n", idxContent.Icon))
		}
		htmlBuilder.WriteString("      </span>\n")
		htmlBuilder.WriteString("      <div class=\"card-body p-0 content\">\n")
		htmlBuilder.WriteString(fmt.Sprintf("        <p class=\"fs-5 fw-semibold card-title mb-1\">%s</p>\n", title))
		htmlBuilder.WriteString("        <p class=\"para card-text mb-0\">\n")
		htmlBuilder.WriteString("        <ul>\n")

		// Get subdirectories for links
		subEntries, err := os.ReadDir(dirPath)
		if err == nil {
			var subDirs []string
			for _, subEntry := range subEntries {
				if subEntry.IsDir() && !strings.HasPrefix(subEntry.Name(), "_") && !strings.HasPrefix(subEntry.Name(), ".") {
					subDirs = append(subDirs, subEntry.Name())
				}
			}
			// Use case-insensitive sort for subdirectories too
			sort.Sort(caseInsensitiveSort(subDirs))

			// Create links
			var links []string
			for _, subDir := range subDirs {
				// URL-friendly subdirectory name (lowercase, replace spaces with dashes)
				urlSubDir := strings.ToLower(strings.ReplaceAll(subDir, " ", "-"))
				linkText := subDir

				// Try to get a better title from _index.md if available
				subIndexPath := filepath.Join(dirPath, subDir, "_index.md")
				if _, err := os.Stat(subIndexPath); !os.IsNotExist(err) {
					subIdxContent, err := readIndexFile(subIndexPath)
					if err == nil && subIdxContent.Title != "" {
						linkText = subIdxContent.Title
					}
				}

				links = append(links, fmt.Sprintf("<a href=\"/docs/%s/%s/\">%s</a>", strings.ToLower(urlDir), urlSubDir, linkText))
			}

			if len(links) > 0 {
				htmlBuilder.WriteString("          " + strings.Join(links, " • ") + "\n")
			}
		}

		htmlBuilder.WriteString("        </ul>\n")
		htmlBuilder.WriteString("        </p>\n")
		htmlBuilder.WriteString("      </div>\n")
		htmlBuilder.WriteString("    </div>\n")
		htmlBuilder.WriteString("  </a>\n")
		htmlBuilder.WriteString("</div>\n\n")

		count++
	}

	htmlBuilder.WriteString("</div>\n")

	return htmlBuilder.String(), nil
}

// updateNavigation updates the navigation section in the given nav file
func updateNavigation(navFilePath, navContent string) error {
	// Read the entire file
	file, err := os.Open(navFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Find the navigation markers and replace content between them
	var newContent []string
	inNavSection := false
	for _, line := range lines {
		if strings.TrimSpace(line) == StartNavMarker {
			newContent = append(newContent, line)
			newContent = append(newContent, "")
			newContent = append(newContent, navContent)
			inNavSection = true
		} else if strings.TrimSpace(line) == StopNavMarker {
			inNavSection = false
			newContent = append(newContent, line)
		} else if !inNavSection {
			newContent = append(newContent, line)
		}
	}

	// Write the updated content back to the file
	return ioutil.WriteFile(navFilePath, []byte(strings.Join(newContent, "\n")), 0644)
}

func main() {
	// Base directory path
	baseDir := "content/docs"

	// Get absolute path
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Check if the base directory exists
	if _, err := os.Stat(absBaseDir); os.IsNotExist(err) {
		fmt.Printf("Error: Directory %s does not exist\n", absBaseDir)
		os.Exit(1)
	}

	// Generate navigation HTML for nav.md
	mainNavHTML, err := generateMainNavHTML(absBaseDir)
	if err != nil {
		fmt.Printf("Error generating main navigation: %v\n", err)
		os.Exit(1)
	}

	// Update nav.md
	mainNavPath := filepath.Join(absBaseDir, "nav.md")
	if _, err := os.Stat(mainNavPath); os.IsNotExist(err) {
		fmt.Printf("Warning: %s does not exist, skipping\n", mainNavPath)
	} else {
		err = updateNavigation(mainNavPath, mainNavHTML)
		if err != nil {
			fmt.Printf("Error updating %s: %v\n", mainNavPath, err)
		} else {
			fmt.Printf("Successfully updated %s\n", mainNavPath)
		}
	}

	// Note: We'll implement sub-nav.md later as mentioned in requirements
	fmt.Println("Navigation update completed")
}
