package csv

import (
	"bytes"
	"github.com/jfyne/csvd"
)

type ErrorFields struct {
	index     string
	errorType string
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// EligibleAgentsCSVReader read eligibile agents csv bytes data and returns list of eligible agent
func EligibleAgentsCSVReader(data []byte) ([]*byte, error, *byte) {

	reader := csvd.NewReader(bytes.NewBuffer(data))
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true
	_, err := reader.Read()
	if err != nil {
		return nil, err, nil
	}

	//records, err := reader.ReadAll()
	//if err != nil {
	//	return nil, err, nil
	//}

	//totalRecords := len(records)
	//if totalRecords > maxRecords {
	//	return nil, errMaxRecordsExceeded, nil
	//}
	//
	//var eligibleAgents []*model.EligibleAgent
	//var cifIDs []string
	//var errorFieldsArray []ErrorFields
	//
	//for index, record := range records {
	//
	//	cifProspera := record[1]
	//	agentName := record[5]
	//	idCardNumber := record[6]
	//	indexString := strconv.Itoa(index)
	//	data, err := model.NewEligibleAgent(record, index+2)
	//	if err != nil {
	//		errorFieldsArray = append(errorFieldsArray, SetErrorFields(indexString, cifProspera, idCardNumber, agentName, err.Error()))
	//		continue
	//	}
	//	cifID := data.CifProspera + data.IDCardNumber
	//	_, found := Find(cifIDs, cifID)
	//	if !found {
	//		if cifID != "" {
	//			cifIDs = append(cifIDs, cifID)
	//			eligibleAgents = append(eligibleAgents, data)
	//		}
	//	}
	//
	//}
	//
	//errorCsv := NewCSVErrorFields(errorFieldsArray)
	//
	//return eligibleAgents, nil, errorCsv
	return nil, nil, nil
}

func BulkRegisterT24CSVReader(data []byte) (string, error) {

	reader := csvd.NewReader(bytes.NewBuffer(data))
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true
	_, err := reader.Read()
	if err != nil {
		return "", err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	totalRecords := len(records)
	if totalRecords > maxRecords {
		return "", errMaxRecordsExceeded
	}

	var bulkUploadData string

	//for index, record := range records {
	//	data := model.NewBulkUploadRegisterT24Line(record, index+2)
	//	if data != "" {
	//		bulkUploadData = bulkUploadData + data + "\r\n"
	//	}
	//}

	return bulkUploadData, nil
}

func BookingT24CSVReader(data []byte) ([]string, error) {

	reader := csvd.NewReader(bytes.NewBuffer(data))
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true
	_, err := reader.Read()
	if err != nil {
		return nil, err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	totalRecords := len(records)
	if totalRecords > maxRecords {
		return nil, errMaxRecordsExceeded
	}

	var t24Booking []string

	//for index, record := range records {
	//	data := model.NewBookingT24Line(record, index+2)
	//	if data != "" {
	//		t24Booking = append(t24Booking, data)
	//	}
	//
	//}

	return t24Booking, nil
}

func PRKCSVReader(data []byte) ([]string, error) {

	reader := csvd.NewReader(bytes.NewBuffer(data))
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true
	_, err := reader.Read()
	if err != nil {
		return nil, err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	totalRecords := len(records)
	if totalRecords > maxRecords {
		return nil, errMaxRecordsExceeded
	}

	var PRK []string

	//for index, record := range records {
	//	data := model.NewPRKLine(record, index+2)
	//	if data != "" {
	//		PRK = append(PRK, data)
	//	}
	//
	//}

	return PRK, nil
}

func PaylaterCSVReader(data []byte) ([]string, error) {

	reader := csvd.NewReader(bytes.NewBuffer(data))
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true
	_, err := reader.Read()
	if err != nil {
		return nil, err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	totalRecords := len(records)
	if totalRecords > maxRecords {
		return nil, errMaxRecordsExceeded
	}

	var Paylater []string

	//for index, record := range records {
	//	data := model.NewPaylaterLine(record, index+2)
	//	if data != "" {
	//		Paylater = append(Paylater, data)
	//	}
	//
	//}

	return Paylater, nil
}

func UploadProsperaCSVReader(data []byte) ([]*byte, error, *byte) {

	reader := csvd.NewReader(bytes.NewBuffer(data))
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true
	_, err := reader.Read()
	if err != nil {
		return nil, err, nil
	}

	//records, err := reader.ReadAll()
	//if err != nil {
	//	return nil, err, nil
	//}

	//totalRecords := len(records)
	//if totalRecords > maxRecords {
	//	return nil, errMaxRecordsExceeded, nil
	//}
	//
	//var prosperaData []*model.BulkUploadRegisterT24
	//var cifIDs []string
	//var errorFieldsArray []ErrorFields
	//
	//for index, record := range records {
	//
	//	cifProspera := record[1]
	//	agentName := record[5]
	//	idCardNumber := record[6]
	//	indexString := strconv.Itoa(index)
	//	data, err := model.NewProsperaData(record, index+2)
	//
	//	if err != nil {
	//		errorFieldsArray = append(errorFieldsArray, SetErrorFields(indexString, cifProspera, idCardNumber, agentName, err.Error()))
	//		continue
	//	}
	//	cifID := data.AltCustID
	//	_, found := Find(cifIDs, cifID)
	//
	//	if !found {
	//		if cifID != "" {
	//			cifIDs = append(cifIDs, cifID)
	//			prosperaData = append(prosperaData, data)
	//		}
	//	}
	//
	//}
	//
	//errorCsv := NewCSVErrorFields(errorFieldsArray)
	//
	//return prosperaData, nil, errorCsv

	return nil, nil, nil
}

func NewCSVErrorFields(data []ErrorFields) *byte {
	//var result []*pb.Field
	//for _, error := range data {
	//	result = append(result, &pb.Field{
	//		Name:  "Index/CIF Prospera/ID Card Number/Agent Name/Error Type",
	//		Value: error.index + "/" + error.cifProspera + "/" + error.idCardNumber + "/" + error.agentName + "/" + error.errorType,
	//	})
	//}
	//errorCount := len(data)
	//return &pb.ErrorCSV{
	//	Line:   int32(errorCount),
	//	Fields: result,
	//}
	return nil
}
