package validator

type MapValidator struct {
	path path
	fvs  map[string]*FieldValidator
}

func NewMapValidator() *MapValidator {
	v := &MapValidator{}
	v.fvs = make(map[string]*FieldValidator)
	return v
}

func NewMapValidatorWithPath(path path) *MapValidator {
	v := &MapValidator{}
	v.path = path
	v.fvs = make(map[string]*FieldValidator)
	return v
}

func (v *MapValidator) Validate(dict map[string]interface{}) (MessageRecord, bool) {
	sumMr := make(MessageRecord)

	for key, fv := range v.fvs {
		data := dict[key]
		mr, passes := fv.Validate(data)
		if !passes {
			sumMr = sumMr.Append(mr)
		}
	}

	passes := sumMr.Empty()
	return sumMr, passes
}

func (v *MapValidator) Field(key string) *FieldValidator {
	fv := &FieldValidator{path: append(v.path, key)}
	v.fvs[key] = fv
	return fv
}

type FieldValidator struct {
	path  path
	rules []rule
}

func (v *FieldValidator) Validate(data interface{}) (MessageRecord, bool) {
	var currData interface{} = data

	for _, rule := range v.rules {
		data, mr, passes := rule(v.path, currData)

		// Stop as soon as the first rule fails. There is no need
		// to check the remaining rules.
		if !passes {
			return mr, false
		}

		currData = data
	}

	return make(MessageRecord), true
}

func (v *FieldValidator) registerRule(rule rule) {
	v.rules = append(v.rules, rule)
}
