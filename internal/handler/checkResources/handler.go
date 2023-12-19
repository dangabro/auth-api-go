package checkResources

import (
	"auth-api-go/internal/data"
	"auth-api-go/internal/handler"
	"auth-api-go/internal/service/dao"
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"net/http"
	"strings"
)

//go:embed template.html
var crTemplate string

type checkResources struct {
	config data.Config
	db     *sql.DB
}

func New(config data.Config, db *sql.DB) http.Handler {
	return &checkResources{
		config: config,
		db:     db,
	}
}

func (h *checkResources) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	res, err := h.process(ctx, writer, r)

	if err != nil {
		handler.Error(err, writer, http.StatusBadRequest)
		return
	}

	_ = handler.WriteHtmlResponse(writer, res)
}

func (h *checkResources) process(ctx context.Context, writer http.ResponseWriter, r *http.Request) (any, error) {
	var list []string

	list = append(list, "<h2>Version</h2>")
	list = append(list, fmt.Sprintf("<div>%s</div>", data.SOFTWARE_VERSION))
	list = append(list, "<h2>Configuration</h2>")
	list = append(list, "<ul>")
	configurationSettings := h.config
	list = append(list, processConfigValue("Port:", configurationSettings.Port))
	list = append(list, processConfigValue("Token validity MS:", configurationSettings.TokenDurationMs))
	list = append(list, processConfigValue("Cors:", configurationSettings.Cors))
	list = append(list, processConfigValue("Log level:", configurationSettings.LogLevel))
	list = append(list, "</ul>")
	list = append(list, "<h2>Database</h2>")
	list = append(list, "<ul>")
	list = append(list, processConfigValue("User:", configurationSettings.Db.User))
	list = append(list, processConfigValue("Database:", configurationSettings.Db.Database))
	list = append(list, processConfigValue("Machine:", configurationSettings.Db.Machine))
	list = append(list, processConfigValue("Port:", configurationSettings.Db.Port))
	list = append(list, processConfigValue("Pool size:", configurationSettings.Db.PoolSize))
	list = append(list, "</ul>")

	response, err := dao.CheckConnection(ctx, h.db)
	var er bool = false
	var message string

	if err != nil {
		message = "error:" + err.Error()
		er = true
	} else {
		message = "success; database timestamp: " + response
	}

	var style string
	if er {
		style = "color: red; font-weight: bold"
	} else {
		style = "color: green; font-weight: bold"
	}

	htmlStuff := "<div style=\"%s\">%s</div>"
	htmlStuff = fmt.Sprintf(htmlStuff, style, message)
	list = append(list, htmlStuff)

	strList := strings.Join(list, "\n")

	res := strings.Replace(crTemplate, "@@@", strList, 1)

	return res, nil
}

func processConfigValue(name string, configParam any) string {
	var format string
	switch configParam.(type) {
	case bool:
		format = "%t"
	case string:
		format = "%s"
	case int64:
		format = "%d"
	case int32:
		format = "%d"
	default:
		format = "%s"
	}

	strFormat := fmt.Sprintf("<li><strong>%%s</strong> - %s</li>", format)
	strVal := fmt.Sprintf(strFormat, name, configParam)
	return strVal
}
