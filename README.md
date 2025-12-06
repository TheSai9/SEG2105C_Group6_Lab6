# SEG2105 Lab 6 - Multithreading (Concurrent Web Scraper)

This project implements a **concurrent web scraper** in Go using **goroutines**, **channels**, and a **worker pool**.  
It processes multiple URLs concurrently, retrieves their HTTP status codes and content sizes, and reports all results in a structured format.

This lab demonstrates:
- Basic concurrency concepts  
- Launching and coordinating multiple goroutines  
- Using channels for communication  
- Implementing a worker pool  
- Waiting for all workers to finish cleanly  

---

## Features

✔️ Worker pool with fixed number of goroutines  
✔️ Jobs channel for input URLs  
✔️ Results channel for structured output  
✔️ Timeout-enabled HTTP client  
✔️ Error handling (timeouts, invalid URLs, network failures)  
✔️ Summary of successes and failures  
✔️ Worker logs showing job start and finish  

---

## How It Works

1. The **main goroutine** creates job and result channels.  
2. A fixed number of **worker goroutines** are launched.  
3. Workers receive URLs from the jobs channel.  
4. Each worker:
   - Sends an HTTP GET request  
   - Records status code and content size  
   - Reports the result through the results channel  
5. Main goroutine listens for results and prints them.  
6. Summary is displayed at the end once all workers complete.

---

## Project Structure
```
├── main.go # Main program with worker pool & scraper logic
├── README.md # Documentation (this file)
```
---

## Prerequisites

Go must be installed (Go 1.20+ recommended).

Check with:
```
go version
```

---

## Running the Program

Clone the repository:
```
git clone https://github.com/TheSai9/SEG2105C_Group6_Lab6.git
cd <your-repo>
```

Run the web scraper:

```
go run main.go
```

---

## Sample Output

```
Starting 5 workers...
Sending jobs...

--- Results (in order of completion) ---
Worker 1 started job: https://www.google.com
Worker 3 started job: https://www.github.com
Worker 2 started job: https://www.uottawa.ca
Worker 5 started job: https://www.reddit.com
Worker 4 started job: https://www.microsoft.com

[+] URL: https://www.google.com | Status: 200 | Size: 15532 bytes
[+] URL: https://www.github.com | Status: 200 | Size: 89244 bytes
[-] Error scraping https://www.reddit.com: 403 Forbidden

--- Summary ---
Total URLs processed: 8
Successful requests: 6
Failed requests: 2
All jobs completed.
```

---

## Demo Video



https://github.com/user-attachments/assets/c0ffe69b-5950-49dc-a5a1-9c197f66d6b4



