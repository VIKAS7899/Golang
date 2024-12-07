package main


import (
	"fmt"
	"time"
	"todo_app/database"
	"todo_app/handler"
	"todo_app/middleware"
	"github.com/gin-gonic/gin"

)

func main(){

	fmt.Println("welcome ")

	fmt.Println("Welcome Vikas , we are here for tasks right , The time is :",time.Now().Format("Jan 02 2006, 3:04;05"))

	database.InitDB()  // Creating DB 

	router := gin.Default()

	router.Use(gin.Recovery())

	//go starttimer()

	public := router.Group("/")
	{
		
		public.POST("/login",handler.Login)
		public.POST("/register",handler.Register)
		
	}

	protected := router.Group("/todo") // list all the savved task 
	protected.Use(middleware.Authenticationmiddleware)
	{
		
		go starttimer()
			
		protected.GET("/task",handler.GetTask)
		protected.POST("/add",handler.AddTask)
		//protected.DELETE("/delete",handler.Delete)
		//time.Sleep(200 * time.Second)
	}


	router.Run(":3000")

}