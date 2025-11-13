package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create [url]",
	Short: "Shorten a URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shortenURL(args[0])
	},
}

func shortenURL(longURL string) {
	apiURL := os.Getenv("SHORTLY_API")
	if apiURL == "" {
		apiURL = "https://go-shortly.onrender.com" // 👈 Your live URL
	}

	payload, _ := json.Marshal(map[string]string{"url": longURL})
	req, _ := http.NewRequest("POST", apiURL+"/shorten", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var result struct {
		ShortURL string `json:"short_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid response\n")
		os.Exit(1)
	}

	if resp.StatusCode != 201 {
		fmt.Fprintf(os.Stderr, "API error\n")
		os.Exit(1)
	}

	fmt.Println(result.ShortURL)
}
