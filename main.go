package main

import (
	"github.com/s992/nfl/ui"
)

func main() {
	ui.Run()
	// games, err := api.ListGames()
	// if err != nil {
	// 	panic(err)
	// }

	// for _, game := range games {
	// 	fmt.Printf("%v: %s\n", game.StartTime.Local().Format("02 Jan 03:04pm"), game.ShortName)
	// }

	// plays, err := api.GetPlays(games[0].Id)
	// if err != nil {
	// 	panic(err)
	// }

	// s, err := api.GetScoreboards(time.Now())
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(s)
}
