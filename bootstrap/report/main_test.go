/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-16 11:11:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-18 11:58:15
 * @FilePath: \go-argus-benchmark\bootstrap\report\main_test.go
 * @Description: 测试报告解析器
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseBenchLineNormalizesGoBenchmarkName(t *testing.T) {
	line := "BenchmarkArgus_FieldSuccess-20  20865500  61.20 ns/op  16 B/op  1 allocs/op"

	name, nsPerOp, bPerOp, allocsOp, ok := parseBenchLine(line)

	if !ok {
		t.Fatal("expected benchmark line to parse")
	}
	if name != "BenchmarkArgus_FieldSuccess" {
		t.Fatalf("name = %q, want %q", name, "BenchmarkArgus_FieldSuccess")
	}
	if nsPerOp != 61.20 {
		t.Fatalf("nsPerOp = %v, want 61.20", nsPerOp)
	}
	if bPerOp != 16 {
		t.Fatalf("bPerOp = %d, want 16", bPerOp)
	}
	if allocsOp != 1 {
		t.Fatalf("allocsOp = %d, want 1", allocsOp)
	}
}

func TestParseBenchmarkFileUnifiedAggregatesNormalizedNames(t *testing.T) {
	content := `goos: windows
BenchmarkArgus_FieldSuccess-20  100  10 ns/op  0 B/op  0 allocs/op
BenchmarkArgus_FieldSuccess-16  100  20 ns/op  2 B/op  2 allocs/op
panic: later benchmark failed
`
	filename := filepath.Join(t.TempDir(), "benchmark_output.txt")
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	results := parseBenchmarkFileUnified(filename)

	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1: %#v", len(results), results)
	}
	got := results[0]
	if got.Name != "BenchmarkArgus_FieldSuccess" {
		t.Fatalf("Name = %q, want %q", got.Name, "BenchmarkArgus_FieldSuccess")
	}
	if got.NsPerOp != 15 {
		t.Fatalf("NsPerOp = %v, want 15", got.NsPerOp)
	}
	if got.BPerOp != 1 {
		t.Fatalf("BPerOp = %d, want 1", got.BPerOp)
	}
	if got.AllocsOp != 1 {
		t.Fatalf("AllocsOp = %d, want 1", got.AllocsOp)
	}
}

func TestSplitArgusPlaygroundTrimsPrefixesAndAveragesMatchingSuffixes(t *testing.T) {
	all := []benchResult{
		{Name: "BenchmarkArgus_FieldSuccess", NsPerOp: 10, BPerOp: 0, AllocsOp: 0},
		{Name: "BenchmarkArgus_FieldSuccess", NsPerOp: 20, BPerOp: 2, AllocsOp: 2},
		{Name: "BenchmarkPlayground_FieldSuccess", NsPerOp: 40, BPerOp: 4, AllocsOp: 4},
		{Name: "BenchmarkPlayground_FieldSuccess", NsPerOp: 60, BPerOp: 6, AllocsOp: 6},
	}

	argus, playground := splitArgusPlayground(all)

	if len(argus) != 1 {
		t.Fatalf("len(argus) = %d, want 1: %#v", len(argus), argus)
	}
	if argus[0].Name != "FieldSuccess" || argus[0].NsPerOp != 15 || argus[0].BPerOp != 1 || argus[0].AllocsOp != 1 {
		t.Fatalf("argus[0] = %#v", argus[0])
	}
	if len(playground) != 1 {
		t.Fatalf("len(playground) = %d, want 1: %#v", len(playground), playground)
	}
	if playground[0].Name != "FieldSuccess" || playground[0].NsPerOp != 50 || playground[0].BPerOp != 5 || playground[0].AllocsOp != 5 {
		t.Fatalf("playground[0] = %#v", playground[0])
	}
}

func TestBuildVarStringVsVarComparisonsTrimsPrefixesAndAveragesMatchingSuffixes(t *testing.T) {
	all := []benchResult{
		{Name: "BenchmarkArgus_VarStringPath", NsPerOp: 20, BPerOp: 0, AllocsOp: 0},
		{Name: "BenchmarkArgus_VarStringPath", NsPerOp: 30, BPerOp: 2, AllocsOp: 2},
		{Name: "BenchmarkArgus_VarPath", NsPerOp: 40, BPerOp: 2, AllocsOp: 1},
		{Name: "BenchmarkArgus_VarPath", NsPerOp: 60, BPerOp: 4, AllocsOp: 3},
		{Name: "BenchmarkArgus_VarString_Email", NsPerOp: 50, BPerOp: 0, AllocsOp: 0},
		{Name: "BenchmarkArgus_Var_Email", NsPerOp: 100, BPerOp: 4, AllocsOp: 2},
	}

	rows, allocRows := buildVarStringVsVarComparisons(all)

	if len(rows) != 2 {
		t.Fatalf("len(rows) = %d, want 2: %#v", len(rows), rows)
	}
	if rows[0].Rule != "Email" || rows[1].Rule != "Path" {
		t.Fatalf("rules = %q, %q; want Email, Path", rows[0].Rule, rows[1].Rule)
	}
	if rows[1].Ratio != 2 {
		t.Fatalf("Path ratio = %v, want 2", rows[1].Ratio)
	}
	if len(allocRows) != 2 {
		t.Fatalf("len(allocRows) = %d, want 2", len(allocRows))
	}
	if allocRows[1].Rule != "Path" || allocRows[1].VarStringAllocs != 1 || allocRows[1].VarAllocs != 2 ||
		allocRows[1].VarStringBytes != 1 || allocRows[1].VarBytes != 3 {
		t.Fatalf("Path alloc row = %#v", allocRows[1])
	}
}
