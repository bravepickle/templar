package main

// TODO: either remove this file or make new useful functions (e.g. unescape raw values output) and merge with sprig map. For now use sprig instead

import (
	// "bytes"
	// "text/template"
	"github.com/Masterminds/sprig"
)

// here we store all additional functions for rendering templates

// // sum sums 2 values
// func sum(x int, y int) int {
// 	return x + y
// }

// // sub substracts one value from another
// func sub(x int, y int) int {
// 	return x - y
// }

// // repeat repeats chunk n times
// // func repeat(raw interface{}, repeat int) []byte {
// func repeat(repeat int, raw string) string {
// 	// chunk2, ok := raw.(string)
// 	// // chunk, ok := raw.([]byte)
// 	// if !ok {
// 	// 	log.Fatal(`Failed to parse value: `, raw)
// 	// }
// 	// chunk := []byte(chunk2)
// 	chunk := []byte(raw)
// 	buf := bytes.NewBuffer(chunk)
// 	for i := 0; i < repeat; i++ {
// 		buf.Write(chunk)
// 	}

// 	return buf.String()
// }

// func noescape(str string) template.HTML {
// 	return template.HTML(str)
// }

// var funcMap = template.FuncMap{
// 	"sum":    sum,
// 	"sub":    sub,
// 	"repeat": repeat,
// }

var funcMap = sprig.TxtFuncMap()

// funcMap[`raw`] = append(noescape)
