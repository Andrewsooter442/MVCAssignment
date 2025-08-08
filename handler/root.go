package handler

import (
	"fmt"
	"net/http"
)

type User struct {
}

func (app *Application) HandleRootRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handle root request")

}
