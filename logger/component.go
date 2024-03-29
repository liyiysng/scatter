package logger

import (
	"fmt"
)

// componentData records the settings for a component.
type componentData struct {
	name string
}

const componentDepth = 3

var cache = map[string]*componentData{}

func (c *componentData) InfoDepth(depth int, args ...interface{}) {
	args = append([]interface{}{"[" + string(c.name) + "]"}, args...)
	InfoDepth(depth+1, args...)
}

func (c *componentData) WarningDepth(depth int, args ...interface{}) {
	args = append([]interface{}{"[" + string(c.name) + "]"}, args...)
	WarningDepth(depth+1, args...)
}

func (c *componentData) ErrorDepth(depth int, args ...interface{}) {
	args = append([]interface{}{"[" + string(c.name) + "]"}, args...)
	ErrorDepth(depth+1, args...)
}

func (c *componentData) FatalDepth(depth int, args ...interface{}) {
	args = append([]interface{}{"[" + string(c.name) + "]"}, args...)
	FatalDepth(depth+1, args...)
}

func (c *componentData) Info(args ...interface{}) {
	c.InfoDepth(componentDepth, args...)
}

func (c *componentData) Warning(args ...interface{}) {
	c.WarningDepth(componentDepth, args...)
}

func (c *componentData) Error(args ...interface{}) {
	c.ErrorDepth(componentDepth, args...)
}

func (c *componentData) Fatal(args ...interface{}) {
	c.FatalDepth(componentDepth, args...)
}

func (c *componentData) Infof(format string, args ...interface{}) {
	c.InfoDepth(componentDepth, fmt.Sprintf(format, args...))
}

func (c *componentData) Warningf(format string, args ...interface{}) {
	c.WarningDepth(componentDepth, fmt.Sprintf(format, args...))
}

func (c *componentData) Errorf(format string, args ...interface{}) {
	c.ErrorDepth(componentDepth, fmt.Sprintf(format, args...))
}

func (c *componentData) Fatalf(format string, args ...interface{}) {
	c.FatalDepth(componentDepth, fmt.Sprintf(format, args...))
}

func (c *componentData) Infoln(args ...interface{}) {
	c.InfoDepth(componentDepth, args...)
}

func (c *componentData) Warningln(args ...interface{}) {
	c.WarningDepth(componentDepth, args...)
}

func (c *componentData) Errorln(args ...interface{}) {
	c.ErrorDepth(componentDepth, args...)
}

func (c *componentData) Fatalln(args ...interface{}) {
	c.FatalDepth(componentDepth, args...)
}

func (c *componentData) V(l int) bool {
	return V(l)
}

// Component creates a new component and returns it for logging. If a component
// with the name already exists, nothing will be created and it will be
// returned. SetLogger will panic if it is called with a logger created by
// Component.
func Component(componentName string) DepthLogger {
	if cData, ok := cache[componentName]; ok {
		return cData
	}
	c := &componentData{componentName}
	cache[componentName] = c
	return c
}
