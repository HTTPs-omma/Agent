package Model

import (
	"fmt"
	"os"
	"testing"
)

func Test_db(t *testing.T) {
	os.Chdir("../")

	fmt.Printf("gid : %d \n", os.Getgid())

	tests := []struct {
		name string
	}{
		// To DO
		{name: "Test case 1"},
		//fmt.Println("db open in test_db"),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("db open in test_db")
		})
	}
}
