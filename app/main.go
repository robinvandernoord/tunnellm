package main

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/andybalholm/brotli"
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

const max_verbosity = 2

func (c AppConfig) load() AppConfig {
	_ = godotenv.Load()

	api_key, ok := os.LookupEnv("OPENROUTER_API_KEY")
	if !ok {
		panic("missing OPENROUTER_API_KEY")
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

	return AppConfig{
		port:      fmt.Sprintf(":%s", port),
		verbosity: int8(verbosity_int),
		api_key:   api_key,
		version:   version,
		endpoint:  endpoint,
	}
}

func (c AppConfig) Bearer() string {
	return fmt.Sprintf("Bearer %s", c.api_key)
}

func (c AppConfig) UserAgent() string {
	return fmt.Sprintf("tunnellm / %s", c.version)
}

func decompress_body(resp *http.Response, body_bytes []byte) ([]byte, error) {
	encoding := resp.Header.Get("Content-Encoding")
	if encoding == "" {
		return body_bytes, nil
	}

	reader := bytes.NewReader(body_bytes)

	switch encoding {
	case "gzip":
		gz_reader, err := gzip.NewReader(reader)
		if err != nil {
			return nil, err
		}
		defer gz_reader.Close()
		return io.ReadAll(gz_reader)

	case "deflate":
		deflate_reader := flate.NewReader(reader)
		defer deflate_reader.Close()
		return io.ReadAll(deflate_reader)

	case "br":
		br_reader := brotli.NewReader(reader)
		return io.ReadAll(br_reader)

	default:
		return nil, fmt.Errorf("unknown content encoding: %s", encoding)
	}
}

func main() {
	cfg := AppConfig{}.load()
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
			if cfg.verbosity > 1 {
				body_bytes, err := io.ReadAll(resp.Body)
				if err != nil {
					return err
				}

				// fmt.Printf("%s", string(body_bytes))
				decompressed, err := decompress_body(resp, body_bytes)
				if err == nil {
					// todo: this breaks streaming, we should be using a tee‑reader
					fmt.Printf("%s\n", decompressed)
				}

				// put body back so the proxy can forward it
				resp.Body = io.NopCloser(bytes.NewBuffer(body_bytes))

			}
			return nil
		},
	}

	http.HandleFunc("/", proxy.ServeHTTP)

	fmt.Printf("TunneLLM running on localhost%s\n", cfg.port)

	log.Fatal(http.ListenAndServe(cfg.port, nil))
}
