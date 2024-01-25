package parser

import (
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"strings"
)

type Define struct {
	Name  string
	Value string
}

type Target struct {
	Name      string
	Depends   []*Target
	IsRaw     bool
	RawValues []string
}

type Parsed struct {
	Targets []*Target
	Defines []*Define
}

var (
	defineRegex = regexp.MustCompile("^([\\w]+)=(.*)$")
	targetRegex = regexp.MustCompile("^([^:]+):(.*)$")
)

type state int

const (
	next state = iota
	define
	in_target
)

func removeEmpty(n []string) []string {
	result := make([]string, 0, len(n))
	for _, s := range n {
		s := strings.TrimSpace(s)
		if s != "" {
			result = append(result, s)
		}
	}

	return result
}

func replaceLineCon(unlined string) string {
	return strings.ReplaceAll(unlined, "\\\n", "")
}

func Parse(r io.Reader) (*Parsed, error) {
	linesbytes, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read all content of the make file: %w", err)
	}
	content := replaceLineCon(string(linesbytes))

	result := &Parsed{
		Targets: make([]*Target, 0),
		Defines: make([]*Define, 0),
	}

	definemap := make(map[string]*Define)
	targetmap := make(map[string]*Target)

	s := next
	targetname := ""
	var target *Target
	definename := ""

parserloop:
	for _, l := range strings.Split(content, "\n") {
		l := strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}

		switch s {
		case next:
			defines := defineRegex.FindStringSubmatch(l)
			if len(defines) > 0 {
				definename = defines[1]
				d := &Define{
					Name: defines[1], Value: defines[1],
				}
				definemap[definename] = d
				result.Defines = append(result.Defines, d)
				s = next
				continue parserloop
			}
			targets := targetRegex.FindStringSubmatch(l)
			if len(targets) > 0 {
				targetname = targets[1]
				target = &Target{Name: targetname, RawValues: []string{targets[2]}}
				result.Targets = append(result.Targets, target)
				targetmap[targetname] = target
				s = in_target
				continue parserloop
			}
			slog.Error("Unknown", "line", l)
		case in_target:
			defines := defineRegex.FindStringSubmatch(l)
			if len(defines) > 0 {
				definename = defines[1]
				d := &Define{
					Name: defines[1], Value: defines[1],
				}
				definemap[definename] = d
				result.Defines = append(result.Defines, d)
				s = next
				continue parserloop
			}
			targets := targetRegex.FindStringSubmatch(l)
			if len(targets) > 0 {
				targetname = targets[1]
				target = &Target{Name: targetname, RawValues: []string{targets[2]}}
				result.Targets = append(result.Targets, target)
				targetmap[targetname] = target
				s = in_target
				continue parserloop
			}
			target.RawValues = append(target.RawValues, l)
		}
	}

	return result, nil
}
