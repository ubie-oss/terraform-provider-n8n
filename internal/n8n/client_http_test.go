package n8n

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestClientRespectsMaxConcurrent(t *testing.T) {
	var peak atomic.Int32
	var active atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cur := active.Add(1)
		for {
			old := peak.Load()
			if cur <= old {
				break
			}
			if peak.CompareAndSwap(old, cur) {
				break
			}
		}
		time.Sleep(40 * time.Millisecond)
		active.Add(-1)
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	client, err := New(srv.URL, "secret", &Options{
		MaxConcurrent: 3,
		RPS:           500,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	const workers = 25
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, client.Endpoint, nil)
			if err != nil {
				t.Errorf("NewRequest: %v", err)
				return
			}
			resp, err := client.HTTP.Do(req)
			if err != nil {
				t.Errorf("Do: %v", err)
				return
			}
			_ = resp.Body.Close()
		}()
	}
	wg.Wait()

	if got := peak.Load(); got > 3 {
		t.Fatalf("expected peak concurrent <= 3, got %d", got)
	}
}

func TestClientRateLimitsSerialRequests(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	client, err := New(srv.URL, "secret", &Options{
		MaxConcurrent: 10,
		RPS:           2,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	start := time.Now()
	for i := 0; i < 6; i++ {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, client.Endpoint, nil)
		if err != nil {
			t.Fatalf("NewRequest: %v", err)
		}
		resp, err := client.HTTP.Do(req)
		if err != nil {
			t.Fatalf("Do: %v", err)
		}
		_ = resp.Body.Close()
	}
	elapsed := time.Since(start)
	if elapsed < 1200*time.Millisecond {
		t.Fatalf("expected serial requests to take at least ~1.2s with RPS=2, got %v", elapsed)
	}
}

func TestClientRefusesCrossHostRedirect(t *testing.T) {
	var targetHit atomic.Bool
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetHit.Store(true)
		if r.Header.Get("X-N8N-API-KEY") != "" {
			t.Error("API key forwarded to cross-host redirect target")
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(target.Close)

	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL+"/landed", http.StatusFound)
	}))
	t.Cleanup(origin.Close)

	client, err := New(origin.URL, "leak-me-key", &Options{RPS: 100})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, client.Endpoint, nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	resp, err := client.HTTP.Do(req)
	if err == nil {
		_ = resp.Body.Close()
		t.Fatal("expected error refusing cross-host redirect")
	}
	if targetHit.Load() {
		t.Fatal("cross-host redirect target was contacted")
	}
}
