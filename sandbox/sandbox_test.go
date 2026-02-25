package sandbox_test

import (
	"testing"

	"github.com/kosuke9809/yadokari/sandbox"
)

func TestRiskString(t *testing.T) {
	tests := []struct {
		risk sandbox.Risk
		want string
	}{
		{sandbox.RiskStrict, "🔒"},
		{sandbox.RiskDev, "🟡"},
		{sandbox.RiskHigh, "🔴"},
	}
	for _, tt := range tests {
		if got := tt.risk.String(); got != tt.want {
			t.Errorf("Risk(%d).String() = %q, want %q", tt.risk, got, tt.want)
		}
	}
}
