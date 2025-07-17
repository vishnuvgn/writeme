package cmd

import (
	"fmt"
	"os"
	"strings"
	"writeme/config"
	"writeme/helpers"

	"github.com/spf13/cobra"
)

var useAI bool

var noteCmd = &cobra.Command{
	Use:   "note {the note}",
	Short: "Add a note to NOTES.md",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		note := args[0]

		var cfg *config.Config
		if useAI {
			fmt.Println("AI flag is set. Validating config...")

			path, err := config.ResolveConfigPath()
			if err != nil {
				return fmt.Errorf("could not resolve config path: %w", err)
			}

			cfg, err = config.LoadConfig(path)
			if err != nil {
				return fmt.Errorf("could not load config: %w", err)
			}
		}

		// 1. Read NOTES.md
		content, err := os.ReadFile("NOTES.md")
		if err != nil {
			return fmt.Errorf("could not read NOTES.md: %w", err)
		}

		// 2. Ensure top-level heading
		contentStr, err := helpers.EnsureTopLevelHeading(string(content))
		if err != nil {
			return err
		}

		// 3. Parse headings
		tree := helpers.ParseHeadings(contentStr)
		// fmt.Printf("Parsed headings: %+v\n", tree)

		// 4. Select placement
		placement, err := helpers.SelectPlacement(tree)
		if err != nil {
			return fmt.Errorf("could not select placement: %w", err)
		}
		fmt.Printf("User selected placement: %v\n", placement)

		// 5. AI rewording
		if useAI {
			fmt.Println("AI flag is set. Processing note with AI...")
			note, err = helpers.RewordNote(cfg, note)
			if err != nil {
				return fmt.Errorf("could not reword note: %w", err)
			}
		}

		fmt.Printf("Final note to insert: %s\n", note)

		// 6. Insert note — but get back both:
		// - new content
		// - insertAt index for preview snippet
		newContent, insertAt := helpers.InsertNote(contentStr, placement, note)

		// 7. Build preview snippet around insertAt
		lines := strings.Split(newContent, "\n")

		start := insertAt - 2
		if start < 0 {
			start = 0
		}
		end := insertAt + 2
		if end >= len(lines) {
			end = len(lines) - 1
		}

		linesAbove := []string{}
		linesBelow := []string{}
		for i := start; i < insertAt; i++ {
			linesAbove = append(linesAbove, lines[i])
		}
		for i := insertAt + 1; i <= end; i++ {
			linesBelow = append(linesBelow, lines[i])
		}

		// Extract the inserted line WITHOUT the "- " so user edits just the text
		insertedLine := strings.TrimSpace(lines[insertAt])
		insertedLine = strings.TrimPrefix(insertedLine, "- ")

		finalNote, confirmed, err := helpers.RunPreviewWithEdit(linesAbove, linesBelow, insertedLine)
		if err != nil {
			return fmt.Errorf("preview failed: %w", err)
		}

		if !confirmed {
			fmt.Println("Note insertion cancelled.")
			return nil
		}

		// 8. Rebuild the note line with "- " prefix
		lines[insertAt] = "- " + finalNote
		newContent = strings.Join(lines, "\n")

		// 9. Write to file if confirmed
		err = os.WriteFile("NOTES.md", []byte(newContent), 0644)
		if err != nil {
			return fmt.Errorf("could not write NOTES.md: %w", err)
		}

		fmt.Println("Note inserted!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(noteCmd)
	noteCmd.Flags().BoolVarP(&useAI, "ai", "a", false, "Use AI to process the note")
}
