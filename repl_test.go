package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "ALLCAPS nocaps mIXofCAPS",
			expected: []string{"allcaps", "nocaps", "mixofcaps"},
		},
		{
			input:    "        12345 number test             ",
			expected: []string{"12345", "number", "test"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Length of slice incorrect: %v (actual) %v (expected)",
				len(actual), len(c.expected),
			)
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Expected word does not match. \nExpected: %s\nUnexpected: %s",
					expectedWord, word,
				)
			}
		}
	}
}
