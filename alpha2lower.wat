(module

  ;; bytes     0 -  65535: original bytes
  ;; bytes 65536 - 131071: lowered  bytes
  (memory (export "memory") 2)

  (func (export "lower64")  (param $original i64) (result i64)
    (local $loa v128)
    (local $hiz v128)
    (local $off v128)
    (local $z v128)
    (local $is_upper v128)
    (local $add v128)

    (local $v v128)

    (local.set $loa (
      v128.const i32x4 0x41414141 0x41414141 0x41414141 0x41414141
    ))

    (local.set $hiz (
      v128.const i32x4 0x5a5a5a5a 0x5a5a5a5a 0x5a5a5a5a 0x5a5a5a5a
    ))

    (local.set $off (
      v128.const i32x4 0x20202020 0x20202020 0x20202020 0x20202020
    ))

    (local.set $z (
      v128.const i64x2 0 0
    ))

    local.get $z
    local.get $original
    i64x2.replace_lane 0
    local.set $v

    local.get $v
    local.get $loa
    i8x16.ge_u

    local.get $v
    local.get $hiz
    i8x16.le_u

    v128.and
    local.set $is_upper

    local.get $is_upper
    local.get $off
    v128.and
    local.set $add

    local.get $v
    local.get $add

    i8x16.add
    i64x2.extract_lane 0
  )

  (func (export "lowerpage")
    (local $loa v128)
    (local $hiz v128)
    (local $off v128)
    (local $z v128)
    (local $is_upper v128)
    (local $add v128)

    (local $original v128)
    (local $converted v128)

    (local $v v128)

    (local $iptr i32)
    (local $optr i32)

    (local.set $loa (
      v128.const i32x4 0x41414141 0x41414141 0x41414141 0x41414141
    ))

    (local.set $hiz (
      v128.const i32x4 0x5a5a5a5a 0x5a5a5a5a 0x5a5a5a5a 0x5a5a5a5a
    ))

    (local.set $off (
      v128.const i32x4 0x20202020 0x20202020 0x20202020 0x20202020
    ))

    (local.set $z (
      v128.const i64x2 0 0
    ))

    (local.set $iptr (i32.const 0))
    (local.set $optr (i32.const 65536))

    (block $exit
      (loop $process
        (br_if $exit (i32.ge_u (local.get $iptr) (i32.const 65536)))

        (local.set $original (v128.load (local.get $iptr)))

        local.get $original
        local.get $loa
        i8x16.ge_u

        local.get $original
        local.get $hiz
        i8x16.le_u

        ;; is upper
        v128.and

        ;; offset to add if upper else 0
        local.get $off
        v128.and

        ;; converted
        local.get $original
        i8x16.add
        local.set $converted

        (v128.store
          (local.get $optr)
          (local.get $converted)
        )

        (local.set $iptr (i32.add (local.get $iptr) (i32.const 16)))
        (local.set $optr (i32.add (local.get $optr) (i32.const 16)))
        (br $process)
      )
    )
  )

)
