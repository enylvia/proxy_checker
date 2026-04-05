package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

// ANSI Color Codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
)

type Result struct {
	Proxy    string
	Protocol string
	Alive    bool
	Elapsed  time.Duration
	Origin   string
}

func main() {
	// Command Line Flags
	workerCount := flag.Int("w", 50, "Jumlah worker paralel")
	proxyFile := flag.String("f", "proxy.txt", "File sumber list proxy")
	timeoutSec := flag.Int("t", 10, "Timeout dalam detik")
	flag.Parse()

	activeFile := "active.txt"
	timeout := time.Duration(*timeoutSec) * time.Second
	targetURL := "http://httpbin.org/ip"

	proxies, err := readLines(*proxyFile)
	if err != nil {
		fmt.Printf("%s[ERROR]%s Gagal membaca %s: %v\n", ColorRed, ColorReset, *proxyFile, err)
		return
	}

	jobs := make(chan string, len(proxies))
	results := make(chan Result, len(proxies)*2)

	var wg sync.WaitGroup

	// Start Worker Pool
	for w := 1; w <= *workerCount; w++ {
		wg.Add(1)
		go worker(&wg, jobs, results, targetURL, timeout)
	}

	// Send Jobs
	go func() {
		for _, p := range proxies {
			jobs <- p
		}
		close(jobs)
	}()

	// Result collector
	go func() {
		wg.Wait()
		close(results)
	}()

	// Output file
	outputFile, err := os.Create(activeFile)
	if err != nil {
		fmt.Printf("%s[ERROR]%s Gagal membuat %s\n", ColorRed, ColorReset, activeFile)
		return
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	aliveCount := 0
	deadCount := 0

	fmt.Printf("%s[INFO]%s Memulai pengecekan %d proxy dengan %d worker...\n\n", ColorCyan, ColorReset, len(proxies), *workerCount)

	for res := range results {
		if res.Alive {
			fmt.Printf("%s[ALIVE]%s [%-6s] %-20s | Respon: %-8v | IP: %s\n", ColorGreen, ColorReset, res.Protocol, res.Proxy, res.Elapsed, res.Origin)
			writer.WriteString(fmt.Sprintf("%s://%s\n", res.Protocol, res.Proxy))
			aliveCount++
		} else {
			fmt.Printf("%s[DEAD ]%s [%-6s] %-20s\n", ColorRed, ColorReset, res.Protocol, res.Proxy)
			deadCount++
		}
	}

	fmt.Printf("\n%s--- RINGKASAN ---%s\n", ColorYellow, ColorReset)
	fmt.Printf("Total Proxies : %d\n", len(proxies))
	fmt.Printf("Alive         : %s%d%s (tersimpan di %s)\n", ColorGreen, aliveCount, ColorReset, activeFile)
	fmt.Printf("Dead          : %s%d%s\n", ColorRed, deadCount, ColorReset)
}

func worker(wg *sync.WaitGroup, jobs <-chan string, results chan<- Result, target string, timeout time.Duration) {
	defer wg.Done()
	for proxyAddr := range jobs {
		check(proxyAddr, "http", target, timeout, results)
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
	if err != nil {
		results <- Result{Proxy: proxyAddr, Protocol: protocol, Alive: false}
		return
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)

	if resp.StatusCode == 200 {
		var result struct {
			Origin string `json:"origin"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
			results <- Result{Proxy: proxyAddr, Protocol: protocol, Alive: true, Elapsed: elapsed, Origin: result.Origin}
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
		if text := scanner.Text(); text != "" {
			lines = append(lines, text)
		}
	}
	return lines, scanner.Err()
}
