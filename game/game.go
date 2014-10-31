package game

import (
	"errors"
	"fmt"
	"log"
	"time"
)

const (
	STACK_TYPE_CENTER uint8 = iota
	STACK_TYPE_SIDE
	STACK_TYPE_SELF
	STACK_TYPE_OPPONENT
)

type player struct {
	name string
}

func (p player) Name() string {
	return p.name
}

type Game struct {
	P1           *player
	P2           *player
	timeStarted  time.Time
	centerStacks []*cardstack
	sideStacks   []*cardstack
	p1Stacks     []*cardstack
	p2Stacks     []*cardstack
}

func New(d *deck) Game {
	g := Game{}
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

func (g *Game) Start() error {
	if g.timeStarted.IsZero() {
		g.timeStarted = time.Now()
		return nil
	}
	return errors.New("Game already started")
}

func (g *Game) Duration() time.Duration {
	return time.Since(g.timeStarted)
}

func (g *Game) Grab(p *player, typ uint8, idx int) {
	log.Printf("[game] Got GRAB from %s, type=%d, index=%d\n", p, typ, idx)
}

func (g *Game) Drop(p *player, typ uint8, idx int) {
	log.Printf("[game] Got DROP from %s, type=%d, index=%d\n", p, typ, idx)
}
