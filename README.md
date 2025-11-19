🌐 Proxy Finder

Proxy Finder is a high-performance, multi-threaded command-line tool written in Go (Golang). It scrapes, validates, and benchmarks public proxies from over 30+ distinct sources, supporting HTTP, HTTPS, SOCKS4, and SOCKS5 protocols.

✨ Features

🚀 Ultra Fast: Uses Go routines for concurrent checking (default 500 threads).

🛡️ Multi-Protocol: Supports HTTP, HTTPS, SOCKS4, and SOCKS5.

🌍 Global Sources: Scrapes from 30+ reliable public repositories and APIs (ProxyScrape, Proxifly, TheSpeedX, etc.).

💾 Dual Save Mode: Automatically saves two types of lists:

proxies_detailed.txt: Includes Latency info (e.g., 1.1.1.1:80 | 150ms).

proxies_raw.txt: Clean IP:Port format for use in other tools (e.g., 1.1.1.1:80).

🎲 Random Generator: Experimental feature to generate and test random IP ranges (Brute-force mode).

⚙️ Configurable: Adjust thread count, timeout, and target URL directly from the menu.

📦 Installation & Usage

You need Go (Golang) installed on your machine.

Clone the repository (or download the file):

git clone [https://github.com/yourusername/proxy-finder.git](https://github.com/yourusername/proxy-finder.git)
cd proxy-finder


Run directly:

go run main.go


Build (Optional):
To create a standalone executable (.exe on Windows or binary on Linux/Mac):

go build -o proxy-finder main.go


🖥️ Menu System

When you launch the tool, you will see the following interactive menu:

START SCAN: Begins the scraping and checking process immediately.

CONFIGURATION: Allows you to change settings:

Target URL: The website used to test if the proxy is working (Default: http://www.google.com).

Thread Count: How many proxies to check at the same time (Default: 500).

Timeout: Maximum time to wait for a response (Default: 8s).

Random Generate: Number of random IPs to generate and test (Default: 0 / Disabled).

EXIT: Closes the application.

📂 Output Files

After a scan is completed, two files are generated in the same directory:

File Name

Format

Description

proxies_detailed.txt

`IP:PORT

LATENCY`

proxies_raw.txt

IP:PORT

Good for copy-pasting into other bots/tools.

⚠️ Disclaimer

This tool is for educational purposes and security research only.

The developer is not responsible for any misuse of this tool.

Scanning random IP addresses (Brute-force mode) may be illegal in some jurisdictions. Use with caution and only on networks you own or have permission to test.

The public proxies scraped are not controlled by this tool; do not use them for sensitive transactions (banking, login, etc.).
