package database

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"
	"time"

	"os"
	"todo_app/model"
	

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)


var db *sql.DB
var mu sync.Mutex

func InitDB(){

	file,err := os.Create("database.db")
	if(err != nil){
		panic(err)
	}
	defer file.Close()

	db,err := sql.Open("sqlite3","./database.db")
	if(err != nil){
		panic(err)
	}
	defer db.Close()

	CreateUserTable(db)
	CreateTaskTable(db)
}

func CreateUserTable(db *sql.DB){
	
	fmt.Println("Creating a user table IF not exist")

	query := `CREATE TABLE IF NOT EXISTS USER (
		USER_ID INT AUTO_INCREMENT,
		NAME NVARCHAR(255),
		TASK_ID INT AUTO_INCREMENT
	);`

	_,err := db.Exec(query)

	if( err != nil){
		panic(err)
	}
}

func CreateTaskTable(db *sql.DB){
	
	fmt.Println("Creating a Task table IF not exist")

	query := `CREATE TABLE IF NOT EXISTS TASKS (
		TASK_ID FLOAT(53) NOT NULL,
		TASK NVARCHAR(100),
		EXECUTION_TIME DATETIME
	);`

	_,err := db.Exec(query)

	if( err != nil){
		panic(err)
	}
}

func CheckRow(user model.User , c *gin.Context) (int ,string , error){

	db,_ := sql.Open("sqlite3","./database.db")
	defer db.Close()

	row := db.QueryRow("SELECT USER_ID, NAME, TASK_ID FROM USER WHERE USER_ID = ? and NAME = ?", user.UID,user.Name)
	
	var users model.User

	// Scan the result into the user variable
	err := row.Scan(&users.UID, &users.Name, &users.TID)
	

	// for row.Next(){
	// 	var user model.User	

	// 	err := row.Scan(&user.UID,&user.Name,&user.TID)

	// 	if(err != nil){
	// 		return 0,"",err
	// 	}

	// 	return user.TID,user.Name,err
	// }
	if err != nil {

		if err == sql.ErrNoRows {
			// No rows were returned
			return 0, "", fmt.Errorf("user not present")
		}
		// Return an error if scanning fails
		return 0, "", err
	}

	// Return the user details and no error
	return users.TID, users.Name, nil
	
}

func NewRecord(user model.User,c *gin.Context) error{

	db,_ := sql.Open("sqlite3","./database.db")
	defer db.Close()

	if db == nil {
        err := fmt.Errorf("database connection is not initialized")
        fmt.Println(err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return err
    }
		result,err := db.Exec(
		"INSERT INTO USER VALUES(?,?,?);",
		user.UID,user.Name,user.TID)

		if (err != nil) {
			fmt.Println("error inserting new user", err.Error())
			return err
		   }
	  
		lastInsertId, _ := result.LastInsertId()
		
		fmt.Printf("Last Inserted ID: %d\n", lastInsertId)

	c.JSON(200,gin.H{"message":"New user inserted"})
	return nil
}

func UserTask(taskid float64) ([]string, error){

	db,_ := sql.Open("sqlite3","./database.db")
	defer db.Close()

	query := "SELECT TASK_ID,TASK FROM TASKS WHERE TASK_ID = ?"
	//fmt.Println(taskid)

	row,err := db.Query(query,taskid)

	if (err != nil) {
		fmt.Printf("Error executing query %v/n",err)
		//fmt.Println("error getting rows", err.Error())
		 
	   }

	var thing []string 

	for row.Next(){
		var t model.Task
		err := row.Scan(&t.Taskid,&t.Tasks)
		
		if(err != nil){

			fmt.Printf("Error executing query %v/n",err)
		//	fmt.Println("Error getting the rows")
					
		}
	thing = append(thing, t.Tasks)

	}
	return thing, err
}

func AddedTask(Task model.Task)(bool ,error){

	db,_ := sql.Open("sqlite3","./database.db")
	//defer db.Close()
	query := "INSERT INTO TASKS (TASK_ID,TASK,EXECUTION_TIME) VALUES(?,?,?) "

	fmt.Println(Task)

	result,err := db.Exec(query, Task.Taskid, Task.Tasks,time.Now().Add(time.Second * time.Duration(Task.Timestamp) ))

	if(err != nil){
		fmt.Printf("Error executing query %v/n",err)
		return false,err
	}

	lastinsertedId,_ := result.LastInsertId()

	fmt.Printf("Last inserted ID %d\n" , lastinsertedId)
	
	return true, nil
}

func DeleteTask(now time.Time){

	mu.Lock()
    defer mu.Unlock()

	query := "DELETE FROM TASK WHERE EXECUTION_TIME  <= ?"
	_,err := db.Exec(query,now)

	if(err != nil){
		fmt.Printf("Error deleting row %v",err)
	}

}


func Checktask(){
		
	mu.Lock()
    defer mu.Unlock()

	db,_ := sql.Open("sqlite3","./database.db")
	//defer db.Close()
	

	query := `SELECT TASK FROM TASKS WHERE EXECUTION_TIME <= ?  `

	row,err := db.Query(query,time.Now())

	if(err != nil){
		fmt.Printf("Error in executing quer in table %v",err)
		return
	}

	var notif []string
	
	for row.Next(){
		var t model.Task
		err := row.Scan(&t.Tasks)

		if(err != nil){
			fmt.Printf("Error in row mapping %v",err)
		}

		notif = append(notif,t.Tasks)

	}

	//fmt.Println(len(notif))
	
	if(len(notif) >0 ){

		fmt.Println("Hi Vikas , The time is ,",time.Now(),"I got a message for you")
		for _,msg := range notif{
			fmt.Println(msg)
		}

		//DeleteTask(now)
	}

}