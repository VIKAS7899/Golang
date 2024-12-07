package utils

import (
	"fmt"
	
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	
)

var privatekey = []byte("secreatecode")
var Token string = ""

func GenerateToken(taskid int , username string , c *gin.Context){

	fmt.Println("Good User found , Creating token for the same ",taskid,username)

	claim := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{	
		"nam" : username,
		"tid" : taskid,              
		"exp": time.Now().Add(time.Hour*10).Unix(),
		"iat": time.Now().Unix(),                 // Issued at
		
	})

	token , err := claim.SignedString(privatekey)

	if(err != nil){
		c.JSON(http.StatusBadGateway,gin.H{"error":"Error generating token"})
	}

	Token = string(token)

	fmt.Println("The token for the string ",Token)

	c.Request.Header.Add("Authorization", "Bearer " + Token) // Add token as  a header to Context 

	c.JSON(200,gin.H{"token":token})

}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
    // Parse the token
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Check the signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
        return privatekey, nil
    })

    // Check for errors	
    if err != nil {
        return nil, err
    }

    // Validate the token
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("Invalid token")
}

func GetTaskID(c *gin.Context) float64{
		
	task_id,exist := c.Get("taskid")

	fmt.Println(c.Get("taskid"))
	
	if( !exist){
		
		c.JSON(http.StatusBadRequest,gin.H{"error":"Error in getting  task id"})
		 return 0
	}

	Taskid, _ := task_id.(float64)

	return Taskid
}