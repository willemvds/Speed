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
	name    string
	holding *card
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
	g.P1 = &player{name: "Nobody"}
	g.P2 = &player{name: "Somebody"}

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
	if p.holding != nil {
		log.Println("[game] Player is already holding something, try again? (N/n)")
		return
	}
	if typ == STACK_TYPE_SELF {
		card, err := g.p1Stacks[idx].Pop()
		if err != nil {
			log.Println("[game] No more cards on that stack, peace")
			return
		}
		log.Println("[game] Player is now holding", card)
		p.holding = card
	}
}

func (g *Game) Drop(p *player, typ uint8, idx int) {
	log.Printf("[game] Got DROP from %s, type=%d, index=%d\n", p, typ, idx)
	defer func() { p.holding = nil }()
	if typ == STACK_TYPE_CENTER {
		if p.holding == nil {
			log.Println("[game] Player not holding anything, pointless!")
			return
		}
		top, _ := g.centerStacks[idx].Top()
		log.Printf("[game] top=%s, holding=%s, nextto=%b\n", top, p.holding, p.holding.NextTo(top))
		if p.holding.NextTo(top) {
			g.centerStacks[idx].Push(p.holding)
			log.Println("Someone actually made a legit move")
		}
	}
	g.CheckWinConditions()
}

func (g *Game) Discard(p *player) {
	p.holding = nil
}

func (g *Game) CheckWinConditions() {
	p1CardsLeft := 0
	for _, stack := range g.p1Stacks {
		p1CardsLeft += stack.Len()
	}
	log.Println("Player 1 cards left:", p1CardsLeft)
	if p1CardsLeft == 0 {
		log.Println("Player 1 won somehow...")
	}
}
