package conditions

type AuthCondition interface {
	EvmContractCondition | AccessControlCondition
}

type ReturnValueTest struct {
	Key        string      `json:"key"`
	Comparator string      `json:"comparator"`
	Value      interface{} `json:"value"`
}
