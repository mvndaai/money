package money

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
func CurrencyDecimals(currencyCode string) int {
	v, ok := currencyDecimals[currencyCode]
	if !ok {
		return 2
	}
	return v
}
