package Core

import (
	"testing"
)

func TestAgentFactory_getSystemTime(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			af := &AgentFactory{}
			af.getSystemTime()
		})
	}
}
