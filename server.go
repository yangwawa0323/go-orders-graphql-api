package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	// _ "github.com/go-sql-driver/mysql"
	"github.com/yangwawa0323/go-orders-graphql-api/graph"
	"github.com/yangwawa0323/go-orders-graphql-api/graph/generated"
	"github.com/yangwawa0323/go-orders-graphql-api/graph/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const defaultPort = "8080"
const dbName = "test_db4"

var db *gorm.DB

func initDB() *gorm.DB {
	var err error
	dataSourceName := "root:redhat@tcp(localhost:3306)/?parseTime=True"
	//db, err = gorm.Open("mysql", dataSourceName)
	db, err = gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	// db.LogMode(true)
	db.Debug()
	// Create the database. This is a one-time step.
	// Comment out if running multiple times - You may see an error otherwise
	db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	db.Exec(fmt.Sprintf("USE %s", dbName))
	// Migration to create tables for Order and Item schema
	db.AutoMigrate(&model.Order{}, &model.Item{})
	return db
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	initDB()

	//log.Fatal(db)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{
			DB: db,
		}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
