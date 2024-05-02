package prints

import (
	"fmt"
	"os"
	"os/exec"
	"sirius/Repository/entities"
)

func PrintRequestsToFriend(users []entities.User) {
	fmt.Println("Your requests to friend peers:")
	fmt.Println()
	for _, user := range users {
		fmt.Printf("\033[96m%s\033[97m [%s]\n\n", user.Login, user.IP)
	}
}

func PrintWaitToFriend(users []entities.User) {
	fmt.Println("Your requests to friend peers:")
	fmt.Println()
	for _, user := range users {
		fmt.Printf("\033[96m%s\033[97m [%s]\n\n", user.Login, user.IP)
	}
}

func PrintConfiguration(user entities.User) {
	fmt.Printf("\033[96mLogin\033[97m: %s\n\n", user.Login)
	fmt.Printf("\033[96mOpenKey\033[97m: %s\n\n", user.OpenKey)
	fmt.Printf("\033[96mIP\033[97m: %s\n\n", user.IP)
	fmt.Println("\033[96mShare friendly peers\033[97m: On")
	fmt.Println()
	fmt.Println("\033[96mDatabase\033[97m: Postgres")

}

func PrintMainPage() {
	// ClearConsole()
	PrintLogo()
	fmt.Print("[\033[96m1\033[97m] Friendly peers\t[\033[96m4\033[97m] Connect to peer for register\n\n")
	fmt.Print("[\033[96m2\033[97m] Requests to friend\t[\033[96m5\033[97m] Connect to peer for chating\n\n")
	fmt.Print("[\033[96m3\033[97m] Wait to friend\t[\033[96m6\033[97m] Peer configurations\n\n")
	// fmt.Print("Set the number: ")
}

func PrintFriendlyPeers(users []entities.User) {
	fmt.Println("Your friendly peers:")
	fmt.Println()
	for _, user := range users {
		fmt.Printf("\033[96m%s\033[97m [%s]\n\n", user.Login, user.IP)
	}
}

func PrintLogo() {
	fmt.Println("\033[96m" + "  _                 _       _              ")
	fmt.Println("\033[96m" + " | |               (_)     | |             ")
	fmt.Println("\033[96m" + " | |     __ _ _ __  _  __ _| | _____  __ _ ")
	fmt.Println("\033[96m" + " | |    / _` | '_ \\| |/ _` | |/ / _ \\/ _` |")
	fmt.Println("\033[97m" + " | |___| (_| | | | | | (_| |   <  __/ (_| |")
	fmt.Println("\033[97m" + " |______\\__,_|_| |_|_|\\__,_|_|\\_\\___|\\__,_| v1.1.0")
	fmt.Println()
}

func ClearConsole() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
