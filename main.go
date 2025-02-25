package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Destination struct
type Destination struct {
	ID              uint   `gorm:"primaryKey;autoIncrement"`
	DestinationId   uint   `json:"id"`
	Label           string `json:"label"`
	ProvinceName    string `json:"province_name"`
	CityName        string `json:"city_name"`
	DistrictName    string `json:"district_name"`
	SubdistrictName string `json:"subdistrict_name"`
	ZipCode         string `json:"zip_code"`
}

var db *gorm.DB

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s user=%s password='%s' dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	fmt.Println("DSN:", dsn) // Debugging output

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate
	db.AutoMigrate(&Destination{})

	// Create router
	r := gin.Default()
	r.GET("/destinations", getDestinations)

	// Run server
	r.Run(":8080")
}

func getDestinations(c *gin.Context) {
	query := c.Query("search")

	// Fetch from DB
	var destinations []Destination
	// Search across multiple columns using OR conditions
	searchQuery := db.Where(
		"label ILIKE ? OR province_name ILIKE ? OR city_name ILIKE ? OR district_name ILIKE ? OR subdistrict_name ILIKE ? OR zip_code ILIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%",
	)

	// Apply limit & offset for pagination
	result := searchQuery.Limit(100).Offset(0).Find(&destinations)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}

	if len(destinations) > 0 {
		c.JSON(http.StatusOK, gin.H{"data": destinations})
		return
	}
	// If no results, fetch from RajaOngkir
	newDestinations, err := fetchFromRajaOngkir(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from RajaOngkir"})
		return
	}
	// Store fetched data
	tx := db.Begin()
	if err := tx.Create(&newDestinations).Error; err != nil {
		tx.Rollback() // Rollback on error
	} else {
		tx.Commit() // Commit changes
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{"data": newDestinations})
}

func fetchFromRajaOngkir(query string) ([]Destination, error) {
	apiKey := os.Getenv("RAJAONGKIR_API_KEY")
	client := resty.New()

	resp, err := client.R().
		SetHeader("key", apiKey).
		SetQueryParam("search", query).
		Get("https://rajaongkir.komerce.id/api/v1/destination/domestic-destination")

	if err != nil {
		return nil, err
	}

	log.Println("Response Body:", resp.String())
	// Parse response and map to Destination struct (modify based on RajaOngkir's response)
	var result struct {
		Data []Destination `json:"data"`
	}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}
