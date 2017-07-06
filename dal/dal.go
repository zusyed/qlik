package dal

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //this is the postgres dialect for gorm"github.com/jinzhu/gorm"
	"github.com/zusyed/qlik/models"
)

//The constants for the database connection should be in a config file
const (
	//Host is the host of the Postgres database
	Host = "localhost"

	//User is the username of the Postgres database
	User = "postgres"

	//Password is the password of the Postgres database
	Password = "postgres"

	//DbName is the database name for the messages table
	DbName = "messages"

	//NotFound returned from database when a record is not found
	NotFound = "not found"
)

type (
	//MessageDB provides functionality to interact with the messages database
	MessageDB interface {
		GetMessages() ([]models.Message, error)
		GetMessage(id int) (models.Message, error)
		AddMessage(message models.Message) (models.Message, error)
		DeleteMessage(id int) error
	}

	//PostgresDB provides functionality to interact with the postgres messages database
	PostgresDB struct {
		db *gorm.DB
	}
)

//NewMessageDB connects to the message database and creates the necessary tables
func NewMessageDB() (MessageDB, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}

	db.createTables()

	return db, nil
}

//Connect creates a new connection to the database
func Connect() (*PostgresDB, error) {

	conn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", Host, User, DbName, Password)
	gDB, err := gorm.Open("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to %s: %+v", DbName, err)
	}

	err = gDB.DB().Ping()
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to %s: %+v", DbName, err)
	}

	gDB.DB().SetMaxIdleConns(20)
	gDB.DB().SetMaxOpenConns(200)

	pdb := &PostgresDB{
		db: gDB,
	}

	return pdb, err
}

//createTables sets up the postgres schema
func (pdb *PostgresDB) createTables() {
	pdb.db.AutoMigrate(&models.Message{})
}

//GetMessages returns all the messages from the database
func (pdb *PostgresDB) GetMessages() ([]models.Message, error) {
	messages := []models.Message{}
	err := pdb.db.Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("Unable to get messages from db: %s", err)
	}

	return messages, nil
}

//GetMessage returns the message with the specified id from the database
func (pdb *PostgresDB) GetMessage(id int) (models.Message, error) {
	message := models.Message{}
	result := pdb.db.First(&message, id)
	err := result.Error
	if err != nil {
		if result.RecordNotFound() {
			return message, fmt.Errorf(NotFound)
		}

		return message, fmt.Errorf("Unable to get message with the specified id %d: %s", id, err)
	}

	return message, nil
}

//AddMessage adds the message to the database
func (pdb *PostgresDB) AddMessage(message models.Message) (models.Message, error) {
	err := pdb.db.Create(&message).Error
	if err != nil {
		return message, fmt.Errorf("Unable to add message to the db:  %s", err)
	}

	return message, nil
}

//DeleteMessage deletes the message with the specified id from the database
func (pdb *PostgresDB) DeleteMessage(id int) error {
	message := models.Message{
		ID: id,
	}
	result := pdb.db.Delete(&message)
	err := result.Error
	if err != nil {
		return fmt.Errorf("Unable to delete message with the specified id %d: %s", id, err)
	}

	if result.RowsAffected < 1 {
		return fmt.Errorf(NotFound)
	}

	return nil
}
