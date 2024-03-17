package entities

type User struct {
	Id      int    `db:"id"`
	Login   string `bson:"login" db:"login"`
	IP      string `bson:"IP" db:"ip"`
	OpenKey string `bson:"open_key" db:"open_key"`
	Ship    string `db:"ship"`
}

type Table struct {
	TableName string `db:"table_name"`
}