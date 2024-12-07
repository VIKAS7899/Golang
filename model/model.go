package model


type User struct {
	UID  int    `json:"uid"  binding:"required"`
	Name string `json:"name" binding:"required"`
	TID  int    `json:"tid"  binding:"required"`
}

type Task struct {
	Taskid    float64   `json:"tid"`
	Tasks     string    `json:"tasks"`
	Timestamp int		`json:"time"`
}