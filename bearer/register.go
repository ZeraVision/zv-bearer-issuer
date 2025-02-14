package bearer

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type RegisterBearerTokenResponse struct {
	Token      string `json:"token"`
	ValidUntil int64  `json:"validUntil"`
	MaxRPS     int16  `json:"maxRequestsPerSecond"`
	MaxBPS     int16  `json:"maxBurstPerSecond"`
	MaxPerDay  int    `json:"maxRequestsPerDay"`
}

// Struct to represent the error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// Function to extract the error message from the JSON response
func extractErrorMessage(body []byte) (string, error) {
	var errorResponse ErrorResponse
	err := json.Unmarshal(body, &errorResponse)
	if err != nil {
		return "", errors.New("error unmarshaling error response")
	}
	return errorResponse.Error, nil
}

func Register() (RegisterBearerTokenResponse, error) {
	// Pull your api key data to authorize the bearer
	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")

	//* Determine the request parameters for your bearer (Modify values as Needed) */
	var maxRequestPerSecond, maxBurstPerSecond int16 // the rate at which the bearer bucket refills per second || the max amount it can fill | [token bucket algorithm]
	var maxRequestPerDay, expiresInSeconds int       // resets to 0 one time per day || expires <set> seconds after current time. Less time is generally more secure and better practice to mitigate abuse of your associated api key.
	var issuer, subject, audience *string            // optional data, many cases will leave this blank

	//* Sample values
	maxRequestPerSecond = 1
	maxBurstPerSecond = 5
	maxRequestPerDay = 1000
	expiresInSeconds = 86400

	// Construct the query parameters
	params := url.Values{}
	params.Set("requestType", "createBearer")
	params.Set("apiKey", apiKey)
	params.Set("apiSecret", apiSecret)
	params.Set("maxRequestPerSecond", fmt.Sprintf("%d", maxRequestPerSecond))
	params.Set("maxBurstPerSecond", fmt.Sprintf("%d", maxBurstPerSecond))
	params.Set("maxRequestPerDay", fmt.Sprintf("%d", maxRequestPerDay))
	params.Set("expiresInSeconds", fmt.Sprintf("%d", expiresInSeconds))

	// Add optional parameters only if they are set
	if issuer != nil {
		params.Set("issuer", *issuer)
	}
	if subject != nil {
		params.Set("subject", *subject)
	}
	if audience != nil {
		params.Set("audience", *audience)
	}

	// API URL
	apiURL := os.Getenv("INDEXER_URL") + "/store?" + params.Encode()

	// Create HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer([]byte{}))
	if err != nil {
		return RegisterBearerTokenResponse{}, errors.New("error creating request")
	}

	// Set headers
	req.Header.Set("Target", "apiAccess")

	// Execute the request
	//client := &http.Client{}

	// Create HTTP client with TLS configuration to skip certificate verification //! For testing
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return RegisterBearerTokenResponse{}, errors.New("error making request")
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RegisterBearerTokenResponse{}, errors.New("error reading response")
	} else if resp.StatusCode != http.StatusOK {
		// Extract the error message from the response body
		errorMessage, err := extractErrorMessage(body)
		if err != nil {
			return RegisterBearerTokenResponse{}, err
		}
		return RegisterBearerTokenResponse{}, errors.New(errorMessage)
	}

	// Unmarshal the response body into the struct
	var response RegisterBearerTokenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return RegisterBearerTokenResponse{}, errors.New("error unmarshaling response")
	}

	return response, nil
}
