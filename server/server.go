package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	repository "sirius/Repository"
	"sirius/Repository/entities"
	"sirius/Repository/postgresrepo"
	grpcserver "sirius/grpc_server"
	"sirius/proto"

	"google.golang.org/grpc"
)

type Server struct {
	Repo repository.Repository
}

func (s *Server) ServerRun(port string) {
	if port == "" {
		port = "8000"
	}
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()

	ss := grpcserver.GrpcServer{}
	ss.ConnectRepository(s.Repo)

	proto.RegisterServicesServer(grpcServer, &ss)
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalln(err)
	}
}

func NewServer(userLogin, userIP, user, password, dbPort, sslMode string) (*Server, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Error to generated RSA keys:", err)
		return nil, err
	}
	openKey := &privateKey.PublicKey
	pubKeyBytes := x509.MarshalPKCS1PublicKey(openKey)
	pubKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}
	pemPublicKey := pem.EncodeToMemory(pubKeyBlock)

	// Кодирование PEM-блока в Base64
	publicKeyString := base64.StdEncoding.EncodeToString(pemPublicKey)
	userData := entities.User{
		Login:   userLogin,
		IP:      userIP,
		OpenKey: publicKeyString,
	}
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	pemPrivateKey := pem.EncodeToMemory(privateKeyBlock)

	// Кодирование PEM-блока в Base64
	privateKeyString := base64.StdEncoding.EncodeToString(pemPrivateKey)

	repo, err := postgresrepo.NewPostgresDriver(userData, privateKeyString, user, password, dbPort, sslMode)
	if err != nil {
		return nil, err
	}
	return &Server{
		Repo: repo,
	}, nil
}
