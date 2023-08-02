package controller

import (
	"crud-todo/todo/model"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

/*
 * Global Variables
 */
var (
	view = template.Must(template.ParseFiles("./views/index.html"))
	edit = template.Must(template.ParseFiles("./views/edit.html"))
)

/*
 * This function is to show the available data and form to the user
 */
func Show(w http.ResponseWriter, r *http.Request) {
	var data = model.Show()

	view.Execute(w, data)
}

/*
 * This function is called on clicking the add button
 */
func Add(w http.ResponseWriter, r *http.Request) {
	item := r.FormValue("item")

	model.Add(item)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

/*
 * This function is called on clicking the edit button
 */
func Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var data = model.Edit(id)

	edit.Execute(w, data)
}

/*
 * This function is called on clicking the update button
 */
func Update(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	item := r.FormValue("item")

	model.Update(id, item)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

/*
 * This function is called on clicking the delete button
 */
func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	model.Delete(id)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

/*
 * This function is called on clicking the complete button
 */
func Complete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	model.Complete(id)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
