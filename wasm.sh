#!/bin/sh

iname=./alpha2lower.wat
oname=./alpha2lower.wasm

wat2wasm \
	"${iname}" \
	-o "${oname}" \
	--enable-relaxed-simd \
	--enable-tail-call || exec sh -c '
		echo unable to compile.
		exit 1
	'

ls -l "${oname}"
