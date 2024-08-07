package Core

import (
	"errors"
	"github.com/google/uuid"
	"net"

	"agent/Extension"
	"agent/Type"
)

/*
*
Agent 는 한 빌더 패턴 및 추상 팩토리 패턴으로 생성됩니다.
*/
type Agent struct {
	status           string // Const 변수로 다시 변경해야함
	hostIP           string // Const 변수로 다시 변경해야함
	Protocol         string // Const 변수로 다시 변경해야함
	os_name          string
	os_platform      Type.OSPLATFORM
	os_vsersion      string        // ex. 10.4.5
	connect_protocol Type.PROTOCOL // ex. TCP, UDP, HTTP
	UUID             string        // DB 에서 가져와서 할당
}

func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func NewAgent(hostIP string) (*Agent, error) {
	if isValidIP(hostIP) {
		return &Agent{}, errors.New("newAgent 호출시 유효하지 않은 IP 주소가 넘어옴 " + string(hostIP))
		// hostIP 가 string 함수로 문자열 치환이 안되는 경우 새로운 오류가 발생할 수 있음. 추후에는 대비할 것
	}

	// =============
	agent := &Agent{}
	agent.connect_protocol = "TCP"
	// =============

	UUID := uuid.New() // 1. UUID 만들기
	agent.UUID = UUID.String()
	agent.hostIP = hostIP // 동적으로 heap 에 메모리 할당해줌 ( go 언어 특징, 메모리 해제 걱정 필요 없음 )

	sysinfo, err := Extension.NewSysutils()
	if err != nil {
		return &Agent{}, err
	}

	/**
	구조체를 만들어서 매개변수로 만들자
	DSystemInfo
	매개변수로 전달
	*/
	agent.os_vsersion = sysinfo.GetOsVersion()
	agent.os_name = sysinfo.GetOsName()

	return agent, nil
}
