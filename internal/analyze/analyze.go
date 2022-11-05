package analyze

import (
	"encoding/hex"
	"net/mail"
	"net/url"
	"strings"

	"go.arsenm.dev/lure-repo-bot/internal/spdx"
	"golang.org/x/exp/slices"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type Finding struct {
	ItemType string
	ItemName string
	Line     uint
	Index    any
	Msg      string
	ExtraMsg string
}

func AnalyzeScript(r *interp.Runner, fl *syntax.File) ([]Finding, error) {
	var findings []Finding

	if _, ok := r.Vars["name"]; !ok {
		findings = append(findings, Finding{
			ItemType: "variable",
			ItemName: "name",
			Msg:      "The %s is required",
		})
	}

	if _, ok := r.Vars["version"]; !ok {
		findings = append(findings, Finding{
			ItemType: "variable",
			ItemName: "version",
			Msg:      "The %s is required",
		})
	}

	if _, ok := r.Vars["release"]; !ok {
		findings = append(findings, Finding{
			ItemType: "variable",
			ItemName: "release",
			Msg:      "The %s is required",
		})
	}

	if _, ok := r.Funcs["package"]; !ok {
		findings = append(findings, Finding{
			ItemType: "function",
			ItemName: "package",
			Msg:      "The %s is required",
		})
	}

	for name, scriptVar := range r.Vars {
		_, scriptVar = scriptVar.Resolve(r.Env)
		val := getVal(&scriptVar)

		// Remove any override suffix, so that we
		// check all the overrides as well
		cutName, _, _ := strings.Cut(name, "_")

		// build_vars has an underscore, and thus is a special case
		// that must be accounted for
		if strings.HasPrefix("name", "build_vars") {
			cutName = "build_vars"
		}

		switch cutName {
		case "release":
			valStr, ok := mustBeStr(val, name, &findings)
			if !ok {
				continue
			}

			if !isNumeric(strings.TrimPrefix(valStr, "-")) {
				findings = append(findings, Finding{
					ItemType: "variable",
					ItemName: name,
					Msg:      "The %s must be an integer",
				})
				continue
			}
		case "epoch":
			valStr, ok := mustBeStr(val, name, &findings)
			if !ok {
				continue
			}

			if !isNumeric(valStr) {
				findings = append(findings, Finding{
					ItemType: "variable",
					ItemName: name,
					Msg:      "The %s must be a positive integer",
				})
				continue
			}
		case "homepage":
			valStr, ok := mustBeStr(val, name, &findings)
			if !ok {
				continue
			}

			_, err := url.ParseRequestURI(valStr)
			if err != nil {
				findings = append(findings, Finding{
					ItemType: "variable",
					ItemName: name,
					Msg:      "The %s must be a valid URL",
				})
				continue
			}
		case "maintainer":
			valStr, ok := mustBeStr(val, name, &findings)
			if !ok {
				continue
			}

			addr, err := mail.ParseAddress(valStr)
			if err != nil {
				findings = append(findings, Finding{
					ItemType: "variable",
					ItemName: name,
					Msg:      "The %s must be a valid RFC 5322 address",
				})
				continue
			}

			if addr.Name == "" || addr.Address == "" {
				findings = append(findings, Finding{
					ItemType: "variable",
					ItemName: name,
					Msg:      "The %s must contain a name and email (e.g. Arsen Musayelyan <arsen@arsenm.dev>)",
				})
				continue
			}
		case "architectures":
			valSlice, ok := mustBeArray(val, name, &findings)
			if !ok {
				continue
			}

			if slices.Contains(valSlice, "noarch") || slices.Contains(valSlice, "any") {
				findings = append(findings, Finding{
					ItemType: "variable",
					ItemName: name,
					Msg:      "The %s must be set to 'all' to represent noarch/any",
				})
				continue
			}
		case "license":
			valSlice, ok := mustBeArray(val, name, &findings)
			if !ok {
				continue
			}

			for _, val := range valSlice {
				if strings.Contains(strings.ToLower(val), "custom") {
					continue
				}

				license := spdx.Licenses.License(val)
				if license == nil {
					similar := spdx.FindSimilarLicense(val)

					f := Finding{
						ItemType: "variable",
						ItemName: name,
						Msg:      "The %s contains an invalid SPDX license identifier: '" + val + "'.",
						ExtraMsg: "A list of SPDX license identifiers can be found at https://spdx.org/licenses/.",
					}

					if similar != "" {
						f.Msg += " Did you mean '" + similar + "'?"
					}

					findings = append(findings, f)
					continue
				}
			}
		case "provides":
			mustBeArray(val, name, &findings)
		case "conflicts":
			mustBeArray(val, name, &findings)
		case "deps":
			mustBeArray(val, name, &findings)
		case "build_deps":
			mustBeArray(val, name, &findings)
		case "replaces":
			mustBeArray(val, name, &findings)
		case "sources":
			valSlice, ok := mustBeArray(val, name, &findings)
			if !ok {
				continue
			}

			for i, val := range valSlice {
				u, err := url.ParseRequestURI(val)
				if err != nil {
					findings = append(findings, Finding{
						ItemType: "element",
						Index:    i,
						ItemName: name,
						Msg:      "The %s must be a valid URL",
					})
					continue
				}
				query := u.Query()

				var validParams []string
				if strings.HasPrefix(u.Scheme, "git+") {
					validParams = []string{"tag", "branch", "commit", "depth", "name"}
				} else {
					validParams = []string{"archive"}
				}

				for paramName := range query {
					if strings.HasPrefix(paramName, "~") {
						paramName = strings.TrimPrefix(paramName, "~")
					} else {
						continue
					}

					if !slices.Contains(validParams, paramName) {
						findings = append(findings, Finding{
							ItemType: "element",
							ItemName: name,
							Index:    i,
							Msg:      "The %s contains an invalid parameter name '~" + paramName + "'",
						})
						continue
					}
				}
			}
		case "checksums":
			valSlice, ok := mustBeArray(val, name, &findings)
			if !ok {
				continue
			}

			sourcesName := strings.Replace(name, "checksums", "sources", 1)
			srcs, ok := r.Vars[sourcesName]
			if !ok || len(srcs.List) != len(valSlice) {
				findings = append(findings, Finding{
					ItemType: "array",
					ItemName: name,
					Msg:      "The %s is not the same size as its corresponding sources array",
				})
			}

			for i, val := range valSlice {
				if val != "SKIP" && len(val) != 64 {
					findings = append(findings, Finding{
						ItemType: "element",
						ItemName: name,
						Index:    i,
						Msg:      "The %s contains an invalid SHA256 checksum. SHA256 hashes must be 64 characters in length.",
					})
					continue
				}

				if val != "SKIP" {
					_, err := hex.DecodeString(val)
					if err != nil {
						findings = append(findings, Finding{
							ItemType: "element",
							ItemName: name,
							Index:    i,
							Msg:      "The %s contains an invalid SHA256 checksum. SHA256 hashes must be valid hexadecimal.",
						})
						continue
					}
				}
			}
		case "backup":
			mustBeArray(val, name, &findings)
		case "scripts":
			mustBeMap(val, name, &findings)
		}
	}

	lns := FindLines(fl)
	for i, finding := range findings {
		if finding.ItemType == "function" {
			ln, ok := lns.Funcs[finding.ItemName]
			if ok {
				findings[i].Line = ln
			}
		} else {
			ln, ok := lns.Vars[finding.ItemName]
			if ok {
				findings[i].Line = ln
			}
		}
	}

	return findings, nil
}

func mustBeStr(val any, name string, findings *[]Finding) (string, bool) {
	valStr, ok := val.(string)
	if !ok {
		*findings = append(*findings, Finding{
			ItemType: "variable",
			ItemName: name,
			Msg:      "The %s must be a string",
		})
	}
	return valStr, ok
}

func mustBeArray(val any, name string, findings *[]Finding) ([]string, bool) {
	valSlice, ok := val.([]string)
	if !ok {
		*findings = append(*findings, Finding{
			ItemType: "variable",
			ItemName: name,
			Msg:      "The %s must be an array",
		})
	}
	return valSlice, ok
}

func mustBeMap(val any, name string, findings *[]Finding) (map[string]string, bool) {
	valMap, ok := val.(map[string]string)
	if !ok {
		*findings = append(*findings, Finding{
			ItemType: "variable",
			ItemName: name,
			Msg:      "The %s must be a map",
		})
	}
	return valMap, ok
}

func getVal(v *expand.Variable) any {
	if v.Str != "" {
		return v.Str
	} else if v.List != nil {
		return v.List
	} else if v.Map != nil {
		return v.Map
	}

	return nil
}

func isNumeric(s string) bool {
	index := strings.IndexFunc(s, func(r rune) bool {
		return r < '0' || r > '9'
	})
	return index == -1
}

type Lines struct {
	Vars  map[string]uint
	Funcs map[string]uint
}

func FindLines(fl *syntax.File) Lines {
	out := Lines{map[string]uint{}, map[string]uint{}}

	for _, stmt := range fl.Stmts {
		switch cmd := stmt.Cmd.(type) {
		case *syntax.CallExpr:
			if len(cmd.Assigns) == 0 {
				continue
			}

			name := cmd.Assigns[0].Name.Value
			line := cmd.Assigns[0].Pos().Line()
			out.Vars[name] = line
		case *syntax.FuncDecl:
			out.Funcs[cmd.Name.Value] = cmd.Pos().Line()
		}
	}

	return out
}
