package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Env string

func (e Env) String() string {
	return string(e)
}

func (e Env) Split(sep string) (r []string) {
	for _, el := range strings.Split(string(e), sep) {
		if el != "" {
			r = append(r, el)
		}
	}
	return r
}

func (e Env) Int() int {
	i, err := strconv.Atoi(string(e))
	if err != nil {
		panic(fmt.Sprintf("Error converting env %v to int: %s", string(e), err.Error()))
	}
	return i
}

func getEnv(key string) Env {
	val, _ := os.LookupEnv(key)
	return Env(val)
}