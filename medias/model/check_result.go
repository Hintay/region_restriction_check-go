package model

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

const (
	CheckResultYes           = "yes"
	CheckResultNo            = "no"
	CheckResultUnexpected    = "unexpected"
	CheckResultFailed        = "failed"
	CheckResultOverseaOnly   = "oversea_only"
	CheckResultOriginalsOnly = "originals_only"

	printPadding = 48
)

type CheckResult struct {
	Result  string
	Media   string
	Region  string
	Type    string
	Task    string
	Message string
}

func (c *CheckResult) Yes(msg ...interface{}) {
	c.Result = CheckResultYes
	if len(msg) > 0 {
		c.Message = fmt.Sprint(msg...)
	}
}

func (c *CheckResult) No(msg ...interface{}) {
	c.Result = CheckResultNo
	if len(msg) > 0 {
		c.Message = fmt.Sprint(msg...)
	}
}

func (c *CheckResult) Oversea() {
	c.Result = CheckResultOverseaOnly
}

func (c *CheckResult) OriginalsOnly() {
	c.Result = CheckResultOriginalsOnly
}

func (c *CheckResult) Unexpected(msg ...interface{}) {
	c.Result = CheckResultUnexpected
	c.Message = fmt.Sprint(msg...)
}

func (c *CheckResult) UnexpectedStatusCode(code int) {
	c.Result = CheckResultUnexpected
	c.Message = fmt.Sprintf("status code: %d", code)
}

func (c *CheckResult) Failed(msg ...interface{}) {
	c.Result = CheckResultFailed
	c.Message = fmt.Sprint(msg...)
}

type CheckResultSlice []*CheckResult

func (c *CheckResultSlice) Len() int {
	return len(*c)
}

func (c *CheckResultSlice) Swap(i, j int) {
	(*c)[i], (*c)[j] = (*c)[j], (*c)[i]
}

func (c *CheckResultSlice) Less(i, j int) bool {
	if (*c)[i].Region < (*c)[j].Region {
		return true
	} else if (*c)[i].Region == (*c)[j].Region {
		if (*c)[i].Type < (*c)[j].Type {
			return true
		} else if (*c)[i].Type == (*c)[j].Type {
			return (*c)[i].Media < (*c)[j].Media
		}
	}
	return false
}

func (c *CheckResultSlice) PrintTo(writer io.Writer) {
	w := tabwriter.NewWriter(writer, 8, 8, 0, ' ', 0)
	lastRegion := ""
	lastOttType := ""
	for _, res := range *c {
		region, ok := HumanReadableRegions[res.Region]
		if !ok {
			region = res.Region
		}

		if lastRegion != region {
			w.Flush()
			fmt.Fprintf(writer, "\n==========[ %s ]==========\n", region)
			lastRegion = region
		}
		if lastOttType != res.Type {
			lastOttType = res.Type
			if lastOttType != "" {
				w.Flush()
				fmt.Fprintf(writer, "\n------< %s - %s >------\n", region, res.Type)
			}
		}

		s := HumanReadableNames[res.Media]
		if len(s) == 0 {
			s = res.Media
		}
		pad := printPadding - len(s)
		for i := 0; i < pad; i++ {
			s += " "
		}
		s += "\t"
		s += strings.ToUpper(res.Result)

		if res.Message != "" {
			s += fmt.Sprintf(" (%s)", res.Message)
		}
		fmt.Fprintln(w, s)
	}
	w.Flush()
}
