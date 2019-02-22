package main

// COMTEXT

// type Context
// map[string]interface{}

type ContextStack []map[string]interface{}

func (c ContextStack) PushContext(v map[string]interface{}) ContextStack {
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
