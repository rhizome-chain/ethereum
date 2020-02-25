package log

type EthLogJobInfo struct {
	SourceJobID string `json:"source"`
	DataType    string `json:"dataType"`
	LogType     string `json:"logType"`
}
