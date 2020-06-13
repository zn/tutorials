package main

import(
	"testing"
)

func Test_buildStringSlow(t *testing.T){
	if buildStringSlow("lee(t(c)o)de)") != "lee(t(c)o)de" {
		t.Fail()
	}
	if buildStringSlow("((((((()") != "()" {
		t.Fail()
	}
}

func Test_buildStringFast(t *testing.T){
	if buildStringFast("lee(t(c)o)de)") != "lee(t(c)o)de" {
		t.Fail()
	}
	if buildStringFast("((((((()") != "()" {
		t.Fail()
	}
}

var result string

func Benchmark_buildStringSlow(b *testing.B){
	var r string
	for i := 0; i < b.N; i++ {
		r = buildStringSlow("l(((((()()()()))(ee(t(c)o)d(((e(((((((((((()))fdsafdf))))")
	}
	result = r
}

func Benchmark_buildStringFast(b *testing.B){
	var r string
	for i := 0; i < b.N; i++ {
		r = buildStringFast("l(((((()()()()))(ee(t(c)o)d(((e(((((((((((()))fdsafdf))))")
	}
	result = r
}