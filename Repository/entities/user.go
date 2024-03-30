package entities

type User struct {
	Id         int    `db:"id"`
	Login      string `bson:"login" db:"login"`
	IP         string `bson:"IP" db:"ip"`
	OpenKey    string `bson:"open_key" db:"open_key"`
	PrivateKey string `bson:"open_key" db:"private_key"`
	Ship       string `db:"ship"`
	Gfp        string `db:"gfp"`
}

type Table struct {
	TableName string `db:"table_name"`
}

type Message struct {
	Id        int    `db:"id"`
	IsOwner   string `db:"is_owner"`
	From      string `db:"from"`
	Data      string `db:"data"`
	TimeStamp string `db:"timestamp"`
}
