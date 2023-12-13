package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)



func main() {
  file := os.Args[1]
  
  p := tea.NewProgram(InitialModel(file), tea.WithAltScreen(),tea.WithMouseCellMotion())
  if _, err := p.Run(); err != nil {
    fmt.Printf("Alas, there's been an error: %v", err)
    os.Exit(1)
  }
  
} 


