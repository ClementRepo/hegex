package hegex

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Hegexp struct {
	// expression of hegex
	expression string
	// consecutive * found in expression
	as []asterisk
	// curly brace pairs found in expression
	bs []cBrace
	// expression is translated into regex under the hood
	re                     *regexp.Regexp
	regexGroupToHegexGroup map[string]string
}

// Compile parses a hegex expression and returns, if successful,
// a Hegexp object that can be used to match text.
func Compile(expression string) (*Hegexp, error) {
	h, err := newHegex(expression)
	if err != nil {
		return nil, err
	}
	return h, nil
}

// MustCompile is like Compile but panics if the expression cannot be parsed.
func MustCompile(str string) *Hegexp {
	h, err := newHegex(str)
	if err != nil {
		panic(`hegexp: Compile(` + strconv.Quote(str) + `): ` + err.Error())
	}
	return h
}

// MatchString reports whether the string s contains any match of the Hegexp.
func (h *Hegexp) MatchString(s string) bool {
	return h.re.MatchString(s)
}

// MatchString reports whether the string s contains any match of the Hegexp.
// return an error if expression can't be compiled into a Hegexp
func MatchString(expression string, s string) (ok bool, err error) {
	h, err := Compile(expression)
	if err != nil {
		return false, err
	}
	return h.MatchString(s), nil
}

// MatchAndFindStringSubmatch returns a map of strings holding the text of the
// leftmost match of the hegex expression in s.
// The key of map corresponds to the name of the submatch within the match.
// A return value of (nil,false) means that the expression and s did not match
// A return value of (nil,true) means that the expression and s matched, but no submatches were found
func (h *Hegexp) MatchAndFindStringSubmatch(s string) (group map[string]string, ok bool) {
	submatches := h.re.FindStringSubmatch(s)
	if len(submatches) == 0 {
		return nil, false
	}
	if len(submatches) == 1 {
		return nil, true
	}
	submatchMap := map[string]string{}
	for i, sm := range h.re.SubexpNames() {
		if i != 0 {
			// i==0 or sm=="" means this group match the whole string instead of part of the string
			hName := h.regexGroupToHegexGroup[sm]
			submatchMap[hName] = submatches[i]
		}
	}
	return submatchMap, true
}

// MatchAndSubstitute returns a new string where {} and * in template have been substituted.
// If the Hegexp matches 's', it returns the substituted template and 'true'.
// Otherwise, it returns empty string and 'false'.
func (h *Hegexp) MatchAndSubstitute(s string, template string) (substituted string, ok bool) {
	group, match := h.MatchAndFindStringSubmatch(s)
	if !match {
		return "", false
	}
	substituted = template
	// group with longer name should be used applied to template earlier
	// this ensures ** in template has higher priority in template
	var keys []string
	for k := range group {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})

	for _, k := range keys {
		v := group[k]
		if len(k) == 0 || strings.Count(k, "*") != len(k) {
			// if group name relates to curly brace
			k = fmt.Sprintf("{%s}", k)
		}
		substituted = strings.ReplaceAll(substituted, k, v)
	}
	return substituted, true
}

func newHegex(expression string) (*Hegexp, error) {
	as := findAsterix(expression)

	cbs, err := findCBrace(expression)
	if err != nil {
		return nil, err
	}

	exprCopy := expression
	var parts []string
	i, j := len(as)-1, len(cbs)-1
	for i >= 0 && j >= 0 {
		a := as[i]
		cb := cbs[j]
		if a.endExclusive > cb.endExclusive {
			parts = append([]string{a.groupRegex(), regexp.QuoteMeta(exprCopy[a.endExclusive:])}, parts...)
			exprCopy = exprCopy[:a.start]
			i--
		} else {
			parts = append([]string{cb.groupRegex(), regexp.QuoteMeta(exprCopy[cb.endExclusive:])}, parts...)
			exprCopy = exprCopy[:cb.start]
			j--
		}
	}
	for i >= 0 {
		a := as[i]
		parts = append([]string{a.groupRegex(), regexp.QuoteMeta(exprCopy[a.endExclusive:])}, parts...)
		exprCopy = exprCopy[:a.start]
		i--
	}
	for j >= 0 {
		cb := cbs[j]
		parts = append([]string{cb.groupRegex(), regexp.QuoteMeta(exprCopy[cb.endExclusive:])}, parts...)
		exprCopy = exprCopy[:cb.start]
		j--
	}

	parts = append([]string{exprCopy}, parts...)
	join := strings.Join(parts, "")
	join = fmt.Sprintf("^%s$", join)
	re, err := regexp.Compile(join)
	if err != nil {
		return nil, err
	}

	regexGroupToHegexGroup := map[string]string{}
	for _, a := range as {
		regexGroupToHegexGroup[a.regexGroupName()] = strings.Repeat("*", a.len())
	}
	for _, cb := range cbs {
		regexGroupToHegexGroup[cb.regexGroupName()] = cb.hegexGroupName
	}

	return &Hegexp{expression, as, cbs, re, regexGroupToHegexGroup}, nil
}

const asteriskGroupPrefix = "hegexAsteriskGroup"
const emptyGroupName = "hegexEmptyGroupName"

type asterisk struct {
	start        int
	endExclusive int
}

func (a asterisk) len() int {
	return a.endExclusive - a.start
}

func (a asterisk) regexGroupName() string {
	return fmt.Sprintf("%s%d", asteriskGroupPrefix, a.len())
}

func (a asterisk) groupRegex() string {
	return fmt.Sprintf("(?P<%s>%s)", a.regexGroupName(), ".*")
}

type cBrace struct {
	start          int
	endExclusive   int
	hegexGroupName string
	options        []string
}

func (cb cBrace) Less(bb cBrace) bool {
	return cb.start < bb.start
}

func (cb cBrace) groupRegex() string {
	if len(cb.options) == 0 {
		// this has a drawback that /{path}/subpath cannot match "/a.txt/subpath"
		// but I'm not going to fix it now
		return fmt.Sprintf("(?P<%s>%s)", cb.regexGroupName(), "[^\\s\\./]+")
	} else {
		join := strings.Join(cb.options, "|")
		return fmt.Sprintf("(?P<%s>%s)", cb.regexGroupName(), join)
	}
}

func (cb cBrace) regexGroupName() string {
	if cb.hegexGroupName == "" {
		return emptyGroupName
	}
	return cb.hegexGroupName
}

func findAsterix(expression string) []asterisk {
	var as []asterisk
	// find all asterisk in expr
	for i := 0; i < len(expression); i++ {
		if expression[i] == '*' {
			j := i + 1
			for ; j < len(expression) && expression[j] == '*'; j++ {
			}
			a := asterisk{start: i, endExclusive: j}
			as = append(as, a)
			i = j
		}
	}
	return as
}

const cBraceFormatRegex = `{[A-Za-z0-9-]*(\[([A-Za-z0-9-]+\|)*[A-Za-z0-9-]+\])?}`

func findCBrace(expression string) ([]cBrace, error) {
	var cbs []cBrace
	var stack []int
	// find cBrace range in expr
	for i := 0; i < len(expression); i++ {
		if expression[i] == '{' {
			stack = append(stack, i)
		} else if expression[i] == '}' && len(stack) > 0 {
			// pop stack
			if top := stack[len(stack)-1]; expression[top] == '{' {
				b := cBrace{start: top, endExclusive: i + 1}
				cbs = append(cbs, b)
			}
		}
	}

	if len(cbs) > 0 {
		// ensure curly brace does not overlap
		sort.Slice(cbs, func(i, j int) bool {
			return cbs[i].start < cbs[j].start
		})
		b := cbs[0]
		for i := 1; i < len(cbs); i++ {
			bb := cbs[i]
			if b.endExclusive-1 == bb.start {
				msg := fmt.Sprintf("curly brace pairs should not overlap\n"+
					"curly brace pairs overlap in expression: %s\n"+
					"curly brace pair 1: %s\n"+
					"curly brace pair 2: %s\n"+
					"", expression, expression[b.start:b.endExclusive], expression[bb.start:bb.endExclusive])
				return nil, &Error{msg, expression}
			}
		}
	}

	for i := range cbs {
		content := expression[cbs[i].start+1 : cbs[i].endExclusive-1]
		ok, err := regexp.MatchString(cBraceFormatRegex, expression)
		if err != nil {
			return nil, &Error{err.Error(), expression}
		}
		if !ok {
			return nil, &Error{"bad format, didn't match " + cBraceFormatRegex, expression}
		}

		if strings.Contains(content, "[") {
			// extract options
			sp := strings.Split(content, "[")
			cbs[i].hegexGroupName = sp[0]
			opts := strings.TrimSuffix(sp[1], "]")
			cbs[i].options = strings.Split(opts, "|")
		} else {
			cbs[i].hegexGroupName = content
		}
	}
	return cbs, nil
}

type Error struct {
	Reason string
	Expr   string
}

func (e *Error) Error() string {
	return "error parsing regexp: " + e.Reason + ": `" + e.Expr + "`"
}
