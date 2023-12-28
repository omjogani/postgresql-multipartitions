package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	helper "github.com/omjogani/postgre-multipartitions/helpers"
)

type Stock struct {
	Name  string
	Price float64
}

func fetchRecords(db *sql.DB) []Stock {
	var stocksMaster0 []Stock

	// load From Master 0
	rowsMaster0, err := db.Query("SELECT name, price FROM stocks")
	if err != nil {
		log.Fatalln(err)
		return stocksMaster0
	}

	defer rowsMaster0.Close()

	for rowsMaster0.Next() {
		var stock Stock
		err = rowsMaster0.Scan(&stock.Name, &stock.Price)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		stocksMaster0 = append(stocksMaster0, stock)
	}

	return stocksMaster0
}

func postHandler(c *fiber.Ctx, db []*sql.DB) error {
	newStock := Stock{}
	if err := c.BodyParser(&newStock); err != nil {
		log.Printf("An Error: %v", err)
		return c.SendString(err.Error())
	}
	fmt.Printf("%v", newStock)
	if newStock.Name != "" {

		if strings.ToLower(newStock.Name)[0] >= 'a' && strings.ToLower(newStock.Name)[0] <= 'm' {
			_, errorQuery := db[0].Exec("INSERT into stocks VALUES ($1, $2)", newStock.Name, newStock.Price)
			if errorQuery != nil {
				log.Fatalf("An error while executing query: %v", errorQuery)
			}
		} else {
			_, errorQuery := db[1].Exec("INSERT into stocks VALUES ($1, $2)", newStock.Name, newStock.Price)
			if errorQuery != nil {
				log.Fatalf("An error while executing query: %v", errorQuery)
			}
		}

	}
	return c.Redirect("/")
}

func indexHandler(c *fiber.Ctx, db []*sql.DB) error {
	stocksMaster0 := fetchRecords(db[0])
	stocksMaster1 := fetchRecords(db[1])
	return c.Render("index", fiber.Map{
		"M0Stocks": stocksMaster0,
		"M1Stocks": stocksMaster1,
	})
}

func main() {
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
		return indexHandler(c, dbConnections)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		return postHandler(c, dbConnections)
	})

	if port == "" {
		port = "3000"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}
