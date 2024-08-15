package Model

import (
	"fmt"
	"testing"
)

func TestApplicationDB_createTable(t1 *testing.T) {
	type fields struct {
		dbName string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{name: "Test case 1"},
	}

	for _, tt := range tests {
		createQeruy()
		fmt.Println(tt.name)
	}
}
