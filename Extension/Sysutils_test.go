package Extension

import (
	"fmt"
	"testing"
)

func TestSysutils_getOsName(t1 *testing.T) {
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
		sys, err := NewSysutils()
		if err != nil {
			t1.Fatal(err)
		}

		t1.Run(tt.name, func(t1 *testing.T) {
			fmt.Println("getOsName : ", sys.GetOsName())
			fmt.Println("getHostName : ", sys.GetHostName())
			fmt.Println("getArchitecture : ", sys.GetArchitecture())
			fmt.Println("getContainerized : ", sys.GetContainerized())
			fmt.Println("getBootTime : ", sys.GetBootTime())
			fmt.Println("getIPs : ", sys.GetIPs())
			fmt.Println("getMACs : ", sys.GetMACs())
			fmt.Println("getFamily : ", sys.GetFamily())
			fmt.Println("getOsVersion : ", sys.GetOsVersion())
		})
	}
}
