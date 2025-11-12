package helpers

import "github.com/charmbracelet/glamour"

func RenderMarkdown(markdown string) (string, error) {
    // Create renderer with auto-detected terminal width
    r, err := glamour.NewTermRenderer(
        glamour.WithAutoStyle(), // Uses best style for your terminal
    )
    if err != nil {
        return "", err
    }
    
    return r.Render(markdown)
}