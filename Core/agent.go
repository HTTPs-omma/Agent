package Core

import (
	"time"

	"agent/Type"
)

/*
*
Agent 는 한 빌더 패턴 및 추상 팩토리 패턴으로 생성됩니다.
*/
type Agent struct {
	status          string      // Const 변수로 다시 변경행야함
	hostIP          string      // Const 변수로 다시 변경행야함
	Protocol        string      // Const 변수로 다시 변경행야함
	os_name         Type.OSName // ex. window
	os_vsersion     string      // ex. 10.4.5
	os_install_date time.Time   //
	os_SystemTime   time.Time   //
	connProto       Type.Protocol
	UUID            string //
}

func (agent *Agent) init() {

}

/*
*
refernce : https://pkg.go.dev/time
What the fuck code!
os_SystemTime 의 존재 의미가 있을까요? 일단 만들어 놓고 삭제를 논의 해봅시다.
*/
func (agent *Agent) getSystemTime() time.Time {
	agent.os_SystemTime = time.Now()
	return agent.os_SystemTime
}

func (agent *Agent) getOSName() string {
	return agent.os_name
}

func (agent *Agent) getOs_install_date() string {
	return agent.getOs_install_date()
}
