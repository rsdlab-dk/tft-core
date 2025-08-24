package riot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

func (c *Client) GetSummonerByName(ctx context.Context, region, summonerName string) (*Summoner, error) {
	encodedName := url.QueryEscape(summonerName)
	endpoint := fmt.Sprintf("%s/tft/summoner/v1/summoners/by-name/%s", 
		c.getRegionURL("summoner", region), encodedName)
	
	body, err := c.makeRequest(ctx, "GET", endpoint)
	if err != nil {
		return nil, fmt.Errorf("get summoner by name %s: %w", summonerName, err)
	}

	var summoner Summoner
	if err := json.Unmarshal(body, &summoner); err != nil {
		return nil, fmt.Errorf("unmarshaling summoner: %w", err)
	}

	return &summoner, nil
}

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