package repository

import "sirius/Repository/entities"

type Repository interface {
	AddToRequestToFriendList(entities.User) error
	AddToFriendList(entities.User) error
	AddToWaitToFriendList(entities.User) error
	GetUserFromWaitList(entities.User) (entities.User, error)
	GetFriendlyPeers() ([]entities.User, error)
	GetRequestsToFriend() ([]entities.User, error)
	GetWaitToFriend() ([]entities.User, error)
	GetOwnerUser() (entities.User, error)
	DeleteUser(entities.User) error
}
