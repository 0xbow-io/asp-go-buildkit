package chainalysis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	APIV2Endpoint = "https://api.chainalysis.com/api/kyt/v2/"
)

type APIV2 interface {
	RegisterTransfer(ctx context.Context, userID string, body TransferRegisterReq) (TransferRegisterResp, error)
	GetTransferSummary(ctx context.Context, transferID string) (TransferRegisterResp, error)
	GetDirectExposure(ctx context.Context, transferID string) (Exposure, error)
	GetAlerts(ctx context.Context, transferID string) ([]Alert, error)
}

type apiV2 struct {
	apiKey string
	client *http.Client
}

func executeAndDecode[T any](client *http.Client, req *http.Request) (T, error) {
	var out T
	resp, err := client.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return out, err
	}

	// For debugging: print the raw response
	fmt.Printf("Raw response: %s\n", string(bodyBytes))

	err = json.Unmarshal(bodyBytes, &out)
	if err != nil {
		return out, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return out, nil
}

func NewAPIV2(apiKey string) APIV2 {
	return &apiV2{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (a apiV2) request(
	ctx context.Context, method, path string, body any) (*http.Request, error) {

	uri, err := url.JoinPath(APIV2Endpoint, path)
	if err != nil {
		return nil, err
	}

	var buff bytes.Buffer
	if body != nil {
		err = json.NewEncoder(&buff).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, uri, &buff)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Token", a.apiKey)
	return req, nil
}

func (a apiV2) RegisterTransfer(
	ctx context.Context,
	userID string,
	body TransferRegisterReq) (TransferRegisterResp, error) {

	path := fmt.Sprintf("users/%s/transfers", userID)
	req, err := a.request(ctx, http.MethodPost, path, body)
	if err != nil {
		return TransferRegisterResp{}, err
	}
	return executeAndDecode[TransferRegisterResp](a.client, req)
}

func (a apiV2) GetTransferSummary(ctx context.Context, externalId string) (TransferRegisterResp, error) {
	path := fmt.Sprintf("transfers/%s", externalId)

	req, err := a.request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return TransferRegisterResp{}, err
	}
	return executeAndDecode[TransferRegisterResp](a.client, req)
}

func (a apiV2) GetDirectExposure(ctx context.Context, externalId string) (Exposure, error) {
	path := fmt.Sprintf("transfers/%s/exposures", externalId)

	req, err := a.request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return Exposure{}, err
	}
	out, err := executeAndDecode[map[string]Exposure](a.client, req)
	if err != nil {
		return Exposure{}, err
	}
	exposure, ok := out["direct"]
	if !ok {
		return Exposure{}, fmt.Errorf("no direct exposure found")
	}
	return exposure, nil
}

func (a apiV2) GetAlerts(ctx context.Context, externalId string) ([]Alert, error) {
	path := fmt.Sprintf("transfers/%s/alerts", externalId)

	req, err := a.request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	out, err := executeAndDecode[map[string][]Alert](a.client, req)
	if err != nil {
		return nil, err
	}
	alerts, ok := out["alerts"]
	if !ok {
		return nil, fmt.Errorf("no alerts found")
	}
	return alerts, nil
}
