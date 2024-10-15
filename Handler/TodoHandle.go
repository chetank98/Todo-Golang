package Handle

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"todo/Database/dbHelper"
	"todo/Middleware"
	"todo/Models"
	"todo/Utils"
)

//func GetAll(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	var allPost []Models.Todos
//	//defer Database.CloseDatabase()
//
//	GetQuery := `select todos.todoid,todos.userid,todos.name, todos.note,todos.iscompleted from public.todos`
//
//	//GetQuery := `select t.name, t.note, t.iscompleted from todos t inner join regisuser r ON r.userid = t.userid`
//
//	rows, err := Database.DBConnection.Query(GetQuery)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for rows.Next() {
//		var data Models.Todos
//
//		err := rows.Scan(&data.NoteId, &data.UserId, &data.Name, &data.Note, &data.IsCompleted)
//		if err != nil {
//			log.Print(err)
//		}
//
//		allPost = append(allPost, data)
//	}
//	fmt.Println(allPost)
//	err = json.NewEncoder(w).Encode(allPost)
//	if err != nil {
//		return
//	}
//}

func CreateNote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create a todo")

	var note Models.Todos

	userCtx := Middleware.UserContext(r)
	note.UserId = userCtx.UserID

	// use parseBody()
	dataErr := json.NewDecoder(r.Body).Decode(&note)
	fmt.Println(r.Body)
	if dataErr != nil {
		Utils.RespondError(w, http.StatusBadRequest, dataErr, "invalid payload")
	}

	// use Validator
	if note.Name == "" || note.Note == "" {
		Utils.RespondError(w, http.StatusBadRequest, nil, "Put some data")
		return
	}

	err := dbHelper.CreateTodo(note.Name, note.Note, note.UserId)
	if err != nil {
		Utils.RespondError(w, http.StatusInternalServerError, err, "failed to save the todo")
		return
	}

	Utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{"Todo Creation Successful"})

}

func GetTodoByName(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get todo by name")
	//todo use this in model
	data := struct {
		Title string `json:"title"`
	}{}

	userCtx := Middleware.UserContext(r)
	userId := userCtx.UserID

	dataErr := json.NewDecoder(r.Body).Decode(data)
	if dataErr != nil {
		Utils.RespondError(w, http.StatusBadRequest, dataErr, "Unable to Decode")
	}

	//TODO use validator
	if data.Title == "" {
		Utils.RespondError(w, http.StatusBadRequest, nil, "Enter the todo title")
		return
	}

	todos, getErr := dbHelper.GetTodoByName(data.Title, userId)
	if getErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to search todo")
		return
	}

	Utils.RespondJSON(w, http.StatusFound, todos)

}

func MarkCompleted(w http.ResponseWriter, r *http.Request) {
	//todo take this in query param
	data := struct {
		Id string `json:"id"`
	}{}

	userCtx := Middleware.UserContext(r)
	userId := userCtx.UserID

	dataErr := json.NewDecoder(r.Body).Decode(&data)
	if dataErr != nil {
		Utils.RespondError(w, http.StatusBadRequest, dataErr, "unable to extract data")
		return
	}

	if data.Id == "" {
		Utils.RespondError(w, http.StatusBadRequest, nil, "enter the correct title")
	}

	updateErr := dbHelper.MarkComplete(userId, data.Id)
	if updateErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, updateErr, "failed to update todo")
	}

	Utils.RespondJSON(w, http.StatusAccepted, struct {
		Message string `json:"message"`
	}{"Todo updated Successfully"})

}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {

	var note Models.Todos

	userCtx := Middleware.UserContext(r)
	note.UserId = userCtx.UserID

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		Utils.RespondError(w, http.StatusBadRequest, err, "invalid request payload")
		return
	}

	if note.Name == "" || note.Note == "" {
		Utils.RespondError(w, http.StatusBadRequest, nil, "Put some data")
		return
	}

	if err := dbHelper.UpdateTodo(note.Name, note.Note, note.UserId); err != nil {
		Utils.RespondError(w, http.StatusInternalServerError, err, "error updating todo")
		return
	}

	Utils.RespondJSON(w, http.StatusAccepted, struct {
		Message string `json:"message"`
	}{"Todo updated Successfully"})

}

func TodoDeleted(w http.ResponseWriter, r *http.Request) {

	//TODO use this in query param
	//data := struct {
	//	ID string `json:"id"`
	//}{}

	data := chi.URLParam(r, "id")

	userCtx := Middleware.UserContext(r)
	userID := userCtx.UserID

	dataErr := json.NewDecoder(r.Body).Decode(&data)
	if dataErr != nil {
		Utils.RespondError(w, http.StatusBadRequest, dataErr, "unable to extract data")
		return
	}

	saveErr := dbHelper.DeleteTodo("id", userID)
	if saveErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to delete todo")
		return
	}

	Utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"todo deleted successfully"})

}

//func UpdateTodo(w http.ResponseWriter, r *http.Request) {
//
//	fmt.Println("Update the user data")
//	var todo Models.Todos
//	if r.Body == nil {
//		w.Write([]byte("Insert the data for update..."))
//		return
//	}
//	json.NewDecoder(r.Body).Decode(&todo)
//
//	UpdateQuery := `Update todos Set name=$1,note=$2,iscompleted=$3,updatedat=$4 where todoid=$5`
//	res, err := Database.DBConnection.Exec(UpdateQuery, todo.Name, todo.Note, todo.IsCompleted, todo.UpdatedAt, todo.NoteId)
//	if err == nil {
//		count, err := res.RowsAffected()
//		if err == nil {
//			if count == 1 {
//				json.NewEncoder(w).Encode("To to updated Sucessfully")
//			} else {
//				json.NewEncoder(w).Encode("No todo found")
//			}
//		}
//	}
//
//}

//func DeleteNote(w http.ResponseWriter, r *http.Request) {
//
//	fmt.Println("Delete the todo")
//	Param := chi.URLParam(r, "id")
//
//	fmt.Println(Param)
//
//	deleteQuery := `Delete from Todos where todoid=$1`
//	res, err := Database.DBConnection.Exec(deleteQuery, Param)
//	if err == nil {
//		count, err := res.RowsAffected()
//		if err == nil {
//			if count == 1 {
//				json.NewEncoder(w).Encode("Deleted")
//				fmt.Println("Deleted")
//			} else {
//				json.NewEncoder(w).Encode("No rows found")
//			}
//		}
//	}
//
//}
