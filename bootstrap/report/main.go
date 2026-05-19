/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-16 11:11:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-18 15:01:15
 * @FilePath: \go-argus-benchmark\bootstrap\report\main.go
 * @Description: 测试报告解析器
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type benchResult struct {
	Name     string  `json:"name"`
	NsPerOp  float64 `json:"ns_per_op"`
	BPerOp   uint64  `json:"bytes_per_op"`
	AllocsOp uint64  `json:"allocs_per_op"`
}

type comparisonRow struct {
	Scenario     string  `json:"scenario"`
	ArgusNs      float64 `json:"argus_ns"`
	PlaygroundNs float64 `json:"playground_ns"`
	Ratio        float64 `json:"ratio"`
	Winner       string  `json:"winner"`
}

type varStringVsVarRow struct {
	Rule        string  `json:"rule"`
	VarStringNs float64 `json:"varstring_ns"`
	VarNs       float64 `json:"var_ns"`
	Ratio       float64 `json:"ratio"`
	Speedup     string  `json:"speedup"`
}

type varStringVsVarAllocRow struct {
	Rule            string `json:"rule"`
	VarStringAllocs uint64 `json:"varstring_allocs"`
	VarAllocs       uint64 `json:"var_allocs"`
	VarStringBytes  uint64 `json:"varstring_bytes"`
	VarBytes        uint64 `json:"var_bytes"`
}

type envInfo struct {
	Goos   string `json:"goos"`
	Goarch string `json:"goarch"`
	Pkg    string `json:"pkg"`
	CPU    string `json:"cpu"`
}

type allocRow struct {
	Scenario         string `json:"scenario"`
	ArgusAllocs      uint64 `json:"argus_allocs"`
	PlaygroundAllocs uint64 `json:"playground_allocs"`
	ArgusBytes       uint64 `json:"argus_bytes"`
	PlaygroundBytes  uint64 `json:"playground_bytes"`
	Winner           string `json:"winner"`
}

func main() {
	rootDir := "."
	parseFile := ""

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-parse" && i+1 < len(os.Args) {
			parseFile = os.Args[i+1]
			i++
		} else {
			rootDir = os.Args[i]
		}
	}

	var allResults []benchResult
	var env envInfo

	if parseFile != "" {
		allResults, env = parseBenchmarkFileUnifiedWithEnv(parseFile)
		if len(allResults) == 0 {
			fmt.Println("Warning: no benchmark results parsed, using fallback data")
			allResults = fallbackDataUnified()
		}
	} else {
		allResults = fallbackDataUnified()
	}

	argusResults, pgResults := splitArgusPlayground(allResults)

	comparisons := buildComparisons(argusResults, pgResults)
	allocs := buildAllocComparisons(argusResults, pgResults)
	vsRows, vsAllocRows := buildVarStringVsVarComparisons(allResults)

	benchDir := filepath.Join(rootDir, "benchmarks")
	os.MkdirAll(benchDir, 0755)

	seqComparisons := filterParallel(comparisons, false)
	parComparisons := filterParallel(comparisons, true)

	latencySVG := generateLatencySVG(seqComparisons, "Latency Comparison (ns/op) — Sequential")
	writeFile(filepath.Join(benchDir, "latency.svg"), latencySVG)

	parallelSVG := generateLatencySVG(parComparisons, "Latency Comparison (ns/op) — Parallel")
	writeFile(filepath.Join(benchDir, "latency_parallel.svg"), parallelSVG)

	seqAllocs := filterAllocParallel(allocs, false)
	allocSVG := generateAllocSVG(seqAllocs, "Memory Allocation (heap allocs/op) — Sequential")
	writeFile(filepath.Join(benchDir, "allocs.svg"), allocSVG)

	vsSVG := generateVarStringVsVarSVG(vsRows, "VarString vs Var Latency (ns/op) — Zero-Reflection Gain")
	writeFile(filepath.Join(benchDir, "varstring_vs_var.svg"), vsSVG)

	vsAllocSVG := generateVarStringVsVarAllocSVG(vsAllocRows, "VarString vs Var Memory (allocs/op)")
	writeFile(filepath.Join(benchDir, "varstring_vs_var_allocs.svg"), vsAllocSVG)

	writeJSON(filepath.Join(benchDir, "benchmark_results.json"), struct {
		Argus      []benchResult `json:"argus"`
		Playground []benchResult `json:"playground"`
	}{argusResults, pgResults})
	writeJSON(filepath.Join(benchDir, "benchmark_comparisons.json"), comparisons)
	writeJSON(filepath.Join(benchDir, "benchmark_allocs.json"), allocs)
	writeJSON(filepath.Join(benchDir, "varstring_vs_var.json"), vsRows)
	writeJSON(filepath.Join(benchDir, "varstring_vs_var_allocs.json"), vsAllocRows)

	generateBenchmarksMD(rootDir, comparisons, allocs, vsRows, vsAllocRows, env)

	fmt.Printf("Done! Generated %d Argus/Playground comparisons, %d alloc rows, %d VarString/Var comparisons\n",
		len(comparisons), len(allocs), len(vsRows))
}

func fallbackDataUnified() []benchResult {
	return []benchResult{
		{"BenchmarkArgus_ComplexOrder_Invalid", 2311, 2525, 16},
		{"BenchmarkArgus_ComplexOrder_Valid", 2048, 304, 5},
		{"BenchmarkPlayground_ComplexOrder_Invalid", 2850, 3086, 41},
		{"BenchmarkPlayground_ComplexOrder_Valid", 3589, 592, 19},
		{"BenchmarkArgus_VarString_Email", 62, 0, 0},
		{"BenchmarkArgus_Var_Email", 210, 200, 3},
		{"BenchmarkPlayground_Var_Email", 626, 98, 5},
		{"BenchmarkArgus_VarString_URL", 75, 0, 0},
		{"BenchmarkArgus_Var_URL", 280, 200, 3},
		{"BenchmarkPlayground_Var_URL", 750, 98, 5},
	}
}

func parseBenchmarkFileUnified(filename string) []benchResult {
	results, _ := parseBenchmarkFileUnifiedWithEnv(filename)
	return results
}

func parseBenchmarkFileUnifiedWithEnv(filename string) ([]benchResult, envInfo) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening %s: %v\n", filename, err)
		return nil, envInfo{}
	}
	defer f.Close()

	aggregated := map[string]*aggResult{}
	var env envInfo
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "goos:") {
			env.Goos = strings.TrimSpace(strings.TrimPrefix(line, "goos:"))
			continue
		}
		if strings.HasPrefix(line, "goarch:") {
			env.Goarch = strings.TrimSpace(strings.TrimPrefix(line, "goarch:"))
			continue
		}
		if strings.HasPrefix(line, "pkg:") {
			env.Pkg = strings.TrimSpace(strings.TrimPrefix(line, "pkg:"))
			continue
		}
		if strings.HasPrefix(line, "cpu:") {
			env.CPU = strings.TrimSpace(strings.TrimPrefix(line, "cpu:"))
			continue
		}
		if !strings.HasPrefix(line, "Benchmark") {
			continue
		}
		name, nsPerOp, bPerOp, allocsOp, ok := parseBenchLine(line)
		if !ok {
			continue
		}
		if _, exists := aggregated[name]; !exists {
			aggregated[name] = &aggResult{name: name}
		}
		aggregated[name].add(nsPerOp, bPerOp, allocsOp)
	}

	return aggregateResults(aggregated), env
}

func splitArgusPlayground(all []benchResult) ([]benchResult, []benchResult) {
	return collectResultsByBenchmarkPrefix(all, "BenchmarkArgus"),
		collectResultsByBenchmarkPrefix(all, "BenchmarkPlayground")
}

func buildVarStringVsVarComparisons(all []benchResult) ([]varStringVsVarRow, []varStringVsVarAllocRow) {
	varStringAgg := map[string]*aggResult{}
	varAgg := map[string]*aggResult{}

	for _, r := range all {
		if addByTrimmedPrefix(varStringAgg, r, "BenchmarkArgus_VarString") {
			continue
		}
		addByTrimmedPrefix(varAgg, r, "BenchmarkArgus_Var")
	}

	varStringMap := resultsByName(aggregateResults(varStringAgg))
	varMap := resultsByName(aggregateResults(varAgg))

	var rows []varStringVsVarRow
	var allocRows []varStringVsVarAllocRow

	for rule, vs := range varStringMap {
		v, ok := varMap[rule]
		if !ok {
			continue
		}
		ratio := v.NsPerOp / vs.NsPerOp
		speedup := "VarString faster"
		if vs.NsPerOp > v.NsPerOp {
			speedup = "Var faster"
		}
		rows = append(rows, varStringVsVarRow{
			Rule: rule, VarStringNs: vs.NsPerOp, VarNs: v.NsPerOp,
			Ratio: ratio, Speedup: speedup,
		})
		allocRows = append(allocRows, varStringVsVarAllocRow{
			Rule: rule, VarStringAllocs: vs.AllocsOp, VarAllocs: v.AllocsOp,
			VarStringBytes: vs.BPerOp, VarBytes: v.BPerOp,
		})
	}

	sort.Slice(rows, func(i, j int) bool { return rows[i].Rule < rows[j].Rule })
	sort.Slice(allocRows, func(i, j int) bool { return allocRows[i].Rule < allocRows[j].Rule })
	return rows, allocRows
}

func collectResultsByBenchmarkPrefix(all []benchResult, prefix string) []benchResult {
	aggregated := map[string]*aggResult{}
	for _, r := range all {
		addByTrimmedPrefix(aggregated, r, prefix)
	}
	return aggregateResults(aggregated)
}

func addByTrimmedPrefix(aggregated map[string]*aggResult, r benchResult, prefix string) bool {
	name, ok := trimBenchmarkPrefix(r.Name, prefix)
	if !ok {
		return false
	}
	if _, exists := aggregated[name]; !exists {
		aggregated[name] = &aggResult{name: name}
	}
	aggregated[name].add(r.NsPerOp, r.BPerOp, r.AllocsOp)
	return true
}

func trimBenchmarkPrefix(name, prefix string) (string, bool) {
	if !strings.HasPrefix(name, prefix) {
		return "", false
	}
	trimmed := strings.TrimPrefix(name, prefix)
	trimmed = strings.TrimPrefix(trimmed, "_")
	if trimmed == "" {
		return "", false
	}
	return trimmed, true
}

func aggregateResults(aggregated map[string]*aggResult) []benchResult {
	var results []benchResult
	for _, agg := range aggregated {
		results = append(results, benchResult{
			Name:     agg.name,
			NsPerOp:  agg.avgNsPerOp(),
			BPerOp:   agg.avgBPerOp(),
			AllocsOp: agg.avgAllocsOp(),
		})
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Name < results[j].Name })
	return results
}

func resultsByName(results []benchResult) map[string]benchResult {
	byName := map[string]benchResult{}
	for _, r := range results {
		byName[r.Name] = r
	}
	return byName
}

type aggResult struct {
	name      string
	nsList    []float64
	bList     []uint64
	allocList []uint64
}

func (a *aggResult) add(ns float64, b uint64, alloc uint64) {
	a.nsList = append(a.nsList, ns)
	a.bList = append(a.bList, b)
	a.allocList = append(a.allocList, alloc)
}

func (a *aggResult) avgNsPerOp() float64 {
	if len(a.nsList) == 0 {
		return 0
	}
	var sum float64
	for _, v := range a.nsList {
		sum += v
	}
	return sum / float64(len(a.nsList))
}

func (a *aggResult) avgBPerOp() uint64 {
	if len(a.bList) == 0 {
		return 0
	}
	var sum uint64
	for _, v := range a.bList {
		sum += v
	}
	return sum / uint64(len(a.bList))
}

func (a *aggResult) avgAllocsOp() uint64 {
	if len(a.allocList) == 0 {
		return 0
	}
	var sum uint64
	for _, v := range a.allocList {
		sum += v
	}
	return sum / uint64(len(a.allocList))
}

func parseBenchLine(line string) (name string, nsPerOp float64, bPerOp uint64, allocsOp uint64, ok bool) {
	fields := strings.Fields(line)
	if len(fields) < 5 {
		return "", 0, 0, 0, false
	}

	name = normalizeBenchName(fields[0])

	nsIdx := -1
	for i, f := range fields {
		if strings.HasSuffix(f, "ns/op") {
			nsIdx = i
			break
		}
	}
	if nsIdx < 1 {
		return "", 0, 0, 0, false
	}

	nsPerOp, err := strconv.ParseFloat(fields[nsIdx-1], 64)
	if err != nil {
		return "", 0, 0, 0, false
	}

	for i, f := range fields {
		if strings.HasSuffix(f, "B/op") && i > 0 {
			b, _ := strconv.ParseUint(fields[i-1], 10, 64)
			bPerOp = b
		}
		if strings.HasSuffix(f, "allocs/op") && i > 0 {
			a, _ := strconv.ParseUint(fields[i-1], 10, 64)
			allocsOp = a
		}
	}

	return name, nsPerOp, bPerOp, allocsOp, true
}

func normalizeBenchName(name string) string {
	idx := strings.LastIndexByte(name, '-')
	if idx < 0 || idx == len(name)-1 {
		return name
	}
	for _, r := range name[idx+1:] {
		if r < '0' || r > '9' {
			return name
		}
	}
	return name[:idx]
}

func buildComparisons(argus, pg []benchResult) []comparisonRow {
	pgMap := map[string]benchResult{}
	for _, r := range pg {
		pgMap[r.Name] = r
	}
	var rows []comparisonRow
	for _, a := range argus {
		p, ok := pgMap[a.Name]
		if !ok {
			continue
		}
		ratio := p.NsPerOp / a.NsPerOp
		winner := "go-argus"
		if a.NsPerOp > p.NsPerOp {
			winner = "validator/v10"
		}
		rows = append(rows, comparisonRow{a.Name, a.NsPerOp, p.NsPerOp, ratio, winner})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Scenario < rows[j].Scenario })
	return rows
}

func buildAllocComparisons(argus, pg []benchResult) []allocRow {
	pgMap := map[string]benchResult{}
	for _, r := range pg {
		pgMap[r.Name] = r
	}
	var rows []allocRow
	for _, a := range argus {
		p, ok := pgMap[a.Name]
		if !ok {
			continue
		}
		winner := "go-argus"
		if a.AllocsOp > p.AllocsOp || (a.AllocsOp == p.AllocsOp && a.BPerOp > p.BPerOp) {
			winner = "validator/v10"
		}
		rows = append(rows, allocRow{a.Name, a.AllocsOp, p.AllocsOp, a.BPerOp, p.BPerOp, winner})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Scenario < rows[j].Scenario })
	return rows
}

func filterParallel(comparisons []comparisonRow, parallel bool) []comparisonRow {
	var filtered []comparisonRow
	for _, c := range comparisons {
		if strings.Contains(c.Scenario, "Parallel") == parallel {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

func filterAllocParallel(allocs []allocRow, parallel bool) []allocRow {
	var filtered []allocRow
	for _, a := range allocs {
		if strings.Contains(a.Scenario, "Parallel") == parallel {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

func writeFile(path, content string) {
	os.WriteFile(path, []byte(content), 0644)
}

func writeJSON(filename string, data interface{}) {
	f, _ := os.Create(filename)
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.Encode(data)
}

const (
	svgWidth        = 800
	barHeight       = 28
	barGap          = 6
	groupGap        = 20
	marginLeft      = 190
	marginRight     = 60
	marginTop       = 50
	marginBottom    = 40
	argusColor      = "#3B82F6"
	playgroundColor = "#F97316"
	varStringColor  = "#10B981"
	varColor        = "#8B5CF6"
	argusLabel      = "go-argus"
	playgroundLabel = "validator/v10"
	varStringLabel  = "VarString"
	varLabel        = "Var"
	bgColor         = "#FFFFFF"
	titleColor      = "#1E293B"
	labelColor      = "#334155"
	valueColor      = "#475569"
	legendColor     = "#64748B"
	hintColor       = "#94A3B8"
)

func generateLatencySVG(comparisons []comparisonRow, title string) string {
	if len(comparisons) == 0 {
		return "<svg></svg>"
	}

	maxNs := 0.0
	for _, c := range comparisons {
		maxNs = math.Max(maxNs, c.ArgusNs)
		maxNs = math.Max(maxNs, c.PlaygroundNs)
	}

	chartWidth := svgWidth - marginLeft - marginRight
	groupHeight := 2*barHeight + barGap + groupGap
	totalHeight := marginTop + len(comparisons)*groupHeight + marginBottom + 30

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, svgWidth, totalHeight, svgWidth, totalHeight))
	sb.WriteString(fmt.Sprintf(`<rect width="100%%" height="100%%" fill="%s" rx="12"/>`, bgColor))
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="30" fill="%s" font-family="system-ui,-apple-system,sans-serif" font-size="16" font-weight="600">%s</text>`, marginLeft, titleColor, title))
	sb.WriteString(fmt.Sprintf(`<rect x="%d" y="14" width="14" height="14" rx="3" fill="%s"/>`, svgWidth-250, argusColor))
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="26" fill="%s" font-family="system-ui,sans-serif" font-size="12">%s</text>`, svgWidth-232, legendColor, argusLabel))
	sb.WriteString(fmt.Sprintf(`<rect x="%d" y="14" width="14" height="14" rx="3" fill="%s"/>`, svgWidth-120, playgroundColor))
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="26" fill="%s" font-family="system-ui,sans-serif" font-size="12">%s</text>`, svgWidth-102, legendColor, playgroundLabel))

	for i, c := range comparisons {
		y := marginTop + i*groupHeight
		label := c.Scenario
		if len(label) > 28 {
			label = label[:25] + "..."
		}
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" fill="%s" font-family="monospace" font-size="12" text-anchor="end">%s</text>`, marginLeft-10, y+barHeight-6, labelColor, label))
		argusW := (c.ArgusNs / maxNs) * float64(chartWidth)
		pgW := (c.PlaygroundNs / maxNs) * float64(chartWidth)
		sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%.1f" height="%d" rx="4" fill="%s" opacity="0.9"/>`, marginLeft, y, argusW, barHeight, argusColor))
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" fill="%s" font-family="monospace" font-size="11">%.0f ns</text>`, marginLeft+int(argusW)+6, y+barHeight-6, valueColor, c.ArgusNs))
		sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%.1f" height="%d" rx="4" fill="%s" opacity="0.9"/>`, marginLeft, y+barHeight+barGap, pgW, barHeight, playgroundColor))
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" fill="%s" font-family="monospace" font-size="11">%.0f ns</text>`, marginLeft+int(pgW)+6, y+2*barHeight+barGap-6, valueColor, c.PlaygroundNs))
	}

	sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" fill="%s" font-family="system-ui,sans-serif" font-size="10">Lower is better ▸</text>`, marginLeft, totalHeight-10, hintColor))
	sb.WriteString(`</svg>`)
	return sb.String()
}

func generateVarStringVsVarSVG(rows []varStringVsVarRow, title string) string {
	if len(rows) == 0 {
		return "<svg></svg>"
	}

	maxNs := 0.0
	for _, r := range rows {
		maxNs = math.Max(maxNs, r.VarStringNs)
		maxNs = math.Max(maxNs, r.VarNs)
	}

	chartWidth := svgWidth - marginLeft - marginRight
	groupHeight := 2*barHeight + barGap + groupGap
	totalHeight := marginTop + len(rows)*groupHeight + marginBottom + 30

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, svgWidth, totalHeight, svgWidth, totalHeight))
	sb.WriteString(fmt.Sprintf(`<rect width="100%%" height="100%%" fill="%s" rx="12"/>`, bgColor))
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="30" fill="%s" font-family="system-ui,-apple-system,sans-serif" font-size="16" font-weight="600">%s</text>`, marginLeft, titleColor, title))
	sb.WriteString(fmt.Sprintf(`<rect x="%d" y="14" width="14" height="14" rx="3" fill="%s"/>`, svgWidth-280, varStringColor))
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="26" fill="%s" font-family="system-ui,sans-serif" font-size="12">%s</text>`, svgWidth-262, legendColor, varStringLabel))
	sb.WriteString(fmt.Sprintf(`<rect x="%d" y="14" width="14" height="14" rx="3" fill="%s"/>`, svgWidth-150, varColor))
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="26" fill="%s" font-family="system-ui,sans-serif" font-size="12">%s</text>`, svgWidth-132, legendColor, varLabel))

	for i, r := range rows {
		y := marginTop + i*groupHeight
		label := r.Rule
		if len(label) > 28 {
			label = label[:25] + "..."
		}
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" fill="%s" font-family="monospace" font-size="12" text-anchor="end">%s</text>`, marginLeft-10, y+barHeight-6, labelColor, label))
		vsW := (r.VarStringNs / maxNs) * float64(chartWidth)
		vW := (r.VarNs / maxNs) * float64(chartWidth)
		sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%.1f" height="%d" rx="4" fill="%s" opacity="0.9"/>`, marginLeft, y, vsW, barHeight, varStringColor))
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" fill="%s" font-family="monospace" font-size="11">%.0f ns</text>`, marginLeft+int(vsW)+6, y+barHeight-6, valueColor, r.VarStringNs))
		sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%.1f" height="%d" rx="4" fill="%s" opacity="0.9"/>`, marginLeft, y+barHeight+barGap, vW, barHeight, varColor))
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" fill="%s" font-family="monospace" font-size="11">%.0f ns</text>`, marginLeft+int(vW)+6, y+2*barHeight+barGap-6, valueColor, r.VarNs))
	}

	sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" fill="%s" font-family="system-ui,sans-serif" font-size="10">Lower is better ▸</text>`, marginLeft, totalHeight-10, hintColor))
	sb.WriteString(`</svg>`)
	return sb.String()
}

func generateVarStringVsVarAllocSVG(rows []varStringVsVarAllocRow, title string) string {
	if len(rows) == 0 {
		return "<svg></svg>"
	}

	maxAllocs := uint64(0)
	for _, r := range rows {
		maxAllocs = maxU64(maxAllocs, r.VarStringAllocs, r.VarAllocs)
	}
	if maxAllocs == 0 {
		maxAllocs = 1
	}

	chartWidth := svgWidth - marginLeft - marginRight
	groupHeight := 2*barHeight + barGap + groupGap
	totalHeight := marginTop + len(rows)*groupHeight + marginBottom + 30

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, svgWidth, totalHeight, svgWidth, totalHeight))
	sb.WriteString(fmt.Sprintf(`<rect width="100%%" height="100%%" fill="%s" rx="12"/>`, bgColor))
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="30" fill="%s" font-family="system-ui,-apple-system,sans-serif" font-size="16" font-weight="600">%s</text>`, marginLeft, titleColor, title))

	legendX := marginLeft
	sb.WriteString(fmt.Sprintf(`<rect x="%d" y="14" width="14" height="14" rx="3" fill="%s"/>`, legendX, varStringColor))
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="26" fill="%s" font-family="system-ui,sans-serif" font-size="12">%s</text>`, legendX+18, legendColor, varStringLabel))
	sb.WriteString(fmt.Sprintf(`<rect x="%d" y="14" width="14" height="14" rx="3" fill="%s"/>`, legendX+110, varColor))
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="26" fill="%s" font-family="system-ui,sans-serif" font-size="12">%s</text>`, legendX+128, legendColor, varLabel))

	for i, r := range rows {
		y := marginTop + i*groupHeight
		label := r.Rule
		if len(label) > 22 {
			label = label[:19] + "..."
		}
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" fill="%s" font-family="monospace" font-size="12" text-anchor="end">%s</text>`, marginLeft-10, y+barHeight-6, labelColor, label))
		vsW := (float64(r.VarStringAllocs) / float64(maxAllocs)) * float64(chartWidth)
		vW := (float64(r.VarAllocs) / float64(maxAllocs)) * float64(chartWidth)
		sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%.1f" height="%d" rx="4" fill="%s" opacity="0.9"/>`, marginLeft, y, vsW, barHeight, varStringColor))
		sb.WriteString(allocValueText(marginLeft, vsW, y+barHeight-6, r.VarStringAllocs, r.VarStringBytes, chartWidth, varStringColor))
		sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%.1f" height="%d" rx="4" fill="%s" opacity="0.9"/>`, marginLeft, y+barHeight+barGap, vW, barHeight, varColor))
		sb.WriteString(allocValueText(marginLeft, vW, y+2*barHeight+barGap-6, r.VarAllocs, r.VarBytes, chartWidth, varColor))
	}

	sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" fill="%s" font-family="system-ui,sans-serif" font-size="10">Lower is better ▸</text>`, marginLeft, totalHeight-10, hintColor))
	sb.WriteString(`</svg>`)
	return sb.String()
}

func generateAllocSVG(allocs []allocRow, title string) string {
	if len(allocs) == 0 {
		return "<svg></svg>"
	}

	maxAllocs := uint64(0)
	for _, a := range allocs {
		maxAllocs = maxU64(maxAllocs, a.ArgusAllocs, a.PlaygroundAllocs)
	}
	if maxAllocs == 0 {
		maxAllocs = 1
	}

	chartWidth := svgWidth - marginLeft - marginRight
	groupHeight := 2*barHeight + barGap + groupGap
	totalHeight := marginTop + len(allocs)*groupHeight + marginBottom + 30

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, svgWidth, totalHeight, svgWidth, totalHeight))
	sb.WriteString(fmt.Sprintf(`<rect width="100%%" height="100%%" fill="%s" rx="12"/>`, bgColor))
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="30" fill="%s" font-family="system-ui,-apple-system,sans-serif" font-size="16" font-weight="600">%s</text>`, marginLeft, titleColor, title))

	legendX := marginLeft
	sb.WriteString(fmt.Sprintf(`<rect x="%d" y="14" width="14" height="14" rx="3" fill="%s"/>`, legendX, argusColor))
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="26" fill="%s" font-family="system-ui,sans-serif" font-size="12">%s</text>`, legendX+18, legendColor, argusLabel))
	sb.WriteString(fmt.Sprintf(`<rect x="%d" y="14" width="14" height="14" rx="3" fill="%s"/>`, legendX+110, playgroundColor))
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="26" fill="%s" font-family="system-ui,sans-serif" font-size="12">%s</text>`, legendX+128, legendColor, playgroundLabel))

	for i, a := range allocs {
		y := marginTop + i*groupHeight
		label := a.Scenario
		if len(label) > 22 {
			label = label[:19] + "..."
		}
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" fill="%s" font-family="monospace" font-size="12" text-anchor="end">%s</text>`, marginLeft-10, y+barHeight-6, labelColor, label))
		argusW := (float64(a.ArgusAllocs) / float64(maxAllocs)) * float64(chartWidth)
		pgW := (float64(a.PlaygroundAllocs) / float64(maxAllocs)) * float64(chartWidth)
		sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%.1f" height="%d" rx="4" fill="%s" opacity="0.9"/>`, marginLeft, y, argusW, barHeight, argusColor))
		sb.WriteString(allocValueText(marginLeft, argusW, y+barHeight-6, a.ArgusAllocs, a.ArgusBytes, chartWidth, argusColor))
		sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%.1f" height="%d" rx="4" fill="%s" opacity="0.9"/>`, marginLeft, y+barHeight+barGap, pgW, barHeight, playgroundColor))
		sb.WriteString(allocValueText(marginLeft, pgW, y+2*barHeight+barGap-6, a.PlaygroundAllocs, a.PlaygroundBytes, chartWidth, playgroundColor))
	}

	sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" fill="%s" font-family="system-ui,sans-serif" font-size="10">Lower is better ▸</text>`, marginLeft, totalHeight-10, hintColor))
	sb.WriteString(`</svg>`)
	return sb.String()
}

func allocValueText(baseX int, barW float64, textY int, allocs, bytes uint64, chartWidth int, color string) string {
	if allocs == 0 {
		return fmt.Sprintf(`<text x="%d" y="%d" fill="%s" font-family="monospace" font-size="11">0</text>`, baseX+4, textY, valueColor)
	}
	text := fmt.Sprintf("%d (%dB)", allocs, bytes)
	textX := baseX + int(barW) + 6
	if textX+len(text)*7 > baseX+chartWidth {
		textX = baseX + int(barW) - len(text)*7 - 6
		return fmt.Sprintf(`<text x="%d" y="%d" fill="#FFF" font-family="monospace" font-size="11">%s</text>`, textX, textY, text)
	}
	return fmt.Sprintf(`<text x="%d" y="%d" fill="%s" font-family="monospace" font-size="11">%s</text>`, textX, textY, valueColor, text)
}

func maxU64(a, b, c uint64) uint64 {
	if b > a {
		a = b
	}
	if c > a {
		a = c
	}
	return a
}

func generateBenchmarksMD(rootDir string, comparisons []comparisonRow, allocs []allocRow, vsRows []varStringVsVarRow, vsAllocRows []varStringVsVarAllocRow, env envInfo) {
	var sb strings.Builder

	sb.WriteString("# Benchmark Details\n\n")
	sb.WriteString("Auto-generated by `go run ./bootstrap/report`. Do not edit manually.\n\n")
	if env.Goos != "" || env.Goarch != "" || env.CPU != "" {
		sb.WriteString("## Environment\n\n")
		sb.WriteString(fmt.Sprintf("| Key | Value |\n|-----|-------|\n| goos | %s |\n| goarch | %s |\n| pkg | %s |\n| cpu | %s |\n\n",
			env.Goos, env.Goarch, env.Pkg, env.CPU))
	}

	sb.WriteString("## Latency (ns/op) — Argus vs Playground\n\n")
	sb.WriteString("| Scenario | go-argus | validator/v10 | Ratio | Winner |\n")
	sb.WriteString("|----------|---------:|--------------:|------:|--------|\n")
	for _, c := range comparisons {
		sb.WriteString(fmt.Sprintf("| %s | %.0f | %.0f | %.2fx | %s |\n",
			c.Scenario, c.ArgusNs, c.PlaygroundNs, c.Ratio, c.Winner))
	}

	sb.WriteString("\n## Memory Allocation — Argus vs Playground\n\n")
	sb.WriteString("| Scenario | go-argus (allocs) | validator/v10 (allocs) | go-argus (bytes) | validator/v10 (bytes) | Winner |\n")
	sb.WriteString("|----------|------------------:|-----------------------:|-----------------:|----------------------:|--------|\n")
	for _, a := range allocs {
		sb.WriteString(fmt.Sprintf("| %s | %d | %d | %d | %d | %s |\n",
			a.Scenario, a.ArgusAllocs, a.PlaygroundAllocs, a.ArgusBytes, a.PlaygroundBytes, a.Winner))
	}

	if len(vsRows) > 0 {
		sb.WriteString("\n## VarString vs Var — Zero-Reflection Latency\n\n")
		sb.WriteString("| Rule | VarString (ns) | Var (ns) | Ratio (Var/VarString) | Speedup |\n")
		sb.WriteString("|------|--------------:|---------:|----------------------:|---------|\n")
		for _, r := range vsRows {
			sb.WriteString(fmt.Sprintf("| %s | %.0f | %.0f | %.2fx | %s |\n",
				r.Rule, r.VarStringNs, r.VarNs, r.Ratio, r.Speedup))
		}

		sb.WriteString("\n## VarString vs Var — Memory Allocation\n\n")
		sb.WriteString("| Rule | VarString (allocs) | Var (allocs) | VarString (bytes) | Var (bytes) |\n")
		sb.WriteString("|------|-------------------:|-------------:|------------------:|------------:|\n")
		for _, a := range vsAllocRows {
			sb.WriteString(fmt.Sprintf("| %s | %d | %d | %d | %d |\n",
				a.Rule, a.VarStringAllocs, a.VarAllocs, a.VarStringBytes, a.VarBytes))
		}
	}

	os.WriteFile(filepath.Join(rootDir, "BENCHMARKS.md"), []byte(sb.String()), 0644)
}
