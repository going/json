package vm

import (
	"fmt"
	"io"

	"github.com/going/json/internal/encoder"
)

func DebugRun(ctx *encoder.RuntimeContext, b []byte, codeSet *encoder.OpcodeSet) ([]byte, error) {
	defer func() {
		var code *encoder.Opcode
		if (ctx.Option.Flag & encoder.HTMLEscapeOption) != 0 {
			code = codeSet.EscapeKeyCode
		} else {
			code = codeSet.NoescapeKeyCode
		}
		if wc := ctx.Option.DebugDOTOut; wc != nil {
			_, _ = io.WriteString(wc, code.DumpDOT())
			wc.Close()
			ctx.Option.DebugDOTOut = nil
		}

		if err := recover(); err != nil {
			w := ctx.Option.DebugOut
			fmt.Fprintln(w, "=============[DEBUG]===============")
			fmt.Fprintln(w, "* [TYPE]")
			fmt.Fprintln(w, codeSet.Type)
			fmt.Fprintf(w, "\n")
			fmt.Fprintln(w, "* [ALL OPCODE]")
			fmt.Fprintln(w, code.Dump())
			fmt.Fprintf(w, "\n")
			fmt.Fprintln(w, "* [CONTEXT]")
			fmt.Fprintf(w, "%+v\n", ctx)
			fmt.Fprintln(w, "===================================")
			panic(err)
		}
	}()

	return Run(ctx, b, codeSet)
}
