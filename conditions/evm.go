package conditions

type AbiIO struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type AbiMember struct {
	Name            string  `json:"name"`
	Inputs          []AbiIO `json:"inputs"`
	Outputs         []AbiIO `json:"outputs"`
	Constant        bool    `json:"constant"`
	StateMutability string  `json:"stateMutability"`
}

type EvmContractCondition struct {
	ContractAddress string          `json:"contractAddress"`
	FunctionName    string          `json:"functionName"`
	FunctionParams  []string        `json:"functionParams"`
	FunctionAbi     AbiMember       `json:"functionAbi"`
	Chain           string          `json:"chain"`
	ReturnValueTest ReturnValueTest `json:"returnValueTest"`
}

type AccessControlCondition struct {
	ContractAddress      string          `json:"contractAddress"`
	Chain                string          `json:"chain"`
	StandardContractType string          `json:"standardContractType"`
	Method               string          `json:"method"`
	Parameters           []string        `json:"parameters"`
	ReturnValueTest      ReturnValueTest `json:"returnValueTest"`
}
