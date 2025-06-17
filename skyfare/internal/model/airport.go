package model

type Coordinates struct {
	Lat string `json:"latitude"`
	Lng string `json:"longitude"`
}

type Airport struct {
	Code        string      `json:"code"`
	Name        string      `json:"name"`
	Country     string      `json:"country"`
	Coordinates Coordinates `json:"coordinates"`
}
