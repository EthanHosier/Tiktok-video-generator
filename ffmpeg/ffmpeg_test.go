package ffmpeg

import "testing"

func TestIsValidTime(t *testing.T) {
	testCases := []struct {
		time     string
		expected bool
	}{
		// Valid cases
		{"00:01:57", true},
		{"12:34:56", true},
		{"23:59:59", true},
		{"00:00:00", true},

		// Invalid cases
		{"25:00:00", false},
		{"01:60:00", false},
		{"00:00:60", false},
		{"123:00:00", false},
		{"00:123:00", false},
		{"00:00:123", false},
		{"invalid", false},
		{"00:00", false},
		{"", false},
	}

	for _, tc := range testCases {
		result := isValidTime(tc.time)
		if result != tc.expected {
			t.Errorf("isValidTime(%q) = %v; expected %v", tc.time, result, tc.expected)
		}
	}
}
