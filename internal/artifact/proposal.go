package artifact

import (
	"bufio"
	"strings"
)

type Proposal struct {
	Sections      []string
	Why           string
	WhatChanges   string
	Goals         string
	NonGoals      string
	Assumptions   string
	OpenQuestions string
}

func ParseProposalFile(path string) (*Proposal, error) {
	content, err := ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseProposalContent(content), nil
}

func ParseProposalContent(content string) *Proposal {
	prop := &Proposal{}
	var currentSection string
	var sectionBuf *strings.Builder

	flushSection := func() {
		if sectionBuf == nil {
			return
		}
		text := strings.TrimSpace(sectionBuf.String())
		switch currentSection {
		case "Why":
			prop.Why = text
		case "What Changes":
			prop.WhatChanges = text
		case "Goals":
			prop.Goals = text
		case "Non-Goals":
			prop.NonGoals = text
		case "Key Assumptions":
			prop.Assumptions = text
		case "Open Questions":
			prop.OpenQuestions = text
		}
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		if m := sectionRe.FindStringSubmatch(line); m != nil {
			flushSection()
			currentSection = m[1]
			sectionBuf = &strings.Builder{}
			prop.Sections = append(prop.Sections, m[1])
			continue
		}

		if sectionBuf != nil {
			if sectionBuf.Len() > 0 {
				sectionBuf.WriteString("\n")
			}
			sectionBuf.WriteString(line)
		}
	}

	flushSection()

	return prop
}

func (p *Proposal) SectionNames() []string {
	return p.Sections
}

func (p *Proposal) HasSection(name string) bool {
	for _, s := range p.Sections {
		if s == name {
			return true
		}
	}
	return false
}
