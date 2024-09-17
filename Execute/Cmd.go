package Execute

import (
	"fmt"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os/exec"
	"strings"
)

type Cmd struct {
	isAvailable bool
}

/*
파워셀을 빠르게 사용해보자 : https://stackoverflow.com/questions/65331558/how-to-call-powershell-from-go-faster
*/
func (c *Cmd) Execute(command string) (string, error) {
	// setting
	fmt.Println("cmd : " + command)
	cmd := exec.Command("cmd", "/C", command)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	decodedOutput, err := decodeCP949(output)
	if err != nil {
		return "", err
	}

	return string(decodedOutput), nil
}
func decodeCP949(input []byte) (string, error) {
	reader := transform.NewReader(strings.NewReader(string(input)), korean.EUCKR.NewDecoder())
	decoded, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
