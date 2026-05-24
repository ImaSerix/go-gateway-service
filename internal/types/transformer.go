package types

type TransformerName string

var (
	Headers     TransformerName = "header"
	BodyFields  TransformerName = "body_fields"
	QueryParams TransformerName = "query_params"
)
