/*
  http_tengo.go,client_tengo.go,request_tengo.go和response_tengo.go 定义了http的tengo接口，
  能够在tengo脚本里构造http请求，获得http响应内容，便于用脚本操作http,现将相关的接口函数简要
  说明如下,tengo的模块名称http：
   1. 构建一个http client
       client := http.newHttpClient(host,port,isSSL,timeout)

   2. 构建一个http request
        req := http.newHttpRequest(method,uri)

		2.1 授权函数
   		req.auth(user,passwd)

		2.2 添加请求头部
   		req.addHeader(key,value)

		2.3 通过字符串格式设置post/put等方法请求体
		req.putString(content,isFromFile)

		2.4 通过十六进制格式设置post/put等方法请求体
		req.putHex(content)
		2.5 文件上传
		req.upload(fname,filepath,boundary)

	3. 发送http request 方法
		res := client.send(req)

		3.1. 获取响应状态码
   		res.getStatusCode()

		3.2 获取二进制形式响应内容
   		res.getBodyAsByte()

		3.3 获取字符串形式响应内容
   		res.getBodyAsString()

		3.4 获取http协议版本
   		res.getProtocol()

		3.5 获取指定key的响应头部的value值
   		res.getHeader(key)

		3.6	获取指定key的响应头部所有value值
   		res.getHeaders(key)

	4. url 编码
		http.urlEncode(uri)

	5. url 解码
		http.urlDecode(uri)

*/

package http

import (
	stengo "common/scripts/tengo"
	"errors"
	"net/url"

	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
)

var (
	errNewClient = errors.New("New Http Client Invalid args,must provide <host><port><isSSL><timeout>")
	errNewRequst = errors.New("New Http Request Invalid args,must provide <method><uri>")
)

type HttpTengo struct {
	stengo.TengoObj
}

func newHttpClient(args ...objects.Object) (objects.Object, error) {

	if len(args) != 4 {

		return nil, errNewClient
	}

	host, ok := objects.ToString(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "host",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	port, ok := objects.ToInt(args[1])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "port",
			Expected: "int(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	isSSL, ok := objects.ToBool(args[2])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "isSSL",
			Expected: "bool(compatible)",
			Found:    args[2].TypeName(),
		}
	}

	timeOut, ok := objects.ToInt64(args[3])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "port",
			Expected: "int64(compatible)",
			Found:    args[3].TypeName(),
		}
	}

	return &HttpClientTengo{
		TengoObj: stengo.TengoObj{Name: "HttpClient"},
		client:   NewHttpClient(host, port, isSSL, timeOut),
	}, nil
}

//tengo call method example: http.newHttpRequest("get","/")
func newHttpRequest(args ...objects.Object) (objects.Object, error) {

	if len(args) != 2 {

		return nil, errNewRequst
	}

	method, ok := objects.ToString(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "method",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	uri, ok := objects.ToString(args[1])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "uri",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	return &HttpRequestTengo{
		TengoObj: stengo.TengoObj{Name: "HttpRequest"},
		req:      NewHttpRequest(method, uri),
	}, nil
}

func urlEncode(args ...objects.Object) (objects.Object, error) {

	if len(args) != 1 {

		return nil, tengo.ErrWrongNumArguments
	}

	urlRaw, ok := objects.ToString(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "url",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	return objects.FromInterface(url.QueryEscape(urlRaw))
}

func urlDecode(args ...objects.Object) (objects.Object, error) {

	if len(args) != 1 {

		return nil, tengo.ErrWrongNumArguments
	}

	urlRaw, ok := objects.ToString(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "url",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	durl, err := url.QueryUnescape(urlRaw)

	if err != nil {
		return objects.FromInterface(urlRaw)
	}

	return objects.FromInterface(durl)
}

var moduleMap objects.Object = &objects.ImmutableMap{

	Value: map[string]objects.Object{
		newHttpClientMethod: &objects.UserFunction{
			Name:  newHttpClientMethod,
			Value: newHttpClient,
		},
		newHttpRequestMethod: &objects.UserFunction{
			Name:  newHttpRequestMethod,
			Value: newHttpRequest,
		},

		urlEncodeMethod: &objects.UserFunction{
			Name:  urlEncodeMethod,
			Value: urlEncode,
		},

		urlDecodeMethod: &objects.UserFunction{
			Name:  urlDecodeMethod,
			Value: urlDecode,
		},
	},
}

func (HttpTengo) Import(moduleName string) (interface{}, error) {

	switch moduleName {
	case httpModuleName:
		return moduleMap, nil
	default:
		return nil, errors.New("undefined module" + moduleName)
	}
}
