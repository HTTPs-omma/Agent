package Core

import (
	"gopkg.in/yaml.v3"
)

type InstructionData struct {
	ID               string `yaml:"id"`
	MITREID          string `yaml:"MITRE_ID"`
	Description      string `yaml:"Description"`
	Escalation       bool   `yaml:"Escalation"`
	Tool             string `yaml:"tool"`
	RequisiteCommand string `yaml:"requisite_command"`
	Command          string `yaml:"command"`
	Cleanup          string `yaml:"cleanup"`
}

// ToBytes는 InstructionData 구조체를 YAML 바이트 슬라이스로 변환하는 함수입니다.
func (cd *InstructionData) ToBytes() ([]byte, error) {
	// YAML로 직렬화
	data, err := yaml.Marshal(cd)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (cd *InstructionData) GetInstData(bdata []byte) (*InstructionData, error) {
	// 새로운 InstructionData 인스턴스를 생성
	var instData InstructionData
	//fmt.Println("bytes : ", bdata)

	// YAML 데이터를 InstructionData 구조체로 역직렬화
	err := yaml.Unmarshal(bdata, &instData)
	if err != nil {
		return nil, err
	}

	return &instData, nil
}
