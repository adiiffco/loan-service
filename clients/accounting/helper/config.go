package helper

type AccountingChoices int

const (
	XERO AccountingChoices = iota
	MYOB
)

var (
	SUPPORTED_ACCOUNTING = map[AccountingChoices]string{
		XERO: "XERO service",
		MYOB: "MYOB service",
	}
)
