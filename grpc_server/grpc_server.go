package grpcserver

import (
	"context"
	repository "sirius/Repository"
	"sirius/Repository/entities"
	"sirius/proto"

	"google.golang.org/grpc/peer"
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
	p, _ := peer.FromContext(ctx)
	user := entities.User{
		IP:      p.Addr.String(),
		Login:   userData.GetLogin(),
		OpenKey: userData.GetLogin(),
	}
	if user.OpenKey == "" && user.Login == "" {
		waitUser, err := gs.repo.GetUserFromWaitList(user)
		if err != nil {
			return &proto.StatusCode{
				Status: "100",
			}, err
		}
		if waitUser != user {
			return &proto.StatusCode{
				Status: "401",
			}, nil
		}
		err = gs.repo.DeleteFromWaitToFriendList(user)
		if err != nil {
			return &proto.StatusCode{
				Status: "100",
			}, err
		}
		return &proto.StatusCode{
			Status: "200",
		}, err
	}
	waitUser, err := gs.repo.GetUserFromWaitList(user)
	if err != nil {
		return &proto.StatusCode{
			Status: "100",
		}, err
	}
	if waitUser != user {
		return &proto.StatusCode{
			Status: "401",
		}, nil
	}
	err = gs.repo.DeleteFromWaitToFriendList(user)
	if err != nil {
		return &proto.StatusCode{
			Status: "100",
		}, err
	}
	err = gs.repo.AddToFriendList(user)
	if err != nil {
		return &proto.StatusCode{
			Status: "100",
		}, err
	}
	return &proto.StatusCode{
		Status: "200",
	}, nil
}

func (gs *GrpcServer) RegisterUser(ctx context.Context, userData *proto.UserData) (*proto.StatusCode, error) {
	p, _ := peer.FromContext(ctx)
	user := entities.User{
		IP:      p.Addr.String(),
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
