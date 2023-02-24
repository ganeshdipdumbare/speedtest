package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ganeshdipdumbare/speedtest/internal/speed"
	"github.com/ganeshdipdumbare/speedtest/internal/speed/fast"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	unexpectedFatalError = "Alas, there's been an error"
	speedCheckCompleted  = "completed"

	colorRed    = "205"
	colorYellow = "220"
)

type speedCheckMsg struct {
	value string
	unit  string
}

func checkSpeed(m model, ch <-chan speed.NetSpeed) tea.Cmd {
	return func() tea.Msg {
		v, ok := <-ch
		if !ok {
			return speedCheckMsg{
				value: speedCheckCompleted,
			}
		}

		if v.Err != nil {
			fmt.Printf("%v: %v", unexpectedFatalError, v.Err)
			os.Exit(1)
		}

		return speedCheckMsg{
			value: v.Value,
			unit:  v.Unit,
		}
	}
}

type model struct {
	downloadSpeed           speedCheckMsg
	uploadSpeed             speedCheckMsg
	isDownloadCheckFinished bool
	isUploadCheckFinished   bool
	downloadspinner         spinner.Model
	uploadSpinner           spinner.Model
	downloadChan            <-chan speed.NetSpeed
	uploadChan              <-chan speed.NetSpeed
}

func getSpinner(colorCode string) spinner.Model {
	ds := spinner.New()
	ds.Spinner = spinner.Pulse
	ds.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(colorCode))
	return ds
}

func InitialModel() model {
	fastResp, err := fast.NewSpeedChecker().GetSpeed()
	if err != nil {
		fmt.Printf("%v: %v", unexpectedFatalError, err)
		os.Exit(1)
	}

	return model{
		downloadSpeed: speedCheckMsg{
			value: "0",
			unit:  "Mbps",
		},
		downloadspinner: getSpinner(colorRed),
		uploadSpinner:   getSpinner(colorYellow),
		downloadChan:    fastResp.DownloadSpeedChannel,
		uploadChan:      fastResp.UploadSpeedChannel,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		checkSpeed(m, m.downloadChan),
		m.downloadspinner.Tick,
		m.uploadSpinner.Tick,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.downloadspinner = getSpinner(colorRed)
			m.uploadSpinner = getSpinner(colorYellow)
			return m, tea.Quit

		}
	case speedCheckMsg:
		message := (speedCheckMsg)(msg)

		if m.isDownloadCheckFinished {
			if message.value == speedCheckCompleted {
				m.isUploadCheckFinished = true
				return m, nil
			}
			m.uploadSpeed = message
			return m, checkSpeed(m, m.uploadChan)
		}

		if message.value == speedCheckCompleted {
			m.isDownloadCheckFinished = true
			return m, checkSpeed(m, m.uploadChan)
		}

		m.downloadSpeed = message
		return m, checkSpeed(m, m.downloadChan)

	case spinner.TickMsg:
		if msg.ID == m.downloadspinner.ID() {
			if !m.isDownloadCheckFinished {
				var cmd tea.Cmd
				m.downloadspinner, cmd = m.downloadspinner.Update(msg)
				return m, cmd
			}
			m.downloadspinner = getSpinner(colorRed)
		}

		if !m.isUploadCheckFinished {
			var cmd tea.Cmd
			m.uploadSpinner, cmd = m.uploadSpinner.Update(msg)
			return m, cmd
		}

		m.uploadSpinner = getSpinner(colorYellow)
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	// The header
	s := ""

	// Iterate over our choices
	s += fmt.Sprintf("%s  Download speed is: %v %v\n", m.downloadspinner.View(), strings.TrimSpace(m.downloadSpeed.value), strings.TrimSpace(m.downloadSpeed.unit))

	if m.isDownloadCheckFinished {
		s += fmt.Sprintf("%s  Upload speed is: %v %v", m.uploadSpinner.View(), strings.TrimSpace(m.uploadSpeed.value), strings.TrimSpace(m.uploadSpeed.unit))
	}
	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
