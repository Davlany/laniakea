package grpcserver

import (
	"context"
	"database/sql"
	"fmt"
	repository "sirius/Repository"
	"sirius/Repository/entities"
	"sirius/proto"
)

const GFP = true

type GrpcServer struct {
	repo repository.Repository
	proto.ServicesServer
}

func (gs *GrpcServer) ConnectRepository(repo repository.Repository) {
	gs.repo = repo
}

func (gs *GrpcServer) GetFriendlyPeers(userData *proto.UserData, stream proto.Services_GetFriendlyPeersServer) error {
	if GFP {
		friends, err := gs.repo.GetFriendlyPeers()
		if err != nil {
			return err
		}
		for _, friend := range friends {
			user := proto.UserIPData{
				IP:      friend.IP,
				Login:   friend.Login,
				OpenKey: friend.OpenKey,
			}
			err := stream.Send(&user)
			if err != nil {
				return err
			}
		}
	} else {
		stream.Send(&proto.UserIPData{})
	}
	return nil

}

func (gs *GrpcServer) Answer(ctx context.Context, userData *proto.UserData) (*proto.StatusCode, error) {
	user := entities.User{
		IP:      userData.GetIp(),
		Login:   userData.GetLogin(),
		OpenKey: userData.GetOpenKey(),
	}
	fmt.Println(user)
	if user.OpenKey == "" && user.Login == "" {
		_, err := gs.repo.GetUserFromWaitList(user)
		if err == sql.ErrNoRows {
			return &proto.StatusCode{
				Status: "401",
			}, err
		}

		if err != nil {
			return &proto.StatusCode{
				Status: "400",
			}, err
		}
		err = gs.repo.DeleteUser(user)
		if err != nil {
			return &proto.StatusCode{
				Status: "400",
			}, err
		}
		return &proto.StatusCode{
			Status: "200",
		}, nil
	}
	fmt.Println("work")
	err := gs.repo.AddToFriendList(user)
	if err != nil {
		return &proto.StatusCode{
			Status: "100",
		}, err
	}
	return &proto.StatusCode{
		Status: "201",
	}, nil
}

func (gs *GrpcServer) AddToWaitUser(ctx context.Context, userData *proto.UserData) (*proto.StatusCode, error) {
	user := entities.User{
		IP:      userData.GetIp(),
		Login:   userData.GetLogin(),
		OpenKey: userData.GetOpenKey(),
	}
	err := gs.repo.AddToWaitToFriendList(user)
	if err != nil {
		return &proto.StatusCode{Status: "400"}, err
	}
	return &proto.StatusCode{Status: "200"}, nil
}

func (gs *GrpcServer) RegisterUser(ctx context.Context, userData *proto.UserData) (*proto.StatusCode, error) {
	user := entities.User{
		IP:      userData.GetIp(),
		Login:   userData.GetLogin(),
		OpenKey: userData.GetOpenKey(),
	}

	err := gs.repo.AddToRequestToFriendList(user)
	if err != nil {
		return &proto.StatusCode{
			Status: "100",
		}, err
	}
	return &proto.StatusCode{
		Status: "200",
	}, err
}
