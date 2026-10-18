package option

import (
	"reflect"
	"testing"
)

func TestNormalizedFillsDefaultsAndClamps(t *testing.T) {
	got := Settings{Groups: []string{" A ", "", "A", "B"}, FailureThreshold: 99, Interval: 1, Timeout: -1}.Normalized()
	want := Settings{
		Groups:           []string{"A", "B"},
		FailureThreshold: 10,
		Interval:         10,
		Timeout:          1000,
		Concurrency:      DefaultConcurrency,
		URL:              DefaultURL,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Normalized() = %+v, want %+v", got, want)
	}
}

func TestValidate(t *testing.T) {
	for name, test := range map[string]struct {
		settings Settings
		valid    bool
	}{
		"disabled without groups": {Settings{}, true},
		"enabled without groups":  {Settings{Enable: true, Groups: []string{" "}}, false},
		"invalid filter":          {Settings{PreferredFilter: "("}, false},
		"invalid expected status": {Settings{ExpectedStatus: "abc"}, false},
		"enabled with lookbehind": {Settings{Enable: true, Groups: []string{"G"}, PreferredFilter: `(?i)(?<![A-Za-z])UK(?![A-Za-z])`}, true},
	} {
		if err := test.settings.Validate(); (err == nil) != test.valid {
			t.Errorf("%s: Validate() = %v", name, err)
		}
	}
}

func TestPreferredFilterIsTokenAware(t *testing.T) {
	filter, err := Settings{PreferredFilter: `(?i)(?<![A-Za-z])UK(?![A-Za-z])`}.CompileFilter()
	if err != nil {
		t.Fatal(err)
	}
	for name, want := range map[string]bool{"UK 01": true, "🇬🇧 uk-2": true, "Ukraine 01": false, "Duke": false} {
		if got, _ := filter.MatchString(name); got != want {
			t.Errorf("MatchString(%q) = %v, want %v", name, got, want)
		}
	}
}
