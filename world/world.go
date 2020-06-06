package world

import (
	"math/rand"
	"time"
)

type World struct {
	Turn    string `json:"turn"`
	Result  string `json:"result"`
	Player1 Player `json:"player1"`
	Player2 Player `json:"player2"`
	Moves   []Move `json:"moves"`
}

type Player struct {
	Name string `json:"name"`
	Mark string `json:"mark"`
}

type Action struct {
	Pos int `json:"pos"`
}

type Move struct {
	Pos  int    `json:"pos"`
	Mark string `json:"mark"`
}

type Result struct {
	Winner string `json:"winner"`
}

func NewWolrd(playerName1, playerName2 string) World {
	rand.Seed(time.Now().UnixNano())
	rnd := rand.Intn(2)
	var moves = make([]Move, 0, 9)
	player1 := NewPlayer(playerName1, "X")
	player2 := NewPlayer(playerName2, "O")
	if rnd == 1 {
		player1 = NewPlayer(playerName2, "X")
		player2 = NewPlayer(playerName1, "O")
	}
	world := World{
		Turn:    player1.Name,
		Player1: player1,
		Player2: player2,
		Moves:   moves,
	}
	return world
}

func NewPlayer(name string, mark string) Player {
	return Player{
		Name: name,
		Mark: mark,
	}
}

func (w *World) OtherPlayer(playerName string) string {
	if playerName == w.Player1.Name {
		return w.Player2.Name
	}
	return w.Player1.Name
}
