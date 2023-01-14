package hegex

import (
	"fmt"
	"regexp"
	"runtime"
	"testing"
)

// simple expressions
var regexSimpleExpression = `(?P<site>[^\s\./]+)\.example\.com`
var hegexSimpleExpression = `{site}.example.com`

func BenchmarkSimple(b *testing.B) {
	cpuModel := runtime.GOARCH
	fmt.Println("CPU Model:", cpuModel)
	rp := regexp.MustCompile(regexSimpleExpression)

	b.Run("regex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s := fmt.Sprintf("site-%d.example.com", i)
			ok := rp.MatchString(s)
			if !ok {
				b.Fatalf("%s doesn't match %s", regexSimpleExpression, s)
			}
		}
	})

	hp := MustCompile(hegexSimpleExpression)

	b.Run("hegex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s := fmt.Sprintf("site-%d.example.com", i)
			ok := hp.MatchString(s)
			if !ok {
				b.Fatalf("%s doesn't match %s", hegexSimpleExpression, s)
			}
		}
	})

}

// complex expressions
var regexComplexExpression = `(?P<prefix>[^\s\./]+)\.(?P<options>opt-000|opt-100|opt-200)\.example\.(?P<asteriskgroup1>.*)\.whAt-a-C0mplex-str_ing\.(?P<asteriskgroup2>.*)\.(?P<suffix>[^\s\./]+)`
var hegexComplexExpression = `{prefix}.{options[opt-000|opt-100|opt-200]}.example.*.whAt-a-C0mplex-str_ing.**.{suffix}`

func BenchmarkComplex(b *testing.B) {
	rp := regexp.MustCompile(regexComplexExpression)

	b.Run("regex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			remainder := i % 6
			s := fmt.Sprintf("pref1x100.opt-%d00.example.anything01.02Anything.whAt-a-C0mplex-str_ing..com", remainder)
			ok := rp.MatchString(s)
			if remainder < 3 && !ok {
				b.Fatalf("%s should match %s", regexComplexExpression, s)
			} else if remainder >= 3 && ok {
				b.Fatalf("%s shouldn't match %s", regexComplexExpression, s)
			}
		}
	})

	hp := MustCompile(hegexComplexExpression)

	b.Run("hegex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			remainder := i % 6
			s := fmt.Sprintf("pref1x100.opt-%d00.example.anything01.02Anything.whAt-a-C0mplex-str_ing..com", remainder)
			ok := hp.MatchString(s)
			if remainder < 3 && !ok {
				b.Fatalf("%s should match %s", hegexComplexExpression, s)
			} else if remainder >= 3 && ok {
				b.Fatalf("%s shouldn't match %s", hegexComplexExpression, s)
			}
		}
	})

}
