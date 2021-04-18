package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)
func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(heguicai.cool:3306)/test")
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()
	// 连接数据库
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}
	// entity
	type Person struct {
		Id         int
		First_Name string
		Last_Name  string
	}
	router := gin.Default()
	// GET 通过id获取数据  该方法id不可省略
	router.GET("/person/:id", func(c *gin.Context) {
		var (
			person Person
			result gin.H
		)
		id := c.Param("id")
		row := db.QueryRow("select id, first_name, last_name from person where id = ?;", id)
		err = row.Scan(&person.Id, &person.First_Name, &person.Last_Name)
		if err != nil {
			// If no results send null
			result = gin.H{
				"result": nil,
				"count":  0,
			}
		} else {
			result = gin.H{
				"result": person,
				"count":  1,
			}
		}
		c.JSON(http.StatusOK, result)
	})
	// GET all persons
	router.GET("/persons", func(c *gin.Context) {
		var (
			person  Person
			persons []Person
		)
		rows, err := db.Query("select id, first_name, last_name from person;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&person.Id, &person.First_Name, &person.Last_Name)
			persons = append(persons, person)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"result": persons,
			"count":  len(persons),
		})
	})
	// POST new person details
	router.POST("/person", func(c *gin.Context) {
		person :=&Person{}
		if c.Bind(person) ==nil{
			stmt, err := db.Prepare("insert into person (first_name, last_name) values(?,?);")
			if err != nil {
				fmt.Print(err.Error())
			}
			_, err = stmt.Exec(person.First_Name, person.Last_Name)
			if err != nil {
				fmt.Print(err.Error())
				c.JSON(http.StatusOK, gin.H{
					"error": err,
				})
			} else {defer stmt.Close()
			name := person.First_Name+" "+person.Last_Name
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf(" %s successfully created", name),
			})
			}
		}else {
			c.JSON(200, gin.H{
				"status":  "error",
				"message":nil,
			})
		}
	})
	// PUT - update a person details
	router.PUT("/person", func(c *gin.Context) {
		person :=&Person{}
		if c.Bind(person) ==nil{
			id := person.Id
			first_name := person.First_Name
			last_name := person.Last_Name
			stmt, err := db.Prepare("update person set first_name= ?, last_name= ? where id= ?;")
			if err != nil {
				fmt.Print(err.Error())
			}
			_, err = stmt.Exec(first_name, last_name, id)
			if err != nil {
				fmt.Print(err.Error())
				c.JSON(http.StatusOK, gin.H{
					"error": err,
				})
			}else {
				// Fastest way to append strings
				defer stmt.Close()
				name := first_name + " " + last_name
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("Successfully updated to %s", name),
				})
			}
		}else {
			c.JSON(200, gin.H{
				"status":  "error",
				"message":nil,
			})
		}
	})
	// Delete resources
	router.DELETE("/person", func(c *gin.Context) {
		id := c.Query("id")
		stmt, err := db.Prepare("delete from person where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(id)
		if err != nil {
			fmt.Print(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully deleted user: %s", id),
		})
	})
	router.Run(":3000")
}