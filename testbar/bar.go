package main

import (
	"fmt"
	"log"
	"os"
	"sirius/Repository/entities"
	"sirius/server"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices       []string
	cursor        int
	selected      map[int]struct{}
	user          entities.User
	serv          *server.Server
	flag          int
	friendlyPeers []entities.User
}

func initialModel() model {
	server, err := server.NewServer("Alice", "127.0.0.1:8000", "postgres", "123456", "5432", "disable")
	if err != nil {
		log.Fatalln(err)
	}
	go server.ServerRun("8000")
	return model{
		choices:  []string{"\033[97m" + "[\033[96m1\033[97m] Friendly peers\n", "[\033[96m2\033[97m] Requests to friend\n", "[\033[96m3\033[97m] Wait to friend\n", "[\033[96m4\033[97m] Connect to peer for register\n", "[\033[96m5\033[97m] Connect to peer for chating\n", "[\033[96m6\033[97m] Peer configurations\n"},
		cursor:   0,
		selected: make(map[int]struct{}),
		serv:     server,
		flag:     -1,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.cursor == 0 {
				m.cursor = len(m.choices)
			}
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor == len(m.choices)-1 {
				m.cursor = 0
			} else {
				m.cursor++
			}
		case "b":
			m.flag = -1
		case "enter", " ":
			if m.cursor == 0 {
				users, err := m.serv.Repo.GetFriendlyPeers()
				if err != nil {
					log.Fatalln(err)
				}
				m.friendlyPeers = users
				m.flag = 1
			} else if m.cursor == 5 {
				user, err := m.serv.Repo.GetOwnerUser()
				if err != nil {
					log.Fatalln(err)
				}
				m.user = user
				m.flag = 5
			}

		}
	}
	return m, nil
}

func (m model) View() string {
	s := "\033[96m" + "  _                 _       _              \n" + "\033[96m" + " | |               (_)     | |             \n" + "\033[96m" + " | |     __ _ _ __  _  __ _| | _____  __ _ \n" + "\033[96m" + " | |    / _` | '_ \\| |/ _` | |/ / _ \\/ _` |\n" + "\033[97m" + " | |___| (_| | | | | | (_| |   <  __/ (_| |\n" + "\033[97m" + " |______\\__,_|_| |_|_|\\__,_|_|\\_\\___|\\__,_| v1.1.0\n\n" + "\033[97m"

	switch m.flag {
	case -1:
		for i, choice := range m.choices {

			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}
		return s
	case 1:
		s += "Your friendly peers:\n\n"
		for _, user := range m.friendlyPeers {
			s += fmt.Sprintf("\033[96m%s\033[97m [%s]\n\n", user.Login, user.IP)
		}
		return s
	case 5:
		s += fmt.Sprintf("\033[96mLogin\033[97m: %s\n\n\033[96mOpenKey\033[97m: %s\n\n\033[96mIP\033[97m: %s\n\n\033[96mShare friendly peers\033[97m: On\n\n\033[96mDatabase\033[97m: Postgres\n\n", m.user.Login, m.user.OpenKey, m.user.IP)
		return s
	}
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
