package game

import (
	"errors"
	"fmt"
)

var ErrNoMoreCards = errors.New("No cards left on stack")
var ErrTooManyCards = errors.New("Too many cards on stack")

type cardrank uint8

const MinRank = 1
const MaxRank = 13

var WILDCARD = cardrank(77)

func (cr cardrank) Next() cardrank {
	if uint8(cr) < MaxRank {
		return cr + 1
	}
	return cardrank(MinRank)
}

func (cr cardrank) Prev() cardrank {
	if uint8(cr) > MinRank {
		return cr - 1
	}
	return cardrank(MaxRank)
}

type card struct {
	rank cardrank
}

func NewCard(rank uint8) *card {
	c := card{}
	c.rank = cardrank(rank)
	return &c
}

func (c *card) String() string {
	if c == nil {
		return "I appear to be a <nil> '*card'"
	}
	return fmt.Sprintf("Card(%d)", c.rank)
}

func (c *card) NextTo(targetCard *card) bool {
	if c.rank == targetCard.rank.Next() || c.rank == targetCard.rank.Prev() {
		return true
	}
	return false
}

type cardstack struct {
	stack []*card
	size  int
}

func NewCardStack(size int) *cardstack {
	cs := cardstack{}
	cs.stack = make([]*card, 0, size)
	cs.size = size
	return &cs
}

func (cs *cardstack) Push(c *card) error {
	if len(cs.stack) < cs.size {
		cs.stack = append(cs.stack, c)
		return nil
	}
	return ErrNoMoreCards
}

func (cs *cardstack) droplast() {
	cs.stack = cs.stack[:len(cs.stack)-1]
}

func (cs *cardstack) Pop() (*card, error) {
	if len(cs.stack) > 0 {
		defer cs.droplast()
		return cs.stack[len(cs.stack)-1], nil
	}
	return nil, ErrTooManyCards
}

func (cs *cardstack) Top() (*card, error) {
	if len(cs.stack) > 0 {
		return cs.stack[len(cs.stack)-1], nil
	}
	return nil, ErrNoMoreCards
}

func (cs *cardstack) Len() int {
	return len(cs.stack)
}

type deck struct {
	cards []*card
	idx   int
}

func NewDeck() *deck {
	d := deck{}
	d.idx = 0
	d.cards = make([]*card, 0, 52)
	var i uint8
	for i = MinRank; i <= MaxRank; i++ {
		for range 4 {
			d.cards = append(d.cards, NewCard(i))
		}
	}
	d.cards = append(d.cards, NewCard(uint8(WILDCARD)))
	d.cards = append(d.cards, NewCard(uint8(WILDCARD)))
	fmt.Println("Got ", len(d.cards), "cards")
	return &d
}

func (d *deck) GetCards() []*card {
	return d.cards
}

func (d *deck) GetNextCard() *card {
	if d.idx < len(d.cards) {
		defer func() { d.idx++ }()
		return d.cards[d.idx]
	}
	return nil
}
