package main

import (
	"errors"
	"net"
	"testing"
)

func TestShouldColdFallbackHot(t *testing.T) {
	if shouldColdFallbackHot(nil, `{"ok":true}`) {
		t.Fatal("ok body should not fallback")
	}
	if !shouldColdFallbackHot(errors.New("dial tcp 127.0.0.1:443: connect: connection refused"), "") {
		t.Fatal("connection refused should fallback")
	}
	var ne net.OpError
	_ = ne
	if !shouldColdFallbackHot(errors.New("upstream HTTP 502"), "") {
		t.Fatal("502 should fallback")
	}
}
