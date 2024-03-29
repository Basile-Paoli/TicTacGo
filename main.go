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

const BoardSize = 3
const Player1 Player = 0
const Player2 Player = 1
const NullPlayer Player = -1

type Gamestate struct {
	Board  Board  `json:"b"`
	ToPlay Player `json:"to_play"`
	Score  [2]int `json:"score"`
}

const Reset = "\033[0m"
const Red = "\033[31m"
const Blue = "\033[34m"

// resetBoard initialise l'état de la partie
func resetBoard(gs *Gamestate) {
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			gs.Board[i][j] = NullPlayer
		}
	}

	gs.ToPlay = Player1
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
	fmt.Printf("%s\n", strings.Repeat("*", 21))
	for i := 0; i < BoardSize; i++ {
		line := GetLine(i, b)
		print(strings.Repeat(" ", 8))
		for _, c := range line {
			switch c {
			case 'o':
				fmt.Print(Red + "o" + Reset)
				break
			case 'x':
				fmt.Print(Blue + "x" + Reset)

			default:
				fmt.Printf("%c", c)
			}
		}
		println()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 20))
	return nil
}

// TakeTurn demande au joueur son coup et l'applique au plateau. Retourne -1 si l'action n'est pas possible, 0 si le joueur souhaite sauvegarder, et 1 s'il a joué un coup légal
func TakeTurn(gs *Gamestate) int {
	PrintBoard(gs.Board)
	fmt.Printf("Joueur %d : Quelle case voulez-vous jouer ? (1-9 | 0 pour sauvegarder et quitter)\n", gs.ToPlay+1)
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
		saveGame(gs)
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
	if (gs.Board)[x][y] != NullPlayer {
		println("Case non vide")
		return -1
	}
	(gs.Board)[x][y] = gs.ToPlay
	return 1
}

// isOver retourne -1 si la partie n'est pas fini, 0 ou 1 selon le joueur qui a gagné
func isOver(b Board) Player {
	if hasWon(Player1, b) {
		return Player1
	}
	if hasWon(Player2, b) {
		return Player2
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
		n := TakeTurn(gs)
		if n == -1 {
			continue
		}
		if isOver(gs.Board) != NullPlayer {
			PrintBoard(gs.Board)
			winningPlayer := isOver(gs.Board)
			fmt.Printf("Le joueur %d a gagné, ", winningPlayer+1)
			return winningPlayer
		}
		if isNull(gs.Board) {
			PrintBoard(gs.Board)
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
		resetBoard(&gs)
	}
	os.WriteFile("save.json", []byte{}, 0644)
	for {
		switch playGame(&gs) {
		case Player1:
			gs.Score[0] += 1
		case Player2:
			gs.Score[1] += 1
		}
		println("score : ", gs.Score[0], " - ", gs.Score[1])
	out:
		for {
			println("1 - Rejouer\n" +
				"2 - Sauvegarder\n" +
				"3 - Quitter")
			var rep string
			fmt.Scanln(&rep)
			switch rep {

			case "1":
				resetBoard(&gs)
				break out
			case "2":
				resetBoard(&gs)
				saveGame(&gs)
				println("Partie sauvegardée")
				fmt.Scanln()
				os.Exit(0)
			case "3":
				os.Exit(0)
			}

		}
	}

}
