package Execute

import (
	"bytes"
	"os/exec"
)

type Cmd struct {
	cmd *exec.Cmd
}

func NewCmd() *Cmd {
	return &Cmd{
		cmd: exec.Command("cmd", "/C"),
	}
}

func (c *Cmd) Execute(command string) (string, error) {
	// 프로세스의 stdin, stdout, stderr에 파이프를 연결
	stdin, err := c.cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	var outBuf, errBuf bytes.Buffer
	c.cmd.Stdout = &outBuf
	c.cmd.Stderr = &errBuf

	// 프로세스 시작
	err = c.cmd.Start()
	if err != nil {
		return "", err
	}

	// 명령어 입력을 파이프로 보내기
	_, err = stdin.Write([]byte(command + "\n"))
	if err != nil {
		return "", err
	}
	stdin.Close() // 입력 스트림 종료

	// 프로세스 실행 완료 대기
	err = c.cmd.Wait()
	if err != nil {
		return "", err
	}

	// 명령 실행 후 결과 반환
	if errBuf.Len() > 0 {
		return errBuf.String(), nil
	}
	return outBuf.String(), nil
}
