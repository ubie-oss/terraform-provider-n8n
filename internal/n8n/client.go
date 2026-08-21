package n8n

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/sync/semaphore"
	"golang.org/x/time/rate"
)

// Defaults for rate and concurrency when the provider omits optional attributes.
const (
	DefaultMaxConcurrent = int64(10)
	DefaultRPS           = 10.0
)

// Options configures HTTP rate limiting and concurrency. Zero values mean "use defaults".
type Options struct {
	MaxConcurrent int64   // max in-flight HTTP requests (default DefaultMaxConcurrent)
	RPS           float64 // sustained requests per second for the token bucket (default DefaultRPS)
}

// Client is the n8n Public API client. Use HTTP for REST calls; transport applies auth, rate, and concurrency limits.
// Endpoint is the API base including /api/v1 (for example https://n8n.example.com/api/v1).
type Client struct {
	Endpoint string
	HTTP     *http.Client
}

// New validates configuration and returns a client for use as provider ResourceData/DataSourceData.
// endpoint is the n8n instance URL (with or without /api/v1). opts may be nil; zero fields select defaults.
func New(endpoint, apiKey string, opts *Options) (*Client, error) {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, fmt.Errorf("api_key must not be empty")
	}

	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint must not be empty")
	}

	normalized, err := NormalizeEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	endpointURL, err := url.Parse(normalized)
	if err != nil {
		return nil, fmt.Errorf("invalid normalized endpoint: %w", err)
	}

	maxC, rps := DefaultMaxConcurrent, DefaultRPS
	if opts != nil {
		if opts.MaxConcurrent < 0 {
			return nil, fmt.Errorf("max concurrent requests cannot be negative")
		}
		if opts.RPS < 0 {
			return nil, fmt.Errorf("requests per second cannot be negative")
		}
		if opts.MaxConcurrent > 0 {
			maxC = opts.MaxConcurrent
		}
		if opts.RPS > 0 {
			rps = opts.RPS
		}
	}

	burst := int(math.Ceil(rps))
	if burst > 100 {
		burst = 100
	}

	lim := rate.NewLimiter(rate.Limit(rps), burst)
	sem := semaphore.NewWeighted(maxC)

	rt := &roundTripper{
		apiKey:  apiKey,
		apiHost: endpointURL.Host,
		lim:     lim,
		sem:     sem,
		base:    http.DefaultTransport,
	}

	return &Client{
		Endpoint: normalized,
		HTTP: &http.Client{
			Timeout:   60 * time.Second,
			Transport: rt,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return fmt.Errorf("stopped after 10 redirects")
				}
				if len(via) > 0 && !sameHTTPHost(req.URL, via[0].URL) {
					return fmt.Errorf("refusing cross-host redirect from %q to %q", via[0].URL.Host, req.URL.Host)
				}
				return nil
			},
		},
	}, nil
}

// NormalizeEndpoint trims trailing slashes and ensures the path ends with /api/v1.
func NormalizeEndpoint(endpoint string) (string, error) {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return "", fmt.Errorf("endpoint must not be empty")
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("invalid endpoint URL: %w", err)
	}
	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("endpoint must include scheme and host")
	}

	path := strings.TrimSuffix(u.Path, "/")
	if !strings.HasSuffix(path, "/api/v1") {
		if path == "" {
			path = "/api/v1"
		} else {
			path = path + "/api/v1"
		}
	}
	u.Path = path
	u.RawQuery = ""
	u.Fragment = ""

	return strings.TrimSuffix(u.String(), "/"), nil
}

func sameHTTPHost(a, b *url.URL) bool {
	return strings.EqualFold(a.Host, b.Host)
}

// roundTripper applies concurrency limit, then rate limit, then API key auth, then delegates.
type roundTripper struct {
	apiKey  string
	apiHost string
	lim     *rate.Limiter
	sem     *semaphore.Weighted
	base    http.RoundTripper
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	// Acquire concurrency first so rate tokens are not spent while queued.
	if err := rt.sem.Acquire(ctx, 1); err != nil {
		return nil, err
	}
	defer rt.sem.Release(1)
	if err := rt.lim.Wait(ctx); err != nil {
		return nil, err
	}

	r2 := req.Clone(ctx)
	// Only attach the API key for the configured API host (defense in depth vs redirects).
	if strings.EqualFold(r2.URL.Host, rt.apiHost) {
		r2.Header.Set("X-N8N-API-KEY", rt.apiKey)
	}
	return rt.base.RoundTrip(r2)
}
