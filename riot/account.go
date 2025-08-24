package riot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

func (c *Client) GetAccountByRiotID(ctx context.Context, cluster, gameName, tagLine string) (*Account, error) {
	encodedGameName := url.QueryEscape(gameName)
	encodedTagLine := url.QueryEscape(tagLine)
	endpoint := fmt.Sprintf("%s/riot/account/v1/accounts/by-riot-id/%s/%s",
		c.getClusterURL("account", cluster), encodedGameName, encodedTagLine)

	body, err := c.makeRequest(ctx, "GET", endpoint)
	if err != nil {
		return nil, fmt.Errorf("get account by riot id %s#%s: %w", gameName, tagLine, err)
	}

	var account Account
	if err := json.Unmarshal(body, &account); err != nil {
		return nil, fmt.Errorf("unmarshaling account: %w", err)
	}

	return &account, nil
}

func (c *Client) GetAccountByPUUID(ctx context.Context, cluster, puuid string) (*Account, error) {
	endpoint := fmt.Sprintf("%s/riot/account/v1/accounts/by-puuid/%s",
		c.getClusterURL("account", cluster), puuid)

	body, err := c.makeRequest(ctx, "GET", endpoint)
	if err != nil {
		return nil, fmt.Errorf("get account by puuid %s: %w", puuid, err)
	}

	var account Account
	if err := json.Unmarshal(body, &account); err != nil {
		return nil, fmt.Errorf("unmarshaling account: %w", err)
	}

	return &account, nil
}

func RegionToCluster(region string) string {
	clusterMap := map[string]string{
		"br1":  "americas",
		"la1":  "americas",
		"la2":  "americas",
		"na1":  "americas",
		"oc1":  "americas",
		"eun1": "europe",
		"euw1": "europe",
		"tr1":  "europe",
		"ru":   "europe",
		"jp1":  "asia",
		"kr":   "asia",
		"sg2":  "asia",
		"tw2":  "asia",
		"vn2":  "asia",
	}

	if cluster, exists := clusterMap[region]; exists {
		return cluster
	}
	return "americas"
}
