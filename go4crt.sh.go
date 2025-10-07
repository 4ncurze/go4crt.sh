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

// 🌸 Each domain has a story — here begins our quest to unveil them.
type CrtshResult struct {
	NameValue string `json:"name_value"`
}

// 🌟 A little spinner of life — turning time into motion, progress into joy.
func fancySpinner(stopChan chan struct{}) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	steps := []string{
		"🔍 Searching crt.sh logs",
		"📡 Gathering subdomains",
		"📂 Cleaning duplicates",
	}

	frameIndex := 0
	stepIndex := 0
	ticker := time.NewTicker(220 * time.Millisecond)
	stepTicker := time.NewTicker(2 * time.Second) // change message every 2s

	for {
		select {
		case <-stopChan:
			fmt.Print("\r\033[K") // clear line
			ticker.Stop()
			stepTicker.Stop()
			return
		case <-ticker.C:
			fmt.Printf("\r%s %s...", frames[frameIndex], steps[stepIndex])
			frameIndex = (frameIndex + 1) % len(frames)
		case <-stepTicker.C:
			stepIndex = (stepIndex + 1) % len(steps)
		}
	}
}

// 📖 A friendly guide for our traveler — how to use the map of go4crt.sh.
func usage() {
	fmt.Println(`
╔════════════════════════════════════════════╗
║      go4crt.sh Subdomain Finder Tool       ║
║                 by @4n_curze...            ║
╚════════════════════════════════════════════╝

Usage:
  go4crt.sh -d <domain> -o <output_file_Path_with_filename>

Example:
  go4crt.sh -d example.com -o /home/kali/target.txt

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

	// 🎯 Collect the clues before we begin our exploration.
	flagDomain := flag.String("d", "", "Target domain (e.g., example.com)")
	flagOutput := flag.String("o", "", "Output file path (required)")
	flag.Usage = usage
	flag.Parse()

	// 🌐 Allow both flags and direct input — flexibility for every explorer.
	domain := *flagDomain
	if domain == "" && flag.NArg() > 0 {
		domain = flag.Arg(0)
	}

	// 🚨 Don’t start the journey without your map and destination.
	if domain == "" || *flagOutput == "" {
		flag.Usage()
		os.Exit(1)
	}

	// 🌀 Set the spinner in motion — the journey begins.
	stopChan := make(chan struct{})
	go fancySpinner(stopChan)

	// 🌍 Dive into crt.sh — uncover hidden domains beneath the surface.
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

	// 🛑 Time to rest — our spinner’s dance is complete.
	close(stopChan)
	fmt.Print("\r\033[K") // clear spinner line
	
	// 📦 Transforming the raw data into meaningful insights.
	var crtshResults []CrtshResult
	if err := json.Unmarshal(body, &crtshResults); err != nil {
		log.Fatalf("❌ Error parsing JSON: %v", err)
	}

	// 💎 Extract unique gems (subdomains) from the trove.
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

	// 🌿 Sort and organize — clarity brings peace.
	uniqueSubdomains := make([]string, 0, len(subdomains))
	for sub := range subdomains {
		uniqueSubdomains = append(uniqueSubdomains, sub)
	}
	sort.Strings(uniqueSubdomains)

	// 🖋️ Seal the findings in a scroll — safe for future discoveries.
	if err := ioutil.WriteFile(*flagOutput, []byte(strings.Join(uniqueSubdomains, "\n")), 0644); err != nil {
		log.Fatalf("❌ Error writing file: %v", err)
	}

	// 🎉 The exploration concludes — knowledge earned, not just found.
	fmt.Printf("✅ Scan completed successfully!\n")
	fmt.Printf("📁 Results saved in: %s\n", *flagOutput)
	fmt.Printf("🔢 Total subdomains found: %d\n", len(uniqueSubdomains))
}
