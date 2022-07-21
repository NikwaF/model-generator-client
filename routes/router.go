package routes

import (
	"echo-api/controllers"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TemplateRenderer ...
type TemplateRenderer struct {
	templates *template.Template
}

// Render TemplateRenderer
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func add(x, y int) int {
	return x + y
}

// Init Routes
func Init() *echo.Echo {
	e := echo.New()
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	//templates
	funcs := template.FuncMap{
		"add": add,
	}
	tmp := template.Must(template.New("main").Funcs(funcs).ParseGlob("view/*.html"))
	renderer := &TemplateRenderer{
		templates: template.Must(tmp.ParseGlob("view/template/*.html")),
	}
	e.Renderer = renderer

	//Asset Files
	e.Static("/assets", "assets")

	e.GET("/login", controllers.LoginPages)
	e.POST("/login-process", controllers.LoginProcess)
	e.GET("/logout", controllers.LogoutProcess)

	e.GET("/", controllers.Dashboard)
	e.GET("/api-generator", controllers.GenV1)
	e.GET("/generate-csv-getdata", controllers.GenV2)
	// e.GET("/list-endpoint-api", controllers.ListApi)
	// e.GET("/edit-endpoint-api/:id", controllers.ListEditApi)
	// e.POST("/preview-endpoint", controllers.PreviewEndpoint)
	// e.POST("/save-endpoint", controllers.SaveEndpoint)
	// e.POST("/update-endpoint", controllers.UpdateEndpoint)
	// e.POST("/test-query-endpoint", controllers.TestQuery)
	// e.GET("/delete-endpoint/:id", controllers.DeleteEndpoint)
	// e.GET("/execute-api", controllers.ExecApi)
	// e.GET("/execute-api/:id", controllers.ExecApi2)
	// e.POST("/execute-api", controllers.ExecQuery)
	e.GET("/get-connections", controllers.FetchConnection)
	e.POST("/generate-getdata", controllers.GenerateGetData)
	e.POST("/generate-csv-getdata", controllers.GenerateCsvGetData)

	//niko

	return e
}
