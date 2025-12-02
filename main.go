package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Result struct holds the output of a scraping job
type Result struct {
	URL        string
	StatusCode int
	Size       int64
	Error      error
}

// worker function processes jobs from the jobs channel and sends results to results channel
func worker(id int, jobs <-chan string, results chan<- Result, wg *sync.WaitGroup) {
	// Decrease WaitGroup counter when worker exits
	defer wg.Done()

	for url := range jobs {
		fmt.Printf("Worker %d started job: %s\n", id, url)

		// Start timer for demonstration
		start := time.Now()

		// Create HTTP client with timeout to prevent hanging
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		// Perform the HTTP request
		resp, err := client.Get(url)

		result := Result{URL: url}

		if err != nil {
			result.Error = err
		} else {
			result.StatusCode = resp.StatusCode
			// Read body to get Content Size (bytes)
			body, err := io.ReadAll(resp.Body)
			if err == nil {
				result.Size = int64(len(body))
			}
			resp.Body.Close()
		}

		// Send result to the results channel
		results <- result

		fmt.Printf("Worker %d finished job: %s (took %v)\n", id, url, time.Since(start))
	}
}

func main() {
	// List of URLs to scrape (hardcoded as per deliverables)
	urls := []string{
		"https://www.uottawa.ca",
		"https://www.google.com",
		"https://www.github.com",
		"https://golang.org",
		"https://www.stackoverflow.com",
		"https://www.reddit.com",
		"https://www.wikipedia.org",
		"https://www.microsoft.com",
	}

	// Configuration - using constants as requested
	const numWorkers = 5         // Fixed number of workers
	const jobChannelSize = 10    // Reasonable buffer size for job channel
	const resultChannelSize = 10 // Reasonable buffer size for result channel

	// Create channels with reasonable buffer sizes
	// jobs: to send URLs to workers
	// results: to collect Result structs from workers
	jobs := make(chan string, jobChannelSize)
	results := make(chan Result, resultChannelSize)

	// WaitGroup to track worker completion
	var wg sync.WaitGroup

	// 1. Start Workers
	// Launch 'numWorkers' goroutines using the constant
	fmt.Printf("Starting %d workers...\n", numWorkers)
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// 2. Send Jobs
	// Push all URLs into the jobs channel in a separate goroutine
	// This prevents blocking if job channel buffer is full
	go func() {
		fmt.Println("Sending jobs...")
		for _, url := range urls {
			jobs <- url
		}
		// Close jobs channel to signal to workers that no more jobs are coming
		close(jobs)
	}()

	// 3. Wait for workers to finish in a separate goroutine
	// This ensures we can process results while waiting
	done := make(chan bool)
	go func() {
		wg.Wait()
		close(results) // Close results channel when all workers are done
		done <- true
	}()

	// 4. Process Results
	// Read from results channel until it is closed
	fmt.Println("\n--- Results (in order of completion) ---")

	// Create a slice to store results for ordered output if needed
	var allResults []Result

	for res := range results {
		if res.Error != nil {
			fmt.Printf("[-] Error scraping %s: %s\n", res.URL, res.Error)
		} else {
			fmt.Printf("[+] URL: %s | Status: %d | Size: %d bytes\n",
				res.URL, res.StatusCode, res.Size)
		}
		allResults = append(allResults, res)
	}

	// Wait for the done signal
	<-done

	// Optional: Display summary
	fmt.Println("\n--- Summary ---")
	fmt.Printf("Total URLs processed: %d\n", len(urls))
	fmt.Printf("Successful requests: %d\n", countSuccessful(allResults))
	fmt.Printf("Failed requests: %d\n", countFailed(allResults))
	fmt.Println("All jobs completed.")
}

// Helper function to count successful requests
func countSuccessful(results []Result) int {
	count := 0
	for _, r := range results {
		if r.Error == nil {
			count++
		}
	}
	return count
}

// Helper function to count failed requests
func countFailed(results []Result) int {
	count := 0
	for _, r := range results {
		if r.Error != nil {
			count++
		}
	}
	return count
}
