package env

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
	Env = Env.Create(files[0])
	err := Env.Read()
	if err != nil {
		return nil
	}
	err = Env.Set()
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
	_ = ctx.Env.UnSet()
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

type IEnvironment interface {
	Set() error
	UnSet() error
}

func (env *Environment) Create(path string) *Environment {
	return &Environment{branch: strings.SplitAfter(path, "env.")[1], path: path, variables: make([]*EnvironmentVariable, 0)}
}

func (env *Environment) Set() error {
	err := os.Setenv("local-environment", env.branch)
	for _, variable := range env.variables {
		err = os.Setenv(variable.key, variable.value)
		if err != nil {
			return err
		}
	}
	return err
}

func (env *Environment) UnSet() error {
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

func (env *Environment) Add(variable *EnvironmentVariable) {
	env.variables = append(env.variables, variable)
}

func (env *Environment) Read() error {
	var err error

	content, err := os.ReadFile(env.path)
	if err != nil {
		return err
	}

	var unmarshalled map[string]string

	err = json.Unmarshal(content, &unmarshalled)

	for key, value := range unmarshalled {
		var newEnv *EnvironmentVariable
		env.Add(newEnv.Create(key, value))
	}

	return err
}

// Variable definition and operations

type EnvironmentVariable struct {
	key   string
	value string
}

type IEnvironmentVariable interface {
	Create(string, string) *EnvironmentVariable
	Value() string
}

func (ev *EnvironmentVariable) Create(key string, val string) *EnvironmentVariable {
	return &EnvironmentVariable{key: key, value: val}
}

func (ev *EnvironmentVariable) Value() (string, error) {
	if len(ev.value) == 0 {
		return "", errors.New("no Value set in variable")
	} else {
		return ev.value, nil
	}
}
