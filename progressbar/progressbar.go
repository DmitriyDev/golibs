package progressbar

import (
	"fmt"
	"math"
	"strings"
)

const LineLength = 140
const PercentLengthMultiplier = 1
const percentGrade = 10

type InlinePrint struct {
	percentLengthMultiplier int
	percentGrade            int
	lineLength              int
}

//
func (ip InlinePrint) New(percentLengthMultiplier int, lineLength int) InlinePrint {

	if percentLengthMultiplier <= 0 {
		percentLengthMultiplier = PercentLengthMultiplier
	}

	if lineLength <= 0 {
		lineLength = LineLength
	}

	return InlinePrint{
		percentLengthMultiplier: percentLengthMultiplier,
		percentGrade:            percentGrade,
		lineLength:              lineLength,
	}
}

func (ip InlinePrint) NewDefault() InlinePrint {
	return InlinePrint{
		percentLengthMultiplier: PercentLengthMultiplier,
		percentGrade:            percentGrade,
		lineLength:              LineLength,
	}
}

func (ip *InlinePrint) PrintPercentLine(p int) {

	full := int(math.Ceil((float64(p) / float64(ip.percentGrade)) * float64(ip.percentLengthMultiplier)))
	emt := (ip.percentGrade*ip.percentLengthMultiplier - full)
	fStr := strings.Repeat("▓", full)
	if emt <= 0 {
		emt = 0
	}
	eStr := strings.Repeat("░", emt)

	txt := fmt.Sprintf("[%s%s] %d", fStr, eStr, p) + "%"
	ip.print(txt)

}

func (ip *InlinePrint) PrintStringLine(t string) {
	ip.print(t)
}

func (ip *InlinePrint) print(t string) {

	t = strings.Trim(t, "\n")
	l := len(t)

	if l > ip.lineLength {
		t = t[0:ip.lineLength-5] + "..."
	} else {
		t = t + strings.Repeat(" ", ip.lineLength-2-l)
	}

	fmt.Print("\r" + t)
}
