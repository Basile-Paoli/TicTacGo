package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Board [BoardSize][BoardSize]Player
type Player int

const BoardSize = 4
const Player1 Player = 0
const Player2 Player = 1
const Empty Player = -1

type Gamestate struct {
	B      Board  `json:"b"`
	ToPlay Player `json:"to_play"`
}

// GetLine renvoie une ligne du plateau entré en paramètre sous forme de string.
func GetLine(n int, b Board) string {
	s := ""
	for i := 0; i < BoardSize; i++ {
		switch b[n][i] {
		case Player1:
			s += "o "
			break
		case Player2:
			s += "x "
			break
		case Empty:
			s += "- "
			break
		}

	}
	return s

}

// PrintBoard affiche l'état du plateau entré en paramètre
func PrintBoard(b Board) error {
	fmt.Printf("%s\n", strings.Repeat("*", 20))
	for i := 0; i < BoardSize; i++ {
		line := GetLine(i, b)
		fmt.Printf("%14s\n", line)
	}
	fmt.Printf("%s\n", strings.Repeat("*", 20))
	return nil
}

// TakeTurn demande au joueur son coup et l'applique au plateau. Retourne -1 si l'action n'est pas possible, 0 si le joueur souhaite sauvegarder, et 1 s'il a joué un coup légal
func TakeTurn(p Player, b *Board) int {
	PrintBoard(*b)
	fmt.Printf("Joueur %d : Quelle case voulez-vous jouer ? (0 pour sauvegarder et quitter)\n", p+1)
	var tile int
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return -1
	}
	tile, err = strconv.Atoi(input)
	if err != nil {
		return -1
	}
	if tile == 0 {
		gs := Gamestate{
			B:      *b,
			ToPlay: p,
		}
		saveGame(&gs)
		println("Partie sauvegardée")
		return 0
	}
	if tile > BoardSize*BoardSize || tile < 1 {
		return -1
	}
	x := (tile - 1) / BoardSize
	y := (tile - 1) % BoardSize
	if (*b)[x][y] != Empty {
		return -1
	}
	(*b)[x][y] = p
	return 1
}

// isOver retourne 2 si la partie n'est pas fini, 0 ou 1 selon le joueur qui a gagné
func isOver(b Board) int {
	if hasWon(1, b) {
		return 1
	}
	if hasWon(0, b) {
		return 0
	}
	return 2
}

// hasWon retourne true si le joueur p a gagné
func hasWon(p Player, b Board) bool {
	diagonal1Win := 1
	diagonal2Win := 1
	for i := 0; i < BoardSize; i++ {
		horizontalWin := 1
		verticalWin := 1
		for j := 0; j < BoardSize; j++ {
			if b[i][j] != p {
				horizontalWin = 0
			}
			if b[j][i] != p {
				verticalWin = 0
			}
		}
		if horizontalWin == 1 || verticalWin == 1 {
			return true
		}
		if b[i][i] != p {
			diagonal1Win = 0
		}
		if b[i][BoardSize-i-1] != p {
			diagonal2Win = 0
		}
	}
	if diagonal1Win == 1 || diagonal2Win == 1 {
		return true
	}
	return false
}

// isNull renvoie true si la partie est nulle (i.e le plateau est rempli)
func isNull(b Board) bool {
	for _, l := range b {
		for _, s := range l {
			if s == Empty {
				return false
			}
		}

	}
	return true
}


func saveGame(gs *Gamestate) {
	jsonGamestate, _ := json.MarshalIndent(*gs, "", "  ")
	os.WriteFile("save.json", jsonGamestate, 0644)
}
func loadGame(gs *Gamestate) error {
	jsonGamestate, err := os.ReadFile("save.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonGamestate, &gs)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	gs := Gamestate{
		B:      Board{},
		ToPlay: 0,
	}
	err := loadGame(&gs)
	os.WriteFile("save.json", []byte{}, 0644)
	if err != nil {
		for i := 0; i < BoardSize; i++ {
			for j := 0; j < BoardSize; j++ {
				gs.B[i][j] = Empty
			}
		}

		gs.ToPlay = Player1
	}

	for {
		n := TakeTurn(gs.ToPlay, &gs.B)
		if n == -1 {
			continue
		}
		if n == 0 {
			break
		}
		if isOver(gs.B) != 2 {
			PrintBoard(gs.B)
			fmt.Printf("Le joueur %d a gagné\n", isOver(gs.B)+1)
			break
		}
		if isNull(gs.B) {
			PrintBoard(gs.B)
			println("Partie nulle")
			break
		}
		gs.ToPlay = 1 - gs.ToPlay
	}
	fmt.Scanln()

}
