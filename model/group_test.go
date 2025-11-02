package model

import (
	"testing"
)

func TestGroupStructFields(t *testing.T) {
	group := &Group{
		GroupName: "Test Group",
		GroupURL:  "https://example.com",
		GroupLogo: "logo.png",
		Verified:  true,
		Active:    true,
		Handle:    "testgroup",
	}

	if group.GroupName != "Test Group" {
		t.Errorf("GroupName = %s, want Test Group", group.GroupName)
	}
	if group.GroupURL != "https://example.com" {
		t.Errorf("GroupURL = %s, want https://example.com", group.GroupURL)
	}
	if group.GroupLogo != "logo.png" {
		t.Errorf("GroupLogo = %s, want logo.png", group.GroupLogo)
	}
	if !group.Verified {
		t.Error("Verified should be true")
	}
	if !group.Active {
		t.Error("Active should be true")
	}
	if group.Handle != "testgroup" {
		t.Errorf("Handle = %s, want testgroup", group.Handle)
	}
}

func TestGroupWithIbans(t *testing.T) {
	group := &Group{
		GroupName: "Test Group",
		Handle:    "testgroup",
		Ibans: []Iban{
			{
				Text:   "TR320010009999901234567890",
				Handle: "iban1",
			},
			{
				Text:   "TR420010009999901234567891",
				Handle: "iban2",
			},
		},
	}

	if len(group.Ibans) != 2 {
		t.Errorf("Expected 2 IBANs, got %d", len(group.Ibans))
	}

	if group.Ibans[0].Text != "TR320010009999901234567890" {
		t.Errorf("First IBAN text = %s, want TR320010009999901234567890", group.Ibans[0].Text)
	}
	if group.Ibans[1].Handle != "iban2" {
		t.Errorf("Second IBAN handle = %s, want iban2", group.Ibans[1].Handle)
	}
}

func TestGroupDefaultValues(t *testing.T) {
	group := &Group{}

	if group.GroupName != "" {
		t.Errorf("Default GroupName should be empty, got %s", group.GroupName)
	}
	if group.Verified {
		t.Error("Default Verified should be false")
	}
	if group.Active {
		t.Error("Default Active should be false")
	}
}
