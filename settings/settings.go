package settings

type Settings struct {
	PORT                 int
	API_V1               string
	Check_Health         string
	Account_Route        string
	Create_Account_Route string
}

var AppSettings *Settings

func init() {
	AppSettings = &Settings{
		PORT:                 8080,
		API_V1:               "/v1",
		Check_Health:         "/health",
		Account_Route:        "/account/{accountId}",
		Create_Account_Route: "/account/create",
	}
}
