package vm_color_indent

import (
	"encoding/json"
	"fmt"
	"unsafe"

	"github.com/going/json/internal/encoder"
	"github.com/going/json/internal/runtime"
)

const uintptrSize = 4 << (^uintptr(0) >> 63)

var (
	appendIndent        = encoder.AppendIndent
	appendStructEnd     = encoder.AppendStructEndIndent
	errUnsupportedValue = encoder.ErrUnsupportedValue
	errUnsupportedFloat = encoder.ErrUnsupportedFloat
	mapiterinit         = encoder.MapIterInit
	mapiterkey          = encoder.MapIterKey
	mapitervalue        = encoder.MapIterValue
	mapiternext         = encoder.MapIterNext
	maplen              = encoder.MapLen
)

type emptyInterface struct {
	typ *runtime.Type
	ptr unsafe.Pointer
}

type nonEmptyInterface struct {
	itab *struct {
		ityp *runtime.Type // static interface type
		typ  *runtime.Type // dynamic concrete type
		// unused fields...
	}
	ptr unsafe.Pointer
}

func errUnimplementedOp(op encoder.OpType) error {
	return fmt.Errorf("encoder (indent): opcode %s has not been implemented", op)
}

func load(base uintptr, idx uint32) uintptr {
	addr := base + uintptr(idx)
	return **(**uintptr)(unsafe.Pointer(&addr))
}

func store(base uintptr, idx uint32, p uintptr) {
	addr := base + uintptr(idx)
	**(**uintptr)(unsafe.Pointer(&addr)) = p
}

func loadNPtr(base uintptr, idx uint32, ptrNum uint8) uintptr {
	addr := base + uintptr(idx)
	p := **(**uintptr)(unsafe.Pointer(&addr))
	for i := uint8(0); i < ptrNum; i++ {
		if p == 0 {
			return 0
		}
		p = ptrToPtr(p)
	}
	return p
}

func ptrToUint64(p uintptr, bitSize uint8) uint64 {
	switch bitSize {
	case 8:
		return (uint64)(**(**uint8)(unsafe.Pointer(&p)))
	case 16:
		return (uint64)(**(**uint16)(unsafe.Pointer(&p)))
	case 32:
		return (uint64)(**(**uint32)(unsafe.Pointer(&p)))
	case 64:
		return **(**uint64)(unsafe.Pointer(&p))
	}
	return 0
}

func ptrToFloat32(p uintptr) float32            { return **(**float32)(unsafe.Pointer(&p)) }
func ptrToFloat64(p uintptr) float64            { return **(**float64)(unsafe.Pointer(&p)) }
func ptrToBool(p uintptr) bool                  { return **(**bool)(unsafe.Pointer(&p)) }
func ptrToBytes(p uintptr) []byte               { return **(**[]byte)(unsafe.Pointer(&p)) }
func ptrToNumber(p uintptr) json.Number         { return **(**json.Number)(unsafe.Pointer(&p)) }
func ptrToString(p uintptr) string              { return **(**string)(unsafe.Pointer(&p)) }
func ptrToSlice(p uintptr) *runtime.SliceHeader { return *(**runtime.SliceHeader)(unsafe.Pointer(&p)) }
func ptrToPtr(p uintptr) uintptr {
	return uintptr(**(**unsafe.Pointer)(unsafe.Pointer(&p)))
}
func ptrToNPtr(p uintptr, ptrNum uint8) uintptr {
	for i := uint8(0); i < ptrNum; i++ {
		if p == 0 {
			return 0
		}
		p = ptrToPtr(p)
	}
	return p
}

func ptrToUnsafePtr(p uintptr) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&p))
}
func ptrToInterface(code *encoder.Opcode, p uintptr) interface{} {
	return *(*interface{})(unsafe.Pointer(&emptyInterface{
		typ: code.Type,
		ptr: *(*unsafe.Pointer)(unsafe.Pointer(&p)),
	}))
}

func appendInt(ctx *encoder.RuntimeContext, b []byte, p uintptr, code *encoder.Opcode) []byte {
	format := ctx.Option.ColorScheme.Int
	b = append(b, format.Header...)
	b = encoder.AppendInt(ctx, b, p, code)
	return append(b, format.Footer...)
}

func appendUint(ctx *encoder.RuntimeContext, b []byte, p uintptr, code *encoder.Opcode) []byte {
	format := ctx.Option.ColorScheme.Uint
	b = append(b, format.Header...)
	b = encoder.AppendUint(ctx, b, p, code)
	return append(b, format.Footer...)
}

func appendFloat32(ctx *encoder.RuntimeContext, b []byte, v float32) []byte {
	format := ctx.Option.ColorScheme.Float
	b = append(b, format.Header...)
	b = encoder.AppendFloat32(ctx, b, v)
	return append(b, format.Footer...)
}

func appendFloat64(ctx *encoder.RuntimeContext, b []byte, v float64) []byte {
	format := ctx.Option.ColorScheme.Float
	b = append(b, format.Header...)
	b = encoder.AppendFloat64(ctx, b, v)
	return append(b, format.Footer...)
}

func appendString(ctx *encoder.RuntimeContext, b []byte, v string) []byte {
	format := ctx.Option.ColorScheme.String
	b = append(b, format.Header...)
	b = encoder.AppendString(ctx, b, v)
	return append(b, format.Footer...)
}

func appendByteSlice(ctx *encoder.RuntimeContext, b []byte, src []byte) []byte {
	format := ctx.Option.ColorScheme.Binary
	b = append(b, format.Header...)
	b = encoder.AppendByteSlice(ctx, b, src)
	return append(b, format.Footer...)
}

func appendNumber(ctx *encoder.RuntimeContext, b []byte, n json.Number) ([]byte, error) {
	format := ctx.Option.ColorScheme.Int
	b = append(b, format.Header...)
	bb, err := encoder.AppendNumber(ctx, b, n)
	if err != nil {
		return nil, err
	}
	return append(bb, format.Footer...), nil
}

func appendBool(ctx *encoder.RuntimeContext, b []byte, v bool) []byte {
	format := ctx.Option.ColorScheme.Bool
	b = append(b, format.Header...)
	if v {
		b = append(b, "true"...)
	} else {
		b = append(b, "false"...)
	}
	return append(b, format.Footer...)
}

func appendNull(ctx *encoder.RuntimeContext, b []byte) []byte {
	format := ctx.Option.ColorScheme.Null
	b = append(b, format.Header...)
	b = append(b, "null"...)
	return append(b, format.Footer...)
}

func appendComma(_ *encoder.RuntimeContext, b []byte) []byte {
	return append(b, ',', '\n')
}

func appendNullComma(ctx *encoder.RuntimeContext, b []byte) []byte {
	format := ctx.Option.ColorScheme.Null
	b = append(b, format.Header...)
	b = append(b, "null"...)
	return append(append(b, format.Footer...), ',', '\n')
}

func appendColon(_ *encoder.RuntimeContext, b []byte) []byte {
	return append(b[:len(b)-2], ':', ' ')
}

func appendMapKeyValue(ctx *encoder.RuntimeContext, code *encoder.Opcode, b, key, value []byte) []byte {
	b = appendIndent(ctx, b, code.Indent+1)
	b = append(b, key...)
	b[len(b)-2] = ':'
	b[len(b)-1] = ' '
	return append(b, value...)
}

func appendMapEnd(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte) []byte {
	b = b[:len(b)-2]
	b = append(b, '\n')
	b = appendIndent(ctx, b, code.Indent)
	return append(b, '}', ',', '\n')
}

func appendArrayHead(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte) []byte {
	b = append(b, '[', '\n')
	return appendIndent(ctx, b, code.Indent+1)
}

func appendArrayEnd(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte) []byte {
	b = b[:len(b)-2]
	b = append(b, '\n')
	b = appendIndent(ctx, b, code.Indent)
	return append(b, ']', ',', '\n')
}

func appendEmptyArray(_ *encoder.RuntimeContext, b []byte) []byte {
	return append(b, '[', ']', ',', '\n')
}

func appendEmptyObject(_ *encoder.RuntimeContext, b []byte) []byte {
	return append(b, '{', '}', ',', '\n')
}

func appendObjectEnd(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte) []byte {
	last := len(b) - 1
	// replace comma to newline
	b[last-1] = '\n'
	b = appendIndent(ctx, b[:last], code.Indent)
	return append(b, '}', ',', '\n')
}

func appendMarshalJSON(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte, v interface{}) ([]byte, error) {
	return encoder.AppendMarshalJSONIndent(ctx, code, b, v)
}

func appendMarshalText(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte, v interface{}) ([]byte, error) {
	format := ctx.Option.ColorScheme.String
	b = append(b, format.Header...)
	bb, err := encoder.AppendMarshalTextIndent(ctx, code, b, v)
	if err != nil {
		return nil, err
	}
	return append(bb, format.Footer...), nil
}

func appendStructHead(_ *encoder.RuntimeContext, b []byte) []byte {
	return append(b, '{', '\n')
}

func appendStructKey(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte) []byte {
	b = appendIndent(ctx, b, code.Indent)

	format := ctx.Option.ColorScheme.ObjectKey
	b = append(b, format.Header...)
	b = append(b, code.Key[:len(code.Key)-1]...)
	b = append(b, format.Footer...)

	return append(b, ':', ' ')
}

func appendStructEndSkipLast(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte) []byte {
	last := len(b) - 1
	if b[last-1] == '{' {
		b[last] = '}'
	} else {
		if b[last] == '\n' {
			// to remove ',' and '\n' characters
			b = b[:len(b)-2]
		}
		b = append(b, '\n')
		b = appendIndent(ctx, b, code.Indent-1)
		b = append(b, '}')
	}
	return appendComma(ctx, b)
}

func restoreIndent(ctx *encoder.RuntimeContext, code *encoder.Opcode, ctxptr uintptr) {
	ctx.BaseIndent = uint32(load(ctxptr, code.Length))
}

func storeIndent(ctxptr uintptr, code *encoder.Opcode, indent uintptr) {
	store(ctxptr, code.Length, indent)
}

func appendArrayElemIndent(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte) []byte {
	return appendIndent(ctx, b, code.Indent+1)
}

func appendMapKeyIndent(ctx *encoder.RuntimeContext, code *encoder.Opcode, b []byte) []byte {
	return appendIndent(ctx, b, code.Indent)
}
