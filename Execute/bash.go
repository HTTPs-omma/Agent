package Execute

import "os/exec"

type bash struct {
	isAvailable bool
}

/*
파워셀을 빠르게 사용해보자 : https://stackoverflow.com/questions/65331558/how-to-call-powershell-from-go-faster
*/
func (c *bash) execute(command string) (string, error) {
	// setting
	cmd := exec.Command("sh", "-Command", command)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
