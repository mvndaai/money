package money

import (
	"fmt"
	"strings"
)

// 0 Decimal Currency
const (
	CurrencyCodeBIF = "BIF"
	CurrencyCodeCLP = "CLP"
	CurrencyCodeDJF = "DJF"
	CurrencyCodeGNF = "GNF"
	CurrencyCodeISK = "ISK"
	CurrencyCodeJPY = "JPY"
	CurrencyCodeKMF = "KMF"
	CurrencyCodeKRW = "KRW"
	CurrencyCodePYG = "PYG"
	CurrencyCodeRWF = "RWF"
	CurrencyCodeUGX = "UGX"
	CurrencyCodeUYI = "UYI"
	CurrencyCodeVND = "VND"
	CurrencyCodeVUV = "VUV"
	CurrencyCodeXAF = "XAF"
	CurrencyCodeXOF = "XOF"
	CurrencyCodeXPF = "XPF"
)

// 2 Decimal Currency
const (
	CurrencyCodeAED = "AED"
	CurrencyCodeAFN = "AFN"
	CurrencyCodeALL = "ALL"
	CurrencyCodeAMD = "AMD"
	CurrencyCodeANG = "ANG"
	CurrencyCodeAOA = "AOA"
	CurrencyCodeARS = "ARS"
	CurrencyCodeAUD = "AUD"
	CurrencyCodeAWG = "AWG"
	CurrencyCodeAZN = "AZN"
	CurrencyCodeBAM = "BAM"
	CurrencyCodeBBD = "BBD"
	CurrencyCodeBDT = "BDT"
	CurrencyCodeBGN = "BGN"
	CurrencyCodeBMD = "BMD"
	CurrencyCodeBND = "BND"
	CurrencyCodeBOB = "BOB"
	CurrencyCodeBOV = "BOV"
	CurrencyCodeBRL = "BRL"
	CurrencyCodeBSD = "BSD"
	CurrencyCodeBTN = "BTN"
	CurrencyCodeBWP = "BWP"
	CurrencyCodeBYN = "BYN"
	CurrencyCodeBZD = "BZD"
	CurrencyCodeCAD = "CAD"
	CurrencyCodeCDF = "CDF"
	CurrencyCodeCHE = "CHE"
	CurrencyCodeCHF = "CHF"
	CurrencyCodeCHW = "CHW"
	CurrencyCodeCOP = "COP"
	CurrencyCodeCOU = "COU"
	CurrencyCodeCRC = "CRC"
	CurrencyCodeCUC = "CUC"
	CurrencyCodeCUP = "CUP"
	CurrencyCodeCVE = "CVE"
	CurrencyCodeCZK = "CZK"
	CurrencyCodeDKK = "DKK"
	CurrencyCodeDOP = "DOP"
	CurrencyCodeDZD = "DZD"
	CurrencyCodeEGP = "EGP"
	CurrencyCodeERN = "ERN"
	CurrencyCodeETB = "ETB"
	CurrencyCodeEUR = "EUR"
	CurrencyCodeFJD = "FJD"
	CurrencyCodeFKP = "FKP"
	CurrencyCodeGBP = "GBP"
	CurrencyCodeGEL = "GEL"
	CurrencyCodeGHS = "GHS"
	CurrencyCodeGIP = "GIP"
	CurrencyCodeGMD = "GMD"
	CurrencyCodeGTQ = "GTQ"
	CurrencyCodeGYD = "GYD"
	CurrencyCodeHKD = "HKD"
	CurrencyCodeHNL = "HNL"
	CurrencyCodeHTG = "HTG"
	CurrencyCodeHUF = "HUF"
	CurrencyCodeIDR = "IDR"
	CurrencyCodeILS = "ILS"
	CurrencyCodeINR = "INR"
	CurrencyCodeIRR = "IRR"
	CurrencyCodeJMD = "JMD"
	CurrencyCodeKES = "KES"
	CurrencyCodeKGS = "KGS"
	CurrencyCodeKHR = "KHR"
	CurrencyCodeKPW = "KPW"
	CurrencyCodeKYD = "KYD"
	CurrencyCodeKZT = "KZT"
	CurrencyCodeLAK = "LAK"
	CurrencyCodeLBP = "LBP"
	CurrencyCodeLKR = "LKR"
	CurrencyCodeLRD = "LRD"
	CurrencyCodeLSL = "LSL"
	CurrencyCodeMAD = "MAD"
	CurrencyCodeMDL = "MDL"
	CurrencyCodeMGA = "MGA"
	CurrencyCodeMKD = "MKD"
	CurrencyCodeMMK = "MMK"
	CurrencyCodeMNT = "MNT"
	CurrencyCodeMOP = "MOP"
	CurrencyCodeMRU = "MRU"
	CurrencyCodeMUR = "MUR"
	CurrencyCodeMVR = "MVR"
	CurrencyCodeMWK = "MWK"
	CurrencyCodeMXN = "MXN"
	CurrencyCodeMXV = "MXV"
	CurrencyCodeMYR = "MYR"
	CurrencyCodeMZN = "MZN"
	CurrencyCodeNAD = "NAD"
	CurrencyCodeNGN = "NGN"
	CurrencyCodeNIO = "NIO"
	CurrencyCodeNOK = "NOK"
	CurrencyCodeNPR = "NPR"
	CurrencyCodeNZD = "NZD"
	CurrencyCodePAB = "PAB"
	CurrencyCodePEN = "PEN"
	CurrencyCodePGK = "PGK"
	CurrencyCodePHP = "PHP"
	CurrencyCodePKR = "PKR"
	CurrencyCodePLN = "PLN"
	CurrencyCodeQAR = "QAR"
	CurrencyCodeRON = "RON"
	CurrencyCodeRSD = "RSD"
	CurrencyCodeCNY = "CNY"
	CurrencyCodeRUB = "RUB"
	CurrencyCodeSAR = "SAR"
	CurrencyCodeSBD = "SBD"
	CurrencyCodeSCR = "SCR"
	CurrencyCodeSDG = "SDG"
	CurrencyCodeSEK = "SEK"
	CurrencyCodeSGD = "SGD"
	CurrencyCodeSHP = "SHP"
	CurrencyCodeSLE = "SLE"
	CurrencyCodeSLL = "SLL"
	CurrencyCodeSOS = "SOS"
	CurrencyCodeSRD = "SRD"
	CurrencyCodeSSP = "SSP"
	CurrencyCodeSTN = "STN"
	CurrencyCodeSVC = "SVC"
	CurrencyCodeSYP = "SYP"
	CurrencyCodeSZL = "SZL"
	CurrencyCodeTHB = "THB"
	CurrencyCodeTJS = "TJS"
	CurrencyCodeTMT = "TMT"
	CurrencyCodeTOP = "TOP"
	CurrencyCodeTRY = "TRY"
	CurrencyCodeTTD = "TTD"
	CurrencyCodeTWD = "TWD"
	CurrencyCodeTZS = "TZS"
	CurrencyCodeUAH = "UAH"
	CurrencyCodeUSD = "USD"
	CurrencyCodeUSN = "USN"
	CurrencyCodeUYU = "UYU"
	CurrencyCodeUZS = "UZS"
	CurrencyCodeVED = "VED"
	CurrencyCodeVES = "VES"
	CurrencyCodeWST = "WST"
	CurrencyCodeXCD = "XCD"
	CurrencyCodeYER = "YER"
	CurrencyCodeZAR = "ZAR"
	CurrencyCodeZMW = "ZMW"
	CurrencyCodeZWL = "ZWL"
)

// 3 Decimal Currency
const (
	CurrencyCodeBHD = "BHD"
	CurrencyCodeIQD = "IQD"
	CurrencyCodeJOD = "JOD"
	CurrencyCodeKWD = "KWD"
	CurrencyCodeLYD = "LYD"
	CurrencyCodeOMR = "OMR"
	CurrencyCodeTND = "TND"
)

// 4 Decimal Currency
const (
	CurrencyCodeCLF = "CLF"
	CurrencyCodeUYW = "UYW"
)

var currencyDecimals = map[string]int{
	CurrencyCodeBIF: 0,
	CurrencyCodeCLP: 0,
	CurrencyCodeDJF: 0,
	CurrencyCodeGNF: 0,
	CurrencyCodeISK: 0,
	CurrencyCodeJPY: 0,
	CurrencyCodeKMF: 0,
	CurrencyCodeKRW: 0,
	CurrencyCodePYG: 0,
	CurrencyCodeRWF: 0,
	CurrencyCodeUGX: 0,
	CurrencyCodeUYI: 0,
	CurrencyCodeVND: 0,
	CurrencyCodeVUV: 0,
	CurrencyCodeXAF: 0,
	CurrencyCodeXOF: 0,
	CurrencyCodeXPF: 0,

	CurrencyCodeAED: 2,
	CurrencyCodeAFN: 2,
	CurrencyCodeALL: 2,
	CurrencyCodeAMD: 2,
	CurrencyCodeANG: 2,
	CurrencyCodeAOA: 2,
	CurrencyCodeARS: 2,
	CurrencyCodeAUD: 2,
	CurrencyCodeAWG: 2,
	CurrencyCodeAZN: 2,
	CurrencyCodeBAM: 2,
	CurrencyCodeBBD: 2,
	CurrencyCodeBDT: 2,
	CurrencyCodeBGN: 2,
	CurrencyCodeBMD: 2,
	CurrencyCodeBND: 2,
	CurrencyCodeBOB: 2,
	CurrencyCodeBOV: 2,
	CurrencyCodeBRL: 2,
	CurrencyCodeBSD: 2,
	CurrencyCodeBTN: 2,
	CurrencyCodeBWP: 2,
	CurrencyCodeBYN: 2,
	CurrencyCodeBZD: 2,
	CurrencyCodeCAD: 2,
	CurrencyCodeCDF: 2,
	CurrencyCodeCHE: 2,
	CurrencyCodeCHF: 2,
	CurrencyCodeCHW: 2,
	CurrencyCodeCOP: 2,
	CurrencyCodeCOU: 2,
	CurrencyCodeCRC: 2,
	CurrencyCodeCUC: 2,
	CurrencyCodeCUP: 2,
	CurrencyCodeCVE: 2,
	CurrencyCodeCZK: 2,
	CurrencyCodeDKK: 2,
	CurrencyCodeDOP: 2,
	CurrencyCodeDZD: 2,
	CurrencyCodeEGP: 2,
	CurrencyCodeERN: 2,
	CurrencyCodeETB: 2,
	CurrencyCodeEUR: 2,
	CurrencyCodeFJD: 2,
	CurrencyCodeFKP: 2,
	CurrencyCodeGBP: 2,
	CurrencyCodeGEL: 2,
	CurrencyCodeGHS: 2,
	CurrencyCodeGIP: 2,
	CurrencyCodeGMD: 2,
	CurrencyCodeGTQ: 2,
	CurrencyCodeGYD: 2,
	CurrencyCodeHKD: 2,
	CurrencyCodeHNL: 2,
	CurrencyCodeHTG: 2,
	CurrencyCodeHUF: 2,
	CurrencyCodeIDR: 2,
	CurrencyCodeILS: 2,
	CurrencyCodeINR: 2,
	CurrencyCodeIRR: 2,
	CurrencyCodeJMD: 2,
	CurrencyCodeKES: 2,
	CurrencyCodeKGS: 2,
	CurrencyCodeKHR: 2,
	CurrencyCodeKPW: 2,
	CurrencyCodeKYD: 2,
	CurrencyCodeKZT: 2,
	CurrencyCodeLAK: 2,
	CurrencyCodeLBP: 2,
	CurrencyCodeLKR: 2,
	CurrencyCodeLRD: 2,
	CurrencyCodeLSL: 2,
	CurrencyCodeMAD: 2,
	CurrencyCodeMDL: 2,
	CurrencyCodeMGA: 2,
	CurrencyCodeMKD: 2,
	CurrencyCodeMMK: 2,
	CurrencyCodeMNT: 2,
	CurrencyCodeMOP: 2,
	CurrencyCodeMRU: 2,
	CurrencyCodeMUR: 2,
	CurrencyCodeMVR: 2,
	CurrencyCodeMWK: 2,
	CurrencyCodeMXN: 2,
	CurrencyCodeMXV: 2,
	CurrencyCodeMYR: 2,
	CurrencyCodeMZN: 2,
	CurrencyCodeNAD: 2,
	CurrencyCodeNGN: 2,
	CurrencyCodeNIO: 2,
	CurrencyCodeNOK: 2,
	CurrencyCodeNPR: 2,
	CurrencyCodeNZD: 2,
	CurrencyCodePAB: 2,
	CurrencyCodePEN: 2,
	CurrencyCodePGK: 2,
	CurrencyCodePHP: 2,
	CurrencyCodePKR: 2,
	CurrencyCodePLN: 2,
	CurrencyCodeQAR: 2,
	CurrencyCodeRON: 2,
	CurrencyCodeRSD: 2,
	CurrencyCodeCNY: 2,
	CurrencyCodeRUB: 2,
	CurrencyCodeSAR: 2,
	CurrencyCodeSBD: 2,
	CurrencyCodeSCR: 2,
	CurrencyCodeSDG: 2,
	CurrencyCodeSEK: 2,
	CurrencyCodeSGD: 2,
	CurrencyCodeSHP: 2,
	CurrencyCodeSLE: 2,
	CurrencyCodeSLL: 2,
	CurrencyCodeSOS: 2,
	CurrencyCodeSRD: 2,
	CurrencyCodeSSP: 2,
	CurrencyCodeSTN: 2,
	CurrencyCodeSVC: 2,
	CurrencyCodeSYP: 2,
	CurrencyCodeSZL: 2,
	CurrencyCodeTHB: 2,
	CurrencyCodeTJS: 2,
	CurrencyCodeTMT: 2,
	CurrencyCodeTOP: 2,
	CurrencyCodeTRY: 2,
	CurrencyCodeTTD: 2,
	CurrencyCodeTWD: 2,
	CurrencyCodeTZS: 2,
	CurrencyCodeUAH: 2,
	CurrencyCodeUSD: 2,
	CurrencyCodeUSN: 2,
	CurrencyCodeUYU: 2,
	CurrencyCodeUZS: 2,
	CurrencyCodeVED: 2,
	CurrencyCodeVES: 2,
	CurrencyCodeWST: 2,
	CurrencyCodeXCD: 2,
	CurrencyCodeYER: 2,
	CurrencyCodeZAR: 2,
	CurrencyCodeZMW: 2,
	CurrencyCodeZWL: 2,

	CurrencyCodeBHD: 3,
	CurrencyCodeIQD: 3,
	CurrencyCodeJOD: 3,
	CurrencyCodeKWD: 3,
	CurrencyCodeLYD: 3,
	CurrencyCodeOMR: 3,
	CurrencyCodeTND: 3,

	CurrencyCodeCLF: 4,
	CurrencyCodeUYW: 4,
}

// Decimal places based on https://en.wikipedia.org/wiki/ISO_4217
func CurrencyDecimals(currencyCode string) (int, error) {
	v, ok := currencyDecimals[strings.ToUpper(currencyCode)]
	if !ok {
		return 0, fmt.Errorf("unknown currency code '%s'", currencyCode)
	}
	return v, nil
}
