package main

import (
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	winWidth = 800
	winHeight = 600
)

//Game State
type gameState int
const (
	start gameState = iota
	play
)

type aiState int
const (
	wait aiState = iota
	move
)

var state = start
var aistate = move

type color struct {
	r, g, b byte
}

type position struct {
	x, y float32
}

type ball struct {
	position
	radius int
	xVelocity float32
	yVelocity float32
	color color
}

var numberFonts = [][]byte{
	{1,1,1,
	1,0,1,
	1,0,1,
	1,0,1,
	1,1,1},

	{1,1,0,
	0,1,0,
	0,1,0,
	0,1,0,
	0,1,0,},

	{1,1,1,
	0,0,1,
	1,1,1,
	1,0,0,
	1,1,1},

	{1,1,1,
	0,0,1,
	1,1,1,
	0,0,1,
	1,1,1,},
}

func lerp(a, b, percentage float32) float32 {
	return a+percentage*(b-a)
}

func drawNumber(pos position, color color, size, num int, pixels []byte) {
	startX := int(pos.x) - (size*3)/2
	startY := int(pos.y) - (size*5)/2

	for i, v := range numberFonts[num] {
		if v == 1{
			for y := startY; y < startY+size; y++ {
				for x := startX; x < startX+size; x++ {
					setPixel(x, y, color, pixels)
				}
			}
		}
		startX += size
		if (i+1) % 3 == 0 {
			startY += size
			startX -= size*3
		}
	}
}

func getCenter() position {
	return position{float32(winWidth)/2, float32(winHeight)/2}
}

func (ball *ball) update(leftPaddle, rightPaddle *paddle, elapsedTime float32) {
	ball.x += ball.xVelocity * elapsedTime
	ball.y += ball.yVelocity * elapsedTime
	if int(ball.y)-ball.radius < 0 || int(ball.y)+ball.radius > winHeight {
		ball.yVelocity *= -1
	}
	if ball.x < 0 {
		rightPaddle.score++
		ball.position = getCenter()
		state = start
		aistate = move
		ball.xVelocity = 300
		ball.yVelocity = 300
		leftPaddle.position = position{100,300+rand.Float32()*(500-300)}
		rightPaddle.position = position{700,300+rand.Float32()*(500-300)}
	} else if ball.x > winWidth {
		leftPaddle.score++
		ball.position = getCenter()
		state = start
		aistate = move
		ball.xVelocity = -300
		ball.yVelocity = 300
		leftPaddle.position = position{100,300+rand.Float32()*(500-300)}
		rightPaddle.position = position{700,300+rand.Float32()*(500-300)}
	}

	if ball.x-float32(ball.radius) < leftPaddle.x+leftPaddle.width/2 && ball.x-float32(ball.radius) > leftPaddle.x-leftPaddle.width/2 {
		if ball.y > leftPaddle.y-leftPaddle.height/2 && ball.y < leftPaddle.y+leftPaddle.height/2 {
			ball.xVelocity *= -1
			if (ball.xVelocity < 0) {
				ball.xVelocity -= 25
			} else {
				ball.xVelocity += 25
			}
			ball.yVelocity += 25
			aistate = move
		}
	}

	if ball.x + float32(ball.radius) > rightPaddle.x-rightPaddle.width/2 && ball.x + float32(ball.radius) < rightPaddle.x+rightPaddle.width/2 {
		if ball.y > rightPaddle.y-rightPaddle.height/2 && ball.y < rightPaddle.y+rightPaddle.height/2 {
			ball.xVelocity *= -1
			if (ball.xVelocity < 0) {
				ball.xVelocity -= 25
			} else {
				ball.xVelocity += 25
			}
			ball.yVelocity += 25
			aistate = wait
		}
	}
}

func (ball *ball) draw(pixels []byte) {
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x)+x, int(ball.y)+y, ball.color, pixels)
			}
		}
	}
}

type paddle struct {
	position
	width float32
	height float32
	speed float32
	score int
	color color
}

func (paddle *paddle) update(keyState []uint8, elapsedTime float32) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.y -= paddle.speed * elapsedTime
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y += paddle.speed * elapsedTime
	}
}

func aiMoveUp(paddle *paddle, elapsedTime float32) {
	paddle.y -= paddle.speed * elapsedTime
}

func aiMoveDown(paddle *paddle, elapsedTime float32) {
	paddle.y += paddle.speed * elapsedTime
}

func (paddle *paddle) aiUpdate(ball *ball, elapsedTime float32) {
	if aistate == move {
		paddle.speed = 300
		if ball.y > paddle.y {
			aiMoveDown(paddle, elapsedTime)
		} else {
			aiMoveUp(paddle,elapsedTime)
		}
	} else {
		paddle.speed = 100 + rand.Float32()*(300-100)
		random := rand.Intn(2)
		if random == 0 {
			aiMoveDown(paddle, elapsedTime)
		} else {
			aiMoveUp(paddle, elapsedTime)
		}
	}
}

func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x - (paddle.width/2))
	startY := int(paddle.y - (paddle.height/2))

	for y := 0; y < int(paddle.height); y++ {
		for x := 0; x < int(paddle.width); x++ {
			setPixel(startX+x, startY+y, paddle.color, pixels)
		}
	}

	numX := lerp(paddle.x, getCenter().x, 0.2)
	drawNumber(position{numX,35},paddle.color, 10, paddle.score, pixels)
}

func setPixel(x, y int, c color, pixels []byte) {
	index := ((y*winWidth)+x)*4

	if index < len(pixels) && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	window, err := sdl.CreateWindow("Test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		panic(err)
	}
	defer texture.Destroy()

	pixels := make([]byte, winWidth*winHeight*4)

	player1 := paddle{
		position: position{100,300+rand.Float32()*(500-300)},
		width: 20,
		height: 100,
		color: color{255,255,255},
		score: 0,
		speed: 300,
	}

	player2 := paddle{
		position: position{700,300+rand.Float32()*(500-300)},
		width: 20,
		height: 100,
		color: color{255,255,255},
		score: 0,
		speed: 300,
	}

	ball := ball{
		position: position{300,300},
		radius: 20,
		xVelocity: 300,
		yVelocity: 300,
		color: color{255,255,255},
	}

	keyState := sdl.GetKeyboardState()

	var timeStart time.Time
	var elapsedTime float32

	for {
		timeStart = time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		
		if (state == play) {
			player1.update(keyState, elapsedTime)
			player2.aiUpdate(&ball, elapsedTime)
			ball.update(&player1, &player2, elapsedTime)
		} else if (state == start) {
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				if player1.score == 3 || player2.score == 3 {
					player1.score = 0
					player2.score = 0
				}
				state = play
			}
		}

		clear(pixels)
		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		texture.Update(nil, pixels, winWidth*4)
		renderer.Copy(texture, nil, nil)
		renderer.Present()

		elapsedTime = float32(time.Since(timeStart).Seconds())
		if elapsedTime < .005 {
			sdl.Delay(5-uint32(elapsedTime*1000.0))
			elapsedTime = float32(time.Since(timeStart).Seconds())
		}
	}
}