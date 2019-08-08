package models

import (
	
)

type Scrum_poker struct {
	Id int `json:"id"`
	Url_code string `json:"url_code"`
	User_id int `json:"user_id"`
	Poker string `json:"poker"`
}

func GetOne(where map[string]interface{}) (test Scrum_poker) {
	db.Where(where).First(&test)

	return
}

func UpdPoker(where map[string]interface{},data Scrum_poker) bool {
	db.Model(&Scrum_poker{}).Where(where).Updates(data)

	return true
}
func AddPoker(poker Scrum_poker) bool {
	db.Create(&poker)

	return true
}


