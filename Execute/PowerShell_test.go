package Execute

//
//import (
//	"fmt"
//	"testing"
//)
//
//func TestPowerShell_Osexecute(t *testing.T) {
//	type fields struct {
//		isAvailable bool
//	}
//	type args struct {
//		commands []string
//	}
//
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		{name: "Powershell Execute Test",
//			fields: fields{isAvailable: true},
//			args:   args{commands: []string{"dir", "cd ../", "dir"}},
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			p := &PowerShell{
//				IsAvailable: tt.fields.isAvailable,
//			}
//
//			p.execute_osExec(tt.args.commands)
//		})
//	}
//}
//
//func TestPowerShell_CreatePrc(t *testing.T) {
//	type fields struct {
//		isAvailable bool
//	}
//	type args struct {
//		commands []string
//	}
//
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		{name: "Powershell Execute Test",
//			fields: fields{isAvailable: true},
//			args:   args{commands: []string{"dir", "cd ../", "dir"}},
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			p := &PowerShell{
//				IsAvailable: tt.fields.isAvailable,
//			}
//
//			output, err := p.execute_createPrc(tt.args.commands[0])
//			if err != nil {
//				t.Fatal(err)
//			}
//			fmt.Print(output)
//		})
//	}
//}
