/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-16 11:11:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-18 11:58:15
 * @FilePath: \go-argus-benchmark\bench_alloc_test.go
 * @Description: Argus vs Playground 性能基准测试 — 基线 / Field / Struct 场景
 *
 * 每个场景按 Argus → Playground 紧挨排列，便于上下对比查阅。
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package bench

import (
	"reflect"
	"testing"
	"time"

	playground "github.com/go-playground/validator/v10"
	argus "github.com/kamalyes/go-argus"
)

var (
	argusV      *argus.Validate
	playgroundV *playground.Validate
)

func init() {
	argusV = argus.New()
	playgroundV = playground.New()
}

func benchVarArgus(b *testing.B, val, tag string) {
	b.Helper()
	for i := 0; i < b.N; i++ {
		_ = argusV.Var(val, tag)
	}
}

func benchVarStringArgus(b *testing.B, val, tag string) {
	b.Helper()
	for i := 0; i < b.N; i++ {
		_ = argusV.VarString(val, tag)
	}
}

func benchVarPlayground(b *testing.B, val, tag string) {
	b.Helper()
	guardPlaygroundBenchmark(b, func() {
		_ = playgroundV.Var(val, tag)
	})
	for i := 0; i < b.N; i++ {
		_ = playgroundV.Var(val, tag)
	}
}

func guardPlaygroundBenchmark(b *testing.B, call func()) {
	b.Helper()
	b.StopTimer()
	defer func() {
		if r := recover(); r != nil {
			b.Skipf("skip Playground benchmark: %v", r)
		}
	}()
	call()
	b.ResetTimer()
	b.StartTimer()
}

// ══════════════════════════════════════════════════════════════════════════════
// Reflection Baseline
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_ReflectValueOfInterface(b *testing.B) {
	s := "hello"
	for i := 0; i < b.N; i++ {
		_ = reflect.ValueOf(s)
	}
}

func BenchmarkArgus_ReflectValueOfString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = reflect.ValueOf("hello")
	}
}

// ══════════════════════════════════════════════════════════════════════════════
// VarString vs Var vs Playground Path
// ══════════════════════════════════════════════════════════════════════════════

// ─── Sequential ───

func BenchmarkArgus_VarStringPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = argusV.VarString("hello", "required")
	}
}

func BenchmarkArgus_VarPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = argusV.Var("hello", "required")
	}
}

func BenchmarkPlayground_VarPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = playgroundV.Var("hello", "required")
	}
}

// ─── Parallel ───

func BenchmarkArgus_VarStringPathParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.VarString("hello", "required")
		}
	})
}

func BenchmarkArgus_VarPathParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Var("hello", "required")
		}
	})
}

func BenchmarkPlayground_VarPathParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Var("hello", "required")
		}
	})
}

// ══════════════════════════════════════════════════════════════════════════════
// Field Var Tests
// ══════════════════════════════════════════════════════════════════════════════

// ─── Field Success ───

func BenchmarkArgus_FieldSuccess(b *testing.B) {
	s := "1"
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = argusV.Var(&s, "len=1")
	}
}

func BenchmarkPlayground_FieldSuccess(b *testing.B) {
	s := "1"
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = playgroundV.Var(&s, "len=1")
	}
}

// ─── Field Success Parallel ───

func BenchmarkArgus_FieldSuccessParallel(b *testing.B) {
	s := "1"
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Var(&s, "len=1")
		}
	})
}

func BenchmarkPlayground_FieldSuccessParallel(b *testing.B) {
	s := "1"
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Var(&s, "len=1")
		}
	})
}

// ─── Field Failure ───

func BenchmarkArgus_FieldFailure(b *testing.B) {
	s := "12"
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = argusV.Var(&s, "len=1")
	}
}

func BenchmarkPlayground_FieldFailure(b *testing.B) {
	s := "12"
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = playgroundV.Var(&s, "len=1")
	}
}

// ─── Field Failure Parallel ───

func BenchmarkArgus_FieldFailureParallel(b *testing.B) {
	s := "12"
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Var(&s, "len=1")
		}
	})
}

func BenchmarkPlayground_FieldFailureParallel(b *testing.B) {
	s := "12"
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Var(&s, "len=1")
		}
	})
}

// ══════════════════════════════════════════════════════════════════════════════
// Field Array Dive Tests
// ══════════════════════════════════════════════════════════════════════════════

// ─── ArrayDive Success ───

func BenchmarkArgus_FieldArrayDiveSuccess(b *testing.B) {
	m := []string{"val1", "val2", "val3"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = argusV.Var(m, "required,dive,required")
	}
}

func BenchmarkPlayground_FieldArrayDiveSuccess(b *testing.B) {
	m := []string{"val1", "val2", "val3"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = playgroundV.Var(m, "required,dive,required")
	}
}

// ─── ArrayDive Success Parallel ───

func BenchmarkArgus_FieldArrayDiveSuccessParallel(b *testing.B) {
	m := []string{"val1", "val2", "val3"}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Var(m, "required,dive,required")
		}
	})
}

func BenchmarkPlayground_FieldArrayDiveSuccessParallel(b *testing.B) {
	m := []string{"val1", "val2", "val3"}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Var(m, "required,dive,required")
		}
	})
}

// ─── ArrayDive Failure ───

func BenchmarkArgus_FieldArrayDiveFailure(b *testing.B) {
	m := []string{"val1", "", "val3"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = argusV.Var(m, "required,dive,required")
	}
}

func BenchmarkPlayground_FieldArrayDiveFailure(b *testing.B) {
	m := []string{"val1", "", "val3"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = playgroundV.Var(m, "required,dive,required")
	}
}

// ─── ArrayDive Failure Parallel ───

func BenchmarkArgus_FieldArrayDiveFailureParallel(b *testing.B) {
	m := []string{"val1", "", "val3"}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Var(m, "required,dive,required")
		}
	})
}

func BenchmarkPlayground_FieldArrayDiveFailureParallel(b *testing.B) {
	m := []string{"val1", "", "val3"}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Var(m, "required,dive,required")
		}
	})
}

// ══════════════════════════════════════════════════════════════════════════════
// Field Map Dive Tests
// ══════════════════════════════════════════════════════════════════════════════

// ─── MapDive Success ───

func BenchmarkArgus_FieldMapDiveSuccess(b *testing.B) {
	m := map[string]string{"val1": "val1", "val2": "val2", "val3": "val3"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = argusV.Var(m, "required,dive,required")
	}
}

func BenchmarkPlayground_FieldMapDiveSuccess(b *testing.B) {
	m := map[string]string{"val1": "val1", "val2": "val2", "val3": "val3"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = playgroundV.Var(m, "required,dive,required")
	}
}

// ─── MapDive Success Parallel ───

func BenchmarkArgus_FieldMapDiveSuccessParallel(b *testing.B) {
	m := map[string]string{"val1": "val1", "val2": "val2", "val3": "val3"}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Var(m, "required,dive,required")
		}
	})
}

func BenchmarkPlayground_FieldMapDiveSuccessParallel(b *testing.B) {
	m := map[string]string{"val1": "val1", "val2": "val2", "val3": "val3"}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Var(m, "required,dive,required")
		}
	})
}

// ─── MapDive Failure ───

func BenchmarkArgus_FieldMapDiveFailure(b *testing.B) {
	m := map[string]string{"": "", "val3": "val3"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = argusV.Var(m, "required,dive,required")
	}
}

func BenchmarkPlayground_FieldMapDiveFailure(b *testing.B) {
	m := map[string]string{"": "", "val3": "val3"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = playgroundV.Var(m, "required,dive,required")
	}
}

// ─── MapDive Failure Parallel ───

func BenchmarkArgus_FieldMapDiveFailureParallel(b *testing.B) {
	m := map[string]string{"": "", "val3": "val3"}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Var(m, "required,dive,required")
		}
	})
}

func BenchmarkPlayground_FieldMapDiveFailureParallel(b *testing.B) {
	m := map[string]string{"": "", "val3": "val3"}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Var(m, "required,dive,required")
		}
	})
}

// ══════════════════════════════════════════════════════════════════════════════
// Field Map Dive With Keys Tests
// ══════════════════════════════════════════════════════════════════════════════

// ─── MapDiveWithKeys Success ───

func BenchmarkArgus_FieldMapDiveWithKeysSuccess(b *testing.B) {
	m := map[string]string{"val1": "val1", "val2": "val2", "val3": "val3"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = argusV.Var(m, "required,dive,keys,required,endkeys,required")
	}
}

func BenchmarkPlayground_FieldMapDiveWithKeysSuccess(b *testing.B) {
	m := map[string]string{"val1": "val1", "val2": "val2", "val3": "val3"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = playgroundV.Var(m, "required,dive,keys,required,endkeys,required")
	}
}

// ─── MapDiveWithKeys Success Parallel ───

func BenchmarkArgus_FieldMapDiveWithKeysSuccessParallel(b *testing.B) {
	m := map[string]string{"val1": "val1", "val2": "val2", "val3": "val3"}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Var(m, "required,dive,keys,required,endkeys,required")
		}
	})
}

func BenchmarkPlayground_FieldMapDiveWithKeysSuccessParallel(b *testing.B) {
	m := map[string]string{"val1": "val1", "val2": "val2", "val3": "val3"}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Var(m, "required,dive,keys,required,endkeys,required")
		}
	})
}

// ─── MapDiveWithKeys Failure ───

func BenchmarkArgus_FieldMapDiveWithKeysFailure(b *testing.B) {
	m := map[string]string{"": "", "val3": "val3"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = argusV.Var(m, "required,dive,keys,required,endkeys,required")
	}
}

func BenchmarkPlayground_FieldMapDiveWithKeysFailure(b *testing.B) {
	m := map[string]string{"": "", "val3": "val3"}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = playgroundV.Var(m, "required,dive,keys,required,endkeys,required")
	}
}

// ─── MapDiveWithKeys Failure Parallel ───

func BenchmarkArgus_FieldMapDiveWithKeysFailureParallel(b *testing.B) {
	m := map[string]string{"": "", "val3": "val3"}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Var(m, "required,dive,keys,required,endkeys,required")
		}
	})
}

func BenchmarkPlayground_FieldMapDiveWithKeysFailureParallel(b *testing.B) {
	m := map[string]string{"": "", "val3": "val3"}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Var(m, "required,dive,keys,required,endkeys,required")
		}
	})
}

// ══════════════════════════════════════════════════════════════════════════════
// Field Or Tag Tests
// ══════════════════════════════════════════════════════════════════════════════

// ─── OrTag Success ───

func BenchmarkArgus_FieldOrTagSuccess(b *testing.B) {
	s := "rgba(0,0,0,1)"
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = argusV.Var(s, "rgb|rgba")
	}
}

func BenchmarkPlayground_FieldOrTagSuccess(b *testing.B) {
	s := "rgba(0,0,0,1)"
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = playgroundV.Var(s, "rgb|rgba")
	}
}

// ─── OrTag Success Parallel ───

func BenchmarkArgus_FieldOrTagSuccessParallel(b *testing.B) {
	s := "rgba(0,0,0,1)"
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Var(s, "rgb|rgba")
		}
	})
}

func BenchmarkPlayground_FieldOrTagSuccessParallel(b *testing.B) {
	s := "rgba(0,0,0,1)"
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Var(s, "rgb|rgba")
		}
	})
}

// ─── OrTag Failure ───

func BenchmarkArgus_FieldOrTagFailure(b *testing.B) {
	s := "#000"
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = argusV.Var(s, "rgb|rgba")
	}
}

func BenchmarkPlayground_FieldOrTagFailure(b *testing.B) {
	s := "#000"
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = playgroundV.Var(s, "rgb|rgba")
	}
}

// ─── OrTag Failure Parallel ───

func BenchmarkArgus_FieldOrTagFailureParallel(b *testing.B) {
	s := "#000"
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Var(s, "rgb|rgba")
		}
	})
}

func BenchmarkPlayground_FieldOrTagFailureParallel(b *testing.B) {
	s := "#000"
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Var(s, "rgb|rgba")
		}
	})
}

// ══════════════════════════════════════════════════════════════════════════════
// Struct Tests — Simple
// ══════════════════════════════════════════════════════════════════════════════

type benchSubTest struct {
	Test string `validate:"required"`
}

type benchTestString struct {
	Required  string `validate:"required"`
	Len       string `validate:"len=10"`
	Min       string `validate:"min=1"`
	Max       string `validate:"max=10"`
	MinMax    string `validate:"min=1,max=10"`
	Lt        string `validate:"lt=10"`
	Lte       string `validate:"lte=10"`
	Gt        string `validate:"gt=10"`
	Gte       string `validate:"gte=10"`
	OmitEmpty string `validate:"omitempty,min=1,max=10"`
	Sub       *benchSubTest
	SubIgnore *benchSubTest `validate:"-"`
	Anonymous struct {
		A string `validate:"required"`
	}
}

// ─── StructSimple Success ───

func BenchmarkArgus_StructSimpleSuccess(b *testing.B) {
	type Foo struct {
		StringValue string `validate:"min=5,max=10"`
		IntValue    int    `validate:"min=5,max=10"`
	}
	validFoo := &Foo{StringValue: "Foobar", IntValue: 7}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = argusV.Struct(validFoo)
	}
}

func BenchmarkPlayground_StructSimpleSuccess(b *testing.B) {
	type Foo struct {
		StringValue string `validate:"min=5,max=10"`
		IntValue    int    `validate:"min=5,max=10"`
	}
	validFoo := &Foo{StringValue: "Foobar", IntValue: 7}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = playgroundV.Struct(validFoo)
	}
}

// ─── StructSimple Success Parallel ───

func BenchmarkArgus_StructSimpleSuccessParallel(b *testing.B) {
	type Foo struct {
		StringValue string `validate:"min=5,max=10"`
		IntValue    int    `validate:"min=5,max=10"`
	}
	validFoo := &Foo{StringValue: "Foobar", IntValue: 7}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Struct(validFoo)
		}
	})
}

func BenchmarkPlayground_StructSimpleSuccessParallel(b *testing.B) {
	type Foo struct {
		StringValue string `validate:"min=5,max=10"`
		IntValue    int    `validate:"min=5,max=10"`
	}
	validFoo := &Foo{StringValue: "Foobar", IntValue: 7}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Struct(validFoo)
		}
	})
}

// ─── StructSimple Failure ───

func BenchmarkArgus_StructSimpleFailure(b *testing.B) {
	type Foo struct {
		StringValue string `validate:"min=5,max=10"`
		IntValue    int    `validate:"min=5,max=10"`
	}
	invalidFoo := &Foo{StringValue: "Fo", IntValue: 3}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = argusV.Struct(invalidFoo)
	}
}

func BenchmarkPlayground_StructSimpleFailure(b *testing.B) {
	type Foo struct {
		StringValue string `validate:"min=5,max=10"`
		IntValue    int    `validate:"min=5,max=10"`
	}
	invalidFoo := &Foo{StringValue: "Fo", IntValue: 3}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = playgroundV.Struct(invalidFoo)
	}
}

// ─── StructSimple Failure Parallel ───

func BenchmarkArgus_StructSimpleFailureParallel(b *testing.B) {
	type Foo struct {
		StringValue string `validate:"min=5,max=10"`
		IntValue    int    `validate:"min=5,max=10"`
	}
	invalidFoo := &Foo{StringValue: "Fo", IntValue: 3}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Struct(invalidFoo)
		}
	})
}

func BenchmarkPlayground_StructSimpleFailureParallel(b *testing.B) {
	type Foo struct {
		StringValue string `validate:"min=5,max=10"`
		IntValue    int    `validate:"min=5,max=10"`
	}
	invalidFoo := &Foo{StringValue: "Fo", IntValue: 3}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Struct(invalidFoo)
		}
	})
}

// ══════════════════════════════════════════════════════════════════════════════
// Struct Tests — Complex
// ══════════════════════════════════════════════════════════════════════════════

func newComplexSuccess() *benchTestString {
	return &benchTestString{
		Required:  "Required",
		Len:       "length==10",
		Min:       "min=1",
		Max:       "1234567890",
		MinMax:    "12345",
		Lt:        "012345678",
		Lte:       "0123456789",
		Gt:        "01234567890",
		Gte:       "0123456789",
		OmitEmpty: "",
		Sub:       &benchSubTest{Test: "1"},
		SubIgnore: &benchSubTest{Test: ""},
		Anonymous: struct {
			A string `validate:"required"`
		}{A: "1"},
	}
}

func newComplexFailure() *benchTestString {
	return &benchTestString{
		Required:  "",
		Len:       "",
		Min:       "",
		Max:       "12345678901",
		MinMax:    "",
		Lt:        "0123456789",
		Lte:       "01234567890",
		Gt:        "1",
		Gte:       "1",
		OmitEmpty: "12345678901",
		Sub:       &benchSubTest{Test: ""},
		Anonymous: struct {
			A string `validate:"required"`
		}{A: ""},
	}
}

// ─── StructComplex Success ───

func BenchmarkArgus_StructComplexSuccess(b *testing.B) {
	t := newComplexSuccess()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = argusV.Struct(t)
	}
}

func BenchmarkPlayground_StructComplexSuccess(b *testing.B) {
	t := newComplexSuccess()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = playgroundV.Struct(t)
	}
}

// ─── StructComplex Success Parallel ───

func BenchmarkArgus_StructComplexSuccessParallel(b *testing.B) {
	t := newComplexSuccess()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Struct(t)
		}
	})
}

func BenchmarkPlayground_StructComplexSuccessParallel(b *testing.B) {
	t := newComplexSuccess()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Struct(t)
		}
	})
}

// ─── StructComplex Failure ───

func BenchmarkArgus_StructComplexFailure(b *testing.B) {
	t := newComplexFailure()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = argusV.Struct(t)
	}
}

func BenchmarkPlayground_StructComplexFailure(b *testing.B) {
	t := newComplexFailure()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = playgroundV.Struct(t)
	}
}

// ─── StructComplex Failure Parallel ───

func BenchmarkArgus_StructComplexFailureParallel(b *testing.B) {
	t := newComplexFailure()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Struct(t)
		}
	})
}

func BenchmarkPlayground_StructComplexFailureParallel(b *testing.B) {
	t := newComplexFailure()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Struct(t)
		}
	})
}

// ══════════════════════════════════════════════════════════════════════════════
// Struct Tests — Cross Field
// ══════════════════════════════════════════════════════════════════════════════

// ─── CrossField Success ───

func BenchmarkArgus_StructSimpleCrossFieldSuccess(b *testing.B) {
	type Test struct {
		Start time.Time
		End   time.Time `validate:"gtfield=Start"`
	}
	now := time.Now().UTC()
	then := now.Add(time.Hour * 5)
	test := &Test{Start: now, End: then}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = argusV.Struct(test)
	}
}

func BenchmarkPlayground_StructSimpleCrossFieldSuccess(b *testing.B) {
	type Test struct {
		Start time.Time
		End   time.Time `validate:"gtfield=Start"`
	}
	now := time.Now().UTC()
	then := now.Add(time.Hour * 5)
	test := &Test{Start: now, End: then}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = playgroundV.Struct(test)
	}
}

// ─── CrossField Success Parallel ───

func BenchmarkArgus_StructSimpleCrossFieldSuccessParallel(b *testing.B) {
	type Test struct {
		Start time.Time
		End   time.Time `validate:"gtfield=Start"`
	}
	now := time.Now().UTC()
	then := now.Add(time.Hour * 5)
	test := &Test{Start: now, End: then}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Struct(test)
		}
	})
}

func BenchmarkPlayground_StructSimpleCrossFieldSuccessParallel(b *testing.B) {
	type Test struct {
		Start time.Time
		End   time.Time `validate:"gtfield=Start"`
	}
	now := time.Now().UTC()
	then := now.Add(time.Hour * 5)
	test := &Test{Start: now, End: then}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Struct(test)
		}
	})
}

// ─── CrossField Failure ───

func BenchmarkArgus_StructSimpleCrossFieldFailure(b *testing.B) {
	type Test struct {
		Start time.Time
		End   time.Time `validate:"gtfield=Start"`
	}
	now := time.Now().UTC()
	then := now.Add(time.Hour * -5)
	test := &Test{Start: now, End: then}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = argusV.Struct(test)
	}
}

func BenchmarkPlayground_StructSimpleCrossFieldFailure(b *testing.B) {
	type Test struct {
		Start time.Time
		End   time.Time `validate:"gtfield=Start"`
	}
	now := time.Now().UTC()
	then := now.Add(time.Hour * -5)
	test := &Test{Start: now, End: then}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = playgroundV.Struct(test)
	}
}

// ─── CrossField Failure Parallel ───

func BenchmarkArgus_StructSimpleCrossFieldFailureParallel(b *testing.B) {
	type Test struct {
		Start time.Time
		End   time.Time `validate:"gtfield=Start"`
	}
	now := time.Now().UTC()
	then := now.Add(time.Hour * -5)
	test := &Test{Start: now, End: then}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Struct(test)
		}
	})
}

func BenchmarkPlayground_StructSimpleCrossFieldFailureParallel(b *testing.B) {
	type Test struct {
		Start time.Time
		End   time.Time `validate:"gtfield=Start"`
	}
	now := time.Now().UTC()
	then := now.Add(time.Hour * -5)
	test := &Test{Start: now, End: then}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Struct(test)
		}
	})
}

// ══════════════════════════════════════════════════════════════════════════════
// Struct Tests — Cross Struct Cross Field
// ══════════════════════════════════════════════════════════════════════════════

type benchCrossStructInner struct {
	Start time.Time
}

type benchCrossStructOuter struct {
	Inner     *benchCrossStructInner
	CreatedAt time.Time `validate:"eqcsfield=Inner.Start"`
}

// ─── CrossStructCrossField Success ───

func BenchmarkArgus_StructSimpleCrossStructCrossFieldSuccess(b *testing.B) {
	now := time.Now().UTC()
	inner := &benchCrossStructInner{Start: now}
	outer := &benchCrossStructOuter{Inner: inner, CreatedAt: now}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = argusV.Struct(outer)
	}
}

func BenchmarkPlayground_StructSimpleCrossStructCrossFieldSuccess(b *testing.B) {
	now := time.Now().UTC()
	inner := &benchCrossStructInner{Start: now}
	outer := &benchCrossStructOuter{Inner: inner, CreatedAt: now}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = playgroundV.Struct(outer)
	}
}

// ─── CrossStructCrossField Success Parallel ───

func BenchmarkArgus_StructSimpleCrossStructCrossFieldSuccessParallel(b *testing.B) {
	now := time.Now().UTC()
	inner := &benchCrossStructInner{Start: now}
	outer := &benchCrossStructOuter{Inner: inner, CreatedAt: now}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Struct(outer)
		}
	})
}

func BenchmarkPlayground_StructSimpleCrossStructCrossFieldSuccessParallel(b *testing.B) {
	now := time.Now().UTC()
	inner := &benchCrossStructInner{Start: now}
	outer := &benchCrossStructOuter{Inner: inner, CreatedAt: now}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Struct(outer)
		}
	})
}

// ─── CrossStructCrossField Failure ───

func BenchmarkArgus_StructSimpleCrossStructCrossFieldFailure(b *testing.B) {
	now := time.Now().UTC()
	then := now.Add(time.Hour * 5)
	inner := &benchCrossStructInner{Start: then}
	outer := &benchCrossStructOuter{Inner: inner, CreatedAt: now}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = argusV.Struct(outer)
	}
}

func BenchmarkPlayground_StructSimpleCrossStructCrossFieldFailure(b *testing.B) {
	now := time.Now().UTC()
	then := now.Add(time.Hour * 5)
	inner := &benchCrossStructInner{Start: then}
	outer := &benchCrossStructOuter{Inner: inner, CreatedAt: now}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = playgroundV.Struct(outer)
	}
}

// ─── CrossStructCrossField Failure Parallel ───

func BenchmarkArgus_StructSimpleCrossStructCrossFieldFailureParallel(b *testing.B) {
	now := time.Now().UTC()
	then := now.Add(time.Hour * 5)
	inner := &benchCrossStructInner{Start: then}
	outer := &benchCrossStructOuter{Inner: inner, CreatedAt: now}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Struct(outer)
		}
	})
}

func BenchmarkPlayground_StructSimpleCrossStructCrossFieldFailureParallel(b *testing.B) {
	now := time.Now().UTC()
	then := now.Add(time.Hour * 5)
	inner := &benchCrossStructInner{Start: then}
	outer := &benchCrossStructOuter{Inner: inner, CreatedAt: now}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Struct(outer)
		}
	})
}

// ══════════════════════════════════════════════════════════════════════════════
// Struct Tests — Oneof / Noneof
// ══════════════════════════════════════════════════════════════════════════════

type benchOneof struct {
	Color string `validate:"oneof=red green"`
}

type benchNoneof struct {
	Color string `validate:"noneof=red green"`
}

// ─── Oneof ───

func BenchmarkArgus_Oneof(b *testing.B) {
	w := &benchOneof{Color: "green"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = argusV.Struct(w)
	}
}

func BenchmarkPlayground_Oneof(b *testing.B) {
	w := &benchOneof{Color: "green"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = playgroundV.Struct(w)
	}
}

// ─── Oneof Parallel ───

func BenchmarkArgus_OneofParallel(b *testing.B) {
	w := &benchOneof{Color: "green"}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Struct(w)
		}
	})
}

func BenchmarkPlayground_OneofParallel(b *testing.B) {
	w := &benchOneof{Color: "green"}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Struct(w)
		}
	})
}

// ─── Noneof ───

func BenchmarkArgus_Noneof(b *testing.B) {
	w := &benchNoneof{Color: "blue"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = argusV.Struct(w)
	}
}

func BenchmarkPlayground_Noneof(b *testing.B) {
	w := &benchNoneof{Color: "blue"}
	guardPlaygroundBenchmark(b, func() {
		_ = playgroundV.Struct(w)
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = playgroundV.Struct(w)
	}
}

// ─── Noneof Parallel ───

func BenchmarkArgus_NoneofParallel(b *testing.B) {
	w := &benchNoneof{Color: "blue"}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = argusV.Struct(w)
		}
	})
}

func BenchmarkPlayground_NoneofParallel(b *testing.B) {
	w := &benchNoneof{Color: "blue"}
	guardPlaygroundBenchmark(b, func() {
		_ = playgroundV.Struct(w)
	})
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = playgroundV.Struct(w)
		}
	})
}

// ══════════════════════════════════════════════════════════════════════════════
// Struct Tests — Model (SimpleUser / ComplexOrder / NestedWorkspace)
// ══════════════════════════════════════════════════════════════════════════════

func validSimpleUser() SimpleUser {
	return SimpleUser{Name: "Alice", Email: "alice@example.com", Age: 30}
}

func invalidSimpleUser() SimpleUser {
	return SimpleUser{Name: "", Email: "not-an-email", Age: -1}
}

func validComplexOrder() ComplexOrder {
	return ComplexOrder{
		OrderID:     "550e8400-e29b-41d4-a716-446655440000",
		CustomerID:  "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		Email:       "customer@example.com",
		Phone:       "+14155552671",
		Currency:    "USD",
		TotalAmount: 99.99,
		Status:      "pending",
		Items: []OrderItem{
			{ProductID: "6ba7b811-9dad-11d1-80b4-00c04fd430c8", SKU: "SKU001", Quantity: 2, UnitPrice: 49.99},
			{ProductID: "6ba7b812-9dad-11d1-80b4-00c04fd430c8", SKU: "SKU002", Quantity: 1, UnitPrice: 0.01},
		},
		Shipping: ShippingInfo{Country: "US", State: "California", City: "San Francisco", Street: "123 Main St", ZipCode: "94102"},
	}
}

func invalidComplexOrder() ComplexOrder {
	return ComplexOrder{
		OrderID: "not-a-uuid", CustomerID: "", Email: "bad-email", Phone: "123",
		Currency: "US", TotalAmount: -10, Status: "unknown", Items: []OrderItem{},
		Shipping: ShippingInfo{Country: "USA", State: "", City: "", Street: "", ZipCode: ""},
	}
}

func validNestedWorkspace() NestedWorkspace {
	return NestedWorkspace{
		WorkspaceID: "550e8400-e29b-41d4-a716-446655440000",
		Code:        "myworkspace",
		DisplayName: "My Workspace",
		AdminEmail:  "admin@example.com",
		Envs: []EnvironmentConfig{
			{EnvID: "6ba7b810-9dad-11d1-80b4-00c04fd430c8", DisplayName: "Production", Region: "US", SupportedLocales: []string{"zh", "en", "ja"}, Theme: "dark"},
			{EnvID: "6ba7b811-9dad-11d1-80b4-00c04fd430c8", DisplayName: "Staging", Region: "EU", SupportedLocales: []string{"en", "es"}, Theme: "light"},
		},
	}
}

func invalidNestedWorkspace() NestedWorkspace {
	return NestedWorkspace{
		WorkspaceID: "not-uuid", Code: "a", DisplayName: "", AdminEmail: "bad-email",
		Envs: []EnvironmentConfig{
			{EnvID: "", DisplayName: "", Region: "XX", SupportedLocales: []string{"invalid"}, Theme: "neon"},
		},
	}
}

// ─── SimpleUser Valid ───

func BenchmarkArgus_SimpleUser_Valid(b *testing.B) {
	data := validSimpleUser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		argusV.Struct(data)
	}
}

func BenchmarkPlayground_SimpleUser_Valid(b *testing.B) {
	data := validSimpleUser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		playgroundV.Struct(data)
	}
}

// ─── SimpleUser Invalid ───

func BenchmarkArgus_SimpleUser_Invalid(b *testing.B) {
	data := invalidSimpleUser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		argusV.Struct(data)
	}
}

func BenchmarkPlayground_SimpleUser_Invalid(b *testing.B) {
	data := invalidSimpleUser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		playgroundV.Struct(data)
	}
}

// ─── ComplexOrder Valid ───

func BenchmarkArgus_ComplexOrder_Valid(b *testing.B) {
	data := validComplexOrder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		argusV.Struct(data)
	}
}

func BenchmarkPlayground_ComplexOrder_Valid(b *testing.B) {
	data := validComplexOrder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		playgroundV.Struct(data)
	}
}

// ─── ComplexOrder Invalid ───

func BenchmarkArgus_ComplexOrder_Invalid(b *testing.B) {
	data := invalidComplexOrder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		argusV.Struct(data)
	}
}

func BenchmarkPlayground_ComplexOrder_Invalid(b *testing.B) {
	data := invalidComplexOrder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		playgroundV.Struct(data)
	}
}

// ─── NestedWorkspace Valid ───

func BenchmarkArgus_NestedWorkspace_Valid(b *testing.B) {
	data := validNestedWorkspace()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		argusV.Struct(data)
	}
}

func BenchmarkPlayground_NestedWorkspace_Valid(b *testing.B) {
	data := validNestedWorkspace()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		playgroundV.Struct(data)
	}
}

// ─── NestedWorkspace Invalid ───

func BenchmarkArgus_NestedWorkspace_Invalid(b *testing.B) {
	data := invalidNestedWorkspace()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		argusV.Struct(data)
	}
}

func BenchmarkPlayground_NestedWorkspace_Invalid(b *testing.B) {
	data := invalidNestedWorkspace()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		playgroundV.Struct(data)
	}
}

// ══════════════════════════════════════════════════════════════════════════════
// Struct Tests — RequiredIf / EqField (cross-field struct)
// ══════════════════════════════════════════════════════════════════════════════

// ─── RequiredIf ───

func BenchmarkArgus_RequiredIf(b *testing.B) {
	type S struct {
		Flag string `validate:"required,oneof=a b"`
		Val  string `validate:"required_if=Flag a"`
	}
	s := S{Flag: "a", Val: "x"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		argusV.Struct(s)
	}
}

func BenchmarkPlayground_RequiredIf(b *testing.B) {
	type S struct {
		Flag string `validate:"required,oneof=a b"`
		Val  string `validate:"required_if=Flag a"`
	}
	s := S{Flag: "a", Val: "x"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		playgroundV.Struct(s)
	}
}

// ─── EqField ───

func BenchmarkArgus_EqField(b *testing.B) {
	type S struct {
		Password string `validate:"required"`
		Confirm  string `validate:"required,eqfield=Password"`
	}
	s := S{Password: "secret", Confirm: "secret"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		argusV.Struct(s)
	}
}

func BenchmarkPlayground_EqField(b *testing.B) {
	type S struct {
		Password string `validate:"required"`
		Confirm  string `validate:"required,eqfield=Password"`
	}
	s := S{Password: "secret", Confirm: "secret"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		playgroundV.Struct(s)
	}
}
