package dbHelper

import (
	"todo/Database"
	"todo/Models"
)

func CreateTodo(title, description, userId string) error {

	SqlQuery := `INSERT INTO todos (title, description, id)
								VALUES (TRIM($1),TRIM($2),$3)`

	_, creatErr := Database.DBConnection.Exec(SqlQuery, title, description, userId)
	if creatErr != nil {
		return creatErr
	}
	return nil
}

func GetTodoByName(title, userId string) ([]Models.Todos, error) {

	SqlQuery := `SELECT id,todo_id,title, description,is_completed 
							FROM todos
							WHERE id=$1 
							AND title ILIKE  '%'|| $2 || '%'`

	todos := make([]Models.Todos, 0)
	searchErr := Database.DBConnection.Select(&todos, SqlQuery, userId, title)
	return todos, searchErr
}

func UpdateTodo(title, description, userId string) error {

	SqlQuery := `UPDATE todos SET title = $1, description = $2  WHERE id = $3`

	_, updaErr := Database.DBConnection.Exec(SqlQuery, title, description, userId)
	if updaErr != nil {
		return updaErr
	}
	return nil

}

func MarkComplete(id, UserId string) error {

	SqlQuery := `UPDATE todos
              SET is_completed = true        
              WHERE todo_id = $1                  
                AND id = $2`

	_, updateErr := Database.DBConnection.Exec(SqlQuery, id, UserId)
	if updateErr != nil {
		return updateErr
	}
	return nil
}

func DeleteTodo(userId, todoId string) error {

	SqlQuery := `UPDATE todos
			  SET archieved_at = NOW()
			  WHERE todo_id = $1
			    AND id = $2
			    AND archieved_at IS NULL`

	_, delErr := Database.DBConnection.Exec(SqlQuery, todoId, userId)
	if delErr != nil {
		return delErr
	}
	return nil

}
