package http

const (
	httpModuleName      = "http"
	httpClientUDName    = "http.udclient"
	httpClientApiName   = "http.clientApis"
	httpRequestUDName   = "http.udreq"
	httpRequestApiName  = "http.reqApis"
	httpResponseApiName = "http.resApis"
	httpResponseUDName  = "http.udres"

	//http client send method
	sendMethod = "send"

	//http methods
	newHttpClientMethod  = "newHttpClient"
	newHttpRequestMethod = "newHttpRequest"
	urlEncodeMethod      = "urlEncode"
	urlDecodeMethod      = "urlDecode"

	//for http request methods
	authMethod       = "auth"
	addHeaderMethod  = "addHeader"
	addHeadersMethod = "addHeaders"
	putStringMethod  = "putString"
	putHexMethod     = "putHex"
	uploadMethod     = "upload"

	//for http response methods
	getStatusCodeMethod   = "getStatusCode"
	getBodyAsByteMethod   = "getBodyAsByte"
	getBodyAsStringMethod = "getBodyAsString"
	getProtocolMethod     = "getProtocol"
	getHeaderMethod       = "getHeader"
	getHeadersMethod      = "getHeaders"
)
