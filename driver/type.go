package driver

/*import (
	"encoding/binary"
)*/

//we ignore go conventions here as it's basically a copypaste

func IXGBE_BY_MAC(hw, r int) int { return 0 }

const (
	IXGBE_CTRL            = 0x00000
	IXGBE_STATUS          = 0x00008
	IXGBE_CTRL_EXT        = 0x00018
	IXGBE_ESDP            = 0x00020
	IXGBE_EODSDP          = 0x00028
	IXGBE_I2CCTL_82599    = 0x00028
	IXGBE_I2CCTL          = IXGBE_I2CCTL_82599
	IXGBE_I2CCTL_X540     = IXGBE_I2CCTL_82599
	IXGBE_I2CCTL_X550     = 0x15F5C
	IXGBE_I2CCTL_X550EM_x = IXGBE_I2CCTL_X550
	IXGBE_I2CCTL_X550EM_a = IXGBE_I2CCTL_X550
)

const IXGBE_EEC = 0x10010

func IXGBE_I2CCTL_BY_MAC(hw int) int { return IXGBE_BY_MAC(hw, IXGBE_I2CCTL) } //not sure if correct second arg

const (
	IXGBE_PHY_GPIO       = 0x00028
	IXGBE_MAC_GPIO       = 0x00030
	IXGBE_PHYINT_STATUS0 = 0x00100
	IXGBE_PHYINT_STATUS1 = 0x00104
	IXGBE_PHYINT_STATUS2 = 0x00108
	IXGBE_LEDCTL         = 0x00200
	IXGBE_FRTIMER        = 0x00048
	IXGBE_TCPTIMER       = 0x0004C
	IXGBE_CORESPARE      = 0x00600
	IXGBE_EXVET          = 0x05078
)

// MAC Registers
const (
	IXGBE_PCS1GCFIG      = 0x04200
	IXGBE_PCS1GLCTL      = 0x04208
	IXGBE_PCS1GLSTA      = 0x0420C
	IXGBE_PCS1GDBG0      = 0x04210
	IXGBE_PCS1GDBG1      = 0x04214
	IXGBE_PCS1GANA       = 0x04218
	IXGBE_PCS1GANLP      = 0x0421C
	IXGBE_PCS1GANNP      = 0x04220
	IXGBE_PCS1GANLPNP    = 0x04224
	IXGBE_HLREG0         = 0x04240
	IXGBE_HLREG1         = 0x04244
	IXGBE_PAP            = 0x04248
	IXGBE_MACA           = 0x0424C
	IXGBE_APAE           = 0x04250
	IXGBE_ARD            = 0x04254
	IXGBE_AIS            = 0x04258
	IXGBE_MSCA           = 0x0425C
	IXGBE_MSRWD          = 0x04260
	IXGBE_MLADD          = 0x04264
	IXGBE_MHADD          = 0x04268
	IXGBE_MAXFRS         = 0x04268
	IXGBE_TREG           = 0x0426C
	IXGBE_PCSS1          = 0x04288
	IXGBE_PCSS2          = 0x0428C
	IXGBE_XPCSS          = 0x04290
	IXGBE_MFLCN          = 0x04294
	IXGBE_SERDESC        = 0x04298
	IXGBE_MAC_SGMII_BUSY = 0x04298
	IXGBE_MACS           = 0x0429C
	IXGBE_AUTOC          = 0x042A0
	IXGBE_LINKS          = 0x042A4
	IXGBE_LINKS2         = 0x04324
	IXGBE_AUTOC2         = 0x042A8
	IXGBE_AUTOC3         = 0x042AC
	IXGBE_ANLP1          = 0x042B0
	IXGBE_ANLP2          = 0x042B4
	IXGBE_MACC           = 0x04330
	IXGBE_ATLASCTL       = 0x04800
	IXGBE_MMNGC          = 0x042D0
	IXGBE_ANLPNP1        = 0x042D4
	IXGBE_ANLPNP2        = 0x042D8
	IXGBE_KRPCSFC        = 0x042E0
	IXGBE_KRPCSS         = 0x042E4
	IXGBE_FECS1          = 0x042E8
	IXGBE_FECS2          = 0x042EC
	IXGBE_SMADARCTL      = 0x14F10
	IXGBE_MPVC           = 0x04318
	IXGBE_SGMIIC         = 0x04314
)

/* AUTOC Bit Masks */
const (
	IXGBE_AUTOC_KX4_KX_SUPP_MASK    = 0xC0000000
	IXGBE_AUTOC_KX4_SUPP            = 0x80000000
	IXGBE_AUTOC_KX_SUPP             = 0x40000000
	IXGBE_AUTOC_PAUSE               = 0x30000000
	IXGBE_AUTOC_ASM_PAUSE           = 0x20000000
	IXGBE_AUTOC_SYM_PAUSE           = 0x10000000
	IXGBE_AUTOC_RF                  = 0x08000000
	IXGBE_AUTOC_PD_TMR              = 0x06000000
	IXGBE_AUTOC_AN_RX_LOOSE         = 0x01000000
	IXGBE_AUTOC_AN_RX_DRIFT         = 0x00800000
	IXGBE_AUTOC_AN_RX_ALIGN         = 0x007C0000
	IXGBE_AUTOC_FECA                = 0x00040000
	IXGBE_AUTOC_FECR                = 0x00020000
	IXGBE_AUTOC_KR_SUPP             = 0x00010000
	IXGBE_AUTOC_AN_RESTART          = 0x00001000
	IXGBE_AUTOC_FLU                 = 0x00000001
	IXGBE_AUTOC_LMS_SHIFT           = 13
	IXGBE_AUTOC_LMS_10G_SERIAL      = (0x3 << IXGBE_AUTOC_LMS_SHIFT)
	IXGBE_AUTOC_LMS_KX4_KX_KR       = (0x4 << IXGBE_AUTOC_LMS_SHIFT)
	IXGBE_AUTOC_LMS_SGMII_1G_100M   = (0x5 << IXGBE_AUTOC_LMS_SHIFT)
	IXGBE_AUTOC_LMS_KX4_KX_KR_1G_AN = (0x6 << IXGBE_AUTOC_LMS_SHIFT)
	IXGBE_AUTOC_LMS_KX4_KX_KR_SGMII = (0x7 << IXGBE_AUTOC_LMS_SHIFT)
	IXGBE_AUTOC_LMS_MASK            = (0x7 << IXGBE_AUTOC_LMS_SHIFT)
	IXGBE_AUTOC_LMS_1G_LINK_NO_AN   = (0x0 << IXGBE_AUTOC_LMS_SHIFT)
	IXGBE_AUTOC_LMS_10G_LINK_NO_AN  = (0x1 << IXGBE_AUTOC_LMS_SHIFT)
	IXGBE_AUTOC_LMS_1G_AN           = (0x2 << IXGBE_AUTOC_LMS_SHIFT)
	IXGBE_AUTOC_LMS_KX4_AN          = (0x4 << IXGBE_AUTOC_LMS_SHIFT)
	IXGBE_AUTOC_LMS_KX4_AN_1G_AN    = (0x6 << IXGBE_AUTOC_LMS_SHIFT)
	IXGBE_AUTOC_LMS_ATTACH_TYPE     = (0x7 << IXGBE_AUTOC_10G_PMA_PMD_SHIFT)

	IXGBE_AUTOC_1G_PMA_PMD_MASK   = 0x00000200
	IXGBE_AUTOC_1G_PMA_PMD_SHIFT  = 9
	IXGBE_AUTOC_10G_PMA_PMD_MASK  = 0x00000180
	IXGBE_AUTOC_10G_PMA_PMD_SHIFT = 7
	IXGBE_AUTOC_10G_XAUI          = (0x0 << IXGBE_AUTOC_10G_PMA_PMD_SHIFT)
	IXGBE_AUTOC_10G_KX4           = (0x1 << IXGBE_AUTOC_10G_PMA_PMD_SHIFT)
	IXGBE_AUTOC_10G_CX4           = (0x2 << IXGBE_AUTOC_10G_PMA_PMD_SHIFT)
	IXGBE_AUTOC_1G_BX             = (0x0 << IXGBE_AUTOC_1G_PMA_PMD_SHIFT)
	IXGBE_AUTOC_1G_KX             = (0x1 << IXGBE_AUTOC_1G_PMA_PMD_SHIFT)
	IXGBE_AUTOC_1G_SFI            = (0x0 << IXGBE_AUTOC_1G_PMA_PMD_SHIFT)
	IXGBE_AUTOC_1G_KX_BX          = (0x1 << IXGBE_AUTOC_1G_PMA_PMD_SHIFT)

	IXGBE_AUTOC2_UPPER_MASK               = 0xFFFF0000
	IXGBE_AUTOC2_10G_SERIAL_PMA_PMD_MASK  = 0x00030000
	IXGBE_AUTOC2_10G_SERIAL_PMA_PMD_SHIFT = 16
	IXGBE_AUTOC2_10G_KR                   = (0x0 << IXGBE_AUTOC2_10G_SERIAL_PMA_PMD_SHIFT)
	IXGBE_AUTOC2_10G_XFI                  = (0x1 << IXGBE_AUTOC2_10G_SERIAL_PMA_PMD_SHIFT)
	IXGBE_AUTOC2_10G_SFI                  = (0x2 << IXGBE_AUTOC2_10G_SERIAL_PMA_PMD_SHIFT)
	IXGBE_AUTOC2_LINK_DISABLE_ON_D3_MASK  = 0x50000000
	IXGBE_AUTOC2_LINK_DISABLE_MASK        = 0x70000000
)

/* Receive Config masks */
const (
	IXGBE_RXCTRL_RXEN      = 0x00000001 /* Enable Receiver */
	IXGBE_RXCTRL_DMBYPS    = 0x00000002 /* Desc Monitor Bypass */
	IXGBE_RXDCTL_ENABLE    = 0x02000000 /* Ena specific Rx Queue */
	IXGBE_RXDCTL_SWFLSH    = 0x04000000 /* Rx Desc wr-bk flushing */
	IXGBE_RXDCTL_RLPMLMASK = 0x00003FFF /* X540 supported only */
	IXGBE_RXDCTL_RLPML_EN  = 0x00008000
	IXGBE_RXDCTL_VME       = 0x40000000 /* VLAN mode enable */
)

const (
	/* Transmit DMA registers */
	IXGBE_DTXCTL = 0x07E00

	IXGBE_DMATXCTL   = 0x04A80
	IXGBE_PFDTXGSWC  = 0x08220
	IXGBE_DTXMXSZRQ  = 0x08100
	IXGBE_DTXTCPFLGL = 0x04A88
	IXGBE_DTXTCPFLGH = 0x04A8C
	IXGBE_LBDRPEN    = 0x0CA00

	IXGBE_DMATXCTL_TE       = 0x1  /* Transmit Enable */
	IXGBE_DMATXCTL_NS       = 0x2  /* No Snoop LSO hdr buffer */
	IXGBE_DMATXCTL_GDV      = 0x8  /* Global Double VLAN */
	IXGBE_DMATXCTL_MDP_EN   = 0x20 /* Bit 5 */
	IXGBE_DMATXCTL_MBINTEN  = 0x40 /* Bit 6 */
	IXGBE_DMATXCTL_VT_SHIFT = 16   /* VLAN EtherType */
)

func IXGBE_TDBAL(i int) int {
	return 0x06000 + i*0x40 /* 32 of them (0-31)*/
}
func IXGBE_TDBAH(i int) int {
	return 0x06004 + i*0x40
}
func IXGBE_TDLEN(i int) int {
	return 0x06008 + i*0x40
}
func IXGBE_TDH(i int) int {
	return 0x06010 + i*0x40
}
func IXGBE_TDT(i int) int {
	return 0x06018 + i*0x40
}
func IXGBE_TXDCTL(i int) int {
	return 0x06028 + i*0x40
}
func IXGBE_TDWBAL(i int) int {
	return 0x06038 + i*0x40
}
func IXGBE_TDWBAH(i int) int {
	return 0x0603C + i*0x40
}
func IXGBE_PFVFSPOOF(i int) int {
	return 0x08200 + i*4 /* 8 of these 0 - 7 */
}
func IXGBE_TXPBTHRESH(i int) int {
	return 0x04950 + i*4 /* 8 of these 0 - 7 */
}

const (
	IXGBE_TXDCTL_ENABLE        = 0x02000000 /* Ena specific Tx Queue */
	IXGBE_TXDCTL_SWFLSH        = 0x04000000 /* Tx Desc. wr-bk flushing */
	IXGBE_TXDCTL_WTHRESH_SHIFT = 16         /* shift to WTHRESH bits */
)

//Recieve DMA Registers
//Todo: Copypaste für den Rest
func IXGBE_RDBAL(i int) int {
	if i < 64 {
		return 0x01000 + i*0x40
	}
	return 0x0D000 + (i-64)*0x40
}
func IXGBE_RDBAH(i int) int {
	if i < 64 {
		return 0x01004 + i*0x40
	}
	return 0x0D004 + (i-64)*0x40
}
func IXGBE_RDLEN(i int) int {
	if i < 64 {
		return 0x01008 + i*0x40
	}
	return 0x0D008 + (i-64)*0x40
}
func IXGBE_RDH(i int) int {
	if i < 64 {
		return 0x01010 + i*0x40
	}
	return 0x0D010 + (i-64)*0x40
}
func IXGBE_RDT(i int) int {
	if i < 64 {
		return 0x01018 + i*0x40
	}
	return 0x0D018 + (i-64)*0x40
}
func IXGBE_RXDCTL(i int) int {
	if i < 64 {
		return 0x01028 + i*0x40
	}
	return 0x0D028 + (i-64)*0x40
}
func IXGBE_RSCCTL(i int) int {
	if i < 64 {
		return 0x0102C + i*0x40
	}
	return 0x0D02C + (i-64)*0x40
}

/*
 * Split and Replication Receive Control Registers
 * 00-15 : 0x02100 + n*4
 * 16-64 : 0x01014 + n*0x40
 * 64-127: 0x0D014 + (n-64)*0x40
 */
func IXGBE_SRRCTL(i int) int {
	if i <= 15 {
		return 0x02100 + i*4
	}
	if i < 64 {
		return 0x01014 + i*0x40
	}
	return 0x0D014 + (i-64)*0x40
}

/*
 * Rx DCA Control Register:
 * 00-15 : 0x02200 + n*4
 * 16-64 : 0x0100C + n*0x40
 * 64-127: 0x0D00C + (n-64)*0x40
 */
func IXGBE_DCA_RXCTRL(i int) int {
	if i <= 15 {
		return 0x02200 + i*4
	}
	if i < 64 {
		return 0x0100C + i*0x40
	}
	return 0x0D00C + (i-64)*0x40
}

/* 8 of these 0x03C00 - 0x03C1C */
func IXGBE_RXPBSIZE(i int) int {
	return 0x03C00 + i*4
}

const (
	IXGBE_RDRXCTL        = 0x02F00
	IXGBE_RXCTRL         = 0x03000
	IXGBE_DROPEN         = 0x03D04
	IXGBE_RXPBSIZE_SHIFT = 10
	IXGBE_RXPBSIZE_MASK  = 0x000FFC00

	/* Receive Registers */
	IXGBE_RXCSUM           = 0x05000
	IXGBE_RFCTL            = 0x05008
	IXGBE_DRECCCTL         = 0x02F08
	IXGBE_DRECCCTL_DISABLE = 0
	IXGBE_DRECCCTL2        = 0x02F8C
)

/* Packet Buffer Initialization */
const (
	IXGBE_MAX_PACKET_BUFFERS = 8

	IXGBE_TXPBSIZE_20KB  = 0x00005000 /* 20KB Packet Buffer */
	IXGBE_TXPBSIZE_40KB  = 0x0000A000 /* 40KB Packet Buffer */
	IXGBE_RXPBSIZE_48KB  = 0x0000C000 /* 48KB Packet Buffer */
	IXGBE_RXPBSIZE_64KB  = 0x00010000 /* 64KB Packet Buffer */
	IXGBE_RXPBSIZE_80KB  = 0x00014000 /* 80KB Packet Buffer */
	IXGBE_RXPBSIZE_128KB = 0x00020000 /* 128KB Packet Buffer */
	IXGBE_RXPBSIZE_MAX   = 0x00080000 /* 512KB Packet Buffer */
	IXGBE_TXPBSIZE_MAX   = 0x00028000 /* 160KB Packet Buffer */

	IXGBE_TXPKT_SIZE_MAX = 0xA /* Max Tx Packet size */
	IXGBE_MAX_PB         = 8
)

func IXGBE_TXPBSIZE(i int) int {
	return 0x0CC00 + i*4 /* 8 of these */
}

/* HLREG0 Bit Masks */
const (
	IXGBE_HLREG0_TXCRCEN      = 0x00000001 /* bit  0 */
	IXGBE_HLREG0_RXCRCSTRP    = 0x00000002 /* bit  1 */
	IXGBE_HLREG0_JUMBOEN      = 0x00000004 /* bit  2 */
	IXGBE_HLREG0_TXPADEN      = 0x00000400 /* bit 10 */
	IXGBE_HLREG0_TXPAUSEEN    = 0x00001000 /* bit 12 */
	IXGBE_HLREG0_RXPAUSEEN    = 0x00004000 /* bit 14 */
	IXGBE_HLREG0_LPBK         = 0x00008000 /* bit 15 */
	IXGBE_HLREG0_MDCSPD       = 0x00010000 /* bit 16 */
	IXGBE_HLREG0_CONTMDC      = 0x00020000 /* bit 17 */
	IXGBE_HLREG0_CTRLFLTR     = 0x00040000 /* bit 18 */
	IXGBE_HLREG0_PREPEND      = 0x00F00000 /* bits 20-23 */
	IXGBE_HLREG0_PRIPAUSEEN   = 0x01000000 /* bit 24 */
	IXGBE_HLREG0_RXPAUSERECDA = 0x06000000 /* bits 25-26 */
	IXGBE_HLREG0_RXLNGTHERREN = 0x08000000 /* bit 27 */
	IXGBE_HLREG0_RXPADSTRIPEN = 0x10000000 /* bit 28 */
)

/* RDRXCTL Bit Masks */
const (
	IXGBE_RDRXCTL_RDMTS_1_2    = 0x00000000 /* Rx Desc Min THLD Size */
	IXGBE_RDRXCTL_CRCSTRIP     = 0x00000002 /* CRC Strip */
	IXGBE_RDRXCTL_PSP          = 0x00000004 /* Pad Small Packet */
	IXGBE_RDRXCTL_MVMEN        = 0x00000020
	IXGBE_RDRXCTL_RSC_PUSH_DIS = 0x00000020
	IXGBE_RDRXCTL_DMAIDONE     = 0x00000008 /* DMA init cycle done */
	IXGBE_RDRXCTL_RSC_PUSH     = 0x00000080
	IXGBE_RDRXCTL_AGGDIS       = 0x00010000 /* Aggregation disable */
	IXGBE_RDRXCTL_RSCFRSTSIZE  = 0x003E0000 /* RSC First packet size */
	IXGBE_RDRXCTL_RSCLLIDIS    = 0x00800000 /* Disable RSC compl on LLI*/
	IXGBE_RDRXCTL_RSCACKC      = 0x02000000 /* must set 1 when RSC ena */
	IXGBE_RDRXCTL_FCOE_WRFIX   = 0x04000000 /* must set 1 when RSC ena */
	IXGBE_RDRXCTL_MBINTEN      = 0x10000000
	IXGBE_RDRXCTL_MDP_EN       = 0x20000000
)
const (
	IXGBE_FCTRL = 0x05080

	IXGBE_FCTRL_SBP  = 0x00000002 /* Store Bad Packet */
	IXGBE_FCTRL_MPE  = 0x00000100 /* Multicast Promiscuous Ena*/
	IXGBE_FCTRL_UPE  = 0x00000200 /* Unicast Promiscuous Ena */
	IXGBE_FCTRL_BAM  = 0x00000400 /* Broadcast Accept Mode */
	IXGBE_FCTRL_PMCF = 0x00001000 /* Pass MAC Control Frames */
	IXGBE_FCTRL_DPF  = 0x00002000 /* Discard Pause Frame */
	/* Receive Priority Flow Control Enable */
	IXGBE_FCTRL_RPFCE       = 0x00004000
	IXGBE_FCTRL_RFCE        = 0x00008000 /* Receive Flow Control Ena */
	IXGBE_MFLCN_PMCF        = 0x00000001 /* Pass MAC Control Frames */
	IXGBE_MFLCN_DPF         = 0x00000002 /* Discard Pause Frame */
	IXGBE_MFLCN_RPFCE       = 0x00000004 /* Receive Priority FC Enable */
	IXGBE_MFLCN_RFCE        = 0x00000008 /* Receive FC Enable */
	IXGBE_MFLCN_RPFCE_MASK  = 0x00000FF4 /* Rx Priority FC bitmap mask */
	IXGBE_MFLCN_RPFCE_SHIFT = 4          /* Rx Priority FC bitmap shift */
)

/* SRRCTL bit definitions */
const (
	IXGBE_SRRCTL_BSIZEPKT_SHIFT                     = 10 /* so many KBs */
	IXGBE_SRRCTL_BSIZEHDRSIZE_SHIFT                 = 2  /* 64byte resolution (>> 6) + at bit 8 offset (<< 8) = (<< 2) */
	IXGBE_SRRCTL_RDMTS_SHIFT                        = 22
	IXGBE_SRRCTL_RDMTS_MASK                         = 0x01C00000
	IXGBE_SRRCTL_DROP_EN                            = 0x10000000
	IXGBE_SRRCTL_BSIZEPKT_MASK                      = 0x0000007F
	IXGBE_SRRCTL_BSIZEHDR_MASK                      = 0x00003F00
	IXGBE_SRRCTL_DESCTYPE_LEGACY                    = 0x00000000
	IXGBE_SRRCTL_DESCTYPE_ADV_ONEBUF                = 0x02000000
	IXGBE_SRRCTL_DESCTYPE_HDR_SPLIT                 = 0x04000000
	IXGBE_SRRCTL_DESCTYPE_HDR_REPLICATION_LARGE_PKT = 0x08000000
	IXGBE_SRRCTL_DESCTYPE_HDR_SPLIT_ALWAYS          = 0x0A000000
	IXGBE_SRRCTL_DESCTYPE_MASK                      = 0x0E000000
)

/* Extended Device Control */
const (
	IXGBE_CTRL_EXT_PFRSTD   = 0x00004000 /* Physical Function Reset Done */
	IXGBE_CTRL_EXT_NS_DIS   = 0x00010000 /* No Snoop disable */
	IXGBE_CTRL_EXT_RO_DIS   = 0x00020000 /* Relaxed Ordering disable */
	IXGBE_CTRL_EXT_DRV_LOAD = 0x10000000 /* Driver loaded bit for FW */
)

/* DCB registers */
const (
	IXGBE_RTRPCS          = 0x02430
	IXGBE_RTTDCS          = 0x04900
	IXGBE_RTTDCS_ARBDIS   = 0x00000040 /* DCB arbiter disable */
	IXGBE_RTTPCS          = 0x0CD00
	IXGBE_RTRUP2TC        = 0x03020
	IXGBE_RTTUP2TC        = 0x0C800
	IXGBE_RTTDQSEL        = 0x04904
	IXGBE_RTTDT1C         = 0x04908
	IXGBE_RTTDT1S         = 0x0490C
	IXGBE_RTTDTECC        = 0x04990
	IXGBE_RTTDTECC_NO_BCN = 0x00000100

	IXGBE_RTTBCNRC              = 0x04984
	IXGBE_RTTBCNRC_RS_ENA       = 0x80000000
	IXGBE_RTTBCNRC_RF_DEC_MASK  = 0x00003FFF
	IXGBE_RTTBCNRC_RF_INT_SHIFT = 14
	IXGBE_RTTBCNRC_RF_INT_MASK  = IXGBE_RTTBCNRC_RF_DEC_MASK << IXGBE_RTTBCNRC_RF_INT_SHIFT
	IXGBE_RTTBCNRM              = 0x04980

	/* BCN (for DCB) Registers */
	IXGBE_RTTBCNRS  = 0x04988
	IXGBE_RTTBCNCR  = 0x08B00
	IXGBE_RTTBCNACH = 0x08B04
	IXGBE_RTTBCNACL = 0x08B08
	IXGBE_RTTBCNTG  = 0x04A90
	IXGBE_RTTBCNIDX = 0x08B0C
	IXGBE_RTTBCNCP  = 0x08B10
	IXGBE_RTFRTIMER = 0x08B14
	IXGBE_RTTBCNRTT = 0x05150
	IXGBE_RTTBCNRD  = 0x0498C
)

func IXGBE_RTRPT4C(i int) int {
	return 0x02140 + i*4 /* 8 of these (0-7) */
}
func IXGBE_TXLLQ(i int) int {
	return 0x082E0 + i*4 /* 4 of these (0-3) */
}
func IXGBE_RTRPT4S(i int) int {
	return 0x02160 + i*4 /* 8 of these (0-7) */
}
func IXGBE_RTTDT2C(i int) int {
	return 0x04910 + i*4 /* 8 of these (0-7) */
}
func IXGBE_RTTDT2S(i int) int {
	return 0x04930 + i*4 /* 8 of these (0-7) */
}
func IXGBE_RTTPT2C(i int) int {
	return 0x0CD20 + i*4 /* 8 of these (0-7) */
}
func IXGBE_RTTPT2S(i int) int {
	return 0x0CD40 + i*4 /* 8 of these (0-7) */
}

/* Interrupt Registers */
const (
	IXGBE_EICR = 0x00800
	IXGBE_EICS = 0x00808
	IXGBE_EIMS = 0x00880
	IXGBE_EIMC = 0x00888
	IXGBE_EIAC = 0x00810
	IXGBE_EIAM = 0x00890
)

func IXGBE_EICS_EX(i int) int {
	return 0x00A90 + i*4
}
func IXGBE_EIMS_EX(i int) int {
	return 0x00AA0 + i*4
}
func IXGBE_EIMC_EX(i int) int {
	return 0x00AB0 + i*4
}
func IXGBE_EIAM_EX(i int) int {
	return 0x00AD0 + i*4
}

/* CTRL Bit Masks */
const (
	IXGBE_CTRL_GIO_DIS  = 0x00000004 /* Global IO Master Disable bit */
	IXGBE_CTRL_LNK_RST  = 0x00000008 /* Link Reset. Resets everything. */
	IXGBE_CTRL_RST      = 0x04000000 /* Reset (SW) */
	IXGBE_CTRL_RST_MASK = IXGBE_CTRL_LNK_RST | IXGBE_CTRL_RST
)

/* EEC Register */
const (
	IXGBE_EEC_SK        = 0x00000001 /* EEPROM Clock */
	IXGBE_EEC_CS        = 0x00000002 /* EEPROM Chip Select */
	IXGBE_EEC_DI        = 0x00000004 /* EEPROM Data In */
	IXGBE_EEC_DO        = 0x00000008 /* EEPROM Data Out */
	IXGBE_EEC_FWE_MASK  = 0x00000030 /* FLASH Write Enable */
	IXGBE_EEC_FWE_DIS   = 0x00000010 /* Disable FLASH writes */
	IXGBE_EEC_FWE_EN    = 0x00000020 /* Enable FLASH writes */
	IXGBE_EEC_FWE_SHIFT = 4
	IXGBE_EEC_REQ       = 0x00000040 /* EEPROM Access Request */
	IXGBE_EEC_GNT       = 0x00000080 /* EEPROM Access Grant */
	IXGBE_EEC_PRES      = 0x00000100 /* EEPROM Present */
	IXGBE_EEC_ARD       = 0x00000200 /* EEPROM Auto Read Done */
	IXGBE_EEC_FLUP      = 0x00800000 /* Flash update command */
	IXGBE_EEC_SEC1VAL   = 0x02000000 /* Sector 1 Valid */
	IXGBE_EEC_FLUDONE   = 0x04000000 /* Flash update done */
)

/* LINKS Bit Masks */
const (
	IXGBE_LINKS_KX_AN_COMP    = 0x80000000
	IXGBE_LINKS_UP            = 0x40000000
	IXGBE_LINKS_SPEED         = 0x20000000
	IXGBE_LINKS_MODE          = 0x18000000
	IXGBE_LINKS_RX_MODE       = 0x06000000
	IXGBE_LINKS_TX_MODE       = 0x01800000
	IXGBE_LINKS_XGXS_EN       = 0x00400000
	IXGBE_LINKS_SGMII_EN      = 0x02000000
	IXGBE_LINKS_PCS_1G_EN     = 0x00200000
	IXGBE_LINKS_1G_AN_EN      = 0x00100000
	IXGBE_LINKS_KX_AN_IDLE    = 0x00080000
	IXGBE_LINKS_1G_SYNC       = 0x00040000
	IXGBE_LINKS_10G_ALIGN     = 0x00020000
	IXGBE_LINKS_10G_LANE_SYNC = 0x00017000
	IXGBE_LINKS_TL_FAULT      = 0x00001000
	IXGBE_LINKS_SIGNAL        = 0x00000F00

	IXGBE_LINKS_SPEED_NON_STD     = 0x08000000
	IXGBE_LINKS_SPEED_82599       = 0x30000000
	IXGBE_LINKS_SPEED_10G_82599   = 0x30000000
	IXGBE_LINKS_SPEED_1G_82599    = 0x20000000
	IXGBE_LINKS_SPEED_100_82599   = 0x10000000
	IXGBE_LINKS_SPEED_10_X550EM_A = 0x00000000
	IXGBE_LINK_UP_TIME            = 90 /* 9.0 Seconds */
	IXGBE_AUTO_NEG_TIME           = 45 /* 4.5 Seconds */

	IXGBE_LINKS2_AN_SUPPORTED = 0x00000040
)

/* Stats registers */
const (
	IXGBE_CRCERRS    = 0x04000
	IXGBE_ILLERRC    = 0x04004
	IXGBE_ERRBC      = 0x04008
	IXGBE_MSPDC      = 0x04010
	IXGBE_MLFC       = 0x04034
	IXGBE_MRFC       = 0x04038
	IXGBE_RLEC       = 0x04040
	IXGBE_LXONTXC    = 0x03F60
	IXGBE_LXONRXC    = 0x0CF60
	IXGBE_LXOFFTXC   = 0x03F68
	IXGBE_LXOFFRXC   = 0x0CF68
	IXGBE_LXONRXCNT  = 0x041A4
	IXGBE_LXOFFRXCNT = 0x041A8
	IXGBE_PRC64      = 0x0405C
	IXGBE_PRC127     = 0x04060
	IXGBE_PRC255     = 0x04064
	IXGBE_PRC511     = 0x04068
	IXGBE_PRC1023    = 0x0406C
	IXGBE_PRC1522    = 0x04070
	IXGBE_GPRC       = 0x04074
	IXGBE_BPRC       = 0x04078
	IXGBE_MPRC       = 0x0407C
	IXGBE_GPTC       = 0x04080
	IXGBE_GORCL      = 0x04088
	IXGBE_GORCH      = 0x0408C
	IXGBE_GOTCL      = 0x04090
	IXGBE_GOTCH      = 0x04094
	IXGBE_RUC        = 0x040A4
	IXGBE_RFC        = 0x040A8
	IXGBE_ROC        = 0x040AC
	IXGBE_RJC        = 0x040B0
	IXGBE_MNGPRC     = 0x040B4
	IXGBE_MNGPDC     = 0x040B8
	IXGBE_MNGPTC     = 0x0CF90
	IXGBE_TORL       = 0x040C0
	IXGBE_TORH       = 0x040C4
	IXGBE_TPR        = 0x040D0
	IXGBE_TPT        = 0x040D4
	IXGBE_PTC64      = 0x040D8
	IXGBE_PTC127     = 0x040DC
	IXGBE_PTC255     = 0x040E0
	IXGBE_PTC511     = 0x040E4
	IXGBE_PTC1023    = 0x040E8
	IXGBE_PTC1522    = 0x040EC
	IXGBE_MPTC       = 0x040F0
	IXGBE_BPTC       = 0x040F4
	IXGBE_XEC        = 0x04120
	IXGBE_SSVPC      = 0x08780

	IXGBE_FCCRC           = 0x05118    /* Num of Good Eth CRC w/ Bad FC CRC */
	IXGBE_FCOERPDC        = 0x0241C    /* FCoE Rx Packets Dropped Count */
	IXGBE_FCLAST          = 0x02424    /* FCoE Last Error Count */
	IXGBE_FCOEPRC         = 0x02428    /* Number of FCoE Packets Received */
	IXGBE_FCOEDWRC        = 0x0242C    /* Number of FCoE DWords Received */
	IXGBE_FCOEPTC         = 0x08784    /* Number of FCoE Packets Transmitted */
	IXGBE_FCOEDWTC        = 0x08788    /* Number of FCoE DWords Transmitted */
	IXGBE_FCCRC_CNT_MASK  = 0x0000FFFF /* CRC_CNT: bit 0 - 15 */
	IXGBE_FCLAST_CNT_MASK = 0x0000FFFF /* Last_CNT: bit 0 - 15 */
	IXGBE_O2BGPTC         = 0x041C4
	IXGBE_O2BSPC          = 0x087B0
	IXGBE_B2OSPC          = 0x041C0
	IXGBE_B2OGPRC         = 0x02F90
	IXGBE_BUPRC           = 0x04180
	IXGBE_BMPRC           = 0x04184
	IXGBE_BBPRC           = 0x04188
	IXGBE_BUPTC           = 0x0418C
	IXGBE_BMPTC           = 0x04190
	IXGBE_BBPTC           = 0x04194
	IXGBE_BCRCERRS        = 0x04198
	IXGBE_BXONRXC         = 0x0419C
	IXGBE_BXOFFRXC        = 0x041E0
	IXGBE_BXONTXC         = 0x041E4
	IXGBE_BXOFFTXC        = 0x041E8
)

func IXGBE_MPC(i int) int {
	return 0x03FA0 + i*4 /* 8 of these 3FA0-3FBC*/
}
func IXGBE_PXONRXCNT(i int) int {
	return 0x04140 + i*4 /* 8 of these */
}
func IXGBE_PXOFFRXCNT(i int) int {
	return 0x04160 + i*4 /* 8 of these */
}
func IXGBE_PXON2OFFCNT(i int) int {
	return 0x03240 + i*4 /* 8 of these */
}
func IXGBE_PXONTXC(i int) int {
	return 0x03F00 + i*4 /* 8 of these 3F00-3F1C*/
}
func IXGBE_PXONRXC(i int) int {
	return 0x0CF00 + i*4 /* 8 of these CF00-CF1C*/
}
func IXGBE_PXOFFTXC(i int) int {
	return 0x03F20 + i*4 /* 8 of these 3F20-3F3C*/
}
func IXGBE_PXOFFRXC(i int) int {
	return 0x0CF20 + i*4 /* 8 of these CF20-CF3C*/
}
func IXGBE_RNBC(i int) int {
	return 0x03FC0 + i*4 /* 8 of these 3FC0-3FDC*/
}
func IXGBE_RQSMR(i int) int {
	return 0x02300 + i*4
}
func IXGBE_TQSMR(i int) int {
	if i < 7 {
		return 0x07300 + i*4
	}
	return 0x08600 + i*4
}
func IXGBE_TQSM(i int) int {
	return 0x08600 + i*4
}
func IXGBE_QPRC(i int) int {
	return 0x1030 + i*0x40 /* 16 of these */
}
func IXGBE_QPTC(i int) int {
	return 0x06030 + i*0x40 /* 16 of these */
}
func IXGBE_QBRC(i int) int {
	return 0x01034 + i*0x40 /* 16 of these */
}
func IXGBE_QBTC(i int) int {
	return 0x06032 + i*0x40 /* 16 of these */
}
func IXGBE_QBRC_L(i int) int {
	return 0x01034 + i*0x40 /* 16 of these */
}
func IXGBE_QBRC_H(i int) int {
	return 0x01038 + i*0x40 /* 16 of these */
}
func IXGBE_QPRDC(i int) int {
	return 0x0143 + i*0x40 /* 16 of these */
}
func IXGBE_QBTC_L(i int) int {
	return 0x08700 + i*0x8 /* 16 of these */
}
func IXGBE_QBTC_H(i int) int {
	return 0x08704 + i*0x8 /* 16 of these */
}

/* Receive Descriptor bit definitions */
const (
	IXGBE_RXD_STAT_DD          = 0x01       /* Descriptor Done */
	IXGBE_RXD_STAT_EOP         = 0x02       /* End of Packet */
	IXGBE_RXD_STAT_FLM         = 0x04       /* FDir Match */
	IXGBE_RXD_STAT_VP          = 0x08       /* IEEE VLAN Packet */
	IXGBE_RXDADV_NEXTP_MASK    = 0x000FFFF0 /* Next Descriptor Index */
	IXGBE_RXDADV_NEXTP_SHIFT   = 0x00000004
	IXGBE_RXD_STAT_UDPCS       = 0x10       /* UDP xsum calculated */
	IXGBE_RXD_STAT_L4CS        = 0x20       /* L4 xsum calculated */
	IXGBE_RXD_STAT_IPCS        = 0x40       /* IP xsum calculated */
	IXGBE_RXD_STAT_PIF         = 0x80       /* passed in-exact filter */
	IXGBE_RXD_STAT_CRCV        = 0x100      /* Speculative CRC Valid */
	IXGBE_RXD_STAT_OUTERIPCS   = 0x100      /* Cloud IP xsum calculated */
	IXGBE_RXD_STAT_VEXT        = 0x200      /* 1st VLAN found */
	IXGBE_RXD_STAT_UDPV        = 0x400      /* Valid UDP checksum */
	IXGBE_RXD_STAT_DYNINT      = 0x800      /* Pkt caused INT via DYNINT */
	IXGBE_RXD_STAT_LLINT       = 0x800      /* Pkt caused Low Latency Interrupt */
	IXGBE_RXD_STAT_TSIP        = 0x08000    /* Time Stamp in packet buffer */
	IXGBE_RXD_STAT_TS          = 0x10000    /* Time Stamp */
	IXGBE_RXD_STAT_SECP        = 0x20000    /* Security Processing */
	IXGBE_RXD_STAT_LB          = 0x40000    /* Loopback Status */
	IXGBE_RXD_STAT_ACK         = 0x8000     /* ACK Packet indication */
	IXGBE_RXD_ERR_CE           = 0x01       /* CRC Error */
	IXGBE_RXD_ERR_LE           = 0x02       /* Length Error */
	IXGBE_RXD_ERR_PE           = 0x08       /* Packet Error */
	IXGBE_RXD_ERR_OSE          = 0x10       /* Oversize Error */
	IXGBE_RXD_ERR_USE          = 0x20       /* Undersize Error */
	IXGBE_RXD_ERR_TCPE         = 0x40       /* TCP/UDP Checksum Error */
	IXGBE_RXD_ERR_IPE          = 0x80       /* IP Checksum Error */
	IXGBE_RXDADV_ERR_MASK      = 0xfff00000 /* RDESC.ERRORS mask */
	IXGBE_RXDADV_ERR_SHIFT     = 20         /* RDESC.ERRORS shift */
	IXGBE_RXDADV_ERR_OUTERIPER = 0x04000000 /* CRC IP Header error */
	IXGBE_RXDADV_ERR_RXE       = 0x20000000 /* Any MAC Error */
	IXGBE_RXDADV_ERR_FCEOFE    = 0x80000000 /* FCEOFe/IPE */
	IXGBE_RXDADV_ERR_FCERR     = 0x00700000 /* FCERR/FDIRERR */
	IXGBE_RXDADV_ERR_FDIR_LEN  = 0x00100000 /* FDIR Length error */
	IXGBE_RXDADV_ERR_FDIR_DROP = 0x00200000 /* FDIR Drop error */
	IXGBE_RXDADV_ERR_FDIR_COLL = 0x00400000 /* FDIR Collision error */
	IXGBE_RXDADV_ERR_HBO       = 0x00800000 /*Header Buffer Overflow */
	IXGBE_RXDADV_ERR_CE        = 0x01000000 /* CRC Error */
	IXGBE_RXDADV_ERR_LE        = 0x02000000 /* Length Error */
	IXGBE_RXDADV_ERR_PE        = 0x08000000 /* Packet Error */
	IXGBE_RXDADV_ERR_OSE       = 0x10000000 /* Oversize Error */
	IXGBE_RXDADV_ERR_USE       = 0x20000000 /* Undersize Error */
	IXGBE_RXDADV_ERR_TCPE      = 0x40000000 /* TCP/UDP Checksum Error */
	IXGBE_RXDADV_ERR_IPE       = 0x80000000 /* IP Checksum Error */
	IXGBE_RXD_VLAN_ID_MASK     = 0x0FFF     /* VLAN ID is in lower 12 bits */
	IXGBE_RXD_PRI_MASK         = 0xE000     /* Priority is in upper 3 bits */
	IXGBE_RXD_PRI_SHIFT        = 13
	IXGBE_RXD_CFI_MASK         = 0x1000 /* CFI is bit 12 */
	IXGBE_RXD_CFI_SHIFT        = 12

	IXGBE_RXDADV_STAT_DD            = IXGBE_RXD_STAT_DD  /* Done */
	IXGBE_RXDADV_STAT_EOP           = IXGBE_RXD_STAT_EOP /* End of Packet */
	IXGBE_RXDADV_STAT_FLM           = IXGBE_RXD_STAT_FLM /* FDir Match */
	IXGBE_RXDADV_STAT_VP            = IXGBE_RXD_STAT_VP  /* IEEE VLAN Pkt */
	IXGBE_RXDADV_STAT_MASK          = 0x000fffff         /* Stat/NEXTP: bit 0-19 */
	IXGBE_RXDADV_STAT_FCEOFS        = 0x00000040         /* FCoE EOF/SOF Stat */
	IXGBE_RXDADV_STAT_FCSTAT        = 0x00000030         /* FCoE Pkt Stat */
	IXGBE_RXDADV_STAT_FCSTAT_NOMTCH = 0x00000000         /* 00: No Ctxt Match */
	IXGBE_RXDADV_STAT_FCSTAT_NODDP  = 0x00000010         /* 01: Ctxt w/o DDP */
	IXGBE_RXDADV_STAT_FCSTAT_FCPRSP = 0x00000020         /* 10: Recv. FCP_RSP */
	IXGBE_RXDADV_STAT_FCSTAT_DDP    = 0x00000030         /* 11: Ctxt w/ DDP */
	IXGBE_RXDADV_STAT_TS            = 0x00010000         /* IEEE1588 Time Stamp */
	IXGBE_RXDADV_STAT_TSIP          = 0x00008000         /* Time Stamp in packet buffer */
)

const (
	IXGBE_TXD_POPTS_IXSM = 0x01       /* Insert IP checksum */
	IXGBE_TXD_POPTS_TXSM = 0x02       /* Insert TCP/UDP checksum */
	IXGBE_TXD_CMD_EOP    = 0x01000000 /* End of Packet */
	IXGBE_TXD_CMD_IFCS   = 0x02000000 /* Insert FCS (Ethernet CRC) */
	IXGBE_TXD_CMD_IC     = 0x04000000 /* Insert Checksum */
	IXGBE_TXD_CMD_RS     = 0x08000000 /* Report Status */
	IXGBE_TXD_CMD_DEXT   = 0x20000000 /* Desc extension (0 = legacy) */
	IXGBE_TXD_CMD_VLE    = 0x40000000 /* Add VLAN tag */
	IXGBE_TXD_STAT_DD    = 0x00000001 /* Descriptor Done */
)

/* Adv Transmit Descriptor Config Masks */
const (
	IXGBE_ADVTXD_DTALEN_MASK         = 0x0000FFFF         /* Data buf length(bytes) */
	IXGBE_ADVTXD_MAC_LINKSEC         = 0x00040000         /* Insert LinkSec */
	IXGBE_ADVTXD_MAC_TSTAMP          = 0x00080000         /* IEEE1588 time stamp */
	IXGBE_ADVTXD_IPSEC_SA_INDEX_MASK = 0x000003FF         /* IPSec SA index */
	IXGBE_ADVTXD_IPSEC_ESP_LEN_MASK  = 0x000001FF         /* IPSec ESP length */
	IXGBE_ADVTXD_DTYP_MASK           = 0x00F00000         /* DTYP mask */
	IXGBE_ADVTXD_DTYP_CTXT           = 0x00200000         /* Adv Context Desc */
	IXGBE_ADVTXD_DTYP_DATA           = 0x00300000         /* Adv Data Descriptor */
	IXGBE_ADVTXD_DCMD_EOP            = IXGBE_TXD_CMD_EOP  /* End of Packet */
	IXGBE_ADVTXD_DCMD_IFCS           = IXGBE_TXD_CMD_IFCS /* Insert FCS */
	IXGBE_ADVTXD_DCMD_RS             = IXGBE_TXD_CMD_RS   /* Report Status */
	IXGBE_ADVTXD_DCMD_DDTYP_ISCSI    = 0x10000000         /* DDP hdr type or iSCSI */
	IXGBE_ADVTXD_DCMD_DEXT           = IXGBE_TXD_CMD_DEXT /* Desc ext 1=Adv */
	IXGBE_ADVTXD_DCMD_VLE            = IXGBE_TXD_CMD_VLE  /* VLAN pkt enable */
	IXGBE_ADVTXD_DCMD_TSE            = 0x80000000         /* TCP Seg enable */
	IXGBE_ADVTXD_STAT_DD             = IXGBE_TXD_STAT_DD  /* Descriptor Done */
	IXGBE_ADVTXD_STAT_SN_CRC         = 0x00000002         /* NXTSEQ/SEED pres in WB */
	IXGBE_ADVTXD_STAT_RSV            = 0x0000000C         /* STA Reserved */
	IXGBE_ADVTXD_IDX_SHIFT           = 4                  /* Adv desc Index shift */
	IXGBE_ADVTXD_CC                  = 0x00000080         /* Check Context */
	IXGBE_ADVTXD_POPTS_SHIFT         = 8                  /* Adv desc POPTS shift */
	IXGBE_ADVTXD_POPTS_IXSM          = IXGBE_TXD_POPTS_IXSM << IXGBE_ADVTXD_POPTS_SHIFT
	IXGBE_ADVTXD_POPTS_TXSM          = IXGBE_TXD_POPTS_TXSM << IXGBE_ADVTXD_POPTS_SHIFT
	IXGBE_ADVTXD_POPTS_ISCO_1ST      = 0x00000000 /* 1st TSO of iSCSI PDU */
	IXGBE_ADVTXD_POPTS_ISCO_MDL      = 0x00000800 /* Middle TSO of iSCSI PDU */
	IXGBE_ADVTXD_POPTS_ISCO_LAST     = 0x00001000 /* Last TSO of iSCSI PDU */
	/* 1st&Last TSO-full iSCSI PDU */
	IXGBE_ADVTXD_POPTS_ISCO_FULL        = 0x00001800
	IXGBE_ADVTXD_POPTS_RSV              = 0x00002000       /* POPTS Reserved */
	IXGBE_ADVTXD_PAYLEN_SHIFT           = 14               /* Adv desc PAYLEN shift */
	IXGBE_ADVTXD_MACLEN_SHIFT           = 9                /* Adv ctxt desc mac len shift */
	IXGBE_ADVTXD_VLAN_SHIFT             = 16               /* Adv ctxt vlan tag shift */
	IXGBE_ADVTXD_TUCMD_IPV4             = 0x00000400       /* IP Packet Type: 1=IPv4 */
	IXGBE_ADVTXD_TUCMD_IPV6             = 0x00000000       /* IP Packet Type: 0=IPv6 */
	IXGBE_ADVTXD_TUCMD_L4T_UDP          = 0x00000000       /* L4 Packet TYPE of UDP */
	IXGBE_ADVTXD_TUCMD_L4T_TCP          = 0x00000800       /* L4 Packet TYPE of TCP */
	IXGBE_ADVTXD_TUCMD_L4T_SCTP         = 0x00001000       /* L4 Packet TYPE of SCTP */
	IXGBE_ADVTXD_TUCMD_L4T_RSV          = 0x00001800       /* RSV L4 Packet TYPE */
	IXGBE_ADVTXD_TUCMD_MKRREQ           = 0x00002000       /* req Markers and CRC */
	IXGBE_ADVTXD_POPTS_IPSEC            = 0x00000400       /* IPSec offload request */
	IXGBE_ADVTXD_TUCMD_IPSEC_TYPE_ESP   = 0x00002000       /* IPSec Type ESP */
	IXGBE_ADVTXD_TUCMD_IPSEC_ENCRYPT_EN = 0x00004000       /* ESP Encrypt Enable */
	IXGBE_ADVTXT_TUCMD_FCOE             = 0x00008000       /* FCoE Frame Type */
	IXGBE_ADVTXD_FCOEF_EOF_MASK         = (0x3 << 10)      /* FC EOF index */
	IXGBE_ADVTXD_FCOEF_SOF              = ((1 << 2) << 10) /* FC SOF index */
	IXGBE_ADVTXD_FCOEF_PARINC           = ((1 << 3) << 10) /* Rel_Off in F_CTL */
	IXGBE_ADVTXD_FCOEF_ORIE             = ((1 << 4) << 10) /* Orientation End */
	IXGBE_ADVTXD_FCOEF_ORIS             = ((1 << 5) << 10) /* Orientation Start */
	IXGBE_ADVTXD_FCOEF_EOF_N            = (0x0 << 10)      /* 00: EOFn */
	IXGBE_ADVTXD_FCOEF_EOF_T            = (0x1 << 10)      /* 01: EOFt */
	IXGBE_ADVTXD_FCOEF_EOF_NI           = (0x2 << 10)      /* 10: EOFni */
	IXGBE_ADVTXD_FCOEF_EOF_A            = (0x3 << 10)      /* 11: EOFa */
	IXGBE_ADVTXD_L4LEN_SHIFT            = 8                /* Adv ctxt L4LEN shift */
	IXGBE_ADVTXD_MSS_SHIFT              = 16               /* Adv ctxt MSS shift */

	IXGBE_ADVTXD_OUTER_IPLEN       = 16 /* Adv ctxt OUTERIPLEN shift */
	IXGBE_ADVTXD_TUNNEL_LEN        = 24 /* Adv ctxt TUNNELLEN shift */
	IXGBE_ADVTXD_TUNNEL_TYPE_SHIFT = 16 /* Adv Tx Desc Tunnel Type shift */
	IXGBE_ADVTXD_OUTERIPCS_SHIFT   = 17 /* Adv Tx Desc OUTERIPCS Shift */
	IXGBE_ADVTXD_TUNNEL_TYPE_NVGRE = 1  /* Adv Tx Desc Tunnel Type NVGRE */
	/* Adv Tx Desc OUTERIPCS Shift for X550EM_a */
	IXGBE_ADVTXD_OUTERIPCS_SHIFT_X550EM_a = 26
	IXGBE_LINK_SPEED_UNKNOWN              = 0
	IXGBE_LINK_SPEED_10_FULL              = 0x0002
	IXGBE_LINK_SPEED_100_FULL             = 0x0008
	IXGBE_LINK_SPEED_1GB_FULL             = 0x0020
	IXGBE_LINK_SPEED_2_5GB_FULL           = 0x0400
	IXGBE_LINK_SPEED_5GB_FULL             = 0x0800
	IXGBE_LINK_SPEED_10GB_FULL            = 0x0080
	IXGBE_LINK_SPEED_82598_AUTONEG        = IXGBE_LINK_SPEED_1GB_FULL | IXGBE_LINK_SPEED_10GB_FULL
	IXGBE_LINK_SPEED_82599_AUTONEG        = IXGBE_LINK_SPEED_100_FULL | IXGBE_LINK_SPEED_1GB_FULL | IXGBE_LINK_SPEED_10GB_FULL
)

/* Statistics Registers */
const (
	IXGBE_RXNFGPC     = 0x041B0
	IXGBE_RXNFGBCL    = 0x041B4
	IXGBE_RXNFGBCH    = 0x041B8
	IXGBE_RXDGPC      = 0x02F50
	IXGBE_RXDGBCL     = 0x02F54
	IXGBE_RXDGBCH     = 0x02F58
	IXGBE_RXDDGPC     = 0x02F5C
	IXGBE_RXDDGBCL    = 0x02F60
	IXGBE_RXDDGBCH    = 0x02F64
	IXGBE_RXLPBKGPC   = 0x02F68
	IXGBE_RXLPBKGBCL  = 0x02F6C
	IXGBE_RXLPBKGBCH  = 0x02F70
	IXGBE_RXDLPBKGPC  = 0x02F74
	IXGBE_RXDLPBKGBCL = 0x02F78
	IXGBE_RXDLPBKGBCH = 0x02F7C
	IXGBE_TXDGPC      = 0x087A0
	IXGBE_TXDGBCL     = 0x087A4
	IXGBE_TXDGBCH     = 0x087A8
)

//enum: const ( a = iota \n b = iota ...)

//"As Go doesn't have support for C's union type in the general case, C's union types are represented as a Go byte array with the same length."
//https://golang.org/cmd/cgo/
//instead we define functions on the types

/* Receive Descriptor - Advanced */
/*union ixgbe_adv_rx_desc {
	struct {
		__le64 pkt_addr;  //Packet buffer address
		__le64 hdr_addr;  //Header buffer address
	} read;
	struct {
		struct {
			union {
				__le32 data;
				struct {
					__le16 pkt_info;  //RSS, Pkt type
					__le16 hdr_info;  //Splithdr, hdrlen
				} hs_rss;
			} lo_dword;
			union {
				__le32 rss;  //RSS Hash
				struct {
					__le16 ip_id;  //IP id
					__le16 csum;  //Packet Checksum
				} csum_ip;
			} hi_dword;
		} lower;
		struct {
			__le32 status_error;  //ext status/error
			__le16 length;  //Packet length
			__le16 vlan;  //VLAN tag
		} upper;
	} wb;   writeback
};*/
type IxgbeAdvRxDesc struct {
	raw []byte
}

/*
// we set the read struct
func (desc *IxgbeAdvRxDesc) read_pktAddr(pktAddr uint64) {
	if isBig {
		binary.BigEndian.PutUint64(desc.raw[:8], pktAddr)
	} else {
		binary.LittleEndian.PutUint64(desc.raw[:8], pktAddr)
	}
}

func (desc *IxgbeAdvRxDesc) read_hdrAddr(hdrAddr uint64) {
	if isBig {
		binary.BigEndian.PutUint64(desc.raw[8:], hdrAddr)
	} else {
		binary.LittleEndian.PutUint64(desc.raw[8:], hdrAddr)
	}
}

// the NIC writes back the descriptors
func (desc *IxgbeAdvRxDesc) wb_pktInfo() uint16 {	//RSS, Pkt type
	if isBig {
		return binary.BigEndian.Uint16(desc.raw[:2])
	}
	return binary.LittleEndian.Uint16(desc.raw[:2])
}

func (desc *IxgbeAdvRxDesc) wb_hdrInfo() uint16 {	 //Splithdr, hdrlen
	if isBig {
		return binary.BigEndian.Uint16(desc.raw[2:4])
	}
	return binary.LittleEndian.Uint16(desc.raw[2:4])
}

func (desc *IxgbeAdvRxDesc) wb_data() uint32 {	 //data
	if isBig {
		return binary.BigEndian.Uint32(desc.raw[0:4])
	}
	return binary.LittleEndian.Uint32(desc.raw[0:4])
}

func (desc *IxgbeAdvRxDesc) wb_rss() uint32 {	 //RSS Hash
	if isBig {
		return binary.BigEndian.Uint32(desc.raw[4:8])
	}
	return binary.LittleEndian.Uint32(desc.raw[4:8])
}

func (desc *IxgbeAdvRxDesc) wb_ipID() uint16 {	 //IP ID
	if isBig {
		return binary.BigEndian.Uint16(desc.raw[4:6])
	}
	return binary.LittleEndian.Uint16(desc.raw[4:6])
}

func (desc *IxgbeAdvRxDesc) wb_csum() uint16 {	 //Packet Checksum
	if isBig {
		return binary.BigEndian.Uint16(desc.raw[6:8])
	}
	return binary.LittleEndian.Uint16(desc.raw[6:8])
}

func (desc *IxgbeAdvRxDesc) wb_statusError() uint32 {	 //ext status/error
	if isBig {
		return binary.BigEndian.Uint32(desc.raw[8:12])
	}
	return binary.LittleEndian.Uint32(desc.raw[8:12])
}

func (desc *IxgbeAdvRxDesc) wb_length() uint16 {	 //Packet length
	if isBig {
		return binary.BigEndian.Uint16(desc.raw[12:14])
	}
	return binary.LittleEndian.Uint16(desc.raw[12:14])
}

func (desc *IxgbeAdvRxDesc) wb_vlan() uint16 {	 //VLAN tag
	if isBig {
		return binary.BigEndian.Uint16(desc.raw[14:])
	}
	return binary.LittleEndian.Uint16(desc.raw[14:])
}
*/

/* Transmit Descriptor - Advanced */
/*union ixgbe_adv_tx_desc {
	struct {
		__le64 buffer_addr; // Address of descriptor's data buf
		__le32 cmd_type_len;
		__le32 olinfo_status;
	} read;
	struct {
		__le64 rsvd; // Reserved
		__le32 nxtseq_seed;
		__le32 status;
	} wb;
};*/
type IxgbeAdvTxDesc struct {
	raw []byte
}

/*func (desc *IxgbeAdvTxDesc) read_bufferAddr(bufAddr uint64) {	//Address of descriptor's data buf
	if isBig {
		binary.BigEndian.PutUint64(desc.raw[:8], bufAddr)
	}
	binary.LittleEndian.PutUint64(desc.raw[:8], bufAddr)
}

func (desc *IxgbeAdvTxDesc) read_cmdTypeLen(cmdTypeLen uint32) {
	if isBig {
		binary.BigEndian.PutUint32(desc.raw[8:12], cmdTypeLen)
	}
	binary.LittleEndian.PutUint32(desc.raw[8:12], cmdTypeLen)
}

func (desc *IxgbeAdvTxDesc) read_olinfoStatus(olinfo uint32) {
	if isBig {
		binary.BigEndian.PutUint32(desc.raw[12:], olinfo)
	}
	binary.LittleEndian.PutUint32(desc.raw[12:], olinfo)
}

func (desc *IxgbeAdvTxDesc) wb_rsvd() uint64 {	//reserved
	if isBig {
		return binary.BigEndian.Uint64(desc.raw[:8])
	}
	return binary.LittleEndian.Uint64(desc.raw[:8])
}

func (desc *IxgbeAdvTxDesc) wb_nxtseqSeed() uint32 {
	if isBig {
		return binary.BigEndian.Uint32(desc.raw[8:12])
	}
	return binary.LittleEndian.Uint32(desc.raw[8:12])
}

func (desc *IxgbeAdvTxDesc) wb_status() uint32 {
	if isBig {
		return binary.BigEndian.Uint32(desc.raw[12:])
	}
	return binary.LittleEndian.Uint32(desc.raw[12:])
}*/
