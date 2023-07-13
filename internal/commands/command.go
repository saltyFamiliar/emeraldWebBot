package commands

import (
	"fmt"
	"reflect"
	"strconv"
)

type CmdFunc struct {
	Raw     interface{}
	Wrapper func(...interface{}) (string, error)
}

var CmdMap = map[string]CmdFunc{
	"scan": {
		Raw: ScanPorts,
		Wrapper: func(args ...interface{}) (string, error) {
			return ScanPorts(args[0].(string), args[1].(int), args[2].(int), args[3].(int))
		},
	},
}

type Command struct {
	CastParams []interface{}
	Params     []string
	Fn         CmdFunc
	Result     string
	ErrorMsg   string
}

func (cmd *Command) Execute() (string, error) {
	return cmd.Fn.Wrapper(cmd.CastParams...)
}

func (cmd *Command) ValidateParams() bool {
	fnType := reflect.TypeOf(cmd.Fn.Raw)
	numParams := fnType.NumIn()
	if numParams != len(cmd.Params) {
		fmt.Printf("Wrong num. Given: %d, Takes: %d", len(cmd.Params), numParams)
		return false
	}
	castParams := make([]interface{}, numParams)
	for i := 0; i < numParams; i++ {
		paramType := fnType.In(i)
		switch paramType.Kind() {
		case reflect.Int:
			asInt, err := strconv.Atoi(cmd.Params[i])
			if err != nil {
				fmt.Println("Wrong type")
				return false
			}
			castParams[i] = asInt
		case reflect.String:
			castParams[i] = cmd.Params[i]
		case reflect.Float64:
			asFloat, err := strconv.ParseFloat(cmd.Params[i], 64)
			if err != nil {
				fmt.Println("Wrong type")
				return false
			}
			castParams[i] = asFloat
		case reflect.Float32:
			asFloat, err := strconv.ParseFloat(cmd.Params[i], 32)
			if err != nil {
				fmt.Println("Wrong type")
				return false
			}
			castParams[i] = asFloat
		default:
			fmt.Println("default hit")
		}

	}
	cmd.CastParams = castParams
	fmt.Println("Is valid")
	return true
}

func NewCommand(funcName string, params []string) (*Command, error) {
	cmdFn, ok := CmdMap[funcName]
	if !ok {
		return nil, fmt.Errorf("invalid command name: %s", funcName)
	}

	cmd := Command{
		Params: params,
		Fn:     cmdFn,
	}

	if cmd.ValidateParams() {
		return &cmd, nil
	}
	return nil, fmt.Errorf("invalid arg list: %s", params)
}
