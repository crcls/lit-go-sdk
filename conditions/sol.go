package conditions

type SolRpcCondition struct {
	Method          string          `json:"method"`
	Params          []string        `json:"params"`
	Chain           string          `json:"chain"`
	ReturnValueTest ReturnValueTest `json:"returnValueTest"`
}