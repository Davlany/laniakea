package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sirius/Repository/entities"
	"sirius/client/prints"
	"sirius/proto"
	"sirius/server"

	"google.golang.org/grpc"
)

func main() {
	prints.PrintMainPage()
	login := os.Getenv("login")
	ip := os.Getenv("ip")
	port := os.Getenv("port")
	dbUser := os.Getenv("dbuser")
	server, err := server.NewServer(login, ip, dbUser, "123456", "5432", "disable")
	if err != nil {
		log.Fatalln(err)
	}
	go server.ServerRun(port)
	var flag int
	for {
		var input string
		fmt.Scanln(&input)
		if flag != 0 {
			if input == "\\b" {
				prints.PrintMainPage()
				flag = 0
			}
		}
		if input == "1" {
			prints.ClearConsole()
			prints.PrintLogo()
			users, err := server.Repo.GetFriendlyPeers()
			if err != nil {
				log.Fatalln(err)
			}
			prints.PrintFriendlyPeers(users)
			fmt.Print("Input command: ")
			flag = 1
		} else if input == "2" {
			flag = 2
			prints.ClearConsole()
			prints.PrintLogo()
			users, err := server.Repo.GetRequestsToFriend()
			if err != nil {
				log.Fatalln(err)
			}
			prints.PrintRequestsToFriend(users)
		} else if input == "3" {
			flag = 3
			prints.ClearConsole()
			prints.PrintLogo()
			users, err := server.Repo.GetWaitToFriend()
			if err != nil {
				log.Fatalln(err)
			}
			prints.PrintWaitToFriend(users)
		} else if input == "4" {
			flag = 4
			prints.ClearConsole()
			prints.PrintLogo()
			fmt.Println("Input ip and port to request")
			fmt.Print("IP: ")
			var ip string
			fmt.Scanln(&ip)
			conn, err := grpc.Dial(ip, grpc.WithInsecure())
			if err != nil {
				log.Fatalln(err)
			}
			c := proto.NewServicesClient(conn)
			user, err := server.Repo.GetOwnerUser()
			if err != nil {
				log.Fatalln(err)
			}
			res, err := c.RegisterUser(context.Background(), &proto.UserData{
				Login:   user.Login,
				OpenKey: user.OpenKey,
			})
			if err != nil {
				log.Fatal(err)
			}
			conn.Close()
			fmt.Print("\nEnter login for this peer: ")
			var name string
			fmt.Scanln(&name)
			err = server.Repo.AddToWaitToFriendList(entities.User{Login: name, IP: ip})
			if err != nil {
				log.Fatalln(err)
			}

			if res.GetStatus() == "200" {
				fmt.Println("Request succesfully!")
			}

		} else if input == "6" {
			flag = 6
			prints.ClearConsole()
			prints.PrintLogo()
			user, err := server.Repo.GetOwnerUser()
			if err != nil {
				log.Fatalln(err)
			}
			prints.PrintConfiguration(user)

		}
	}
}
