package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	helper "github.com/omjogani/postgre-multipartitions/helpers"
)

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	var stockName string
	var stocks []string

	rows, err := db.Query("SELECT name FROM stocks")
	if err != nil {
		log.Fatalln(err)
		return c.JSON("An error occurred")
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&stockName)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		stocks = append(stocks, stockName)
	}

	return c.Render("index", fiber.Map{
		"Stocks": stocks,
	})
}

func navigateRequest(c *fiber.Ctx, db []*sql.DB) error {
	return indexHandler(c, db[0])
}

func main() {
	fmt.Println("PostgreSQL Multi-Partitions")

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// load .env and it's values
	envError := godotenv.Load(".env")
	if envError != nil {
		fmt.Println("ERROR: ", envError)
		log.Fatal("Could not load Environment Variables...")
	}
	port := os.Getenv("PORT")
	POSTGRE_USERNAME := os.Getenv("POSTGRE_USERNAME")
	POSTGRE_PASSWORD := os.Getenv("POSTGRE_PASSWORD")
	POSTGRE_URL := os.Getenv("POSTGRE_URL")
	POSTGRE_DBNAME := os.Getenv("POSTGRE_DBNAME")

	// 2 Connection to the PostgreSQL
	connStrMaster0 := "postgresql://" + POSTGRE_USERNAME + ":" + POSTGRE_PASSWORD + "@" + POSTGRE_URL + "/" + POSTGRE_DBNAME + "0?sslmode=disable"
	connStrMaster1 := "postgresql://" + POSTGRE_USERNAME + ":" + POSTGRE_PASSWORD + "@" + POSTGRE_URL + "/" + POSTGRE_DBNAME + "1?sslmode=disable"

	var dbConnections []*sql.DB
	dbConnections = append(dbConnections, helper.ConnectToDb(connStrMaster0))
	dbConnections = append(dbConnections, helper.ConnectToDb(connStrMaster1))

	app.Get("/", func(c *fiber.Ctx) error {
		return navigateRequest(c, dbConnections)
	})

	if port == "" {
		port = "3000"
	}
	app.Static("/", "./public")
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}
