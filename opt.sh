iwasi=./alpha2lower.wasm

wasm-opt \
	-O3 \
	-o opt.wasm \
	--enable-simd \
	--enable-relaxed-simd \
	--enable-bulk-memory \
	--enable-nontrapping-float-to-int \
	"${iwasi}"
