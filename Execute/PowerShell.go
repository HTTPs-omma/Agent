package Execute

import (
	"fmt"
	"os/exec"
)

type PowerShell struct {
	Name        string
	IsAvailable bool
}

func NewPowerShell() *PowerShell {
	return &PowerShell{
		Name:        "PowerShell",
		IsAvailable: true,
	}
}

func (p *PowerShell) Execute(command string) (string, error) {
	cmd := exec.Command("powershell", "-Command", command)

	// 표준 출력 및 표준 에러를 함께 캡처
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %s, error: %w", string(output), err)
	}

	return string(output), nil
}

///*
//파워셀을 빠르게 사용해보자 : https://stackoverflow.com/questions/65331558/how-to-call-powershell-from-go-faster
//*/
//func (p *PowerShell) execute_osExec(commands []string) ([]string, error) {
//	cmd := exec.Command("powershell", "-nologo", "-noprofile")
//	stdin, err := cmd.StdinPipe()
//	if err != nil {
//		log.Fatal(err)
//		p.isAvailable = false
//		return nil, err
//	}
//	p.isAvailable = true
//
//	go func() {
//		defer stdin.Close()
//		for _, command := range commands {
//			fmt.Fprintf(stdin, "%s\n", command)
//		}
//	}()
//
//	// 이렇게 안하면 속도가 너무 느림
//	out, err := cmd.CombinedOutput() //비동기적으로 동작한 함수가 실행된 코루틴의 작업이 완료될 때 까지 대기한다.
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Println(string(out))
//
//	return nil, nil
//}
