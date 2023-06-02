package main

import (
	arm "github.com/zyedidia/isolator/arm64/arm64asm"
)

const (
	branchReg  = arm.X21
	resReg     = arm.X20
	retReg     = arm.X30
	segmentId  = 0xffc0
	bundleMask = uint64(0x07)
)

var dataRegs = map[arm.Reg]bool{
	resReg: true,
	arm.SP: true,
}

var ctrlRegs = map[arm.Reg]bool{
	branchReg: true,
	retReg:    true,
}

var restrictedRegs = map[arm.Reg]bool{
	arm.X30: true,
	arm.W30: true,
	arm.H30: true,
	arm.B30: true,

	arm.X21: true,
	arm.W21: true,
	arm.H21: true,
	arm.B21: true,

	arm.X20: true,
	arm.W20: true,
	arm.H20: true,
	arm.B20: true,

	arm.SP:  true,
	arm.WSP: true,
}

// list of permitted instructions
var allowed = map[arm.Op]bool{
	arm.ABS:    true,
	arm.ADC:    true,
	arm.ADCS:   true,
	arm.ADD:    true,
	arm.ADDHN:  true,
	arm.ADDHN2: true,
	arm.ADDP:   true,
	arm.ADDS:   true,
	arm.ADDV:   true,
	arm.ADR:    true,
	arm.ADRP:   true,
	arm.AESD:   true,
	arm.AESE:   true,
	arm.AESIMC: true,
	arm.AESMC:  true,
	arm.AND:    true,
	arm.ANDS:   true,
	arm.ASR:    true,
	arm.ASRV:   true,
	// arm.AT:     true,
	arm.B:     true,
	arm.BFI:   true,
	arm.BFM:   true,
	arm.BFXIL: true,
	arm.BIC:   true,
	arm.BICS:  true,
	arm.BIF:   true,
	arm.BIT:   true,
	arm.BL:    true,
	arm.BLR:   true,
	arm.BR:    true,
	// arm.BRK:       true,
	arm.BSL:  true,
	arm.CBNZ: true,
	arm.CBZ:  true,
	arm.CCMN: true,
	arm.CCMP: true,
	arm.CINC: true,
	arm.CINV: true,
	// arm.CLREX:     true,
	arm.CLS:     true,
	arm.CLZ:     true,
	arm.CMEQ:    true,
	arm.CMGE:    true,
	arm.CMGT:    true,
	arm.CMHI:    true,
	arm.CMHS:    true,
	arm.CMLE:    true,
	arm.CMLT:    true,
	arm.CMN:     true,
	arm.CMP:     true,
	arm.CMTST:   true,
	arm.CNEG:    true,
	arm.CNT:     true,
	arm.CRC32B:  true,
	arm.CRC32CB: true,
	arm.CRC32CH: true,
	arm.CRC32CW: true,
	arm.CRC32CX: true,
	arm.CRC32H:  true,
	arm.CRC32W:  true,
	arm.CRC32X:  true,
	arm.CSEL:    true,
	arm.CSET:    true,
	arm.CSETM:   true,
	arm.CSINC:   true,
	arm.CSINV:   true,
	arm.CSNEG:   true,
	// arm.DC:      true,
	// debug state instructions
	// arm.DCPS1:     true,
	// arm.DCPS2:     true,
	// arm.DCPS3:     true,
	arm.DMB: true,
	// arm.DRPS:      true,
	arm.DSB: true,
	arm.DUP: true,
	arm.EON: true,
	arm.EOR: true,
	// arm.ERET:      true,
	arm.EXT:     true,
	arm.EXTR:    true,
	arm.FABD:    true,
	arm.FABS:    true,
	arm.FACGE:   true,
	arm.FACGT:   true,
	arm.FADD:    true,
	arm.FADDP:   true,
	arm.FCCMP:   true,
	arm.FCCMPE:  true,
	arm.FCMEQ:   true,
	arm.FCMGE:   true,
	arm.FCMGT:   true,
	arm.FCMLE:   true,
	arm.FCMLT:   true,
	arm.FCMP:    true,
	arm.FCMPE:   true,
	arm.FCSEL:   true,
	arm.FCVT:    true,
	arm.FCVTAS:  true,
	arm.FCVTAU:  true,
	arm.FCVTL:   true,
	arm.FCVTL2:  true,
	arm.FCVTMS:  true,
	arm.FCVTMU:  true,
	arm.FCVTN:   true,
	arm.FCVTN2:  true,
	arm.FCVTNS:  true,
	arm.FCVTNU:  true,
	arm.FCVTPS:  true,
	arm.FCVTPU:  true,
	arm.FCVTXN:  true,
	arm.FCVTXN2: true,
	arm.FCVTZS:  true,
	arm.FCVTZU:  true,
	arm.FDIV:    true,
	arm.FMADD:   true,
	arm.FMAX:    true,
	arm.FMAXNM:  true,
	arm.FMAXNMP: true,
	arm.FMAXNMV: true,
	arm.FMAXP:   true,
	arm.FMAXV:   true,
	arm.FMIN:    true,
	arm.FMINNM:  true,
	arm.FMINNMP: true,
	arm.FMINNMV: true,
	arm.FMINP:   true,
	arm.FMINV:   true,
	arm.FMLA:    true,
	arm.FMLS:    true,
	arm.FMOV:    true,
	arm.FMSUB:   true,
	arm.FMUL:    true,
	arm.FMULX:   true,
	arm.FNEG:    true,
	arm.FNMADD:  true,
	arm.FNMSUB:  true,
	arm.FNMUL:   true,
	arm.FRECPE:  true,
	arm.FRECPS:  true,
	arm.FRECPX:  true,
	arm.FRINTA:  true,
	arm.FRINTI:  true,
	arm.FRINTM:  true,
	arm.FRINTN:  true,
	arm.FRINTP:  true,
	arm.FRINTX:  true,
	arm.FRINTZ:  true,
	arm.FRSQRTE: true,
	arm.FRSQRTS: true,
	arm.FSQRT:   true,
	arm.FSUB:    true,
	arm.HINT:    true,
	// arm.HLT:       true,
	// arm.HVC:       true,
	// arm.IC:        true,
	arm.INS: true,
	arm.ISB: true,
	// arm.LD1:       true,
	// arm.LD1R:      true,
	// arm.LD2:       true,
	// arm.LD2R:      true,
	// arm.LD3:       true,
	// arm.LD3R:      true,
	// arm.LD4:       true,
	// arm.LD4R:      true,
	// arm.LDAR:      true,
	// arm.LDARB:     true,
	// arm.LDARH:     true,
	// arm.LDAXP:     true,
	// arm.LDAXR:     true,
	// arm.LDAXRB:    true,
	// arm.LDAXRH:    true,
	// arm.LDNP:      true,
	arm.LDP:   true,
	arm.LDPSW: true,
	arm.LDR:   true,
	arm.LDRB:  true,
	arm.LDRH:  true,
	arm.LDRSB: true,
	arm.LDRSH: true,
	arm.LDRSW: true,
	// arm.LDTR:      true,
	// arm.LDTRB:     true,
	// arm.LDTRH:     true,
	// arm.LDTRSB:    true,
	// arm.LDTRSH:    true,
	// arm.LDTRSW:    true,
	arm.LDUR:   true,
	arm.LDURB:  true,
	arm.LDURH:  true,
	arm.LDURSB: true,
	arm.LDURSH: true,
	arm.LDURSW: true,
	// arm.LDXP:      true,
	// arm.LDXR:      true,
	// arm.LDXRB:     true,
	// arm.LDXRH:     true,
	arm.LSL:  true,
	arm.LSLV: true,
	arm.LSR:  true,
	arm.LSRV: true,
	arm.MADD: true,
	arm.MLA:  true,
	arm.MLS:  true,
	arm.MNEG: true,
	arm.MOV:  true,
	arm.MOVI: true,
	arm.MOVK: true,
	arm.MOVN: true,
	arm.MOVZ: true,
	// arm.MRS:       true,
	// arm.MSR:       true,
	arm.MSUB:   true,
	arm.MUL:    true,
	arm.MVN:    true,
	arm.MVNI:   true,
	arm.NEG:    true,
	arm.NEGS:   true,
	arm.NGC:    true,
	arm.NGCS:   true,
	arm.NOP:    true,
	arm.NOT:    true,
	arm.ORN:    true,
	arm.ORR:    true,
	arm.PMUL:   true,
	arm.PMULL:  true,
	arm.PMULL2: true,
	// arm.PRFM:      true,
	// arm.PRFUM:     true,
	arm.RADDHN:  true,
	arm.RADDHN2: true,
	arm.RBIT:    true,
	arm.RET:     true,
	arm.REV:     true,
	arm.REV16:   true,
	arm.REV32:   true,
	arm.REV64:   true,
	arm.ROR:     true,
	arm.RORV:    true,
	arm.RSHRN:   true,
	arm.RSHRN2:  true,
	arm.RSUBHN:  true,
	arm.RSUBHN2: true,
	arm.SABA:    true,
	arm.SABAL:   true,
	arm.SABAL2:  true,
	arm.SABD:    true,
	arm.SABDL:   true,
	arm.SABDL2:  true,
	arm.SADALP:  true,
	arm.SADDL:   true,
	arm.SADDL2:  true,
	arm.SADDLP:  true,
	arm.SADDLV:  true,
	arm.SADDW:   true,
	arm.SADDW2:  true,
	arm.SBC:     true,
	arm.SBCS:    true,
	arm.SBFIZ:   true,
	arm.SBFM:    true,
	arm.SBFX:    true,
	arm.SCVTF:   true,
	arm.SDIV:    true,
	// arm.SEV:       true,
	// arm.SEVL:      true,
	arm.SHA1C:     true,
	arm.SHA1H:     true,
	arm.SHA1M:     true,
	arm.SHA1P:     true,
	arm.SHA1SU0:   true,
	arm.SHA1SU1:   true,
	arm.SHA256H:   true,
	arm.SHA256H2:  true,
	arm.SHA256SU0: true,
	arm.SHA256SU1: true,
	arm.SHADD:     true,
	arm.SHL:       true,
	arm.SHLL:      true,
	arm.SHLL2:     true,
	arm.SHRN:      true,
	arm.SHRN2:     true,
	arm.SHSUB:     true,
	arm.SLI:       true,
	arm.SMADDL:    true,
	arm.SMAX:      true,
	arm.SMAXP:     true,
	arm.SMAXV:     true,
	// arm.SMC:       true,
	arm.SMIN:      true,
	arm.SMINP:     true,
	arm.SMINV:     true,
	arm.SMLAL:     true,
	arm.SMLAL2:    true,
	arm.SMLSL:     true,
	arm.SMLSL2:    true,
	arm.SMNEGL:    true,
	arm.SMOV:      true,
	arm.SMSUBL:    true,
	arm.SMULH:     true,
	arm.SMULL:     true,
	arm.SMULL2:    true,
	arm.SQABS:     true,
	arm.SQADD:     true,
	arm.SQDMLAL:   true,
	arm.SQDMLAL2:  true,
	arm.SQDMLSL:   true,
	arm.SQDMLSL2:  true,
	arm.SQDMULH:   true,
	arm.SQDMULL:   true,
	arm.SQDMULL2:  true,
	arm.SQNEG:     true,
	arm.SQRDMULH:  true,
	arm.SQRSHL:    true,
	arm.SQRSHRN:   true,
	arm.SQRSHRN2:  true,
	arm.SQRSHRUN:  true,
	arm.SQRSHRUN2: true,
	arm.SQSHL:     true,
	arm.SQSHLU:    true,
	arm.SQSHRN:    true,
	arm.SQSHRN2:   true,
	arm.SQSHRUN:   true,
	arm.SQSHRUN2:  true,
	arm.SQSUB:     true,
	arm.SQXTN:     true,
	arm.SQXTN2:    true,
	arm.SQXTUN:    true,
	arm.SQXTUN2:   true,
	arm.SRHADD:    true,
	arm.SRI:       true,
	arm.SRSHL:     true,
	arm.SRSHR:     true,
	arm.SRSRA:     true,
	arm.SSHL:      true,
	arm.SSHLL:     true,
	arm.SSHLL2:    true,
	arm.SSHR:      true,
	arm.SSRA:      true,
	arm.SSUBL:     true,
	arm.SSUBL2:    true,
	arm.SSUBW:     true,
	arm.SSUBW2:    true,
	// arm.ST1:       true,
	// arm.ST2:       true,
	// arm.ST3:       true,
	// arm.ST4:       true,
	// arm.STLR:      true,
	// arm.STLRB:     true,
	// arm.STLRH:     true,
	// arm.STLXP:     true,
	// arm.STLXR:     true,
	// arm.STLXRB:    true,
	// arm.STLXRH:    true,
	// arm.STNP:      true,
	arm.STP:  true,
	arm.STR:  true,
	arm.STRB: true,
	arm.STRH: true,
	// arm.STTR:      true,
	// arm.STTRB:     true,
	// arm.STTRH:     true,
	arm.STUR:  true,
	arm.STURB: true,
	arm.STURH: true,
	// arm.STXP:      true,
	// arm.STXR:      true,
	// arm.STXRB:     true,
	// arm.STXRH:     true,
	arm.SUB:    true,
	arm.SUBHN:  true,
	arm.SUBHN2: true,
	arm.SUBS:   true,
	arm.SUQADD: true,
	// arm.SVC:       true,
	arm.SXTB:  true,
	arm.SXTH:  true,
	arm.SXTL:  true,
	arm.SXTL2: true,
	arm.SXTW:  true,
	// arm.SYS:       true,
	// arm.SYSL:      true,
	arm.TBL:  true,
	arm.TBNZ: true,
	arm.TBX:  true,
	arm.TBZ:  true,
	// arm.TLBI:      true,
	arm.TRN1:     true,
	arm.TRN2:     true,
	arm.TST:      true,
	arm.UABA:     true,
	arm.UABAL:    true,
	arm.UABAL2:   true,
	arm.UABD:     true,
	arm.UABDL:    true,
	arm.UABDL2:   true,
	arm.UADALP:   true,
	arm.UADDL:    true,
	arm.UADDL2:   true,
	arm.UADDLP:   true,
	arm.UADDLV:   true,
	arm.UADDW:    true,
	arm.UADDW2:   true,
	arm.UBFIZ:    true,
	arm.UBFM:     true,
	arm.UBFX:     true,
	arm.UCVTF:    true,
	arm.UDIV:     true,
	arm.UHADD:    true,
	arm.UHSUB:    true,
	arm.UMADDL:   true,
	arm.UMAX:     true,
	arm.UMAXP:    true,
	arm.UMAXV:    true,
	arm.UMIN:     true,
	arm.UMINP:    true,
	arm.UMINV:    true,
	arm.UMLAL:    true,
	arm.UMLAL2:   true,
	arm.UMLSL:    true,
	arm.UMLSL2:   true,
	arm.UMNEGL:   true,
	arm.UMOV:     true,
	arm.UMSUBL:   true,
	arm.UMULH:    true,
	arm.UMULL:    true,
	arm.UMULL2:   true,
	arm.UQADD:    true,
	arm.UQRSHL:   true,
	arm.UQRSHRN:  true,
	arm.UQRSHRN2: true,
	arm.UQSHL:    true,
	arm.UQSHRN:   true,
	arm.UQSHRN2:  true,
	arm.UQSUB:    true,
	arm.UQXTN:    true,
	arm.UQXTN2:   true,
	arm.URECPE:   true,
	arm.URHADD:   true,
	arm.URSHL:    true,
	arm.URSHR:    true,
	arm.URSQRTE:  true,
	arm.URSRA:    true,
	arm.USHL:     true,
	arm.USHLL:    true,
	arm.USHLL2:   true,
	arm.USHR:     true,
	arm.USQADD:   true,
	arm.USRA:     true,
	arm.USUBL:    true,
	arm.USUBL2:   true,
	arm.USUBW:    true,
	arm.USUBW2:   true,
	arm.UXTB:     true,
	arm.UXTH:     true,
	arm.UXTL:     true,
	arm.UXTL2:    true,
	arm.UZP1:     true,
	arm.UZP2:     true,
	// arm.WFE:       true,
	// arm.WFI:       true,
	arm.XTN:  true,
	arm.XTN2: true,
	// arm.YIELD:     true,
	arm.ZIP1: true,
	arm.ZIP2: true,
}

var stores = map[arm.Op]bool{
	arm.STP:   true,
	arm.STR:   true,
	arm.STRB:  true,
	arm.STRH:  true,
	arm.STUR:  true,
	arm.STURB: true,
	arm.STURH: true,
}

var branches = map[arm.Op]bool{
	arm.RET:  true,
	arm.BLR:  true,
	arm.BL:   true,
	arm.BR:   true,
	arm.TBZ:  true,
	arm.TBNZ: true,
	arm.CBZ:  true,
	arm.CBNZ: true,
	arm.B:    true,
}

var nomodify = map[arm.Op]bool{
	arm.NOP: true,
}