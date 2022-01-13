package log

import(
	"encoding/json"
	"io/ioutil"
	"time"

	"payment/models/Log"
)

func CreateLogHistory(error bool, message, data string) {
	// open file log.json
	fileName := "json/log.json"
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}

	// assign value to struct Log
	object := Log.Log{error, message, data}
	jsonString, _ := json.Marshal(object)

	// assign value to struct LogHistory
	LogData := []Log.LogHistory{}
	json.Unmarshal(file, &LogData)
	newStruct := &Log.LogHistory{
		Id:         len(LogData) + 1,
		Log:        string(jsonString),
		Created_at: time.Now().String(),
		Updated_at: time.Now().String(),
	}
	LogData = append(LogData, *newStruct)

	// Preparing the data to be marshalled and written.
	dataBytes, err := json.MarshalIndent(LogData, "", " ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(fileName, dataBytes, 0644)
	if err != nil {
		return
	}

}