package main

import (
	"fmt"
	"log"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type EventRegion struct {
	Rect *sdl.Rect
}

func NewEventRegion(x, y, width, height int32) *EventRegion {
	er := EventRegion{Rect: &sdl.Rect{X: x, Y: y, W: width, H: height}}
	return &er
}

func (er *EventRegion) SDLRect() *sdl.Rect {
	return er.Rect
}

func (er *EventRegion) HitWhat(x int32, y int32) *EventRegion {
	if x < er.Rect.X || x > er.Rect.X+er.Rect.W {
		return nil
	}
	if y < er.Rect.Y || y > er.Rect.Y+er.Rect.H {
		return nil
	}
	return er
}

type RegionList []*EventRegion

func (rl RegionList) HitWhat(x int32, y int32) *EventRegion {
	for _, region := range rl {
		if what := region.HitWhat(x, y); what != nil {
			return what
		}
	}
	return nil
}

const CARD_WIDTH = 100
const CARD_HEIGHT = 180

func main() {
	res := sdl.Init(sdl.INIT_VIDEO)
	log.Println(res)

	//eventRegions := make([]*EventRegion, 0)
	eventRegions := make(RegionList, 0)
	var x int32
	var y int32 = 20
	for i := 0; i < 8; i++ {
		x = int32(80 + (i*CARD_WIDTH + i*10))
		eventRegions = append(eventRegions, NewEventRegion(x, y, CARD_WIDTH, CARD_HEIGHT))
	}

	y = 520
	for i := 0; i < 8; i++ {
		x = int32(80 + (i*CARD_WIDTH + i*10))
		eventRegions = append(eventRegions, NewEventRegion(x, y, CARD_WIDTH, CARD_HEIGHT))
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
				fmt.Println(eventRegions.HitWhat(t.X, t.Y))
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
