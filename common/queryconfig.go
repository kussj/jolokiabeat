package common

type QueryConfig struct {
	Domain 		string
	Attributes []string
}

func(qc QueryConfig) GetDomain() string {
	return qc.Domain
}

func(qc QueryConfig) GetAttributes() []string {
	return qc.Attributes
}