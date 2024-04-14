package main

import (
	"database/sql"
	"fmt"      // formatting and printing values to the console.
	"log"      // logging messages to the console.
	"net/http" // Used for build HTTP servers and clients.

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

// Port we listen on.
const portNum string = ":8080"
const connStr = "user=postgres password=mysecretpassword dbname=postgres sslmode=disable"

var db *sql.DB

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Id   int
}

// Handler functions.
func Home(c echo.Context) error {
	// fmt.Fprintf(w, "Homepage")
	return c.String(http.StatusAccepted, "hello home")
}

func CreateUser(c echo.Context) error {

	var user User
	err := c.Bind(&user)

	if err != nil {
		return err
	}

	tx, err := db.Begin()

	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO public.user (name, age) VALUES($1,$2) RETURNING id;`)

	if err != nil {
		tx.Rollback()
		return err
	}

	defer stmt.Close()

	row := stmt.QueryRow(user.Name, user.Age)
	var id string

	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return fmt.Errorf("user error : %v", err)
	}

	tx.Commit()
	return c.JSON(http.StatusCreated, id)

}

func GetUser(c echo.Context) error {
	var users []User
	rows, err := db.Query("SELECT * FROM public.user ")

	if err != nil {
		return fmt.Errorf("user error : %v", err)
	}

	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var user User

		if err := rows.Scan(&user.Name, &user.Age, &user.Id); err != nil {
			msg := fmt.Errorf("user error : %v", err)
			fmt.Println("error theere", msg)
			// return fmt.Errorf("user error : %v", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("albumsByArtist : %v", err)
	}

	return c.JSON(http.StatusAccepted, users)

}

func DisplayHello(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		fmt.Println("hello console")
		return next(c)

	})
}

func HandleBasicAuth(username string, password string, c echo.Context) (bool, error) {
	if username == "joe" && password == "secret" {
		return true, nil
	}
	return false, nil
}

func main() {
	log.Println("Starting our simple http server.")

	var err error
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()

	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected! to db")

	e := echo.New()
	e.GET("/", Home)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(DisplayHello)

	g := e.Group("/api")

	g.Use(middleware.BasicAuth(HandleBasicAuth))

	g.GET("/users", GetUser)
	g.POST("/users", CreateUser)

	e.Logger.Fatal(e.Start(portNum))
	fmt.Println("To close connection CTRL+C :-)")

}
