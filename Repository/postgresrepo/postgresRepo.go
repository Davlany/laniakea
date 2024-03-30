package postgresrepo

import (
	"database/sql"
	"fmt"
	"log"
	repository "sirius/Repository"
	"sirius/Repository/entities"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	repository.Repository
	conn *sqlx.DB
}

func (ps *PostgresRepo) AddToRequestToFriendList(user entities.User) error {
	isInFriendListQuery := fmt.Sprintf("SELECT * FROM users WHERE open_key = '%s'", user.OpenKey)
	err := ps.conn.Get(&entities.User{}, isInFriendListQuery)
	if err != nil {
		if err == sql.ErrNoRows {
			query := fmt.Sprintf("INSERT INTO users(login, open_key, ip, ship) VALUES('%s','%s','%s', 'rtf')", user.Login, user.OpenKey, user.IP)
			_, err := ps.conn.Query(query)
			if err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}
	log.Println("User already in friendly list")
	return nil
}

func (ps *PostgresRepo) AddToWaitToFriendList(user entities.User) error {
	query := fmt.Sprintf("INSERT INTO users(login, open_key, ip, ship) VALUES('%s','%s','%s','wtf')", user.Login, user.OpenKey, user.IP)
	_, err := ps.conn.Query(query)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostgresRepo) AddToFriendList(user entities.User) error {
	query := fmt.Sprintf("UPDATE users SET ship = 'f' WHERE ip = '%s'", user.IP)
	_, err := ps.conn.Query(query)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostgresRepo) GetPrivateKey() (string, error) {
	var res string
	query := "SELECT private_key FROM owners"
	err := ps.conn.Get(&res, query)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (ps *PostgresRepo) DeleteUser(user entities.User) error {
	query := fmt.Sprintf("DELETE FROM users WHERE open_key = '%s';", user.OpenKey)
	_, err := ps.conn.Query(query)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostgresRepo) GetUserFromWaitList(user entities.User) (entities.User, error) {
	var res entities.User

	getQuery := fmt.Sprintf("SELECT * FROM users WHERE open_key = '%s' AND ship = 'wtf'", user.OpenKey)
	err := ps.conn.Get(&res, getQuery)
	if err != nil {
		return entities.User{}, err
	}
	return res, nil
}

func (ps *PostgresRepo) GetFriendlyPeers() ([]entities.User, error) {
	var res []entities.User
	getQuery := "SELECT * FROM users WHERE ship = 'f'"
	err := ps.conn.Select(&res, getQuery)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (ps *PostgresRepo) GetRequestsToFriend() ([]entities.User, error) {
	var res []entities.User
	getQuery := "SELECT * FROM users WHERE ship = 'rtf'"
	err := ps.conn.Select(&res, getQuery)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (ps *PostgresRepo) GetWaitToFriend() ([]entities.User, error) {
	var res []entities.User
	getQuery := "SELECT * FROM users WHERE ship = 'wtf'"
	err := ps.conn.Select(&res, getQuery)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (ps *PostgresRepo) GetAllMessages(openKey string) ([]entities.Message, error) {
	var msgs []entities.Message
	query := fmt.Sprintf("SELECT * FROM WHERE from = '%s' OR to = '%s'", openKey, openKey)
	err := ps.conn.Select(&msgs, query)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (ps *PostgresRepo) AddMessage(msg entities.Message) error {
	query := fmt.Sprintf("INSERT INTO messages(is_owner, from, data, timestamp) VALUES (%s,'%s','%s','%s')", msg.IsOwner, msg.From, msg.Data, msg.TimeStamp)
	_, err := ps.conn.Query(query)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostgresRepo) GetOwnerUser() (entities.User, error) {
	var res entities.User
	getQuery := "SELECT * FROM owners"
	err := ps.conn.Get(&res, getQuery)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (ps *PostgresRepo) InitializeDatabase() error {

	query1 := "SELECT table_name FROM information_schema.tables WHERE table_name = 'users' AND table_schema = 'public'"
	err := ps.conn.Get(&entities.Table{}, query1)
	if err == sql.ErrNoRows {
		createQuery := "CREATE TABLE users(id SERIAL, login TEXT, open_key TEXT, ip TEXT, ship TEXT);"
		createOwnerQuery := "CREATE TABLE owners(id, login TEXT, open_key TEXT, ip TEXT, private_key text, gfp TEXT);"
		createMessageTable := "CREATE TABLE messages(id SERIAL PRIMARY KEY, from TEXT, to TEXT, data TEXT, timestamp string);"
		_, err := ps.conn.Query(createQuery)
		if err != nil {
			log.Println("Initialize error")
			return err
		}
		_, err = ps.conn.Query(createOwnerQuery)
		if err != nil {
			log.Println("Initialize error")
			return err
		}
		_, err = ps.conn.Query(createMessageTable)
		if err != nil {
			log.Println("Initialize error")
			return err
		}

	}
	return nil
}

func (ps *PostgresRepo) InitUser(user entities.User, privateKey string) error {
	query1 := "SELECT * FROM owners;"
	err := ps.conn.Get(&entities.User{}, query1)
	if err == sql.ErrNoRows {
		createQuery := fmt.Sprintf("INSERT INTO owners(login, open_key, ip, private_key, gfp) VALUES ('%s','%s','%s','%s', true)", user.Login, user.OpenKey, user.IP, privateKey)
		_, err = ps.conn.Query(createQuery)
		if err != nil {
			log.Println("Initialize user error", err)
			return nil
		}
	}
	return err
}

func NewPostgresDriver(userData entities.User, privateKey string, user, password, port, sslMode string) (*PostgresRepo, error) {
	var repo PostgresRepo
	conn, err := sqlx.Connect("postgres", fmt.Sprintf("user = %s password = %s dbname = sirius2 sslmode = %s port = %s", user, password, sslMode, port))
	if err != nil {
		return nil, err
	}
	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	repo.conn = conn
	err = repo.InitializeDatabase()

	if err != nil {
		return &PostgresRepo{}, err
	}

	err = repo.InitUser(userData, privateKey)
	if err != nil {
		return nil, err
	}
	//log.Println("Postgres database connected succesfully!")
	return &PostgresRepo{
		conn: conn,
	}, nil
}
