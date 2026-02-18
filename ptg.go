package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type Game struct {
	Name            string `json:"name"`
	PlaytimeForever int    `json:"playtime_forever"`
}

type Response struct {
	Games []Game `json:"games"`
}

type SteamResponse struct {
	Response Response `json:"response"`
}

func main() {
	fmt.Println(`
    ╔════════════════════════════════════════════════╗
    ║  ✦ ✧ ✦ ✧ ✦ ✧ ✦ ✧ ✦ ✧ ✦ ✧ ✦ ✧ ✦ ✧ ✦ ✧ ✦ ✧ ✦ ✧ ✦ ║
    ╠════════════════════════════════════════════════╣
    ║                                                ║
    ║     ██████╗ ██╗      █████╗ ██╗   ██╗          ║
    ║     ██╔══██╗██║     ██╔══██╗╚██╗ ██╔╝          ║
    ║     ██████╔╝██║     ███████║ ╚████╔╝           ║
    ║     ██╔═══╝ ██║     ██╔══██║  ╚██╔╝            ║
    ║     ██║     ███████╗██║  ██║   ██║             ║
    ║     ╚═╝     ╚══════╝╚═╝  ╚═╝   ╚═╝             ║
    ║                                                ║
    ║        ████████╗██╗  ██╗ █████╗ ████████╗      ║
    ║           ██╔══╝██║  ██║██╔══██╗╚══██╔══╝      ║
    ║           ██║   ███████║███████║   ██║         ║
    ║           ██║   ██╔══██║██╔══██║   ██║         ║
    ║           ██║   ██║  ██║██║  ██║   ██║         ║
    ║           ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝         ║
    ║                                                ║
    ║     ██████╗  █████╗ ███╗   ███╗███████╗ ██╗    ║
    ║    ██╔════╝ ██╔══██╗████╗ ████║██╔════╝ ██║    ║
    ║    ██║  ███╗███████║██╔████╔██║█████╗   ██║    ║
    ║    ██║   ██║██╔══██║██║╚██╔╝██║██╔══╝   ╚═╝    ║
    ║    ╚██████╔╝██║  ██║██║ ╚═╝ ██║███████╗ ██╗    ║
    ║     ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝ ╚═╝    ║
    ║                                                ║
    ╠════════════════════════════════════════════════╣
    ║  ✦ ✧ ✦ ✧ ✦ ✧ ✦ ✧ ✦ ✧ ✦ ✧ ✦ ✧ ✦ ✧ ✦ ✧ ✦ ✧ ✦ ✧ ✦ ║
    ╚════════════════════════════════════════════════╝`)
	fmt.Println()
	fmt.Println("Get a Steam API key from https://steamcommunity.com/dev/apikey")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Your Steam API key: ")
	scanner.Scan()
	apiKey := scanner.Text()

	fmt.Print("Steam ID: ")
	scanner.Scan()
	steamId := scanner.Text()

	apiUrl := fmt.Sprintf("https://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json&include_appinfo=1", apiKey, steamId)

	resp, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println("Error making request: ", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Uh-oh: Got a %d instead of 200 OK from the Steam API\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading request: ", err)
		return
	}

	steamResp := SteamResponse{}
	err = json.Unmarshal(body, &steamResp)
	if err != nil {
		fmt.Println("Error parsing JSON: ", err)
		return
	}

	if len(steamResp.Response.Games) == 0 {
		fmt.Println("This user is hiding their games! Make sure the profile is public.")
		return
	}

	for {
		fmt.Println(`
 /)/)
( . .)
( づ Rolling...`)
		fmt.Println()

		time.Sleep(3 * time.Second)

		randomIndex := rand.Intn(len(steamResp.Response.Games))
		game := steamResp.Response.Games[randomIndex]

		if game.PlaytimeForever == 0 {
			fmt.Printf("Play %s! (No hours on record)\n", game.Name)
		} else {
			hours := game.PlaytimeForever / 60
			minutes := game.PlaytimeForever % 60
			fmt.Printf("Play %s! (%dh %dm on record)\n", game.Name, hours, minutes)
		}

		fmt.Println()
		fmt.Print("Reroll? (Y/n): ")
		scanner.Scan()
		answer := strings.ToLower(scanner.Text())

		if answer == "n" || answer == "no" {
			return
		}
	}
}
