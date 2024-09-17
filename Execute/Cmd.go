package Execute

import "os/exec"

type Cmd struct {
	isAvailable bool
}

/*
파워셀을 빠르게 사용해보자 : https://stackoverflow.com/questions/65331558/how-to-call-powershell-from-go-faster
*/
func (c *Cmd) Execute(command string) (string, error) {
	// setting
	cmd := exec.Command("cmd", "-Command", command)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
