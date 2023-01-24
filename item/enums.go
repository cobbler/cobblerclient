package item

type ResourceAction int64

const (
	raCREATE = iota
	raREMOVE
)

func (r ResourceAction) String() string {
	switch r {
	case raCREATE:
		return "create"
	case raREMOVE:
		return "remove"
	}
	return "unknown"
}

type Architecture int64

const (
	aI386 = iota
	aX8664
	aIA64
	aPPC
	aPPC64
	aPPC64LE
	aPPC64EL
	aS390
	aS390X
	aARM
	aAARCH64
)

func (a Architecture) String() string {
	switch a {
	case aI386:
		return "i386"
	case aX8664:
		return "x86_64"
	case aIA64:
		return "ia64"
	case aPPC:
		return "ppc"
	case aPPC64:
		return "ppc64"
	case aPPC64LE:
		return "ppc64le"
	case aPPC64EL:
		return "ppc64el"
	case aS390:
		return "s390"
	case aS390X:
		return "s390x"
	case aARM:
		return "arm"
	case aAARCH64:
		return "aarch64"
	}
	return "unknown"
}

type ImageType int64

const (
	itDIRECT = iota
	itISO
	itMEMDISK
	itVIRTCLONE
)

func (i ImageType) String() string {
	switch i {
	case itDIRECT:
		return "direct"
	case itISO:
		return "iso"
	case itMEMDISK:
		return "memdisk"
	case itVIRTCLONE:
		return "virt-clone"
	}
	return "unknown"
}

type VirtType int64

const (
	itINHERITED = iota
	vtQEMU
	itKVM
	itXENPV
	itXENFV
	itVMWARE
	itVMWAREW
	itOPENVZ
	itAUTO
)

func (v VirtType) String() string {
	switch v {
	case itINHERITED:
		return "<<inherit>>"
	case vtQEMU:
		return "qemu"
	case itKVM:
		return "kvm"
	case itXENPV:
		return "xenpv"
	case itXENFV:
		return "xenfv"
	case itVMWARE:
		return "vmware"
	case itVMWAREW:
		return "vmwarew"
	case itOPENVZ:
		return "openvz"
	case itAUTO:
		return "auto"
	}
	return "unknown"
}

type VirtDiskDriver int64

const (
	vddINHERITED = iota
	vddRAW
	vddQCOW2
	vddQED
	vddVDI
	vddVDMK
)

func (v VirtDiskDriver) String() string {
	switch v {
	case vddINHERITED:
		return "<<inherit>>"
	case vddRAW:
		return "raw"
	case vddQCOW2:
		return "qcow2"
	case vddQED:
		return "qed"
	case vddVDI:
		return "vdi"
	case vddVDMK:
		return "vdmk"
	}
	return "unknown"
}
