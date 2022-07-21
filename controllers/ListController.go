package controllers

import (
	"fmt"
	"bytes"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"echo-api/models"

	"github.com/labstack/echo/v4"
)

var newQuotation = strings.NewReplacer("\"", "'", "%P", "%p", "%D", "%d")

func ListEdit(id string) (models.ResponseData, error) {
	var responseData models.ResponseData
	response, err := http.Get(target_url+"/helper_api/"+id)
	if err != nil { fmt.Println(err.Error()) }

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil { fmt.Println(err.Error()) }

	err = json.Unmarshal(body, &responseData)
	if err != nil { fmt.Println(err.Error()) }

	return responseData, nil
}

func TestQuery(c echo.Context) error {
	var testData models.TestData
	var response []map[string]interface{}
	var response_error models.StatusResponse

	err := c.Bind(&testData)
	if err != nil { fmt.Println(err.Error()) }

	if testData.Type_endpoint == "GET" {
		testData.Type_endpoint = "READ"
	}
	if testData.Type_endpoint == "INSERT" {
		testData.Type_endpoint = "CREATE"
	}

	//check query param data
	testData.Query = newQuotation.Replace(testData.Query)
	countParams := strings.Count(testData.Query, "%p")
	countField := strings.Count(testData.Query, "%d")

	// if count not match
	if (countField != len(testData.Data)) || (countParams != len(testData.Params)) {
		errorCount := models.StatusResponse {
			ERROR : "Jumlah Parameter atau Field tidak sesuai dengan Query",
			SUCCESS : false,
		}
	
		return c.JSON(http.StatusOK, errorCount)
	}
	//end check

	jsonData , _ := json.Marshal(testData)
	request, err := http.Post(target_url+"/Happy/TestCall", "application/json", bytes.NewBuffer(jsonData))
	if err != nil { fmt.Println("request - " + err.Error()) }

	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil { fmt.Println("body - " + err.Error()) }

	err = json.Unmarshal(body, &response)
	if err != nil { 
		fmt.Println("response - " + err.Error())
		json.Unmarshal(body, &response_error)

		return c.JSON(http.StatusOK, response_error)
	}

	return c.JSON(http.StatusOK, response)
}