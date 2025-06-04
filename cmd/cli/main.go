package main

import (
	"fmt"
	"os"

	"github.com/BrunodsLilly/Summarizer/pkg/core"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	menu state = iota
	inputURL
	processing
	displayResult
)

type model struct {
	state        state
	choices      []string
	cursor       int
	urlInput     textinput.Model
	viewport     viewport.Model
	result       string
	renderedMD   string
	width        int
	height       int
}

func initModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter YouTube video URL"
	ti.CharLimit = 256
	ti.Width = 50

	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	return model{
		state:    menu,
		choices:  []string{"Generate content from YouTube video", "Exit"},
		urlInput: ti,
		viewport: vp,
		width:    80,
		height:   24,
	}
}

type resultMsg struct {
	content string
	err     error
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.state == displayResult {
			m.viewport.Width = msg.Width - 4
			m.viewport.Height = msg.Height - 6
		}
		return m, nil
	
	case resultMsg:
		if msg.err != nil {
			m.result = fmt.Sprintf("Error: %v", msg.err)
		} else {
			m.result = msg.content
		}
		
		renderer, _ := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(m.width-6),
		)
		
		rendered, err := renderer.Render(m.result)
		if err != nil {
			m.renderedMD = m.result
		} else {
			m.renderedMD = rendered
		}
		
		m.viewport.SetContent(m.renderedMD)
		m.state = displayResult
		return m, nil
	}

	switch m.state {
	case menu:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "enter":
				if m.cursor == 0 {
					m.state = inputURL
					m.urlInput.Focus()
					return m, textinput.Blink
				} else {
					return m, tea.Quit
				}
			}
		}

	case inputURL:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "esc":
				m.state = menu
				m.urlInput.Blur()
				m.urlInput.SetValue("")
				return m, nil
			case "enter":
				if m.urlInput.Value() != "" {
					m.state = processing
					return m, fetchSummary(m.urlInput.Value())
				}
			}
		}
		m.urlInput, cmd = m.urlInput.Update(msg)

	case processing:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		}

	case displayResult:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "esc":
				m.state = menu
				m.urlInput.Blur()
				m.urlInput.SetValue("")
				return m, nil
			case "r":
				m.state = menu
				m.urlInput.Blur()
				m.urlInput.SetValue("")
				return m, nil
			}
		}
		m.viewport, cmd = m.viewport.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	switch m.state {
	case menu:
		title := headerStyle.Render("YouTube Video Summarizer")
		
		s := fmt.Sprintf("%s\n\n", title)
		s += "Select an option:\n\n"
		for i, choice := range m.choices {
			cursor := " "
			if m.cursor == i {
				cursor = "▶"
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}
		s += "\nUse ↑/↓ arrows to navigate, Enter to select, q to quit.\n"
		return s

	case inputURL:
		title := headerStyle.Render("Enter URL")
		return fmt.Sprintf(
			"%s\n\nEnter YouTube video URL:\n\n%s\n\n%s",
			title,
			m.urlInput.View(),
			"Press Enter to submit, Esc to go back, q to quit.",
		)

	case processing:
		title := headerStyle.Render("Processing")
		spinner := "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"
		return fmt.Sprintf("%s\n\nProcessing your request... %s", title, string(spinner[0]))

	case displayResult:
		title := headerStyle.Render("Summary Results")
		
		helpStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Render("• Use ↑/↓ arrows to scroll • Press 'r' or 'Esc' to return to menu • Press 'q' to quit")
		
		return fmt.Sprintf("%s\n\n%s\n\n%s", 
			title, 
			m.viewport.View(), 
			helpStyle)
	}

	return ""
}

func fetchSummary(url string) tea.Cmd {
	return func() tea.Msg {
		resp, err := core.SummarizeURL(url)
		return resultMsg{
			content: resp,
			err:     err,
		}
	}
}

func main() {
	p := tea.NewProgram(initModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
