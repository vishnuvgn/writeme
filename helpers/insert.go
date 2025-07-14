package helpers

import (
	"fmt"
	"strings"
)

func InsertNote(content string, placement []string, note string) (string, int) {
	lines := strings.Split(content, "\n")
	var newLines []string
	var currentPath []string

	insertAtIndex := -1

	i := 0
	for i < len(lines) {
		line := lines[i]

		// Update path
		if strings.HasPrefix(line, "#") {
			level := 0
			for level < len(line) && line[level] == '#' {
				level++
			}

			title := strings.TrimSpace(line[level+1:])

			if level <= len(currentPath) {
				currentPath = currentPath[:level-1]
			}
			currentPath = append(currentPath, title)
		}

		newLines = append(newLines, line)

		if samePath(currentPath, placement) {
			insertAt := len(newLines)
			foundBullet := false

			j := i + 1
			for j < len(lines) {
				nextLine := strings.TrimSpace(lines[j])
				if strings.HasPrefix(nextLine, "#") {
					break
				}
				if strings.HasPrefix(nextLine, "-") {
					foundBullet = true
					insertAt = len(newLines) + 1
				}
				newLines = append(newLines, lines[j])
				i++
				j++
			}

			if !foundBullet {
				if len(newLines) == 0 || strings.TrimSpace(newLines[len(newLines)-1]) != "" {
					newLines = append(newLines, "")
				}
			}

			noteLine := fmt.Sprintf("- %s", note)
			newLines = append(newLines[:insertAt], append([]string{noteLine}, newLines[insertAt:]...)...)

			insertAtIndex = insertAt

			// 🔑 After inserting, add the remaining lines!
			for k := i + 1; k < len(lines); k++ {
				newLines = append(newLines, lines[k])
			}

			break // done inserting
		}

		i++
	}

	return strings.Join(newLines, "\n"), insertAtIndex
}

func samePath(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
