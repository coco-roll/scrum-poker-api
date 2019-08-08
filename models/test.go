package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Test struct {
	Model
	NAME string `json:"name"`
}

func GetTestModel(id int) (test Test) {
	db.Where("id = ?", id).First(&test)

	return
}

func EditTestModel(id int, data interface{}) bool {
	db.Model(&Test{}).Where("id = ?", id).Updates(data)

	return true
}

func DelTestModel(id int, data interface{}) bool {
	db.Model(&Test{}).Where("id = ?", id).Delete(data)

	return true
}

func AddTest(data map[string]interface{}) bool {
	db.Create(&Test{
		NAME: data["name"].(string),
	})

	return true
}

func (t *Test) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())

	return nil
}

func (t *Test) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())

	return nil
}
