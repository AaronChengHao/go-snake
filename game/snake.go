package game

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"math/rand"
	"snake/util"
	"time"
)

var (
	x = 0
	y = 0
)

var (
	up      int8  = 1
	down    int8  = 2
	left    int8  = 3
	right   int8  = 4
	stepNum int64 = 20
	nodeW   int64 = 20
	nodeH   int64 = 20
)

var direction int8 = right

const (
	WindowW    = 640
	WindowH    = 640
	sampleRate = 48000
)

// 食物结构体
type food struct {
	x     int64
	y     int64
	color color.Color
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		direction = down
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		direction = left
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		direction = up
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		direction = right
	}

	//if inpututil.IsKeyJustPressed(ebiten.KeyP) {
	//	// As audioPlayer has one stream and remembers the playing position,
	//	// rewinding is needed before playing when reusing audioPlayer.
	//	err := g.audioPlayer.Rewind()
	//	if err != nil {
	//		return err
	//	}
	//	g.audioPlayer.Play()
	//}
	return nil
}

var (
	ebitenImage *ebiten.Image
	colors      = []color.RGBA{
		{0xff, 0xff, 0xff, 0xff},
		{0xff, 0xff, 0x0, 0xff},
	}
)

func drawText(screen *ebiten.Image) {

	if x > WindowW || y > WindowH {
		x = 0
		y = 0
	}
	ebitenutil.DebugPrintAt(screen, "程睿曦 程睿懿", x, y)

}

func (g *Game) DrawFood(screen *ebiten.Image) {
	if g.food == nil {
		randW := rand.Intn(WindowW)
		randH := rand.Intn(WindowH)
		spanW := int64(randW) / nodeW
		spanH := int64(randH) / nodeH
		g.food = &food{
			x: spanW * nodeW,
			y: spanH * nodeH,
		}
	}

	foodImg := ebiten.NewImage(20, 20)
	foodImg.Fill(color.White)
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(g.food.x), float64(g.food.y))
	screen.DrawImage(foodImg, op)
}

func (g *Game) isEatFood(screen *ebiten.Image) bool {
	if (g.head.X >= g.food.x && g.head.X <= g.food.x+nodeW) && (g.head.Y >= g.food.y && g.head.Y <= g.food.y+nodeH) {
		util.AddNode(g.head)
		g.food = nil
		g.DrawFood(screen)

		_ = g.audioPlayer.Rewind()
		g.audioPlayer.Play()
	}
	return true
}

func printNodeCount(node *util.Node) int {
	if node.Child != nil {
		return printNodeCount(node.Child) + 1
	}
	return 1
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 设置背景色
	screen.Fill(color.NRGBA{G: 0x40, B: 0x80, A: 0xff})
	time.Sleep(time.Millisecond * 100)
	// 画食物
	g.DrawFood(screen)
	// 判断蛇头师傅吃到食物了
	g.isEatFood(screen)

	// 画文本
	drawText(screen)

	node := g.head

	node.OldX = node.X
	node.OldY = node.Y
	switch direction {
	case up:
		node.Y -= stepNum
	case down:
		node.Y += stepNum
	case left:
		node.X -= stepNum
	case right:
		node.X += stepNum
	}

	if node.X > WindowW ||
		node.X < 0 ||
		node.Y < 0 ||
		node.Y > WindowH {
		node.X = 0
		node.Y = 0
		direction = right
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(node.X), float64(node.Y))

	nodeImg := ebiten.NewImage(int(nodeW), int(nodeH))
	nodeImg.Fill(node.Color)
	fmt.Println(node.Color)
	screen.DrawImage(nodeImg, op)
	drawSnake(g.head.Child, screen)
}

func drawSnake(node *util.Node, screen *ebiten.Image) bool {
	if node != nil {
		node.OldX = node.X
		node.OldY = node.Y

		node.X = node.Parent.OldX
		node.Y = node.Parent.OldY

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(node.X), float64(node.Y))

		nodeImg := ebiten.NewImage(int(nodeW), int(nodeH))
		nodeImg.Fill(node.Color)
		// 画
		screen.DrawImage(nodeImg, op)

		return drawSnake(node.Child, screen)
	}

	return true
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowW, WindowH
}

type Game struct {
	head         *util.Node
	keys         []ebiten.Key
	food         *food
	audioContext *audio.Context
	audioPlayer  *audio.Player
}

func NewGame() (*Game, error) {
	var err error

	g := &Game{}
	g.head = &util.Node{X: 0, Y: 0, OldX: 0, OldY: 0, Color: color.White}

	g.audioContext = audio.NewContext(sampleRate)
	//d, err := wav.Decode(g.audioContext, bytes.NewReader(raudio.Jab_wav))
	file, err := ioutil.ReadFile("assets/eat.wav")
	if err != nil {
		return nil, err
	}
	d, err := wav.Decode(g.audioContext, bytes.NewReader(file))
	if err != nil {
		return nil, err
	}

	// Create an audio.Player that has one stream.
	g.audioPlayer, err = g.audioContext.NewPlayer(d)
	if err != nil {
		return nil, err
	}

	return g, nil
}
