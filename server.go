package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	// _ "github.com/go-sql-driver/mysql"

	"github.com/yangwawa0323/go-orders-graphql-api/graph/handlers"
	"github.com/yangwawa0323/go-orders-graphql-api/graph/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const defaultPort = "8080"
const dbName = "test_db4"

var db *gorm.DB

func initDB() *gorm.DB {
	var err error
	dataSourceName := "root:redhat@tcp(localhost:3306)/" + dbName + "?parseTime=True"
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

// Defining the Graphql handler
// func graphqlHandler() gin.HandlerFunc {
// 	// NewExecutableSchema and Config are in the generated.go file
// 	// Resolver is in the resolver.go file
// 	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
// 		DB: db,
// 	}}))

// 	return func(c *gin.Context) {
// 		h.ServeHTTP(c.Writer, c.Request)
// 	}
// }

// // Defining the Playground handler
// func playgroundHandler() gin.HandlerFunc {
// 	h := playground.Handler("GraphQL", "/query")

// 	return func(c *gin.Context) {
// 		h.ServeHTTP(c.Writer, c.Request)
// 	}
// }

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	initDB()

	// srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
	// 	Resolvers: &graph.Resolver{
	// 		DB: db,
	// 	}}))

	// http.Handle("/", playground.Handler("Query playground", "/query"))
	// http.Handle("/query", srv)

	// log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))
	// // router.Run(":" + defaultPort)

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.POST("/api/v1/", handlers.GraphqlHandler(db))
	r.GET("/", handlers.PlaygroundHandler())

	r.GET("/js/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "js.html", nil)
	})

	// Router must load before gin.Engine is running

	r.Run(":" + defaultPort)

}
