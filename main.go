package main

import (
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	var template Template
	db.First(&template, 1)
	templateXML := template.XML

	var objects []Object
	db.Where("templates_ID = ?", template.ID).Find(&objects)

	for _, object := range objects {
		expression, err := govaluate.NewEvaluableExpression(object.Formula)

		parameters := make(map[string]interface{}, 8)
		parameters["Q1"] = 2
		parameters["Q2"] = 1
		parameters["Q3"] = 3

		result, err := expression.Evaluate(parameters)

		if err != nil {
			fmt.Println(err)
		}

		objectXML := object.XML

		objectXML = strings.ReplaceAll(
			objectXML,
			"{value}",
			fmt.Sprintf("%v", result),
		)

		templateXML = strings.ReplaceAll(
			templateXML,
			fmt.Sprintf("{O:%d}", object.ID),
			objectXML,
		)
	}

	fmt.Println(templateXML)
}

// Calculator is a representation of calcultaors table instance
type Calculator struct {
	ID   int
	Name string
}

// Template is a representation of templates table instance
type Template struct {
	ID  int
	XML string
}

// Object is a representation of objects table instance
type Object struct {
	ID          int
	XML         string
	Condition   string
	Formula     string
	TemplatesID int
}
