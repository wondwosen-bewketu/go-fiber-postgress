package main

import (
	"fmt"
	"log"
	"os"

	"net/http"

	"github.com/wondwosen-bewketu/go-fiber-postgress/storage" // Keep as is
	"github.com/wondwosen-bewketu/go-fiber-postgress/models"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	// "golang.org/x/mod/sumdb/storage" // Alias this import
	"gorm.io/gorm"
)

type Book struct {
	Author       string     `json:"author"`
	Title        string     `json:"title"`
	Publisher    string     `json:"publisher"`  
}
type Respostory struct{
	DB *gorm.DB
}

func(r *Respostory) CreateBook(context *fiber.Ctx) error{
	book := Book{}
	err := context.BodyParser(&book)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message":"request field"})
			return err
	
	}
	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not create book"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message":"book has been added"})
	return nil 
}
func (r *Respostory) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(bookModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete book",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book delete successfully",
	})
	return nil
}

func(r *Respostory) GetBooks(context *fiber.Ctx) error{
	bookModels := &[]models.Books{}

err := r.DB.Find(bookModels).Error
if err != nil {
	context.Status(http.StatusBadRequest).JSON(
		&fiber.Map{"message":"Could Not Get Books"})
	return err
}
context.Status(http.StatusOK).JSON(&fiber.Map{"message":"books fetched successfully",
"data": bookModels})
	return nil
}

func(r *Respostory) GetBookByID(context *fiber.Ctx) error{
	id := context.Params("id")
	bookModel := &models.Books{}
	if id == ""{
		context.Status(http.StatusInternalServerError).JSON(
		&fiber.Map{"message":"ID can not be empty"})

return nil
	}
fmt.Println("The ID is", id)
err := r.DB.Where("id  = ?",id).First(bookModel).Error

	if err != nil {
context.Status(http.StatusBadRequest).JSON(
		&fiber.Map{"message":"Could Not Get The Book"})
	return err
	
}
context.Status(http.StatusOK).JSON(&fiber.Map{"message":"Book Id Fetched Successfully",
"data": bookModel,
})
return nil
}

func (r *Respostory) SetupRoutes(app *fiber.App){
	api := app.Group("/api")
	api.Post("/create_books",r.CreateBook)
	api.Delete("delete_book/:id",r.DeleteBook)
	api.Get("get_books/:id",r.GetBookByID)
	api.Get("books",r.GetBooks)
}
func main() {
	err := godotenv.Load(".env")
	if err != nil{
		log.Fatal(err)
	}
	config := &storage.Config{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User: os.Getenv("DB_USER"),
		SSLMode: os.Getenv("DB_SSLMODE"),
		DBName: os.Getenv("DB_NAME"),
	}
	db, err := storage.NewConnection(config)

	if err != nil{
		log.Fatal(("Could not load the database"))

	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("Could Not Migrate DB")
	}
	r := Respostory{
		DB: db,
		
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
