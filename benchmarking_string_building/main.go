package main

import "strings"

func buildStringSlow(s string) string{
	removeMap := minRemoveToMakeValid(s)
	newString := ""
	for i := 0; i < len(s); i++ {
		if _, ok := removeMap[i]; !ok {
			newString += string(s[i])
		}
	}
	return newString
}

func buildStringFast(s string) string {
	removeMap := minRemoveToMakeValid(s)
	var builder strings.Builder
	builder.Grow(len(s) - len(removeMap))
	for i := 0; i < len(s); i++ {
		if _, ok := removeMap[i]; !ok {
			builder.WriteByte(s[i])
		}
	}
	return builder.String()
}

func minRemoveToMakeValid(s string) map[int]*struct{} {
	stack := make([]int, 0)
	toRemove := make([]int, 0)

	for i := 0; i < len(s); i++ {
		if string(s[i]) == "(" {
			stack = append(stack, i)
		} else if string(s[i]) == ")" {
			if len(stack) == 0 {
				toRemove = append(toRemove, i)
			} else {
				stack = stack[:len(stack)-1]
			}
		}
	}

	toRemove = append(toRemove, stack...)
	removeMap := make(map[int]*struct{})
	for _, i := range toRemove {
		removeMap[i] = nil
	}
	return removeMap
}



func main() {
}
