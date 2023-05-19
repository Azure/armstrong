package azidx

import (
	"fmt"
	"regexp"
	"strings"
)

type Matcher struct {
	PrefixSep bool
	Separater string
	Segments  []MatchSegment
}

type MatchSegment struct {
	Value      string
	IsWildcard bool
	IsAny      bool
}

func (m Matcher) Match(input string) bool {
	regstrs := []string{}
	for _, seg := range m.Segments {
		if !seg.IsWildcard {
			regstrs = append(regstrs, seg.Value)
			continue
		}
		if seg.IsAny {
			regstrs = append(regstrs, ".+")
		} else {
			regstrs = append(regstrs, fmt.Sprintf("[^%s]+", m.Separater))
		}
	}
	regstr := strings.Join(regstrs, m.Separater)
	if m.PrefixSep {
		regstr = m.Separater + regstr
	}
	return regexp.MustCompile("^" + regstr + "$").MatchString(input)
}

func (m Matcher) Less(om Matcher) bool {
	var fixCnt1, fixCnt2, wildcardCnt1, wildcardCnt2, anyCnt1, anyCnt2 int
	for _, seg := range m.Segments {
		if !seg.IsWildcard {
			fixCnt1++
			continue
		}
		if !seg.IsAny {
			wildcardCnt1++
			continue
		}
		anyCnt1++
	}
	for _, seg := range om.Segments {
		if !seg.IsWildcard {
			fixCnt2++
			continue
		}
		if !seg.IsAny {
			wildcardCnt2++
			continue
		}
		anyCnt2++
	}

	if anyCnt1 != anyCnt2 {
		return anyCnt1 < anyCnt2
	}

	if wildcardCnt1 != wildcardCnt2 {
		return wildcardCnt1 < wildcardCnt2
	}

	if fixCnt1 != fixCnt2 {
		return fixCnt1 < fixCnt2
	}

	// If all are the same, then compare per segment
	for idx := 0; idx < len(m.Segments); idx++ {
		seg1, seg2 := m.Segments[idx], om.Segments[idx]
		if seg1.IsWildcard != seg2.IsWildcard {
			return !seg1.IsWildcard
		}
		// Both are not wildcard
		if !seg1.IsWildcard && seg1.Value != seg2.Value {
			return seg1.Value < seg2.Value
		}
		// Both are wildcard
		if seg1.IsAny != seg2.IsAny {
			return !seg1.IsAny
		}
	}
	return false
}

type Matchers []Matcher

func (m Matchers) Len() int {
	return len(m)
}

func (m Matchers) Less(i int, j int) bool {
	return m[i].Less(m[j])
}

func (m Matchers) Swap(i int, j int) {
	m[i], m[j] = m[j], m[i]
}
