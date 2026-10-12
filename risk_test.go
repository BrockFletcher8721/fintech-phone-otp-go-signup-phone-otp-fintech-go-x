package main

import "testing"

func TestAllowPayment(t *testing.T) {
	tests := []struct {
		name   string
		amount float64
		risk   string
		want   bool
	}{
		{"normal transfer", 250, "low", true},
		{"high risk", 250, "high", false},
		{"limit", 1000, "low", true},
		{"over limit", 1000.01, "low", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := allowPayment(tt.amount, tt.risk); got != tt.want {
				t.Fatalf("allowPayment() = %v, want %v", got, tt.want)
			}
		})
	}
}
