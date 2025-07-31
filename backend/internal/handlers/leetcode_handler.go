package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// LeetCodeGraphQLRequest represents the GraphQL request structure
type LeetCodeGraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// LeetCodeProxyHandler handles proxying requests to LeetCode's GraphQL API
func LeetCodeProxyHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse the request
	var graphqlRequest LeetCodeGraphQLRequest
	if err := json.Unmarshal(body, &graphqlRequest); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Create a new request to the LeetCode GraphQL API
	leetcodeURL := "https://leetcode.com/graphql"
	proxyReq, err := http.NewRequest(http.MethodPost, leetcodeURL, bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}

	// Copy necessary headers
	proxyReq.Header.Set("Content-Type", "application/json")
	proxyReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	proxyReq.Header.Set("Referer", "https://leetcode.com/contest/")

	// Make the request to LeetCode
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Error fetching data from LeetCode", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response from LeetCode", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	// Write the response body
	w.Write(respBody)
}
