package helpers

import (
	"fmt"
	"github.com/VishalTanwani/gochat/apiserver/internal/config"
	"net/http"
	"runtime/debug"
)

var app *config.AppConfig

//NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

//ServerError will handle server errors
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
