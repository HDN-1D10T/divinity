package main

import "testing"

func TestGetIPsFromCIDR(t *testing.T) {
	ips, err := getIPsFromCIDR("192.0.2.0/30")
	if err != nil {
		t.Fatalf("getIPsFromCIDR returned error: %v", err)
	}
	want := []string{"192.0.2.1", "192.0.2.2"}
	if len(ips) != len(want) {
		t.Fatalf("expected %d IPs, got %d: %v", len(want), len(ips), ips)
	}
	for i := range want {
		if ips[i] != want[i] {
			t.Fatalf("expected IP %d to be %q, got %q", i, want[i], ips[i])
		}
	}
}

func TestGetIPsFromCIDRRejectsZeroRange(t *testing.T) {
	_, err := getIPsFromCIDR("0.0.0.0/0")
	if err == nil {
		t.Fatal("expected /0 CIDR to return an error")
	}
}

func TestIsStdin(t *testing.T) {
	if !isStdin("-") {
		t.Fatal("expected '-' to mean stdin")
	}
	if !isStdin("stdin") || !isStdin("STDIN") {
		t.Fatal("expected stdin sentinel to be case-insensitive")
	}
	if isStdin("a") {
		t.Fatal("did not expect a one-character filename to mean stdin")
	}
}
