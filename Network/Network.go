package Network

import (
	"agent/Model"
	"fmt"
	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
)

func SendPacket(hs HSProtocol.HS) (*HSProtocol.HS, error) {
	agdb, err := Model.NewAgentStatusDB()
	if err != nil {
		return nil, err
	}

	statusRecords, err := agdb.SelectAllRecords()
	record := statusRecords[0]

	switch record.Protocol {
	case HSProtocol.TCP:
		return SendPacketByTcp(hs)
	case HSProtocol.HTTP:
		return sendPacketByHttp(hs)
	}

	return nil, fmt.Errorf("protocol not support")
}
