#!/bin/sh

wpath=./alpha2lower.wasm

input(){
	local rate
	rate=$1

	dd \
		if=/dev/zero \
		bs=1048576 \
		status=none \
		count=${rate:-1024} |
		openssl \
			enc \
			-nosalt \
			-aes-256-ctr \
			-pass pass:non-secure-random-data
}

bench(){
	input |
		./cmd/alpha2lower/alpha2lower -wasm-path "${wpath}" |
		dd \
			of=/dev/null \
			bs=1048576 \
			status=progress
}

check_tr(){
	input 128 |
		LC_ALL=C tr '[:upper:]' '[:lower:]' |
		shasum -a 256
}

check_simd(){
	input 128 |
		./cmd/alpha2lower/alpha2lower -wasm-path "${wpath}" |
		shasum -a 256
}

#check_tr
#check_simd
time bench
