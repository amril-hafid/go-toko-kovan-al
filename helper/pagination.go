package helper

import "html/template"

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"seq": func(start, end int) []int {
			var s []int
			for i := start; i <= end; i++ {
				s = append(s, i)
			}
			return s
		},
		"min": func(a, b int) int {
			if a < b {
				return a
			}
			return b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"add": func(a, b int) int {
			return a + b
		},
	}
}
