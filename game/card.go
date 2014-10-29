package game

import (
	"errors"
)

var ErrNoMoreCards = errors.New("No cards left on stack")
var ErrTooManyCards = errors.New("Too many cards on stack")

type cardrank uint8

func (cr cardrank) Next() cardrank {
	if uint8(cr) < 9 {
		return cr + 1
	}
	return cardrank(0)
}

func (cr cardrank) Prev() cardrank {
	if uint8(cr) > 0 {
		return cr - 1
	}
	return cardrank(9)
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
	return "Hello, I am a Card."
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
