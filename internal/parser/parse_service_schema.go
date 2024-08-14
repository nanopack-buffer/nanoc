package parser

import (
	"errors"
	"fmt"
	"nanoc/internal/errs"
	"nanoc/internal/npschema"
	"nanoc/internal/symbol"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

var funcHeaderRegex = regexp.MustCompile("(?P<name>^\\w+)(?P<params>\\((?:\\w+ \\w+, |\\w+ \\w+)*\\))$")

func parseServiceSchema(header string, body yaml.MapSlice) (*npschema.PartialService, error) {
	ps := strings.Split(header, " ")
	if len(ps) != 2 {
		return nil, errors.New(fmt.Sprintf("invalid service header declaration received: %v. A valid example: service MyService.", header))
	}

	schema := npschema.PartialService{
		Name: ps[1],
	}

	for _, e := range body {
		k, ok := e.Key.(string)
		if !ok {
			return nil, errs.NewNanocError("Invalid function declaration", fmt.Sprintf("service %v", schema.Name))
		}

		var returnTypeName string
		if e.Value == nil {
			returnTypeName = ""
		} else {
			v, ok := e.Value.(string)
			if !ok {
				fmt.Println("asdkjkasjdk;sdj")
				return nil, errs.NewNanocError("Invalid function declaration", fmt.Sprintf("service %v", schema.Name))
			}
			returnTypeName = v
		}

		f, err := parseFunction(k, returnTypeName)
		if err != nil {
			var syntaxErr *errs.NanocError
			if errors.As(err, &syntaxErr) {
				syntaxErr.PrependToStack(fmt.Sprintf("service %v", schema.Name))
				return nil, syntaxErr
			}
			return nil, err
		}
		schema.DeclaredFunctions = append(schema.DeclaredFunctions, *f)
	}

	return &schema, nil
}

func parseFunction(funcHeader, returnType string) (*npschema.PartialDeclaredFunction, error) {
	stripped, found := strings.CutPrefix(funcHeader, symbol.Func+" ")
	if !found {
		return nil, errs.NewNanocError("Function must start with a 'func' keyword", funcHeader)
	}

	matches := funcHeaderRegex.FindStringSubmatch(stripped)
	if len(matches) != 3 {
		return nil, errs.NewNanocError("Invalid function declaration", funcHeader)
	}

	funcName := ""
	funcParamsStr := ""
	for i, name := range funcHeaderRegex.SubexpNames() {
		switch name {
		case "name":
			funcName = matches[i]
		case "params":
			funcParamsStr = strings.Trim(matches[i], "()")
		}
	}

	if funcName == "" {
		return nil, errs.NewNanocError("Function name cannot be empty", funcHeader)
	}
	if funcParamsStr == "" {
		return nil, errs.NewNanocError("Parenthesis must follow function name", funcName)
	}

	f := npschema.PartialDeclaredFunction{
		Name: funcName,
	}

	for _, s := range strings.Split(funcParamsStr, ",") {
		s = strings.TrimSpace(s)
		ps := strings.Split(s, " ")
		if len(ps) != 2 {
			return nil, errs.NewNanocError("Invalid parameter syntax", s, funcParamsStr, funcHeader)
		}

		paramName := ps[0]
		paramType := ps[1]

		f.Parameters = append(f.Parameters, npschema.PartialFunctionParameter{
			Name:     paramName,
			TypeName: paramType,
		})
	}

	ps := strings.Split(returnType, " "+symbol.FuncThrows+" ")
	switch len(ps) {
	case 1:
		f.ReturnTypeName = ps[0]
	case 2:
		f.ReturnTypeName = ps[0]
		f.ErrorTypeName = ps[1]
	default:
		return nil, errs.NewNanocError("Invalid return type syntax. Valid example: int, int throws MyError.", returnType, funcHeader)
	}

	return &f, nil
}
