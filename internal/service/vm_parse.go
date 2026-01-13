package service

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type vmResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []any             `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

type instantSample struct {
	Metric     map[string]string
	ValueFloat float64
}

func parseInstantVector(raw []byte) ([]instantSample, error) {
	var resp vmResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	if resp.Status != "success" {
		return nil, fmt.Errorf("vm status=%s", resp.Status)
	}

	out := make([]instantSample, 0, len(resp.Data.Result))
	for _, r := range resp.Data.Result {
		if len(r.Value) < 2 {
			continue
		}
		valStr, ok := r.Value[1].(string)
		if !ok {
			continue
		}
		f, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			continue
		}
		out = append(out, instantSample{
			Metric:     r.Metric,
			ValueFloat: f,
		})
	}
	return out, nil
}
