package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	// Read output.html
	html, err := os.ReadFile("output.html")
	if err != nil {
		log.Fatalf("failed to read output.html: %v", err)
	}
	fmt.Println("HTML content:\n", string(html))

	// Read output.md
	md, err := os.ReadFile("output.md")
	if err != nil {
		log.Fatalf("failed to read output.md: %v", err)
	}
	mdContent := string(md)
	fmt.Println("Markdown content:\n", mdContent)

	// Get all of ![][image{index}] using regex
	mdPattern := `\[image\d+\]:\s*<[^>]+>`
	mdRegex := regexp.MustCompile(mdPattern)
	mdMatches := mdRegex.FindAllString(mdContent, -1)

	fmt.Println("Markdown image references:")
	for i, match := range mdMatches {
		fmt.Printf("%d: %s\n", i+1, match)
	}

	// Get all HTML image tags in the HTML file using regex
	htmlPattern := `<img\s+[^>]*src="([^"]*)"`
	htmlRegex := regexp.MustCompile(htmlPattern)
	htmlMatches := htmlRegex.FindAllStringSubmatch(string(html), -1)

	fmt.Println("HTML image URLs:")
	for i, match := range htmlMatches {
		fmt.Printf("%d: %s\n", i+1, match[1]) // match[1] contains the src URL
	}

	// Replace [image1]: <url> with the corresponding URL from HTML in the markdown content
	for i, match := range mdMatches {
		if i < len(htmlMatches) {
			// Construct the replacement string [image{index}]: <url>
			imageReference := fmt.Sprintf("[image%d]: %s", i+1, htmlMatches[i][1])
			mdContent = strings.Replace(mdContent, match, imageReference, 1)
		}
	}

	// Save the modified markdown content to a new file
	err = os.WriteFile("updated_output.md", []byte(mdContent), 0644)
	if err != nil {
		log.Fatalf("failed to write updated_output.md: %v", err)
	}
	fmt.Println("Updated markdown content saved to updated_output.md")
}
