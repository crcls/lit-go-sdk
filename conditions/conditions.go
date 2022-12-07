package conditions

type ReturnValueTest struct {
	Key        string      `json:"key"`
	Comparator string      `json:"comparator"`
	Value      interface{} `json:"value"`
}
