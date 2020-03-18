package base

const (
	//StatusOK status ok
	StatusOK = "ok"
	//StatusError status error
	StatusError = "error"
	//StatusProxy status error
	StatusProxy = "proxy"
	//StatusNoMatch status no match
	StatusNoMatch = "noMatch"
	//StatusNoFound status no found
	StatusNoFound = "notFound"

	//StatusPartial partial match (waiting for done marker)
	StatusPartial = "partial"

	//StatusDisabled status disabled
	StatusDisabled = "disabled"

	//StatusUnProcess status for unprocessed file
	StatusUnProcess = "unprocessed"

	//SourceAttribute dest attribute
	SourceAttribute = "Source"

	//UnclassifiedStatus
	UnclassifiedStatus = "unclassified"

	//ConfigEnvKey config env key
	ConfigEnvKey = "CONFIG"

	//DestEnvKey destination env key
	DestEnvKey = "DEST"


	//LambdaScheme represents lambda schem
	LambdaScheme = "lambda"

	//YAMLExt yaml extension
	YAMLExt = ".yaml"

	//CloudFunctionScheme represents clud function scheme
	CloudFunctionScheme = "cloudfunction"

	//InMemoryStorageBaseURL in memory storage URL
	InMemoryStorageBaseURL = "mem://localhost/"

)
