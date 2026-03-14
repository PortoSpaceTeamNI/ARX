package valve

type Valve byte

const (
	Pressurizing Valve = iota
	Vent
	Abort
	Main
	N2OFill
	N2OPurge
	N2Fill
	N2Purge
	N2OQuickDc
	N2QuickDc
)

func (v Valve) ToString() string {
	switch v {
	case Pressurizing:
		return "Pressurizing"
	case Vent:
		return "Vent"
	case Abort:
		return "Abort"
	case Main:
		return "Main"
	case N2OFill:
		return "N2O Fill"
	case N2OPurge:
		return "N2O Purge"
	case N2Fill:
		return "N2 Fill"
	case N2Purge:
		return "N2 Purge"
	case N2OQuickDc:
		return "N2O Quick Disconnect"
	case N2QuickDc:
		return "N2 Quick Disconnect"
	default:
		return "(This valve doesn't exist)"
	}
}
