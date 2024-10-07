package Handle

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"todo/Database"
	"todo/Models"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var allPost []Models.Todos
	//defer Database.CloseDatabase()

	//GetQuery := `select todos.todoid,todos.userid,todos.name, todos.note,todos.iscompleted from public.todos`

	GetQuery := `select t.name, t.note, t.iscompleted from todos t inner join regisuser r ON r.userid = t.userid`

	rows, err := Database.DBConnection.Query(GetQuery)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var data Models.Todos

		err := rows.Scan(&data.NoteId, &data.UserId, &data.Name, &data.Note, &data.IsCompleted)
		if err != nil {
			log.Print(err)
		}

		allPost = append(allPost, data)
	}
	fmt.Println(allPost)
	err = json.NewEncoder(w).Encode(allPost)
	if err != nil {
		return
	}
}

func GetById(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")

	fmt.Println(param)
	var note Models.Todos
	queryById := `Select todoid,userid,name,note,iscompleted from Todos where todoid=$1`

	err := Database.DBConnection.Get(&note, queryById, param)
	if err != nil {
		w.Write([]byte("no data found"))
		return
	}
	//fmt.Println(res)
	json.NewEncoder(w).Encode(note)

}

func CreateNote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create a todo")

	if r.Body == nil {
		w.Write([]byte("Insert the data..."))
		return
	}
	var note Models.Todos

	_ = json.NewDecoder(r.Body).Decode(&note)

	fmt.Println(&note)

	SQLQuery := `Insert into Todos ( name, note, iscompleted) values ($1,$2,$3) returning todoid`

	res, err := Database.DBConnection.Exec(SQLQuery, note.Name, note.Note, note.IsCompleted)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(res)

}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Update the user data")
	var todo Models.Todos
	if r.Body == nil {
		w.Write([]byte("Insert the data for update..."))
		return
	}
	json.NewDecoder(r.Body).Decode(&todo)

	UpdateQuery := `Update todos Set name=$1,note=$2,iscompleted=$3,updatedat=$4 where todoid=$5`
	res, err := Database.DBConnection.Exec(UpdateQuery, todo.Name, todo.Note, todo.IsCompleted, todo.UpdatedAt, todo.NoteId)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			if count == 1 {
				json.NewEncoder(w).Encode("To to updated Sucessfully")
			} else {
				json.NewEncoder(w).Encode("No todo found")
			}
		}
	}

}

func DeleteNote(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Delete the todo")
	Param := chi.URLParam(r, "id")

	fmt.Println(Param)

	// need to check archived_at us null or not
	// we don't do hard delete just insert the timestamp in archived_at that is soft delete
	// send proper response in JSON format to client

	deleteQuery := `Delete from Todos where todoid=$1`
	res, err := Database.DBConnection.Exec(deleteQuery, Param)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			if count == 1 {
				json.NewEncoder(w).Encode("Deleted")
				fmt.Println("Deleted")
			} else {
				json.NewEncoder(w).Encode("No rows found")
			}
		}
	}

}
