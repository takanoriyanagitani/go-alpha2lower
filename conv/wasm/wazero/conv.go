package conv

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/tetratelabs/wazero"
	wa "github.com/tetratelabs/wazero/api"
)

var (
	ErrNilMem  error = errors.New("nil memory")
	ErrNilFunc error = errors.New("nil function")

	ErrUnableToRead  error = errors.New("unable to read the converted")
	ErrUnableToWrite error = errors.New("unable to write the original")
)

type WasmFn struct{ wa.Function }

func (f WasmFn) Call(ctx context.Context) error {
	_, err := f.Function.Call(ctx)
	return err
}

type WasmMem struct{ wa.Memory }

func (m WasmMem) ReadPage(offset uint32) ([]byte, bool) {
	return m.Memory.Read(offset, 65536)
}

func (m WasmMem) ReadConverted() ([]byte, error) {
	conv, ok := m.ReadPage(65536)
	if !ok {
		return nil, ErrUnableToRead
	}
	return conv, nil
}

type WasmMod struct{ wa.Module }

func (m WasmMod) Close(ctx context.Context) error {
	if nil == m.Module {
		return nil
	}
	return m.Module.Close(ctx)
}

func (m WasmMod) Memory() (WasmMem, error) {
	var mem wa.Memory = m.Module.Memory()
	if nil == mem {
		return WasmMem{}, ErrNilMem
	}
	return WasmMem{Memory: mem}, nil
}

func (m WasmMod) GetFunction(name string) (WasmFn, error) {
	var fnc wa.Function = m.Module.ExportedFunction(name)
	if nil == fnc {
		return WasmFn{}, ErrNilFunc
	}
	return WasmFn{Function: fnc}, nil
}

func (m WasmMod) GetConverter() (WasmFn, error) {
	return m.GetFunction("lowerpage")
}

type Compiled struct{ wazero.CompiledModule }

func (c Compiled) Close(ctx context.Context) error {
	if nil == c.CompiledModule {
		return nil
	}
	return c.CompiledModule.Close(ctx)
}

type WasmRuntime struct{ wazero.Runtime }

func (r WasmRuntime) Close(ctx context.Context) error {
	if nil == r.Runtime {
		return nil
	}
	return r.Runtime.Close(ctx)
}

func (r WasmRuntime) Compile(
	ctx context.Context,
	wasm []byte,
) (Compiled, error) {
	cmod, err := r.Runtime.CompileModule(ctx, wasm)
	return Compiled{CompiledModule: cmod}, err
}

func (r WasmRuntime) Instantiate(
	ctx context.Context,
	compiled Compiled,
	cfg wazero.ModuleConfig,
) (WasmMod, error) {
	amod, err := r.Runtime.InstantiateModule(
		ctx,
		compiled.CompiledModule,
		cfg,
	)

	return WasmMod{Module: amod}, err
}

type WasmConfig struct{ wazero.RuntimeConfig }

type Converter struct {
	WasmRuntime
	Compiled
	WasmMod
	WasmMem
	WasmFn
}

func (c Converter) Close(ctx context.Context) error {
	return errors.Join(
		c.WasmMod.Close(ctx),
		c.Compiled.Close(ctx),
		c.WasmRuntime.Close(ctx),
	)
}

//nolint:cyclop
func (c Converter) Lower(
	ctx context.Context,
	rdr io.Reader,
	wtr io.Writer,
) error {
	var ibuf [65536]byte
	for {
		cnt, err := io.ReadFull(rdr, ibuf[:])
		if 0 < cnt {
			woriginal, bok := c.WasmMem.Memory.Read(0, 65536)
			if !bok {
				return ErrUnableToRead
			}
			copy(woriginal, ibuf[:cnt])

			err = c.WasmFn.Call(ctx)
			if nil != err {
				return err
			}

			var clen int = cnt & 0xffff_ffff
			var ulen uint32 = uint32(clen) //nolint:gosec
			wlower, bok := c.WasmMem.Memory.Read(65536, ulen)
			if !bok {
				return ErrUnableToWrite
			}

			_, err = wtr.Write(wlower)
			if nil != err {
				return err
			}
		}

		switch {
		case nil == err:
			continue
		case io.EOF == err: //nolint:errorlint
			return nil
		case errors.Is(err, io.ErrUnexpectedEOF):
			return nil
		default:
			return err
		}
	}
}

func (c Converter) LowerStdinToStdout(ctx context.Context) error {
	return c.Lower(
		ctx,
		os.Stdin,
		os.Stdout,
	)
}

type WasmBytes []byte

func (b WasmBytes) ToConverter(
	ctx context.Context,
	rcfg wazero.RuntimeConfig,
	mcfg wazero.ModuleConfig,
) (Converter, error) {
	var rtm wazero.Runtime = wazero.NewRuntimeWithConfig(
		ctx,
		rcfg,
	)
	var conv Converter
	conv.WasmRuntime = WasmRuntime{Runtime: rtm}

	compiled, err := rtm.CompileModule(ctx, b)
	if nil != err {
		return conv, err
	}
	conv.Compiled = Compiled{CompiledModule: compiled}

	instance, err := rtm.InstantiateModule(
		ctx,
		conv.Compiled.CompiledModule,
		mcfg,
	)
	if nil != err {
		return conv, err
	}
	conv.WasmMod = WasmMod{Module: instance}

	lower, err := conv.WasmMod.GetConverter()
	if nil != err {
		return conv, err
	}
	conv.WasmFn = lower

	wmem, err := conv.WasmMod.Memory()
	if nil != err {
		return conv, err
	}
	conv.WasmMem = wmem

	return conv, nil
}
