package game

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/willemvds/Speed/cards"
)

const (
	STACK_TYPE_CENTER uint8 = iota
	STACK_TYPE_SIDE
	STACK_TYPE_SELF
	STACK_TYPE_OPPONENT
)

var ErrPlayersNotReady = errors.New("Not all players are ready")
var ErrGameAlreadyStarted = errors.New("Game already started")
var ErrNoPlayerSlotsAvailable = errors.New("No more player slots available")
var ErrPlayerNotPresent = errors.New("That player is not present in the game")

type player struct {
	name    string
	holding cards.Card
}

func NewPlayer(name string) *player {
	p := player{}
	p.name = name
	return &p
}

func (p player) Name() string {
	return p.name
}

type gameState uint8

const (
	STATE_PRE_GAME gameState = iota
	STATE_PLAY
	STATE_POST_GAME
)

type Game struct {
	P1           *player
	p1Ready      bool
	P2           *player
	p2Ready      bool
	centerStacks []*cards.Stack
	sideStacks   []*cards.Stack
	p1Stacks     []*cards.Stack
	p2Stacks     []*cards.Stack
	state        gameState
	startedAt    time.Time
}

func New(deck *cards.Stack) Game {
	g := Game{}
	g.state = STATE_PRE_GAME
	//g.P1 = &player{name: "Nobody"}
	//g.P2 = &player{name: "Somebody"}

	//for _, card := range deck.GetCards() {
	//	fmt.Printf("%s, ", card)
	//}
	fmt.Println(".")

	g.p1Stacks = []*cards.Stack{
		cards.NewStack(4),
		cards.NewStack(4),
		cards.NewStack(4),
		cards.NewStack(4),
		cards.NewStack(4),
		cards.NewStack(4),
	}
	g.p2Stacks = []*cards.Stack{
		cards.NewStack(4),
		cards.NewStack(4),
		cards.NewStack(4),
		cards.NewStack(4),
		cards.NewStack(4),
		cards.NewStack(4),
	}
	for i := range 6 {
		for range 4 {
			cardForP1, _ := deck.Pop()
			g.p1Stacks[i].Push(cardForP1)
			cardForP2, _ := deck.Pop()
			g.p2Stacks[i].Push(cardForP2)
		}
	}
	g.centerStacks = []*cards.Stack{
		cards.NewStack(1),
		cards.NewStack(1),
	}
	g.sideStacks = []*cards.Stack{
		cards.NewStack(1),
		cards.NewStack(1),
	}

	fmt.Println("should be out of cards:")
	fmt.Println(deck.Pop())

	return g
}

func (g *Game) Join(p player) error {
	if g.P1 == nil {
		g.P1 = &p
		return nil
	} else if g.P2 == nil {
		g.P2 = &p
		return nil
	} else {
		return ErrNoPlayerSlotsAvailable
	}
}

func (g *Game) Ready(p *player) error {
	if g.P1 == p {
		if g.P1 == nil {
			return ErrPlayerNotPresent
		}
		g.p1Ready = true
		return nil
	}
	if g.P2 == p {
		if g.P2 == nil {
			return ErrPlayerNotPresent
		}
		g.p2Ready = true
		return nil
	}
	return ErrPlayerNotPresent
}

func (g *Game) P1Ready(p *player) error {
	g.p1Ready = true
	return nil
}

func (g *Game) P2Ready(p *player) error {
	g.p2Ready = true
	return nil
}

func (g *Game) Start() error {
	if !g.p1Ready || !g.p2Ready {
		return ErrPlayersNotReady
	}
	if g.startedAt.IsZero() {
		g.startedAt = time.Now()
		g.state = STATE_PLAY
		return nil
	}
	return ErrGameAlreadyStarted
}

func (g *Game) Duration() time.Duration {
	return time.Since(g.startedAt)
}

func (g *Game) Grab(p *player, typ uint8, idx int) {
	if g.state != STATE_PLAY {
		return
	}
	log.Printf("[game] Got GRAB from %s, type=%d, index=%d\n", p, typ, idx)
	if p.holding == cards.Nothing {
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
	if g.state != STATE_PLAY {
		return
	}
	log.Printf("[game] Got DROP from %s, type=%d, index=%d\n", p, typ, idx)
	card := p.holding
	p.holding = cards.Nothing
	if typ == STACK_TYPE_CENTER {
		if card == cards.Nothing {
			log.Println("[game] Player not holding anything, pointless!")
			return
		}
		top, _ := g.centerStacks[idx].Top()
		log.Printf("[game] top=%s, holding=%s, nextto=%b\n", top, p.holding, p.holding.NextTo(&top))
		if p.holding.NextTo(&top) {
			g.centerStacks[idx].Push(p.holding)
			log.Println("Someone actually made a legit move")
		}
	}
	g.CheckWinConditions()
}

func (g *Game) Discard(p *player) {
	if g.state != STATE_PLAY {
		return
	}
	p.holding = cards.Nothing
}

func (g *Game) CheckWinConditions() {
	cardsLeft := 0
	for _, stack := range g.p1Stacks {
		cardsLeft += stack.Size()
	}
	log.Println("Player 1 cards left:", cardsLeft)
	if cardsLeft == 0 {
		log.Println("Player 1 won somehow...")
		g.state = STATE_POST_GAME
		return
	}
	for _, stack := range g.p2Stacks {
		cardsLeft += stack.Size()
	}
	log.Println("Player 2 cards left:", cardsLeft)
	if cardsLeft == 0 {
		log.Println("Player 2 won somehow...")
		g.state = STATE_POST_GAME
		return
	}
}
