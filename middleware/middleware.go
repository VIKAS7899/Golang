package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"todo_app/utils"

	"github.com/gin-gonic/gin"
)

func Authenticationmiddleware(c *gin.Context){

	token := c.GetHeader("Authorization")

	
	tokenstring := strings.Split(token," " )

	 index := len(tokenstring[1])-1

	 tokenstring[1] = tokenstring[1][1:index]

	fmt.Println(tokenstring[1])

	if(tokenstring[0] != "Bearer" || tokenstring[1] == ""){
		c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid token and session"})
		return
	}

	if(token == ""){
		c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid token and session"})
		return
	}
	

	claim ,err := utils.VerifyToken(tokenstring[1])

	if(err != nil){
		fmt.Printf("Authorization error %v/n",err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error":"Token cannot created"})
		c.Abort()
	}

	//fmt.Println(claim["tid"])

	exp,_ := claim["exp"].(float64)

	expirationTime := time.Unix(int64(exp), 0)
	fmt.Println("Token expires at:", expirationTime)

	c.Set("taskid", claim["tid"])
	// c.Set("name", claim["nam"])
	c.Next()

	fmt.Println("User is legal and TaskIS is set using Token")

	
}