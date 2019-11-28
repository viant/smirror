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

	//ProxyMethod proxy method env key
	ProxyMethod = "METHOD"

	//ProxyMethodCall call operation
	ProxyMethodCall = "call"
	//ProxyMethodCopy copy operation
	ProxyMethodCopy = "copy"
	//ProxyMethodMove move operation
	ProxyMethodMove = "move"

	//App
	App = "StorageMirror"

	//LambdaScheme represents lambda schem
	LambdaScheme = "lambda"

	//CloudFunctionScheme represents clud function scheme
	CloudFunctionScheme = "cloudfunction"
)
