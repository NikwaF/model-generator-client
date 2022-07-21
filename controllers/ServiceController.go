package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

var base_url = "http://localhost:8082"
var target_url = "http://localhost:8081"
var store = sessions.NewCookieStore([]byte("SD4-GO"))

func LoginProcess(c echo.Context) error {
	// username := c.FormValue("Username")
	// password := c.FormValue("Password")

	// postBody, _ := json.Marshal(map[string]string{
	// 	"Username": username,
	// 	"Password": password,
	// })

	// request, err := http.Post(target_url+"/SignIn", "application/json", bytes.NewBuffer(postBody))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// defer request.Body.Close()

	// body, err := ioutil.ReadAll(request.Body)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	checkUnautorized := false

	if checkUnautorized == true {
		c.Redirect(http.StatusFound, base_url+"/login")
	} else {
		session, _ := store.Get(c.Request(), "SessionLogin")
		session.Options.MaxAge = 21600
		session.Values["Status"] = "Login"
		session.Save(c.Request(), c.Response())

		c.Redirect(http.StatusFound, base_url+"/")
	}

	return nil
}

func CheckLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, _ := store.Get(c.Request(), "SessionLogin")
		if len(session.Values) == 0 {
			fmt.Println("Empty Session")
		}

		statusLogin := session.Values["Status"]
		if statusLogin != "Login" {
			fmt.Println("Unauthorized")
			c.Redirect(http.StatusFound, base_url+"/login")
		} else {
			next(c)
		}

		return nil
	}
}

func LogoutProcess(c echo.Context) error {
	session, _ := store.Get(c.Request(), "SessionLogin")
	session.Options.MaxAge = -1
	session.Save(c.Request(), c.Response())

	sessionMenu, _ := store.Get(c.Request(), "SessionMenu")
	sessionMenu.Options.MaxAge = -1
	sessionMenu.Save(c.Request(), c.Response())

	c.Redirect(http.StatusFound, base_url+"/login")
	return nil
}
