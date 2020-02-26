package main

import (
	"github.com/piquette/finance-go"
	"testing"
)

func TestGetBTCUSDCNY(t *testing.T) {
	finance.NewBackends(newHttpClient("socks5://127.0.0.1:1080"))
	t.Log(GetBTCUSDCNY())
}
