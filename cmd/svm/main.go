package main

import (
	"encoding/binary"
	"io"
	"log"
	"fmt"
	"os"
)

const (
	KB uint16 = 256
	MEMORY uint32 = 8 * uint32(KB)
	DSTACK uint8 = 32
	RSTACK uint8 = 32
	omask uint32 = 31 << 2
)

var(
	ram [MEMORY]uint32
	dstack [DSTACK]uint32
	rstack [RSTACK]uint32
	dp uint8
	rp uint8
	ip uint32 
	a uint32
	i uint32
	slot uint8
	ops [32]func() = [32]func(){nop, ret, ex, jump, call, unext, next, pif,
					pnif, fetchinca, fetcha, fetchip, storeinca, storea, storeip, halt,
					iinteract, mult, dmod, lshift, rshift, neg, add, and,
					or, drop, dup, over, afetch, astore, push, pop}
	ios []func() = []func(){nil, cout, block, cin}
	file *os.File
)

func vm() {
	dp--
	op := dstack[dp]
	switch op {
	case 0:
		dstack[dp] = MEMORY
	case 1:
		dstack[dp] = uint32(len(ios))
	default:
		fmt.Printf("Illegal io vm op: %d", op)
		os.Exit(1)
	}
	dp++
}

func cout() {
	dp--
	fmt.Printf("%c", uint8(dstack[dp]))
}

func block() {
	dp--
	op := dstack[dp]
	dp--
	b := ram[dstack[dp]]
	data := ram[dstack[dp]+1:dstack[dp]+(1*uint32(KB))+1]
	file.Seek(int64(b * uint32(KB)), 0)
	switch op {
	case 0:
		binary.Read(file, binary.LittleEndian, data)
	case 1:
		binary.Write(file, binary.LittleEndian, data) 
	default:
		fmt.Printf("Block op error: %d", op)
	}
}

func cin() {
	c := make([]byte, 1)
	os.Stdin.Read(c)
	dstack[dp] = uint32(c[0])
	dp++
}


func readin(r io.Reader) {
	var buf [4]byte
	j := uint32(0)

	var err error
	
	io.ReadFull(r, buf[:])
	for err == nil && j < MEMORY {
		ram[j] = binary.LittleEndian.Uint32(buf[:])
		j++
		io.ReadFull(r, buf[:])
	}
}

func stat() {
	fmt.Printf("i0x%x ip%d a0x%x d%d r%d\n", i, ip, a, dp, rp)
	fmt.Printf("% d\n", dstack[:dp])
}

func main() {
	ios[0] = vm
	if _, err := os.Stat(".block"); err != nil {
		file, err = os.Create(".block")
		binary.Write(file, binary.LittleEndian, image)
		file.Close()
	}
	file, err := os.OpenFile(".block", os.O_RDWR, 0655)

	if err != nil {
		log.Fatal(err)
	}
	
	readin(file)

	for {
		i = ram[ip]
		ip++
		execute()
	}
}

func execute() {
	var op uint32
	for slot = 0; slot < 6; {
		op = omask << ((5-slot)*5)
		op &= i
		op >>= 2+(5-slot)*5
//		fmt.Printf("oX%x\n", op)
		ops[op]()
//		stat()
	}
}

func nop() {slot++}

func ret() {
	rp--
	ip = rstack[rp]
	slot = 6
}

func ex() {
	x := ip
	ip = rstack[rp-1]
	rstack[rp-1] = x
	slot = 6
}

func jump() {
	addr := i
	addr <<= 5*(slot+1)
	addr >>= 5*(slot+1)
	ip = addr
	slot = 6
}

func call() {
	addr := i
	addr <<= 5*(slot+1)
	addr >>= 5*(slot+1)
	rstack[rp] = ip
	rp++
	ip = addr
	slot = 6
}

func unext() {
	if rstack[rp-1] == 0 {
		rp--
		slot++
		return
	}
	rstack[rp-1]--
	slot = 0
}

func next() {
	if rstack[rp-1] == 0 {
		rp--
		slot++
		return
	}
	addr := i
	addr <<= 5*(slot+1)
	addr >>= 5*(slot+1)
	ip = addr
	rstack[rp-1]--
	slot = 6
}

func pif() {
	if dstack[dp-1] != 0 {
		slot = 6
		return
	}
	addr := i
	addr <<= 5*(slot+1)
	addr >>= 5*(slot+1)
	ip = addr
	slot = 6
}

func pnif() {
	if dstack[dp-1] < 0 {
		slot = 6
		return
	}
	addr := i
	addr <<= 5*(slot+1)
	addr >>= 5*(slot+1)
	ip = addr
	slot = 6
}

func fetchinca() {
	dstack[dp] = ram[a]
	dp++
	a++
	slot++
}

func fetcha() {
	dstack[dp] = ram[a]
	dp++
	slot++
}

func fetchip() {
	dstack[dp] = ram[ip]
	dp++
	ip++
	slot++
}

func storeinca() {
	dp--
	ram[a] = dstack[dp]
	a++
	slot++
}

func storea() {
	dp--
	ram[a] = dstack[dp]
	slot++
}

func storeip() {
	dp--
	ram[ip] = dstack[dp]
	ip++
	slot++
}

func halt() {
	os.Exit(0)
}


func iinteract() {
	ios[a]()
	slot++
}

func mult() {
	dp--
	dstack[dp-1] *= dstack[dp]
	slot++
}

func dmod() {
	d := dstack[dp-2] / dstack[dp-1]
	dstack[dp-2] %= dstack[dp-1]
	dstack[dp-1] = d
	slot++
}

func lshift() {
	dstack[dp-1] <<= 1
	slot++
}

func rshift() {
	dstack[dp-1] = uint32(int32(dstack[dp-1]) >> 1)
	slot++
}

func neg() {
	dstack[dp-1] = ^dstack[dp-1]
	slot++
}

func add() {
	dp--
	dstack[dp-1] += dstack[dp]
	slot++
}

func and() {
	dp--
	dstack[dp-1] &= dstack[dp]
	slot++
}

func or() {
	dp--
	dstack[dp-1] ^= dstack[dp]
	slot++
}

func drop() {
	dp--
	slot++
}

func dup() {
	dstack[dp] = dstack[dp-1]
	dp++
	slot++
}

func over() {
	dstack[dp] = dstack[dp-2]
	dp++
	slot++
}

func afetch() {
	dstack[dp] = a
	dp++
	slot++
}

func astore() {
	dp--
	a = dstack[dp]
	slot++
}

func push() {
	dp--
	rstack[rp] = dstack[dp]
	rp++
	slot++
}

func pop() {
	rp--
	dstack[dp] = rstack[rp]
	dp++
	slot++
}
	
