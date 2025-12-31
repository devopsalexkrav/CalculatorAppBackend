package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/google/uuid"
	"github.com/Knetic/govaluate"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	//TODO : discuss with gpt &Calculation{} in AutoMigrate
	if err := db.AutoMigrate(&Calculation{}); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
}

type Calculation struct {
	ID         string `gorm:"primaryKey" json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

type CalculationRequest struct {
	Expression string `json:"expression"`
}

func calculateExpression(expression string) (string, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)

	if err != nil {
		return "", err
	}

	result, err := expr.Evaluate(nil)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", result), err
}

func getCalculations(c echo.Context) error{
	var calculations []Calculation

	if err := db.Find(&calculations).Error; err != nil{
		return c.JSON(http.StatusInternalServerError, map[string]string{"error":"Could not get calculations"})
	}

	return c.JSON(http.StatusOK, calculations)
}

func postCalculations(c echo.Context) error{
	var req CalculationRequest

	if err := c.Bind(&req); err != nil{
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Request"})
	}

	result, err := calculateExpression(req.Expression)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
	}

	calc := Calculation{
		ID: uuid.NewString(),
		Expression: req.Expression,
		Result: result,
	}

	if err := db.Create(&calc).Error; err != nil{
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not add calculation"})
	}

	return c.JSON(http.StatusOK, calc)
}

func patchCalculations(c echo.Context) error{
	id := c.Param("id")

	var req CalculationRequest

	if err := c.Bind(&req); err != nil{
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Request"})
	}

	result, error := calculateExpression(req.Expression)

	if error != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not find expression"})
	}	

	var calc Calculation

	if err := db.First(&calc, "id = ?", id).Error; err != nil{
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Calculation not found"})
	}

	calc.Expression = req.Expression
	calc.Result = result

	if err := db.Save(&calc).Error; err != nil{
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update calculation"})
	}

	return c.JSON(http.StatusOK, calc)
}

func deleteCalculations(c echo.Context) error{
	id := c.Param("id")

	//TODO : discuss with gpt &Calculation{} in delete
	if err := db.Delete(&Calculation{}, "id = ?", id).Error; err != nil{
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete calculation"})
	}	

	return c.NoContent(http.StatusNoContent)
}	  

func main() {
	initDB()

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/calculations", getCalculations)
	e.POST("/calculations", postCalculations)
	e.PATCH("/calculations/:id", patchCalculations)
	e.DELETE("/calculations/:id", deleteCalculations)

	e.Start("localhost:8000")
}
