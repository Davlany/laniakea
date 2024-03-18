package main

import (
	"log"
	"net"
	"os"
	"sirius/Repository/postgresrepo"
	grpcserver "sirius/grpc_server"
	"sirius/proto"

	"google.golang.org/grpc"
)

const URI = "mongodb+srv://euler:xbLK6uPRlNdN0JY3@sirius1.bd5egub.mongodb.net/?retryWrites=true&w=majority"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Listenning @ : " + port)

	repo, err := postgresrepo.NewPostgresDriver("postgres", "123456", "5432", "disable")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()

	ss := grpcserver.GrpcServer{}
	ss.ConnectRepository(repo)

	proto.RegisterServicesServer(grpcServer, &ss)
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalln(err)
	}
}
