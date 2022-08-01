package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

func RegisterT24BulkWrite(w io.Writer, data []string) error {
	var msgs []string
	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write(data)

	if len(msgs) > 0 {
		return fmt.Errorf("%s. %w", strings.Join(msgs, ","), errMaxRecordsExceeded)
	}

	return nil
}

func ReconciliationReportWriter(w io.Writer, datas []*byte) {
	//writer := csv.NewWriter(w)
	//
	//defer writer.Flush()
	//
	//titleRow := []string{
	//	"No",
	//	"Transaction Date",
	//	"T24 Date",
	//	"Transaction Type",
	//	"T24 Transaction Amount",
	//	"Wow Syariah Transaction Amount",
	//	"FT T24 Number",
	//	"Loan Qardh Number",
	//	"Cif T24 Number",
	//	"Wow Syariah Agent Number",
	//}
	//writer.Write(titleRow)
	//
	//for index, data := range datas {
	//	wowAmount := "-"
	//	t24Amount := "-"
	//	transactionAmount := ""
	//	wowSyariahTransactionAmount := ""
	//	if data.T24TransactionAmount != nil {
	//		transactionAmount = strconv.Itoa(int(*data.T24TransactionAmount))
	//	}
	//	if data.WowSyariahTransactionAmount != nil {
	//		wowSyariahTransactionAmount = strconv.Itoa(int(*data.WowSyariahTransactionAmount))
	//	}
	//	switch data.TransactionType {
	//	case "Posting Credit", "Reversal", "Reversal - Murabahah", "Reversal - Qardh", "Reversal - Ujrah", "Reversal - Payoff", "Reversal - Unknown source", "Posting Credit - Murabahah", "Posting Credit - Qardh", "Posting Credit - Ujrah", "Posting Credit - Payoff", "Posting Credit - Unknown source":
	//		wowAmount = "(" + wowSyariahTransactionAmount + ")"
	//	case "Payment", "Payment - Murabahah", "Payment - Qardh", "Payment - Ujrah", "Payment - Payoff", "Payment - Unknown source":
	//		t24Amount = "(" + transactionAmount + ")"
	//	case "Posting Debit", "Posting Debit - Murabahah", "Posting Debit - Qardh", "Posting Debit - Ujrah", "Posting Debit - Payoff", "Posting Debit - Unknown source":
	//		wowAmount = wowSyariahTransactionAmount
	//	case "Disburse", "Disburse - Murabahah", "Disburse - Qardh", "Disburse - Unknown type":
	//		t24Amount = transactionAmount
	//	default:
	//	}
	//	record := []string{
	//		strVal(index + 1),
	//		data.TransactionDate.Format(model.ShortDateFormat),
	//		data.T24_date.Format(model.ShortDateFormat),
	//		data.TransactionType,
	//		t24Amount,
	//		wowAmount,
	//		data.FtT24Number,
	//		data.LoanQardhNumber,
	//		data.CifT24Number,
	//		data.WowSyariahAgentNumber,
	//	}
	//	writer.Write(record)
	//}
}

func AmountRangeWriter(w io.Writer, datas []*byte) {
	//writer := csv.NewWriter(w)
	//
	//defer writer.Flush()
	//
	//titleRow := []string{
	//	"No",
	//	"Type",
	//	"Min Value",
	//	"Max Value",
	//	"Created At",
	//	"Updated At",
	//}
	//writer.Write(titleRow)
	//
	//for index, data := range datas {
	//	minValue := strconv.Itoa(int(data.MinValue))
	//	maxValue := strconv.Itoa(int(data.MaxValue))
	//	createdAt := data.CreatedAt.Format(model.ShortDateFormat)
	//	UpdatedAt := data.UpdatedAt.Format(model.ShortDateFormat)
	//	record := []string{
	//		strVal(index + 1),
	//		data.Type,
	//		minValue,
	//		maxValue,
	//		createdAt,
	//		UpdatedAt,
	//	}
	//	writer.Write(record)
	//}
}

func SummaryBillsWriter(w io.Writer, datas []*byte) {
	//writer := csv.NewWriter(w)
	//
	//defer writer.Flush()
	//
	//titleRow := []string{
	//	"No",
	//	"Cif Prospera",
	//	"Cif T24",
	//	"Total Tunggakan Paylater",
	//	"Total Tunggakan Angsuran Paylater",
	//	"Total Tunggakan Biaya Paylater",
	//	"Total Tunggakan Qardh Revolving",
	//	"Total Tunggakan Angsuran Qardh Revolving",
	//	"Total Tunggakan Biaya Qardh Revolving",
	//	"Created At",
	//	"Created By",
	//	"Updated At",
	//	"Updated By",
	//	"Ujroh Payment",
	//}
	//writer.Write(titleRow)
	//
	//for index, data := range datas {
	//	ujrohPayment := strconv.Itoa(int(*data.UjrohPayment))
	//	totalTunggakanPaylater := strconv.Itoa(int(*data.TotalTunggakanPaylater))
	//	totalTunggakanAngsuranPaylater := strconv.Itoa(int(*data.TotalTunggakanAngsuranPaylater))
	//	totalTunggakanBiayaPaylater := strconv.Itoa(int(*data.TotalTunggakanBiayaPaylater))
	//	totalTunggakanQardhRevolving := strconv.Itoa(int(*data.TotalTunggakanQardhRevolving))
	//	totalTunggakanAngsuranQardhRevolving := strconv.Itoa(int(*data.TotalTunggakanAngsuranQardhRevolving))
	//	totalTunggakanBiayaQardhRevolving := strconv.Itoa(int(*data.TotalTunggakanBiayaQardhRevolving))
	//	createdAt := data.CreatedAt.Format(model.ShortDateFormat)
	//	createdBy := strconv.Itoa(int(*&data.CreatedBy))
	//	updatedBy := strconv.Itoa(int(*&data.UpdatedBy))
	//	UpdatedAt := data.UpdatedAt.Format(model.ShortDateFormat)
	//	record := []string{
	//		strVal(index + 1),
	//		data.CifProspera,
	//		data.CifT24,
	//		totalTunggakanPaylater,
	//		totalTunggakanAngsuranPaylater,
	//		totalTunggakanBiayaPaylater,
	//		totalTunggakanQardhRevolving,
	//		totalTunggakanAngsuranQardhRevolving,
	//		totalTunggakanBiayaQardhRevolving,
	//		createdAt,
	//		createdBy,
	//		UpdatedAt,
	//		updatedBy,
	//		ujrohPayment,
	//	}
	//	writer.Write(record)
	//}
}

func AgentReportWriter(w io.Writer, datas []*byte) {
	//writer := csv.NewWriter(w)
	//
	//defer writer.Flush()
	//
	//if len(datas) < 1 {
	//	writer.Write([]string{
	//		"Data Tidak Tersedia",
	//	})
	//} else {
	//	titleRow := []string{
	//		"No",
	//		"Data Date",
	//		"Cif Prospera",
	//		"Cif T24",
	//		"Account Number Wow Syariah",
	//		"Account Number T24",
	//		"Agent Name",
	//		"ID Card Number",
	//		"Total Limit",
	//		"Qardh Limit",
	//		"Murabahah Limit",
	//		"Expired Offer Date",
	//		"Tenor Facility",
	//		"Notes",
	//		"Approval Limit Number",
	//		"Agent Status",
	//		"Approval Date",
	//		"Expired Limit Date",
	//		"Upload Date",
	//		"Created By",
	//		"Created At",
	//		"Updated By",
	//		"Updated At",
	//		"Previous Status",
	//		"Daily Counter",
	//		"Booking ID",
	//		"Is Paid Fully",
	//		"Is Paylater Payment Checked",
	//	}
	//	writer.Write(titleRow)
	//
	//	for index, data := range datas {
	//		totalLimit := strconv.Itoa(int(data.TotalLimit))
	//		qardhLimit := strconv.Itoa(int(data.QardhLimit))
	//		murabahahLimit := strconv.Itoa(int(data.MurabahahLimit))
	//		tenorFacility := strconv.Itoa(int(data.TenorFacility))
	//		createdBy := strconv.Itoa(int(data.CreatedBy))
	//		updatedBy := strconv.Itoa(int(data.UpdatedBy))
	//		data.DataDate = strings.Replace(data.DataDate, "T00:00:00Z", "", -1)
	//		data.ExpiredOfferDate = strings.Replace(data.ExpiredOfferDate, "T00:00:00Z", "", -1)
	//		dailyCounter := strconv.Itoa(int(data.DailyCounter))
	//		approvalDate := data.ApprovalDate.Time.Format(model.ShortDateFormat)
	//		if data.ApprovalDate.Time.IsZero() {
	//			approvalDate = ""
	//		}
	//		expiredLimitDate := data.ExpiredLimitDate.Time.Format(model.ShortDateFormat)
	//		if data.ExpiredLimitDate.Time.IsZero() {
	//			expiredLimitDate = ""
	//		}
	//		uploadDate := data.UploadDate.Format(model.ShortDateFormat)
	//		if data.UploadDate.IsZero() {
	//			uploadDate = ""
	//		}
	//		createdAt := data.CreatedAt.Format(model.ShortDateFormat)
	//		if data.CreatedAt.IsZero() {
	//			createdAt = ""
	//		}
	//		updatedAt := data.UpdatedAt.Format(model.ShortDateFormat)
	//		if data.UpdatedAt.IsZero() {
	//			updatedAt = ""
	//		}
	//		record := []string{
	//			strVal(index + 1),
	//			data.DataDate,
	//			data.CifProspera,
	//			data.CifT24,
	//			data.AccountNumberWowSyariah,
	//			data.AccountNumberT24,
	//			data.AgentName,
	//			data.IDCardNumber,
	//			totalLimit,
	//			qardhLimit,
	//			murabahahLimit,
	//			data.ExpiredOfferDate,
	//			tenorFacility,
	//			data.Notes,
	//			data.ApprovalLimitNumber,
	//			data.AgentStatus.Enum().String(),
	//			approvalDate,
	//			expiredLimitDate,
	//			uploadDate,
	//			createdBy,
	//			createdAt,
	//			updatedBy,
	//			updatedAt,
	//			data.PreviousStatus.Enum().String(),
	//			dailyCounter,
	//			data.BookingID,
	//			strconv.FormatBool(data.IsPaidFully),
	//			strconv.FormatBool(data.IsPaylaterPaymentChecked),
	//		}
	//		writer.Write(record)
	//	}
	//}
}

func CifDataWriter(w io.Writer, datas []*byte) {
	//writer := csv.NewWriter(w)
	//
	//defer writer.Flush()
	//
	//if len(datas) < 1 {
	//	writer.Write([]string{
	//		"Data Tidak Tersedia",
	//	})
	//} else {
	//	titleRow := []string{
	//		"No",
	//		"Header ID",
	//		"Cif T24",
	//		"Alt Cust ID",
	//		"Alt Ac ID",
	//		"Mnemonic",
	//		"Cust Title 1",
	//		"Short Name",
	//		"Cust Title 2",
	//		"Name 1",
	//		"Name 2",
	//		"Given Names",
	//		"Gender",
	//		"Place Birth",
	//		"Moth Maiden",
	//		"Legal Type",
	//		"Legal ID No",
	//		"Expiry Date ID",
	//		"Reside Y N",
	//		"Nationality",
	//		"Taxable",
	//		"Tax ID",
	//		"Religion",
	//		"Marital Status",
	//		"Education",
	//		"Education Other",
	//		"Sector",
	//		"Industry",
	//		"Target",
	//		"Account Officer",
	//		"Cust Type",
	//		"Language",
	//		"Street",
	//		"Address",
	//		"RT_RW",
	//		"Province",
	//		"Country",
	//		"Town_Country",
	//		"Residence",
	//		"District Code",
	//		"Post Code",
	//		"Residence Status",
	//		"Other Res Status",
	//		"Addr Phone No",
	//		"Off Phone",
	//		"Fax 1",
	//		"Sms 1",
	//		"Email 1",
	//		"Add Type",
	//		"Add Street",
	//		"Add RT RW",
	//		"Add Province",
	//		"Add Suburb Twin",
	//		"Add Municipal",
	//		"Add Country",
	//		"Add District",
	//		"Add Post Code",
	//		"Job Title",
	//		"Employment Status",
	//		"Occupation",
	//		"Employer Buss",
	//		"Employer Name",
	//		"Employers Add",
	//		"Employment Start",
	//		"Fund Prov Name",
	//		"Fund Prov Job",
	//		"Fund Prov Addr",
	//		"Fund Prov Phone",
	//		"Fund Source",
	//		"Oth Fund Source",
	//		"Fund Source Amt",
	//		"Oth Acct Type",
	//		"Oth Acct No",
	//		"Oth Ac Branch",
	//		"Oth Ac Bank Name",
	//		"Oth Ac Opened",
	//		"Oth Remarks",
	//		"Contact Name",
	//		"Contact Street",
	//		"Contact Homtel",
	//		"Contat Rel Cus",
	//		"No Debit Trans",
	//		"Value Dr Trans",
	//		"No Credit Trans",
	//		"Value Cr Trans",
	//		"No Credit Trans",
	//		"Value Dr Trans",
	//		"No Credit Trans",
	//		"Value Cr Trans",
	//		"High Risk",
	//		"Guarantor Code",
	//		"Sid Relati Bank",
	//		"Din Number",
	//		"BMPK Violation",
	//		"BMPK Exceeding",
	//		"LBU Cust Type",
	//		"Customer Rating",
	//		"Cu Rate Date",
	//		"LBBU Cust Type",
	//		"Customer Since",
	//		"Upload Company",
	//		"Res Status",
	//		"Res Year",
	//		"Res Month",
	//		"Total Employee",
	//		"Pur Relati Bank",
	//		"Total Liability",
	//		"Att Status",
	//		"Relation Code",
	//		"Rel Customer",
	//		"Portofolio Categ",
	//		"Addr Phone Area",
	//		"Phone AO",
	//		"LLC Y N",
	//		"Ati Llc Bname",
	//		"Ati Llc Acc",
	//		"Ati Llc Oacc",
	//		"Fatca Y N",
	//		"Project Locate",
	//		"Marketing Code",
	//		"No Aplikasi",
	//		"Tgl Pencairan",
	//		"Status Upload",
	//		"Nama File",
	//		"Created By",
	//		"Created At",
	//		"Updated By",
	//		"Updated At",
	//		"Prospera ID",
	//		"Error Log",
	//		"Booking Log",
	//		"Marital Status Description",
	//		"Gender Description",
	//		"Prospera Religion Description",
	//		"Prospera Religion Code",
	//		"Prospera Education Description",
	//		"Prospera Education Code",
	//		"Prospera Gender Description",
	//		"Prospera Gender Code",
	//		"Prospera Fund Source Description",
	//		"Prospera Fund Source Code",
	//		"Prospera Job Title Description",
	//		"Prospera Job Title Code",
	//		"Prospera Trx In Year Description",
	//		"Prospera Trx In Year Code",
	//		"Prospera Residence Description",
	//		"Prospera Residence Code",
	//		"Prospera Province Description",
	//		"Prospera Province Code",
	//		"Prospera Kelurahan Description",
	//		"Prospera District Description",
	//		"Prospera District Code",
	//		"CDC Time",
	//		"Prospera Marital Status Description",
	//		"Prospera Marital Status Code",
	//	}
	//	writer.Write(titleRow)
	//
	//	for index, data := range datas {
	//
	//		record := []string{
	//			strVal(index + 1),
	//			data.HeaderID,
	//			data.Cif_T24,
	//			data.AltCustID,
	//			data.AltAcID,
	//			data.Mnemonic,
	//			data.Cust_Title_1,
	//			data.ShortName,
	//			data.Cust_Title_2,
	//			data.Name_1,
	//			data.Name_2,
	//			data.GivenNames,
	//			data.Gender,
	//			data.PlaceBirth,
	//			data.DateOfBirth.Local().Format(model.ShortDateFormat2),
	//			data.MothMaiden,
	//			data.LegalType,
	//			data.LegalIDNo,
	//			data.ExpiryDateID.Local().Format(model.ShortDateFormat2),
	//			data.Reside_y_n,
	//			data.Nationality,
	//			data.Taxable,
	//			strVal(data.TaxID),
	//			strVal(data.Religion),
	//			data.MaritalStatus,
	//			strVal(data.Education),
	//			strVal(data.EducationOther),
	//			strVal(data.Sector),
	//			strVal(data.Industry),
	//			strVal(data.Target),
	//			strVal(data.AcccountOfficer),
	//			strVal(data.CustType),
	//			strVal(data.Language),
	//			data.Street,
	//			data.Address,
	//			data.Rt_Rw,
	//			strVal(data.Province),
	//			strVal(data.Country),
	//			strVal(data.TownCountry),
	//			strVal(data.Residence),
	//			strVal(data.DistrictCode),
	//			strVal(data.PostCode),
	//			strVal(data.ResidenceStatus),
	//			strVal(data.OthResStatus),
	//			strVal(data.AddrPhoneNo),
	//			strVal(data.OffPhone),
	//			strVal(data.Fax_1),
	//			strVal(data.Sms_1),
	//			strVal(data.Email_1),
	//			strVal(data.AddType),
	//			strVal(data.AddStreet),
	//			strVal(data.AddRtRw),
	//			strVal(data.AddProvince),
	//			strVal(data.AddSuburbTwin),
	//			strVal(data.AddMunicipal),
	//			strVal(data.AddCountry),
	//			strVal(data.AddDistrict),
	//			strVal(data.AddPostCode),
	//			strVal(data.JobTitle),
	//			strVal(data.EmploymentStatus),
	//			strVal(data.Occupation),
	//			strVal(data.EmployersBuss),
	//			strVal(data.EmployersName),
	//			strVal(data.EmployersAdd),
	//			strVal(data.EmploymentStart),
	//			strVal(data.FundProvName),
	//			strVal(data.FundProvJob),
	//			strVal(data.FundProvAddr),
	//			strVal(data.FundSource),
	//			strVal(data.OthFundSource),
	//			strVal(data.FundSourceAmt),
	//			strVal(data.OthAcctType),
	//			strVal(data.OthAcctNo),
	//			strVal(data.OthAcBranch),
	//			strVal(data.OthAcBankName),
	//			strVal(data.OthAcOpened),
	//			strVal(data.OthRemarks),
	//			strVal(data.ContactName),
	//			strVal(data.ContactStreet),
	//			strVal(data.ContactHomtel),
	//			strVal(data.Contact_Rel_Cus),
	//			strVal(data.NoDebitTrans),
	//			strVal(data.ValueDrTrans),
	//			strVal(data.NoCreditTrans),
	//			strVal(data.ValueCrTrans),
	//			strVal(data.HighRisk),
	//			strVal(data.GuarantorCode),
	//			strVal(data.SidRelatiBank),
	//			strVal(data.DinNumber),
	//			strVal(data.BMPKViolation),
	//			strVal(data.BMPKExceeding),
	//			strVal(data.LBUCustType),
	//			strVal(data.CustomerRating),
	//			strVal(data.CURateDate),
	//			strVal(data.LBBUCustType),
	//			strVal(data.CustomerSince),
	//			strVal(data.ResStatus),
	//			strVal(data.ResYear),
	//			strVal(data.ResMonth),
	//			strVal(data.TotalEmployee),
	//			strVal(data.PurRelatiBank),
	//			strVal(data.TotalLiability),
	//			strVal(data.AttStatus),
	//			strVal(data.RelationCode),
	//			strVal(data.RelCustomer),
	//			strVal(data.PortofolioCateg),
	//			strVal(data.AddrPhoneArea),
	//			strVal(data.PhoneAO),
	//			strVal(data.Llc_y_n),
	//			strVal(data.AtiLLCBname),
	//			strVal(data.AtiLLCAcc),
	//			strVal(data.AtiLLCOacc),
	//			strVal(data.Fatca_y_n),
	//			strVal(data.ProjectLocate),
	//			strVal(data.MarketingCode),
	//			strVal(data.NoAplikasi),
	//			strVal(data.TglPencairan),
	//			"",
	//			"",
	//			strVal(data.MaritalStatusDescription),
	//			strVal(data.GenderDescription),
	//			strVal(data.ProsperaReligionDescription),
	//			strVal(data.ProsperaReligionCode),
	//			strVal(data.ProsperaEducationDescription),
	//			strVal(data.ProsperaEducationCode),
	//			strVal(data.ProsperaGenderDescription),
	//			strVal(data.ProsperaGenderCode),
	//			strVal(data.ProsperaFundSourceDescription),
	//			strVal(data.ProsperaFundSourceCode),
	//			strVal(data.ProsperaJobTitleDescription),
	//			strVal(data.ProsperaJobTitleCode),
	//			strVal(data.ProsperaTrxInYearDescription),
	//			strVal(data.ProsperaTrxInYearCode),
	//			strVal(data.ProsperaResidenceStatusDescription),
	//			strVal(data.ProsperaResidenceStatusCode),
	//			strVal(data.ProsperaProvinceName),
	//			strVal(data.ProsperaProvinceCode),
	//			strVal(data.ProsperaKelurahanDescription),
	//			strVal(data.ProsperaDistrictDescription),
	//			strVal(data.ProsperaDistrictCode),
	//			strVal(data.CdcTime),
	//		}
	//		writer.Write(record)
	//	}
	//}
}

func LimitReportWriter(w io.Writer, datas []*byte) {
	//writer := csv.NewWriter(w)
	//
	//defer writer.Flush()
	//
	//if len(datas) < 1 {
	//	writer.Write([]string{
	//		"Data Tidak Tersedia",
	//	})
	//} else {
	//	titleRow := []string{
	//		"No",
	//		"Cif Prospera",
	//		"Cif T24",
	//		"Loan Number Qardh",
	//		"ID Total Limit",
	//		"Total Limit",
	//		"Outstanding Total Limit",
	//		"ID Qardh Limit",
	//		"Qardh Limit",
	//		"Outstanding Qardh Limit",
	//		"Id Murabahah Limit",
	//		"Murabahah Limit",
	//		"Outstanding Murabahah Limit",
	//		"Limit Status",
	//		"Expired Limit Date",
	//	}
	//	writer.Write(titleRow)
	//
	//	for index, data := range datas {
	//		totalLimit := strconv.Itoa(int(data.TotalLimit))
	//		outstandingTotalLimit := strconv.Itoa(int(data.OutstandingTotalLimit))
	//		qardhLimit := strconv.Itoa(int(data.QardhLimit))
	//		outstandingQardhLimit := strconv.Itoa(int(data.OutstandingQardhLimit))
	//		murabahahLimit := strconv.Itoa(int(data.MurabahahLimit))
	//		outstandingMurabahahLimit := strconv.Itoa(int(data.OutstandingMurabahahLimit))
	//		dataExpiredLimit := data.ExpiredLimitDate.Time.Format(model.ShortDateFormat)
	//		if data.ExpiredLimitDate.Time.IsZero() {
	//			dataExpiredLimit = ""
	//		}
	//		record := []string{
	//			strVal(index + 1),
	//			data.CifProspera,
	//			data.CifT24,
	//			data.LoanNumberQardh,
	//			data.IDTotalLimit,
	//			totalLimit,
	//			outstandingTotalLimit,
	//			data.IDQardhLimit,
	//			qardhLimit,
	//			outstandingQardhLimit,
	//			data.IDMurabahahLimit,
	//			murabahahLimit,
	//			outstandingMurabahahLimit,
	//			data.LimitStatus.Enum().String(),
	//			dataExpiredLimit,
	//		}
	//		writer.Write(record)
	//	}
	//}
}

func EmailWriter(w io.Writer, datas []*byte) {
	//writer := csv.NewWriter(w)
	//
	//defer writer.Flush()
	//
	//if len(datas) < 1 {
	//	writer.Write([]string{
	//		"Data Tidak Tersedia",
	//	})
	//} else {
	//	titleRow := []string{
	//		"No",
	//		"Type",
	//		"Addre",
	//		"Created At",
	//		"Updated At",
	//		"To CC",
	//	}
	//	writer.Write(titleRow)
	//
	//	for index, data := range datas {
	//
	//		record := []string{
	//			strVal(index + 1),
	//			strVal(data.Type),
	//			strVal(data.Address),
	//			strVal(data.CreatedAt),
	//			strVal(data.UpdatedAt),
	//			strVal(data.ToCc),
	//		}
	//		writer.Write(record)
	//	}
	//}
}

func CIFT24ErrorReportWriter(w io.Writer, datas []*byte) {
	//writer := csv.NewWriter(w)
	//
	//defer writer.Flush()
	//
	//titleRow := []string{
	//	"header_id",
	//	"cif_t24",
	//	"alt_cust_id",
	//	"alt_ac_id",
	//	"mnemonic",
	//	"cust_title_1",
	//	"short_name",
	//	"cust_title_2",
	//	"name_1",
	//	"name_2",
	//	"given_names",
	//	"gender",
	//	"place_birth",
	//	"date_of_birth",
	//	"moth_maiden",
	//	"legal_type",
	//	"legal_id_no",
	//	"expiry_date_id",
	//	"reside_y_n",
	//	"nationality",
	//	"taxable",
	//	"tax_id",
	//	"religion",
	//	"marital_status",
	//	"education",
	//	"education_other",
	//	"sector",
	//	"industry",
	//	"target",
	//	"acccount_officer",
	//	"cust_type",
	//	"language",
	//	"street",
	//	"address",
	//	"rt_rw",
	//	"province",
	//	"country",
	//	"town_country",
	//	"residence",
	//	"district_code",
	//	"post_code",
	//	"residence_status",
	//	"oth_res_status",
	//	"addr_phone_no",
	//	"off_phone",
	//	"fax_1",
	//	"sms_1",
	//	"email_1",
	//	"add_type",
	//	"add_street",
	//	"add_rt_rw",
	//	"add_province",
	//	"add_suburb_twin",
	//	"add_municipal",
	//	"add_country",
	//	"add_district",
	//	"add_post_code",
	//	"job_title",
	//	"employment_status",
	//	"occupation",
	//	"employers_buss",
	//	"employers_name",
	//	"employers_add",
	//	"employment_start",
	//	"fund_prov_name",
	//	"fund_prov_job",
	//	"fund_prov_addr",
	//	"fund_prov_phone",
	//	"fund_source",
	//	"oth_fund_source",
	//	"fund_source_amt",
	//	"oth_acct_type",
	//	"oth_acct_no",
	//	"oth_ac_branch",
	//	"oth_ac_bank_name",
	//	"oth_ac_opened",
	//	"oth_remarks",
	//	"contact_name",
	//	"contact_street",
	//	"contact_homtel",
	//	"contact_rel_cus",
	//	"no_debit_trans",
	//	"value_dr_trans",
	//	"no_credit_trans",
	//	"value_cr_trans",
	//	"high_risk",
	//	"guarantor_code",
	//	"sid_relati_bank",
	//	"din_number",
	//	"bmpk_violation",
	//	"bmpk_exceeding",
	//	"lbu_cust_type",
	//	"customer_rating",
	//	"cu_rate_date",
	//	"lbbu_cust_type",
	//	"customer_since",
	//	"upload_company",
	//	"res_status",
	//	"res_year",
	//	"res_month",
	//	"total_employee",
	//	"pur_relati_bank",
	//	"total_liability",
	//	"att_status",
	//	"relation_code",
	//	"rel_customer",
	//	"portofolio_categ",
	//	"addr_phone_area",
	//	"phone_ao",
	//	"llc_y_n",
	//	"ati_llc_bname",
	//	"ati_llc_acc",
	//	"ati_llc_oacc",
	//	"fatca_y_n",
	//	"project_locate",
	//	"marketing_code",
	//	"no_aplikasi",
	//	"tgl_pencairan",
	//	"prospera_id",
	//	"Keterangan error CIF",
	//	"Keterangan error Booking",
	//}
	//writer.Write(titleRow)
	//
	//for _, data := range datas {
	//	record := []string{
	//		data.HeaderID,
	//		data.Cif_T24,
	//		data.AltCustID,
	//		data.AltAcID,
	//		data.Mnemonic,
	//		data.Cust_Title_1,
	//		data.ShortName,
	//		data.Cust_Title_2,
	//		data.Name_1,
	//		data.Name_2,
	//		data.GivenNames,
	//		data.Gender,
	//		data.PlaceBirth,
	//		strVal(data.DateOfBirth),
	//		data.MothMaiden,
	//		data.LegalType,
	//		data.LegalIDNo,
	//		strVal(data.ExpiryDateID),
	//		data.Reside_y_n,
	//		data.Nationality,
	//		data.Taxable,
	//		strVal(data.TaxID),
	//		strVal(data.Religion),
	//		data.MaritalStatus,
	//		strVal(data.Education),
	//		data.EducationOther,
	//		strVal(data.Sector),
	//		strVal(data.Industry),
	//		strVal(data.Target),
	//		strVal(data.AcccountOfficer),
	//		data.CustType,
	//		strVal(data.Language),
	//		data.Street,
	//		data.Address,
	//		data.Rt_Rw,
	//		strVal(data.Province),
	//		strVal(data.Country),
	//		strVal(data.TownCountry),
	//		data.Residence,
	//		strVal(data.DistrictCode),
	//		strVal(data.PostCode),
	//		data.ResidenceStatus,
	//		data.OthResStatus,
	//		strVal(data.AddrPhoneNo),
	//		strVal(data.OffPhone),
	//		strVal(data.Fax_1),
	//		strVal(data.Sms_1),
	//		data.Email_1,
	//		data.AddType,
	//		data.AddStreet,
	//		data.AddRtRw,
	//		strVal(data.AddProvince),
	//		strVal(data.AddSuburbTwin),
	//		strVal(data.AddMunicipal),
	//		data.AddCountry,
	//		strVal(data.AddDistrict),
	//		strVal(data.AddPostCode),
	//		strVal(data.JobTitle),
	//		data.EmploymentStatus,
	//		data.Occupation,
	//		data.EmployersBuss,
	//		data.EmployersName,
	//		data.EmployersAdd,
	//		strVal(data.EmploymentStart),
	//		data.FundProvName,
	//		strVal(data.FundProvJob),
	//		data.FundProvAddr,
	//		data.FundProvPhone,
	//		strVal(data.FundSource),
	//		data.OthFundSource,
	//		strVal(data.FundSourceAmt),
	//		data.OthAcctType,
	//		data.OthAcctNo,
	//		data.OthAcBranch,
	//		strVal(data.OthAcBankName),
	//		strVal(data.OthAcOpened),
	//		data.OthRemarks,
	//		data.ContactName,
	//		data.ContactStreet,
	//		strVal(data.ContactHomtel),
	//		strVal(data.Contact_Rel_Cus),
	//		strVal(data.NoDebitTrans),
	//		strVal(data.ValueDrTrans),
	//		strVal(data.NoCreditTrans),
	//		strVal(data.ValueCrTrans),
	//		strVal(data.HighRisk),
	//		strVal(data.GuarantorCode),
	//		strVal(data.SidRelatiBank),
	//		strVal(data.DinNumber),
	//		data.BMPKViolation,
	//		data.BMPKExceeding,
	//		strVal(data.LBUCustType),
	//		strVal(data.CustomerRating),
	//		strVal(data.CURateDate),
	//		strVal(data.LBBUCustType),
	//		strVal(data.CustomerSince),
	//		data.UploadCompany,
	//		data.ResStatus,
	//		strVal(data.ResYear),
	//		strVal(data.ResMonth),
	//		data.TotalEmployee,
	//		data.PurRelatiBank,
	//		strVal(data.TotalLiability),
	//		data.AttStatus,
	//		strVal(data.RelationCode),
	//		strVal(data.RelCustomer),
	//		strVal(data.PortofolioCateg),
	//		strVal(data.AddrPhoneArea),
	//		strVal(data.PhoneAO),
	//		data.Llc_y_n,
	//		strVal(data.AtiLLCBname),
	//		data.AtiLLCAcc,
	//		data.AtiLLCOacc,
	//		data.Fatca_y_n,
	//		strVal(data.ProjectLocate),
	//		strVal(data.MarketingCode),
	//		data.NoAplikasi,
	//		data.TglPencairan,
	//		data.ProsperaID,
	//		data.ErrorLog,
	//		data.BookingLog,
	//	}
	//	writer.Write(record)
	//}
}

func DisburseErrorReportWriter(w io.Writer, datas []*byte) {
	//writer := csv.NewWriter(w)
	//
	//defer writer.Flush()
	//
	//titleRow := []string{
	//
	//	"Account Number Wow Syariah",
	//	"Account Number T24",
	//	"Agent Name",
	//	"Cif Prospera",
	//	"Cif T24",
	//	"Loan Number Qardh",
	//	"ID Total Limit",
	//	"Total Limit",
	//	"OutstandingTotalLimit",
	//	"ID Qardh Limit",
	//	"Qardh Limit",
	//	"Outstanding Qardh Limit",
	//	"ID Murabahah Limit",
	//	"Murabahah Limit",
	//	"Outstanding Murabahah Limit",
	//	"Error Notes",
	//}
	//writer.Write(titleRow)
	//
	//for _, data := range datas {
	//	record := []string{
	//		data.AccountNumberWowSyariah,
	//		data.AccountNumberT24,
	//		data.AgentName,
	//		data.CifProspera,
	//		data.CifT24,
	//		data.LoanNumberQardh,
	//		data.IDTotalLimit,
	//		strVal(data.TotalLimit),
	//		strVal(data.OutstandingTotalLimit),
	//		data.IDQardhLimit,
	//		strVal(data.QardhLimit),
	//		strVal(data.OutstandingQardhLimit),
	//		data.IDMurabahahLimit,
	//		strVal(data.MurabahahLimit),
	//		strVal(data.OutstandingMurabahahLimit),
	//		data.ErrorNotes,
	//	}
	//	writer.Write(record)
	//}
}
