package tests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary //nolint:gochecknoglobals // skip

type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewAPIClient creates a new APIClient with the given baseURL and httpClient.
// If httpClient is nil, http.DefaultClient is used.
func NewAPIClient(
	baseURL string,
	httpClient *http.Client,
) APIClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return APIClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (a APIClient) Get(
	ctx context.Context,
	endpoint string,
	headers http.Header,
	dest any,
	errDest any,
) (*http.Response, error) {
	return a.httpRequest(ctx, http.MethodGet, endpoint, headers, http.NoBody, dest, errDest)
}

func (a APIClient) Post(
	ctx context.Context,
	endpoint string,
	headers http.Header,
	request any,
	dest any,
	errDest any,
) (*http.Response, error) {
	b, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	return a.httpRequest(ctx, http.MethodPost, endpoint, headers, bytes.NewReader(b), dest, errDest)
}

func (a APIClient) PostJSON(
	ctx context.Context,
	endpoint string,
	headers http.Header,
	requestJSON string,
	dest any,
	errDest any,
) (*http.Response, error) {
	b := []byte(requestJSON)
	return a.httpRequest(ctx, http.MethodPost, endpoint, headers, bytes.NewReader(b), dest, errDest)
}

func (a APIClient) Put(
	ctx context.Context,
	endpoint string,
	headers http.Header,
	request any,
	dest any,
	errDest any,
) (*http.Response, error) {
	b, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	return a.httpRequest(ctx, http.MethodPut, endpoint, headers, bytes.NewReader(b), dest, errDest)
}

func (a APIClient) Patch(
	ctx context.Context,
	endpoint string,
	headers http.Header,
	request any,
	dest any,
	errDest any,
) (*http.Response, error) {
	b, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	return a.httpRequest(ctx, http.MethodPatch, endpoint, headers, bytes.NewReader(b), dest, errDest)
}

func (a APIClient) Delete(
	ctx context.Context,
	endpoint string,
	headers http.Header,
	dest any,
	errDest any,
) (*http.Response, error) {
	return a.httpRequest(ctx, http.MethodDelete, endpoint, headers, http.NoBody, dest, errDest)
}

func (a APIClient) DeleteWithBody(
	ctx context.Context,
	endpoint string,
	headers http.Header,
	request any,
	dest any,
	errDest any,
) (*http.Response, error) {
	b, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	return a.httpRequest(ctx, http.MethodDelete, endpoint, headers, bytes.NewReader(b), dest, errDest)
}

func (a APIClient) MultiForm(
	ctx context.Context,
	endpoint string,
	headers http.Header,
	payload io.Reader,
	dest any,
	errDest any,
) (*http.Response, error) {
	return a.httpRequest(ctx, http.MethodPost, endpoint, headers, payload, dest, errDest)
}

func (a APIClient) httpRequest(
	ctx context.Context,
	httpMethod string,
	endpoint string,
	headers http.Header,
	payload io.Reader,
	dest any,
	errDest any,
) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, httpMethod, a.baseURL+endpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	// Set default Content-Type for POST if not provided
	if httpMethod == http.MethodPost {
		if headers == nil || headers.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	for k, values := range headers {
		for _, v := range values {
			req.Header.Add(k, v)
		}
	}

	logRequest(req)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpClient.Do: %w", err)
	}

	logResponse(resp)
	defer resp.Body.Close()

	if err = parseResponse(resp, dest, errDest); err != nil {
		return nil, fmt.Errorf("parseResponse: %w", err)
	}

	return resp, nil
}

func parseResponse(r *http.Response, dest, errDest any) error {
	// Успешный ответ
	if r.StatusCode >= http.StatusOK && r.StatusCode < http.StatusMultipleChoices {
		if dest != nil {
			// Проверим, есть ли вообще тело
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return fmt.Errorf("read body: %w", err)
			}
			if len(body) == 0 {
				return nil // тело пустое — просто возвращаем OK
			}
			if err := json.Unmarshal(body, dest); err != nil {
				return fmt.Errorf("json.Unmarshal(success destination): %w", err)
			}
		}
	} else if errDest != nil {
		// Для ошибок допустим пустой body тоже
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("read error body: %w", err)
		}
		if len(body) == 0 {
			return nil
		}
		if err := json.Unmarshal(body, errDest); err != nil {
			return fmt.Errorf("json.Unmarshal(err destination): %w", err)
		}
	}
	return nil
}


func logRequest(r *http.Request) {
	log.Printf("Request:  %s %s", r.Method, r.URL)
}

func logResponse(r *http.Response) {
	rawResponse, err := httputil.DumpResponse(r, true)
	if err == nil {
		log.Println("Response:", string(rawResponse))
	}
}
