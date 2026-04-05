package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type Result struct {
	Proxy    string
	Protocol string
	Alive    bool
	Elapsed  time.Duration
}

func main() {
	workerCount := 50
	proxyFile := "proxy.txt"
	activeFile := "active.txt"
	timeout := 5 * time.Second
	targetURL := "http://httpbin.org/ip"

	proxies, err := readLines(proxyFile)
	if err != nil {
		fmt.Printf("Gagal membaca %s: %v\n", proxyFile, err)
		fmt.Println("Pastikan file proxy.txt ada dan berisi daftar IP:Port.")
		return
	}

	jobs := make(chan string, len(proxies))
	results := make(chan Result, len(proxies)*2) // *2 karena cek HTTP dan SOCKS5

	var wg sync.WaitGroup

	// Start Worker Pool
	for w := 1; w <= workerCount; w++ {
		wg.Add(1)
		go worker(&wg, jobs, results, targetURL, timeout)
	}

	// Send Jobs
	for _, p := range proxies {
		jobs <- p
	}
	close(jobs)

	// Result collector
	go func() {
		wg.Wait()
		close(results)
	}()

	// Process Results
	outputFile, err := os.Create(activeFile)
	if err != nil {
		fmt.Printf("Gagal membuat %s: %v\n", activeFile, err)
		return
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	aliveCount := 0
	deadCount := 0

	for res := range results {
		if res.Alive {
			fmt.Printf("[ALIVE] [%s] %s - Respon: %v\n", res.Protocol, res.Proxy, res.Elapsed)
			writer.WriteString(fmt.Sprintf("%s://%s\n", res.Protocol, res.Proxy))
			aliveCount++
		} else {
			// Optional: print dead proxies too
			// fmt.Printf("[DEAD]  [%s] %s\n", res.Protocol, res.Proxy)
			deadCount++
		}
	}

	fmt.Printf("\n--- Selesai ---\n")
	fmt.Printf("Total Proxies: %d\n", len(proxies))
	fmt.Printf("Alive: %d (disimpan ke %s)\n", aliveCount, activeFile)
	fmt.Printf("Dead: %d\n", deadCount)
}

func worker(wg *sync.WaitGroup, jobs <-chan string, results chan<- Result, target string, timeout time.Duration) {
	defer wg.Done()
	for proxyAddr := range jobs {
		// Cek HTTP
		check(proxyAddr, "http", target, timeout, results)
		// Cek SOCKS5
		check(proxyAddr, "socks5", target, timeout, results)
	}
}

func check(proxyAddr, protocol, target string, timeout time.Duration, results chan<- Result) {
	proxyURL, _ := url.Parse(fmt.Sprintf("%s://%s", protocol, proxyAddr))
	
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			DialContext: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		},
		Timeout: timeout,
	}

	start := time.Now()
	resp, err := client.Get(target)
	elapsed := time.Since(start)

	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			results <- Result{Proxy: proxyAddr, Protocol: protocol, Alive: true, Elapsed: elapsed}
			return
		}
	}
	results <- Result{Proxy: proxyAddr, Protocol: protocol, Alive: false}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
