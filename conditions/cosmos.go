package conditions

type CosmosCondition struct {
	Path            string `json:"path"`
	Chain           string `json:"chain"`
	ReturnValueTest ReturnValueTest
}
