package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type App struct {
	DB *sql.DB
}

func main() {

	//connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//create the table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT, email TEXT)")

	if err != nil {
		log.Fatal(err)
	}

	app := App{
		DB: db,
	}

	mux := gin.Default()

	mux.GET("/users", app.getUsers)
	mux.POST("/users", app.createUser)
	mux.DELETE("/users/:id", app.deleteUser)

	//start server
	log.Fatal(http.ListenAndServe(":8000", mux))

}

func (a *App) getUsers(c *gin.Context) {

	rows, err := a.DB.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, users)

}

func (a *App) createUser(c *gin.Context) {
	var u User
	err := c.ShouldBindJSON(&u)
	if err != nil {
		fmt.Println(err)
	}

	err = a.DB.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", u.Name, u.Email).Scan(&u.ID)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusCreated, u)

}

func (a *App) deleteUser(c *gin.Context) {
	id := c.Params

	var u User
	err := a.DB.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	} else {
		_, err := a.DB.Exec("DELETE FROM users WHERE id = $1", id)
		if err != nil {
			//todo : fix error handling
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "error deleting",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "user deleted",
		})

	}
}
