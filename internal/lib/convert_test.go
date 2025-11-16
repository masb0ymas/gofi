package lib

import (
	"testing"
	"time"
)

func TestStringToInt32(t *testing.T) {
	stringTo := StringTo{}

	tests := []struct {
		name         string
		input        string
		defaultValue []int32
		expected     int32
	}{
		{
			name:     "valid number",
			input:    "123",
			expected: 123,
		},
		{
			name:         "valid number with default",
			input:        "456",
			defaultValue: []int32{999},
			expected:     456,
		},
		{
			name:         "invalid number uses default",
			input:        "abc",
			defaultValue: []int32{999},
			expected:     999,
		},
		{
			name:     "invalid number uses zero default",
			input:    "abc",
			expected: 0,
		},
		{
			name:     "empty string uses zero default",
			input:    "",
			expected: 0,
		},
		{
			name:     "max int32",
			input:    "2147483647",
			expected: 2147483647,
		},
		{
			name:     "min int32",
			input:    "-2147483648",
			expected: -2147483648,
		},
		{
			name:     "overflow uses default",
			input:    "2147483648",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringTo.ToInt32(tt.input, tt.defaultValue...)
			if got != tt.expected {
				t.Errorf("StringToInt32() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestStringToInt64(t *testing.T) {
	stringTo := StringTo{}

	tests := []struct {
		name         string
		input        string
		defaultValue []int64
		expected     int64
	}{
		{
			name:     "valid number",
			input:    "123",
			expected: 123,
		},
		{
			name:         "valid number with default",
			input:        "456",
			defaultValue: []int64{999},
			expected:     456,
		},
		{
			name:         "invalid number uses default",
			input:        "abc",
			defaultValue: []int64{999},
			expected:     999,
		},
		{
			name:     "invalid number uses zero default",
			input:    "abc",
			expected: 0,
		},
		{
			name:     "empty string uses zero default",
			input:    "",
			expected: 0,
		},
		{
			name:     "max int64",
			input:    "9223372036854775807",
			expected: 9223372036854775807,
		},
		{
			name:     "min int64",
			input:    "-9223372036854775808",
			expected: -9223372036854775808,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringTo.ToInt64(tt.input, tt.defaultValue...)
			if got != tt.expected {
				t.Errorf("StringToInt64() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestStringToTime(t *testing.T) {
	stringTo := StringTo{}

	tests := []struct {
		name     string
		input    string
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: time.Time{},
			wantErr:  false,
		},
		{
			name:     "RFC3339 format",
			input:    "2023-05-15T14:30:00Z",
			expected: time.Date(2023, 5, 15, 14, 30, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "ISO8601 format",
			input:    "2023-05-15T14:30:00+00:00",
			expected: time.Date(2023, 5, 15, 14, 30, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "datetime format",
			input:    "2023-05-15 14:30:05",
			expected: time.Date(2023, 5, 15, 14, 30, 5, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "date only format",
			input:    "2023-05-15",
			expected: time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "invalid format",
			input:    "2023-13-45",
			expected: time.Time{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stringTo.ToTime(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.expected) {
				t.Errorf("StringToTime() = %v, want %v", got, tt.expected)
			}
		})
	}
}
