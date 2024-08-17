package Execute

import "testing"

func TestPowerShell_execute(t *testing.T) {
	type fields struct {
		isAvailable bool
	}
	type args struct {
		commands []string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "Powershell Execute Test",
			fields: fields{isAvailable: true},
			args:   args{commands: []string{"dir", "cd ../", "dir"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PowerShell{
				isAvailable: tt.fields.isAvailable,
			}

			p.execute(tt.args.commands)
		})
	}
}
