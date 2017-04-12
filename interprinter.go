package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func eval(env map[string]interface{}, v interface{}) interface{} {
	if vl, ok := v.([]interface{}); ok {
		return doRun(env, vl)
	}
	return v
}

func doRun(env map[string]interface{}, v []interface{}) interface{} {
	var r interface{}

	mn := v[0].(string)
	switch mn {
	case "step":
		for _, vi := range v[1:] {
			r = doRun(env, vi.([]interface{}))
		}
	case "until":
		for {
			c := eval(env, v[1])
			if c.(bool) == true {
				break
			}
			r = doRun(env, v[2].([]interface{}))
		}
	case "get":
		return env[eval(env, v[1]).(string)]
	case "set":
		env[eval(env, v[1]).(string)] = eval(env, v[2])
		return v[2]
	case "=":
		return eval(env, v[1]).(float64) == eval(env, v[2]).(float64)
	case "+":
		return eval(env, v[1]).(float64) + eval(env, v[2]).(float64)
	case "-":
		panic("Todo")
	case "/":
		panic("Todo")
	case "*":
		panic("Todo")
	default:
		panic("Unknown operation: " + fmt.Sprint(v))
	}
	return r
}

func main() {
	source := `
["step",
  ["set", "i", 5],
  ["set", "sum", 0],
  ["until", ["=", ["get", "i"], 0], [
    "step",
    ["set", "sum", ["+", ["get", "sum"], ["get", "i"]]],
    ["set", "i", ["+", ["get", "i"], -1]]
  ]],
  ["get", "sum"]
]`

	var v interface{}
	err := json.Unmarshal([]byte(source), &v)
	if err != nil {
		log.Fatal(err)
	}

	env := make(map[string]interface{})

	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()
	fmt.Println(doRun(env, v.([]interface{})))
}
