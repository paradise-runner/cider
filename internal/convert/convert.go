package convert

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var (
	// Apple Notes converts H1 to styled divs - we need to convert back
	appleTitleRegex = regexp.MustCompile(`<div><b><span style="font-size: 24px">(.*?)</span></b></div>`)
)

// MarkdownToHTML converts markdown to HTML using pandoc
func MarkdownToHTML(markdown string) (string, error) {
	cmd := exec.Command("pandoc", "-f", "gfm", "-t", "html", "--wrap=none")
	cmd.Stdin = strings.NewReader(markdown)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("pandoc conversion failed: %s", stderr.String())
	}

	return out.String(), nil
}

// HTMLToMarkdown converts HTML to markdown using pandoc with post-processing
func HTMLToMarkdown(html string) (string, error) {
	// Pre-process: Convert Apple Notes title pattern to H1
	html = convertAppleTitleToH1(html)

	// Run pandoc conversion
	cmd := exec.Command("pandoc", "-f", "html", "-t", "gfm-raw_html", "--wrap=none")
	cmd.Stdin = strings.NewReader(html)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("pandoc conversion failed: %s", stderr.String())
	}

	result := out.String()

	// Post-process: clean up pandoc output
	result = removeNbsp(result)
	result = trimTrailingWhitespace(result)
	result = collapseBlankLines(result)

	return result, nil
}

// CheckPandocAvailable checks if pandoc is available in PATH
func CheckPandocAvailable() error {
	_, err := exec.LookPath("pandoc")
	if err != nil {
		return fmt.Errorf("pandoc not found: please install pandoc from https://pandoc.org/installing.html")
	}
	return nil
}

// convertAppleTitleToH1 converts Apple Notes title pattern to H1 tag
// Apple Notes: <div><b><span style="font-size: 24px">Title</span></b></div>
// Target: <h1>Title</h1>
func convertAppleTitleToH1(html string) string {
	return appleTitleRegex.ReplaceAllString(html, "<h1>$1</h1>")
}

// removeNbsp removes standalone &nbsp; lines that pandoc adds between lists
func removeNbsp(markdown string) string {
	lines := strings.Split(markdown, "\n")
	var result []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "&nbsp;" {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

// trimTrailingWhitespace removes trailing whitespace from each line
func trimTrailingWhitespace(markdown string) string {
	lines := strings.Split(markdown, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	return strings.Join(lines, "\n")
}

// collapseBlankLines collapses multiple consecutive blank lines into one
func collapseBlankLines(markdown string) string {
	// Replace sequences of 2+ newlines with just 2 newlines (one blank line)
	re := regexp.MustCompile(`\n{3,}`)
	return re.ReplaceAllString(markdown, "\n\n")
}
