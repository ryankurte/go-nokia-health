package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/jessevdk/go-flags"

	"github.com/ryankurte/go-nokia-health/lib"
)

type Args struct {
	Mode         string `short:"m" long:"mode" description:"Util operating mode"`
	APIKey       string `short:"k" long:"api-key" description:"Health API key"`
	APISecret    string `short:"s" long:"api-secret" description:"Health API secret"`
	UserID       string `short:"u" long:"user-id" description:"Health user ID"`
	AccessToken  string `short:"t" long:"access-token" description:"Health user access token"`
	AccessSecret string `short:"a" long:"access-secret" description:"Health user access secret"`
	Port         uint32 `short:"p" long:"port" description:"Port for local interface binding"`
}

func getAccessTokens(h *nhealth.HealthAPI) (userid, accessToken, accessSecret string, err error) {
	requestToken, requestSecret, url, err := h.Request()
	if err != nil {
		return "", "", "", fmt.Errorf("error requesting authorization: %s", err)
	}

	fmt.Printf("Click the following link to authorize the application\n")
	fmt.Printf("%s\n", url)

	res := make(chan *http.Request)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		res <- r
		w.Write([]byte("<html><head></head><body onload=\"window.close();\"><h2>Authorization Complete</h2></body></html>"))
	})
	go http.ListenAndServe("localhost:9002", handlers.LoggingHandler(os.Stdout, http.DefaultServeMux))
	if err != nil {
		return "", "", "", fmt.Errorf("error requesting authorization: %s", err)
	}

	select {
	case <-time.After(30 * time.Second):
		return "", "", "", fmt.Errorf("timeout awaiting OAuth response")
	case r := <-res:
		userid, accessToken, accessSecret, err = h.Authorize(requestToken, requestSecret, r)
		if err != nil {
			return "", "", "", fmt.Errorf("error requesting access token: %s", err)
		}
	}

	return userid, accessToken, accessSecret, nil
}

func main() {
	fmt.Printf("Nokia Health API util\n")

	// Parse utility arguments
	args := Args{
		Mode:         "test",
		APIKey:       os.Getenv("NOKIA_API_KEY"),
		APISecret:    os.Getenv("NOKIA_API_SECRET"),
		AccessToken:  os.Getenv("NOKIA_ACCESS_TOKEN"),
		AccessSecret: os.Getenv("NOKIA_ACCESS_SECRET"),
		UserID:       os.Getenv("NOKIA_USER_ID"),
	}
	_, err := flags.Parse(&args)
	if err != nil {
		fmt.Printf("Error parsing arguments: %s\n", err)
		os.Exit(-1)
	}

	// Check consumer API key and secret are present
	if args.APIKey == "" || args.APISecret == "" {
		fmt.Printf("APIKey and APISecret arguments are required\n")
		os.Exit(-2)
	}

	// Create new HealthAPI connector
	h := nhealth.NewHealthAPI(args.APIKey, args.APISecret, "http://localhost:9002")

	userID, accessToken, accessSecret := "", "", ""

	// OAuth exchange if access token/secret are not specified
	if args.UserID == "" || args.AccessToken == "" || args.AccessSecret == "" {
		fmt.Println("No UserID, Access token or secret supplied, performing OAuth exchange...")
		userID, accessToken, accessSecret, err = getAccessTokens(&h)
		if err != nil {
			fmt.Printf("Error fetching access tokens: %s\n", err)
			os.Exit(-3)
		}
	} else {
		userID, accessToken, accessSecret = args.UserID, args.AccessToken, args.AccessSecret
	}

	fmt.Printf("User ID: %s Access token: %s Access secret: %s\n", userID, accessToken, accessSecret)

	userIDInt, _ := strconv.Atoi(userID)
	mq := nhealth.MeasureQuery{
		UserID:      uint32(userIDInt),
		MeasureType: nhealth.MeasureTypeWeight,
		StartDate:   uint32(time.Now().AddDate(0, 0, -7).Unix()),
		EndDate:     uint32(time.Now().Unix()),
	}
	res, err := h.GetMeasurement(accessToken, accessSecret, mq)
	if err != nil {
		fmt.Printf("Error fetching measurement data: %s\n", err)
		os.Exit(-4)
	}

	fmt.Printf("Measurements: %+v\n", res)
}
