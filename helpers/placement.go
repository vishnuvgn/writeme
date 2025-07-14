package helpers

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

// EnsureTopLevelHeading checks for a top-level heading and adds it if missing.
func EnsureTopLevelHeading(content string) (string, error) {
	if strings.Contains(content, "# ") {
		// Already has a top-level heading
		return content, nil
	}

	// Get current directory name
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get working directory: %w", err)
	}
	dirName := filepath.Base(wd)

	// Add # {dirname} at the top
	newContent := fmt.Sprintf("# %s\n%s", dirName, content)

	// Write it back immediately
	err = os.WriteFile("NOTES.md", []byte(newContent), 0644)
	if err != nil {
		return "", fmt.Errorf("could not write NOTES.md: %w", err)
	}

	return newContent, nil
}

type HeadingNode struct {
	Level    int
	Title    string
	Children []*HeadingNode
}

func ParseHeadings(content string) *HeadingNode {
	root := &HeadingNode{
		Level: 0,
		Title: "ROOT", // pseudo-root
	}

	stack := []*HeadingNode{root}

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			continue
		}

		level := 0
		for level < len(line) && line[level] == '#' {
			level++
		}

		// Only consider headings like "# title"
		if len(line) <= level || line[level] != ' ' {
			continue
		}

		title := strings.TrimSpace(line[level+1:])

		node := &HeadingNode{
			Level: level,
			Title: title,
		}

		// Pop stack until we find the correct parent
		for len(stack) > 0 && stack[len(stack)-1].Level >= level {
			stack = stack[:len(stack)-1]
		}

		parent := stack[len(stack)-1]
		parent.Children = append(parent.Children, node)

		stack = append(stack, node)
	}

	return root
}

// SelectPlacement walks the tree and returns the selected path.
func SelectPlacement(node *HeadingNode) ([]string, error) {
	var path []string
	current := node

	for {
		if len(current.Children) == 0 {
			if current.Title != "ROOT" {
				path = append(path, current.Title)
			}
			return path, nil
		}

		var options []string
		for _, child := range current.Children {
			options = append(options, child.Title)
		}

		// Only offer INSERT AT THIS LEVEL if not ROOT
		if len(options) > 0 && current.Title != "ROOT" {
			options = append(options, "INSERT AT THIS LEVEL")
		}

		prompt := promptui.Select{
			Label: fmt.Sprintf("Choose section under '%s'", current.Title),
			Items: options,
		}

		_, result, err := prompt.Run()
		if err != nil {
			return nil, fmt.Errorf("prompt failed: %w", err)
		}

		if result == "INSERT AT THIS LEVEL" {
			if current.Title != "ROOT" {
				path = append(path, current.Title)
			}
			return path, nil
		}

		// Drill down
		for _, child := range current.Children {
			if child.Title == result {
				if current.Title != "ROOT" {
					path = append(path, current.Title)
				}
				current = child
				break
			}
		}
	}
}
