package handler

import (
	"fmt"
	"net/http"
	"todo_app/database"
	"todo_app/model"
	"todo_app/utils"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context){

	var user model.User

	err := c.ShouldBindJSON(&user)

	if(err != nil){
		c.JSON(http.StatusBadRequest,gin.H{"error": " Header body is empty , no user login"})
	}
	
	//fmt.Println(user.Name)

	taskid,username,err := database.CheckRow(user,c)

	if(err != nil){
		fmt.Printf("Error inserting user: %v\n", err)
		c.JSON(http.StatusBadRequest,gin.H{"error": " User not found"})
		return
	}

	utils.GenerateToken(taskid,username,c)

	
}


func Register(c *gin.Context){

	var user model.User
	
	if err := c.ShouldBindJSON(&user); err != nil {
		// Log the error details for debugging
		fmt.Printf("Error binding JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	if user.UID <= 0 || user.Name == "" || user.TID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		return
	} 
	fmt.Printf("Received user: %v\n", user) // Print user details for debugging

	// Attempt to insert the new user into the database
	if err := database.NewRecord(user,c); err != nil {
		// Log the error details
		fmt.Printf("Error inserting user: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})

}

func GetTask(c *gin.Context){

	taskid := utils.GetTaskID(c)

	

	AllTasks ,err := database.UserTask(taskid)

	if(err != nil){
		c.JSON(http.StatusNotFound,gin.H{"error":"Error in showing task"})
		return
	}

	c.JSON(200,gin.H{"Tasks" : AllTasks})
}


func AddTask(c *gin.Context){

	var Task model.Task	
	
	err := c.ShouldBindJSON(&Task)

	if(err != nil){
		
		c.JSON(http.StatusBadRequest,gin.H{"error":"Error in getting task"})
		return 
	}

	Task.Taskid  = utils.GetTaskID(c) // Getting the task id using toeken in context

	ok,err := database.AddedTask(Task)

	if(!ok || err != nil){
		fmt.Printf("Error inserting task: %v\n", err)
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Error in inserting task"})
		return 

	}

	if(ok){
		c.JSON(200,gin.H{"Messages" : "Tasks added succesfully"})
	}

}

// func Delete(c * gin.Context){

// 	row := c.Param("taskno")

// 	_,err := database.DeleteTask(c,row)

// 	if(err != nil){
// 		fmt.Printf("Error inserting user: %v\n", err)
// 		c.JSON(http.StatusInternalServerError,gin.H{"error":"Error in inserting task"})

// 	}
// 	c.JSON(200,gin.H{"Messages" : "Tasks added succesfully"})
// }