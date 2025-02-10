package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/templates"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/static", "static")
	e.GET("/", func(c echo.Context) error {
		return Render(c, http.StatusOK, templates.Alerts())
	})
	e.Logger.Fatal(e.Start(config.Env().HTTPAddr + ":" + config.Env().HTTPPort))
}

func Render(ctx echo.Context, statusCode int, template templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)
	template = templates.Base(template)
	if err := template.Render(ctx.Request().Context(), buf); err != nil {
		return RenderError(ctx, http.StatusInternalServerError, err)
	}
	return ctx.HTML(statusCode, buf.String())
}

func RenderError(ctx echo.Context, statusCode int, err error) error {
	return ctx.HTML(statusCode, err.Error())
}
