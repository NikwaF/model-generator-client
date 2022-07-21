package models

type InsertStruct struct {
	Endpoint      string `json: "endpoint"`
	Type_endpoint string `json: "type_endpoint"`
	JmlData       string `json: "jmlData"`
	JmlParam      string `json: "jmlParam"`
	Dbcode        string `json: "dbcode"`
	Query         string `json: "query"`
	Id            string `json: "id"`
}

type GenerateResponse struct {
	STATUS  string       `json: "status"`
	MESSAGE string       `json: "message"`
	JSON    InsertStruct `json : "json"`
	HTML    string       `json: "html"`
}

type StatusResponse struct {
	ERROR   string `json: "error"`
	SUCCESS bool   `json: "success"`
}

type ResponseData struct {
	DATA    InsertStruct `json: "data"`
	MESSAGE string       `json: "message"`
	STATUS  string       `json: "status"`
}

type TestData struct {
	DBCode        string   `json:"DBCode"`
	Query         string   `json:"query"`
	Type_endpoint string   `json:"type_endpoint"`
	Data          []string `json:"data"`
	Params        []string `json:"params"`
}

type ConnectionList struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    []struct {
		ID        int    `json:"ID"`
		NMCON     string `json:"NM_CON"`
		CONSTRING string `json:"CON_STRING"`
		USERNAME  string `json:"USERNAME"`
		PASSWORD  string `json:"PASSWORD"`
	} `json:"data"`
}

type TableInfoList struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    []struct {
		COLUMNNAME string `json:"COLUMN_NAME"`
		DATATYPE   string `json:"DATA_TYPE"`
		DATALENGTH int    `json:"DATA_LENGTH"`
	} `json:"data"`
}

var TypeData = map[string]string{
	"VARCHAR2": "string",
	"NUMBER":   "decimal",
	"CHAR":     "string",
	"DATE":     "string",
}

type BalikanGetData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Text    string `json:"text"`
}

type InputGeneratorGet struct {
	Query  string `json:"query"`
	Idcon  string `json:"idcon"`
	Model  string `json:"model"`
	Fungsi string `json:"fungsi"`
}

type InputGeneratorCsv struct {
	Query  string `json:"query"`
	Idcon  string `json:"idcon"`
	Model  string `json:"model"`
	Fungsi string `json:"fungsi"`
	Nmcsv  string `json:"nmcsv"`
}
