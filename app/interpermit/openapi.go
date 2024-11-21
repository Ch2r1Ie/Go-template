package interpermit

// swagger:route POST /interpermits interpermit permitTransaction
// responses:
//

// swagger:parameters permitTransaction
type swaggerAddUserParamsWrapper struct {
	// in:header
	// required:true
	// example: 7c129eb1-c479-47bb-9c73-d263e2673011
	XRefID string `json:"X-Ref-Id"`

	// in:body
	Body SubmitTransactionRequest
}
