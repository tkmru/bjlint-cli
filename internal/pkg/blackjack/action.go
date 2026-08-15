package blackjack

type Action uint8

const (
	ActionHit Action = iota
	ActionStand
	ActionDouble
	ActionSplit
)

func (a Action) String() string {
	switch a {
	case ActionHit:
		return "HIT"
	case ActionStand:
		return "STAND"
	case ActionDouble:
		return "DOUBLE"
	case ActionSplit:
		return "SPLIT"
	default:
		return "UNKNOWN"
	}
}

type ActionSet uint8

func (s ActionSet) Has(a Action) bool { return s&(1<<a) != 0 }
func (s *ActionSet) Add(a Action)     { *s |= 1 << a }
