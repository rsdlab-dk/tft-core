package riot

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) GetChallengerLeague(ctx context.Context, region string) (*LeagueList, error) {
	endpoint := fmt.Sprintf("%s/tft/league/v1/challenger",
		c.getRegionURL("league", region))

	body, err := c.makeRequest(ctx, "GET", endpoint)
	if err != nil {
		return nil, fmt.Errorf("get challenger league: %w", err)
	}

	var league LeagueList
	if err := json.Unmarshal(body, &league); err != nil {
		return nil, fmt.Errorf("unmarshaling challenger league: %w", err)
	}

	return &league, nil
}

func (c *Client) GetGrandmasterLeague(ctx context.Context, region string) (*LeagueList, error) {
	endpoint := fmt.Sprintf("%s/tft/league/v1/grandmaster",
		c.getRegionURL("league", region))

	body, err := c.makeRequest(ctx, "GET", endpoint)
	if err != nil {
		return nil, fmt.Errorf("get grandmaster league: %w", err)
	}

	var league LeagueList
	if err := json.Unmarshal(body, &league); err != nil {
		return nil, fmt.Errorf("unmarshaling grandmaster league: %w", err)
	}

	return &league, nil
}

func (c *Client) GetMasterLeague(ctx context.Context, region string) (*LeagueList, error) {
	endpoint := fmt.Sprintf("%s/tft/league/v1/master",
		c.getRegionURL("league", region))

	body, err := c.makeRequest(ctx, "GET", endpoint)
	if err != nil {
		return nil, fmt.Errorf("get master league: %w", err)
	}

	var league LeagueList
	if err := json.Unmarshal(body, &league); err != nil {
		return nil, fmt.Errorf("unmarshaling master league: %w", err)
	}

	return &league, nil
}

func (c *Client) GetLeagueEntriesByTier(ctx context.Context, region, tier, division string) ([]LeagueEntry, error) {
	endpoint := fmt.Sprintf("%s/tft/league/v1/entries/%s/%s",
		c.getRegionURL("league", region), tier, division)

	body, err := c.makeRequest(ctx, "GET", endpoint)
	if err != nil {
		return nil, fmt.Errorf("get league entries by tier %s/%s: %w", tier, division, err)
	}

	var entries []LeagueEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, fmt.Errorf("unmarshaling league entries: %w", err)
	}

	return entries, nil
}

func (c *Client) FindPlayerInHighElo(ctx context.Context, region, puuid string) (*LeagueItem, error) {
	leagues := []func(context.Context, string) (*LeagueList, error){
		c.GetChallengerLeague,
		c.GetGrandmasterLeague,
		c.GetMasterLeague,
	}

	for _, getLeague := range leagues {
		league, err := getLeague(ctx, region)
		if err != nil {
			continue
		}

		for _, entry := range league.Entries {
			if entry.PUUID == puuid {
				return &entry, nil
			}
		}
	}

	return nil, fmt.Errorf("player with puuid %s not found in high elo leagues", puuid)
}
