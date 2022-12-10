package hegex

import (
	"fmt"
	"regexp"
	"runtime"
	"testing"
)

// simple patterns
var regexSimplePattern = `(?P<site>[^\s\./]+)\.example\.com`
var hegexSimplePattern = `{site}.example.com`

func BenchmarkSimple(b *testing.B) {
	cpuModel := runtime.GOARCH
	fmt.Println("CPU Model:", cpuModel)
	rp := regexp.MustCompile(regexSimplePattern)

	b.Run("regex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s := fmt.Sprintf("site-%d.example.com", i)
			ok := rp.MatchString(s)
			if !ok {
				b.Fatalf("%s doesn't match %s", regexSimplePattern, s)
			}
		}
	})

	hp := MustCompile(hegexSimplePattern)

	b.Run("hegex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s := fmt.Sprintf("site-%d.example.com", i)
			ok := hp.MatchString(s)
			if !ok {
				b.Fatalf("%s doesn't match %s", hegexSimplePattern, s)
			}
		}
	})

}

// complex patterns
var regexComplexPattern = `(?P<prefix>[^\s\./]+)\.(?P<candidates>c-000|c-100|c-200)\.example\.(?P<asteriskgroup1>.*)\.whAt-a-C0mplex-str_ing\.(?P<asteriskgroup2>.*)\.(?P<suffix>[^\s\./]+)`
var hegexComplexPattern = `{prefix}.{candidates[c-000|c-100|c-200]}.example.*.whAt-a-C0mplex-str_ing.**.{suffix}`

func BenchmarkComplex(b *testing.B) {
	rp := regexp.MustCompile(regexComplexPattern)

	b.Run("regex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			remainder := i % 6
			s := fmt.Sprintf("pref1x100.c-%d00.example.anything01.02Anything.whAt-a-C0mplex-str_ing..com", remainder)
			ok := rp.MatchString(s)
			if remainder < 3 && !ok {
				b.Fatalf("%s should match %s", regexComplexPattern, s)
			} else if remainder >= 3 && ok {
				b.Fatalf("%s shouldn't match %s", regexComplexPattern, s)
			}
		}
	})

	hp := MustCompile(hegexComplexPattern)

	b.Run("hegex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			remainder := i % 6
			s := fmt.Sprintf("pref1x100.c-%d00.example.anything01.02Anything.whAt-a-C0mplex-str_ing..com", remainder)
			ok := hp.MatchString(s)
			if remainder < 3 && !ok {
				b.Fatalf("%s should match %s", hegexComplexPattern, s)
			} else if remainder >= 3 && ok {
				b.Fatalf("%s shouldn't match %s", hegexComplexPattern, s)
			}
		}
	})

}
