package external

import (
	"testing"
)

func TestChecklyCheck(t *testing.T) {

	data := Check{
		Name:            "foo",
		Namespace:       "bar",
		Frequency:       15,
		MaxResponseTime: 2000,
		Locations:       []string{"basement"},
		Endpoint:        "https://foo.bar/baz",
		SuccessCode:     "200",
	}

	testData := checklyCheck(data)

	if testData.Name != data.Name {
		t.Errorf("Expected %s, got %s", data.Name, testData.Name)
	}

	if testData.Frequency != data.Frequency {
		t.Errorf("Expected %d, got %d", data.Frequency, testData.Frequency)
	}

	if testData.MaxResponseTime != data.MaxResponseTime {
		t.Errorf("Expected %d, got %d", data.Frequency, testData.Frequency)
	}

	data = Check{
		Name:        "foo",
		Namespace:   "bar",
		Endpoint:    "https://foo.bar/baz",
		SuccessCode: "200",
	}

	testData = checklyCheck(data)

	if testData.Frequency != 5 {
		t.Errorf("Expected %d, got %d", data.Frequency, testData.Frequency)
	}

	if testData.MaxResponseTime != 15000 {
		t.Errorf("Expected %d, got %d", data.Frequency, testData.Frequency)
	}

	return
}
