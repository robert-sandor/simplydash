package internal

import (
	"reflect"
	"testing"
)

type testData struct {
	name     string
	initial  []App
	expected []App
	toInsert App
}

func Test_insertOrdered(t *testing.T) {
	cases := []testData{
		{
			name:     "insert into empty list",
			initial:  []App{},
			toInsert: App{Name: "Test"},
			expected: []App{{Name: "Test"}},
		},
		{
			name:     "append to list",
			initial:  []App{{Name: "App"}},
			toInsert: App{Name: "Test"},
			expected: []App{{Name: "App"}, {Name: "Test"}},
		},
		{
			name:     "prepend to list",
			initial:  []App{{Name: "Zig"}},
			toInsert: App{Name: "Test"},
			expected: []App{{Name: "Test"}, {Name: "Zig"}},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			result := insertOrdered(testCase.initial, testCase.toInsert)
			if !reflect.DeepEqual(result, testCase.expected) {
				t.Errorf("Expected %+v - Got %+v", testCase.expected, result)
			}
		})
	}
}
