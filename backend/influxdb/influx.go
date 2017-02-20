package influxdb

import (
	"context"
	"encoding/json"
	"fmt"
	client "github.com/influxdata/influxdb/client/v2"
	"sync"
)

var clients map[string]client.Client
var lock sync.RWMutex

func Backend(ctx context.Context, query string) ([]float64, error) {
	// Pull config from ctx
	config, ok := FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("Context missing InfluxDB configuration.")
	}

	lock.RLock()
	c, ok := clients[config.HTTPConfig.Addr]
	lock.RUnlock()
	if !ok {
		lock.Lock()
		if clients == nil {
			clients = map[string]client.Client{}
		}
		var err error
		c, err = client.NewHTTPClient(config.HTTPConfig)
		if err != nil || c == nil {
			lock.Unlock()
			return nil, fmt.Errorf("Unable to create new InfluxDB client with provided configuration: %v", err)
		}
		clients[config.HTTPConfig.Addr] = c
		lock.Unlock()
	}

	resp, err := c.Query(client.NewQueryWithParameters(query, config.Database, config.Precision, config.Parameters))
	if err != nil {
		return nil, err
	}
	if err = resp.Error(); err != nil {
		return nil, err
	}

	if len(resp.Results) != 1 {
		return nil, fmt.Errorf(`The query does not result in a single series`)
	}
	r := resp.Results[0]
	if len(resp.Results) != 1 {
		return nil, fmt.Errorf(`The query does not result in a single series`)
	}

	data := []float64{}
	for _, row := range r.Series {
		for _, pt := range row.Values {
			n, _ := pt[1].(json.Number).Float64()
			data = append(data, n)
		}
	}
	return data, nil

}

const configKey = `influxconfig`

func NewContext(ctx context.Context, config BackendConfig) context.Context {
	return context.WithValue(ctx, configKey, config)
}

func FromContext(ctx context.Context) (BackendConfig, bool) {
	v, ok := ctx.Value(configKey).(BackendConfig)
	return v, ok
}

type BackendConfig struct {
	HTTPConfig client.HTTPConfig
	Database   string
	Precision  string
	Parameters map[string]interface{}
}
