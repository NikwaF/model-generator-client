package utils

import (
	"bytes"
	"echo-api/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var baseURL = "http://appapi2.indomaret.lan:3000"

func GetTableInformation(params *[]byte) *models.TableInfoList {
	var datalisttable models.TableInfoList
	req, err := http.NewRequest("POST", baseURL+"/testpost", bytes.NewBuffer(*params))
	// req.Header.Set("X-Custom-Header", "myvalue")
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&datalisttable)

	if err != nil {
		datalisttable.Status = "error"
		datalisttable.Message = string(err.Error())
		return &datalisttable
	}

	return &datalisttable
}

func TampunganModelGetFungsi(infotable *models.TableInfoList, model string, fungsi string, query *string) string {
	var nmModel string = model
	var nmGetFunction string = fungsi

	var tampungan string

	tampungan = fmt.Sprintf(`public class %s`, nmModel)
	tampungan += "\n"
	tampungan += `{`
	tampungan += "\n"
	// fmt.Println(tampungan)
	// fmt.Println(infotable.Data)

	for _, item := range infotable.Data {
		typeDataCsharp := models.TypeData[item.DATATYPE]
		tampungan += fmt.Sprintf(`	public %s %s { get; set; }`, typeDataCsharp, item.COLUMNNAME)
		tampungan += "\n"
	}
	tampungan += "\n"
	tampungan += "\n"

	tampungan += fmt.Sprintf(`	public List<%s> %s(){`, nmModel, nmGetFunction)
	tampungan += "\n"
	tampungan += fmt.Sprintf(`		string query = "%s"; //isi querynya disini`, *query)
	tampungan += "\n"
	tampungan += fmt.Sprintf(`		List<%s> list = new List<%s>();`, nmModel, nmModel)
	tampungan += "\n"

	tampungan += fmt.Sprintf(`
		try
		{
			using (var conn = new OracleConnection(SiteSettings.ConnectionString)) //ganti ke string koneksi yang bener ya 
			{
				using (var cmd = conn.CreateCommand())
				{
					conn.Open();
					cmd.CommandText = query;
					using (var reader = cmd.ExecuteReader())
					{
						while (reader.Read())
						{
								%s record = new %s();`, nmModel, nmModel)
	tampungan += "\n"

	for i, item := range infotable.Data {
		switch item.DATATYPE {
		case "VARCHAR2":
			tampungan += fmt.Sprintf(`								record.%s = reader.IsDBNull(%d) ? string.Empty : reader.GetString(%d); //%s`, item.COLUMNNAME, i, i, item.COLUMNNAME)
		case "CHAR":
			tampungan += fmt.Sprintf(`								record.%s = reader.IsDBNull(%d) ? string.Empty : reader.GetString(%d); //%s`, item.COLUMNNAME, i, i, item.COLUMNNAME)
		case "NUMBER":
			tampungan += fmt.Sprintf(`								record.%s = reader.IsDBNull(%d) ? 0 : reader.GetDecimal(%d); //%s`, item.COLUMNNAME, i, i, item.COLUMNNAME)
		case "DATE":
			tampungan += fmt.Sprintf(`								record.%s = reader.IsDBNull(%d) ? string.Empty : reader.GetDateTime(%d).ToString("dd-MMM-yyyy"); //%s`, item.COLUMNNAME, i, i, item.COLUMNNAME)
		}
		tampungan += "\n"
	}
	tampungan += "\n"
	tampungan += `								list.Add(record);
						}
					}

					conn.Close();
				}
			}
		}
		catch (Exception ex)
		{
				Console.WriteLine(ex.InnerException != null ? ex.InnerException.Message : ex.Message);
				var message = ex.InnerException != null ? ex.InnerException.Message : ex.Message;
				list.Clear();
		}
`
	tampungan += "\n"
	tampungan += `		return list;`
	tampungan += "\n"
	tampungan += `	}`
	tampungan += "\n"
	tampungan += `}`

	return tampungan
}

func TampunganGetDataCsv(infotable *models.TableInfoList, model string, fungsi string, nmcsv string, tampungan *string) {
	*tampungan += "public virtual ActionResult ExportCsv()"
	*tampungan += "\n"
	*tampungan += "{"
	*tampungan += fmt.Sprintf(`%s obj = new %s();`, model, model)
	*tampungan += "\n"
	*tampungan += fmt.Sprintf(`var list = obj.%s();`, fungsi)
	*tampungan += "\n"
	*tampungan += "var folderPath = HttpContext.Server.MapPath(SiteSettings.TempDirectory);"
	*tampungan += "\n"
	*tampungan += fmt.Sprintf(`string csvname = "%s" + DateTime.Now.ToString("ddMMyyHHmmss") + ".csv";`, nmcsv)
	*tampungan += "\n"
	*tampungan += "using (StreamWriter sw = new StreamWriter(new FileStream(Path.Combine(folderPath, csvname), FileMode.Create))){"
	*tampungan += "\n"
	var header string = ""
	var isibody string = ""
	for i, item := range infotable.Data {
		if header == "" {
			header += item.COLUMNNAME
		} else {
			header += "|" + item.COLUMNNAME
		}
		if item.DATATYPE != "VARCHAR2" && item.DATATYPE != "CHAR" {
			isibody += fmt.Sprintf(`s[%d] = list[i].%s.ToString();`, i, item.COLUMNNAME)
		} else {
			isibody += fmt.Sprintf(`s[%d] = list[i].%s;`, i, item.COLUMNNAME)
		}
		isibody += "\n"
	}
	*tampungan += fmt.Sprintf(`sw.WriteLine("%s");`, header)
	*tampungan += "for (Int32 i = 0; i < list.Count(); i++){"
	isibody += "\n"
	*tampungan += fmt.Sprintf(`string[] s = new string[%d];`, len(infotable.Data))
	isibody += "\n"
	*tampungan += isibody
	isibody += "\n"
	*tampungan += `sw.WriteLine(string.Join("|", s));`
	*tampungan += "\n"
	*tampungan += "}"
	*tampungan += "}"
	*tampungan += `return File(Path.Combine(folderPath, csvname), System.Net.Mime.MediaTypeNames.Application.Octet, csvname);`
	*tampungan += "\n"
	*tampungan += "}"
	*tampungan += "\n"
	// sw.WriteLine("KD_PROMO|KODEBIN|DCI|TOKO|TGLTRANS|TGLAWAL|TGLAKHIR|NO_KARTU|NOTRANSAKSI|NO_STRUK_PENJUALAN|PLU|DESKRIPSI|QTY|NILAI|NILAICLAIM|TOTALBELANJA");
}

func CleansingQuery(query *string) {
	queryKirim := strings.ReplaceAll(*query, "\n", " ")
	fixQueryKirim := strings.ReplaceAll(queryKirim, `"`, "")

	reg1 := regexp.MustCompile(`\s+`)
	*query = reg1.ReplaceAllString(fixQueryKirim, " ")
}
