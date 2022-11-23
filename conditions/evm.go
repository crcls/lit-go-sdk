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
