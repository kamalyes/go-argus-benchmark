/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-16 11:11:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-18 15:58:30
 * @FilePath: \go-argus-benchmark\bench_rules_test.go
 * @Description: Argus vs Playground 字符串规则性能基准测试
 *
 * 每个规则按 VarString → Var → Playground 三连排列，便于对比：
 *   - VarString: Argus 零反射快速路径
 *   - Var:        Argus 反射路径
 *   - Playground: validator/v10 反射路径
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package bench

import "testing"

// ══════════════════════════════════════════════════════════════════════════════
// required
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_Required(b *testing.B) { benchVarStringArgus(b, "hello", "required") }
func BenchmarkArgus_Var_Required(b *testing.B)       { benchVarArgus(b, "hello", "required") }
func BenchmarkPlayground_Var_Required(b *testing.B)  { benchVarPlayground(b, "hello", "required") }

// ══════════════════════════════════════════════════════════════════════════════
// min / max / len
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_Min(b *testing.B) { benchVarStringArgus(b, "abcdef", "min=3") }
func BenchmarkArgus_Var_Min(b *testing.B)       { benchVarArgus(b, "abcdef", "min=3") }
func BenchmarkPlayground_Var_Min(b *testing.B)  { benchVarPlayground(b, "abcdef", "min=3") }

func BenchmarkArgus_VarString_Max(b *testing.B) { benchVarStringArgus(b, "ab", "max=5") }
func BenchmarkArgus_Var_Max(b *testing.B)       { benchVarArgus(b, "ab", "max=5") }
func BenchmarkPlayground_Var_Max(b *testing.B)  { benchVarPlayground(b, "ab", "max=5") }

func BenchmarkArgus_VarString_Len(b *testing.B) { benchVarStringArgus(b, "abcde", "len=5") }
func BenchmarkArgus_Var_Len(b *testing.B)       { benchVarArgus(b, "abcde", "len=5") }
func BenchmarkPlayground_Var_Len(b *testing.B)  { benchVarPlayground(b, "abcde", "len=5") }

// ══════════════════════════════════════════════════════════════════════════════
// eq / ne
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_Eq(b *testing.B) { benchVarStringArgus(b, "hello", "eq=hello") }
func BenchmarkArgus_Var_Eq(b *testing.B)       { benchVarArgus(b, "hello", "eq=hello") }
func BenchmarkPlayground_Var_Eq(b *testing.B)  { benchVarPlayground(b, "hello", "eq=hello") }

func BenchmarkArgus_VarString_Ne(b *testing.B) { benchVarStringArgus(b, "hello", "ne=world") }
func BenchmarkArgus_Var_Ne(b *testing.B)       { benchVarArgus(b, "hello", "ne=world") }
func BenchmarkPlayground_Var_Ne(b *testing.B)  { benchVarPlayground(b, "hello", "ne=world") }

// ══════════════════════════════════════════════════════════════════════════════
// gt / gte / lt / lte
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_Gt(b *testing.B) { benchVarStringArgus(b, "10", "gt=5") }
func BenchmarkArgus_Var_Gt(b *testing.B)       { benchVarArgus(b, "10", "gt=5") }
func BenchmarkPlayground_Var_Gt(b *testing.B)  { benchVarPlayground(b, "10", "gt=5") }

func BenchmarkArgus_VarString_Gte(b *testing.B) { benchVarStringArgus(b, "10", "gte=10") }
func BenchmarkArgus_Var_Gte(b *testing.B)       { benchVarArgus(b, "10", "gte=10") }
func BenchmarkPlayground_Var_Gte(b *testing.B)  { benchVarPlayground(b, "10", "gte=10") }

func BenchmarkArgus_VarString_Lt(b *testing.B) { benchVarStringArgus(b, "3", "lt=5") }
func BenchmarkArgus_Var_Lt(b *testing.B)       { benchVarArgus(b, "3", "lt=5") }
func BenchmarkPlayground_Var_Lt(b *testing.B)  { benchVarPlayground(b, "3", "lt=5") }

func BenchmarkArgus_VarString_Lte(b *testing.B) { benchVarStringArgus(b, "5", "lte=5") }
func BenchmarkArgus_Var_Lte(b *testing.B)       { benchVarArgus(b, "5", "lte=5") }
func BenchmarkPlayground_Var_Lte(b *testing.B)  { benchVarPlayground(b, "5", "lte=5") }

// ══════════════════════════════════════════════════════════════════════════════
// oneof / oneofci / noneof / noneofci
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_Oneof(b *testing.B) {
	benchVarStringArgus(b, "red", "oneof=red green blue")
}
func BenchmarkArgus_Var_Oneof(b *testing.B) { benchVarArgus(b, "red", "oneof=red green blue") }
func BenchmarkPlayground_Var_Oneof(b *testing.B) {
	benchVarPlayground(b, "red", "oneof=red green blue")
}

func BenchmarkArgus_VarString_Oneofci(b *testing.B) {
	benchVarStringArgus(b, "Red", "oneofci=red green blue")
}
func BenchmarkArgus_Var_Oneofci(b *testing.B) { benchVarArgus(b, "Red", "oneofci=red green blue") }
func BenchmarkPlayground_Var_Oneofci(b *testing.B) {
	benchVarPlayground(b, "Red", "oneofci=red green blue")
}

func BenchmarkArgus_VarString_Noneof(b *testing.B) {
	benchVarStringArgus(b, "yellow", "noneof=red green blue")
}
func BenchmarkArgus_Var_Noneof(b *testing.B) { benchVarArgus(b, "yellow", "noneof=red green blue") }

func BenchmarkArgus_VarString_Noneofci(b *testing.B) {
	benchVarStringArgus(b, "Yellow", "noneofci=red green blue")
}
func BenchmarkArgus_Var_Noneofci(b *testing.B) { benchVarArgus(b, "Yellow", "noneofci=red green blue") }

// ══════════════════════════════════════════════════════════════════════════════
// email
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_Email(b *testing.B) {
	benchVarStringArgus(b, "test@example.com", "required,email")
}
func BenchmarkArgus_Var_Email(b *testing.B) { benchVarArgus(b, "test@example.com", "required,email") }
func BenchmarkPlayground_Var_Email(b *testing.B) {
	benchVarPlayground(b, "test@example.com", "required,email")
}

// ══════════════════════════════════════════════════════════════════════════════
// url / uri / http_url / https_url
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_URL(b *testing.B) {
	benchVarStringArgus(b, "https://example.com/path", "required,url")
}
func BenchmarkArgus_Var_URL(b *testing.B) {
	benchVarArgus(b, "https://example.com/path", "required,url")
}
func BenchmarkPlayground_Var_URL(b *testing.B) {
	benchVarPlayground(b, "https://example.com/path", "required,url")
}

func BenchmarkArgus_VarString_URI(b *testing.B) {
	benchVarStringArgus(b, "mailto:test@example.com", "uri")
}
func BenchmarkArgus_Var_URI(b *testing.B) { benchVarArgus(b, "mailto:test@example.com", "uri") }
func BenchmarkPlayground_Var_URI(b *testing.B) {
	benchVarPlayground(b, "mailto:test@example.com", "uri")
}

func BenchmarkArgus_VarString_HTTPURL(b *testing.B) {
	benchVarStringArgus(b, "http://example.com/path", "http_url")
}
func BenchmarkArgus_Var_HTTPURL(b *testing.B) {
	benchVarArgus(b, "http://example.com/path", "http_url")
}

func BenchmarkArgus_VarString_HTTPSURL(b *testing.B) {
	benchVarStringArgus(b, "https://example.com/path", "https_url")
}
func BenchmarkArgus_Var_HTTPSURL(b *testing.B) {
	benchVarArgus(b, "https://example.com/path", "https_url")
}

// ══════════════════════════════════════════════════════════════════════════════
// semver
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_Semver(b *testing.B) {
	benchVarStringArgus(b, "1.2.3-alpha.1+build.123", "semver")
}
func BenchmarkArgus_Var_Semver(b *testing.B) { benchVarArgus(b, "1.2.3-alpha.1+build.123", "semver") }

// ══════════════════════════════════════════════════════════════════════════════
// ip / ipv4 / ipv6 / cidr / cidrv4 / cidrv6
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_IP(b *testing.B) { benchVarStringArgus(b, "192.168.1.1", "ip") }
func BenchmarkArgus_Var_IP(b *testing.B)       { benchVarArgus(b, "192.168.1.1", "ip") }
func BenchmarkPlayground_Var_IP(b *testing.B)  { benchVarPlayground(b, "192.168.1.1", "ip") }

func BenchmarkArgus_VarString_IPv4(b *testing.B) { benchVarStringArgus(b, "192.168.1.1", "ipv4") }
func BenchmarkArgus_Var_IPv4(b *testing.B)       { benchVarArgus(b, "192.168.1.1", "ipv4") }
func BenchmarkPlayground_Var_IPv4(b *testing.B)  { benchVarPlayground(b, "192.168.1.1", "ipv4") }

func BenchmarkArgus_VarString_IPv6(b *testing.B) { benchVarStringArgus(b, "::1", "ipv6") }
func BenchmarkArgus_Var_IPv6(b *testing.B)       { benchVarArgus(b, "::1", "ipv6") }
func BenchmarkPlayground_Var_IPv6(b *testing.B)  { benchVarPlayground(b, "::1", "ipv6") }

func BenchmarkArgus_VarString_CIDR(b *testing.B) { benchVarStringArgus(b, "192.168.1.0/24", "cidr") }
func BenchmarkArgus_Var_CIDR(b *testing.B)       { benchVarArgus(b, "192.168.1.0/24", "cidr") }
func BenchmarkPlayground_Var_CIDR(b *testing.B)  { benchVarPlayground(b, "192.168.1.0/24", "cidr") }

func BenchmarkArgus_VarString_CIDRv4(b *testing.B) {
	benchVarStringArgus(b, "192.168.1.0/24", "cidrv4")
}
func BenchmarkArgus_Var_CIDRv4(b *testing.B) { benchVarArgus(b, "192.168.1.0/24", "cidrv4") }

func BenchmarkArgus_VarString_CIDRv6(b *testing.B) { benchVarStringArgus(b, "2001:db8::/32", "cidrv6") }
func BenchmarkArgus_Var_CIDRv6(b *testing.B)       { benchVarArgus(b, "2001:db8::/32", "cidrv6") }

// ══════════════════════════════════════════════════════════════════════════════
// uuid / uuid3 / uuid4 / uuid5
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_UUID(b *testing.B) {
	benchVarStringArgus(b, "6ba7b810-9dad-11d1-80b4-00c04fd430c8", "uuid")
}
func BenchmarkArgus_Var_UUID(b *testing.B) {
	benchVarArgus(b, "6ba7b810-9dad-11d1-80b4-00c04fd430c8", "uuid")
}
func BenchmarkPlayground_Var_UUID(b *testing.B) {
	benchVarPlayground(b, "6ba7b810-9dad-11d1-80b4-00c04fd430c8", "uuid")
}

func BenchmarkArgus_VarString_UUID3(b *testing.B) {
	benchVarStringArgus(b, "6ba7b810-9dad-11d1-80b4-00c04fd430c8", "uuid3")
}
func BenchmarkArgus_Var_UUID3(b *testing.B) {
	benchVarArgus(b, "6ba7b810-9dad-11d1-80b4-00c04fd430c8", "uuid3")
}
func BenchmarkPlayground_Var_UUID3(b *testing.B) {
	benchVarPlayground(b, "6ba7b810-9dad-11d1-80b4-00c04fd430c8", "uuid3")
}

func BenchmarkArgus_VarString_UUID4(b *testing.B) {
	benchVarStringArgus(b, "6ba7b810-9dad-41d1-80b4-00c04fd430c8", "uuid4")
}
func BenchmarkArgus_Var_UUID4(b *testing.B) {
	benchVarArgus(b, "6ba7b810-9dad-41d1-80b4-00c04fd430c8", "uuid4")
}
func BenchmarkPlayground_Var_UUID4(b *testing.B) {
	benchVarPlayground(b, "6ba7b810-9dad-41d1-80b4-00c04fd430c8", "uuid4")
}

func BenchmarkArgus_VarString_UUID5(b *testing.B) {
	benchVarStringArgus(b, "6ba7b810-9dad-51d1-80b4-00c04fd430c8", "uuid5")
}
func BenchmarkArgus_Var_UUID5(b *testing.B) {
	benchVarArgus(b, "6ba7b810-9dad-51d1-80b4-00c04fd430c8", "uuid5")
}
func BenchmarkPlayground_Var_UUID5(b *testing.B) {
	benchVarPlayground(b, "6ba7b810-9dad-51d1-80b4-00c04fd430c8", "uuid5")
}

// ══════════════════════════════════════════════════════════════════════════════
// hostname / fqdn / hostname_port
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_Hostname(b *testing.B) {
	benchVarStringArgus(b, "example.com", "hostname")
}
func BenchmarkArgus_Var_Hostname(b *testing.B)      { benchVarArgus(b, "example.com", "hostname") }
func BenchmarkPlayground_Var_Hostname(b *testing.B) { benchVarPlayground(b, "example.com", "hostname") }

func BenchmarkArgus_VarString_FQDN(b *testing.B) { benchVarStringArgus(b, "example.com.", "fqdn") }
func BenchmarkArgus_Var_FQDN(b *testing.B)       { benchVarArgus(b, "example.com.", "fqdn") }
func BenchmarkPlayground_Var_FQDN(b *testing.B)  { benchVarPlayground(b, "example.com.", "fqdn") }

func BenchmarkArgus_VarString_HostnamePort(b *testing.B) {
	benchVarStringArgus(b, "example.com:8080", "hostname_port")
}
func BenchmarkArgus_Var_HostnamePort(b *testing.B) {
	benchVarArgus(b, "example.com:8080", "hostname_port")
}
func BenchmarkPlayground_Var_HostnamePort(b *testing.B) {
	benchVarPlayground(b, "example.com:8080", "hostname_port")
}

// ══════════════════════════════════════════════════════════════════════════════
// alpha / alphaspace / alphanum / alphanumspace / alphaunicode / alphanumunicode
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_Alpha(b *testing.B) { benchVarStringArgus(b, "helloworld", "alpha") }
func BenchmarkArgus_Var_Alpha(b *testing.B)       { benchVarArgus(b, "helloworld", "alpha") }
func BenchmarkPlayground_Var_Alpha(b *testing.B)  { benchVarPlayground(b, "helloworld", "alpha") }

func BenchmarkArgus_VarString_AlphaSpace(b *testing.B) {
	benchVarStringArgus(b, "hello world", "alphaspace")
}
func BenchmarkArgus_Var_AlphaSpace(b *testing.B) { benchVarArgus(b, "hello world", "alphaspace") }

func BenchmarkArgus_VarString_Alphanum(b *testing.B) { benchVarStringArgus(b, "hello123", "alphanum") }
func BenchmarkArgus_Var_Alphanum(b *testing.B)       { benchVarArgus(b, "hello123", "alphanum") }
func BenchmarkPlayground_Var_Alphanum(b *testing.B)  { benchVarPlayground(b, "hello123", "alphanum") }

func BenchmarkArgus_VarString_AlphanumSpace(b *testing.B) {
	benchVarStringArgus(b, "hello 123", "alphanumspace")
}
func BenchmarkArgus_Var_AlphanumSpace(b *testing.B) { benchVarArgus(b, "hello 123", "alphanumspace") }

func BenchmarkArgus_VarString_AlphaUnicode(b *testing.B) {
	benchVarStringArgus(b, "héllo", "alphaunicode")
}
func BenchmarkArgus_Var_AlphaUnicode(b *testing.B) { benchVarArgus(b, "héllo", "alphaunicode") }

func BenchmarkArgus_VarString_AlphanumUnicode(b *testing.B) {
	benchVarStringArgus(b, "héllo123", "alphanumunicode")
}
func BenchmarkArgus_Var_AlphanumUnicode(b *testing.B) {
	benchVarArgus(b, "héllo123", "alphanumunicode")
}

// ══════════════════════════════════════════════════════════════════════════════
// ascii / printascii / multibyte
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_ASCII(b *testing.B) { benchVarStringArgus(b, "hello world 123", "ascii") }
func BenchmarkArgus_Var_ASCII(b *testing.B)       { benchVarArgus(b, "hello world 123", "ascii") }
func BenchmarkPlayground_Var_ASCII(b *testing.B)  { benchVarPlayground(b, "hello world 123", "ascii") }

func BenchmarkArgus_VarString_PrintASCII(b *testing.B) {
	benchVarStringArgus(b, "hello world 123", "printascii")
}
func BenchmarkArgus_Var_PrintASCII(b *testing.B) { benchVarArgus(b, "hello world 123", "printascii") }

func BenchmarkArgus_VarString_Multibyte(b *testing.B) {
	benchVarStringArgus(b, "你好世界", "multibyte")
}
func BenchmarkArgus_Var_Multibyte(b *testing.B) { benchVarArgus(b, "你好世界", "multibyte") }
func BenchmarkPlayground_Var_Multibyte(b *testing.B) {
	benchVarPlayground(b, "你好世界", "multibyte")
}

// ══════════════════════════════════════════════════════════════════════════════
// hexadecimal / hexcolor / rgb / rgba / hsl / hsla
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_Hexadecimal(b *testing.B) {
	benchVarStringArgus(b, "deadbeef", "hexadecimal")
}
func BenchmarkArgus_Var_Hexadecimal(b *testing.B) { benchVarArgus(b, "deadbeef", "hexadecimal") }
func BenchmarkPlayground_Var_Hexadecimal(b *testing.B) {
	benchVarPlayground(b, "deadbeef", "hexadecimal")
}

func BenchmarkArgus_VarString_HexColor(b *testing.B) { benchVarStringArgus(b, "#ff0033", "hexcolor") }
func BenchmarkArgus_Var_HexColor(b *testing.B)       { benchVarArgus(b, "#ff0033", "hexcolor") }
func BenchmarkPlayground_Var_HexColor(b *testing.B)  { benchVarPlayground(b, "#ff0033", "hexcolor") }

func BenchmarkArgus_VarString_RGB(b *testing.B) { benchVarStringArgus(b, "rgb(255,0,51)", "rgb") }
func BenchmarkArgus_Var_RGB(b *testing.B)       { benchVarArgus(b, "rgb(255,0,51)", "rgb") }
func BenchmarkPlayground_Var_RGB(b *testing.B)  { benchVarPlayground(b, "rgb(255,0,51)", "rgb") }

func BenchmarkArgus_VarString_RGBA(b *testing.B) {
	benchVarStringArgus(b, "rgba(255,0,51,0.5)", "rgba")
}
func BenchmarkArgus_Var_RGBA(b *testing.B)      { benchVarArgus(b, "rgba(255,0,51,0.5)", "rgba") }
func BenchmarkPlayground_Var_RGBA(b *testing.B) { benchVarPlayground(b, "rgba(255,0,51,0.5)", "rgba") }

func BenchmarkArgus_VarString_HSL(b *testing.B) { benchVarStringArgus(b, "hsl(0,100%,50%)", "hsl") }
func BenchmarkArgus_Var_HSL(b *testing.B)       { benchVarArgus(b, "hsl(0,100%,50%)", "hsl") }
func BenchmarkPlayground_Var_HSL(b *testing.B)  { benchVarPlayground(b, "hsl(0,100%,50%)", "hsl") }

func BenchmarkArgus_VarString_HSLA(b *testing.B) {
	benchVarStringArgus(b, "hsla(0,100%,50%,0.5)", "hsla")
}
func BenchmarkArgus_Var_HSLA(b *testing.B) { benchVarArgus(b, "hsla(0,100%,50%,0.5)", "hsla") }
func BenchmarkPlayground_Var_HSLA(b *testing.B) {
	benchVarPlayground(b, "hsla(0,100%,50%,0.5)", "hsla")
}

// ══════════════════════════════════════════════════════════════════════════════
// base32 / base64 / base64url / base64rawurl
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_Base32(b *testing.B) {
	benchVarStringArgus(b, "JBSWY3DPEHPK3PXP", "base32")
}
func BenchmarkArgus_Var_Base32(b *testing.B) { benchVarArgus(b, "JBSWY3DPEHPK3PXP", "base32") }

func BenchmarkArgus_VarString_Base64(b *testing.B) { benchVarStringArgus(b, "aGVsbG8=", "base64") }
func BenchmarkArgus_Var_Base64(b *testing.B)       { benchVarArgus(b, "aGVsbG8=", "base64") }
func BenchmarkPlayground_Var_Base64(b *testing.B)  { benchVarPlayground(b, "aGVsbG8=", "base64") }

func BenchmarkArgus_VarString_Base64URL(b *testing.B) {
	benchVarStringArgus(b, "aGVsbG8=", "base64url")
}
func BenchmarkArgus_Var_Base64URL(b *testing.B)      { benchVarArgus(b, "aGVsbG8=", "base64url") }
func BenchmarkPlayground_Var_Base64URL(b *testing.B) { benchVarPlayground(b, "aGVsbG8=", "base64url") }

func BenchmarkArgus_VarString_Base64RawURL(b *testing.B) {
	benchVarStringArgus(b, "aGVsbG8", "base64rawurl")
}
func BenchmarkArgus_Var_Base64RawURL(b *testing.B) { benchVarArgus(b, "aGVsbG8", "base64rawurl") }

// ══════════════════════════════════════════════════════════════════════════════
// json / mongodb
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_JSON(b *testing.B) { benchVarStringArgus(b, `{"key":"value"}`, "json") }
func BenchmarkArgus_Var_JSON(b *testing.B)       { benchVarArgus(b, `{"key":"value"}`, "json") }
func BenchmarkPlayground_Var_JSON(b *testing.B)  { benchVarPlayground(b, `{"key":"value"}`, "json") }

func BenchmarkArgus_VarString_MongoDB(b *testing.B) {
	benchVarStringArgus(b, "507f1f77bcf86cd799439011", "mongodb")
}
func BenchmarkArgus_Var_MongoDB(b *testing.B) {
	benchVarArgus(b, "507f1f77bcf86cd799439011", "mongodb")
}

// ══════════════════════════════════════════════════════════════════════════════
// bic / cron / datauri / bcp47
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_BIC(b *testing.B) { benchVarStringArgus(b, "DEUTDEFF", "bic") }
func BenchmarkArgus_Var_BIC(b *testing.B)       { benchVarArgus(b, "DEUTDEFF", "bic") }

func BenchmarkArgus_VarString_Cron(b *testing.B) { benchVarStringArgus(b, "*/5 * * * *", "cron") }
func BenchmarkArgus_Var_Cron(b *testing.B)       { benchVarArgus(b, "*/5 * * * *", "cron") }

func BenchmarkArgus_VarString_DataURI(b *testing.B) {
	benchVarStringArgus(b, "data:text/plain;base64,SGVsbG8=", "datauri")
}
func BenchmarkArgus_Var_DataURI(b *testing.B) {
	benchVarArgus(b, "data:text/plain;base64,SGVsbG8=", "datauri")
}

func BenchmarkArgus_VarString_BCP47(b *testing.B) { benchVarStringArgus(b, "en-US", "bcp47") }
func BenchmarkArgus_Var_BCP47(b *testing.B)       { benchVarArgus(b, "en-US", "bcp47") }

// ══════════════════════════════════════════════════════════════════════════════
// eth_addr / btc_addr
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_EthAddr(b *testing.B) {
	benchVarStringArgus(b, "0x71C7656EC7ab88b098defB751B7401B5f6d8976F", "eth_addr")
}
func BenchmarkArgus_Var_EthAddr(b *testing.B) {
	benchVarArgus(b, "0x71C7656EC7ab88b098defB751B7401B5f6d8976F", "eth_addr")
}

func BenchmarkArgus_VarString_BtcAddr(b *testing.B) {
	benchVarStringArgus(b, "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", "btc_addr")
}
func BenchmarkArgus_Var_BtcAddr(b *testing.B) {
	benchVarArgus(b, "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", "btc_addr")
}

// ══════════════════════════════════════════════════════════════════════════════
// isbn10 / isbn13 / issn
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_ISBN10(b *testing.B) { benchVarStringArgus(b, "0471958697", "isbn10") }
func BenchmarkArgus_Var_ISBN10(b *testing.B)       { benchVarArgus(b, "0471958697", "isbn10") }

func BenchmarkArgus_VarString_ISBN13(b *testing.B) {
	benchVarStringArgus(b, "978-0-471-95869-0", "isbn13")
}
func BenchmarkArgus_Var_ISBN13(b *testing.B) { benchVarArgus(b, "978-0-471-95869-0", "isbn13") }

func BenchmarkArgus_VarString_ISSN(b *testing.B) { benchVarStringArgus(b, "0317-8471", "issn") }
func BenchmarkArgus_Var_ISSN(b *testing.B)       { benchVarArgus(b, "0317-8471", "issn") }

// ══════════════════════════════════════════════════════════════════════════════
// e164 / mac / port
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_E164(b *testing.B) { benchVarStringArgus(b, "+12345678901", "e164") }
func BenchmarkArgus_Var_E164(b *testing.B)       { benchVarArgus(b, "+12345678901", "e164") }
func BenchmarkPlayground_Var_E164(b *testing.B)  { benchVarPlayground(b, "+12345678901", "e164") }

func BenchmarkArgus_VarString_MAC(b *testing.B) { benchVarStringArgus(b, "00:11:22:33:44:55", "mac") }
func BenchmarkArgus_Var_MAC(b *testing.B)       { benchVarArgus(b, "00:11:22:33:44:55", "mac") }
func BenchmarkPlayground_Var_MAC(b *testing.B)  { benchVarPlayground(b, "00:11:22:33:44:55", "mac") }

func BenchmarkArgus_VarString_Port(b *testing.B) { benchVarStringArgus(b, "8080", "port") }
func BenchmarkArgus_Var_Port(b *testing.B)       { benchVarArgus(b, "8080", "port") }

// ══════════════════════════════════════════════════════════════════════════════
// latitude / longitude / datetime / timezone
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_Latitude(b *testing.B) { benchVarStringArgus(b, "39.9042", "latitude") }
func BenchmarkArgus_Var_Latitude(b *testing.B)       { benchVarArgus(b, "39.9042", "latitude") }
func BenchmarkPlayground_Var_Latitude(b *testing.B)  { benchVarPlayground(b, "39.9042", "latitude") }

func BenchmarkArgus_VarString_Longitude(b *testing.B) {
	benchVarStringArgus(b, "116.4074", "longitude")
}
func BenchmarkArgus_Var_Longitude(b *testing.B)      { benchVarArgus(b, "116.4074", "longitude") }
func BenchmarkPlayground_Var_Longitude(b *testing.B) { benchVarPlayground(b, "116.4074", "longitude") }

func BenchmarkArgus_VarString_Datetime(b *testing.B) {
	benchVarStringArgus(b, "2026-05-18T10:00:00Z", "datetime")
}
func BenchmarkArgus_Var_Datetime(b *testing.B) { benchVarArgus(b, "2026-05-18T10:00:00Z", "datetime") }
func BenchmarkPlayground_Var_Datetime(b *testing.B) {
	benchVarPlayground(b, "2026-05-18T10:00:00Z", "datetime")
}

func BenchmarkArgus_VarString_Timezone(b *testing.B) {
	benchVarStringArgus(b, "America/New_York", "timezone")
}
func BenchmarkArgus_Var_Timezone(b *testing.B) { benchVarArgus(b, "America/New_York", "timezone") }

// ══════════════════════════════════════════════════════════════════════════════
// boolean / number / numeric / credit_card / luhn_checksum
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_Boolean(b *testing.B) { benchVarStringArgus(b, "true", "boolean") }
func BenchmarkArgus_Var_Boolean(b *testing.B)       { benchVarArgus(b, "true", "boolean") }
func BenchmarkPlayground_Var_Boolean(b *testing.B)  { benchVarPlayground(b, "true", "boolean") }

func BenchmarkArgus_VarString_Number(b *testing.B) { benchVarStringArgus(b, "12345", "number") }
func BenchmarkArgus_Var_Number(b *testing.B)       { benchVarArgus(b, "12345", "number") }
func BenchmarkPlayground_Var_Number(b *testing.B)  { benchVarPlayground(b, "12345", "number") }

func BenchmarkArgus_VarString_Numeric(b *testing.B) { benchVarStringArgus(b, "12345", "numeric") }
func BenchmarkArgus_Var_Numeric(b *testing.B)       { benchVarArgus(b, "12345", "numeric") }

func BenchmarkArgus_VarString_CreditCard(b *testing.B) {
	benchVarStringArgus(b, "4111111111111111", "credit_card")
}
func BenchmarkArgus_Var_CreditCard(b *testing.B) { benchVarArgus(b, "4111111111111111", "credit_card") }

func BenchmarkArgus_VarString_LuhnChecksum(b *testing.B) {
	benchVarStringArgus(b, "4111111111111111", "luhn_checksum")
}
func BenchmarkArgus_Var_LuhnChecksum(b *testing.B) {
	benchVarArgus(b, "4111111111111111", "luhn_checksum")
}

// ══════════════════════════════════════════════════════════════════════════════
// lowercase / uppercase
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_Lowercase(b *testing.B) { benchVarStringArgus(b, "hello", "lowercase") }
func BenchmarkArgus_Var_Lowercase(b *testing.B)       { benchVarArgus(b, "hello", "lowercase") }
func BenchmarkPlayground_Var_Lowercase(b *testing.B)  { benchVarPlayground(b, "hello", "lowercase") }

func BenchmarkArgus_VarString_Uppercase(b *testing.B) { benchVarStringArgus(b, "HELLO", "uppercase") }
func BenchmarkArgus_Var_Uppercase(b *testing.B)       { benchVarArgus(b, "HELLO", "uppercase") }
func BenchmarkPlayground_Var_Uppercase(b *testing.B)  { benchVarPlayground(b, "HELLO", "uppercase") }

// ══════════════════════════════════════════════════════════════════════════════
// startswith / endswith / startsnotwith / endsnotwith
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_StartsWith(b *testing.B) {
	benchVarStringArgus(b, "hello world", "startswith=hello")
}
func BenchmarkArgus_Var_StartsWith(b *testing.B) { benchVarArgus(b, "hello world", "startswith=hello") }
func BenchmarkPlayground_Var_StartsWith(b *testing.B) {
	benchVarPlayground(b, "hello world", "startswith=hello")
}

func BenchmarkArgus_VarString_EndsWith(b *testing.B) {
	benchVarStringArgus(b, "hello world", "endswith=world")
}
func BenchmarkArgus_Var_EndsWith(b *testing.B) { benchVarArgus(b, "hello world", "endswith=world") }
func BenchmarkPlayground_Var_EndsWith(b *testing.B) {
	benchVarPlayground(b, "hello world", "endswith=world")
}

func BenchmarkArgus_VarString_StartsNotWith(b *testing.B) {
	benchVarStringArgus(b, "hello world", "startsnotwith=xyz")
}
func BenchmarkArgus_Var_StartsNotWith(b *testing.B) {
	benchVarArgus(b, "hello world", "startsnotwith=xyz")
}

func BenchmarkArgus_VarString_EndsNotWith(b *testing.B) {
	benchVarStringArgus(b, "hello world", "endsnotwith=xyz")
}
func BenchmarkArgus_Var_EndsNotWith(b *testing.B) { benchVarArgus(b, "hello world", "endsnotwith=xyz") }

// ══════════════════════════════════════════════════════════════════════════════
// contains / containsany / containsrune / excludes / excludesall / excludesrune
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_Contains(b *testing.B) {
	benchVarStringArgus(b, "hello world", "contains=world")
}
func BenchmarkArgus_Var_Contains(b *testing.B) { benchVarArgus(b, "hello world", "contains=world") }
func BenchmarkPlayground_Var_Contains(b *testing.B) {
	benchVarPlayground(b, "hello world", "contains=world")
}

func BenchmarkArgus_VarString_ContainsAny(b *testing.B) {
	benchVarStringArgus(b, "hello", "containsany=aeiou")
}
func BenchmarkArgus_Var_ContainsAny(b *testing.B) { benchVarArgus(b, "hello", "containsany=aeiou") }

func BenchmarkArgus_VarString_ContainsRune(b *testing.B) {
	benchVarStringArgus(b, "hello", "containsrune=ö")
}
func BenchmarkArgus_Var_ContainsRune(b *testing.B) { benchVarArgus(b, "hello", "containsrune=ö") }

func BenchmarkArgus_VarString_Excludes(b *testing.B) { benchVarStringArgus(b, "hello", "excludes=xyz") }
func BenchmarkArgus_Var_Excludes(b *testing.B)       { benchVarArgus(b, "hello", "excludes=xyz") }

func BenchmarkArgus_VarString_ExcludesAll(b *testing.B) {
	benchVarStringArgus(b, "hello", "excludesall=xyz")
}
func BenchmarkArgus_Var_ExcludesAll(b *testing.B) { benchVarArgus(b, "hello", "excludesall=xyz") }

func BenchmarkArgus_VarString_ExcludesRune(b *testing.B) {
	benchVarStringArgus(b, "hello", "excludesrune=ö")
}
func BenchmarkArgus_Var_ExcludesRune(b *testing.B) { benchVarArgus(b, "hello", "excludesrune=ö") }

// ══════════════════════════════════════════════════════════════════════════════
// unique / eq_ignore_case / ne_ignore_case / isdefault
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_Unique(b *testing.B) { benchVarStringArgus(b, "abcdef", "unique") }
func BenchmarkArgus_Var_Unique(b *testing.B)       { benchVarArgus(b, "abcdef", "unique") }

func BenchmarkArgus_VarString_EqIgnoreCase(b *testing.B) {
	benchVarStringArgus(b, "Hello", "eq_ignore_case=hello")
}
func BenchmarkArgus_Var_EqIgnoreCase(b *testing.B) { benchVarArgus(b, "Hello", "eq_ignore_case=hello") }

func BenchmarkArgus_VarString_NeIgnoreCase(b *testing.B) {
	benchVarStringArgus(b, "Hello", "ne_ignore_case=world")
}
func BenchmarkArgus_Var_NeIgnoreCase(b *testing.B) { benchVarArgus(b, "Hello", "ne_ignore_case=world") }

func BenchmarkArgus_VarString_IsDefault(b *testing.B) { benchVarStringArgus(b, "", "isdefault") }
func BenchmarkArgus_Var_IsDefault(b *testing.B)       { benchVarArgus(b, "", "isdefault") }

// ══════════════════════════════════════════════════════════════════════════════
// html / html_encoded / url_encoded
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_HTML(b *testing.B) { benchVarStringArgus(b, "<div>hello</div>", "html") }
func BenchmarkArgus_Var_HTML(b *testing.B)       { benchVarArgus(b, "<div>hello</div>", "html") }
func BenchmarkPlayground_Var_HTML(b *testing.B)  { benchVarPlayground(b, "<div>hello</div>", "html") }

func BenchmarkArgus_VarString_HTMLEncoded(b *testing.B) {
	benchVarStringArgus(b, "&lt;div&gt;", "html_encoded")
}
func BenchmarkArgus_Var_HTMLEncoded(b *testing.B) { benchVarArgus(b, "&lt;div&gt;", "html_encoded") }
func BenchmarkPlayground_Var_HTMLEncoded(b *testing.B) {
	benchVarPlayground(b, "&lt;div&gt;", "html_encoded")
}

func BenchmarkArgus_VarString_URLEncoded(b *testing.B) {
	benchVarStringArgus(b, "hello%20world", "url_encoded")
}
func BenchmarkArgus_Var_URLEncoded(b *testing.B) { benchVarArgus(b, "hello%20world", "url_encoded") }
func BenchmarkPlayground_Var_URLEncoded(b *testing.B) {
	benchVarPlayground(b, "hello%20world", "url_encoded")
}

// ══════════════════════════════════════════════════════════════════════════════
// dns_rfc1035_label / file / filepath / dir / dirpath
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_DNSRFC1035Label(b *testing.B) {
	benchVarStringArgus(b, "my-label", "dns_rfc1035_label")
}
func BenchmarkArgus_Var_DNSRFC1035Label(b *testing.B) {
	benchVarArgus(b, "my-label", "dns_rfc1035_label")
}

func BenchmarkArgus_VarString_File(b *testing.B) { benchVarStringArgus(b, "test.txt", "file") }
func BenchmarkArgus_Var_File(b *testing.B)       { benchVarArgus(b, "test.txt", "file") }

func BenchmarkArgus_VarString_FilePath(b *testing.B) {
	benchVarStringArgus(b, "/tmp/test.txt", "filepath")
}
func BenchmarkArgus_Var_FilePath(b *testing.B) { benchVarArgus(b, "/tmp/test.txt", "filepath") }

func BenchmarkArgus_VarString_Dir(b *testing.B) { benchVarStringArgus(b, "testdata", "dir") }
func BenchmarkArgus_Var_Dir(b *testing.B)       { benchVarArgus(b, "testdata", "dir") }

func BenchmarkArgus_VarString_DirPath(b *testing.B) { benchVarStringArgus(b, "/tmp", "dirpath") }
func BenchmarkArgus_Var_DirPath(b *testing.B)       { benchVarArgus(b, "/tmp", "dirpath") }

// ══════════════════════════════════════════════════════════════════════════════
// MultiRule
// ══════════════════════════════════════════════════════════════════════════════

func BenchmarkArgus_VarString_MultiRule(b *testing.B) {
	benchVarStringArgus(b, "hello world", "required,min=1,max=100")
}
func BenchmarkArgus_Var_MultiRule(b *testing.B) {
	benchVarArgus(b, "hello world", "required,min=1,max=100")
}
func BenchmarkPlayground_Var_MultiRule(b *testing.B) {
	benchVarPlayground(b, "hello world", "required,min=1,max=100")
}
