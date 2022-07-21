package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type DataDb struct {
	Id            string `json:"id"`
	Dbcode        string `json:"dbcode"`
	Endpoint      string `json:"endpoint"`
	Type_endpoint string `json:"type_endpoint"`
	Query         string `json:"query"`
	JmlData       string `json:"jmlData"`
	JmlParam      string `json:"jmlParam"`
}

type ResponseDb struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    DataDb
}

type QueryBuilder struct {
	DBCode        string   `json:"DBCode"`
	Query         string   `json:"query"`
	Type_endpoint string   `json:"type_endpoint"`
	Data          []string `json:"data"`
	Params        []string `json:"params"`
}

type AjaxRequest struct {
	Tipe  string       `json:"tipe"`
	Query QueryBuilder `json:"data"`
}

const (
	base_link      string = "http://localhost:8081/"
	base_happy_url string = base_link + "Happy/"
	info_db_url    string = base_link + "helper_api/"
	test_call_url  string = base_happy_url + "TestCall"
)

func GetInfoDbcode(id string) (ResponseDb, error) {
	var data ResponseDb
	url := info_db_url + id
	resp, err := http.Get(url)
	var (
		error_api       = errors.New("get api error")
		error_not_found = errors.New("no data")
	)

	if err != nil {
		return data, error_api
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, error_api
	}

	err = json.Unmarshal(body, &data)

	fmt.Println(err)
	if err != nil {
		return data, err
	}

	if data.Status == "error" {
		return data, error_not_found
	}

	data.Data.Type_endpoint = strings.ToUpper(data.Data.Type_endpoint)
	data.Data.Query = strings.ToUpper(data.Data.Query)

	return data, nil
}

func MilihHtml(tipe, jmlparam, jmldata string) string {
	param, _ := strconv.Atoi(jmlparam)
	data, _ := strconv.Atoi(jmldata)

	if param < 1 && data < 1 {
		return ""
	}

	return generateInputForm(param, data)
}

func generateInputForm(jmlparam int, jmldata int) string {
	const tag_input string = `<input class="form-control" value=""`
	var tag_param string = ""
	var tag_data string = ""

	for i := 1; i <= jmlparam; i++ {
		urutan := strconv.Itoa(i)
		tag_param += fmt.Sprintf(`<label>PARAM %s</label> %s name="params[]" %s />`, urutan, tag_input, urutan)
	}

	for i := 1; i <= jmldata; i++ {
		urutan := strconv.Itoa(i)
		tag_data += fmt.Sprintf(`<label>DATA %s</label> %s name="datas[]" %s />`, urutan, tag_input, urutan)
	}

	if jmlparam > 0 {
		tag_param = `<div class="mt-3"><h3>Parameter</h3></div><div>` + tag_param + `</div>`
	}

	if jmldata > 0 {
		tag_data = "<hr><div><h3>Data</h3></div><div>" + tag_data + "</div>"
	}

	tag_param += tag_data
	return tag_param
}

func htmlQuery(query string) string {
	tagQuery := `<code style="word-break:break-all;">`
	r := regexp.MustCompile(`(%P)`)
	query = r.ReplaceAllString(query, `<mark style="background:yellow;"><span>${1}</span></mark>`)
	r = regexp.MustCompile(`(%D)`)
	query = r.ReplaceAllString(query, `<mark style="background:rgb(198, 233, 157) none repeat scroll 0% 0%;"><span>${1}</span></mark>`)
	tagQuery += query + "</code>"
	return tagQuery
}

func ExecQuery(c echo.Context) error {
	var builder AjaxRequest
	err := c.Bind(&builder)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err})
	}

	if builder.Tipe == "GET" {
		builder.Tipe = "READ"
	}
	if builder.Tipe == "INSERT" {
		builder.Tipe = "CREATE"
	}

	params := QueryBuilder{Params: builder.Query.Params, Data: builder.Query.Data, DBCode: builder.Query.DBCode, Query: builder.Query.Query, Type_endpoint: builder.Tipe}

	data, err := LemparQuery(test_call_url, &params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err})
	}

	return c.JSON(http.StatusOK, data)
}

func LemparQuery(url string, hola *QueryBuilder) (interface{}, error) {
	var response []map[string]interface{}
	var response_error map[string]interface{}
	jsonVal, _ := json.Marshal(hola)

	request, err := http.Post(url, "application/json", bytes.NewBuffer(jsonVal))

	if err != nil {
		fmt.Println("181", err)
		return response, err
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println("187", err)
		return response, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		json.Unmarshal(body, &response_error)
		return response_error, nil
	}

	return response, nil
}
