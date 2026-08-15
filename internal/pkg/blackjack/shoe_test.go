package blackjack

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestSixDeckShoe(t *testing.T) {
	s := NewShoe(6, rand.New(rand.NewSource(1)))
	if s.Size() != 312 {
		t.Fatalf("size=%d", s.Size())
	}
	counts := map[Card]int{}
	for _, c := range s.Cards() {
		counts[c]++
	}
	for _, suit := range AllSuits {
		for _, rank := range AllRanks {
			if counts[Card{rank, suit}] != 6 {
				t.Fatalf("%s %s count=%d", rank, suit, counts[Card{rank, suit}])
			}
		}
	}
	before := s.Remaining()
	if _, ok := s.Draw(); !ok || s.Remaining() != before-1 {
		t.Fatal("draw did not reduce remaining")
	}
}
func TestShoeDeterministic(t *testing.T) {
	a := NewShoe(6, rand.New(rand.NewSource(42))).Cards()
	b := NewShoe(6, rand.New(rand.NewSource(42))).Cards()
	if !reflect.DeepEqual(a, b) {
		t.Fatal("same seed produced different shoes")
	}
}

func TestSixDeckPenetrationThreshold(t *testing.T) {
	s := NewShoe(6, rand.New(rand.NewSource(7)))
	for i := 0; i < 233; i++ {
		if _, ok := s.Draw(); !ok {
			t.Fatal("shoe exhausted before cut card")
		}
	}
	if s.NeedsShuffle() {
		t.Fatal("cut card reached before 234 cards")
	}
	if _, ok := s.Draw(); !ok || !s.NeedsShuffle() {
		t.Fatal("cut card not reached at 234 cards")
	}
}

func TestCutCardOnlyReshufflesBetweenRounds(t *testing.T) {
	cards := []Card{tc(Two), tc(Ten), tc(Two), tc(Six), tc(Two), tc(Two), tc(Ten), tc(Ten)}
	shoe := NewShoeFromCards(cards)
	g, err := NewGameWithShoe(DefaultRules(), shoe)
	if err != nil {
		t.Fatal(err)
	}
	if err = g.Act(ActionHit); err != nil {
		t.Fatal(err)
	}
	if err = g.Act(ActionHit); err != nil {
		t.Fatal(err)
	}
	if g.Shoe != shoe || !shoe.NeedsShuffle() {
		t.Fatal("shoe reshuffled during round or threshold not reached")
	}
	if err = g.Act(ActionStand); err != nil {
		t.Fatal(err)
	}
	if g.Phase != PhaseRoundComplete {
		t.Fatal(g.Phase)
	}
	if err = g.StartRound(); err != nil {
		t.Fatal(err)
	}
	if g.Shoe == shoe || g.Shoe.Size() != 312 {
		t.Fatal("shoe was not replaced between rounds")
	}
}
