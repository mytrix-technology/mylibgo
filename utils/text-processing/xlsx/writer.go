package xlsx

import (
	"fmt"
	"github.com/xuri/excelize/v2"
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

func EglsReportWriter(totalAmount int64, datas []*byte) (*excelize.File, error) {

	//currentDate := time.Now().In(model.Location).Format(model.ShortDateFormat)
	//desc := fmt.Sprintf("Manual Adjustment EGLS 15421 trx tgl %s (Autodebet PMD)", currentDate)
	idr := "IDR"
	debit := "D"
	credit := "C"
	fixedBranchID := "ID0019002"
	accountID := "215421000000"
	profitCostCenterVal := "1"
	dimension02Val := "9002"
	dimension04Val := "1"
	xlsx := excelize.NewFile()

	bold, err := xlsx.NewStyle(`{
    	"font": {
        	"bold": true
    	}
	}`)
	if err != nil {
		return nil, err
	}

	sheet1Name := "Jurnal Trx Wow"
	xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)

	//
	xlsx.SetCellValue(sheet1Name, "B2", "Batch Name")
	xlsx.SetCellStyle(sheet1Name, "B2", "B2", bold)
	//xlsx.SetCellValue(sheet1Name, "C2", desc)

	//
	xlsx.SetCellValue(sheet1Name, "B4", "Transaction")
	xlsx.MergeCell(sheet1Name, "B4", "I4")
	xlsx.SetCellValue(sheet1Name, "J4", "Dimensions")
	xlsx.MergeCell(sheet1Name, "J4", "AC4")

	//
	xlsx.SetCellValue(sheet1Name, "B5", "Type")
	xlsx.MergeCell(sheet1Name, "B5", "B6")

	xlsx.SetCellValue(sheet1Name, "C5", "Description")
	xlsx.MergeCell(sheet1Name, "C5", "C6")

	xlsx.SetCellValue(sheet1Name, "D5", "Account")
	xlsx.MergeCell(sheet1Name, "D5", "F5")
	xlsx.SetCellValue(sheet1Name, "D6", "CCY")
	xlsx.SetCellValue(sheet1Name, "E6", "Branch ID")
	xlsx.SetCellValue(sheet1Name, "F6", "Account ID")

	xlsx.SetCellValue(sheet1Name, "G5", "Transaction")
	xlsx.MergeCell(sheet1Name, "G5", "H5")
	xlsx.SetCellValue(sheet1Name, "G6", "CCY")
	xlsx.SetCellValue(sheet1Name, "H6", "Amount")

	xlsx.SetCellValue(sheet1Name, "I5", "Profit Cost Center")
	xlsx.MergeCell(sheet1Name, "I5", "I6")

	xlsx.SetCellValue(sheet1Name, "J5", "Dimension01")
	xlsx.MergeCell(sheet1Name, "J5", "J6")
	xlsx.SetCellValue(sheet1Name, "K5", "Dimension02")
	xlsx.MergeCell(sheet1Name, "K5", "K6")
	xlsx.SetCellValue(sheet1Name, "L5", "Dimension03")
	xlsx.MergeCell(sheet1Name, "L5", "L6")
	xlsx.SetCellValue(sheet1Name, "M5", "Dimension04")
	xlsx.MergeCell(sheet1Name, "M5", "M6")
	xlsx.SetCellValue(sheet1Name, "N5", "Dimension05")
	xlsx.MergeCell(sheet1Name, "N5", "N6")
	xlsx.SetCellValue(sheet1Name, "O5", "Dimension06")
	xlsx.MergeCell(sheet1Name, "O5", "O6")
	xlsx.SetCellValue(sheet1Name, "P5", "Dimension07")
	xlsx.MergeCell(sheet1Name, "P5", "P6")
	xlsx.SetCellValue(sheet1Name, "Q5", "Dimension08")
	xlsx.MergeCell(sheet1Name, "Q5", "Q6")
	xlsx.SetCellValue(sheet1Name, "R5", "Dimension09")
	xlsx.MergeCell(sheet1Name, "R5", "R6")
	xlsx.SetCellValue(sheet1Name, "S5", "Dimension10")
	xlsx.MergeCell(sheet1Name, "S5", "S6")
	xlsx.SetCellValue(sheet1Name, "T5", "Dimension11")
	xlsx.MergeCell(sheet1Name, "T5", "T6")
	xlsx.SetCellValue(sheet1Name, "U5", "Dimension12")
	xlsx.MergeCell(sheet1Name, "U5", "U6")
	xlsx.SetCellValue(sheet1Name, "V5", "Dimension13")
	xlsx.MergeCell(sheet1Name, "V5", "V6")
	xlsx.SetCellValue(sheet1Name, "W5", "Dimension14")
	xlsx.MergeCell(sheet1Name, "W5", "W6")
	xlsx.SetCellValue(sheet1Name, "X5", "Dimension15")
	xlsx.MergeCell(sheet1Name, "X5", "X6")
	xlsx.SetCellValue(sheet1Name, "Y5", "Dimension16")
	xlsx.MergeCell(sheet1Name, "Y5", "Y6")
	xlsx.SetCellValue(sheet1Name, "Z5", "Dimension17")
	xlsx.MergeCell(sheet1Name, "Z5", "Z6")
	xlsx.SetCellValue(sheet1Name, "AA5", "Dimension18")
	xlsx.MergeCell(sheet1Name, "AA5", "AA6")
	xlsx.SetCellValue(sheet1Name, "AB5", "Dimension19")
	xlsx.MergeCell(sheet1Name, "AB5", "AB6")
	xlsx.SetCellValue(sheet1Name, "AC5", "Dimension20")
	xlsx.MergeCell(sheet1Name, "AC5", "AC6")

	//
	xlsx.SetCellValue(sheet1Name, "B7", debit)
	//xlsx.SetCellValue(sheet1Name, "C7", desc)
	xlsx.SetCellValue(sheet1Name, "D7", idr)
	xlsx.SetCellValue(sheet1Name, "E7", fixedBranchID)
	xlsx.SetCellValue(sheet1Name, "F7", accountID)
	xlsx.SetCellValue(sheet1Name, "G7", idr)
	xlsx.SetCellValue(sheet1Name, "H7", RemoveDecimal(StrVal(totalAmount)))
	xlsx.SetCellValue(sheet1Name, "I7", profitCostCenterVal)
	xlsx.SetCellValue(sheet1Name, "K7", dimension02Val)
	xlsx.SetCellValue(sheet1Name, "M7", dimension04Val)

	rowNumber := 8
	//for _, data := range datas {
	//	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("B%d", rowNumber), credit)
	//	//xlsx.SetCellValue(sheet1Name, fmt.Sprintf("C%d", rowNumber), desc)
	//	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("D%d", rowNumber), idr)
	//	//xlsx.SetCellValue(sheet1Name, fmt.Sprintf("E%d", rowNumber), getBranchID(data.BranchID))
	//	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("F%d", rowNumber), accountID)
	//	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("G%d", rowNumber), idr)
	//	//xlsx.SetCellValue(sheet1Name, fmt.Sprintf("H%d", rowNumber), RemoveDecimal(StrVal(data.Amount)))
	//	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("I%d", rowNumber), profitCostCenterVal)
	//	//xlsx.SetCellValue(sheet1Name, fmt.Sprintf("K%d", rowNumber), data.BranchID)
	//	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("M%d", rowNumber), dimension04Val)
	//	rowNumber++
	//}

	//
	rowNumber++
	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("E%d", rowNumber), credit)
	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("F%d", rowNumber), RemoveDecimal(StrVal(totalAmount)))

	rowNumber++
	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("E%d", rowNumber), debit)
	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("F%d", rowNumber), RemoveDecimal(StrVal(totalAmount)))

	//
	rowNumber += 2
	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("E%d", rowNumber), "Diajukan Oleh")
	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("F%d", rowNumber), "Mengetahui")

	rowNumber += 2
	xlsx.MergeCell(sheet1Name, fmt.Sprintf("E%d", rowNumber), fmt.Sprintf("E%d", rowNumber+2))
	xlsx.MergeCell(sheet1Name, fmt.Sprintf("F%d", rowNumber), fmt.Sprintf("F%d", rowNumber+2))

	rowNumber += 3
	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("E%d", rowNumber), "Nama :")
	xlsx.SetCellValue(sheet1Name, fmt.Sprintf("F%d", rowNumber), "Nama :")

	rowNumber++
	//xlsx.SetCellValue(sheet1Name, fmt.Sprintf("E%d", rowNumber), currentDate)
	//xlsx.SetCellValue(sheet1Name, fmt.Sprintf("F%d", rowNumber), currentDate)

	return xlsx, nil
}

func getBranchID(branchID string) string {
	if branchID == "" {
		return ""
	}
	return fmt.Sprintf("ID001%s", branchID)
}
