package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Board [BoardSize][BoardSize]Player

func emptyGamestate(gs *Gamestate) {
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			gs.B[i][j] = NullPlayer
		}
	}

	gs.ToPlay = Player1
}

type Player int

const BoardSize = 3
const Player1 Player = 0
const Player2 Player = 1
const NullPlayer Player = -1

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
		case NullPlayer:
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
	fmt.Printf("Joueur %d : Quelle case voulez-vous jouer ? (1-9 | 0 pour sauvegarder et quitter)\n", p+1)
	var tile int
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return -1
	}
	tile, err = strconv.Atoi(input)
	if err != nil {
		println("Saisie invalide")
		return -1
	}
	if tile == 0 {
		gs := Gamestate{
			B:      *b,
			ToPlay: p,
		}
		saveGame(&gs)
		println("Partie sauvegardée")
		fmt.Scanln()
		os.Exit(0)
	}
	if tile > BoardSize*BoardSize || tile < 1 {
		println("Saisie invalide")
		return -1
	}
	x := (tile - 1) / BoardSize
	y := (tile - 1) % BoardSize
	if (*b)[x][y] != NullPlayer {
		println("Case non vide")
		return -1
	}
	(*b)[x][y] = p
	return 1
}

// isOver retourne -1 si la partie n'est pas fini, 0 ou 1 selon le joueur qui a gagné
func isOver(b Board) Player {
	if hasWon(Player1, b) {
		return Player1
	}
	if hasWon(Player1, b) {
		return Player1
	}
	return NullPlayer
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
			if s == NullPlayer {
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

// playGame fait jouer une partie et renvoie le numéro du joueurs gagnant, ou -1 en cas de nulle
func playGame(gs *Gamestate) Player {
	for {
		n := TakeTurn(gs.ToPlay, &gs.B)
		if n == -1 {
			continue
		}
		if isOver(gs.B) != NullPlayer {
			PrintBoard(gs.B)
			winningPlayer := isOver(gs.B)
			fmt.Printf("Le joueur %d a gagné, ", winningPlayer+1)
			return winningPlayer
		}
		if isNull(gs.B) {
			PrintBoard(gs.B)
			print("Partie nulle, ")
			return NullPlayer
		}
		gs.ToPlay = 1 - gs.ToPlay
	}

}

func main() {
	gs := Gamestate{}
	err := loadGame(&gs)
	if err != nil {
		emptyGamestate(&gs)
	}
	os.WriteFile("save.json", []byte{}, 0644)
	score := []int{0, 0}
	for {
		switch playGame(&gs) {
		case Player1:
			score[0] += 1
		case Player2:
			score[1] += 1
		}
		println("score : ", score[0], " - ", score[1])
	out:
		for {
			println("1 - Rejouer\n2 - Quitter")
			var rep string
			fmt.Scanln(&rep)
			switch rep {
			case "1":
				emptyGamestate(&gs)
				break out
			case "2":
				os.Exit(0)
			}

		}
	}

}
