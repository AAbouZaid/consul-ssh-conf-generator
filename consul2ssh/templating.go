package consul2ssh

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"reflect"
	"sort"
	"text/template"
)

type SSHNodeConf struct {
	Host string `json:"node"`
	Main MapInterface
}

// fmtSSHElems - format SSH config elements.
func fmtSSHElems(m MapInterface) []string {
	output := []string{}
	for key, value := range m {
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice, reflect.Array:
			for _, item := range value.([]interface{}) {
				output = append(output, fmt.Sprintf("%v %v", key, item))
			}
		default:
			output = append(output, fmt.Sprintf("%v %v", key, value))
		}
	}
	sort.Strings(output)
	return output
}

// Functions that will be used in the template.
var templFuncs = template.FuncMap{
	"fmtSSHElems": fmtSSHElems,
}

// buildSSHTemplate - Make SSH config template.
func (c *SSHNodeConf) buildTemplate(
	w http.ResponseWriter,
	templateFile string,
) error {

	// Generate the template.
	tmplFileBase := path.Base(templateFile)
	tmplFile := getFilePath(templateFile)
	tmpl, err := template.New(tmplFileBase).Funcs(templFuncs).ParseFiles(tmplFile)
	if err != nil {
		log.Print(err)
		return err
	}
	if err := tmpl.Execute(w, c); err != nil {
		log.Print(err)
		return err
	}

	return nil
}
