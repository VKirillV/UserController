package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

type userService struct {
}

type Data struct {
	username string
}

type userType string

const (
	client userType = "Client" //enum
	admin  userType = "Admin"
)

var (
	Userservice userService
)

var data Data
var token string

// .env Setup
func init() {
	e := godotenv.Load()
	if e != nil {
		log.Println(e)
	}
}

func main() {

	r := gin.Default()

	// Start admin server
	r.GET("kript/:adminname/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
		var adminname string = c.Param("adminname")

		// Functions call
		getDatabaseAdmin(adminname)
		TelebotAdminname(adminname)

	})
	// Start client server
	r.GET("/:username/client/:guid", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
		var username string = c.Param("username")
		var guid string = c.Param("guid")

		// Functions call
		getDatabaseUser(username, guid)
		TelebotUsername(username, guid)

	})

	r.Run(":8080")

}

func TelebotAdminname(adminname string) {

	// Telegram bot token
	token = os.Getenv("TOKEN")

	// Telegram bot setup
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	// MySQL database address
	db, err := sqlx.Open("mysql", "root:!KV54691123@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err.Error())
	}

	bot.Debug = true

	// Bot update timing
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {

		// Ignore any non-Message Updates
		if update.Message == nil {
			continue
		}

		//  Get the user from Telegram
		rows, err := db.Queryx("Select username FROM user WHERE username = (?)", update.Message.From.UserName)
		if err != nil {
			log.Panic(err)
		}

		// Scan the resulting variable
		for rows.Next() {
			err := rows.Scan(&data.username)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Checking for the presence of a user in the database
		if data.username == update.Message.From.UserName {

			// Bot message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are registered")
			bot.Send(msg)

		} else if data.username != update.Message.From.UserName {

			// Add user in database
			insert, err := db.Query("INSERT INTO user(username, usertype) VALUES(?, ?)", update.Message.From.UserName, client)
			if err != nil {
				panic(err.Error())
			}

			// Bot message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You subscribed")
			bot.Send(msg)

			defer insert.Close()
		}
	}
}

func TelebotUsername(username, guid string) {

	// Token telegram bot
	token = os.Getenv("TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	// MySQL database address
	db, err := sqlx.Open("mysql", "root:!KV54691123@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err.Error())
	}

	bot.Debug = true

	// Bot update timing
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {

		// Ignore any non-Message Updates
		if update.Message == nil {
			continue
		}

		//  Get the user from Telegram
		rows, err := db.Queryx("Select username FROM user WHERE username = (?)", update.Message.From.UserName)
		if err != nil {
			log.Panic(err)
		}

		// Scan the resulting variable
		for rows.Next() {
			err := rows.Scan(&data.username)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Checking for the presence of a user in the database
		if data.username == update.Message.From.UserName {

			// Bot message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are registered")
			bot.Send(msg)

		} else if data.username != update.Message.From.UserName {

			// Add user in database
			insert, err := db.Query("INSERT INTO user(username, usertype) VALUES(?, ?)", update.Message.From.UserName, client)
			if err != nil {
				panic(err.Error())
			}

			// Message from bot
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You subscribed")
			bot.Send(msg)

			defer insert.Close()
		}

	}

}

func getDatabaseUser(username, guid string) {

	// MySQL database address
	db, err := sqlx.Open("mysql", "root:!KV54691123@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err.Error())
	}

	// Add user in database
	insert, err := db.Query("INSERT INTO user(username, usertype) VALUES(?, ?)", username, client)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	insert, err = db.Query("INSERT INTO client_user(client_guid) VALUES(?)", guid)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	fmt.Println("Succesfully")

}

func getDatabaseAdmin(adminname string) {

	// MySQL database address
	db, err := sqlx.Open("mysql", "root:!KV54691123@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err.Error())
	}

	insert, err := db.Query("INSERT INTO user(username, usertype) VALUES(?, ?)", adminname, admin)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	fmt.Println("Succesfully")

}

// NOT USE
func unSubscribe() (username, guid string) {
	r := gin.Default()
	r.DELETE("/:username/client/:guid", func(c *gin.Context) {
		c.Param("username")
		c.Param("guid")
		var username = c.Param("username")
		var guid = c.Param("guid")
		fmt.Println(username, guid)
	})
	r.Run(":8080")
	return
}

// NOT USE
func unSubscribeAdmin() (username string) {
	r := gin.Default()
	r.DELETE("/:username/admin", func(c *gin.Context) {
		c.Param("username")
		username = c.Param("username")
	})
	r.Run(":8080")
	return

}

// NOT USE
func getAdmins() (clients []string) {
	r := gin.Default()
	r.POST("/admin", func(c *gin.Context) {
	})
	r.Run(":8080")
	return

}
