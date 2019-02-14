package main

// COMTEXT

type Context map[string]interface{}

type ContextList []Context

type ContextStack []Context

type ValueList []interface{}

func (c ContextStack) PushContext(v Context) ContextStack {
	return append(c, v)
}

func (c ContextStack) FindValue(name string) interface{} {
	for i := len(c) - 1; i >= 0; i-- {
		context := c[i]
		if el, ok := context[name]; ok {
			return el
		}
	}
	return nil
}
