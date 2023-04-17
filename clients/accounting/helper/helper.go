package helper

import (
	myob "loanapp/clients/accounting/MYOB"
	loanInterface "loanapp/clients/accounting/interface"
	"loanapp/clients/accounting/xero"
)

func GetSupportedAccountingService() map[AccountingChoices]string {
	return SUPPORTED_ACCOUNTING
}

func GetAccountingService(choice AccountingChoices) loanInterface.Accounter {
	switch choice {
	case XERO:
		return &xero.XeroClient{}
	case MYOB:
		return &myob.MyobClient{}
	default:
		return nil
	}
}
