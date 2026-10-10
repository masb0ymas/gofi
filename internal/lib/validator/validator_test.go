package validator

import "testing"

func TestRules(t *testing.T) {
	cases := []struct {
		name  string
		build func(fv *FieldValidator)
		data  interface{}
		pass  bool
	}{
		{"required absent fails", func(fv *FieldValidator) { fv.Required() }, nil, false},
		{"required empty string fails", func(fv *FieldValidator) { fv.Required() }, "", false},
		{"required value passes", func(fv *FieldValidator) { fv.Required() }, "x", true},
		{"required zero passes", func(fv *FieldValidator) { fv.Required() }, float64(0), true},
		{"string absent passes", func(fv *FieldValidator) { fv.String() }, nil, true},
		{"string rejects number", func(fv *FieldValidator) { fv.String() }, float64(1), false},
		{"minlen short fails", func(fv *FieldValidator) { fv.MinLen(8) }, "1234567", false},
		{"minlen exact passes", func(fv *FieldValidator) { fv.MinLen(8) }, "12345678", true},
		{"minlen absent passes", func(fv *FieldValidator) { fv.MinLen(8) }, nil, true},
		{"email ok", func(fv *FieldValidator) { fv.Email() }, "a@b.co", true},
		{"email bad", func(fv *FieldValidator) { fv.Email() }, "nope", false},
		{"email absent passes", func(fv *FieldValidator) { fv.Email() }, nil, true},
		{"withins ok", func(fv *FieldValidator) { fv.WithinS("a", "b") }, "a", true},
		{"withins bad", func(fv *FieldValidator) { fv.WithinS("a", "b") }, "z", false},
		{"withins non-string passes", func(fv *FieldValidator) { fv.WithinS("a", "b") }, float64(1), true},
		{"regex ok", func(fv *FieldValidator) { fv.Regex(`^\d+$`) }, "42", true},
		{"regex bad", func(fv *FieldValidator) { fv.Regex(`^\d+$`) }, "x", false},
		{"number absent passes", func(fv *FieldValidator) { fv.Number() }, nil, true},
		{"number accepts float", func(fv *FieldValidator) { fv.Number() }, float64(1.5), true},
		{"number rejects string", func(fv *FieldValidator) { fv.Number() }, "1", false},
		{"bool absent passes", func(fv *FieldValidator) { fv.Bool() }, nil, true},
		{"bool accepts false", func(fv *FieldValidator) { fv.Bool() }, false, true},
		{"bool rejects string", func(fv *FieldValidator) { fv.Bool() }, "true", false},
		{"uuid absent passes", func(fv *FieldValidator) { fv.UUID() }, nil, true},
		{"uuid ok", func(fv *FieldValidator) { fv.UUID() }, "019f3766-6e40-7573-84a8-90f1f762a972", true},
		{"uuid bad", func(fv *FieldValidator) { fv.UUID() }, "nope", false},
		{"date absent passes", func(fv *FieldValidator) { fv.Date() }, nil, true},
		{"date ok", func(fv *FieldValidator) { fv.Date() }, "2026-10-01T00:00:00Z", true},
		{"date rejects bare date", func(fv *FieldValidator) { fv.Date() }, "2026-10-01", false},
		{"array absent fails", func(fv *FieldValidator) { fv.Array() }, nil, false},
		{"array empty fails", func(fv *FieldValidator) { fv.Array() }, []interface{}{}, false},
		{"array non-empty passes", func(fv *FieldValidator) { fv.Array() }, []interface{}{"x"}, true},
	}

	for _, tc := range cases {
		v := NewMapValidator()
		tc.build(v.Field("f"))
		mr, pass := v.Validate(map[string]interface{}{"f": tc.data})
		if pass != tc.pass {
			t.Errorf("%s: pass=%v, want %v (mr=%v)", tc.name, pass, tc.pass, mr)
		}
		if pass && !mr.Empty() {
			t.Errorf("%s: passing validation returned messages %v", tc.name, mr)
		}
		if !pass && len(mr["f"]) == 0 {
			t.Errorf("%s: failing validation returned no message for f", tc.name)
		}
	}
}

func TestValidate_CollectsPerFieldMessages(t *testing.T) {
	v := NewMapValidator()
	v.Field("name").Required()
	v.Field("role").Required().WithinS("admin", "member")

	mr, pass := v.Validate(map[string]interface{}{"role": "bogus"})
	if pass {
		t.Fatal("expected validation to fail")
	}
	if len(mr["name"]) != 1 || mr["name"][0] != "name is required" {
		t.Fatalf("expected name required message, got %v", mr["name"])
	}
	if len(mr["role"]) != 1 {
		t.Fatalf("expected one role message, got %v", mr["role"])
	}
}

func TestFieldValidator_StopsAtFirstFailure(t *testing.T) {
	v := NewMapValidator()
	v.Field("role").Required().WithinS("admin")

	mr, pass := v.Validate(map[string]interface{}{})
	if pass || len(mr["role"]) != 1 || mr["role"][0] != "role is required" {
		t.Fatalf("expected only the required message, got %v", mr)
	}
}
