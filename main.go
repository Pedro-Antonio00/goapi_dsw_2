package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "modernc.org/sqlite"
)

var db *gorm.DB
var err error

//Struct da pessoa
type Person struct {
	ID   string `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
	Address string `json:"address"`
	City string `json:"city"`
	Fone string `json:"fone"`
}

func main() {
	//Configurar conexão com o banco de dados SQLite
	db, err = gorm.Open("sqlite", "test.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	defer db.Close()

	//Se não existe, ela é criada
	db.AutoMigrate(&Person{})

	//Inicializar o roteador do gin
	r := gin.Default()

	//Definir os endpoints
	r.GET("/people", GetPeople)
	r.GET("/people/:id", GetPerson)
	r.POST("/people", CreatePerson)
	r.PUT("/people/:id", UpdatePerson)
	r.DELETE("/people/:id", DeletePerson)

	//Configurar um canal para capturar sinais do sistema operacional
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	//Executar o servidor em uma goroutine
	go func() {
		err := r.Run(":8080")
		if err != nil {
			log.Fatal(err)
		}
	}()

	//Aguarde sinais para encerrar o programa
	<-stopChan

	fmt.Println("Encerrando o programa...")
}

//Obter todas as pessoas
func GetPeople(c *gin.Context) {
	var people []Person
	if err := db.Find(&people).Error; err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
	} else {
		c.JSON(200, people)
	}
}

//Obter uma pessoa por ID
func GetPerson(c *gin.Context) {
	id := c.Params.ByName("id")
	var person Person
	if err := db.Where("id = ?", id).First(&person).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, person)
	}
}

//Criar pessoa
func CreatePerson(c *gin.Context) {
	var person Person
	c.BindJSON(&person)

	//Gerar UUID
	person.ID = uuid.New().String()

	db.Create(&person)
	c.JSON(200, person)
}

//Atualizar pessoa
func UpdatePerson(c *gin.Context) {
	id := c.Params.ByName("id")
	var person Person
	if err := db.Where("id = ?", id).First(&person).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	c.BindJSON(&person)
	db.Save(&person)
	c.JSON(200, person)
}

//Excluir pessoa
func DeletePerson(c *gin.Context) {
	id := c.Params.ByName("id")
	var person Person
	d := db.Where("id = ?", id).Delete(&person)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}
