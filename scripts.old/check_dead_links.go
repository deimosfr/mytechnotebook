/*
CLAUDE AI Prompt to generate this code:

I want you to create a go script in scripts/check_dead_links.go that will:
* check all markdown static links in all markdown files present in the "content" folder and all its sub directories
* ensure it points to a file in the "static" folder and that the file exists. Raise warning if it's not found
* ignore all _index.md files
* If files are present in the static folder but not referenced inside a markdown file, then add a --fix argument that will:
 1. Remove broken links by keeping the name anyway. Ex: [Name](http://broken_links) will be converted into Name
 2. If a link is inside the "Resources" contents, then remove the link

* The default behavior should be to only check the internal links
* Add another argument --with-external-links that will fix the links
* if a link is not working in http mode, try https. If it works, update the link in the markdown file

* Add a progress bar on the number of links currently being checked
* If there are no dead links, then exit 0, if there are dead links, exit 1
* In the root folder, use a ignore_links file, where you should read the list of links to ignore as dead
* When you find links inside ignore_links file, do not perform checks
* Add a separate --rename-files flag that will ONLY handle files with special characters (accents, parentheses, etc.):
 1. Create sanitized versions of filenames by replacing accented characters with non-accented equivalents
 2. Replace spaces and parentheses with underscores
 3. Rename the actual files in the static directory
 4. Update all markdown links to point to the renamed files

* Keep the functionalities of --fix and --rename-files completely separate
* At the beginning of the script, add in the comments the prompt you used to generate this code
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode"
)

// Configuration
const (
	maxConcurrentRequests  = 5                      // Maximum number of concurrent HTTP requests
	httpTimeout            = 10 * time.Second       // Timeout for HTTP requests
	ignoreListFile         = "ignore_links"         // File containing links to ignore
	progressUpdateInterval = 500 * time.Millisecond // How often to update the progress bar
)

// Link represents a markdown link with its source file
type Link struct {
	Text       string
	URL        string
	SourceFile string
	Position   int
}

// DeadLink represents a link that failed validation
type DeadLink struct {
	Text            string
	URL             string
	SourceFile      string
	Position        int
	ErrorCode       int
	Error           string
	FixedHTTPSLink  string // Store the HTTPS version if it works
	SanitizedPath   string // Store the sanitized path for renaming
	OriginalFileExt string // Store the original file extension
}

// Store results for unique links to avoid checking the same URL multiple times
var checkedLinks = make(map[string]bool)
var checkedLinksMutex sync.Mutex

// Progress tracking
var (
	totalLinks     int32      // Total number of links to check
	processedLinks int32      // Number of links processed so far
	inProgress     bool       // Whether link checking is in progress
	progressMutex  sync.Mutex // Mutex for progress updates
)

// Regular expression to find markdown links
// This matches both [text](url) and [text](url "title") formats
var linkRegex = regexp.MustCompile(`\[([^\]]*)\]\(([^)"'\s]+)(?:\s+["']([^"']*)["'])?\)`)

// sanitizeFilename removes or replaces special characters to create a file-system safe filename
func sanitizeFilename(filename string) string {
	// Create a mapping for accented characters
	replacements := map[rune]string{
		'é': "e", 'è': "e", 'ê': "e", 'ë': "e",
		'à': "a", 'â': "a", 'ä': "a",
		'î': "i", 'ï': "i",
		'ô': "o", 'ö': "o",
		'ù': "u", 'û': "u", 'ü': "u",
		'ç': "c",
		'ñ': "n",
		'É': "E", 'È': "E", 'Ê': "E", 'Ë': "E",
		'À': "A", 'Â': "A", 'Ä': "A",
		'Î': "I", 'Ï': "I",
		'Ô': "O", 'Ö': "O",
		'Ù': "U", 'Û': "U", 'Ü': "U",
		'Ç': "C",
		'Ñ': "N",
	}

	var sb strings.Builder
	for _, r := range filename {
		if replacement, ok := replacements[r]; ok {
			sb.WriteString(replacement)
		} else if r == '(' || r == ')' || r == ' ' {
			sb.WriteRune('_')
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '-' || r == '_' || r == '/' {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func main() {
	// Parse command line flags
	fixFlag := flag.Bool("fix", false, "Fix broken links automatically")
	withExternalLinks := flag.Bool("with-external-links", false, "Check external links as well as internal links")
	renameFilesFlag := flag.Bool("rename-files", false, "Rename files with special characters and update links")
	flag.Parse()

	// Get the current working directory
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	contentDir := filepath.Join(rootDir, "content")
	staticDir := filepath.Join(rootDir, "static")
	ignoreFilePath := filepath.Join(rootDir, ignoreListFile)

	// Load ignored links
	ignoredLinks := loadIgnoredLinks(ignoreFilePath)
	fmt.Printf("Loaded %d links to ignore\n", len(ignoredLinks))

	// Find all markdown files
	var markdownFiles []string
	err = filepath.WalkDir(contentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Skip _index.md files
		if d.Name() == "_index.md" {
			return nil
		}
		if !d.IsDir() && (strings.HasSuffix(path, ".md") || strings.HasSuffix(path, ".markdown")) {
			markdownFiles = append(markdownFiles, path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking content directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d markdown files to check\n", len(markdownFiles))

	// Extract and check all links
	allLinks := extractLinks(markdownFiles, rootDir)
	fmt.Printf("Found %d unique links to check\n", len(allLinks))

	// Set up progress counter
	totalLinks = int32(len(allLinks))
	processedLinks = 0

	// Start progress display in a separate goroutine
	inProgress = true
	go displayProgress()

	// Check links and collect dead links
	deadLinks := checkLinks(allLinks, ignoredLinks, rootDir, *withExternalLinks, *renameFilesFlag, staticDir)

	// Stop progress display
	progressMutex.Lock()
	inProgress = false
	progressMutex.Unlock()
	// Give progress display time to show 100%
	time.Sleep(100 * time.Millisecond)

	// Find static files that are not referenced
	unreferencedFiles := findUnreferencedStaticFiles(staticDir, allLinks, rootDir)
	if len(unreferencedFiles) > 0 {
		fmt.Printf("\nFound %d unreferenced files in static folder:\n", len(unreferencedFiles))
		for _, file := range unreferencedFiles {
			relPath, _ := filepath.Rel(staticDir, file)
			fmt.Printf("- %s\n", relPath)
		}
	}

	// Separate handling for --fix and --rename-files
	// Handle file renaming if requested
	if *renameFilesFlag && len(deadLinks) > 0 {
		fmt.Printf("\nRenaming files with special characters...\n")
		renamedCount := renameSpecialCharFiles(deadLinks, staticDir, rootDir)
		fmt.Printf("Renamed %d files and updated their links\n", renamedCount)
	}

	// Fix broken links if requested (separate from file renaming)
	if *fixFlag && len(deadLinks) > 0 {
		fmt.Printf("\nFixing %d broken links...\n", len(deadLinks))
		fixedFiles := fixBrokenLinks(deadLinks, false, staticDir, rootDir) // Pass false to disable file renaming in fixBrokenLinks
		fmt.Printf("Fixed links in %d files\n", len(fixedFiles))
	}

	// Print results
	if len(deadLinks) > 0 {
		fmt.Printf("\nFound %d dead links:\n", len(deadLinks))
		for _, deadLink := range deadLinks {
			relPath, _ := filepath.Rel(rootDir, deadLink.SourceFile)
			if deadLink.SanitizedPath != "" {
				fmt.Printf("- Link contains special characters: %s -> %s (in %s)\n",
					deadLink.URL, deadLink.SanitizedPath, relPath)
			} else if deadLink.FixedHTTPSLink != "" {
				fmt.Printf("- HTTP link can be fixed with HTTPS: %s -> %s (in %s)\n",
					deadLink.URL, deadLink.FixedHTTPSLink, relPath)
			} else if deadLink.ErrorCode != 0 {
				fmt.Printf("- [%d] %s (in %s)\n", deadLink.ErrorCode, deadLink.URL, relPath)
			} else {
				fmt.Printf("- %s: %s (in %s)\n", deadLink.Error, deadLink.URL, relPath)
			}
		}
		os.Exit(1)
	} else {
		fmt.Println("\nAll links are valid!")
		os.Exit(0)
	}
}

// loadIgnoredLinks reads the ignore_links file and returns a map of URLs to ignore
func loadIgnoredLinks(ignoreFilePath string) map[string]bool {
	ignoredLinks := make(map[string]bool)

	file, err := os.Open(ignoreFilePath)
	if err != nil {
		// If the file doesn't exist, return an empty map
		if os.IsNotExist(err) {
			return ignoredLinks
		}
		fmt.Printf("Warning: Could not open ignore list file: %v\n", err)
		return ignoredLinks
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		link := strings.TrimSpace(scanner.Text())
		if link != "" && !strings.HasPrefix(link, "#") {
			ignoredLinks[link] = true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Warning: Error reading ignore list file: %v\n", err)
	}

	return ignoredLinks
}

// extractLinks extracts all links from markdown files
func extractLinks(markdownFiles []string, rootDir string) []Link {
	var allLinks []Link
	uniqueLinks := make(map[string]bool)

	for _, mdFile := range markdownFiles {
		// Read the file
		content, err := os.ReadFile(mdFile)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", mdFile, err)
			continue
		}

		// Find all links in the markdown file
		contentStr := string(content)
		matches := linkRegex.FindAllSubmatchIndex(content, -1)
		for _, match := range matches {
			if len(match) < 6 { // FindAllSubmatchIndex returns indexes in pairs
				continue
			}

			// Get the link URL and text using the indices
			text := contentStr[match[2]:match[3]]
			link := contentStr[match[4]:match[5]]
			position := match[0] // Start position of the match

			// Skip empty links
			if link == "" {
				continue
			}

			// Create a unique key for the link to avoid duplicates
			linkKey := link + "|" + mdFile

			// Only add the link if we haven't seen it before
			if !uniqueLinks[linkKey] {
				uniqueLinks[linkKey] = true
				allLinks = append(allLinks, Link{
					Text:       text,
					URL:        link,
					SourceFile: mdFile,
					Position:   position,
				})
			}
		}
	}

	return allLinks
}

// displayProgress shows a progress bar for link checking
func displayProgress() {
	terminalWidth := 80 // Default terminal width

	for {
		progressMutex.Lock()
		if !inProgress {
			progressMutex.Unlock()
			break
		}
		progressMutex.Unlock()

		total := atomic.LoadInt32(&totalLinks)
		processed := atomic.LoadInt32(&processedLinks)

		if total > 0 {
			percentage := float64(processed) / float64(total) * 100

			// Create progress bar
			barWidth := terminalWidth - 30 // Allow space for text
			completedWidth := int(float64(barWidth) * float64(processed) / float64(total))
			if completedWidth > barWidth {
				completedWidth = barWidth
			}

			progressBar := "["
			for i := 0; i < barWidth; i++ {
				if i < completedWidth {
					progressBar += "="
				} else {
					progressBar += " "
				}
			}
			progressBar += "]"

			// Clear line and print progress
			fmt.Printf("\r\033[K%s %.1f%% (%d/%d)", progressBar, percentage, processed, total)
		}

		time.Sleep(progressUpdateInterval)
	}

	// Print newline after progress is complete
	fmt.Println()
}

// findUnreferencedStaticFiles finds files in the static directory that aren't referenced in any markdown file
func findUnreferencedStaticFiles(staticDir string, links []Link, rootDir string) []string {
	var unreferencedFiles []string
	referencedFiles := make(map[string]bool)

	// Build a map of all referenced static files
	for _, link := range links {
		// Only process links that might point to static files
		if strings.HasPrefix(link.URL, "/") ||
			(!strings.HasPrefix(link.URL, "http://") &&
				!strings.HasPrefix(link.URL, "https://") &&
				!strings.HasPrefix(link.URL, "#")) {

			var staticPath string
			if strings.HasPrefix(link.URL, "/") {
				// For absolute paths
				staticPath = filepath.Join(staticDir, strings.TrimPrefix(link.URL, "/"))
			} else {
				// For relative paths, try to resolve to static directory
				staticPath = filepath.Join(staticDir, link.URL)
			}

			// Normalize the path
			staticPath = filepath.Clean(staticPath)
			referencedFiles[staticPath] = true
		}
	}

	// Walk the static directory and find files not in the referencedFiles map
	_ = filepath.WalkDir(staticDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if !d.IsDir() {
			if !referencedFiles[path] {
				unreferencedFiles = append(unreferencedFiles, path)
			}
		}
		return nil
	})

	return unreferencedFiles
}

// checkLinks validates all links and returns a list of dead links
func checkLinks(links []Link, ignoredLinks map[string]bool, rootDir string, checkExternal bool, renameFiles bool, staticDir string) []DeadLink {
	var deadLinks []DeadLink
	var deadLinksMutex sync.Mutex

	// Create a worker pool for HTTP requests
	sem := make(chan struct{}, maxConcurrentRequests)
	var wg sync.WaitGroup

	for _, link := range links {
		// Skip ignored links
		if ignoredLinks[link.URL] {
			// Increment the counter for skipped links too
			atomic.AddInt32(&processedLinks, 1)
			continue
		}

		// Process each link
		wg.Add(1)
		go func(link Link) {
			defer wg.Done()

			// If it's an external link (starts with http:// or https://)
			if strings.HasPrefix(link.URL, "http://") || strings.HasPrefix(link.URL, "https://") {
				if !checkExternal {
					// Skip external links if not checking them
					atomic.AddInt32(&processedLinks, 1)
					return
				}

				// Acquire a token from the semaphore
				sem <- struct{}{}
				defer func() { <-sem }()

				// Check if we've already validated this URL
				checkedLinksMutex.Lock()
				isValid, checked := checkedLinks[link.URL]
				checkedLinksMutex.Unlock()

				if checked {
					if !isValid {
						deadLinksMutex.Lock()
						deadLinks = append(deadLinks, DeadLink{
							Text:       link.Text,
							URL:        link.URL,
							SourceFile: link.SourceFile,
							Position:   link.Position,
							ErrorCode:  0,
							Error:      "Previously found to be invalid",
						})
						deadLinksMutex.Unlock()
					}
					// Increment the counter
					atomic.AddInt32(&processedLinks, 1)
					return
				}

				// Make an HTTP request to check if the link is valid
				client := &http.Client{
					Timeout: httpTimeout,
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						// Allow up to 10 redirects
						if len(via) >= 10 {
							return fmt.Errorf("too many redirects")
						}
						return nil
					},
				}

				resp, err := client.Get(link.URL)
				if err != nil || resp != nil && resp.StatusCode != http.StatusOK {
					// Try HTTPS if this is an HTTP URL that failed
					var httpsURL string
					var httpsWorks bool

					if strings.HasPrefix(link.URL, "http://") {
						httpsURL = "https://" + strings.TrimPrefix(link.URL, "http://")

						// Try the HTTPS version
						httpsResp, httpsErr := client.Get(httpsURL)
						if httpsErr == nil && httpsResp.StatusCode == http.StatusOK {
							httpsWorks = true
							httpsResp.Body.Close()
						}
					}

					// Mark this URL as invalid
					checkedLinksMutex.Lock()
					checkedLinks[link.URL] = false
					checkedLinksMutex.Unlock()

					deadLinksMutex.Lock()
					deadLink := DeadLink{
						Text:       link.Text,
						URL:        link.URL,
						SourceFile: link.SourceFile,
						Position:   link.Position,
					}

					// Set error information
					if err != nil {
						deadLink.Error = err.Error()
					} else if resp != nil {
						deadLink.ErrorCode = resp.StatusCode
						deadLink.Error = resp.Status
						resp.Body.Close()
					}

					// If HTTPS worked, record the fixed URL
					if httpsWorks {
						deadLink.FixedHTTPSLink = httpsURL
					}

					deadLinks = append(deadLinks, deadLink)
					deadLinksMutex.Unlock()

					// Increment the counter
					atomic.AddInt32(&processedLinks, 1)
					return
				}

				if resp != nil {
					// Mark this URL as valid
					checkedLinksMutex.Lock()
					checkedLinks[link.URL] = true
					checkedLinksMutex.Unlock()

					resp.Body.Close()
				}

				// Increment the counter
				atomic.AddInt32(&processedLinks, 1)
			} else if strings.HasPrefix(link.URL, "/") ||
				(!strings.HasPrefix(link.URL, "#") &&
					!strings.Contains(link.URL, ":")) {
				// This is a local file reference
				// For absolute paths starting with /, resolve from project root
				// For relative paths, resolve from the markdown file's directory

				var fullPath string
				if strings.HasPrefix(link.URL, "/") {
					// Remove the leading slash for path resolution
					fullPath = filepath.Join(rootDir, "static", strings.TrimPrefix(link.URL, "/"))
				} else {
					// For relative links, resolve from the markdown file's location
					basePath := filepath.Dir(link.SourceFile)
					fullPath = filepath.Join(basePath, link.URL)
				}

				// Cleanup the path
				fullPath = filepath.Clean(fullPath)

				// Check if the file exists
				_, err := os.Stat(fullPath)
				if os.IsNotExist(err) {
					// Try to find the file in the static directory as a fallback
					staticPath := filepath.Join(rootDir, "static", link.URL)
					_, err = os.Stat(staticPath)
					if os.IsNotExist(err) {
						// Only consider renaming files if they're in the specified directories
						shouldRename := renameFiles && isFilePath(link.URL) &&
							(containsSpecialChars(link.URL) || strings.Contains(link.URL, " "))

						if shouldRename {
							sanitizedURL := sanitizeFilename(link.URL)
							// Store the sanitized path with extension preserved
							extension := filepath.Ext(link.URL)

							deadLinksMutex.Lock()
							deadLinks = append(deadLinks, DeadLink{
								Text:            link.Text,
								URL:             link.URL,
								SourceFile:      link.SourceFile,
								Position:        link.Position,
								ErrorCode:       0,
								Error:           "File not found, special characters in filename",
								SanitizedPath:   sanitizedURL,
								OriginalFileExt: extension,
							})
							deadLinksMutex.Unlock()
						} else {
							deadLinksMutex.Lock()
							deadLinks = append(deadLinks, DeadLink{
								Text:       link.Text,
								URL:        link.URL,
								SourceFile: link.SourceFile,
								Position:   link.Position,
								ErrorCode:  0,
								Error:      "File not found",
							})
							deadLinksMutex.Unlock()
						}
					}
				}

				// Increment the counter for local files
				atomic.AddInt32(&processedLinks, 1)
			} else {
				// Skip anchor links (starting with #) and other protocols
				// Increment the counter for skipped links
				atomic.AddInt32(&processedLinks, 1)
			}
		}(link)
	}

	wg.Wait()
	return deadLinks
}

// containsSpecialChars checks if a string contains special characters
func containsSpecialChars(s string) bool {
	for _, r := range s {
		if r > 127 || r == '(' || r == ')' || r == '\'' || r == '"' || r == '`' || r == ' ' {
			return true
		}
	}
	return false
}

// fixBrokenLinks processes each broken link and fixes it
// - Regular links: [Text](broken-url) becomes Text
// - Links in Resources section: Remove the entire link
// - If renameFiles is true, tries to rename files with special characters
func fixBrokenLinks(deadLinks []DeadLink, renameFiles bool, staticDir, rootDir string) map[string]bool {
	// Track which files we've modified
	fixedFiles := make(map[string]bool)

	// Group dead links by file to process each file once
	linksByFile := make(map[string][]DeadLink)
	for _, link := range deadLinks {
		linksByFile[link.SourceFile] = append(linksByFile[link.SourceFile], link)
	}

	// Process each file
	for filePath, links := range linksByFile {
		// Read the file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s for fixing: %v\n", filePath, err)
			continue
		}

		// Convert to string for easier manipulation
		fileContent := string(content)

		// Sort links by position in reverse order to avoid position shifts
		// when we modify the file
		sort.Slice(links, func(i, j int) bool {
			return links[i].Position > links[j].Position
		})

		// Check if we're in a Resources section
		lines := strings.Split(fileContent, "\n")
		resourcesSectionLines := make(map[int]bool)

		inResourcesSection := false
		for i, line := range lines {
			// Check for headings that indicate a Resources section
			if strings.Contains(strings.ToLower(line), "## resources") ||
				strings.Contains(strings.ToLower(line), "## references") ||
				strings.Contains(strings.ToLower(line), "## liens") ||
				strings.Contains(strings.ToLower(line), "## links") {
				inResourcesSection = true
			} else if strings.HasPrefix(line, "##") && inResourcesSection {
				// End of resources section when we hit another heading
				inResourcesSection = false
			}

			if inResourcesSection {
				resourcesSectionLines[i] = true
			}
		}

		// Process each link
		modified := false
		for _, link := range links {
			// Find the line number of this link
			lineNumber := 0
			for i, char := range fileContent {
				if i == link.Position {
					break
				}
				if char == '\n' {
					lineNumber++
				}
			}

			// Determine if this link is in a Resources section
			inResources := resourcesSectionLines[lineNumber]

			// Create the original link pattern and its replacement
			originalLinkPattern := fmt.Sprintf("[%s](%s)", regexp.QuoteMeta(link.Text), regexp.QuoteMeta(link.URL))

			var replacement string
			if inResources {
				// In resources section, remove the link entirely
				replacement = ""
			} else if link.SanitizedPath != "" && renameFiles {
				// We need to rename a file with special characters

				var oldFile, newFile string
				// Handle absolute paths starting with /
				if strings.HasPrefix(link.URL, "/") {
					// Original file path
					oldFile = filepath.Join(staticDir, strings.TrimPrefix(link.URL, "/"))

					// New sanitized path
					newPathWithExt := link.SanitizedPath
					if link.OriginalFileExt != "" && !strings.HasSuffix(newPathWithExt, link.OriginalFileExt) {
						newPathWithExt = strings.TrimSuffix(newPathWithExt, filepath.Ext(newPathWithExt)) + link.OriginalFileExt
					}

					// Construct new file path
					newFile = filepath.Join(staticDir, strings.TrimPrefix(newPathWithExt, "/"))

					// Create replacement link with sanitized path
					replacement = fmt.Sprintf("[%s](/%s)", link.Text, strings.TrimPrefix(newPathWithExt, "/"))
				} else {
					// Relative paths
					oldFile = filepath.Join(filepath.Dir(link.SourceFile), link.URL)
					newFile = filepath.Join(filepath.Dir(link.SourceFile), link.SanitizedPath)

					// Create replacement link with sanitized path
					replacement = fmt.Sprintf("[%s](%s)", link.Text, link.SanitizedPath)
				}

				// Check if old file exists and try to rename it
				if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
					// Ensure target directory exists
					targetDir := filepath.Dir(newFile)
					if err := os.MkdirAll(targetDir, 0755); err != nil {
						fmt.Printf("Error creating directory %s: %v\n", targetDir, err)
					} else {
						// Rename the file
						if err := os.Rename(oldFile, newFile); err != nil {
							fmt.Printf("Error renaming file from %s to %s: %v\n", oldFile, newFile, err)
						} else {
							fmt.Printf("Renamed file: %s -> %s\n", oldFile, newFile)
						}
					}
				} else {
					// File doesn't exist, just update the link
					fmt.Printf("Original file %s not found, only fixing link\n", oldFile)
				}
			} else if link.FixedHTTPSLink != "" {
				// Replace HTTP with HTTPS if it works
				replacement = fmt.Sprintf("[%s](%s)", link.Text, link.FixedHTTPSLink)
			} else {
				// Regular section, keep the text
				replacement = link.Text
			}

			// Replace the link
			if strings.Contains(fileContent, originalLinkPattern) {
				fileContent = strings.Replace(fileContent, originalLinkPattern, replacement, 1)
				modified = true
			} else {
				// If exact pattern wasn't found, try a more flexible approach using regex
				linkRegex := regexp.MustCompile(fmt.Sprintf(`\[%s\]\(%s(?:\s+["'][^"']*["'])?\)`,
					regexp.QuoteMeta(link.Text), regexp.QuoteMeta(link.URL)))
				fileContent = linkRegex.ReplaceAllString(fileContent, replacement)
				modified = true
			}
		}

		// Write the modified content back to the file
		if modified {
			err = os.WriteFile(filePath, []byte(fileContent), 0644)
			if err != nil {
				fmt.Printf("Error writing fixed content to %s: %v\n", filePath, err)
				continue
			}
			fixedFiles[filePath] = true
		}
	}

	return fixedFiles
}

// renameSpecialCharFiles is specifically for handling files with special characters
// It renames files and updates links without removing broken links
func renameSpecialCharFiles(deadLinks []DeadLink, staticDir, rootDir string) int {
	renamedCount := 0

	// Get only the links with special characters that need to be renamed
	var specialCharLinks []DeadLink
	for _, link := range deadLinks {
		// Only consider links with a sanitized path that point to actual files in our target directories
		if link.SanitizedPath != "" && isFilePath(link.URL) {
			specialCharLinks = append(specialCharLinks, link)
		}
	}

	// Group by file
	linksByFile := make(map[string][]DeadLink)
	for _, link := range specialCharLinks {
		linksByFile[link.SourceFile] = append(linksByFile[link.SourceFile], link)
	}

	// Process each file
	for filePath, links := range linksByFile {
		// Read the file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s for renaming: %v\n", filePath, err)
			continue
		}

		// Convert to string for easier manipulation
		fileContent := string(content)

		// Sort links by position in reverse order to avoid position shifts
		sort.Slice(links, func(i, j int) bool {
			return links[i].Position > links[j].Position
		})

		// Process each link
		modified := false
		for _, link := range links {
			var oldFile, newFile string

			// Create the original link pattern and its replacement
			originalLinkPattern := fmt.Sprintf("[%s](%s)", regexp.QuoteMeta(link.Text), regexp.QuoteMeta(link.URL))

			// Handle absolute paths starting with /
			if strings.HasPrefix(link.URL, "/") {
				// Original file path
				oldFile = filepath.Join(staticDir, strings.TrimPrefix(link.URL, "/"))

				// New sanitized path
				newPathWithExt := link.SanitizedPath
				if link.OriginalFileExt != "" && !strings.HasSuffix(newPathWithExt, link.OriginalFileExt) {
					newPathWithExt = strings.TrimSuffix(newPathWithExt, filepath.Ext(newPathWithExt)) + link.OriginalFileExt
				}

				// Construct new file path
				newFile = filepath.Join(staticDir, strings.TrimPrefix(newPathWithExt, "/"))

				// Create replacement link with sanitized path
				replacement := fmt.Sprintf("[%s](/%s)", link.Text, strings.TrimPrefix(newPathWithExt, "/"))

				// Replace the link
				if strings.Contains(fileContent, originalLinkPattern) {
					fileContent = strings.Replace(fileContent, originalLinkPattern, replacement, 1)
					modified = true
				} else {
					// If exact pattern wasn't found, try a more flexible approach using regex
					linkRegex := regexp.MustCompile(fmt.Sprintf(`\[%s\]\(%s(?:\s+["'][^"']*["'])?\)`,
						regexp.QuoteMeta(link.Text), regexp.QuoteMeta(link.URL)))
					fileContent = linkRegex.ReplaceAllString(fileContent, replacement)
					modified = true
				}
			} else {
				// Relative paths
				oldFile = filepath.Join(filepath.Dir(link.SourceFile), link.URL)
				newFile = filepath.Join(filepath.Dir(link.SourceFile), link.SanitizedPath)

				// Create replacement link with sanitized path
				replacement := fmt.Sprintf("[%s](%s)", link.Text, link.SanitizedPath)

				// Replace the link
				if strings.Contains(fileContent, originalLinkPattern) {
					fileContent = strings.Replace(fileContent, originalLinkPattern, replacement, 1)
					modified = true
				} else {
					// If exact pattern wasn't found, try a more flexible approach using regex
					linkRegex := regexp.MustCompile(fmt.Sprintf(`\[%s\]\(%s(?:\s+["'][^"']*["'])?\)`,
						regexp.QuoteMeta(link.Text), regexp.QuoteMeta(link.URL)))
					fileContent = linkRegex.ReplaceAllString(fileContent, replacement)
					modified = true
				}
			}

			// Check if old file exists and try to rename it
			if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
				// Ensure target directory exists
				targetDir := filepath.Dir(newFile)
				if err := os.MkdirAll(targetDir, 0755); err != nil {
					fmt.Printf("Error creating directory %s: %v\n", targetDir, err)
				} else {
					// Rename the file
					if err := os.Rename(oldFile, newFile); err != nil {
						fmt.Printf("Error renaming file from %s to %s: %v\n", oldFile, newFile, err)
					} else {
						fmt.Printf("Renamed file: %s -> %s\n", oldFile, newFile)
						renamedCount++
					}
				}
			} else {
				// File doesn't exist, just update the link
				fmt.Printf("Original file %s not found, only fixing link\n", oldFile)
			}
		}

		// Write the modified content back to the file
		if modified {
			err = os.WriteFile(filePath, []byte(fileContent), 0644)
			if err != nil {
				fmt.Printf("Error writing renamed content to %s: %v\n", filePath, err)
			}
		}
	}

	return renamedCount
}

// isFilePath checks if a link points to a file in static directory that should be considered for renaming
func isFilePath(linkURL string) bool {
	// Only consider links that point to static files in specific directories
	filePrefixes := []string{
		"/images/",
		"/images",
		"/pdf/",
		"/pdf",
		"/others/",
		"/others",
	}

	for _, prefix := range filePrefixes {
		if strings.HasPrefix(linkURL, prefix) {
			return true
		}
	}

	return false
}
