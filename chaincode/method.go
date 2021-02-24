package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

//TYPE 电子病历的类型标记
const TYPE = "recordObj"

// 添加电子病历到CouchDB数据库
func putRecord(stub shim.ChaincodeStubInterface, record Record) ([]byte, bool) {
	record.ObjectType = TYPE
	bytes, err := json.Marshal(record)
	if err != nil {
		return nil, false
	}
	// 使用病人ID+医生ID作为key，唯一确定一条数据
	err = stub.PutState(record.PatientID+record.DoctorID, bytes)
	if err != nil {
		return nil, false
	}
	return bytes, true
}

// 通过key得到对应的病历信息
func getRecordInfo(stub shim.ChaincodeStubInterface, key string) (Record, bool) {
	var result Record
	bytes, err := stub.GetState(key)
	if err != nil {
		return result, false
	}
	if bytes == nil {
		return result, false
	}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return result, false
	}
	return result, true
}

// 执行查询字符串并返回结果
func getRecordsByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	resultIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultIterator.Close()
	var buf bytes.Buffer
	buf.WriteString("[")
	flag := false
	for resultIterator.HasNext() {
		resp, err := resultIterator.Next()
		if err != nil {
			return nil, err
		}
		if flag {
			buf.WriteString(",")
		}
		buf.Write(resp.Value)
		flag = true
	}
	buf.WriteString("]")
	return buf.Bytes(), nil
}

// 添加病历
func (c *ChainCode) addRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定参数不符合要求。")
	}
	var record Record
	err := json.Unmarshal([]byte(args[0]), &record)
	if err != nil {
		return shim.Error("反序列化时发生错误。")
	}
	// 查重
	_, exist := getRecordInfo(stub, record.PatientID+record.DoctorID)
	if exist {
		return shim.Error("病历已存在。")
	}
	_, success := putRecord(stub, record)
	if !success {
		return shim.Error("保存信息时发生错误。")
	}
	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("添加病历信息成功。"))
}

// 更新病历
func (c *ChainCode) updateRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定参数不符合要求。")
	}
	var toUpdate Record
	err := json.Unmarshal([]byte(args[0]), &toUpdate)
	if err != nil {
		return shim.Error("反序列化病历信息失败。")
	}
	_, success := getRecordInfo(stub, toUpdate.PatientID+toUpdate.DoctorID)
	if !success {
		return shim.Error("该病历记录不存在，无法更新。")
	}
	_, success = putRecord(stub, toUpdate)
	if !success {
		return shim.Error("保存信息时发生错误。")
	}
	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("信息更新成功。"))
}

// 通过病人ID查询病历
func (c *ChainCode) queryRecordByPatientID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定参数个数不符合要求。")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"%s\",\"PatientID\":\"%s\"}}", TYPE, args[0])
	result, err := getRecordsByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("查询信息时发生错误。")
	}
	if result == nil {
		return shim.Error("没有查询到相关信息。")
	}
	return shim.Success(result)
}

// 通过医生ID查询病历
func (c *ChainCode) queryRecordByDoctorID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定参数个数不符合要求。")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"%s\",\"DoctorID\":\"%s\"}}", TYPE, args[0])
	result, err := getRecordsByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("查询信息时发生错误。")
	}
	if result == nil {
		return shim.Error("没有查询到相关信息。")
	}
	return shim.Success(result)
}

// 通过病人和医生的ID查询病历（包含历史信息）
func (c *ChainCode) queryRecordByKey(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定参数个数不符合要求。")
	}
	bytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("根据key查询信息失败。")
	}
	if bytes == nil {
		return shim.Error("根据key没有查询到相关信息。")
	}
	var result Record
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return shim.Error("反序列化病历信息失败")
	}
	// 获取历史变更
	iter, err := stub.GetHistoryForKey(result.PatientID + result.DoctorID)
	if err != nil {
		return shim.Error("查询历史信息失败。")
	}
	defer iter.Close()
	var historys []HistoryItem
	var historyRecord Record
	for iter.HasNext() {
		historyData, err := iter.Next()
		if err != nil {
			return shim.Error("获取历史变更数据失败。")
		}
		json.Unmarshal(historyData.Value, &historyRecord)
		item := HistoryItem{
			TxID:    historyData.TxId,
			History: historyRecord,
		}
		historys = append(historys, item)
	}
	result.Historys = historys
	bytes, err = json.Marshal(result)
	if err != nil {
		return shim.Error("序列化结果时发生错误。")
	}
	return shim.Success(bytes)
}
