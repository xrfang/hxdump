package hxdump

import (
	"fmt"
	"strings"
)

//Style CSS style for hexdump
type Style struct {
	EvenColor   string
	FontSize    string
	Narrow      bool
	OddColor    string
	OffsetColor string
}

//Render dump data as HTML with default style
func Render(data []byte) (style string, html string) {
	return RenderWithStyle(data, Style{})
}

//RenderWithStyle dump data as HTML with given style
func RenderWithStyle(data []byte, s Style) (style string, html string) {
	if s.EvenColor == "" {
		s.EvenColor = "black"
	}
	if s.FontSize == "" {
		s.FontSize = "12px"
	}
	if s.OddColor == "" {
		s.OddColor = "green"
	}
	if s.OffsetColor == "" {
		s.OffsetColor = "#727272"
	}
	width := 32
	if s.Narrow {
		width = 16
	}
	var hex, asc []string //分别保存hex序列和ascii序列
	var row strings.Builder
	renderRow := func(off int, h, a []string) {
		fmt.Fprintln(&row, `<div class="xd">`)
		fmt.Fprintf(&row, "<div class=\"xd-offset\">%05d: </div>", off)
		var hs, as []string
		for k := 0; k < width/4; k++ {
			hs = append(hs, fmt.Sprintf("%s ", strings.Join(h[k*4:(k+1)*4], "")))
			as = append(as, strings.Join(a[k*4:(k+1)*4], ""))
		}
		for k := 0; k < len(hs)/2; k++ {
			if k%2 == 0 {
				fmt.Fprintf(&row, "<div class=\"xd-col-even\">%s</div>", strings.Join(hs[k*2:(k+1)*2], ""))
			} else {
				fmt.Fprintf(&row, "<div class=\"xd-col-odd\">%s</div>", strings.Join(hs[k*2:(k+1)*2], ""))
			}
		}
		row.WriteString("|\n")
		for k := 0; k < len(as)/2; k++ {
			if k%2 == 0 {
				fmt.Fprintf(&row, "<div class=\"xd-col-even\">%s</div>", strings.Join(as[k*2:(k+1)*2], ""))
			} else {
				fmt.Fprintf(&row, "<div class=\"xd-col-odd\">%s</div>", strings.Join(as[k*2:(k+1)*2], ""))
			}
		}
		fmt.Fprintln(&row, `</div>`)
	}
	for i := 0; i < len(data); i++ {
		if data[i] < 32 || data[i] > 126 {
			hex = append(hex, fmt.Sprintf(`<span class="xd-hex">%02x</span>`, data[i]))
			asc = append(asc, `<span class="xd-hex">░</span>`)
		} else if data[i] == 32 {
			hex = append(hex, `<span class="xd-hex">20</span>`)
			asc = append(asc, `<span class="xd-hex">▯</span>`)
		} else {
			hex = append(hex, fmt.Sprintf(`<span class="xd-asc">%02x</span>`, data[i]))
			asc = append(asc, fmt.Sprintf(`<span class="xd-asc">%s</span>`, string(data[i])))
		}
	}
	if len(data)%width != 0 {
		for i := 0; i < width-len(data)%width; i++ {
			hex = append(hex, `<span class="xd-hex">  </span>`)
			asc = append(asc, `<span class="xd-hex"> </span>`)
		}
	}
	for i := 0; i < len(hex)/width; i++ {
		renderRow(i*width, hex[i*width:(i+1)*width], asc[i*width:(i+1)*width])
	}
	return fmt.Sprintf(`<style>
    .xd {text-align: left;font-family:Courier New;font-size:%s}
    .xd-asc {font-weight:900}
    .xd-hex {font-weight:100}
    .xd-col-even {display:inline-block;color:%s;white-space:pre}
    .xd-col-odd {display:inline-block;color:%s;white-space:pre}
    .xd-offset {display:inline-block;color:%s;white-space:pre}
</style>`, s.FontSize, s.EvenColor, s.OddColor, s.OffsetColor), row.String()
}

//Dump dump data as pure text with default style
func Dump(data []byte) string {
	return DumpWithStyle(data, Style{})
}

//DumpWithStyle dump data as pure text with given style
func DumpWithStyle(data []byte, s Style) string {
	width := 32
	if s.Narrow {
		width = 16
	}
	var hex, asc []string //分别保存hex序列和ascii序列
	var row strings.Builder
	renderRow := func(off int, h, a []string) {
		fmt.Fprintf(&row, "%05d: ", off)
		var hs, as []string
		for k := 0; k < width/4; k++ {
			hs = append(hs, fmt.Sprintf("%s ", strings.Join(h[k*4:(k+1)*4], "")))
			as = append(as, strings.Join(a[k*4:(k+1)*4], ""))
		}
		for k := 0; k < len(hs); k++ {
			fmt.Fprint(&row, hs[k])
		}
		row.WriteString(" | ")
		for k := 0; k < len(as); k++ {
			fmt.Fprint(&row, as[k])
		}
		fmt.Fprintln(&row)
	}
	for i := 0; i < len(data); i++ {
		if data[i] < 32 || data[i] > 126 {
			hex = append(hex, fmt.Sprintf(`%02x`, data[i]))
			asc = append(asc, `░`)
		} else if data[i] == 32 {
			hex = append(hex, `20`)
			asc = append(asc, `▯`)
		} else {
			hex = append(hex, fmt.Sprintf(`%02x`, data[i]))
			asc = append(asc, string(data[i]))
		}
	}
	if len(data)%width != 0 {
		for i := 0; i < width-len(data)%width; i++ {
			hex = append(hex, `  `)
			asc = append(asc, ` `)
		}
	}
	for i := 0; i < len(hex)/width; i++ {
		renderRow(i*width, hex[i*width:(i+1)*width], asc[i*width:(i+1)*width])
	}
	return row.String()
}
