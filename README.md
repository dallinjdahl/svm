# svm
simple (or stack) virtual machine with assembler [golang]

## VM

### Structure

svm is a stack based virtual machine with a single register.  This
register is used in autoincrement instructions as well as an IO port
for communication with the outside world, through the `ii` instruction.

### Instructions

The vm has 5 bit opcodes packed into 32 bit words.
Instructions that require an address use the remainder of the word
for the address.  Instructions that require an address include the
`jump`, `call`, `next`, `if` and `-if`.  Instructions of note
include the autoincrement instructions, `@+`, `!+`, `@p` and `!p`.  These
fetch and store through the `a` register and instruction pointer
respectively, and then increment their addresses. This enables
local variables as illustrated in the colorforth link below.

```
00 .		08 -if	10 ii	18 or 
01 ;		09 @+	11 *	19 drop
02 ex		0a @	12 /mod	1a dup
03 jump		0b @p	13 2*	1b over
04 call		0c !+	14 2/	1c a
05 unext	0d !	15 -	1d a!
06 next		0e !p	16 +	1e push
07 if		0f halt	17 and	1f pop
```

### IO

IO is accomplished through the `ii` instruction.  It takes a device id in the `a` register,
and any arguments on the stack.  The devices are as follows:

ID	|	Name	|	Argument	|	Description
---	|	---	|	---	|	---	|
0	|	VM		|	0			|	Pushes the amount of memory available to the stack
0	|	VM		|	1			|	Pushes the number of devices supported to the stack
1	|	cout	|	'c'			|	Outputs the low 8 bits as an ascii character
2	|	block	|	0, addr		|	Read a block from the block file into addr+1 (block # is at addr)
2	|	block	|	1, addr		|	Write a block to the block file from addr+1 (block # is at addr)
3	|	cin		|				| 	Pushes the next byte of input from the keyboard to the stack

#### Blocks
Instead of supporting generic file access, svm supports access to a block file, being treated
as a block storage device.  The contents of memory are mapped into the first blocks of the file,
which allow for reloading the vm.  In the future, access to files provided on the command line
will also be provided, analogous to disk drives in a real machine, still via the block mechanism.
Blocks are typically cached in memory when used, which allows for memory mapped disk access.

## Assembler

The assembler for svm is called sas.  sas is a simple 2 pass assembler with support for lables and comments,
inspired by `muri`, the assembler for retro forth.
Each line in a sas file begins with a directive: `i`, `d`, `r`, `:`, and `/`.  It's called with an input
file and an output file name.  It can compile include files to be compiled into the svm binary, or binary
files to be used as block files loaded by svm, as specified by the `-b` flag (for binary).

`i` lines continue with a set of 2 letter opcodes from the table below.  If an instruction requires an
address, this is specified immediately following the instruction. Note that if there is not room in
the word for a particular address, the assembler will complain and not compile it.  So for example,
at the beginning of every file is the line:
```
ijumain
```
which compiles a jump instruction and a jump to the label main. Lines not ending in a label referencing instruction
(`ju`, `ca`, `nx`, `if`, or `-i`) must be 6 instructions long, padded with nops (`..`) if needed.  

`d` lines specify an unsigned numeric value in either decimal or hexadecimal.  If it's in hexadecimal,
the first character after the d should be x, as in:
```
dxdeadbeef
```

`r` lines compile the addresses of labels, so to compile the address of main would be
```
rmain
```

`:` lines declare labels, so to specify main would be
```
:main
```

`/` lines are comments, and are ignored by the assembler.

### Opcodes

```
00 ..	08 -i	10 ii	18 or 
01 ;;	09 @+	11 **	19 dr
02 ex	0a @a	12 /m	1a du
03 ju	0b @p	13 2*	1b ov
04 ca	0c !+	14 2/	1c a@
05 un	0d !a	15 --	1d a!
06 nx	0e !p	16 ++	1e pu
07 if	0f ha	17 an	1f po
```

## References
These documents helped flesh out the instruction set and may provide added insight into the operation
of the vm.

https://colorforth.github.io/inst.htm

http://retroforth.org/docs/The_Ngaro_Virtual_Machine.html
