package my_env

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Environment Context definition and operations

type EnvironmentContext struct {
	Env *Environment
}

type IEnvironmentContext interface {
	GetVariables() *[]*EnvironmentVariable
	GetEnvironment() string
	Setup()
}

func (ctx *EnvironmentContext) Setup() *EnvironmentContext {
	var Env *Environment
	files, _ := filepath.Glob(".env.*")
	Env = Env.create(files[0])
	err := Env.read()
	if err != nil {
		return nil
	}
	err = Env.set()
	if err != nil {
		return nil
	}
	return &EnvironmentContext{Env: Env}
}

func (ctx *EnvironmentContext) Get(key string) string {
	val, _ := os.LookupEnv(key)
	return val
}

func (ctx *EnvironmentContext) Breakdown() {
	_ = ctx.Env.unSet()
}

func (ctx *EnvironmentContext) GetVariables() *[]*EnvironmentVariable {
	return &ctx.Env.variables
}

func (ctx *EnvironmentContext) GetEnvironment() string {
	return ctx.Env.branch
}

// Environment definition and operations

type Environment struct {
	branch    string
	variables []*EnvironmentVariable
	path      string
}

func (env *Environment) create(path string) *Environment {
	return &Environment{branch: strings.SplitAfter(path, "env.")[1], path: path, variables: make([]*EnvironmentVariable, 0)}
}

func (env *Environment) set() error {
	err := os.Setenv("local-environment", env.branch)
	for _, variable := range env.variables {
		err = os.Setenv(variable.key, variable.value)
		if err != nil {
			return err
		}
	}
	return err
}

func (env *Environment) unSet() error {
	var err error
	for _, variable := range env.variables {
		err = os.Unsetenv(variable.key)
		if err != nil {
			continue
		}
	}
	err = os.Unsetenv("local-environment")
	return err
}

func (env *Environment) add(variable *EnvironmentVariable) {
	env.variables = append(env.variables, variable)
}

func (env *Environment) read() error {
	var err error

	content, err := os.ReadFile(env.path)
	if err != nil {
		return err
	}

	var unmarshalled map[string]string

	err = json.Unmarshal(content, &unmarshalled)

	for key, value := range unmarshalled {
		var newEnv *EnvironmentVariable
		env.add(newEnv.create(key, value))
	}

	return err
}

// Variable definition and operations

type EnvironmentVariable struct {
	key   string
	value string
}

func (ev *EnvironmentVariable) create(key string, val string) *EnvironmentVariable {
	return &EnvironmentVariable{key: key, value: val}
}

func (ev *EnvironmentVariable) getValue() (string, error) {
	if len(ev.value) == 0 {
		return "", errors.New("no Value set in variable")
	} else {
		return ev.value, nil
	}
}

func (ev *EnvironmentVariable) GetKey() string {
	return ev.key
}
