package external

import "testing"

func TestChecklyGroup(t *testing.T) {
	data := Group{
		Name:      "foo",
		Namespace: "bar",
		Locations: []string{"basement"},
	}

	testData := checklyGroup(data)

	if testData.Name != data.Name {
		t.Errorf("Expected %s, got %s", data.Name, testData.Name)
	}
}
