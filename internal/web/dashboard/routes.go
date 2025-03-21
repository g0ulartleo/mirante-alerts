package dashboard

import (
	"log"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	"github.com/g0ulartleo/mirante-alerts/internal/web/dashboard/templates"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, signalService *signal.Service, alarmService *alarm.AlarmService) {

	e.GET("/", func(c echo.Context) error {
		alarmSignals, err := GetAlarmSignals(signalService, alarmService)
		if err != nil {
			log.Printf("Error fetching config signals: %v", err)
			return RenderError(c, http.StatusInternalServerError, err)
		}
		return Render(c, http.StatusOK, templates.Alarms(alarmSignals))
	})

	e.GET("/*", func(c echo.Context) error {
		pathParam := c.Param("*")
		var level int
		var baseURL string
		if pathParam == "" {
			level = 0
			baseURL = "/"
		} else {
			segments := strings.Split(pathParam, "/")
			level = len(segments)
			baseURL = "/" + pathParam
		}
		alarmSignals, err := GetAlarmSignals(signalService, alarmService)
		if err != nil {
			log.Printf("Error fetching config signals: %v", err)
			return RenderError(c, http.StatusInternalServerError, err)
		}

		return Render(c, http.StatusOK, templates.Treemap(alarmSignals, level, baseURL))
	})
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
