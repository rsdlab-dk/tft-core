package riot

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) GetSummonerByPUUID(ctx context.Context, region, puuid string) (*Summoner, error) {
	endpoint := fmt.Sprintf("%s/tft/summoner/v1/summoners/by-puuid/%s",
		c.getRegionURL("summoner", region), puuid)

	body, err := c.makeRequest(ctx, "GET", endpoint)
	if err != nil {
		return nil, fmt.Errorf("get summoner by puuid %s: %w", puuid, err)
	}

	var summoner Summoner
	if err := json.Unmarshal(body, &summoner); err != nil {
		return nil, fmt.Errorf("unmarshaling summoner: %w", err)
	}

	return &summoner, nil
}

func (c *Client) GetSummonerByID(ctx context.Context, region, summonerID string) (*Summoner, error) {
	endpoint := fmt.Sprintf("%s/tft/summoner/v1/summoners/%s",
		c.getRegionURL("summoner", region), summonerID)

	body, err := c.makeRequest(ctx, "GET", endpoint)
	if err != nil {
		return nil, fmt.Errorf("get summoner by id %s: %w", summonerID, err)
	}

	var summoner Summoner
	if err := json.Unmarshal(body, &summoner); err != nil {
		return nil, fmt.Errorf("unmarshaling summoner: %w", err)
	}

	return &summoner, nil
}

func (c *Client) GetSummonerByRiotID(ctx context.Context, region, gameName, tagLine string) (*Summoner, error) {
	cluster := RegionToCluster(region)

	account, err := c.GetAccountByRiotID(ctx, cluster, gameName, tagLine)
	if err != nil {
		return nil, fmt.Errorf("get account for riot id %s#%s: %w", gameName, tagLine, err)
	}

	summoner, err := c.GetSummonerByPUUID(ctx, region, account.PUUID)
	if err != nil {
		return nil, fmt.Errorf("get summoner by puuid %s: %w", account.PUUID, err)
	}

	return summoner, nil
}
