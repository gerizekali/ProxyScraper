package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Version = "4.1.0-Global"
)

const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[97m"
	ColorBold    = "\033[1m"
	
	NeonPink = "\033[38;5;198m"
	NeonCyan = "\033[38;5;51m"
	NeonLime = "\033[38;5;46m"
	DarkGray = "\033[90m"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Firefox/123.0",
}

var defaultSources = []string{
	"https://api.proxyscrape.com/v2/?request=getproxies&protocol=http&timeout=10000&country=all&ssl=all&anonymity=all",
	"https://api.proxyscrape.com/v2/?request=getproxies&protocol=socks4&timeout=10000&country=all",
	"https://api.proxyscrape.com/v2/?request=getproxies&protocol=socks5&timeout=10000&country=all",
	"https://www.proxy-list.download/api/v1/get?type=http",
	"https://www.proxy-list.download/api/v1/get?type=https",
	
	"https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/http/data.txt",
	"https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/https/data.txt",
	"https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/socks4/data.txt",
	"https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/socks5/data.txt",

	"https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/http.txt",
	"https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/https.txt",
	"https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/socks4.txt",
	"https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/socks5.txt",

	"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/http.txt",
	"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/socks4.txt",
	"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/socks5.txt",

	"https://raw.githubusercontent.com/zloi-user/hideip.me/main/http.txt",
	"https://raw.githubusercontent.com/zloi-user/hideip.me/main/https.txt",
	"https://raw.githubusercontent.com/zloi-user/hideip.me/main/socks4.txt",
	"https://raw.githubusercontent.com/zloi-user/hideip.me/main/socks5.txt",

	"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/http.txt",
	"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/socks4.txt",
	"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/socks5.txt",

	"https://raw.githubusercontent.com/prxchk/proxy-list/main/http.txt",
	"https://raw.githubusercontent.com/mmpx12/proxy-list/master/https.txt",
	"https://raw.githubusercontent.com/roosterkid/openproxylist/master/HTTPS_RAW.txt",
	"https://raw.githubusercontent.com/sunny9577/proxy-scraper/master/proxies.txt",
	"https://raw.githubusercontent.com/shiftytr/proxy-list/master/proxy.txt",
	"https://raw.githubusercontent.com/clarketm/proxy-list/master/proxy-list-raw.txt",
	"https://raw.githubusercontent.com/opsxcq/proxy-list/master/list.txt",
	"https://raw.githubusercontent.com/almroot/proxylist/master/list.txt",
	"https://raw.githubusercontent.com/rdavydov/proxy-list/main/proxies/http.txt",
	"https://raw.githubusercontent.com/rdavydov/proxy-list/main/proxies/socks5.txt",
	"https://raw.githubusercontent.com/hookzof/socks5_list/master/proxy.txt",
	"https://raw.githubusercontent.com/saisuiu/Lionkings-Http-Proxys-Proxies/main/free.txt",
	"https://raw.githubusercontent.com/Anonym0usWork1220/Free-Proxies/main/proxy_files/http_proxies.txt",
	"https://raw.githubusercontent.com/officialputuid/KangProxy/KangProxy/http/http.txt",
	"https://raw.githubusercontent.com/officialputuid/KangProxy/KangProxy/https/https.txt",
}

var commonPorts = []int{8080, 3128, 80, 8000, 8888, 9999, 1080, 1081, 5000}

type Config struct {
	TargetURL   string
	Timeout     time.Duration
	Concurrency int
	OutputFile  string
	RawFile     string
	RandomCount int
}

type ProxyResult struct {
	Address   string
	Latency   time.Duration
	Status    bool
	ErrorMsg  string
}

type Stats struct {
	Total     int64
	Checked   int64
	Working   int64
	Failed    int64
	StartTime time.Time
}

var stats Stats
var currentConfig Config

func main() {
	currentConfig = Config{
		TargetURL:   "http://www.google.com",
		Timeout:     8 * time.Second,
		Concurrency: 500,
		OutputFile:  "proxies_detailed.txt",
		RawFile:     "proxies_raw.txt",
		RandomCount: 0,
	}

	for {
		clearScreen()
		printBanner()
		printMenu()
		
		choice := readInput(NeonCyan + " ┌──(Your Selection)\n └─> " + ColorReset)

		switch choice {
		case "1":
			startScanProcess()
			pause()
		case "2":
			settingsMenu()
		case "3":
			clearScreen()
			fmt.Println(NeonPink + "\n See you later space cowboy... " + ColorReset)
			os.Exit(0)
		default:
			fmt.Println(ColorRed + " [!] Invalid selection!" + ColorReset)
			time.Sleep(1 * time.Second)
		}
	}
}

func settingsMenu() {
	for {
		clearScreen()
		fmt.Println(NeonPink + "\n  ╔════════════════════ SETTINGS ════════════════════╗" + ColorReset)
		fmt.Printf("  ║ 1. Target URL      : %-27s ║\n", currentConfig.TargetURL)
		fmt.Printf("  ║ 2. Thread Count    : %-27d ║\n", currentConfig.Concurrency)
		fmt.Printf("  ║ 3. Timeout (sec)   : %-27.0f ║\n", currentConfig.Timeout.Seconds())
		fmt.Printf("  ║ 4. Random Generate : %-27d ║\n", currentConfig.RandomCount)
		fmt.Println("  ║ 9. Back to Main Menu                             ║")
		fmt.Println("  ╚══════════════════════════════════════════════════╝" + ColorReset)

		choice := readInput(NeonCyan + " └─> " + ColorReset)

		switch choice {
		case "1":
			val := readInput(" New Target URL: ")
			if val != "" { currentConfig.TargetURL = val }
		case "2":
			val := readInput(" New Thread Count: ")
			if n, err := strconv.Atoi(val); err == nil { currentConfig.Concurrency = n }
		case "3":
			val := readInput(" New Timeout (sec): ")
			if n, err := strconv.Atoi(val); err == nil { currentConfig.Timeout = time.Duration(n) * time.Second }
		case "4":
			val := readInput(" Random IP Count (0 = Disabled): ")
			if n, err := strconv.Atoi(val); err == nil { currentConfig.RandomCount = n }
		case "9":
			return
		}
	}
}

func startScanProcess() {
	clearScreen()
	fmt.Println(NeonCyan + "\n [🚀] Initiating Scan Protocol..." + ColorReset)
	
	fmt.Printf(ColorYellow+" [1/4] Fetching data from %d sources...%s\n", len(defaultSources), ColorReset)
	scrapedProxies := scrapeProxies(defaultSources)
	
	var generatedProxies []string
	if currentConfig.RandomCount > 0 {
		fmt.Printf(ColorBlue+" [INFO] Generating %d random IPs...%s\n", currentConfig.RandomCount, ColorReset)
		generatedProxies = generateRandomProxies(currentConfig.RandomCount)
	}

	allProxies := append(scrapedProxies, generatedProxies...)
	uniqueProxies := removeDuplicates(allProxies)

	stats = Stats{Total: int64(len(uniqueProxies)), StartTime: time.Now()}

	fmt.Printf(NeonLime+" [OK] Total %d unique targets loaded.\n\n"+ColorReset, stats.Total)

	goodProxies := startWorkerPool(uniqueProxies, currentConfig)

	fmt.Printf("\n"+ColorYellow+" [3/4] Analyzing speed and sorting...%s\n", ColorReset)
	sort.Slice(goodProxies, func(i, j int) bool {
		return goodProxies[i].Latency < goodProxies[j].Latency
	})

	fmt.Printf(ColorYellow+" [4/4] Writing results to disk...%s\n", ColorReset)
	saveDualProxies(goodProxies)
	
	printSummary(goodProxies)
}

func startWorkerPool(proxies []string, cfg Config) []ProxyResult {
	jobs := make(chan string, len(proxies))
	results := make(chan ProxyResult, len(proxies))
	var wg sync.WaitGroup

	stopUI := make(chan bool)
	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-stopUI:
				return
			case <-ticker.C:
				printLiveStats()
			}
		}
	}()

	for w := 0; w < cfg.Concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range jobs {
				res := checkProxy(p, cfg)
				atomic.AddInt64(&stats.Checked, 1)
				if res.Status {
					atomic.AddInt64(&stats.Working, 1)
				} else {
					atomic.AddInt64(&stats.Failed, 1)
				}
				results <- res
			}
		}()
	}

	for _, p := range proxies { jobs <- p }
	close(jobs)
	wg.Wait()
	close(results)
	stopUI <- true

	var working []ProxyResult
	for res := range results {
		if res.Status { working = append(working, res) }
	}
	return working
}

func checkProxy(proxyAddr string, cfg Config) ProxyResult {
	if !strings.HasPrefix(proxyAddr, "http") {
		proxyAddr = "http://" + proxyAddr
	}
	proxyURL, err := url.Parse(proxyAddr)
	if err != nil { return ProxyResult{Address: proxyAddr, Status: false} }

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
		DialContext: (&net.Dialer{Timeout: cfg.Timeout, KeepAlive: 0}).DialContext,
	}
	client := &http.Client{Transport: transport, Timeout: cfg.Timeout}

	req, _ := http.NewRequest("GET", cfg.TargetURL, nil)
	req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
	req.Header.Set("Connection", "close")

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil { return ProxyResult{Address: proxyAddr, Status: false} }
	defer resp.Body.Close()
	
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return ProxyResult{
			Address: strings.TrimPrefix(proxyAddr, "http://"),
			Latency: time.Since(start),
			Status:  true,
		}
	}
	return ProxyResult{Address: proxyAddr, Status: false}
}

func saveDualProxies(proxies []ProxyResult) {
	fDetail, _ := os.Create(currentConfig.OutputFile)
	defer fDetail.Close()
	wDetail := bufio.NewWriter(fDetail)
	
	wDetail.WriteString("### Generated by PROXY FINDER ###\n")
	wDetail.WriteString("### IP:PORT | LATENCY ###\n")
	for _, p := range proxies {
		wDetail.WriteString(fmt.Sprintf("%s | %s\n", p.Address, p.Latency))
	}
	wDetail.Flush()

	fRaw, _ := os.Create(currentConfig.RawFile)
	defer fRaw.Close()
	wRaw := bufio.NewWriter(fRaw)
	
	for _, p := range proxies {
		wRaw.WriteString(fmt.Sprintf("%s\n", p.Address))
	}
	wRaw.Flush()
}


func scrapeProxies(sources []string) []string {
	var all []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	re := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d+`)

	for _, s := range sources {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			c := &http.Client{Timeout: 15 * time.Second}
			if resp, err := c.Get(url); err == nil {
				defer resp.Body.Close()
				if b, err := io.ReadAll(resp.Body); err == nil {
					m := re.FindAllString(string(b), -1)
					mu.Lock()
					all = append(all, m...)
					mu.Unlock()
				}
			}
		}(s)
	}
	wg.Wait()
	return all
}

func generateRandomProxies(count int) []string {
	p := make([]string, count)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < count; i++ {
		p[i] = fmt.Sprintf("%d.%d.%d.%d:%d", rand.Intn(223)+1, rand.Intn(256), rand.Intn(256), rand.Intn(256), commonPorts[rand.Intn(len(commonPorts))])
	}
	return p
}

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}
	for _, v := range elements {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}

func printBanner() {
	banner := `
` + NeonPink + `    ____  ____  ____  _  ____  __    __________  _____    ____  __________ 
   / __ \/ __ \/ __ \| |/ /\ \/ /   / ____/  _/ / /| |   / __ \/ ____/ __ \
  / /_/ / /_/ / / / /   /  \  /   / /_   / / / / | |  / / / / __/ / /_/ /
 / ____/ _, _/ /_/ /   |   / /   / __/ _/ / / /  | | / /_/ / /___/ _, _/ 
/_/   /_/ |_|\____/_/|_|  /_/   /_/   /___/_/   |_|/_____/_____/_/ |_|  
` + ColorReset + `
													 
` + NeonCyan + `       >> PROXY FINDER <<` + ColorReset + `
` + DarkGray + `       -----------------------------------------` + ColorReset + `
`
	fmt.Println(banner)
}

func printMenu() {
	fmt.Println(NeonLime + "  ╔══════════════════════════════════════════════════╗" + ColorReset)
	fmt.Println(NeonLime + "  ║ " + ColorWhite + "1." + NeonCyan + " Start Scan (Start Protocol)                    " + NeonLime + "║" + ColorReset)
	fmt.Println(NeonLime + "  ║ " + ColorWhite + "2." + NeonCyan + " Config (Settings)                       " + NeonLime + "║" + ColorReset)
	fmt.Println(NeonLime + "  ║ " + ColorWhite + "3." + NeonCyan + " Exit (Quit)                                    " + NeonLime + "║" + ColorReset)
	fmt.Println(NeonLime + "  ╚══════════════════════════════════════════════════╝" + ColorReset)
}

func printLiveStats() {
	c := atomic.LoadInt64(&stats.Checked)
	w := atomic.LoadInt64(&stats.Working)
	t := stats.Total
	p := 0.0
	if t > 0 { p = (float64(c) / float64(t)) * 100 }
	
	barLen := 30
	fill := int((float64(c) / float64(t)) * float64(barLen))
	if fill > barLen { fill = barLen }
	bar := strings.Repeat("█", fill) + strings.Repeat("░", barLen-fill)

	fmt.Printf("\r"+NeonPink+" [%s] "+ColorBold+"%.1f%%"+ColorReset+" | "+NeonLime+"Working: %d"+ColorReset+" | "+ColorRed+"Fail: %d"+ColorReset, bar, p, w, c-w)
}

func printSummary(proxies []ProxyResult) {
	fmt.Println("\n\n" + NeonCyan + strings.Repeat("═", 60) + ColorReset)
	fmt.Println(NeonLime + " [✔] PROCESS COMPLETED!" + ColorReset)
	fmt.Printf(" [★] Total Working : %d\n", len(proxies))
	fmt.Printf(" [★] Detailed List : %s (IP | Latency)\n", currentConfig.OutputFile)
	fmt.Printf(" [★] Raw List      : %s (IP:Port Only)\n", currentConfig.RawFile)
	
	if len(proxies) > 0 {
		fmt.Printf(" [★] Fastest Proxy : %s (%s)\n", proxies[0].Address, proxies[0].Latency)
	}
	fmt.Println(NeonCyan + strings.Repeat("═", 60) + ColorReset)
}

func readInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func pause() {
	fmt.Print(DarkGray + "\n[Press Enter to continue]" + ColorReset)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
