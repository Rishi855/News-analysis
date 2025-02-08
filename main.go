package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

// Define struct for the API response
type Article struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	URL      string `json:"url"`
	PubDate  string `json:"pubDate"`
	Source   string `json:"source"`
	ImageURL string `json:"imageUrl"`
}

type Company struct {
	ID     string
	APIURL string
}

func main() {
	// Connect to the database
	connStr := "postgres://postgres:root@localhost/stocks?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get all companies with api_url
	rows, err := db.Query("SELECT id, api_url FROM companies")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over each company and fetch data
	for rows.Next() {
		var company Company
		if err := rows.Scan(&company.ID, &company.APIURL); err != nil {
			log.Fatal(err)
		}

		// Fetch data from API URL
		articles, err := fetchArticles(company.APIURL)
		if err != nil {
			log.Printf("Error fetching data for company %s: %v", company.ID, err)
			continue
		}

		// Insert articles into stock_articles table
		for _, article := range articles {
			// Check if the article already exists
			var exists bool
			err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM stock_articles WHERE id = $1)", article.ID).Scan(&exists)
			if err != nil {
				log.Printf("Error checking article existence for %s: %v", article.ID, err)
				continue
			}

			if !exists {
				// Insert article data
				_, err := db.Exec(`
					INSERT INTO stock_articles (id, title, summary, url, pub_date, source, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
					article.ID, article.Title, article.Summary, article.URL, article.PubDate, article.Source)
				if err != nil {
					log.Printf("Error inserting article %s: %v", article.ID, err)
				}
			}
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

// fetchArticles makes an HTTP request to fetch articles from the given API URL
func fetchArticles(apiURL string) ([]Article, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		Results []Article `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Convert PubDate to a formatted string
	for i := range result.Results {
		pubDate, err := time.Parse(time.RFC3339, result.Results[i].PubDate)
		if err == nil {
			result.Results[i].PubDate = pubDate.Format("2006-01-02 15:04:05")
		}
	}

	return result.Results, nil
}
