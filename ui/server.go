package ui
import (
	"net/http"
)

func Start() {
	http.Handle("/", http.FileServer(http.Dir("./ui/static/")))
	http.ListenAndServe(":9998", nil)
}