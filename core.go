package emu

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

const (
	VECTOR_NMI   uint16 = 0xFFFA
	VECTOR_RESET uint16 = 0xFFFC
	VECTOR_IRQ   uint16 = 0xFFFE
)

type Core struct {
	// Main registers
	A uint8
	X uint8
	Y uint8

	// Other registers
	PC     uint16 // Program counter
	Phlags uint8  // Status flags
	SP     uint8  // Stack pointer

	memory []byte // Slice of loaded memory.  This is only main RAM.
	rom    []byte // ROM image.  Needs to be a multiple of 256.
	wram   []byte

	InstructionLimit uint64 // number of instructions to run
	testing          bool
	testDone         bool
	t                *testing.T
	ticks            uint64

	fullRW bool

	lastPC   uint16
	lastSame int
	lastReadAddr uint16
	checkStuck bool

	// VERY verbose output
	Debug bool
}

func NewRWCore(rom []byte, instrLimit uint64) (*Core, error) {
	if len(rom) != 0x10000 {
		return nil, fmt.Errorf("ROM must be exactly 64k (%X)", len(rom))
	}

	c := &Core{
		A:      0,
		X:      0,
		Y:      0,
		PC:     0,
		Phlags: 0,
		SP:     0,

		//memory: make([]byte, 0x1000), // no registers, no WRAM, no ROM
		rom: rom,

		InstructionLimit: instrLimit,

		fullRW:     true,
		checkStuck: true,
	}

	c.PC = c.ReadWord(VECTOR_RESET)
	return c, nil
}

func NewCore(rom []byte, wram bool, instrLimit uint64) (*Core, error) {
	if len(rom)%256 != 0 {
		return nil, fmt.Errorf("ROM is not divisible by 256: %d", len(rom))
	}

	c := &Core{
		A:      0,
		X:      0,
		Y:      0,
		PC:     0,
		Phlags: 0,
		SP:     0,

		memory: make([]byte, 0x1000), // no registers, no WRAM, no ROM
		rom:    rom,

		InstructionLimit: instrLimit,
	}

	if wram {
		c.wram = make([]byte, 0x2000)
	}

	if len(c.rom) == 0 {
		return nil, fmt.Errorf("No rom!")
	}
	fmt.Printf("Rom length: %X\n", len(c.rom))

	c.PC = c.ReadWord(VECTOR_RESET)

	return c, nil
}

// Read address.  This will read from API registers if needed.
func (c *Core) ReadByte(addr uint16) uint8 {
	c.lastReadAddr = addr
	if c.fullRW {
		return c.rom[addr]
	}

	if addr < 0x1000 {
		return c.memory[addr]
	}

	if addr >= 0x6000 && addr < 0x8000 {
		if c.wram != nil {
			// TODO: make sure this works with variable WRAM sizes (paging?)
			return c.wram[addr%uint16(len(c.wram))]
		}
		return 0
	}

	if addr >= 0x8000 {
		return c.rom[uint(addr)%uint(len(c.rom))]
	}

	// "Open bus"  always return zero.
	return 0
}

func (c *Core) ReadWord(addr uint16) uint16 {
	defer func() { c.lastReadAddr = addr }() // will this fire off correctly? idk
	return uint16(c.ReadByte(addr)) | (uint16(c.ReadByte(addr+1)) << 8)
}

// Write to an address.  This will delegate to API if needed.
func (c *Core) WriteByte(addr uint16, value byte) {
	if c.fullRW {
		c.rom[addr] = value
		return
	}

	if addr < 0x1000 {
		c.memory[addr] = value
	} else if addr < 0x6000 {
		// TODO: software registers
	} else if addr >= 0x6000 && addr < 0x8000 && c.wram != nil {
		c.wram[addr] = value
	}
}

func (c *Core) WriteInt(addr uint16, value uint8) {
	c.WriteByte(addr, byte(value))
}

func (c *Core) Run() error {

	limit := false
	if c.InstructionLimit > 0 {
		//fmt.Printf("Setting instruction limit to %d\n", c.InstructionLimit)
		limit = true
	}

	done := false
	var err error
	for !done {
		err = c.tick()
		if err != nil {
			return err
		}

		if limit {
			c.InstructionLimit -= 1
			if c.InstructionLimit <= 0 {
				if c.testing {
					return fmt.Errorf("Instruction limit hit")
				}
				done = true
			}
		}

		if c.testing {
			done = c.testDone
		}
	}

	return nil
}

func (c *Core) tick() error {
	if c.checkStuck {
		if c.PC == c.lastPC {
			c.lastSame++
		} else {
			c.lastSame = 0
			c.lastPC = c.PC
		}

		if c.lastSame > 5 {
			return fmt.Errorf("Stuck")
		}
	}

	opcode := c.ReadByte(c.PC)
	//if c.fullRW {
	//	fmt.Printf("[%06d] %04X: %02X\n", c.ticks, c.PC, opcode)
	//}

	if opcode == 0xFF && c.testing {
		c.testDone = true
		return nil // 0xFF means end of test
	}

	//fn, ok := opcodes[opcode]
	instr, ok := instructionList[opcode]
	if !ok || instr == nil {
		return fmt.Errorf("OP Code not implemented: [$%04X] $%02X", c.PC, opcode)
	}

	oppc := c.PC

	c.ticks++
	instr.Execute(c)

	if c.Debug {
		l := instr.InstrLength(c)
		ops := []string{}
		for i := uint8(0); i < l; i++ {
			ops = append(ops, fmt.Sprintf("%X", c.ReadByte(oppc+uint16(i))))
		}

		value := uint16(0)
		_, ok := instr.(Branch)
		if ok {
			value = c.addrRelative(c.ReadByte(oppc+1))
		} else if l == 2 {
			value = uint16(c.ReadByte(oppc+1))
		} else if l == 3 {
			value = c.ReadWord(oppc+1)
		}

		fmt.Printf("[%06d] $%04X: %-9s %s %s\n",
			c.ticks,
			oppc,
			strings.Join(ops, " "),
			instr.Name(),
			instr.AddressMeta().Asm(value),
		)
	}

	return nil
}

func (c *Core) DumpMemoryRange(start, end uint16) {
	if end < start {
		fmt.Println("Invalid dump range given")
		return
	}

	fmt.Printf("start: $%02X end: $%02X\n", start, end)

	vals := []byte{}
	current := start

	for current <= end {
		vals = append(vals, c.ReadByte(current))
		current++
	}

	for i, b := range vals {
		fmt.Printf("$%02X: $%02X (%d)\n", i+int(start), b, b)
	}
}

const (
	FLAG_CARRY     uint8 = 0x01
	FLAG_ZERO      uint8 = 0x02
	FLAG_INTERRUPT uint8 = 0x04
	FLAG_DECIMAL   uint8 = 0x08

	FLAG_OVERFLOW uint8 = 0x40
	FLAG_NEGATIVE uint8 = 0x80
)

func flagToString(ph uint8) string {
	switch ph {
	case FLAG_CARRY:
		return "FLAG_CARRY"
	case FLAG_ZERO:
		return "FLAG_ZERO"
	case FLAG_INTERRUPT:
		return "FLAG_INTERRUPT"
	case FLAG_DECIMAL:
		return "FLAG_DECIMAL"
	case FLAG_OVERFLOW:
		return "FLAG_OVERFLOW"
	case FLAG_NEGATIVE:
		return "FLAG_NEGATIVE"
	}
	return "FLAG_UNUSED"
}

func flagsToString(ph uint8) string {
	sc := "-"
	sz := "-"
	si := "-"
	sd := "-"
	sv := "-"
	sn := "-"

	if ph&FLAG_CARRY == 1 {
		sc = "C"
	}

	if ph&FLAG_ZERO == 1 {
		sz = "Z"
	}

	if ph&FLAG_INTERRUPT == 1 {
		si = "I"
	}

	if ph&FLAG_DECIMAL == 1 {
		sd = "D"
	}

	if ph&FLAG_OVERFLOW == 1 {
		sv = "V"
	}

	if ph&FLAG_NEGATIVE == 1 {
		sn = "N"
	}

	return fmt.Sprintf("%s%s--%s%s%s%s", sn, sv, sd, si, sz, sc)
}

func (c *Core) DumpRegisters() {
	fmt.Printf("A: %02X (%d) X: %02X (%d) Y: %02X (%d) PC: %04X Phlags: %s\n",
		c.A,
		c.A,
		c.X,
		c.X,
		c.Y,
		c.Y,
		c.PC,
		flagsToString(c.Phlags),
	)
}

func (c *Core) DumpPage(page uint8) {
	vals := []string{}
	base := uint16(page) << 8
	for i := uint16(0); i < 256; i++ {
		vals = append(vals, fmt.Sprintf("%02X", c.ReadByte(base+i)))
	}

	for i := 0; i < 256; i += 16 {
		fmt.Printf("%04X: %s\n", int(base)+i, strings.Join(vals[i:i+16], " "))
	}
}

func (c Core) DumpMemoryToFile(filename string) error {
	vals := []string{}
	for i := uint(0); i < 0x10000; i++ {
		vals = append(vals, fmt.Sprintf("%02X", c.ReadByte(uint16(i))))
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for i := 0; i < 0x10000; i += 16 {
		fmt.Fprintf(file, "%04X: %s\n", i, strings.Join(vals[i:i+16], " "))
	}
	return nil
}

func (c Core) Ticks() uint64 {
	return c.ticks
}

func (c *Core) tlog(msg string) {
	if c.t != nil {
		c.t.Log(msg)
	}
}

func (c *Core) tlogf(fmt string, args ...interface{}) {
	if c.t != nil {
		c.t.Logf(fmt, args...)
	}
}

func testCore(rom []byte, mem []byte, wram []byte) (*Core, error) {
	core, err := NewCore(rom, false, 1000)
	if err != nil {
		return nil, err
	}
	core.testing = true

	if mem != nil {
		for len(mem) < 0x1000 {
			mem = append(mem, 0x00)
		}
		core.memory = mem
	}

	if wram != nil {
		core.wram = wram
	}

	return core, core.Run()
}

func padWithVectors(rom []byte, nmi, reset, irq uint16) []byte {
	for len(rom)%256 != 0 {
		rom = append(rom, 0xFF)
	}

	addr := len(rom) - 6

	rom[addr] = byte(nmi & 0x00FF)
	rom[addr+1] = byte(nmi >> 8)

	rom[addr+2] = byte(reset & 0x00FF)
	rom[addr+3] = byte(reset >> 8)

	rom[addr+4] = byte(irq & 0x00FF)
	rom[addr+5] = byte(irq >> 8)

	return rom
}

func (c *Core) dbg(format string, args ...interface{}) {
	if c.t != nil {
		c.t.Logf(format, args...)
	}
}

// Set zero and negative flags based on the given value
func (c *Core) setZeroNegative(value uint8) {
	// zero
	if value == 0 {
		c.Phlags = c.Phlags | FLAG_ZERO
	} else {
		c.Phlags = c.Phlags & (FLAG_ZERO ^ 0xFF)
	}

	// negative
	if value&0x80 != 0 {
		c.Phlags = c.Phlags | FLAG_NEGATIVE
	} else {
		c.Phlags = c.Phlags & (FLAG_NEGATIVE ^ 0xFF)
	}
}

// addrRelative works differently than all other addressing functions.
// It takes the value for offset, and uses the PC of the instruction
// as the start point.  Call this before incrementing PC.
func (c *Core) addrRelative(offset uint8) uint16 {
	addr := c.PC + 2
	val, negative := TwosCompInv(offset)

	if negative {
		addr -= uint16(val)
	} else {
		addr += uint16(val)
	}

	return addr
}

func TwosCompInv(value uint8) (uint8, bool) {
	if value&0x80 != 0 {
		return (value ^ 0xFF) + 1, true
	}
	return value, false
}

