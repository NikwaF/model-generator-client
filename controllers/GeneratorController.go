package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"echo-api/models"
	"echo-api/utils"

	"github.com/labstack/echo/v4"
)

var quotation = strings.NewReplacer("'", "''", "\"", "''", "%P", "%p", "%D", "%d")
var baseURL = "http://appapi2.indomaret.lan:3000"

func PreviewEndpoint(c echo.Context) error {
	returnData := models.GenerateResponse{
		STATUS:  "ERROR - 404",
		MESSAGE: "NOT FOUND",
		JSON:    models.InsertStruct{},
		HTML:    "",
	}
	if (c.FormValue("Inpendpoint") != "") && (c.FormValue("Inptype") != "") && (c.FormValue("Inpquery") != "") && (c.FormValue("Inpkoneksi") != "") {
		changeQuery := quotation.Replace(c.FormValue("Inpquery"))

		countParams := strings.Count(changeQuery, "%p")
		countField := strings.Count(changeQuery, "%d")

		params, err := strconv.Atoi(c.FormValue("Inpparam"))
		if err != nil {
			params = 0
		}
		datas, err := strconv.Atoi(c.FormValue("Inpdata"))
		if err != nil {
			datas = 0
		}

		var inpType string
		if (c.FormValue("Inptype") == "GET") || (c.FormValue("Inptype") == "READ") {
			inpType = "READ"
		} else {
			inpType = c.FormValue("Inptype")
		}

		if (countParams == params) && (countField == datas) {
			arrData := models.InsertStruct{
				Endpoint:      c.FormValue("Inpendpoint"),
				Type_endpoint: inpType,
				JmlData:       c.FormValue("Inpdata"),
				JmlParam:      c.FormValue("Inpparam"),
				Dbcode:        c.FormValue("Inpkoneksi"),
				Query:         changeQuery,
			}

			returnData = models.GenerateResponse{
				STATUS:  "SUCCESS",
				MESSAGE: "",
				JSON:    arrData,
				HTML:    GenerateBox(countField, countParams),
			}
		} else {
			returnData = models.GenerateResponse{
				STATUS:  "ERROR - 404",
				MESSAGE: "FIELD ATAU PARAMETER TIDAK SESUAI",
				JSON:    models.InsertStruct{},
				HTML:    "",
			}
		}
	}

	return c.JSON(http.StatusOK, returnData)
}

func GenerateGetData(c echo.Context) error {
	u := new(models.InputGeneratorGet)
	c.Bind(u)

	// queryKirim := strings.ReplaceAll(u.Query, "\n", " ")
	// fixQueryKirim := strings.ReplaceAll(queryKirim, `"`, "")

	// reg1 := regexp.MustCompile(`\s+`)
	res := &u.Query
	utils.CleansingQuery(res)

	fmt.Println(&res)
	// queryKirim = strings.ReplaceAll(u.Query, `"`, "")
	idconn := u.Idcon

	// var datalisttable models.TableInfoList

	var jsonStr = []byte(fmt.Sprintf(`{"idconn":%s , "query":"%s"}`, idconn, *res))

	datalisttable := utils.GetTableInformation(&jsonStr)

	if datalisttable.Status == "error" {
		return c.JSON(http.StatusOK, datalisttable)
	}
	// var nmModel string = u.Model
	// var nmGetFunction string = u.Fungsi

	tampungan := utils.TampunganModelGetFungsi(datalisttable, u.Model, u.Fungsi, res)

	return c.JSON(http.StatusOK, models.BalikanGetData{Status: "success", Message: "", Text: tampungan})
}

func GenerateCsvGetData(c echo.Context) error {
	u := new(models.InputGeneratorCsv)
	c.Bind(u)
	res := &u.Query
	idconn := u.Idcon

	utils.CleansingQuery(res)
	fmt.Println(&res)

	var jsonStr = []byte(fmt.Sprintf(`{"idconn":%s , "query":"%s"}`, idconn, *res))

	datalisttable := utils.GetTableInformation(&jsonStr)

	if datalisttable.Status == "error" {
		return c.JSON(http.StatusOK, datalisttable)
	}

	tampungan := utils.TampunganModelGetFungsi(datalisttable, u.Model, u.Fungsi, res)
	utils.TampunganGetDataCsv(datalisttable, u.Model, u.Fungsi, u.Nmcsv, &tampungan)

	return c.JSON(http.StatusOK, models.BalikanGetData{Status: "success", Message: "", Text: tampungan})
}

func FetchConnection(c echo.Context) error {
	returnData := models.ConnectionList{}

	var client = &http.Client{}
	var data models.ConnectionList

	request, _ := http.NewRequest("GET", baseURL+"/getConnections", nil)

	response, err := client.Do(request)
	if err != nil {
		returnData = models.ConnectionList{
			Status:  "error",
			Message: string(err.Error()),
		}
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&data)

	if err != nil {
		returnData = models.ConnectionList{
			Status:  "error",
			Message: string(err.Error()),
		}
	}

	if returnData.Status == "error" {
		return c.JSON(http.StatusInternalServerError, returnData)
	}

	return c.JSON(http.StatusOK, data)
}

func SaveEndpoint(c echo.Context) error {
	returnData := models.GenerateResponse{
		STATUS:  "ERROR",
		MESSAGE: "404 - NOT FOUND",
		JSON:    models.InsertStruct{},
		HTML:    "",
	}

	if (c.FormValue("Inpendpoint") != "") && (c.FormValue("Inptype") != "") && (c.FormValue("Inpquery") != "") && (c.FormValue("Inpkoneksi") != "") {
		changeQuery := quotation.Replace(c.FormValue("Inpquery"))

		countParams := strings.Count(changeQuery, "%p")
		countField := strings.Count(changeQuery, "%d")

		params, err := strconv.Atoi(c.FormValue("Inpparam"))
		if err != nil {
			params = 0
		}
		datas, err := strconv.Atoi(c.FormValue("Inpdata"))
		if err != nil {
			datas = 0
		}

		var inpType string
		if (c.FormValue("Inptype") == "GET") || (c.FormValue("Inptype") == "READ") {
			inpType = "READ"
		} else {
			inpType = c.FormValue("Inptype")
		}

		if (countParams == params) && (countField == datas) {
			arrData := models.InsertStruct{
				Endpoint:      c.FormValue("Inpendpoint"),
				Type_endpoint: inpType,
				JmlData:       c.FormValue("Inpdata"),
				JmlParam:      c.FormValue("Inpparam"),
				Dbcode:        c.FormValue("Inpkoneksi"),
				Query:         changeQuery,
			}

			// post request
			jsonData, _ := json.Marshal(arrData)
			request, err := http.Post(target_url+"/helper_api/add", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println(err.Error())
				returnData = models.GenerateResponse{
					STATUS:  "ERROR",
					MESSAGE: "500 - INTERNAL SERVER ERROR",
					JSON:    models.InsertStruct{},
					HTML:    "",
				}
			}

			defer request.Body.Close()

			// read response
			body, err := ioutil.ReadAll(request.Body)
			if err != nil {
				fmt.Println(err.Error())
			}

			var responseMessage models.StatusResponse
			err = json.Unmarshal(body, &responseMessage)
			if err != nil {
				fmt.Println(err.Error())
			}

			if (responseMessage.ERROR == "") && (responseMessage.SUCCESS == true) {
				returnData = models.GenerateResponse{
					STATUS:  "SUCCESS",
					MESSAGE: "",
					JSON:    arrData,
					HTML:    "",
				}
			}
		}
	}

	return c.JSON(http.StatusOK, returnData)
}

func GenerateBox(fields int, params int) string {
	var field_box string = ""
	var param_box string = ""
	var returnHtml string = ""

	if fields > 0 {
		for fHtml := 1; fHtml <= fields; fHtml++ {
			field_box += `<div class="form-group row">
			<label class="col-lg-2 col-form-label">Fields</label>
			<div class="col-lg-10">
				<input type="text" name="field-bar[]" placeholder="Field" class="form-control">
			</div>
		</div>`
		}
	}

	if params > 0 {
		for pHtml := 1; pHtml <= params; pHtml++ {
			param_box += `<div class="form-group row">
			<label class="col-lg-2 col-form-label">Params</label>
			<div class="col-lg-10">
				<input type="text" name="param-bar[]" placeholder="Param" class="form-control">
			</div>
		</div>`
		}
	}

	returnHtml = "<label for='' class='mt-3'><b>Validate Query</b></label><div class='row mb-2'><div class='col-md-6'>" + field_box + "</div><div class='col-md-6'>" + param_box + "</div></div>"
	returnHtml += "<div class='row mb-2'><div class='col-md-9'></div><div class='col-md-3'><button class='btn btn-success btn-block'>Test Query</button></div></div>"

	return returnHtml
}

func UpdateEndpoint(c echo.Context) error {
	returnData := models.GenerateResponse{
		STATUS:  "ERROR",
		MESSAGE: "404 - NOT FOUND",
		JSON:    models.InsertStruct{},
		HTML:    "",
	}

	if (c.FormValue("Inpid") != "") && (c.FormValue("Inpendpoint") != "") && (c.FormValue("Inptype") != "") && (c.FormValue("Inpquery") != "") && (c.FormValue("Inpkoneksi") != "") {
		changeQuery := quotation.Replace(c.FormValue("Inpquery"))

		countParams := strings.Count(changeQuery, "%p")
		countField := strings.Count(changeQuery, "%d")

		params, err := strconv.Atoi(c.FormValue("Inpparam"))
		if err != nil {
			params = 0
		}
		datas, err := strconv.Atoi(c.FormValue("Inpdata"))
		if err != nil {
			datas = 0
		}

		var inpType string
		if (c.FormValue("Inptype") == "GET") || (c.FormValue("Inptype") == "READ") {
			inpType = "READ"
		} else {
			inpType = c.FormValue("Inptype")
		}

		if (countParams == params) && (countField == datas) {
			arrData := models.InsertStruct{
				Endpoint:      c.FormValue("Inpendpoint"),
				Type_endpoint: inpType,
				JmlData:       c.FormValue("Inpdata"),
				JmlParam:      c.FormValue("Inpparam"),
				Dbcode:        c.FormValue("Inpkoneksi"),
				Query:         changeQuery,
				Id:            c.FormValue("Inpid"),
			}

			// post request
			jsonData, _ := json.Marshal(arrData)
			request, err := http.Post(target_url+"/helper_api/update", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println(err.Error())
				returnData = models.GenerateResponse{
					STATUS:  "ERROR",
					MESSAGE: "500 - INTERNAL SERVER ERROR",
					JSON:    models.InsertStruct{},
					HTML:    "",
				}
			}

			defer request.Body.Close()

			// read response
			body, err := ioutil.ReadAll(request.Body)
			if err != nil {
				fmt.Println(err.Error())
			}

			var responseMessage models.StatusResponse
			err = json.Unmarshal(body, &responseMessage)
			if err != nil {
				fmt.Println(err.Error())
			}

			if (responseMessage.ERROR == "") && (responseMessage.SUCCESS == true) {
				returnData = models.GenerateResponse{
					STATUS:  "SUCCESS",
					MESSAGE: "",
					JSON:    arrData,
					HTML:    "",
				}
			}
		}
	}

	return c.JSON(http.StatusOK, returnData)
}

func DeleteEndpoint(c echo.Context) error {
	id := c.Param("id")

	response, err := http.Post(target_url+"/helper_api/Delete/"+id, "application/json", nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer response.Body.Close()

	c.Redirect(http.StatusFound, base_url+"/list-endpoint-api")
	return nil
}
