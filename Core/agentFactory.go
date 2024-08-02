package Core

import (
	"agent/Type"
)

/*
	그량 빌더 패턴과 추상 패토리 패턴을 엮어 보고 싶었음.
	이유 없음.

	이거 인터페이스로 기능을 제공할 때 불필요한 내부 함수를 감쳐야함.
*/
type AgentFactory struct {

}



/**

 */
func (af *AgentFactory) getSystemTime() {
 	// os sys

}



func (af *AgentFactory) agentFactory() {
	// 매개 변수로 줘야할 것
	connProto := Type.HTTP


	// =============
	agent := &Agent{}
	agent.connProto = connProto
	agent.os_name = af.getOsName()
	agent.os_vsersion = af.getOsVersion()
	//agent.os_install_date =

}

func (af *AgentFactory) getOsName() Type.OSName {


	return Type.LINUX
}


func (af *AgentFactory) getOsVersion() string {

	return ""
}


func (af *AgentFactory){

}


