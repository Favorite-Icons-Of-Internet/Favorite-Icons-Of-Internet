package main

import (
	"container/list"
	"encoding/json"
)

func isJsonObj(data []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(data, &js) == nil
}

func isBalancedBrackets(data []byte) bool {
	var (
		stack *list.List
		first *list.Element
	)

	stack = list.New()
	for _, ch := range data {
		if ch == '{' || ch == '[' {
			first = stack.PushFront(ch)
		} else if ch == '}' {
			if first == nil || first.Value.(byte) != '{' {
				return false
			} else {
				stack.Remove(first)
				first = stack.Front()
			}
		} else if ch == ']' {
			if first == nil || first.Value.(byte) != '[' {
				return false
			} else {
				stack.Remove(first)
				first = stack.Front()
			}
		}
	}

	return stack.Len() == 0
}