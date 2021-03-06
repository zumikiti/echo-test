package main

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gocraft/dbr"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type (
	UserJson struct {
		Name  string `json:"name"`
		Email string `json:"email`
	}

	User struct {
		Id        int    `db:id`
		Name      string `db:name`
		Email     string `db:email`
		Password  string `db:password`
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	response struct {
		Users []User
	}
)

var (
	tableName = "users"
	seq       = 1
	conn, _   = dbr.Open("postgres", "host=pgsql port=5432 user=sail password=password dbname=db sslmode=disable", nil)
	sess      = conn.NewSession(nil)
)

func main() {
	e := echo.New()

	// routes
	e.GET("/users", getUsers)
	e.POST("/users", storeUser)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	e.Logger.Fatal(e.Start(":1323"))
}

func getUsers(c echo.Context) error {
	var u []User
	sess.Select("*").From(tableName).Load(&u)

	res := new(response)
	res.Users = u

	return c.JSON(http.StatusCreated, res)
}

func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")

	var u User
	sess.Select("*").From(tableName).Where("id = ?", id).Load(&u)

	return c.JSON(http.StatusOK, u)
}

func storeUser(c echo.Context) error {
	u := new(User)

	if err := c.Bind(u); err != nil {
		return err
	}

	p := GetMd5("password")

	u.Name = c.FormValue("name")
	u.Email = c.FormValue("email")
	u.Password = p
	u.CreatedAt = time.Now()

	_, err := sess.InsertInto(tableName).
		Columns("name", "email", "password", "created_at").
		Values(u.Name, u.Email, u.Password, u.CreatedAt).
		Exec()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, u)
}

// GetMd5 - get encoded password with md5
func GetMd5(password string) string {
	hash := md5.New()
	defer hash.Reset()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func updateUser(c echo.Context) error {
	id := c.Param("id")

	attrMap := map[string]interface{}{
		"name":       c.FormValue("name"),
		"email":      c.FormValue("email"),
		"updated_at": time.Now(),
	}
	_, err := sess.Update(tableName).SetMap(attrMap).Where("id = ?", id).Exec()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func deleteUser(c echo.Context) error {
	id := c.Param("id")

	_, err := sess.DeleteFrom(tableName).Where("id = ?", id).Exec()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
