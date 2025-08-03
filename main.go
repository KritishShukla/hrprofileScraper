package main

import (
	"io"
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type TavilyResult struct {
	Results []struct {
		URL string `json:"url"`
	} `json:"results"`
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found. Make sure TAVILY_API_KEY is set.")
	}

	apiKey := os.Getenv("TAVILY_API_KEY")
	if apiKey == "" {
		fmt.Println("❌ Error: TAVILY_API_KEY not found in environment variables.")
		fmt.Println("Please set your Tavily API key in a .env file or as an environment variable.")
		return
	}

	// Get company name from user
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the company name (or press Enter for 'Microsoft'): ")
	input, _ := reader.ReadString('\n')
	companyName := strings.TrimSpace(input)
	if companyName == "" {
		companyName = "Microsoft"
	}

	fmt.Printf("\n🔍 Searching for HR profiles at %s...\n", companyName)

// Define the search query
query := fmt.Sprintf("('Human Resources' OR 'Talent Acquisition' OR 'Recruiter') '%s' site:linkedin.com/in/", companyName)
fmt.Printf("Executing search with query: %s\n\n", query)

// Create request body
reqBody := fmt.Sprintf(`{"query": "%s", "max_results": 20}`, query)
req, err := http.NewRequest("POST", "https://api.tavily.com/search", strings.NewReader(reqBody))
if err != nil {
	fmt.Println("❌ Error creating request:", err)
	return
}
req.Header.Set("Authorization", "Bearer "+apiKey)
req.Header.Set("Content-Type", "application/json")

// Send request
client := &http.Client{}
resp, err := client.Do(req)
if err != nil {
	fmt.Println("❌ Error sending request:", err)
	return
}
defer resp.Body.Close()

// Read full response body
body, err := io.ReadAll(resp.Body)
if err != nil {
	fmt.Println("❌ Error reading response body:", err)
	return
}

// Decode JSON response
var result TavilyResult
if err := json.Unmarshal(body, &result); err != nil {
	fmt.Println("❌ Error parsing JSON response:", err)
	return
}

	// Extract LinkedIn profile URLs
	linkedinSet := make(map[string]bool)
	for _, res := range result.Results {
		if strings.Contains(res.URL, "linkedin.com/in/") {
			cleanURL := strings.Split(res.URL, "?")[0]
			linkedinSet[cleanURL] = true
		}
	}

	// Print result
	if len(linkedinSet) > 0 {
		fmt.Printf("\n✅ Found %d unique LinkedIn profiles:\n", len(linkedinSet))
		i := 1
		for url := range linkedinSet {
			fmt.Printf("%d. %s\n", i, url)
			i++
		}
	} else {
		fmt.Println("No LinkedIn profiles found.")
		fmt.Println("This might be due to:")
		fmt.Println("- No HR professionals found at this company")
		fmt.Println("- LinkedIn blocking results")
		fmt.Println("- Company name not matching")
	}
}

