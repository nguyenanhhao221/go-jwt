package settings

type Settings struct {
	PORT         int
	API_V1       string
	Check_Health string
}

var AppSettings *Settings

func init() {
	AppSettings = &Settings{
		PORT:   8080,
		API_V1: "/v1",
	}
	AppSettings.Check_Health = AppSettings.API_V1 + "/health"
}
