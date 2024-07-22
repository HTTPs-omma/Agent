package Core

/*
*
Agent 는 한 빌더 패턴 및 추상 팩토리 패턴으로 생성됩니다.
*/
type Agent struct {
	status          string // Const 변수로 다시 변경행야함
	hostIP          string // Const 변수로 다시 변경행야함
	Protocol        string // Const 변수로 다시 변경행야함
	os_name         string // ex. window
	os_vsersion     string // ex. 10.4.5
	os_install_date string //
	os_tiemstamp    string //
	UUID            string //
	timstamp        string
}

func (anget *Agent) init() {

}
