package cmds

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"
	"text/template"

	"github.com/pobri19/sqlboiler/dbdrivers"
)

// generateTemplate generates the template associated to the passed in command name.
func generateTemplate(commandName string, data *tplData) []byte {
	template := getTemplate(commandName)

	if template == nil {
		errorQuit(fmt.Errorf("Unable to find the template: %s", commandName+".tpl"))
	}

	output, err := processTemplate(template, data)
	if err != nil {
		errorQuit(fmt.Errorf("Unable to process the template: %s", err))
	}

	return output
}

// getTemplate returns a pointer to the template matching the passed in name
func getTemplate(name string) *template.Template {
	var tpl *template.Template

	// Find the template that matches the passed in template name
	for _, t := range templates {
		if t.Name() == name+".tpl" {
			tpl = t
			break
		}
	}

	return tpl
}

// processTemplate takes a template and returns the output of the template execution.
func processTemplate(t *template.Template, data *tplData) ([]byte, error) {
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, err
	}

	output, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, err
	}

	return output, nil
}

// it into a go styled object variable name of "ColumnName".
// makeGoName also fully uppercases "ID" components of names, for example
// "column_name_id" to "ColumnNameID".
func makeGoName(name string) string {
	s := strings.Split(name, "_")

	for i := 0; i < len(s); i++ {
		if s[i] == "id" {
			s[i] = "ID"
			continue
		}
		s[i] = strings.Title(s[i])
	}

	return strings.Join(s, "")
}

// makeGoVarName takes a variable name in the format of "var_name" and converts
// it into a go styled variable name of "varName".
// makeGoVarName also fully uppercases "ID" components of names, for example
// "var_name_id" to "varNameID".
func makeGoVarName(name string) string {
	s := strings.Split(name, "_")

	for i := 0; i < len(s); i++ {

		if s[i] == "id" && i > 0 {
			s[i] = "ID"
			continue
		}

		if i == 0 {
			continue
		}

		s[i] = strings.Title(s[i])
	}

	return strings.Join(s, "")
}

// makeDBName takes a table name in the format of "table_name" and a
// column name in the format of "column_name" and returns a name used in the
// `db:""` component of an object in the format of "table_name_column_name"
func makeDBName(tableName, colName string) string {
	return tableName + "_" + colName
}

// makeGoInsertParamNames takes a []DBColumn and returns a comma seperated
// list of parameter names for the insert statement template.
func makeGoInsertParamNames(data []dbdrivers.DBColumn) string {
	var paramNames string
	for i := 0; i < len(data); i++ {
		paramNames = paramNames + data[i].Name
		if len(data) != i+1 {
			paramNames = paramNames + ", "
		}
	}
	return paramNames
}

// makeGoInsertParamFlags takes a []DBColumn and returns a comma seperated
// list of parameter flags for the insert statement template.
func makeGoInsertParamFlags(data []dbdrivers.DBColumn) string {
	var paramFlags string
	for i := 0; i < len(data); i++ {
		paramFlags = fmt.Sprintf("%s$%d", paramFlags, i+1)
		if len(data) != i+1 {
			paramFlags = paramFlags + ", "
		}
	}
	return paramFlags
}

// makeSelectParamNames takes a []DBColumn and returns a comma seperated
// list of parameter names with for the select statement template.
// It also uses the table name to generate the "AS" part of the statement, for
// example: var_name AS table_name_var_name, ...
func makeSelectParamNames(tableName string, data []dbdrivers.DBColumn) string {
	var paramNames string
	for i := 0; i < len(data); i++ {
		paramNames = fmt.Sprintf("%s%s AS %s", paramNames, data[i].Name,
			makeDBName(tableName, data[i].Name),
		)
		if len(data) != i+1 {
			paramNames = paramNames + ", "
		}
	}
	return paramNames
}