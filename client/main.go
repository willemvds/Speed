package main

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/jupiterrider/purego-sdl3/sdl"

	"github.com/willemvds/Speed/game"
)

type Region interface {
	X() int32
	Y() int32
	Width() int32
	Height() int32
}

type Intersectable interface {
	Intersect(Region) bool
}

type Event uint8

const (
	EVENT_GRAB Event = iota
	EVENT_DROP
)

type EventCallback func(event Event)

type Point struct {
	x, y int32
}

func (p Point) X() int32 {
	return p.x
}

func (p Point) Y() int32 {
	return p.y
}

func (p Point) Width() int32 {
	return 1
}

func (p Point) Height() int32 {
	return 1
}

func (p *Point) Intersect(target Region) bool {
	if p.X() >= target.X() && p.X() <= target.X()+target.Width() &&
		p.Y() >= target.Y() && p.Y() <= target.Y()+target.Height() {
		return true
	}
	return false
}

type EventRegion struct {
	Rect          *sdl.Rect
	eventCallback EventCallback
}

func NewEventRegion(x, y, width, height int32) *EventRegion {
	er := EventRegion{Rect: &sdl.Rect{X: x, Y: y, W: width, H: height}}
	return &er
}

func (er *EventRegion) X() int32 {
	return er.Rect.X
}

func (er *EventRegion) Y() int32 {
	return er.Rect.Y
}

func (er *EventRegion) Width() int32 {
	return er.Rect.W
}

func (er *EventRegion) Height() int32 {
	return er.Rect.H
}

func (er *EventRegion) SDLRect() *sdl.FRect {
	return &sdl.FRect{
		X: float32(er.Rect.X),
		Y: float32(er.Rect.Y),
		W: float32(er.Rect.W),
		H: float32(er.Rect.H),
	}
}

func (er *EventRegion) SetEventCallback(cb EventCallback) {
	er.eventCallback = cb
}

func (er *EventRegion) Trigger(event Event) {
	if er.eventCallback != nil {
		er.eventCallback(event)
	}
}

func (er *EventRegion) Intersect(target Region) bool {
	if target.X() > er.X()+er.Width() || target.X()+target.Width() < er.X() {
		return false
	}
	if target.Y() > er.Y()+er.Height() || target.Y()+target.Height() < er.Y() {
		return false
	}
	return true
}

type RegionList []*EventRegion

func (rl RegionList) HitWhat(r Region) *EventRegion {
	for _, region := range rl {
		if intersect := region.Intersect(r); intersect {
			return region
		}
	}
	return nil
}

const CARD_WIDTH = 100
const CARD_HEIGHT = 180

const (
	INTERACTION_STATE_DEFAULT  = 0
	INTERACTION_STATE_DRAGGING = 1
)

/*
Event region list structure:
0..1: Center
2..3: Sides
4..9: Self
10..15: Opponent
*/
func setupEventRegions() *RegionList {
	var idx int = 0
	var x int32
	var y int32 = 280
	eventRegions := make(RegionList, 16)
	// Center stacks
	eventRegions[idx] = NewEventRegion(312, y, CARD_WIDTH, CARD_HEIGHT)
	idx++
	eventRegions[idx] = NewEventRegion(612, y, CARD_WIDTH, CARD_HEIGHT)
	idx++
	// Side stacks
	eventRegions[idx] = NewEventRegion(30, y, CARD_WIDTH, CARD_HEIGHT)
	idx++
	eventRegions[idx] = NewEventRegion(894, y, CARD_WIDTH, CARD_HEIGHT)
	idx++
	// Self stacks
	y = 520
	for i := 0; i < 6; i++ {
		x = int32(80 + (i*CARD_WIDTH + i*50))
		eventRegions[idx] = NewEventRegion(x, y, CARD_WIDTH, CARD_HEIGHT)
		idx++
	}
	// Opponent stacks
	y = 20
	for i := 0; i < 6; i++ {
		x = int32(80 + (i*CARD_WIDTH + i*50))
		eventRegions[idx] = NewEventRegion(x, y, CARD_WIDTH, CARD_HEIGHT)
		idx++
	}
	return &eventRegions
}

func setupEventRegionHandler(g *game.Game, er *EventRegion, typ uint8, idx int) {
	er.SetEventCallback(func(ev Event) {
		switch {
		case ev == EVENT_GRAB:
			g.Grab(g.P1, typ, idx)
		case ev == EVENT_DROP:
			g.Drop(g.P1, typ, idx)
		}
	})
}

func setupEventHandlers(g *game.Game, rl *RegionList) {
	for i := range *rl {
		idx := i
		switch {
		// Center stacks
		case i < 2:
			setupEventRegionHandler(g, (*rl)[i], game.STACK_TYPE_CENTER, idx)
		// Side stacks
		case i < 4:
			setupEventRegionHandler(g, (*rl)[i], game.STACK_TYPE_SIDE, idx-2)
		// Self stacks
		case i < 10:
			setupEventRegionHandler(g, (*rl)[i], game.STACK_TYPE_SELF, idx-4)
		// Opponent stacks
		case i < 16:
			setupEventRegionHandler(g, (*rl)[i], game.STACK_TYPE_OPPONENT, idx-10)
		default:
			panic("I feel like the region list is not tip top")
		}
	}
}

func main() {
	runtime.GOMAXPROCS(2)

	TheDeck := game.NewDeck()
	TheGame := game.New(TheDeck)
	log.Println(TheGame)
	log.Println("[client]", TheGame.Start())

	P1 := game.NewPlayer("Nobody")
	P2 := game.NewPlayer("Somebody")
	err := TheGame.Join(*P1)
	log.Println("[client]", err)
	TheGame.Ready(TheGame.P1)
	err = TheGame.Join(*P2)
	log.Println("[client]", err)
	TheGame.Ready(TheGame.P2)

	log.Println("[client]", TheGame.Start())

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			<-ticker.C
			log.Println("[client] Game Duration: ", TheGame.Duration())
		}
	}()

	interactionState := INTERACTION_STATE_DEFAULT
	// not used? wha?
	log.Println(interactionState)

	var draggingWhat Region
	res := sdl.Init(sdl.InitVideo)
	log.Println("[client]", res)

	eventRegions := setupEventRegions()
	setupEventHandlers(&TheGame, eventRegions)

	var window *sdl.Window
	var renderer *sdl.Renderer
	ok := sdl.CreateWindowAndRenderer("Speed", 1024, 768, sdl.WindowResizable, &window, &renderer)
	if !ok {
		err := sdl.GetError()
		log.Println("[client] Failed to create window:", err)
		os.Exit(1)
	}
	defer sdl.DestroyRenderer(renderer)
	defer sdl.DestroyWindow(window)

	sdl.RenderClear(renderer)

	// Draw an image of a card
	imgCard42 := sdl.LoadBMP("card.bmp")
	if imgCard42 == nil {
		err := sdl.GetError()
		log.Println("[client] Failed to load bitmap (card.bmp):", err)
		os.Exit(1)
	}
	texture := sdl.CreateTextureFromSurface(renderer, imgCard42)
	if texture == nil {
		err := sdl.GetError()
		log.Println("[client] Failed to create texture (42):", err)
		os.Exit(1)
	}
	defer sdl.DestroyTexture(texture)
	dst := sdl.FRect{100, 50, 100, 180}
	sdl.RenderTexture(renderer, texture, nil, &dst)

	// Draw our event regions (for now)
	sdl.SetRenderDrawColor(renderer, 255, 255, 255, 255)
	for _, region := range *eventRegions {
		sdl.RenderRect(renderer, region.SDLRect())
	}
	sdl.SetRenderDrawColor(renderer, 0, 0, 0, 255)

	sdl.RenderPresent(renderer)

	var event sdl.Event
	running := true
	for running {
		for sdl.PollEvent(&event) {
			switch event.Type() {
			case sdl.EventQuit:
				running = false
			/*
				case *sdl.MouseMotionEvent:
				fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n", t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
			*/
			case sdl.EventMouseButtonDown:
				//log.Printf("[client] [%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
				//	t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
				ev := event.Button()
				if ev.Button == 1 {
					point := &Point{int32(ev.X), int32(ev.Y)}
					what := eventRegions.HitWhat(point)
					if ev.Down {
						//log.Println("[client] MOUSE DOWN. Grab something.")
						if what != nil {
							what.Trigger(EVENT_GRAB)
							//log.Println("[client] GRABBING", what)
							interactionState = INTERACTION_STATE_DRAGGING
							draggingWhat = what
						}
					} else {
						//log.Println("[client] MOUSE UP. Drop it.")
						if draggingWhat != nil {
							if what != nil {
								what.Trigger(EVENT_DROP)
							} else {
								TheGame.Discard(TheGame.P1)
							}
							//log.Println("[client] DROPPING", draggingWhat, "ON", what)
							interactionState = INTERACTION_STATE_DEFAULT
							draggingWhat = nil
						}
					}
				}
				/*
					case *sdl.MouseWheelEvent:
						fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n", t.Timestamp, t.Type, t.Which, t.X, t.Y)
				*/
				//			case *sdl.KeyUpEvent:
				//				log.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
				//					t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
			}
		}
	}
}
