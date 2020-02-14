package emu

import (
	//"fmt"
)

//type AddressingModeFunc func(c *Core) (uint16, uint8)
type ExecFunc func(c *Core, address uint16)

type Instruction interface {
	Execute(c *Core)
	Name() string
	InstrLength(c *Core) uint8
	AddressMeta() AddressModeMeta
}

var instructionList = map[byte]Instruction{
	OP_BCC: Branch{
		OpCode: OP_BCC,
		Instruction: "BCC",
		Flag: FLAG_CARRY,
		Set: false},
	OP_BCS: Branch{
		OpCode: OP_BCS,
		Instruction: "BCS",
		Flag: FLAG_CARRY,
		Set: true},
	OP_BEQ: Branch{
		OpCode: OP_BEQ,
		Instruction: "BEQ",
		Flag: FLAG_ZERO,
		Set: true},
	OP_BMI: Branch{
		OpCode: OP_BMI,
		Instruction: "BMI",
		Flag: FLAG_NEGATIVE,
		Set: true},
	OP_BNE: Branch{
		OpCode: OP_BNE,
		Instruction: "BNE",
		Flag: FLAG_ZERO,
		Set: false},
	OP_BPL: Branch{
		OpCode: OP_BPL,
		Instruction: "BPL",
		Flag: FLAG_NEGATIVE,
		Set: false},
	OP_BVC: Branch{
		OpCode: OP_BVC,
		Instruction: "BVC",
		Flag: FLAG_OVERFLOW,
		Set: false},
	OP_BVS: Branch{
		OpCode: OP_BVS,
		Instruction: "BVS",
		Flag: FLAG_OVERFLOW,
		Set: true},

	OP_CLD: StandardInstruction{
		OpCode:         OP_CLD,
		Instruction:    "CLD",
		AddressMode: ADDR_Implied,
		Exec:           instr_CLD},

	OP_DEC_AB: ReadWriteModify{
		OpCode:         OP_DEC_AB,
		Instruction:    "DEC",
		AddressMode: ADDR_Absolute,
		Exec:           instr_DEC},
	OP_DEC_AX: ReadWriteModify{
		OpCode:         OP_DEC_AX,
		Instruction:    "DEC",
		AddressMode: ADDR_AbsoluteX,
		Exec:           instr_DEC},
	OP_DEC_ZP: ReadWriteModify{
		OpCode:         OP_DEC_ZP,
		Instruction:    "DEC",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_DEC},
	OP_DEC_ZX: ReadWriteModify{
		OpCode:         OP_DEC_ZX,
		Instruction:    "DEC",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_DEC},

	OP_DEX: StandardInstruction{
		OpCode:         OP_DEX,
		Instruction:    "DEX",
		AddressMode: ADDR_Implied,
		Exec:           instr_DEX},
	OP_DEY: StandardInstruction{
		OpCode:         OP_DEY,
		Instruction:    "DEY",
		AddressMode: ADDR_Implied,
		Exec:           instr_DEY},

	OP_JMP_AB: Jump{
		OpCode: OP_JMP_AB,
		Instruction: "JMP",
		AddressMode: ADDR_Absolute,
		Exec: instr_JMP},
	OP_JMP_ID: Jump{
		OpCode: OP_JMP_ID,
		Instruction: "JMP",
		AddressMode: ADDR_Indirect,
		Exec: instr_JMP},
	OP_JSR: Jump{
		OpCode: OP_JSR,
		Instruction: "JSR",
		AddressMode: ADDR_Absolute,
		Exec: instr_JSR},
	OP_RTS: Jump{
		OpCode: OP_RTS,
		Instruction: "RTS",
		AddressMode: ADDR_Implied,
		Exec: instr_RTS},
	OP_RTI: Jump{
		OpCode: OP_RTI,
		Instruction: "RTI",
		AddressMode: ADDR_Implied,
		Exec: instr_RTI},

	OP_LDA_AB: StandardInstruction{
		OpCode:         OP_LDA_AB,
		Instruction:    "LDA",
		AddressMode: ADDR_Absolute,
		Exec:           instr_LDA},
	OP_LDA_AX: StandardInstruction{
		OpCode:         OP_LDA_AX,
		Instruction:    "LDA",
		AddressMode: ADDR_AbsoluteX,
		Exec:           instr_LDA},
	OP_LDA_AY: StandardInstruction{
		OpCode:         OP_LDA_AY,
		Instruction:    "LDA",
		AddressMode: ADDR_AbsoluteY,
		Exec:           instr_LDA},
	OP_LDA_IM: StandardInstruction{
		OpCode:         OP_LDA_IM,
		Instruction:    "LDA",
		AddressMode: ADDR_Immediate,
		Exec:           instr_LDA},
	OP_LDA_IX: StandardInstruction{
		OpCode:         OP_LDA_IX,
		Instruction:    "LDA",
		AddressMode: ADDR_IndirectX,
		Exec:           instr_LDA},
	OP_LDA_IY: StandardInstruction{
		OpCode:         OP_LDA_IY,
		Instruction:    "LDA",
		AddressMode: ADDR_IndirectY,
		Exec:           instr_LDA},
	OP_LDA_ZP: StandardInstruction{
		OpCode:         OP_LDA_ZP,
		Instruction:    "LDA",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_LDA},
	OP_LDA_ZX: StandardInstruction{
		OpCode:         OP_LDA_ZX,
		Instruction:    "LDA",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_LDA},

	OP_LDX_AB: StandardInstruction{
		OpCode:         OP_LDX_AB,
		Instruction:    "LDX",
		AddressMode: ADDR_Absolute,
		Exec:           instr_LDX},
	OP_LDX_AY: StandardInstruction{
		OpCode:         OP_LDX_AY,
		Instruction:    "LDX",
		AddressMode: ADDR_AbsoluteY,
		Exec:           instr_LDX},
	OP_LDX_IM: StandardInstruction{
		OpCode:         OP_LDX_IM,
		Instruction:    "LDX",
		AddressMode: ADDR_Immediate,
		Exec:           instr_LDX},
	OP_LDX_ZP: StandardInstruction{
		OpCode:         OP_LDX_ZP,
		Instruction:    "LDX",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_LDX},
	OP_LDX_ZY: StandardInstruction{
		OpCode:         OP_LDX_ZY,
		Instruction:    "LDX",
		AddressMode: ADDR_ZeroPageY,
		Exec:           instr_LDX},

	OP_LDY_AB: StandardInstruction{
		OpCode:         OP_LDY_AB,
		Instruction:    "LDY",
		AddressMode: ADDR_Absolute,
		Exec:           instr_LDY},
	OP_LDY_AX: StandardInstruction{
		OpCode:         OP_LDY_AX,
		Instruction:    "LDY",
		AddressMode: ADDR_AbsoluteX,
		Exec:           instr_LDY},
	OP_LDY_IM: StandardInstruction{
		OpCode:         OP_LDY_IM,
		Instruction:    "LDY",
		AddressMode: ADDR_Immediate,
		Exec:           instr_LDY},
	OP_LDY_ZP: StandardInstruction{
		OpCode:         OP_LDY_ZP,
		Instruction:    "LDY",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_LDY},
	OP_LDY_ZX: StandardInstruction{
		OpCode:         OP_LDY_ZX,
		Instruction:    "LDY",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_LDY},

	OP_INC_AB: ReadWriteModify{
		OpCode:         OP_INC_AB,
		Instruction:    "INC",
		AddressMode: ADDR_Absolute,
		Exec:           instr_INC},
	OP_INC_AX: ReadWriteModify{
		OpCode:         OP_INC_AX,
		Instruction:    "INC",
		AddressMode: ADDR_AbsoluteX,
		Exec:           instr_INC},
	OP_INC_ZP: ReadWriteModify{
		OpCode:         OP_INC_ZP,
		Instruction:    "INC",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_INC},
	OP_INC_ZX: ReadWriteModify{
		OpCode:         OP_INC_ZX,
		Instruction:    "INC",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_INC},

	OP_INX: StandardInstruction{
		OpCode:         OP_INX,
		Instruction:    "INX",
		AddressMode: ADDR_Implied,
		Exec:           instr_INX},
	OP_INY: StandardInstruction{
		OpCode:         OP_INY,
		Instruction:    "INY",
		AddressMode: ADDR_Implied,
		Exec:           instr_INY},

	OP_NOP: StandardInstruction{
		OpCode:         OP_NOP,
		Instruction:    "NOP",
		AddressMode: ADDR_Implied,
		Exec:           instr_NOP},

	OP_STA_AB: StandardInstruction{
		OpCode:         OP_STA_AB,
		Instruction:    "STA",
		AddressMode: ADDR_Absolute,
		Exec:           instr_STA},
	OP_STA_AX: StandardInstruction{
		OpCode:         OP_STA_AX,
		Instruction:    "STA",
		AddressMode: ADDR_AbsoluteX,
		Exec:           instr_STA},
	OP_STA_AY: StandardInstruction{
		OpCode:         OP_STA_AY,
		Instruction:    "STA",
		AddressMode: ADDR_AbsoluteY,
		Exec:           instr_STA},
	OP_STA_IX: StandardInstruction{
		OpCode:         OP_STA_IX,
		Instruction:    "STA",
		AddressMode: ADDR_IndirectX,
		Exec:           instr_STA},
	OP_STA_IY: StandardInstruction{
		OpCode:         OP_STA_IY,
		Instruction:    "STA",
		AddressMode: ADDR_IndirectY,
		Exec:           instr_STA},
	OP_STA_ZP: StandardInstruction{
		OpCode:         OP_STA_ZP,
		Instruction:    "STA",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_STA},
	OP_STA_ZX: StandardInstruction{
		OpCode:         OP_STA_ZX,
		Instruction:    "STA",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_STA},

	OP_STX_AB: StandardInstruction{
		OpCode:         OP_STX_AB,
		Instruction:    "STX",
		AddressMode: ADDR_Absolute,
		Exec:           instr_STX},
	OP_STX_ZP: StandardInstruction{
		OpCode:         OP_STX_ZP,
		Instruction:    "STX",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_STX},
	OP_STX_ZY: StandardInstruction{
		OpCode:         OP_STX_ZY,
		Instruction:    "STX",
		AddressMode: ADDR_ZeroPageY,
		Exec:           instr_STX},

	OP_STY_AB: StandardInstruction{
		OpCode:         OP_STY_AB,
		Instruction:    "STY",
		AddressMode: ADDR_Absolute,
		Exec:           instr_STY},
	OP_STY_ZP: StandardInstruction{
		OpCode:         OP_STY_ZP,
		Instruction:    "STY",
		AddressMode: ADDR_ZeroPage,
		Exec:           instr_STY},
	OP_STY_ZX: StandardInstruction{
		OpCode:         OP_STY_ZX,
		Instruction:    "STY",
		AddressMode: ADDR_ZeroPageX,
		Exec:           instr_STY},

	OP_TAX: StandardInstruction{
		OpCode:         OP_TAX,
		Instruction:    "TAX",
		AddressMode: ADDR_Implied,
		Exec:           instr_TAX},
	OP_TAY: StandardInstruction{
		OpCode:         OP_TAY,
		Instruction:    "TAY",
		AddressMode: ADDR_Implied,
		Exec:           instr_TAY},
	OP_TSX: StandardInstruction{
		OpCode:         OP_TSX,
		Instruction:    "TSX",
		AddressMode: ADDR_Implied,
		Exec:           instr_TSX},
	OP_TXA: StandardInstruction{
		OpCode:         OP_TXA,
		Instruction:    "TXA",
		AddressMode: ADDR_Implied,
		Exec:           instr_TXA},
	OP_TXS: StandardInstruction{
		OpCode:         OP_TXS,
		Instruction:    "TXS",
		AddressMode: ADDR_Implied,
		Exec:           instr_TXS},
}

type StandardInstruction struct {
	AddressMode AddressModeMeta
	OpCode      byte
	Instruction string
	Exec        ExecFunc
}

func (i StandardInstruction) AddressMeta() AddressModeMeta {
	return i.AddressMode
}

func (i StandardInstruction) Execute(c *Core) {
	address, size := i.AddressMode.Address(c)
	i.Exec(c, address)
	c.PC += uint16(size)
}

func (i StandardInstruction) InstrLength(c *Core) uint8 {
	_, size := i.AddressMode.Address(c)
	return size
}

func (i StandardInstruction) Name() string {
	return i.Instruction
}

func instr_CLD(c *Core, address uint16) {
	c.Phlags = c.Phlags & (FLAG_DECIMAL ^ 0xFF)
}

func instr_DEX(c *Core, address uint16) {
	c.X -= 1
	c.setZeroNegative(c.X)
}

func instr_DEY(c *Core, address uint16) {
	c.Y -= 1
	c.setZeroNegative(c.Y)
}

func instr_INX(c *Core, address uint16) {
	c.X += 1
	c.setZeroNegative(c.X)
}

func instr_INY(c *Core, address uint16) {
	c.Y += 1
	c.setZeroNegative(c.Y)
}

func instr_LDA(c *Core, address uint16) {
	c.A = c.ReadByte(address)
	c.setZeroNegative(c.A)
}

func instr_LDX(c *Core, address uint16) {
	c.X = c.ReadByte(address)
	c.setZeroNegative(c.X)
}

func instr_LDY(c *Core, address uint16) {
	c.Y = c.ReadByte(address)
	c.setZeroNegative(c.Y)
}

func instr_NOP(c *Core, address uint16) {
	return
}

func instr_STA(c *Core, address uint16) {
	c.WriteByte(address, c.A)
}

func instr_STX(c *Core, address uint16) {
	c.WriteByte(address, c.X)
}

func instr_STY(c *Core, address uint16) {
	c.WriteByte(address, c.Y)
}

func instr_TAX(c *Core, address uint16) {
	c.X = c.A
	c.setZeroNegative(c.X)
}

func instr_TAY(c *Core, address uint16) {
	c.Y = c.A
	c.setZeroNegative(c.Y)
}

func instr_TSX(c *Core, address uint16) {
	c.X = c.SP
	c.setZeroNegative(c.X)
}

func instr_TXA(c *Core, address uint16) {
	c.A = c.X
	c.setZeroNegative(c.A)
}

func instr_TXS(c *Core, address uint16) {
	c.SP = c.X
}

type ReadWriteModify struct {
	OpCode         byte
	Instruction    string
	AddressMode AddressModeMeta
	Exec           func(c *Core, value uint8) uint8
}

func (rwm ReadWriteModify) AddressMeta() AddressModeMeta {
	return rwm.AddressMode
}

func (rwm ReadWriteModify) Execute(c *Core) {
	address, size := rwm.AddressMode.Address(c)
	c.WriteByte(address, rwm.Exec(c, c.ReadByte(address)))
	c.PC += uint16(size)
}

func (rwm ReadWriteModify) Name() string {
	return rwm.Instruction
}

func (rwm ReadWriteModify) InstrLength(c *Core) uint8 {
	_, size := rwm.AddressMode.Address(c)
	return size
}

func instr_DEC(c *Core, value uint8) uint8 {
	value -= 1
	c.setZeroNegative(value)
	return value
}

func instr_INC(c *Core, value uint8) uint8 {
	value += 1
	c.setZeroNegative(value)
	return value
}

type Branch struct {
	OpCode byte
	Instruction string
	Flag uint8
	Set bool
}

func (b Branch) AddressMeta() AddressModeMeta {
	return ADDR_Relative
}

func (b Branch) Name() string {
	return b.Instruction
}

func (b Branch) Execute(c *Core) {
	var v uint8 = 0
	if b.Set {
		v = 1
	}

	//prevPc := c.PC
	if c.Phlags & b.Flag == v {
		c.PC = c.addrRelative(c.ReadByte(c.PC + 1))
	} else {
		c.PC += 2
	}

	//fmt.Printf("%s: %s set: %t [%04X] -> [%04X]\n",
	//	b.Instruction,
	//	flagToString(Flag),
	//	b.Set,
	//	prevPc,
	//	c.PC,
	//)
	//c.DumpRegisters()
}

func (b Branch) InstrLength(c *Core) uint8 {
	return 2
}

// anything that modifies the PC directly, aside form branches
type Jump struct {
	OpCode byte
	Instruction string
	AddressMode AddressModeMeta
	Exec func(c *Core, address uint16) uint16
}

func (j Jump) Name() string {
	return j.Instruction
}

func (j Jump) Execute(c *Core) {
	address, _ := j.AddressMode.Address(c)
	c.PC = j.Exec(c, address)
}

func (j Jump) InstrLength(c *Core) uint8 {
	_, size := j.AddressMode.Address(c)
	return size
}

func (j Jump) AddressMeta() AddressModeMeta {
	return j.AddressMode
}

func instr_JMP(c *Core, address uint16) uint16 {
	return address
}

func instr_JSR(c *Core, address uint16) uint16 {
	c.pushAddress(c.PC + 2)
	return address
}

func instr_RTS(c *Core, address uint16) uint16 {
	return c.pullAddress()
}

func instr_RTI(c *Core, address uint16) uint16 {
	c.Phlags = c.pullByte()
	return c.pullAddress()
}

func (c *Core) pushAddress(addr uint16) {
	c.pushByte(uint8(addr >> 8))
	c.pushByte(uint8(addr & 0xFF))
}

func (c *Core) pullAddress() uint16 {
	return uint16(c.pullByte()) | uint16(c.pullByte()) << 8
}

func (c *Core) pushByte(val uint8) {
	c.WriteByte(uint16(c.SP) | 0x0100, val)
	c.SP -= 1
}

func (c *Core) pullByte() uint8 {
	c.SP += 1
	return c.ReadByte(uint16(c.SP) | 0x0100)
}
