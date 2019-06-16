package main

import (
	"debug/elf"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	ptGnuStack = 0x6474e551
)

func patch(path string, setStackSize bool, stackSize uint64) (uint64, error) {
	openFlag := os.O_RDONLY
	if setStackSize {
		openFlag = os.O_RDWR
	}
	f, err := os.OpenFile(path, openFlag, 0755)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	var fh elf.FileHeader
	var ident [16]uint8
	if _, err = f.Read(ident[0:]); err != nil {
		return 0, err
	}
	if ident[0] != '\x7f' || ident[1] != 'E' || ident[2] != 'L' || ident[3] != 'F' {
		return 0, fmt.Errorf("bad magic number %v", ident[0:4])
	}

	fh.Class = elf.Class(ident[elf.EI_CLASS])
	switch fh.Class {
	case elf.ELFCLASS32:
	case elf.ELFCLASS64:
	default:
		return 0, fmt.Errorf("unknown ELF class %s", fh.Class)
	}

	fh.Data = elf.Data(ident[elf.EI_DATA])
	switch fh.Data {
	case elf.ELFDATA2LSB:
		fh.ByteOrder = binary.LittleEndian
	case elf.ELFDATA2MSB:
		fh.ByteOrder = binary.BigEndian
	default:
		return 0, fmt.Errorf("unknown ELF data encoding %s", fh.Data)
	}

	fh.Version = elf.Version(ident[elf.EI_VERSION])
	if fh.Version != elf.EV_CURRENT {
		return 0, fmt.Errorf("unknown ELF version %s", fh.Version)
	}

	f.Seek(0, io.SeekStart)

	switch fh.Class {
	case elf.ELFCLASS32:
		eh := &elf.Header32{}
		if err := binary.Read(f, fh.ByteOrder, eh); err != nil {
			return 0, err
		}
		phoff := int64(eh.Phoff)
		phentsize := int64(eh.Phentsize)
		phnum := int(eh.Phnum)
		for i := 0; i < phnum; i++ {
			off := phoff + int64(i)*phentsize
			f.Seek(off, io.SeekStart)
			ph := &elf.Prog32{}
			if err := binary.Read(f, fh.ByteOrder, ph); err != nil {
				return 0, err
			}
			if ph.Type == ptGnuStack {
				if setStackSize {
					ph.Memsz = uint32(stackSize)
					f.Seek(off, io.SeekStart)
					if err := binary.Write(f, fh.ByteOrder, ph); err != nil {
						return 0, err
					}
				}
				return uint64(ph.Memsz), nil
			}
		}
	case elf.ELFCLASS64:
		eh := &elf.Header64{}
		if err := binary.Read(f, fh.ByteOrder, eh); err != nil {
			return 0, err
		}
		phoff := int64(eh.Phoff)
		phentsize := int64(eh.Phentsize)
		phnum := int(eh.Phnum)
		for i := 0; i < phnum; i++ {
			off := phoff + int64(i)*phentsize
			f.Seek(off, io.SeekStart)
			ph := &elf.Prog64{}
			if err := binary.Read(f, fh.ByteOrder, ph); err != nil {
				return 0, err
			}
			if ph.Type == ptGnuStack {
				if setStackSize {
					ph.Memsz = stackSize
					f.Seek(off, io.SeekStart)
					if err := binary.Write(f, fh.ByteOrder, ph); err != nil {
						return 0, err
					}
				}
				return ph.Memsz, nil
			}
		}
	}

	return 0, fmt.Errorf("PT_GNU_STACK program header not found")
}

func main() {
	var (
		setStackSize bool
		stackSize    uint64
	)

	flag.Uint64Var(&stackSize, "s", 0, "set default stack size, example: 0x800000 for 8MB")
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Usage:\n")
		fmt.Fprintf(out, "  %s [options] executables ...\n", os.Args[0])
		fmt.Fprintf(out, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Checks if each flag is given or not
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "s" {
			setStackSize = true
		}
	})

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	exitCode := 0
	for _, path := range flag.Args() {
		newSize, err := patch(path, setStackSize, stackSize)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", path, err)
			exitCode = 1
			continue
		}
		fmt.Printf("%s: stackSize: 0x%x\n", path, newSize)
	}

	os.Exit(exitCode)
}
