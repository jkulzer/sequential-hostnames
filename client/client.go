package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	// Define a command-line flag for the URL.
	urlFlag := flag.String("url", "", "URL to send the GET request")
	flag.Parse()

	// Check if the URL flag is provided.
	if *urlFlag == "" {
		fmt.Println("Please provide a URL using the -url flag.")
		os.Exit(1)
	}

	userAgent := "SequentialHostname/1.0"

	// Create a new HTTP client with the desired User-Agent.
	client := &http.Client{}

	// Create a new GET request with the provided URL.
	req, err := http.NewRequest("GET", *urlFlag, nil)
	if err != nil {
		fmt.Printf("Error creating GET request: %v\n", err)
		os.Exit(1)
	}

	// Set the hardcoded User-Agent header in the request.
	req.Header.Set("User-Agent", userAgent)

	// Send the GET request.
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending GET request: %v\n", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	// Check the response status code.
	if response.StatusCode != http.StatusOK {
		fmt.Printf("Error: Received status code %d\n", response.StatusCode)
		os.Exit(1)
	}

	// Read the response body.
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		os.Exit(1)
	}

	// Parse the JSON response.
	var responseData ResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Printf("Error parsing JSON response: %v\n", err)
		os.Exit(1)
	}

	hostname := responseData.Hostname

	switch os := runtime.GOOS; os {
	case "windows":
		changeHostnameWindows(hostname)
	case "linux":
		changeHostnameLinux(hostname)
	default:
		fmt.Println("Unsupported OS")
	}
}

type ResponseData struct {
	Hostname string `json:"hostname"`
}

func changeHostnameLinux(hostname string) error {

	cmd := exec.Command("hostnamectl", "hostname", hostname)

	// Run the Command
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func changeHostnameWindows(hostname string) error {
	// Format the PowerShell command with the desired hostname.
	powerShellCommand := fmt.Sprintf("Rename-Computer -NewName '%s' -Force -Restart", hostname)

	// Create the cmd object to run the PowerShell process.
	cmd := exec.Command("powershell", "-Command", powerShellCommand)

	// Run the PowerShell command.
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
