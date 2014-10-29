package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/willemvds/speed/game"
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

func (er *EventRegion) SDLRect() *sdl.Rect {
	return er.Rect
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

func POINTLESS_EVENT_HANDLER(ev Event) {
	fmt.Println("Apparently this event occured:", ev)
}

func GetEventCallback(id int) EventCallback {
	return func(ev Event) {
		fmt.Printf("[%d] Apparently this event occured: %d\n", id, ev)
	}
}

/*
Planned event region list structure:
0..1: Center
2..3: Sides
4..9: Player 1
10..15: Player 2
*/
func main() {
	runtime.GOMAXPROCS(2)
	TheGame := game.New()
	fmt.Println(TheGame)
	fmt.Println("P1:", TheGame.P1.Name())
	fmt.Println("P2:", TheGame.P2.Name())
	fmt.Println(TheGame.Start())
	fmt.Println(TheGame.Start())
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			<-ticker.C
			fmt.Println("Game Duration: ", TheGame.Duration())
		}
	}()
	interactionState := INTERACTION_STATE_DEFAULT
	fmt.Println(interactionState)
	var draggingWhat Region
	res := sdl.Init(sdl.INIT_VIDEO)
	log.Println(res)

	eventRegions := make(RegionList, 0)
	var x int32
	var y int32 = 20
	for i := 0; i < 8; i++ {
		x = int32(80 + (i*CARD_WIDTH + i*10))
		er := NewEventRegion(x, y, CARD_WIDTH, CARD_HEIGHT)
		er.SetEventCallback(POINTLESS_EVENT_HANDLER)
		eventRegions = append(eventRegions, er)
	}

	y = 520
	for i := 0; i < 8; i++ {
		x = int32(80 + (i*CARD_WIDTH + i*10))
		er := NewEventRegion(x, y, CARD_WIDTH, CARD_HEIGHT)
		er.SetEventCallback(GetEventCallback(i))
		eventRegions = append(eventRegions, er)
	}

	y = 280
	eventRegions = append(eventRegions, NewEventRegion(30, y, CARD_WIDTH, CARD_HEIGHT))
	eventRegions = append(eventRegions, NewEventRegion(312, y, CARD_WIDTH, CARD_HEIGHT))
	eventRegions = append(eventRegions, NewEventRegion(612, y, CARD_WIDTH, CARD_HEIGHT))
	eventRegions = append(eventRegions, NewEventRegion(894, y, CARD_WIDTH, CARD_HEIGHT))

	window := sdl.CreateWindow("Speed", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 1024, 768, sdl.WINDOW_SHOWN)

	renderer := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if renderer == nil {
		log.Println("Failed to create renderer:", sdl.GetError())
		os.Exit(1)
	}

	imgCard42 := sdl.LoadBMP("42.bmp")
	log.Println(imgCard42)

	texture := renderer.CreateTextureFromSurface(imgCard42)
	if texture == nil {
		log.Println("Failed to create texture (42):", sdl.GetError())
	}

	src := sdl.Rect{0, 0, 100, 180}
	dst := sdl.Rect{100, 50, 100, 180}

	renderer.Clear()
	renderer.Copy(texture, &src, &dst)

	renderer.SetDrawColor(255, 255, 255, 255)
	for _, region := range eventRegions {
		renderer.DrawRect(region.SDLRect())
	}
	renderer.SetDrawColor(0, 0, 0, 255)

	renderer.Present()

	var event sdl.Event
	running := true
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			/*case *sdl.MouseMotionEvent:
			fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
				t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)*/
			case *sdl.MouseButtonEvent:
				fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
				if t.Button == 1 {
					point := &Point{t.X, t.Y}
					what := eventRegions.HitWhat(point)
					if t.Type == 1025 {
						log.Println("MOUSE DOWN. Grab something.")
						if what != nil {
							what.Trigger(EVENT_GRAB)
							log.Println("GRABBING", what)
							interactionState = INTERACTION_STATE_DRAGGING
							draggingWhat = what
						}
					} else if t.Type == 1026 {
						log.Println("MOUSE UP. Drop it.")
						if draggingWhat != nil {
							if what != nil {
								what.Trigger(EVENT_DROP)
							}
							log.Println("DROPPING", draggingWhat, "ON", what)
							interactionState = INTERACTION_STATE_DEFAULT
							draggingWhat = nil
						}

					}
				}
			case *sdl.MouseWheelEvent:
				fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y)
			case *sdl.KeyUpEvent:
				fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
					t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
			}
		}
	}

	window.Destroy()
}
