package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type User struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Age       int        `json:"age"`
	Address   string     `json:"address"`
	CreatedAt time.Time  `json:"createdat"`
	UpdatedAt *time.Time `json:"updatedat,omitempty"`
}

func main() {
	r := echo.New()

	db := NewPostgreStore("postgres", "admin", "crud", "localhost", "disable")

	r.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "home")
	})

	//endpoint to get all users
	r.GET("/users", func(c echo.Context) error {

		users, err := db.GetUsers(context.Background())
		if err != nil {
			log.Fatal(err.Error())
			return err
		}

		return c.JSON(http.StatusOK, users)
	})

	//endoint to get specific user with id as param
	r.GET("/user/:id", func(c echo.Context) error {
		// user := new(User)
		id := c.Param("id")

		user, err := db.GetUserByID(context.Background(), id)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, user)
	})

	//endpoint to create user
	r.POST("/user", func(c echo.Context) error {
		param := new(DBParam)

		if err := c.Bind(param); err != nil {
			log.Fatal(err.Error())
			return err
		}

		err := db.CreateUser(context.Background(), param)
		if err != nil {
			log.Fatal(err.Error())
			return err
		}

		return c.JSON(http.StatusOK, "success")
	})

	//endpoint to update user wiht id and json as param
	r.POST("/update/:id", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.String(http.StatusInternalServerError, "must have id for param")
		}
		param := new(DBParam)

		if err := c.Bind(param); err != nil {
			log.Fatal(err.Error())
			return err
		}

		if err := db.UpdateUser(context.Background(), id, param); err != nil {
			log.Fatal(err.Error())
			return err
		}

		return c.JSON(http.StatusOK, "success update")
	})

	//endpoint to delete user
	r.DELETE("/delete/:id", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.String(http.StatusOK, "param must not nil")
		}

		if err := db.DeleteUser(context.Background(), id); err != nil {
			log.Fatal(err.Error())
			return err
		}
		return c.JSON(http.StatusOK, "delete success")
	})

	r.Start(":9000")
}
