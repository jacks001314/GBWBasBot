package tplutils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"
)

func MakeSourceScriptFromTemplate(tplFile string, data interface{}) ([]byte, error) {

	contentBuffer := bytes.NewBuffer(make([]byte, 0))

	content, err := ioutil.ReadFile(tplFile)
	if err != nil {

		return nil, err
	}

	tp := template.New("template")

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

	tp, err = tp.Parse(string(content))

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
