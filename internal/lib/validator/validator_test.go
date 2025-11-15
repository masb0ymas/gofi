package validator

import (
	"strings"
	"testing"
	"time"
)

const (
	fieldName = "test"
)

type validatorTestTable struct {
	name  string
	value interface{}
	want  bool
	setup func() *MapValidator
}

type messageValidatorTestTable struct {
	name            string
	value           interface{}
	expectedMessage string
	setup           func() *MapValidator
}

func validateTestData(t *testing.T, testTable []validatorTestTable, v *MapValidator) {
	t.Helper()

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			_, passes := v.Validate(map[string]interface{}{
				fieldName: tt.value,
			})
			if passes != tt.want {
				t.Errorf("Expected passes to be %v, got %v", tt.want, passes)
			}
		})
	}
}

func validateTestDataWithSetup(t *testing.T, testTable []validatorTestTable) {
	t.Helper()
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.setup()
			_, passes := v.Validate(map[string]interface{}{
				fieldName: tt.value,
			})
			if passes != tt.want {
				t.Errorf("Expected passes to be %v, got %v", tt.want, passes)
			}
		})
	}
}

func validateTestDataMessage(t *testing.T, testTable []messageValidatorTestTable) {
	t.Helper()
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.setup()
			message, _ := v.Validate(map[string]interface{}{
				fieldName: tt.value,
			})
			for _, errorMessages := range message {
				if !strings.Contains(errorMessages[0], tt.expectedMessage) {
					t.Errorf("Expected %s to be in %s", tt.expectedMessage, errorMessages[0])
				}
				break
			}

		})
	}
}

func Test_ValidationRules(t *testing.T) {
	t.Run("Required", func(t *testing.T) {
		exStr := "test"
		v := NewMapValidator()
		v.Field(fieldName).Required()

		testTable := []validatorTestTable{
			{
				name:  "Should pass - string",
				value: "test",
				want:  true,
			},
			{
				name:  "Should fail - string",
				value: "",
				want:  false,
			},
			{
				name:  "Should pass - pointer [string]",
				value: &exStr,
				want:  true,
			},
			{
				name:  "Should pass - pointer [slice]",
				value: &[]string{"test"},
				want:  true,
			},
			{
				name: "Should pass - pointer [map]",
				value: &map[string]interface{}{
					"test": "test",
				},
				want: true,
			},
			{
				name:  "Should fail - pointer",
				value: nil,
				want:  false,
			},
			{
				name:  "Should pass - slice",
				value: []string{"test"},
				want:  true,
			},
			{
				name: "Should pass - map",
				value: map[string]interface{}{
					"test": "test",
				},
				want: true,
			},
		}

		validateTestData(t, testTable, v)
	})

	t.Run("Alpha", func(t *testing.T) {

		v := NewMapValidator()
		v.Field(fieldName).Alpha()

		testTable := []validatorTestTable{
			{
				name:  "Should pass",
				value: "abcdef",
				want:  true,
			},
			{
				name:  "Should fail - includes numbers",
				value: "ab1def",
				want:  false,
			},
		}
		validateTestData(t, testTable, v)
	})

	t.Run("Num", func(t *testing.T) {

		v := NewMapValidator()
		v.Field(fieldName).Num()

		testTable := []validatorTestTable{
			{
				name:  "Should pass - integers",
				value: 12,
				want:  true,
			},
			{
				name:  "Should pass - floats",
				value: 12.2,
				want:  true,
			},
			{
				name:  "Should fail",
				value: "nan",
				want:  false,
			},
		}
		validateTestData(t, testTable, v)
	})

	t.Run("String", func(t *testing.T) {

		v := NewMapValidator()
		v.Field(fieldName).String()

		testTable := []validatorTestTable{
			{
				name:  "Should pass - pure string",
				value: "test",
				want:  true,
			},
			{
				name:  "Should pass - string of integers",
				value: "123",
				want:  true,
			},
			{
				name:  "Should fail",
				value: 123,
				want:  false,
			},
		}
		validateTestData(t, testTable, v)
	})

	t.Run("Regex", func(t *testing.T) {
		v := NewMapValidator()
		v.Field(fieldName).Regex("^[0-9]+$")

		testTable := []validatorTestTable{
			{
				name:  "Should pass",
				value: "123",
				want:  true,
			},
			{
				name:  "Should fail",
				value: "test123",
				want:  false,
			},
		}
		validateTestData(t, testTable, v)
	})

	t.Run("Email", func(t *testing.T) {
		v := NewMapValidator()
		v.Field(fieldName).Email()

		testTable := []validatorTestTable{
			{
				name:  "Should pass",
				value: "test@test.com",
				want:  true,
			},
			{
				name:  "Should fail - invalid email",
				value: "test@test",
				want:  false,
			},
			{
				name:  "Should fail - not a string",
				value: 123,
				want:  false,
			},
		}
		validateTestData(t, testTable, v)
	})

	t.Run("Bool", func(t *testing.T) {

		v := NewMapValidator()
		v.Field(fieldName).Bool()

		testTable := []validatorTestTable{
			{
				name:  "should pass - true",
				value: true,
				want:  true,
			},
			{
				name:  "should pass - false",
				value: false,
				want:  true,
			},
			{
				name:  "should fail",
				value: "true",
				want:  false,
			},
		}
		validateTestData(t, testTable, v)
	})

	t.Run("Date", func(t *testing.T) {

		v := NewMapValidator()
		v.Field(fieldName).Date()

		testTable := []validatorTestTable{
			{
				name:  "should pass- valid date format",
				value: "2024-09-16T08:45:30.123Z",
				want:  true,
			},
			{
				name:  "should fail - not datetime string",
				value: time.Now(),
				want:  false,
			},
			{
				name:  "should fail - invalid date format",
				value: "2006-01-02T15:04:05Z07:00",
				want:  false,
			},
		}
		validateTestData(t, testTable, v)
	})

	t.Run("Min", func(t *testing.T) {

		v := NewMapValidator()
		v.Field(fieldName).Min(10)

		testTable := []validatorTestTable{
			{
				name:  "should pass - integer",
				value: 11,
				want:  true,
			},
			{
				name:  "should pass - float",
				value: 11.11,
				want:  true,
			},
			{
				name:  "should fail - less than min",
				value: 9,
				want:  false,
			},
		}
		validateTestData(t, testTable, v)
	})

	t.Run("Max", func(t *testing.T) {

		v := NewMapValidator()
		v.Field(fieldName).Max(10)

		testTable := []validatorTestTable{
			{
				name:  "should pass - integer",
				value: 9,
				want:  true,
			},
			{
				name:  "should pass - float",
				value: 9.99,
				want:  true,
			},
			{
				name:  "should fail - more than max",
				value: 11,
				want:  false,
			},
		}
		validateTestData(t, testTable, v)
	})

	t.Run("MinS", func(t *testing.T) {

		v := NewMapValidator()
		v.Field(fieldName).MinS(4)

		testTable := []validatorTestTable{
			{
				name:  "should pass - at boundary",
				value: "1234",
				want:  true,
			},
			{
				name:  "should pass - above boundary",
				value: "12345",
				want:  true,
			},
			{
				name:  "should fail - below boundary",
				value: "123",
				want:  false,
			},
		}
		validateTestData(t, testTable, v)
	})

	t.Run("MaxS", func(t *testing.T) {

		v := NewMapValidator()
		v.Field(fieldName).MaxS(4)

		testTable := []validatorTestTable{
			{
				name:  "should pass - at boundary",
				value: "1234",
				want:  true,
			},
			{
				name:  "should pass - below boundary",
				value: "123",
				want:  true,
			},
			{
				name:  "should fail - above boundary",
				value: "12345",
				want:  false,
			},
		}
		validateTestData(t, testTable, v)
	})

	t.Run("Within", func(t *testing.T) {

		v := NewMapValidator()
		v.Field(fieldName).Within(1, 2, 3, 4)

		testTable := []validatorTestTable{
			{
				name:  "should pass",
				value: 1,
				want:  true,
			},
			{
				name:  "should fail",
				value: 5,
				want:  false,
			},
		}
		validateTestData(t, testTable, v)
	})

	t.Run("WithinS", func(t *testing.T) {

		v := NewMapValidator()
		v.Field(fieldName).WithinS("a", "b", "c")

		testTable := []validatorTestTable{
			{
				name:  "should pass",
				value: "b",
				want:  true,
			},
			{
				name:  "should fail",
				value: "d",
				want:  false,
			},
		}
		validateTestData(t, testTable, v)
	})

	t.Run("Base64", func(t *testing.T) {

		v := NewMapValidator()
		v.Field(fieldName).Base64()

		testTable := []validatorTestTable{
			{
				name:  "should pass",
				value: "dGVzdA==",
				want:  true,
			},
			{
				name:  "should fail",
				value: "dGVzdA",
				want:  false,
			},
		}
		validateTestData(t, testTable, v)
	})

	t.Run("MinRune", func(t *testing.T) {

		v := NewMapValidator()
		v.Field(fieldName).MinRune(3)

		testTable := []validatorTestTable{
			{
				name:  "should pass - above boundary",
				value: "世界世界",
				want:  true,
			},
			{
				name:  "should pass - at boundary",
				value: "世界世",
				want:  true,
			},
			{
				name:  "should fail - below boundary",
				value: "世界",
				want:  false,
			},
		}
		validateTestData(t, testTable, v)
	})

	t.Run("MaxRune", func(t *testing.T) {

		v := NewMapValidator()
		v.Field(fieldName).MaxRune(3)

		testTable := []validatorTestTable{
			{
				name:  "should pass - below boundary",
				value: "世界",
				want:  true,
			},
			{
				name:  "should pass - at boundary",
				value: "世界世",
				want:  true,
			},
			{
				name:  "should fail - above boundary",
				value: "世界世界",
				want:  false,
			},
		}
		validateTestData(t, testTable, v)
	})

	t.Run("UUID", func(t *testing.T) {

		v := NewMapValidator()
		v.Field(fieldName).UUID()

		testTable := []validatorTestTable{
			{
				name:  "should pass",
				value: "e4eaaaf2-d142-11e1-b3e4-080027620cdd",
				want:  true,
			},
			{
				name:  "should fail",
				value: "e4eaaaf2",
				want:  false,
			},
		}
		validateTestData(t, testTable, v)
	})

	t.Run("Slice", func(t *testing.T) {

		exSlice := []int{1, 2, 3}
		exBadSlice := []string{"1", "2", "3"}
		v := NewMapValidator()
		v.Field(fieldName).Slice(func(v *FieldValidator) {
			v.Num()
		})

		testTable := []validatorTestTable{
			{
				name:  "should pass - slice of integers",
				value: exSlice,
				want:  true,
			},
			{
				name:  "should fail - slice of strings",
				value: exBadSlice,
				want:  false,
			},
			{
				name:  "should fail - not a slice",
				value: "{}",
				want:  false,
			},
			{
				name:  "should pass - empty value",
				value: nil,
				want:  true,
			},
			{
				name:  "should pass - inner slice is valid",
				value: &exSlice,
				want:  true,
			},
			{
				name:  "should fail - inner slice is invalid",
				value: &exBadSlice,
				want:  false,
			},
		}
		validateTestData(t, testTable, v)
	})

	t.Run("MaxLen", func(t *testing.T) {

		exSlice := []int{1, 2, 3}
		exBadSlice := []int{1, 2, 3, 4, 5, 6}
		v := NewMapValidator()
		v.Field(fieldName).Slice(func(v *FieldValidator) {
			v.Num()
		}).MaxLen(3)

		testTable := []validatorTestTable{
			{
				name:  "should pass - slice of integers",
				value: exSlice,
				want:  true,
			},
			{
				name:  "should fail - slice of integers",
				value: exBadSlice,
				want:  false,
			},
		}
		validateTestData(t, testTable, v)
	})

	t.Run("AnySlice", func(t *testing.T) {

		v := NewMapValidator()
		v.Field(fieldName).AnySlice()

		testTable := []validatorTestTable{
			{
				name:  "should pass - slice of integers",
				value: []int{1, 2, 3},
				want:  true,
			},
			{
				name:  "should pass - slice of strings",
				value: []string{"1", "2", "3"},
				want:  true,
			},
			{
				name:  "should fail - not a slice",
				value: "{}",
				want:  false,
			},
		}
		validateTestData(t, testTable, v)
	})

	t.Run("Map", func(t *testing.T) {

		exMap := map[string]interface{}{
			"status":  true,
			"message": "test",
		}
		v := NewMapValidator()
		v.Field(fieldName).Map(func(v *MapValidator) {
			v.Field("status").Bool()
			v.Field("message").Required().String()
		})

		testTable := []validatorTestTable{
			{
				name:  "should pass - matches schema",
				value: exMap,
				want:  true,
			},
			{
				name: "should pass - accepted missing field",
				value: map[string]interface{}{
					"message": "status is not required",
				},
				want: true,
			},
			{
				name: "should fail - missing field",
				value: map[string]interface{}{
					"status": true,
				},
				want: false,
			},
			{
				name: "should fail - wrong type",
				value: map[string]interface{}{
					"status":  "true",
					"message": "test",
				},
				want: false,
			},
			{
				name:  "should pass - nil value",
				value: nil,
				want:  true,
			},
			{
				name:  "should pass - inner map is valid",
				value: &exMap,
				want:  true,
			},
		}
		validateTestData(t, testTable, v)
	})
}

func Test_NestedMaps(t *testing.T) {
	testTable := []validatorTestTable{
		{
			name: "should pass - depth 1",
			value: map[string]interface{}{
				"year": 2024,
			},
			setup: func() *MapValidator {
				v := NewMapValidator()
				v.Field(fieldName).Map(func(v *MapValidator) {
					v.Field("year").Required().Num()
				})
				return v
			},
			want: true,
		},
		{
			name: "should fail - depth 1 invalid child",
			value: map[string]interface{}{
				"year": "2024",
			},
			setup: func() *MapValidator {
				v := NewMapValidator()
				v.Field(fieldName).Map(func(v *MapValidator) {
					v.Field("year").Required().Num()
				})
				return v
			},
			want: false,
		},
		{
			name:  "should pass - no required fields",
			value: map[string]interface{}{},
			setup: func() *MapValidator {
				v := NewMapValidator()
				v.Field(fieldName).Map(func(v *MapValidator) {
					v.Field("year").Num()
					v.Field("month").Num()
				})
				return v
			},
			want: true,
		},
		{
			name: "should pass - depth 3",
			value: map[string]interface{}{
				"top": map[string]interface{}{
					"middle": map[string]interface{}{
						"bottom": []string{"e4eaaaf2-d142-11e1-b3e4-080027620cdd", "e4eaaaf2-d142-11e1-b3e4-080027620cdd"},
					},
				},
			},
			setup: func() *MapValidator {
				v := NewMapValidator()
				v.Field(fieldName).Map(func(v *MapValidator) {
					v.Field("top").Map(func(v *MapValidator) {
						v.Field("middle").Map(func(v *MapValidator) {
							v.Field("bottom").Required().Slice(func(v *FieldValidator) {
								v.UUID()
							})
						})
					})
				})
				return v
			},
			want: true,
		},
		{
			name: "should pass - depth 2, with two maps",
			value: map[string]interface{}{
				"data": map[string]interface{}{
					"status": true,
				},
				"requestData": map[string]interface{}{
					"ip":        "127.0.0.1",
					"userAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
				},
			},
			setup: func() *MapValidator {
				v := NewMapValidator()
				v.Field(fieldName).Map(func(v *MapValidator) {
					v.Field("data").Map(func(v *MapValidator) {
						v.Field("status").Required().Bool()
					})
					v.Field("requestData").Map(func(v *MapValidator) {
						v.Field("ip").Required().String()
						v.Field("userAgent").Required().String()
					})
				})
				return v
			},
			want: true,
		},
		{
			name: "should fail - depth 2, with two maps missing required fields",
			value: map[string]interface{}{
				"data": map[string]interface{}{
					"status": true,
				},
				"requestData": map[string]interface{}{
					"ip": "127.0.0.1",
				},
			},
			setup: func() *MapValidator {
				v := NewMapValidator()
				v.Field(fieldName).Map(func(v *MapValidator) {
					v.Field("data").Map(func(v *MapValidator) {
						v.Field("status").Required().Bool()
					})
					v.Field("requestData").Map(func(v *MapValidator) {
						v.Field("ip").Required().String()
						v.Field("userAgent").Required().String()
					})
				})
				return v
			},
			want: false,
		},
		{
			name: "should pass - depth 4 with slices that contain maps",
			value: map[string]interface{}{
				"one": map[string]interface{}{
					"two": map[string]interface{}{
						"three": map[string]interface{}{
							"four": map[string]interface{}{
								"children": []map[string]interface{}{
									{
										"id":   "e4eaaaf2-d142-11e1-b3e4-080027620cdd",
										"name": "test",
									},
								},
							},
						},
					},
				},
			},
			setup: func() *MapValidator {
				v := NewMapValidator()
				v.Field(fieldName).Map(func(v *MapValidator) {
					v.Field("one").Map(func(v *MapValidator) {
						v.Field("two").Map(func(v *MapValidator) {
							v.Field("three").Map(func(v *MapValidator) {
								v.Field("four").Map(func(v *MapValidator) {
									v.Field("children").Required().Slice(func(v *FieldValidator) {
										v.Map(func(v *MapValidator) {
											v.Field("id").Required().UUID()
											v.Field("name").Required().String()
										})
									})
								})
							})
						})
					})
				})
				return v
			},
			want: true,
		},
	}

	validateTestDataWithSetup(t, testTable)

}

func Test_NestedSlices(t *testing.T) {
	testTable := []validatorTestTable{
		{
			name:  "should pass - depth 1",
			value: []string{"e4eaaaf2-d142-11e1-b3e4-080027620cdd"},
			setup: func() *MapValidator {
				v := NewMapValidator()
				v.Field(fieldName).Slice(func(v *FieldValidator) {
					v.UUID()
				})
				return v
			},
			want: true,
		},
		{
			name:  "should fail - depth 1 invalid child",
			value: []uint{123},
			setup: func() *MapValidator {
				v := NewMapValidator()
				v.Field(fieldName).Slice(func(v *FieldValidator) {
					v.UUID()
				})
				return v
			},
			want: false,
		},
		{
			name:  "should pass - inner map has no required fields",
			value: []map[string]interface{}{},
			setup: func() *MapValidator {
				v := NewMapValidator()
				v.Field(fieldName).Slice(func(v *FieldValidator) {
					v.Map(func(v *MapValidator) {
						v.Field("year").Num()
						v.Field("month").Num()
					})
				})
				return v
			},
			want: true,
		},
		{
			name: "should pass - contains inner map",
			value: []map[string]interface{}{
				{
					"year":  2024,
					"month": 10,
				},
				{
					"year":  2022,
					"month": 6,
				},
			},
			setup: func() *MapValidator {
				v := NewMapValidator()
				v.Field(fieldName).Slice(func(v *FieldValidator) {
					v.Map(func(v *MapValidator) {
						v.Field("year").Required().Num()
						v.Field("month").Required().Num()
					})
				})
				return v
			},
			want: true,
		},
		{
			name: "should pass - inner map contains slice",
			value: []map[string]interface{}{
				{
					"data": map[string]interface{}{
						"messages": []string{"hi", "hello"},
					},
				},
			},
			setup: func() *MapValidator {
				v := NewMapValidator()
				v.Field(fieldName).Slice(func(v *FieldValidator) {
					v.Map(func(v *MapValidator) {
						v.Field("data").Map(func(v *MapValidator) {
							v.Field("messages").Slice(func(v *FieldValidator) {
								v.Alpha()
							})
						})
					})
				})
				return v
			},
			want: true,
		},
		{
			name:  "should pass - inner slices of depth 4",
			value: [][][][]int{{{{1}}}},
			setup: func() *MapValidator {
				v := NewMapValidator()
				v.Field(fieldName).Slice(func(v *FieldValidator) {
					v.Slice(func(v *FieldValidator) {
						v.Slice(func(v *FieldValidator) {
							v.Slice(func(v *FieldValidator) {
								v.Num()
							})
						})
					})
				})
				return v
			},
			want: true,
		},
		{
			name:  "should fail - inner slices of depth 4, invalid child",
			value: [][][][]string{{{{"1"}}}},
			setup: func() *MapValidator {
				v := NewMapValidator()
				v.Field(fieldName).Slice(func(v *FieldValidator) {
					v.Slice(func(v *FieldValidator) {
						v.Slice(func(v *FieldValidator) {
							v.Slice(func(v *FieldValidator) {
								v.Num()
							})
						})
					})
				})
				return v
			},
			want: false,
		},
	}

	validateTestDataWithSetup(t, testTable)
}

func Test_Chaining(t *testing.T) {
	testTable := []messageValidatorTestTable{
		{
			name:  "simple chain should be required",
			value: "",
			setup: func() *MapValidator {
				v := NewMapValidator()
				v.Field(fieldName).Required().String()

				return v
			},
			expectedMessage: "is required",
		},
		{
			name:  "irrational chain should fail on second rule (uuid)",
			value: 42,
			setup: func() *MapValidator {
				v := NewMapValidator()
				v.Field(fieldName).Required().Num().UUID().String()

				return v
			},
			expectedMessage: "valid UUID",
		},
		{
			name:  "nested map field is required",
			value: map[string]interface{}{},
			setup: func() *MapValidator {
				v := NewMapValidator()
				v.Field(fieldName).Map(func(v *MapValidator) {
					v.Field("status").Required().Bool()
				})

				return v
			},
			expectedMessage: "status is required",
		}, {
			name:  "long chain should pass",
			value: "test!",
			setup: func() *MapValidator {
				v := NewMapValidator()
				v.Field(fieldName).Required().String().MinS(3).MaxS(10).MinRune(3).MaxRune(10)

				return v
			},
			expectedMessage: "",
		},
	}

	validateTestDataMessage(t, testTable)
}
