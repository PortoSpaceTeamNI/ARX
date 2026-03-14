package valvestate

type ValveState byte

const (
	Closed ValveState = iota
	Opened
	Closing
	Opening
	ClosingNotAcked
	OpeningNotAcked
)

func (vs ValveState) ToString() string {
	switch vs {
	case Closed:
		return "Close"
	case Opened:
		return "Open"
	case Closing:
		return "Closing"
	case Opening:
		return "Opening"
	case ClosingNotAcked:
		return "Closing Not Acked"
	case OpeningNotAcked:
		return "Opening Not Acked"
	default:
		return "(This valve state doesn't exist)"
	}
}

/*func (vs ValveState) CanTransitionTo(goal ValveState) bool {
	switch goal {
	case Opened:
		return vs == Closed || vs == Closing

	case Closed:
		return vs == Opened || vs == Opening
	}

	return false
}*/
