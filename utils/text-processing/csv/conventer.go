package csv

import (
	"bytes"
	"fmt"
)

func RemoveDecimal(val string) string {
	if val == "" {
		return "0"
	}
	if val == "0" {
		return "0"
	}
	return val[:len(val)-2]
}

func StrVal(val interface{}) string {
	switch val.(type) {
	case float64:
		return fmt.Sprintf("%.0f", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func BoolVal(val string) bool {
	if val == "1" {
		return true
	}
	return false
}

func TransactionDataCSVConverter(data []*byte, glAccount *byte) (string, []*byte, error) {

	var transactionData string
	//var line int = 1
	//for idx, record := range data {
	//	trxData, err := NewTransactionData(record, glAccount.Code)
	//	if err != nil {
	//		return "", data, err
	//	}
	//	if trxData != "" {
	//		transactionData = transactionData + trxData + "\r\n"
	//		data[idx].LineNumber = strconv.Itoa(line)
	//		line++
	//	}
	//}
	transactionData = transactionData + "999"

	return transactionData, data, nil
}

func NewTransactionData(data *byte, code string) (string, error) {
	var b bytes.Buffer

	//nominal := fmt.Sprintf("%.0f", data.AmtInstallment)
	//remarks := fmt.Sprintf("%s %s #%d", autodebitPMD, data.FinancingCode, data.InstallmentNum)

	//b.WriteString(SVAToGL + ",")
	//b.WriteString(SVA + ",")
	//b.WriteString(data.SavingWOW_IB + ",")
	//b.WriteString(GL + ",")
	//b.WriteString(code + ",")
	//b.WriteString(nominal + ",")
	//b.WriteString(remarks)

	return b.String(), nil
}

func TransactionDataRobotikCSVConverter(data []*byte, glAccount *byte) (string, []*byte, error) {

	var transactionData string
	//var line int = 1
	//for idx, record := range data {
	//	trxData, err := NewTransactionDataRobotik(record, glAccount.Code)
	//	if err != nil {
	//		return "", data, err
	//	}
	//	if trxData != "" {
	//		transactionData = transactionData + trxData + "\r\n"
	//		data[idx].LineNumber = strconv.Itoa(line)
	//		line++
	//	}
	//}

	return transactionData, data, nil
}

func NewTransactionDataRobotik(data *byte, code string) (string, error) {
	var b bytes.Buffer

	//nominal := fmt.Sprintf("%.0f", data.AmtInstallment)
	//remarks := fmt.Sprintf("%s %s #%d", autodebitPMD, data.FinancingCode, data.InstallmentNum)
	//
	//b.WriteString(SVAToGLRobotik + ",")
	//b.WriteString(SVA + ",")
	//b.WriteString(data.SavingWOW_IB + ",")
	//b.WriteString(GL + ",")
	//b.WriteString(code + ",")
	//b.WriteString(nominal + ",")
	//b.WriteString(remarks)

	return b.String(), nil
}

func getDebitCreditCode(code string) (string, string) {
	//if code == GLToGL {
	//	return GL, GL
	//}
	//
	//if code == SVAToGL {
	//	return SVA, GL
	//}
	//
	//if code == GLToSVA {
	//	return GL, SVA
	//}

	return "", ""
}

func ErrorDataCSVConverter(data []string) (string, error) {

	var errorData string
	for _, record := range data {
		if record != "" {
			errorData = errorData + record + "\r\n"
		}
	}

	return errorData, nil
}
