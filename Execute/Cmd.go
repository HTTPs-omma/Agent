package Execute

import (
	"bytes"
	"fmt"
	"os/exec"
)

type Cmd struct {
	cmd *exec.Cmd
}

func NewCmd() *Cmd {
	return &Cmd{}
}

// Execute: 명령어를 실행하고 결과를 반환하는 함수
func (c *Cmd) Execute(command string) (string, error) {
	// 명령어 실행
	fmt.Println("run Cmd Code : " + command)

	// cmd /C 와 함께 명령어 실행
	// 리디렉션 사용을 위해 명령어를 따옴표로 묶어 실행합니다
	cmd := exec.Command("cmd", "/C", command)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	// 명령어 실행
	err := cmd.Run()
	if err != nil {
		return errBuf.String(), fmt.Errorf("failed to execute command: %w", err)
	}

	// 명령어 출력 결과 반환
	if errBuf.Len() > 0 {
		return errBuf.String(), nil
	}
	return outBuf.String(), nil
}
