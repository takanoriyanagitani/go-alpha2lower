package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"

	lwa0 "github.com/takanoriyanagitani/go-alpha2lower/conv/wasm/wazero"
	"github.com/tetratelabs/wazero"
)

func sub(
	ctx context.Context,
	wasmLoc string,
	wasmByteMax int64,
	wasmLimitPage uint32,
) error {
	file, err := os.Open(wasmLoc) //nolint:gosec
	if nil != err {
		return err
	}
	defer file.Close() //nolint:errcheck

	limited := &io.LimitedReader{
		R: file,
		N: wasmByteMax,
	}

	wbytes, err := io.ReadAll(limited)
	if nil != err {
		return err
	}

	var rcfg wazero.RuntimeConfig = wazero.
		NewRuntimeConfig().
		WithMemoryLimitPages(wasmLimitPage)
	var mcfg wazero.ModuleConfig = wazero.NewModuleConfig()

	conv, err := lwa0.
		WasmBytes(wbytes).
		ToConverter(
			ctx,
			rcfg,
			mcfg,
		)
	if err != nil {
		return err
	}
	defer conv.Close(ctx) //nolint:errcheck

	return conv.LowerStdinToStdout(ctx)
}

func main() {
	var wasmLoc string
	var wasmByteMax int64
	var wasmPageMax uint

	flag.StringVar(&wasmLoc, "wasm-path", "", "wasm path")
	flag.Int64Var(&wasmByteMax, "wasm-size-max", 131072, "wasm size max")
	flag.UintVar(&wasmPageMax, "wasm-page-max", 16, "wasm page max")

	flag.Parse()

	if "" == wasmLoc {
		flag.Usage()
		os.Exit(1)
	}

	err := sub(
		context.Background(),
		wasmLoc,
		wasmByteMax,
		uint32(wasmPageMax), //nolint:gosec
	)
	if nil != err {
		log.Printf("%v\n", err)
	}
}
