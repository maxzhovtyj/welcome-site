package logflags

import (
	"flag"
	"fmt"
	"github.com/valyala/bytebufferpool"
	"io"
	"log"
	"runtime"
	"strings"
)

func PrintAllFlags(w io.Writer) {
	fmt.Fprintf(w, "Flag values\n")
	flag.VisitAll(func(f *flag.Flag) {
		k := f.Name
		v := f.Value.String()
		if strings.Contains(k, "assword") || strings.Contains(k, "Key") {
			v = "hidden"
		}
		fmt.Fprintf(w, "\t%s=%q\n", k, v)
	})
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tgoVersion=%s\n", runtime.Version())
	fmt.Fprintf(w, "\tGOMAXPROCS=%d\n", runtime.GOMAXPROCS(-1))
	fmt.Fprintf(w, "\tNumCPU=%d\n", runtime.NumCPU())
}

func LogAllFlags() {
	w := bytebufferpool.Get()
	PrintAllFlags(w)
	log.Printf("%s\n", w.B)
	bytebufferpool.Put(w)
}
