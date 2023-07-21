package settings

import "testing"

func TestAppSettings(t *testing.T) {
	expectedAPIV1 := "/v1"

	if AppSettings.API_V1 != expectedAPIV1 {
		t.Errorf("Expected API_V1 to be %s, but got %s", expectedAPIV1, AppSettings.API_V1)
	}

	// Test Check_Health
	expectedCheckHealth := "/v1/health"
	if AppSettings.Check_Health != expectedCheckHealth {
		t.Errorf("Expected Check_Health to be %s, but got %s", expectedCheckHealth, AppSettings.Check_Health)
	}
}
