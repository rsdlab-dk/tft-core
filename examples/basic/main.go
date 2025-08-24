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

	fmt.Println("=== TFT Core - Exemplo Básico ===")

	summonerName := "Hide on bush"
	region := "kr"

	fmt.Printf("Buscando summoner: %s na região %s\n", summonerName, region)

	summoner, err := client.GetSummonerByName(ctx, region, summonerName)
	if err != nil {
		if riotErr, ok := err.(*riot.RiotError); ok {
			switch {
			case riotErr.IsNotFound():
				fmt.Printf("❌ Summoner '%s' não encontrado\n", summonerName)
			case riotErr.IsRateLimited():
				fmt.Println("⏳ Rate limit excedido, tente novamente mais tarde")
			case riotErr.IsUnauthorized():
				fmt.Println("🔑 API key inválida ou expirada")
			default:
				fmt.Printf("❌ Erro da API Riot: %v\n", err)
			}
		} else {
			fmt.Printf("❌ Erro inesperado: %v\n", err)
		}
		return
	}

	fmt.Println("✅ Summoner encontrado!")
	fmt.Printf("📝 Nome: %s\n", summoner.Name)
	fmt.Printf("🆔 PUUID: %s\n", summoner.PUUID)
	fmt.Printf("⭐ Level: %d\n", summoner.SummonerLevel)
	fmt.Printf("🖼️  Icon ID: %d\n", summoner.ProfileIconID)

	fmt.Println("\n=== Buscando pelo PUUID ===")
	
	summoner2, err := client.GetSummonerByPUUID(ctx, region, summoner.PUUID)
	if err != nil {
		fmt.Printf("❌ Erro ao buscar por PUUID: %v\n", err)
		return
	}

	fmt.Printf("✅ Confirmado! Nome: %s, Level: %d\n", summoner2.Name, summoner2.SummonerLevel)
}