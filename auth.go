package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(c echo.Context) error {
	//check if session already exists (and the user is loggged in)
	sess, _ := session.Get("session", c)
	if auth, ok := sess.Values["authenticated"].(bool); ok && auth {
		log.Print("Redirecting authenticated user to /game")
		return c.Redirect(http.StatusSeeOther, "/game")
	}

	type loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	data := loginData{}

	if err := c.Bind(&data); err != nil {
		return err
	}

	result, err := gorm.G[User](DB).Where("user_name = ?", data.Username).First(c.Request().Context())
	if err != nil {
		return err
	}

	if isPasswordValid(data.Password, result.PasswordHash) {
		sess.Values["authenticated"] = true
		sess.Values["username"] = data.Username
		sess.Save(c.Request(), c.Response())
		log.Printf("Authentication successfull for User %+v!\n", data.Username)
		return c.Redirect(http.StatusSeeOther, "/game")
	}
	return c.String(http.StatusOK, "Not OK")
}

func Register(c echo.Context) error {
	sess, _ := session.Get("session", c)
	if auth, ok := sess.Values["authenticated"].(bool); ok && auth {
		log.Print("Redirecting authenticated user to /game")
		return c.Redirect(http.StatusSeeOther, "/game")
	}

	type RegisterData struct {
		RealName string `json:"realname"`
		Username string `json:"username"`
		Password string `json:"password"`
		Country  string `json:"country"`
	}

	data := &RegisterData{}

	if err := c.Bind(&data); err != nil {
		return err
	}

	err := DisableCountry(data.Country, c)

	if err != nil {
		return err
	}

	data.Password = HashPassword(data.Password)

	DB.Create(&WaitingUser{
		RealName: data.RealName,
		Username: data.Username,
		Password: data.Password,
		Country:  data.Country,
	})

	return c.String(http.StatusOK, "Ok")
}

func DisableCountry(name string, c echo.Context) error {
	country, err := gorm.G[PreviewCountry](DB).Where("name = ?", name).First(c.Request().Context())

	if err != nil {
		return err
	}

	country.Visible = false

	DB.Save(&country)

	return nil
}

func isPasswordValid(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) string {
	hashedByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Could not hash password: %v", err)
	}
	return string(hashedByte)
}
