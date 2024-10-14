package dbHelper

import (
	"database/sql"
	"errors"
	"time"
	"todo/Database"
	"todo/Utils"
)

func AlreadyUser(email string) (bool, error) {
	SQLQuery := `Select count(email) > 0 from regisuser
						where email=TRIM(LOWER($1))
						AND archived_at IS NULL`
	var check bool
	AlreadyUserErr := Database.DBConnection.Get(&check, SQLQuery, email)
	if AlreadyUserErr != nil {
		return false, AlreadyUserErr
	}
	return check, nil
}

func RegisterUser(username, email, password string) error {

	SqlQuery := `INSERT INTO regisuser (username, email, password) 
								VALUES (TRIM($1),Trim($2),$3)`

	_, createErr := Database.DBConnection.Exec(SqlQuery, username, email, password)
	if createErr != nil {
		return createErr
	}
	return nil
}

func GetArchivedAt(sessionID string) (*time.Time, error) {
	var archivedAt *time.Time

	query := `SELECT archived_at 
              FROM user_sessions  
              WHERE id = $1`

	getErr := Database.DBConnection.Get(&archivedAt, query, sessionID)
	if getErr != nil {
		return nil, getErr // Return error if the query fails
	}

	return archivedAt, nil
}

func LoginCheck(email, password string) (string, string, string, error) {

	SqlQuery := `SELECT userId,username, email,password 
						from regisuser
						where archived_at IS NULL 
						   AND email = TRIM($1)`


	// you can make a model for this and you can use Get() instead of QueryRowx()
	var name string
	var userId string
	var Email string
	var hashPassword string

	//var name string
	//var userId string
	//var Email string
	//var hashPassword string

	body := struct {
		name         string
		userId       string
		Email        string
		hashPassword string
	}{}


	//TODO do not use queryRowx use GET method
	err := Database.DBConnection.Get(body, SqlQuery, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", "", nil
		}
		return "", "", "", err
	}

	passwordErr := Utils.CheckPassword(password, body.hashPassword)
	if passwordErr != nil {
		return "", "", "", passwordErr
	}
	return body.userId, body.Email, body.name, nil
}

func DeleteUserSession(sessionId string) error {
	query := `UPDATE user_sessions
			  SET archived_at = NOW()
			  WHERE id = $1
			    AND archived_at IS NULL`

	_, delErr := Database.DBConnection.Exec(query, sessionId)
	if delErr != nil {
		return delErr
	}
	return nil
}

func SessionGenerated(userId string) (string, error) {
	var sessionID string
	query := `INSERT INTO user_sessions(id) 
              VALUES ($1) RETURNING session_id`

	// use Get()
	crtErr := Database.DBConnection.QueryRow(query, userId).Scan(&sessionID)

	if crtErr != nil {
		return "", crtErr
	}
	return sessionID, nil
}

func DeleteUser(userId string) error {

	SqlQuery := `UPDATE regisuser SET
                    archived_at=now() 
                		WHERE userid= $1
                		  AND archived_at IS NULL`
	//res, err := Database.DBconn.Exec("DELETE FROM usertodo WHERE id=$1", param.Id)
	_, Err := Database.DBConnection.Exec(SqlQuery, userId)
	if Err != nil {
		return Err
	}
	return nil

}
