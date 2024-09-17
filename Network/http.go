package Network

import (
	"bytes"
	"fmt"
	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
	"io"
	"net/http"
	"os"
)

//https://agent

func SendHTTPRequest(hs HSProtocol.HS) (*HSProtocol.HS, error) {

	// HS 객체를 직렬화 (예: ToBytes 함수 사용)
	HSMgr := HSProtocol.NewHSProtocolManager()
	data, err := HSMgr.ToBytes(&hs)
	if err != nil {
		fmt.Println("Error serializing HS  object:", err)
		return nil, err
	}

	// HTTP POST 요청 생성
	url := "http://" + os.Getenv("SERVER_IP") + "/api/checkInstReq" // 실제 서버 URL로 변경
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error : ", err)
	}

	hsmgr := HSProtocol.NewHSProtocolManager()
	ack, err := hsmgr.Parsing(body)
	if err != nil {
		fmt.Println("Error : ", err)
	}

	return ack, nil

	//instD, err := inst.GetInstData(ack.Data)
	//if err != nil {
	//	fmt.Println("Error : ", err)
	//}
	//
	//shell := Execute.PowerShell{}
	//cmdLog, err := shell.Execute(instD.Command)
	//if err != nil {
	//
	//}
	//shell.Execute(instD.Cleanup)
}
