package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"gopkg.in/robfig/cron.v2"
)

const (
	apiKey = "" // Replace with your actual NewsAPI key
)

type Article struct {
	Domain      string `json:"domain"`
	Title       string `json:"title"`
	PublishedAt string `json:"published_at"`
	Slug        string `json:"slug"`
}

type NewsResponse struct {
	Results []Article `json:"results"`
}

type RequestBody struct {
	Content string `json:"content"`
	// Entities    []byte `json:"entities"`
	// Permissions []byte `json:"permissions"`
}

func debankPost(
	title string,
	published_at string,
	domain string,
	slug string,
) {
	requestURL := "https://api.debank.com/article/add"

	parsedTime, err := time.Parse(time.RFC3339, published_at)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}

	// Format the parsed time into a meaningful string format
	formattedDate := parsedTime.Format("Monday, January 2, 2006 at 3:04 PM MST")

	url := "https://" + domain + "/" + strings.ToLower(slug)

	// Create the request body
	requestBody := RequestBody{
		Content: title + "\n" + "Published at: " + formattedDate + "\n" + "Link: " + url,
		// Entities:    []byte("{}"),
		// Permissions: []byte("{}"),
	}

	// Convert the request body to JSON
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal("Error encoding JSON:", err)
	}

	fmt.Println(string(requestBodyJSON))

	// Create a new HTTP POST request with the JSON body
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	// Set the headers
	req.Header.Set("authority", "api.debank.com")
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8,vi;q=0.7")
	req.Header.Set("account", "{\"random_at\":1691937042,\"random_id\":\"\",\"session_id\":\"\",\"user_addr\":\"0x1d8197e9c63b1cbf16ea2896f3eb4241c6a29347\",\"wallet_type\":\"metamask\",\"is_verified\":true}")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("origin", "https://debank.com")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("referer", "https://debank.com/")
	req.Header.Set("sec-ch-ua", "\"Not/A)Brand\";v=\"99\", \"Google Chrome\";v=\"115\", \"Chromium\";v=\"115\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("source", "web")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	req.Header.Set("x-api-nonce", "")
	req.Header.Set("x-api-sign", "")
	req.Header.Set("x-api-ts", "1692442082")
	req.Header.Set("x-api-ver", "v2")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	// print the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}
	fmt.Println(string(body))
}

func getDataAndPost() {
	// Construct the request URL
	requestURL := "https://cryptopanic.com/api/v1/posts/" +
		"?auth_token=" + apiKey +
		"&filter=hot"

		// Send the GET request
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Parse the response JSON
	var newsResponse NewsResponse
	if err := json.NewDecoder(resp.Body).Decode(&newsResponse); err != nil {
		log.Fatal("Error decoding response:", err)
	}

	if len(
		newsResponse.Results,
	) > 0 {
		fmt.Println("Title:", newsResponse.Results[0].Slug)
		debankPost(
			newsResponse.Results[0].Title,
			newsResponse.Results[0].PublishedAt,
			newsResponse.Results[0].Domain,
			newsResponse.Results[0].Slug,
		)
	} else {
		fmt.Println("Error retrieving news")
	}
}

func main() {
	c := cron.New()

	cronExpression := "1 */4 * * *"

	c.AddFunc(cronExpression, getDataAndPost)

	c.Start()

	select {}
}
