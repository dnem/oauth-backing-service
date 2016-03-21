package server

import (
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

//NewServer configures and returns a Negroni server
func NewServer(appEnv *cfenv.App) *negroni.Negroni {

	n := negroni.Classic()
	router := mux.NewRouter()
	router.HandleFunc("/api/hello", helloHandler())

	n.UseHandler(router)
	return n
}
