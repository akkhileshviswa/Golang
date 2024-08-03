package routes

import (
	"crud-todo/todo/controller"

	"github.com/gorilla/mux"
)

/*
 * This function initializes the route functions.
 *
 * @return mux.Router
 */
func Init() *mux.Router {
	route := mux.NewRouter()

	route.HandleFunc("/", controller.Show)
	route.HandleFunc("/add", controller.Add).Methods("POST")
	route.HandleFunc("/edit/{id}", controller.Edit)
	route.HandleFunc("/update", controller.Update).Methods("POST")
	route.HandleFunc("/delete/{id}", controller.Delete)
	route.HandleFunc("/complete/{id}", controller.Complete)
	return route
}
