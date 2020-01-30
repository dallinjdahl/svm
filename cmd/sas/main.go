package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

var line = 0
var labels = map[string]uint32{}
var ops = "..;;exjucaunnxif-i@+@a@p!+!a!phaii**/m2*2/--++anordrduova@a!pupo"
var opa []uint16
var inc = true;

func main() {
	var buf []byte

	opa = make([]uint16, 32)
	for i := 0; i < 32; i++ {
		opa[i] = uint16(ops[2*i+1]) + (uint16(ops[2*i]) << 8)
	}

	if len(os.Args) < 3 {
		print("Usage: sas <infile> <outfile> [-b]\n")
		os.Exit(1)
	}

	if len(os.Args) == 4 && os.Args[3] == "-b" {
		inc = false
	}

	file := os.Args[1]
	inf, err := os.Open(file)
	check(err)
	defer inf.Close()
	in := bufio.NewReader(inf)

	file = os.Args[2]
	outf, err := os.Create(file)
	check(err)
	out := bufio.NewWriter(outf)


	defer func() {
		cerr := outf.Close()
		check(cerr)
	}()

	if inc {
		_, err = out.WriteString("package main\n\nvar image = []uint32{")
		check(err)
	}

	label(in)

	_, err = inf.Seek(0, 0)
	check(err)

	in.Reset(inf)

	line = 0

	buf, err = in.ReadSlice(byte('\n'))
	first := true
	for err == nil {
		op, ok := parse(buf)
		if ok {
			write(out, op, first)
			first = false
		}
		buf, err = in.ReadSlice(byte('\n'))
		line++
	}
	if inc {
		_, err = out.WriteString("}\n")
		check(err)
	}
	out.Flush()
}

func label(in *bufio.Reader) {
	line := uint32(0)
	var err error
	var buf []byte
	for err == nil {
		buf, err = in.ReadSlice(byte('\n'))
		if len(buf) > 1 && buf[0] == ':' {
			labels[string(buf[1:])] = line
			continue
		}
		if len(buf) > 1 {
			line++
		}
	}
}

func parse(buf []byte) (uint32, bool) {
	var op uint32

	if len(buf) < 2 {
		return 0, false
	}

	switch(buf[0]) {
	case 'i':
		for i := 0; i < 6; i++ {
			var temp uint16
			temp = uint16(buf[2*i+2])+(uint16(buf[2*i+1])<<8)
			fmt.Printf("opcode: %d:%s %x\n", line, string(buf[2*i+1:2*i+3]), temp)
			j := 0
			for ;j < 32 && temp != opa[j]; j++ { }
			if(j == 32) {
				fmt.Printf("Unrecognized opcode: %d:%s\n", line, string(buf[2*i+1:2*i+3]))
				os.Exit(1)
			}
			
			fmt.Printf("op %x\n", j)
			op <<= 5
			op |= uint32(j)

			switch temp {
			case opa[3]:
				fallthrough
			case opa[4]:
				fallthrough
			case opa[6]:
				fallthrough
			case opa[7]:
				fallthrough
			case opa[8]:
				size := (((5 - i)*5) + 2)
				l := labels[string(buf[2*i+3:])]
				if l >= 1 << size{
					fmt.Printf("Jump too long: %d%s\n", line, buf[2*i+3:])
					os.Exit(1)
				}
				op <<= size
				op |= l
				return op, true
			}
		}
		op <<= 2
		return op, true
	case 'd':
		var err error
		var big int64
		if buf[1] == 'x' {
			big, err = strconv.ParseInt(string(buf[2:len(buf)-1]), 16, 32)
		} else {
			big, err = strconv.ParseInt(string(buf[1:len(buf)-1]), 10, 32)
		}
		check(err)
		op = uint32(big)
		return op, true
	case 's':
		op += uint32(buf[1])
		op += uint32(buf[2]) << 8
		op += uint32(buf[3]) << 16
		op += uint32(buf[4]) << 24
		return op, true
	case 'r':
		op = labels[string(buf[1:])]
		fmt.Printf("ref: %s:%d %v\n", string(buf[1:]), op, labels)
		return op, true
	case ':':
		return 0, false
	default:
		fmt.Printf("Unrecognized directive: %d:%c\n", line, buf[0])
		os.Exit(1)
	}
	return 0, false
}

func write(out *bufio.Writer, op uint32, first bool) {
	var err error
	switch {
	case inc && first:
		_, err = fmt.Fprintf(out, "0x%x", op)
		check(err)
		return
	case inc:
		_, err = fmt.Fprintf(out, ", 0x%x", op)
		check(err)
		return
	}
	err = binary.Write(out, binary.LittleEndian, op)
	check(err)
}
