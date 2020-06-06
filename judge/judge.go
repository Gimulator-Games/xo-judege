package main

import (
	"encoding/json"
	"os"
	"time"

	"gitlab.com/Syfract/Xerac/gimulator-archive/gimulator-test/world"

	"github.com/Gimulator/client-go"
)

var ch = make(chan client.Object, 100)

func main() {
	cl, err := client.NewClient(ch)
	if err != nil {
		panic(err)
	}
	keyRegister := client.Key{
		Type: "register",
	}
	var players []string
	agents, _ := cl.Find(keyRegister)
	for i := 0; i < len(agents); i++ {
		players = append(players, agents[i].Key.Namespace)
	}
	cl.Watch(keyRegister)
	for i := 0; i < 2-len(agents); i++ {
		agentObj := <-ch
		players = append(players, agentObj.Key.Namespace)
	}
	w := world.NewWolrd(players[0], players[1])
	keyAction := client.Key{
		Type: "action",
	}
	var action world.Action
	keyWorld := client.Key{
		Type: "world",
	}
	keyEndGame := client.Key{
		Type: "end_of_game",
	}
	time.Sleep(3 * time.Second)
	jsonW, _ := json.Marshal(w)
	cl.Set(keyWorld, string(jsonW))
	cl.Watch(keyAction)
	for {
		agentObj := <-ch
		data := agentObj.Value
		json.Unmarshal([]byte(data.(string)), &action)
		Update(&w, action, agentObj.Key.Name)
		if len(w.Moves) > 4 {
			switch checkEndGame(&w, agentObj.Key.Name) {
			case "end":
				winner := world.Result{
					Winner: agentObj.Key.Name,
				}
				jsonW, _ := json.Marshal(winner)
				cl.Set(keyEndGame, string(jsonW))
				// log.Println(agentObj.Key.Name)
				os.Exit(0)
			case "draw":
				winner := world.Result{
					Winner: "draw",
				}
				jsonW, _ := json.Marshal(winner)
				cl.Set(keyEndGame, string(jsonW))
				// log.Println(agentObj.Key.Name)
				os.Exit(0)
			}
		}
		jsonW, _ := json.Marshal(w)
		cl.Set(keyWorld, string(jsonW))
	}
}

func Update(w *world.World, action world.Action, playerName string) {
	if playerName == w.Player1.Name {
		w.Moves = append(w.Moves, world.Move{
			Pos:  action.Pos,
			Mark: w.Player1.Mark,
		})
		w.Turn = w.Player2.Name
	} else {
		w.Moves = append(w.Moves, world.Move{
			Pos:  action.Pos,
			Mark: w.Player2.Mark,
		})
		w.Turn = w.Player1.Name
	}
}

func checkEndGame(w *world.World, playerName string) string {
	var m = make([]string, 9, 9)
	for _, move := range w.Moves {
		m[move.Pos] = move.Mark
	}
	if (m[0] != "" && m[0] == m[1] && m[1] == m[2]) ||
		(m[0] != "" && m[0] == m[4] && m[4] == m[8]) ||
		(m[0] != "" && m[0] == m[3] && m[3] == m[6]) ||
		(m[1] != "" && m[1] == m[4] && m[4] == m[7]) ||
		(m[2] != "" && m[2] == m[4] && m[4] == m[6]) ||
		(m[2] != "" && m[2] == m[5] && m[5] == m[8]) ||
		(m[3] != "" && m[3] == m[4] && m[4] == m[5]) ||
		(m[6] != "" && m[6] == m[7] && m[7] == m[8]) {
		w.Result = "end"
		return "end"
	}
	for i := 0; i < 9; i++ {
		if m[i] == "" {
			return ""
		}
	}
	w.Result = "end"
	return "draw"
}
