package settings

type Settings struct {
	API_V1 string
}

var AppSettings *Settings

func init() {
	AppSettings = &Settings{
		API_V1: "/v1",
	}
}
