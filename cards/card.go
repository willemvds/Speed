package cards

import (
	"errors"
	"fmt"
)

var ErrNoMoreCards = errors.New("no cards left on stack")

type rank uint8

const MinRank = 1
const MaxRank = 13

const WildcardRank = rank(77)

func (r rank) next() rank {
	if uint8(r) < MaxRank {
		return r + 1
	}
	return rank(MinRank)
}

func (cr rank) prev() rank {
	if uint8(cr) > MinRank {
		return cr - 1
	}
	return rank(MaxRank)
}

type Card struct {
	rank rank
}

func NewCard(crank uint8) Card {
	c := Card{}
	c.rank = rank(crank)
	return c
}

var NullCard = NewCard(0)
var Nothing = NewCard(0)

func (c Card) String() string {
	return fmt.Sprintf("Card(%d)", c.rank)
}

func (c Card) NextTo(targetCard *Card) bool {
	if c.rank == targetCard.rank.next() || c.rank == targetCard.rank.prev() {
		return true
	}
	return false
}

type Stack struct {
	cards []Card
}

func NewStack(size int) *Stack {
	cs := Stack{}
	cs.cards = make([]Card, 0, size)
	return &cs
}

func (cs *Stack) Push(c Card) {
	cs.cards = append(cs.cards, c)
}

func (cs *Stack) Pop() (Card, error) {
	if len(cs.cards) > 0 {
		card := cs.cards[len(cs.cards)-1]
		cs.cards = cs.cards[:len(cs.cards)-1]
		return card, nil
	}
	return NullCard, ErrNoMoreCards
}

func (cs *Stack) Top() (Card, error) {
	numCards := len(cs.cards)
	if numCards > 0 {
		return cs.cards[numCards-1], nil
	}
	return Nothing, ErrNoMoreCards
}

func (cs *Stack) Size() int {
	return len(cs.cards)
}

func StandardDeck() *Stack {
	deck := NewStack(52)
	var i uint8
	for i = MinRank; i <= MaxRank; i++ {
		for range 4 {
			deck.Push(NewCard(i))
		}
	}
	fmt.Println("Got ", len(deck.cards), "cards")
	return deck
}
