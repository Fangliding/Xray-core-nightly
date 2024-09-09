package splithttp_test

import (
	"context"
	"net/http"
	"testing"

	. "github.com/xtls/xray-core/transport/internet/splithttp"
)

type fakeRoundTripper struct{}

func (c *fakeRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return nil, nil
}

func TestConnections(t *testing.T) {
	config := Multiplexing{
		Connections: &RandRangeConfig{From: 4, To: 4},
	}

	mux := NewMuxManager(config, func() http.RoundTripper {
		return &fakeRoundTripper{}
	})

	clients := make(map[interface{}]struct{})
	for i := 0; i < 8; i++ {
		clients[mux.GetClient(context.Background())] = struct{}{}
	}

	if len(clients) != 4 {
		t.Error("did not get 4 distinct clients, got ", len(clients))
	}
}

func TestRequestsPerConnection(t *testing.T) {
	config := Multiplexing{
		RequestsPerConnection: &RandRangeConfig{From: 2, To: 2},
	}

	mux := NewMuxManager(config, func() http.RoundTripper {
		return &fakeRoundTripper{}
	})

	clients := make(map[interface{}]struct{})
	for i := 0; i < 64; i++ {
		clients[mux.GetClient(context.Background())] = struct{}{}
	}

	if len(clients) != 32 {
		t.Error("did not get 32 distinct clients, got ", len(clients))
	}
}

func TestConcurrency(t *testing.T) {
	config := Multiplexing{
		Concurrency: &RandRangeConfig{From: 2, To: 2},
	}

	mux := NewMuxManager(config, func() http.RoundTripper {
		return &fakeRoundTripper{}
	})

	clients := make(map[interface{}]struct{})
	for i := 0; i < 64; i++ {
		client := mux.GetClient(context.Background())
		client.OpenRequests.Add(1)
		clients[client] = struct{}{}
	}

	if len(clients) != 32 {
		t.Error("did not get 32 distinct clients, got ", len(clients))
	}
}

func TestDefault(t *testing.T) {
	config := Multiplexing{}

	mux := NewMuxManager(config, func() http.RoundTripper {
		return &fakeRoundTripper{}
	})

	clients := make(map[interface{}]struct{})
	for i := 0; i < 64; i++ {
		client := mux.GetClient(context.Background())
		client.OpenRequests.Add(1)
		clients[client] = struct{}{}
	}

	if len(clients) != 1 {
		t.Error("did not get 1 distinct clients, got ", len(clients))
	}
}
