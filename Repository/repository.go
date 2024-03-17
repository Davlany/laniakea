package repository

import "sirius/Repository/entities"

type Repository interface {
	AddToRequestToFriendList(entities.User) error
	AddToFriendList(entities.User) error
	DeleteFromRequestToFriendList()
	DeleteFromWaitToFriendList(entities.User) error
	DeleteFromFriendList()
	GetUserFromWaitList(entities.User) (entities.User, error)
	GetFriendlyPeers() ([]entities.User, error)
}
