// package main

// import (
// 	"bytes"
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"time"

// 	_ "github.com/lib/pq"
// )

// // Database credentials
// const (
// 	DbHost     = "localhost"
// 	DbName     = "stocks"
// 	DbUser     = "postgres"
// 	DbPassword = "root"
// )

// // Groww API URL
// const growwAPIURL = "https://groww.in/v1/api/stocks_data/v1/all_stocks"

// // Request payload
// type RequestPayload struct {
// 	ListFilters map[string][]string       `json:"listFilters"`
// 	ObjFilters  map[string]map[string]any `json:"objFilters"`
// 	Page        string                     `json:"page"`
// 	Size        string                     `json:"size"`
// 	SortBy      string                     `json:"sortBy"`
// 	SortType    string                     `json:"sortType"`
// }

// // API response structure
// type APIResponse struct {
// 	Records []Company `json:"records"`
// }

// // Company structure
// type Company struct {
// 	CompanyName     string `json:"companyName"`
// 	GrowwContractID string `json:"growwContractId"`
// }

// // Function to fetch data from API
// func fetchCompaniesFromAPI(page int) ([]Company, error) {
// 	payload := RequestPayload{
// 		ListFilters: map[string][]string{
// 			"INDEX": {
// 				"Nifty Bank", "BSE 100", "Nifty 100", "Nifty 50",
// 				"Nifty Next 50", "Nifty Midcap 100", "SENSEX",
// 			},
// 			"INDUSTRY": {},
// 		},
// 		ObjFilters: map[string]map[string]any{
// 			"CLOSE_PRICE": {"max": 500000, "min": 0},
// 			"MARKET_CAP":  {"min": 0, "max": 3000000000000000},
// 		},
// 		Page:     fmt.Sprintf("%d", page),
// 		Size:     "15",
// 		SortBy:   "NA",
// 		SortType: "ASC",
// 	}

// 	jsonPayload, err := json.Marshal(payload)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req, err := http.NewRequest("POST", growwAPIURL, bytes.NewBuffer(jsonPayload))
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var apiResponse APIResponse
// 	if err := json.Unmarshal(body, &apiResponse); err != nil {
// 		return nil, err
// 	}

// 	return apiResponse.Records, nil
// }

// // Function to insert or update company in PostgreSQL
// func insertOrUpdateCompany(db *sql.DB, company Company) error {
// 	apiUrl := fmt.Sprintf("https://groww.in/v1/api/groww-news/v2/stocks/news/%s?page=0&size=10", company.GrowwContractID)
// 	currentTime := time.Now()

// 	// Check if the company already exists
// 	var existingID int
// 	err := db.QueryRow("SELECT id FROM companies WHERE name = $1", company.CompanyName).Scan(&existingID)

// 	if err == sql.ErrNoRows {
// 		// Insert new company
// 		_, err := db.Exec(
// 			`INSERT INTO companies (name, groww_company_id, api_url, created_at, updated_at) 
// 			VALUES ($1, $2, $3, $4, $4)`,
// 			company.CompanyName, company.GrowwContractID, apiUrl, currentTime,
// 		)
// 		if err != nil {
// 			log.Printf("Error inserting %s: %v\n", company.CompanyName, err)
// 			return err
// 		}
// 		fmt.Printf("✅ Created: %s (%s)\n", company.CompanyName, company.GrowwContractID)
// 	} else if err != nil {
// 		// Other SQL errors
// 		log.Printf("Database error: %v\n", err)
// 		return err
// 	} else {
// 		// Update existing company (only update groww_company_id, api_url, and updated_at)
// 		_, err := db.Exec(
// 			`UPDATE companies 
// 			SET groww_company_id = $1, api_url = $2, updated_at = $3 
// 			WHERE name = $4`,
// 			company.GrowwContractID, apiUrl, currentTime, company.CompanyName,
// 		)
// 		if err != nil {
// 			log.Printf("Error updating %s: %v\n", company.CompanyName, err)
// 			return err
// 		}
// 		fmt.Printf("🔄 Updated: %s (%s)\n", company.CompanyName, company.GrowwContractID)
// 	}

// 	return nil
// }

// func updateAPIURL() {
// 	// Connect to PostgreSQL
// 	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
// 		DbHost, DbUser, DbPassword, DbName)

// 	db, err := sql.Open("postgres", psqlInfo)
// 	if err != nil {
// 		log.Fatal("❌ Database connection error:", err)
// 	}
// 	defer db.Close()

// 	// Check database connection
// 	if err := db.Ping(); err != nil {
// 		log.Fatal("❌ Cannot connect to database:", err)
// 	}
// 	fmt.Println("✅ Connected to PostgreSQL successfully!")

// 	// Fetch data from API and store in database
// 	page := 1
// 	for {
// 		companies, err := fetchCompaniesFromAPI(page)
// 		if err != nil {
// 			log.Fatal("❌ Error fetching data from API:", err)
// 		}

// 		if len(companies) == 0 {
// 			fmt.Println("✅ All records fetched successfully!")
// 			break
// 		}

// 		// Insert or update each company
// 		for _, company := range companies {
// 			if err := insertOrUpdateCompany(db, company); err != nil {
// 				log.Printf("❌ Error processing %s: %v\n", company.CompanyName, err)
// 			}
// 		}

// 		page++ // Next page
// 	}

// 	fmt.Println("✅ Data fetching and storing completed!")
// }

// func main() {
// 	updateAPIURL() // Call the updateAPIURL function to start the process
// }