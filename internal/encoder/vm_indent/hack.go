package vm_indent

import (
	// HACK: compile order
	// `vm`, `vm_indent`, `vm_color`, `vm_color_indent` packages uses a lot of memory to compile,
	// so forcibly make dependencies and avoid compiling in concurrent.
	// dependency order: vm => vm_indent => vm_color => vm_color_indent
	_ "github.com/going/json/internal/encoder/vm_color"
)
