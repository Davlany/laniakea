package mongorepository

import (
	"context"
	"log"
	repository "sirius/Repository"
	"sirius/Repository/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	repository.Repository
	db *mongo.Database
}

func (mr MongoRepository) AddToRequestToFriendList(user entities.User) error {
	collection := mr.db.Collection("RequestToFriend")
	_, err := collection.InsertOne(context.TODO(), user)
	return err
}

func (mr MongoRepository) AddToWaitToFriendList() {}

func (mr MongoRepository) AddToFriendList(user entities.User) error {
	collection := mr.db.Collection("FriendList")
	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
}

func (mr MongoRepository) DeleteFromRequestToFriendList() {}

func (mr MongoRepository) DeleteFromFriendList() {}

func (mr MongoRepository) DeleteFromWaitToFriendList(user entities.User) error {
	collection := mr.db.Collection("WaitToFriendList")
	filter := bson.D{{"ip", user.IP}}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (mr MongoRepository) GetUserFromWaitList(user entities.User) (entities.User, error) {
	collection := mr.db.Collection("WaitToFriendList")
	var answerUser entities.User
	filter := bson.D{{"ip", user.IP}}
	err := collection.FindOne(context.TODO(), filter).Decode(&answerUser)
	if err != nil {
		return entities.User{}, err
	}
	return answerUser, nil
}

func (mr MongoRepository) GetFriendlyPeers() ([]entities.User, error) {
	collection := mr.db.Collection("FriendList")
	var friends []entities.User
	filter := bson.M{}
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var friend entities.User
		err := cur.Decode(&friend)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return friends, nil
}

func NewMongoRepository(uri string) *MongoRepository {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connected to database succesful")
	database := client.Database("Sirius")
	return &MongoRepository{
		db: database,
	}
}
