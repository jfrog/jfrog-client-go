package utils

import (
	"regexp"
	"sort"
	"strconv"
)

var placeholderRegExp = regexp.MustCompile(`{([^}]*)}`)

// This struct represents the parentheses used for defining Placeholders (Placeholders is a feature supported by File Specs).
type Parentheses struct {
	OpenIndex  int
	CloseIndex int
}

type ParenthesesSlice struct {
	Parentheses []Parentheses
}

func (p *ParenthesesSlice) IsPresent(index int) bool {
	for _, v := range p.Parentheses {
		if v.OpenIndex == index || v.CloseIndex == index {
			return true
		}
	}
	return false
}

func RemovePlaceholderParentheses(pattern, target string) string {
	parentheses := ParenthesesSlice{findParentheses(pattern, target)}

	// Remove parentheses which have a corresponding placeholder.
	var temp string
	for i := 0; i < len(pattern); i++ {
		if (pattern[i] == '(' || pattern[i] == ')') && parentheses.IsPresent(i) {
			continue
		} else {
			temp = temp + string(pattern[i])
		}

	}
	return temp
}

// Escapoing Parentheses with no corresponding placeholder
func AddEscapingParentheses(pattern, target string) string {
	parentheses := ParenthesesSlice{findParentheses(pattern, target)}
	var temp string
	for i := 0; i < len(pattern); i++ {
		if (pattern[i] == '(' || pattern[i] == ')') && !parentheses.IsPresent(i) {
			temp = temp + "\\" + string(pattern[i])
		} else {
			temp = temp + string(pattern[i])
		}
	}
	return temp
}

func getPlaceHoldersValues(target string) []int {
	var placeholderFound []int
	matches := placeholderRegExp.FindAllStringSubmatch(target, -1)
	for _, v := range matches {
		if number, err := strconv.Atoi(v[1]); err == nil {
			placeholderFound = append(placeholderFound, number)
		}
	}
	if placeholderFound != nil {
		sortNoDuplicates(&placeholderFound)
	}
	return placeholderFound
}

// Find the list of Parentheses in the pattern, which correspond to placeholders defined in the target.
func findParentheses(pattern, target string) []Parentheses {
	// Save each parentheses index
	var parentheses []Parentheses
	for i, v := range pattern {
		if v == '(' {
			parentheses = append(parentheses, Parentheses{i, 0})
		}
		if v == ')' {
			for j := len(parentheses) - 1; j >= 0; j-- {
				if parentheses[j].CloseIndex == 0 {
					parentheses[j].CloseIndex = i
					break
				}
			}
		}
	}

	// Remove open parentheses without closing parenthesis
	var temp []Parentheses
	for i := 0; i < len(parentheses); i++ {
		if parentheses[i].CloseIndex != 0 {
			temp = append(temp, parentheses[i])
		}
	}
	// Filter parentheses without placeholders
	var result []Parentheses
	for _, v := range getPlaceHoldersValues(target) {
		if len(temp) > v-1 {
			result = append(result, temp[v-1])
		}
	}
	return result
}

// Sort array and remove duplicates.
func sortNoDuplicates(arg *[]int) {
	sort.Ints(*arg)
	j := 0
	for i := 1; i < len(*arg); i++ {
		if (*arg)[j] == (*arg)[i] {
			continue
		}
		j++
		(*arg)[j] = (*arg)[i]
	}
	*arg = (*arg)[:j+1]
}
