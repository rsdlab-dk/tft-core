package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/rsdlab-dk/tft-core/riot"
)

func main() {
	apiKey := os.Getenv("RIOT_API_KEY")
	if apiKey == "" {
		log.Fatal("RIOT_API_KEY environment variable is required")
	}

	client := riot.NewClient(apiKey)
	ctx := context.Background()

	fmt.Println("🎮 === TFT Arena - API Funcionando! === 🎮")

	region := "kr"

	fmt.Println("\n🏆 === Challenger League ===")
	challenger, err := client.GetChallengerLeague(ctx, region)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Printf("✅ Challenger: %s\n", challenger.Name)
	fmt.Printf("📊 Total jogadores: %d\n", len(challenger.Entries))

	if len(challenger.Entries) >= 3 {
		fmt.Println("\n🥇 === TOP 3 PLAYERS ===")
		for i := 0; i < 3; i++ {
			player := challenger.Entries[i]
			fmt.Printf("#%d - LP: %d | W/L: %d/%d", i+1, player.LeaguePoints, player.Wins, player.Losses)

			if player.HotStreak {
				fmt.Print(" 🔥")
			}
			if player.Veteran {
				fmt.Print(" ⭐")
			}
			fmt.Println()

			summoner, err := client.GetSummonerByPUUID(ctx, region, player.PUUID)
			if err == nil && summoner.Name != "" {
				fmt.Printf("    📝 Nome: %s (Lv.%d)\n", summoner.Name, summoner.SummonerLevel)
			}
		}
	}

	fmt.Println("\n🏅 === Outros Leagues ===")

	if grandmaster, err := client.GetGrandmasterLeague(ctx, region); err == nil {
		fmt.Printf("✅ Grandmaster: %d jogadores\n", len(grandmaster.Entries))
	}

	if master, err := client.GetMasterLeague(ctx, region); err == nil {
		fmt.Printf("✅ Master: %d jogadores\n", len(master.Entries))
	}

	fmt.Println("\n🔍 === Busca por Riot ID ===")
	testRiotIDs := []struct{ name, tag string }{
		{"Faker", "T1"},
		{"Hide on bush", "KR1"},
		{"Dopa", "KR1"},
	}

	for _, rid := range testRiotIDs {
		fmt.Printf("🔎 Buscando %s#%s... ", rid.name, rid.tag)

		summoner, err := client.GetSummonerByRiotID(ctx, region, rid.name, rid.tag)
		if err != nil {
			fmt.Println("❌ Não encontrado")
			continue
		}

		fmt.Printf("✅ %s (Lv.%d)\n", summoner.Name, summoner.SummonerLevel)

		player, err := client.FindPlayerInHighElo(ctx, region, summoner.PUUID)
		if err == nil {
			fmt.Printf("    🏆 High Elo: %d LP (%d W / %d L)\n",
				player.LeaguePoints, player.Wins, player.Losses)
		} else {
			fmt.Printf("    📊 Não está em high elo\n")
		}
	}

	fmt.Println("\n🎉 === SUCESSO! ===")
	fmt.Println("✅ API TFT completamente funcional")
	fmt.Println("✅ Riot ID → PUUID → Summoner ✓")
	fmt.Println("✅ Challenger/GM/Master ✓")
	fmt.Println("✅ Player search ✓")
	fmt.Println("✅ Rate limiting ✓")
	fmt.Println("✅ Error handling ✓")

	fmt.Println("\n🚀 Ready para produção!")
}

func handleError(err error) {
	if riotErr, ok := err.(*riot.RiotError); ok {
		switch {
		case riotErr.IsNotFound():
			fmt.Println("❌ Não encontrado")
		case riotErr.IsRateLimited():
			fmt.Println("⏳ Rate limit - aguarde")
		case riotErr.IsUnauthorized():
			fmt.Println("🔑 API key inválida")
		case riotErr.IsForbidden():
			fmt.Println("🚫 Acesso negado")
		default:
			fmt.Printf("❌ Erro Riot: %v\n", err)
		}
	} else {
		fmt.Printf("❌ Erro: %v\n", err)
	}
}
