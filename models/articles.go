package models

type Article struct {
	Id       int
	Hot_word string
	Title    string
	Info     string
	Img     string
}


type One_content struct {
	Type int
	Text string
}