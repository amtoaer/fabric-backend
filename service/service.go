package service

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (i *internal) AddRecord(record Record) (string, error) {
	eventID := "eventAddRecord"
	reg, notifier := registerEvent(i.client, i.chaincodeID, eventID)
	defer i.client.UnregisterChaincodeEvent(reg)
	bytes, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("指定对象序列化失败。")
	}
	req := channel.Request{ChaincodeID: i.chaincodeID, Fcn: "addRecord", Args: [][]byte{bytes, []byte(eventID)}}
	response, err := i.client.Execute(req)
	if err != nil {
		return "", err
	}
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

func (i *internal) UpdateRecord(record Record) (string, error) {
	eventID := "eventUpdateRecord"
	reg, notifier := registerEvent(i.client, i.chaincodeID, eventID)
	defer i.client.UnregisterChaincodeEvent(reg)
	bytes, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("指定对象序列化失败。")
	}
	req := channel.Request{ChaincodeID: i.chaincodeID, Fcn: "updateRecord", Args: [][]byte{bytes, []byte(eventID)}}
	response, err := i.client.Execute(req)
	if err != nil {
		return "", err
	}
	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

func (i *internal) QueryRecordByPatientID(patientID string) (string, error) {
	req := channel.Request{ChaincodeID: i.chaincodeID, Fcn: "queryRecordByPatientID", Args: [][]byte{[]byte(patientID)}}
	resp, err := i.client.Query(req)
	if err != nil {
		return "", err
	}
	return string(resp.Payload), nil
}

func (i *internal) QueryRecordByDoctorID(doctorID string) (string, error) {
	req := channel.Request{ChaincodeID: i.chaincodeID, Fcn: "queryRecordByDoctorID", Args: [][]byte{[]byte(doctorID)}}
	resp, err := i.client.Query(req)
	if err != nil {
		return "", err
	}
	return string(resp.Payload), nil
}

func (i *internal) QueryRecordByKey(patientID, doctorID string) (string, error) {
	req := channel.Request{ChaincodeID: i.chaincodeID, Fcn: "queryRecordByKey", Args: [][]byte{[]byte(patientID + doctorID)}}
	resp, err := i.client.Query(req)
	if err != nil {
		return "", err
	}
	return string(resp.Payload), nil
}
