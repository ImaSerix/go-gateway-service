package types

type CheckName string

var (
	Policy         CheckName = "policy"
	HeaderRequired CheckName = "header_required"
	IPWhiteList    CheckName = "ip_whitelist"
	QueryRequired  CheckName = "query_required"
)
