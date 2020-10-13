package main

import (
	"fmt"
	"reflect"
	"strings"
)

// IsNotDigit ...
func IsNotDigit(c rune) bool { 
	return (c < '0') || (c > '9' )
}

// IsDigit ...
func IsDigit(c rune) bool { 
	return (c >= '0') && (c <= '9' )
}

// AreAllDigits ...
func AreAllDigits(s string) bool {
	return strings.IndexFunc(s, IsNotDigit) == -1
}

// NoDigits ...
func NoDigits(s string) bool {
	return strings.IndexFunc(s, IsDigit) == -1
}

// RelectExaminer ...
func RelectExaminer(t reflect.Type, depth int) {
	fmt.Println(strings.Repeat("\t", depth), "Type is", t.Name(), "and kind is", t.Kind())
	switch t.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		fmt.Println(strings.Repeat("\t", depth+1), "Contained type:")
		RelectExaminer(t.Elem(), depth+1)
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fmt.Println(strings.Repeat("\t", depth+1), "Field", i+1, "name is", f.Name, "type is", f.Type.Name(), "and kind is", f.Type.Kind())
			if f.Tag != "" {
				fmt.Println(strings.Repeat("\t", depth+2), "Tag is", f.Tag)
				fmt.Println(strings.Repeat("\t", depth+2), "tag1 is", f.Tag.Get("tag1"), "tag2 is", f.Tag.Get("tag2"))
			}
		}
	}
}


func sampleReflectExaminer() {
	sl := []int{1, 2, 3}
	slType := reflect.TypeOf(sl)
	RelectExaminer(slType, 0)
}
