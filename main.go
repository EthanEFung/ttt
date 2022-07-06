package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)
type cursor struct {
	x, y int
}

type model struct {
	grid [][]byte
	cursor cursor
	mark byte
	state string
}
func (m model) Init() tea.Cmd {
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			return initialModel(), nil
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor.y > 0 {
				m.cursor.y--
			}
		case "down", "j":
			if m.cursor.y < 2 {
				m.cursor.y++
			}
		case "left", "h":
			if m.cursor.x > 0 {
				m.cursor.x--
			} 
		case "right", "l":
			if m.cursor.x < 2 {
				m.cursor.x++
			}
		case "enter", " ":
			x, y := m.cursor.x, m.cursor.y
			if m.grid[y][x] != ' ' || m.state != "playable" {
				break;
			}
			m.grid[y][x] = m.mark
			if m.mark == 'x' {
				m.mark = 'o'
			} else {
				m.mark = 'x'
			}
		}
		// evaluate position
		m.state = evaluate(m.grid)
	}
	return m, nil

}

func evaluate(grid [][]byte) string {
	// check the rows for 3 of a kind
	for y := 0; y < 3; y++ {
		mark := grid[y][0]
		if mark == ' ' {
			continue
		}
		wins := true
		for x := 1; x < 3; x++ {
			// check to make sure
			if grid[y][x] != mark {
				wins = false
				break;
			}
		}
		if wins {
			return string(mark) + " wins"
		}
	}

	// check the columns for 3 of a kind
	for x := 0; x < 3; x++ {
		mark := grid[0][x]
		if mark == ' ' {
			continue
		}
		wins := true
		for y := 1; y < 3; y++ {
			if mark != grid[y][x] {
				wins = false
				break
			}
		}
		if wins {
			return string(mark) + " wins"
		}
	}

	// check the principal diagonal
	mark := grid[0][0]
	if mark != ' ' && grid[1][1] == mark && grid[2][2] == mark {
		return string(mark) + " wins"
	}

	// check the secondary diagonal
	mark = grid[2][0]
	if mark != ' ' &&  grid[1][1] == mark && grid[0][2] == mark {
		return string(mark) + " wins"
	}
	// ^ of course these are hard coded, which will need to change in the event of a larger board

	// check for draw
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if grid[y][x] == ' ' {
				return "playable"
			}
		}
	}
	return "draw"
}

var cellStyle = lipgloss.NewStyle()
var cursorStyle = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("61"))

func (m model) View() string {
	s := "move the cursor with h, j, k, or l\n"
	s += "tap the space to mark the board\n"
	s += "turn: " + string(m.mark) + "\n"
	s += "status: " + m.state + "\n\n"

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			if m.cursor.x == x && m.cursor.y == y {
				s += cursorStyle.Render(string(m.grid[y][x]))
			} else {
				s += cellStyle.Render(string(m.grid[y][x]))
			}
			if x != 2 {
				s += "|"
			}
		}
		s += "\n"
		if y != 2 {
			s += "-----\n"

		}
	}
	s += "\npress r to reset"
	s += "\npress q to quit\n"
	
	return s
}

func initialModel() model {
	grid := make([][]byte, 3)
	
	for i := 0; i < 3; i++ {
		grid[i] = []byte{' ', ' ', ' '}
	}
	return model{
		grid: grid,
		cursor: cursor{0, 0},
		mark: 'x',
		state: "playable",
	}
}


func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
}