package base64

import (
	"encoding/base64"

	lua "github.com/yuin/gopher-lua"
)

const (
	Base64ModuleName   = "base64"
	base64EncodingType = "base64.Encoding"
	base64EncoderType  = "base64.Encoder"
	base64DecoderType  = "base64.Decoder"
)

// Loader is the module loader function.
func Loader(L *lua.LState) int {

	registerBase64Encoding(L)
	registerBase64Encoder(L)
	registerBase64Decoder(L)

	// Register the encodings offered by base64 go module.
	t := L.NewTable()
	L.SetField(t, "RawStdEncoding", LVBase64Encoding(L, base64.RawStdEncoding))
	L.SetField(t, "RawURLEncoding", LVBase64Encoding(L, base64.RawURLEncoding))
	L.SetField(t, "StdEncoding", LVBase64Encoding(L, base64.StdEncoding))
	L.SetField(t, "URLEncoding", LVBase64Encoding(L, base64.URLEncoding))
	L.SetFuncs(t, map[string]lua.LGFunction{
		"new_encoding": NewEncoding,
		"new_encoder":  NewEncoder,
		"new_decoder":  NewDecoder,
	})
	L.Push(t)
	return 1
}
