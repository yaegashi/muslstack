# muslstack

## Introduction

muslstack is a tiny utility to binary-patch ELF executable files
to expand their default thread stack size when running with
[musl libc](https://www.musl-libc.org).

musl allocates memory of very small size for thread stacks.
[musl documentation](https://wiki.musl-libc.org/functional-differences-from-glibc.html#Thread-stack-size)
says it's only 80KB ([now increased to 128KB](http://git.musl-libc.org/cgit/musl/commit/?id=c0058ab465e950c2c3302d2b62e21cc0b494224b)) for default,
which is rather small compared to 2-10MB of glibc.
And that's been said to be the major reason
for stack-consuming executable to crash (segmentation fault)
in musl libc environment.

In September 2018 [musl introduced the feature](http://git.musl-libc.org/cgit/musl/commit/?id=7b3348a98c139b4b4238384e52d4b0eb237e4833) to allow users
to mitigate the stack size limitation without modifying source code.
It takes the default stack size at runtime
from ELF program header of the executable,
namely from the memory size value of PT_GNU_STACK header.

It's originally intended that the ELF object linker should set the value
at building time with a special linker flag like `-Wl,-z,stack-size=N`.
Using muslstack, you can modify it on prebuilt executables afterwards,
without any source code modification or rebuild.

## Alpine Linux compatibility

muslstack could be an effective remedy for stack overflow issues
of various executables to be run in [Alpine Linux](https://alpinelinux.org/)
containers, because Alpine Linux utilizes musl as its libc.

However, the feature muslstack rely on is available only in recent musl.
It's not incorporated in the stable release of Alpine Linux
as of June 2019.
You would need to begin your Dockerfile with `FROM alpine:edge`.

## Usage

```console
$ go get github.com/yaegashi/muslstack
$ muslstack
Usage:
  muslstack [options] executables ...
Options:
  -s uint
        set default stack size, example: 0x800000 for 8MB
```

You can see all program headers of ELF executable using `objdump -p`:

```console
$ echo 'package main; func main() {}' >main.go
$ go build main.go
$ objdump -p main

main:     file format elf64-x86-64

Program Header:
    PHDR off    0x0000000000000040 vaddr 0x0000000000400040 paddr 0x0000000000400040 align 2**12
         filesz 0x0000000000000188 memsz 0x0000000000000188 flags r--
    NOTE off    0x0000000000000f9c vaddr 0x0000000000400f9c paddr 0x0000000000400f9c align 2**2
         filesz 0x0000000000000064 memsz 0x0000000000000064 flags r--
    LOAD off    0x0000000000000000 vaddr 0x0000000000400000 paddr 0x0000000000400000 align 2**12
         filesz 0x000000000004f990 memsz 0x000000000004f990 flags r-x
    LOAD off    0x0000000000050000 vaddr 0x0000000000450000 paddr 0x0000000000450000 align 2**12
         filesz 0x000000000006e493 memsz 0x000000000006e493 flags r--
    LOAD off    0x00000000000bf000 vaddr 0x00000000004bf000 paddr 0x00000000004bf000 align 2**12
         filesz 0x00000000000028c0 memsz 0x00000000000203b8 flags rw-
   STACK off    0x0000000000000000 vaddr 0x0000000000000000 paddr 0x0000000000000000 align 2**3
         filesz 0x0000000000000000 memsz 0x0000000000000000 flags rw-
0x65041580 off    0x0000000000000000 vaddr 0x0000000000000000 paddr 0x0000000000000000 align 2**3
         filesz 0x0000000000000000 memsz 0x0000000000000000 flags --- 2a00
```

You can modify or examine memsz value in STACK header
(which is actually PT_GNU_STACK) using muslstack.

```console
$ muslstack -s 0x800000 main
main: stackSize: 0x800000
$ objdump -p main | grep -A1 STACK
   STACK off    0x0000000000000000 vaddr 0x0000000000000000 paddr 0x0000000000000000 align 2**3
         filesz 0x0000000000000000 memsz 0x0000000000800000 flags rw-
$ muslstack main
main: stackSize: 0x800000
```

You can find more test cases in [tests](./tests) folder.

## License

MIT