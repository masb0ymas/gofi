package validator

type MapValidator struct {
	fvs map[string]*FieldValidator
}

func NewMapValidator() *MapValidator {
	return &MapValidator{fvs: make(map[string]*FieldValidator)}
}

// Validate runs every field's rules; the first failing rule per field wins.
func (v *MapValidator) Validate(dict map[string]interface{}) (MessageRecord, bool) {
	mr := make(MessageRecord)
	for key, fv := range v.fvs {
		mr = mr.Append(fv.Validate(dict[key]))
	}
	return mr, mr.Empty()
}

func (v *MapValidator) Field(key string) *FieldValidator {
	fv := &FieldValidator{key: key}
	v.fvs[key] = fv
	return fv
}

type FieldValidator struct {
	key   string
	rules []rule
}

func (v *FieldValidator) Validate(data interface{}) MessageRecord {
	for _, rule := range v.rules {
		if msg, ok := rule(v.key, data); !ok {
			return MessageRecord{v.key: []string{msg}}
		}
	}
	return nil
}

func (v *FieldValidator) registerRule(rule rule) {
	v.rules = append(v.rules, rule)
}
