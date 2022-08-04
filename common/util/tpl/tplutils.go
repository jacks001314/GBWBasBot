package tplutils

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

func MakeSourceScriptFromTemplate(tplFile, fname string, data interface{}) ([]byte, error) {

	contentBuffer := bytes.NewBuffer(make([]byte, 0))

	tp := template.New(fname)

	tp = tp.Funcs(template.FuncMap{
		"toTStrArray": func(arr []string) string {

			arrStr := make([]string, 0)

			for _, a := range arr {

				arrStr = append(arrStr, fmt.Sprintf(`"%s"`, a))
			}

			return fmt.Sprintf("[%s]", strings.Join(arrStr, ","))
		},
		"toLStrArray": func(arr []string) string {

			arrStr := make([]string, 0)

			for _, a := range arr {

				arrStr = append(arrStr, fmt.Sprintf(`"%s"`, a))
			}

			return fmt.Sprintf("{%s}", strings.Join(arrStr, ","))
		},
	})

	tp, err := tp.ParseFiles(tplFile)

	if err != nil {
		errS := fmt.Sprintf("Parse  Template file:%s ,err:[%v]", tplFile, err)
		return nil, fmt.Errorf(errS)
	}

	if err = tp.Execute(contentBuffer, data); err != nil {
		errS := fmt.Sprintf("Parse Attack Target Template file:%s ,err:[%v]", tplFile, err)
		return nil, fmt.Errorf(errS)
	}

	return contentBuffer.Bytes(), nil
}
