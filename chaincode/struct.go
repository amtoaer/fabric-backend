package main

// Record 电子病历结构体（demo）
type Record struct {
	ObjectType     string
	PatientName    string
	PatientID      string
	DoctorName     string
	DoctorID       string
	Content        string
	ContentEncrypt []byte
	Historys       []HistoryItem
}

// HistoryItem 电子病历历史结构体
type HistoryItem struct {
	TxID    string
	History Record
}
