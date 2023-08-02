package model

import (
	"crud-todo/todo/config"
	"fmt"
)

type Todo struct {
	Id        int
	Item      string
	Completed int
}

type View struct {
	Todos []Todo
}

/*
 * Global Variables
 */
var (
	id        int
	item      string
	completed int
	database  = config.Database()
)

/*
 * This function is used to reterive all the entries present in the database
 *
 * @return View
 */
func Show() View {
	statement, err := database.Query(`SELECT * FROM todos`)

	handleError(err)

	var todos []Todo

	for statement.Next() {
		err = statement.Scan(&id, &item, &completed)

		handleError(err)

		todo := Todo{
			Id:        id,
			Item:      item,
			Completed: completed,
		}

		todos = append(todos, todo)
	}

	data := View{
		Todos: todos,
	}

	return data
}

/*
 * This function is used to add the new entry to database
 *
 * @param item string
 * @return bool
 */
func Add(item string) bool {
	if item != "" {
		_, err := database.Exec(`INSERT INTO todos (item) VALUE (?)`, item)

		handleError(err)
	} else {
		fmt.Printf("Item should not be empty")
	}

	return true
}

/*
 * This function is used to get the entry based on the id from the database
 *
 * @param id string
 * @return View
 */
func Edit(id string) View {
	statement, err := database.Query("SELECT * FROM todos WHERE id=?", id)

	handleError(err)

	var todos []Todo

	for statement.Next() {
		var id, completed int
		var item string
		err = statement.Scan(&id, &item, &completed)

		handleError(err)

		todo := Todo{
			Id:        id,
			Item:      item,
			Completed: completed,
		}

		todos = append(todos, todo)
	}

	data := View{
		Todos: todos,
	}

	return data
}

/*
 * This function is used to update the entry in database
 *
 * @param id string
 * @param item string
 * @return bool
 */
func Update(id string, item string) bool {
	if item != "" {
		_, err := database.Exec(`UPDATE todos SET item = ? WHERE id = ?`, item, id)

		handleError(err)

	} else {
		fmt.Printf("Item should not be empty")
	}

	return true
}

/*
 * This function is used to delete the entry in database based on the id
 *
 * @param id string
 * @return bool
 */
func Delete(id string) bool {
	_, err := database.Exec(`DELETE FROM todos WHERE id = ?`, id)

	handleError(err)

	return true
}

/*
 * This function is used to complete the item in database based on the id
 *
 * @param id string
 * @return bool
 */
func Complete(id string) bool {
	_, err := database.Exec(`UPDATE todos SET completed = 1 WHERE id = ?`, id)

	handleError(err)

	return true
}

/*
 * This function is used to handle the error function
 *
 * @param err error
 * @return bool
 */
func handleError(err error) bool {
	if err != nil {
		panic(err.Error())
	}

	return true
}
