package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	port      string
	verbosity int8
	api_key   string
	endpoint  string
	version   string
}

//go:embed VERSION
var version string

const max_verbosity = 1

func (c AppConfig) load() (*AppConfig, error) {
	_ = godotenv.Load()

	api_key, ok := os.LookupEnv("OPENROUTER_API_KEY")
	if !ok {
		return nil, errors.New("missing OPENROUTER_API_KEY")
	}

	endpoint, ok := os.LookupEnv("ENDPOINT")
	if endpoint == "" {
		endpoint = "https://openrouter.ai"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "11434"
	}

	verbosity_str := os.Getenv("VERBOSITY")
	verbosity_int, err := strconv.Atoi(verbosity_str)
	// verbosity_int = 0 if err

	if verbosity_str != "" && err != nil {
		fmt.Printf("WARN: invalid verbosity '%s', defaulting to 0.\n", verbosity_str)
	} else if verbosity_int > max_verbosity {
		fmt.Printf("WARN: specified verbosity %d exceeds max verbosity, lowering to %d.\n", verbosity_int, max_verbosity)
		verbosity_int = max_verbosity
	}

	return &AppConfig{
		port:      fmt.Sprintf(":%s", port),
		verbosity: int8(verbosity_int),
		api_key:   api_key,
		version:   version,
		endpoint:  endpoint,
	}, nil
}

func (c AppConfig) Bearer() string {
	return fmt.Sprintf("Bearer %s", c.api_key)
}

func (c AppConfig) UserAgent() string {
	return fmt.Sprintf("tunnellm / %s", c.version)
}

func print_version() {
	fmt.Printf("TunnelLM %s\n", version)
}

func (c AppConfig) PortAvailable() bool {
	ln, err := net.Listen("tcp", c.port)

	if err != nil {
		fmt.Fprintf(os.Stderr, "! Failed to bind %s: %v\n", c.port, err)
		return false
	} else {
		ln.Close()
		return true
	}
}

func main() {
	// if sys.argv[0] == --version
	show_version := flag.Bool("version", false, "print version and exit")

	if *show_version {
		print_version()
		os.Exit(0)
	}

	fmt.Printf("TunnelLM %s\n", version)

	cfg, err := AppConfig{}.load()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	bearer := cfg.Bearer()
	user_agent := cfg.UserAgent()

	target, _ := url.Parse(cfg.endpoint)

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.Host = target.Host

			if !strings.HasPrefix(req.URL.Path, "/api") {
				req.URL.Path = fmt.Sprintf("/api%s", req.URL.Path)
			}

			if cfg.verbosity > 0 {
				fmt.Printf("%s: %s\n", req.Method, req.URL)
			}

			req.Header.Set("Authorization", bearer)
			req.Header.Set("User-Agent", user_agent)
		},

		ModifyResponse: func(resp *http.Response) error {
			if cfg.verbosity > 0 {
				fmt.Printf("%d: %s\n", resp.StatusCode, resp.Request.URL)
			}
			// verbosity 2 would read the body (decompressed) but that broke streaming
			return nil
		},
	}

	http.HandleFunc("/", proxy.ServeHTTP)

	if cfg.PortAvailable() {
		fmt.Printf("→ Listening on localhost%s\n", cfg.port)

		log.Fatal(http.ListenAndServe(cfg.port, nil))
	} else {
		os.Exit(1)
	}

}
