package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
  Normal = iota
  Insert
  Replace
  VLine
  Line
)

type status int

type model struct {
  choices  []string           // items on the to-do list
  cursorY   int                // which to-do list item our cursor is pointing at
  cursorX   int
  heighSize int
  screen int
  selected map[int]struct{}   // which to-do items are selected
  status status
  filename string
}

var (
  fileToString []string
)

func ReadFile(filename string ) {
  cat := exec.Command("cat", filename)
  output, err := cat.Output()
  if err != nil {
    panic("A error with cat")
  }
  fileToString = strings.Split((string(output)),"\n")

}

func WriteFile(filename string, data []string) {
  fileData, err := os.OpenFile(filename,os.O_WRONLY,os.ModeAppend)
  if err != nil {
    panic("you don't open this file")
  }
  defer fileData.Close()
  _, err = fileData.WriteString(strings.Join(data, "\n"))
  if err != nil {
    panic(err.Error())
  }

}

func InitialModel(filename string) model {
  
  ReadFile(filename)
  
	return model{
		// Our to-do list is a grocery list
		choices:  fileToString,
    cursorY:  0,
    cursorX: 1,
    status:  1,
    screen: 22,
    filename: filename,
    heighSize:  len(fileToString),
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
  // Just return `nil`, which means "no I/O right now, please."
  return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

  switch msg := msg.(type) {

    // Is it a key press?
  case tea.KeyMsg:

    // Cool, what was the actual key pressed?
    switch msg.String() {

    // These keys should exit the program.
    case "alt+c":
      return m, tea.Quit
    case "ctrl+x":
      WriteFile(m.filename, m.choices)
      return m, tea.Quit
 
    // The "up" and "k" keys move the cursor up
    case "up":
  
  
      if m.cursorY >0 {
        m.cursorY--
      } 
      if m.screen >=22 {
        m.screen--
      }
      if m.cursorX == len(m.choices[m.cursorY])  {
        m.cursorX = 0
      }

     
        // The "down" and "j" keys move the cursor down
    case "down":

      if m.cursorY < len(m.choices)-1 {
        m.cursorY++

      }
       if m.screen < len(m.choices)-1 {
        m.screen++
      }
      
      if m.cursorX < len(m.choices[m.cursorY])  {
        m.cursorX = 0
      }


        
    case "right":
      if  m.cursorX == len(m.choices[m.cursorY])  {
        m.cursorY++
        m.cursorX = 0
      }
      if len(m.choices[m.cursorY]) >=0  {
        m.cursorX++
      }
    case "left":
      if  m.cursorX > 0 && m.cursorX < len(m.choices[m.cursorY]) {
        m.cursorX--
    } 
    
    case "alt+r":
      m.status = 3
    case "alt+n":
      m.status = 2

    case "a", "b", "c", "d", "f", "g","h","i","j","k","l","m","n","ñ","o","p","q","r","s","t","u","w","x","y","z", tea.KeySpace.String(): 
      
      switch m.status {
        case 2:
         if len(m.choices[m.cursorY]) > 0 {
            m.choices[m.cursorY] = fmt.Sprintf("%s%s%s",m.choices[m.cursorY][0:m.cursorX],msg.String(),m.choices[m.cursorY][m.cursorX:]) 

          }
        case 3:
          if len(m.choices[m.cursorY]) > 0 {
            m.choices[m.cursorY] = fmt.Sprintf("%s%s%s",m.choices[m.cursorY][0:m.cursorX-1],msg.String(),m.choices[m.cursorY][m.cursorX:]) 
          }
      }
    
    case tea.KeyBackspace.String():
      if len(m.choices[m.cursorY]) > 0 && m.cursorX >0 {
        m.choices[m.cursorY] = fmt.Sprintf("%s%s",m.choices[m.cursorY][0:m.cursorX-1],m.choices[m.cursorY][m.cursorX:])
        m.cursorX--
      }   
        
        
    // The "enter" key and the spacebar (a literal space) toggle
    // the selected state for the item that the cursor is pointing at.
    case "enter", " ":
      _, ok := m.selected[m.cursorY]
      if ok {
        delete(m.selected, m.cursorY)
      } else {
      m.selected[m.cursorY] = struct{}{}
      }
    }
  }

  // Return the updated model to the Bubble Tea runtime for processing.
  // Note that we're not returning a command.
  return m, nil
}

func (m model) View() string {
    // The header
    s := ""
    // Iterate over our choices
    for i, choice := range m.choices {
      s += fmt.Sprintf("%d ",i)
      for e,character := range choice {
        cursor := "\033[0m" // no cursor
        if  m.cursorY == i && m.cursorX == e {
          cursor = "\033[41m" // cursor!
        }
        if character == 3 {
          character = 'c'
        }
        s += fmt.Sprintf("%s%s",cursor,string(character))
      }
      s += " \n"
    }  
    // Send the UI for rendering
  data := strings.Split(s, "\n")
  return fmt.Sprintf("%s", strings.Join(data[0:m.screen-1],"\n"))
}


