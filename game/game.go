package game

import (
	"errors"
	"time"
)

type player struct {
	name string
}

func (p player) Name() string {
	return p.name
}

type game struct {
	P1          *player
	P2          *player
	timeStarted time.Time
}

func New() game {
	g := game{}
	g.P1 = &player{"Nobody"}
	g.P2 = &player{"Somebody"}
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
