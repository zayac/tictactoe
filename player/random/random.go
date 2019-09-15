// Package random implements a random player of tic-tac-toe.
package random

import (
	"context"
	"fmt"
	"html/template"
	"math/rand"
	"time"

	"github.com/zayac/tictactoe"
)

// NewPlayer creates a random player of tic-tac-toe.
func NewPlayer() (tictactoe.Player, error) {
	gophers := []template.URL{
		"https://raw.githubusercontent.com/zayac/tictactoe/master/player/random/gopher-0.png",
		"https://raw.githubusercontent.com/zayac/tictactoe/master/player/random/gopher-1.png",
		"https://raw.githubusercontent.com/zayac/tictactoe/master/player/random/gopher-2.png",
	}
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return player{
		rand:  rand,
		image: gophers[rand.Intn(len(gophers))],
	}, nil
}

// player is a random player of tic-tac-toe.
type player struct {
	rand  *rand.Rand
	image template.URL
}

func (p player) Name() string {
	return "Random Player"
}

func (p player) Image() template.URL {
	return p.image
}

// Play takes a tic-tac-toe board b and returns the next move
// for this player. Its mark is either X or O.
// ctx is expected to have a deadline set, and Play may take time
// to "think" until deadline is reached before returning.
func (p player) Play(ctx context.Context, b tictactoe.Board, mark tictactoe.State) (tictactoe.Move, error) {
	if b.Condition() != tictactoe.NotEnd {
		return tictactoe.Move(-1), fmt.Errorf("board has a finished game")
	}

	stopThinking := time.Now().Add(2 * time.Second)
	if t, ok := ctx.Deadline(); ok {
		t = t.Add(-time.Second)
		if t.Before(stopThinking) {
			stopThinking = t
		}
	}

	// Decide on a move to make... using randomness!
	legalMoves := legalMoves(b)
	move := legalMoves[p.rand.Intn(len(legalMoves))]

	// Take some more time to pretend we're still "thinking".
	time.Sleep(time.Until(stopThinking))

	return move, nil
}

// legalMoves returns all legal moves on board b.
func legalMoves(b tictactoe.Board) []tictactoe.Move {
	var moves []tictactoe.Move
	for i, cell := range b.Cells {
		if cell != tictactoe.F {
			continue
		}
		moves = append(moves, tictactoe.Move(i))
	}
	return moves
}
