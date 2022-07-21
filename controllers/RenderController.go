package controllers

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	"github.com/labstack/echo/v4"
)

//DataParams ...
type DataParams map[string]interface{}

// TemplateRenderer ...
type TemplateRenderer struct {
	templates *template.Template
}

// Render TemplateRenderer
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// LoginPages ...
func LoginPages(c echo.Context) error {
	data := DataParams{
		"title":   "Helper Easy API",
		"message": "Hello World!",
	}
	return c.Render(http.StatusOK, "login.html", data)
}

func Dashboard(c echo.Context) error {
	session, _ := store.Get(c.Request(), "SessionMenu")
	session.Values["Menu"] = "Dashboard"
	session.Save(c.Request(), c.Response())

	sessionMenu := session.Values["Menu"]
	data := DataParams{
		"title":        "Api Generator",
		"randomNumber": rand.Intn(200),
		"navSession":   sessionMenu,
	}
	return c.Render(http.StatusOK, "dashboard", data)
}

// GenV1 Index pages
func GenV1(c echo.Context) error {
	session, _ := store.Get(c.Request(), "SessionMenu")
	session.Values["Menu"] = "C# Get SQL Statement Generator"
	session.Save(c.Request(), c.Response())

	sessionMenu := session.Values["Menu"]
	data := DataParams{
		"title":        "SC Generator",
		"randomNumber": rand.Intn(200),
		"navSession":   sessionMenu,
	}
	return c.Render(http.StatusOK, "index", data)
}

func GenV2(c echo.Context) error {
	session, _ := store.Get(c.Request(), "SessionMenu")
	session.Values["Menu"] = "C# SQL Statement CSV Generator"
	session.Save(c.Request(), c.Response())

	sessionMenu := session.Values["Menu"]
	data := DataParams{
		"title":        "SC Generator",
		"randomNumber": rand.Intn(200),
		"navSession":   sessionMenu,
	}
	return c.Render(http.StatusOK, "csvgen", data)
}

func ListApi(c echo.Context) error {
	resp, err := http.Get(target_url + "/helper_api")

	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	var apiData []map[string]interface{}

	json.Unmarshal(body, &apiData)

	//return c.JSON(http.StatusOK, apiData)
	//Send params to View
	session, _ := store.Get(c.Request(), "SessionMenu")
	session.Values["Menu"] = "List Endpoint API"
	session.Save(c.Request(), c.Response())

	sessionMenu := session.Values["Menu"]

	data := DataParams{
		"title":        "List Endpoint API",
		"randomNumber": rand.Intn(200),
		"target_url":   base_url,
		"apiData":      apiData,
		"navSession":   sessionMenu,
	}

	return c.Render(http.StatusOK, "list", data)
}

func ListEditApi(c echo.Context) error {
	id := c.Param("id")

	detailList, err := ListEdit(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "oops ada masalah")
	}

	session, _ := store.Get(c.Request(), "SessionMenu")
	session.Values["Menu"] = "List Endpoint API"
	session.Save(c.Request(), c.Response())

	sessionMenu := session.Values["Menu"]
	data := DataParams{
		"title":          "List Endpoint API",
		"randomNumber":   rand.Intn(200),
		"navSession":     sessionMenu,
		"idEndpoint":     id,
		"detailEndpoint": detailList,
	}
	return c.Render(http.StatusOK, "list_edit", data)
}

func ExecApi(c echo.Context) error {
	session, _ := store.Get(c.Request(), "SessionMenu")
	session.Values["Menu"] = "Execute API"
	session.Save(c.Request(), c.Response())

	sessionMenu := session.Values["Menu"]

	data := DataParams{
		"title":        "Execute API",
		"randomNumber": rand.Intn(200),
		"navSession":   sessionMenu,
	}

	return c.Render(http.StatusOK, "execute", data)
}

func ExecApi2(c echo.Context) error {
	id := c.Param("id")
	info, err := GetInfoDbcode(id)

	if err != nil {
		return c.String(http.StatusInternalServerError, "oops ada masalah")
	}

	html := MilihHtml(info.Data.Type_endpoint, info.Data.JmlParam, info.Data.JmlData)

	session, _ := store.Get(c.Request(), "SessionMenu")
	session.Values["Menu"] = "Execute API"
	session.Save(c.Request(), c.Response())

	sessionMenu := session.Values["Menu"]

	data := DataParams{
		"title":        "Execute API",
		"randomNumber": rand.Intn(200),
		"info":         info,
		"html":         template.HTML(html),
		"query":        template.HTML(htmlQuery(info.Data.Query)),
		"navSession":   sessionMenu,
	}

	//return c.String(http.StatusOK, info)
	return c.Render(http.StatusOK, "execute", data)
}
