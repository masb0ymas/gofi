package validator

type MessageRecord map[string][]string

func (mr MessageRecord) Append(amr MessageRecord) MessageRecord {
	merged := make(MessageRecord)

	for key, value := range mr {
		merged[key] = append(merged[key], value...)
	}

	for key, value := range amr {
		merged[key] = append(merged[key], value...)
	}

	return merged
}

func (mr MessageRecord) Empty() bool {
	return len(mr) == 0
}
