package Execute

import (
	"fmt"
	"log"
	"os/exec"
)

type PowerShell struct {
	isAvailable bool
}

/*
파워셀을 빠르게 사용해보자 : https://stackoverflow.com/questions/65331558/how-to-call-powershell-from-go-faster
*/
func (p *PowerShell) execute(commands []string) {
	cmd := exec.Command("powershell", "-nologo", "-noprofile")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
		p.isAvailable = false
		return
	}
	p.isAvailable = true

	go func() {
		defer stdin.Close()
		for _, command := range commands {
			fmt.Fprintf(stdin, "%s\n", command)
		}
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", out)
}
