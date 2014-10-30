package game

import (
	"errors"
	"fmt"
	"time"
)

type player struct {
	name string
}

func (p player) Name() string {
	return p.name
}

type game struct {
	P1           *player
	P2           *player
	timeStarted  time.Time
	centerStacks []*cardstack
	sideStacks   []*cardstack
	p1Stacks     []*cardstack
	p2Stacks     []*cardstack
}

func New(d *deck) game {
	g := game{}
	g.P1 = &player{"Nobody"}
	g.P2 = &player{"Somebody"}

	for _, card := range d.GetCards() {
		fmt.Printf("%s, ", card)
	}
	fmt.Println(".")

	g.p1Stacks = make([]*cardstack, 6, 6)
	g.p2Stacks = make([]*cardstack, 6, 6)
	for i := 0; i < 6; i++ {
		g.p1Stacks[i] = NewCardStack(4)
		g.p2Stacks[i] = NewCardStack(4)
		for j := 0; j < 4; j++ {
			g.p1Stacks[i].Push(d.GetNextCard())
			g.p2Stacks[i].Push(d.GetNextCard())
		}
	}
	g.centerStacks = make([]*cardstack, 2, 2)
	g.centerStacks[0] = NewCardStack(51)
	g.centerStacks[0].Push(d.GetNextCard())
	g.centerStacks[1] = NewCardStack(51)
	g.centerStacks[1].Push(d.GetNextCard())
	g.sideStacks = make([]*cardstack, 2, 2)
	g.sideStacks[0] = NewCardStack(1)
	g.sideStacks[0].Push(d.GetNextCard())
	g.sideStacks[1] = NewCardStack(1)
	g.sideStacks[1].Push(d.GetNextCard())

	fmt.Println("should be out of cards:", d.GetNextCard())

	return g
}

func (g *game) Start() error {
	if g.timeStarted.IsZero() {
		g.timeStarted = time.Now()
		return nil
	}
	return errors.New("Game already started")
}

func (g *game) Duration() time.Duration {
	return time.Since(g.timeStarted)
}
