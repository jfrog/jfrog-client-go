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

func NewParenthesesSlice(slice []Parentheses) ParenthesesSlice {
	return ParenthesesSlice{Parentheses: slice}
}

func CreateParenthesesSlice(pattern, target string) ParenthesesSlice {
	return ParenthesesSlice{findParentheses(pattern, target)}
}

func (p *ParenthesesSlice) IsPresent(index int) bool {
	for _, v := range p.Parentheses {
		if v.OpenIndex == index || v.CloseIndex == index {
			return true
		}
	}
	return false
}

// Return true if at least one of the {i} in 'target' has corresponding parentheses in 'pattern'.
func IsPlaceholdersUsed(pattern, target string) bool {
	removedParenthesesTarget := RemovePlaceholderParentheses(pattern, target)
	return removedParenthesesTarget != target
}

func RemovePlaceholderParentheses(pattern, target string) string {
	parentheses := CreateParenthesesSlice(pattern, target)
	// Remove parentheses which have a corresponding placeholder.
	var temp string
	for i, c := range pattern {
		if (c == '(' || c == ')') && parentheses.IsPresent(i) {
			continue
		} else {
			temp = temp + string(c)
		}
	}
	return temp
}

// addEscapingParentheses escapes parentheses with no corresponding placeholder.
func addEscapingParentheses(pattern, target string) string {
	return AddEscapingParentheses(pattern, target, "")
}

// AddEscapingParentheses escapes parentheses with no corresponding placeholder.
// pattern - the pattern in which the parentheses are escaped.
// target - target parameter containing placeholders.
// targetPathInArchive - The target archive path containing placeholders (relevant only for upload commands).
func AddEscapingParentheses(pattern, target, targetPathInArchive string) string {
	parentheses := CreateParenthesesSlice(pattern, target)
	archiveParentheses := CreateParenthesesSlice(pattern, targetPathInArchive)
	var temp string
	for i, c := range pattern {
		if (c == '(' || c == ')') && !parentheses.IsPresent(i) && !archiveParentheses.IsPresent(i) {
			temp = temp + "\\" + string(c)
		} else {
			temp = temp + string(c)
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
	parentheses := getAllParentheses(pattern)
	// Filter out parentheses without placeholders
	var result []Parentheses
	for _, placeHolderValueIndex := range getPlaceHoldersValues(target) {
		if len(parentheses) > placeHolderValueIndex-1 {
			result = append(result, parentheses[placeHolderValueIndex-1])
		}
	}
	return result
}

// Find the list of Parentheses in the pattern.
func getAllParentheses(pattern string) []Parentheses {
	// Save each parentheses index
	var parentheses []Parentheses
	for i, char := range pattern {
		if char == '(' {
			parentheses = append(parentheses, Parentheses{i, 0})
		}
		if char == ')' {
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
	return temp
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
