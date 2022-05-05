package lpfr

type StatusResponse struct {
	IsPinRequired bool `json:"isPinRequired"`
	AuditRequired bool `json:"auditRequired"`
}

type EnvironmentParamsResponse struct {
	OrganizationName string `json:"organizationName"`
	ServerTimeZone   string `json:"serverTimeZone"`
	Street           string `json:"street"`
	City             string `json:"city"`
	Country          string `json:"country"`
	Endpoints        struct {
		TaxpayerAdminPortal string `json:"taxpayerAdminPortal"`
		TaxCoreApi          string `json:"taxCoreApi"`
		Vsdc                string `json:"vsdc"`
		Root                string `json:"root"`
	} `json:"endpoints"`
	EnvironmentName    string   `json:"environmentName"`
	Logo               string   `json:"logo"`
	NtpServer          string   `json:"ntpServer"`
	SupportedLanguages []string `json:"supportedLanguages"`
}
