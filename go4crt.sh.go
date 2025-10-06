package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// Struct for crt.sh JSON
type CrtshResult struct {
	NameValue string `json:"name_value"`
}

// Interactive spinner
func spinner(stopChan chan struct{}) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	for {
		for _, frame := range frames {
			select {
			case <-stopChan:
				fmt.Print("\r") // clear spinner
				return
			default:
				fmt.Printf("\r%s Scanning... please wait", frame)
				time.Sleep(120 * time.Millisecond)
			}
		}
	}
}

func usage() {
	fmt.Println(`
╔════════════════════════════════════════════╗
║         CRT.SH Subdomain Finder Tool       ║
║                 by @4n_curze...            ║
╚════════════════════════════════════════════╝

Usage:
  crtsh -d <domain> -o <output_file>
  crtsh <domain> -o <output_file>

Example:
  go run main.go -d hackerone.com -o /home/kali/Desktop/hacker.txt

Flags:
  -d    Target domain name (e.g., example.com)
  -o    Output file path (required)
  -h    Show help and usage examples
`)
}

func main() {
	fmt.Println(`
            ▗▄                          ▗▖   
            ▟█            ▐▌            ▐▌   
 ▟█▟▌ ▟█▙  ▐▘█  ▟██▖ █▟█▌▐███      ▗▟██▖▐▙██▖
▐▛ ▜▌▐▛ ▜▌▗▛ █ ▐▛  ▘ █▘   ▐▌       ▐▙▄▖▘▐▛ ▐▌
▐▌ ▐▌▐▌ ▐▌▐███▌▐▌    █    ▐▌        ▀▀█▖▐▌ ▐▌
▝█▄█▌▝█▄█▘   █ ▝█▄▄▌ █    ▐▙▄   █  ▐▄▄▟▌▐▌ ▐▌
 ▞▀▐▌ ▝▀▘    ▀  ▝▀▀  ▀     ▀▀   ▀   ▀▀▀ ▝▘ ▝▘
 ▜█▛▘                                        	`+ "\n\n")

	// Flags
	flagDomain := flag.String("d", "", "Target domain (e.g., example.com)")
	flagOutput := flag.String("o", "", "Output file path (required)")
	flag.Usage = usage
	flag.Parse()

	// Support domain from -d or positional first arg
	domain := *flagDomain
	if domain == "" && flag.NArg() > 0 {
		domain = flag.Arg(0)
	}

	if domain == "" || *flagOutput == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Spinner setup
	stopChan := make(chan struct{})
	go spinner(stopChan)

	crtURL := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)
	resp, err := http.Get(crtURL)
	if err != nil {
		close(stopChan)
		log.Fatalf("\n❌ Error fetching data: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		close(stopChan)
		log.Fatalf("\n❌ Error reading response: %v", err)
	}

	// Stop spinner
	close(stopChan)
	fmt.Print("\r") // clear spinner output

	var crtshResults []CrtshResult
	if err := json.Unmarshal(body, &crtshResults); err != nil {
		log.Fatalf("❌ Error parsing JSON: %v", err)
	}

	subdomains := make(map[string]struct{})
	for _, result := range crtshResults {
		for _, line := range strings.Split(result.NameValue, "\n") {
			subdomain := strings.TrimSpace(line)
			subdomain = strings.TrimPrefix(subdomain, "*.")
			if subdomain != "" {
				subdomains[subdomain] = struct{}{}
			}
		}
	}

	uniqueSubdomains := make([]string, 0, len(subdomains))
	for sub := range subdomains {
		uniqueSubdomains = append(uniqueSubdomains, sub)
	}
	sort.Strings(uniqueSubdomains)

	if err := ioutil.WriteFile(*flagOutput, []byte(strings.Join(uniqueSubdomains, "\n")), 0644); err != nil {
		log.Fatalf("❌ Error writing file: %v", err)
	}

	fmt.Printf("✅ Scan completed successfully!\n")
	fmt.Printf("📁 Results saved in: %s\n", *flagOutput)
	fmt.Printf("🔢 Total subdomains found: %d\n", len(uniqueSubdomains))
}
