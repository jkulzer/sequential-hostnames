package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleGet).Methods("GET")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	srv.ListenAndServe()
}

func handleGet(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("User-Agent") != "SequentialHostname/1.0" {
		http.Error(w, "Bad Request: Invalid User-Agent header", http.StatusBadRequest)
		return
	}

	// Handle other requests normally.
	json := simplejson.New()
	hostname := getRandomFromList()

	if hostname == "" {
		http.NotFound(w, req)
		return
	}

	json.Set("hostname", hostname)

	payload, err := json.MarshalJSON()
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func getRandomFromList() (hostname string) {
	rand.Seed(time.Now().UnixNano())

	filePath := "hostnames"

	lines, err := readLines(filePath)
	if err != nil {
		fmt.Printf("Error reading lines from file: %v\n", err)
		return
	}

	if len(lines) == 0 {
		fmt.Println("The file is empty.")
		return
	}

	randomIndex := rand.Intn(len(lines))
	hostname = lines[randomIndex]

	// Remove the selected line from the slice.
	lines = append(lines[:randomIndex], lines[randomIndex+1:]...)

	// Write the updated lines back to the file.
	err = writeLines(filePath, lines)
	if err != nil {
		fmt.Printf("Error writing lines to file: %v\n", err)
	}

	return hostname
}

// readLines reads the lines from a file and returns them as a slice of strings.
func readLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// writeLines writes the lines to a file.
func writeLines(filePath string, lines []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}
