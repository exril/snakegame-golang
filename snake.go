package main

import (
	"fmt"
	"math/rand"
	"time"
	"github.com/nsf/termbox-go"
)

const (
	width  = 20
	height = 20
)

type point struct {
	x, y int
}

type snake struct {
	body       []point
	direction  int
	nextDir    int
	growing    bool
}

var (
	quit       = make(chan bool)
	score      = 0
	gameOver   = false
	snakeColor = termbox.ColorGreen
	foodColor  = termbox.ColorRed
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)

	snake := newSnake()
	placeFood(snake)

	go func() {
		for {
			ev := termbox.PollEvent()
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyArrowUp:
					snake.setDirection(0)
				case termbox.KeyArrowRight:
					snake.setDirection(1)
				case termbox.KeyArrowDown:
					snake.setDirection(2)
				case termbox.KeyArrowLeft:
					snake.setDirection(3)
				case termbox.KeyEsc:
					quit <- true
					return
				}
			}
		}
	}()

	go func() {
		for {
			snake.move()
			if snake.collided() {
				gameOver = true
				quit <- true
				return
			}
			if snake.ateFood() {
				score++
				snake.grow()
				placeFood(snake)
			}
			time.Sleep(150 * time.Millisecond)
		}
	}()

	for {
		draw(snake)
		if gameOver {
			drawGameOver()
			break
		}
		select {
		case <-quit:
			break
		default:
		}
	}
}

func newSnake() *snake {
	return &snake{
		body:      []point{{width / 2, height / 2}},
		direction: 1,
		nextDir:   1,
		growing:   false,
	}
}

func (s *snake) setDirection(dir int) {
	if (dir == 0 && s.direction != 2) || (dir == 1 && s.direction != 3) || (dir == 2 && s.direction != 0) || (dir == 3 && s.direction != 1) {
		s.nextDir = dir
	}
}

func (s *snake) move() {
	head := s.body[0]
	var newHead point
	switch s.nextDir {
	case 0:
		newHead = point{head.x, head.y - 1}
	case 1:
		newHead = point{head.x + 1, head.y}
	case 2:
		newHead = point{head.x, head.y + 1}
	case 3:
		newHead = point{head.x - 1, head.y}
	}

	s.body = append([]point{newHead}, s.body...)

	if !s.growing {
		s.body = s.body[:len(s.body)-1]
	} else {
		s.growing = false
	}
	s.direction = s.nextDir
}

func (s *snake) grow() {
	s.growing = true
}

func (s *snake) collided() bool {
	head := s.body[0]
	if head.x <= 0 || head.x >= width || head.y <= 0 || head.y >= height {
		return true
	}
	for i := 1; i < len(s.body); i++ {
		if head == s.body[i] {
			return true
		}
	}
	return false
}

func (s *snake) ateFood() bool {
	head := s.body[0]
	food := s.body[len(s.body)-1]
	return head == food
}

func placeFood(s *snake) {
	rand.Seed(time.Now().UnixNano())
	for {
		x := rand.Intn(width-2) + 1
		y := rand.Intn(height-2) + 1
		food := point{x, y}
		collision := false
		for _, b := range s.body {
			if b == food {
				collision = true
				break
			}
		}
		if !collision {
			s.body = append(s.body, food)
			break
		}
	}
}

func draw(s *snake) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for _, p := range s.body {
		termbox.SetCell(p.x, p.y, ' ', snakeColor, snakeColor)
	}
	food := s.body[len(s.body)-1]
	termbox.SetCell(food.x, food.y, ' ', foodColor, foodColor)
	termbox.Flush()
}

func drawGameOver() {
	msg := "Game Over! Your score was: " + fmt.Sprintf("%d", score)
	x := (width - len(msg)) / 2
	y := height / 2
	for _, c := range msg {
		termbox.SetCell(x, y, c, termbox.ColorDefault, termbox.ColorDefault)
		x++
	}
	termbox.Flush()
}
