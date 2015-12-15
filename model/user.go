package model

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/naoina/genmai"
)

// User is user DB table
type User struct {
	ID         uint64 `db:"pk"`
	ScreenName string `db:"unique"`
	CryptoPass string
	Name       string
	Age        int
	Intro      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// CreateUser is user create
func CreateUser(screenName, cryptoPass, name, intro string, age int) User {
	db, err := genmai.New(&genmai.MySQLDialect{}, "root@/ginapp?parseTime=true")
	if err != nil {
		log.Errorln("db conn err:", err)
	}

	user := User{
		ScreenName: screenName,
		CryptoPass: cryptoPass,
		Name:       name,
		Age:        age,
		Intro:      intro,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	n, err := db.Insert(&user)
	if err != nil {
		log.Errorln("db insert err:", err)
	}

	log.Debugf("inserted rows: %d\n", n)

	return user
}

// GetUser is get user from screen name
func GetUser(screenName string) User {
	db, err := genmai.New(&genmai.MySQLDialect{}, "root@/ginapp?parseTime=true")
	if err != nil {
		log.Errorln("db conn err:", err)
	}

	db.SetLogOutput(os.Stdout)
	db.SetLogFormat("format string")

	var users []User
	if err := db.Select(&users, db.Where("screen_name", "=", screenName)); err != nil {
		log.Error("select error: ", err)
		return User{}
	}
	log.Debug("screen_name: ", screenName, ", users: ", users)

	// screen name is unique
	return users[0]
}
