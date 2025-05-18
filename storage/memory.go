package storage

import (
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MetricSample struct {
	Name      string
	Value     float64
	Type      string
	Tags      []string
	Timestamp time.Time
}

type CardinalityEntry struct {
	Metric         string  `json:"metric"`
	TagKey         string  `json:"tag_key"`
	UniqueVals     int     `json:"unique_values"`
	TotalSamples   int     `json:"total_samples"`
	CardinalityPct float64 `json:"cardinality_pct"`
}

var (
	metrics        []MetricSample
	tagCardinality = make(map[string]map[string]map[string]struct{})
	mu             sync.Mutex
)

func AddMetric(name, valueStr, typ string, tags []string) {
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return
	}

	sample := MetricSample{
		Name:      name,
		Value:     value,
		Type:      typ,
		Tags:      tags,
		Timestamp: time.Now(),
	}

	mu.Lock()
	metrics = append(metrics, sample)

	for _, tag := range tags {
		print(tag)
		parts := splitTag(tag)
		if len(parts) == 2 {
			key, val := parts[0], parts[1]
			if _, ok := tagCardinality[name]; !ok {
				tagCardinality[name] = make(map[string]map[string]struct{})
			}
			if _, ok := tagCardinality[name][key]; !ok {
				tagCardinality[name][key] = make(map[string]struct{})
			}
			tagCardinality[name][key][val] = struct{}{}
		}
	}
	mu.Unlock()
}

func GetMetrics() []MetricSample {
	mu.Lock()
	defer mu.Unlock()
	return append([]MetricSample(nil), metrics...)
}

func ResetMetrics() {
	mu.Lock()
	metrics = nil
	tagCardinality = make(map[string]map[string]map[string]struct{})
	mu.Unlock()
}

func GetCardinality() []CardinalityEntry {
	mu.Lock()
	defer mu.Unlock()
	var result []CardinalityEntry

	totalByMetric := make(map[string]int)
	for _, sample := range metrics {
		totalByMetric[sample.Name]++
	}

	for metric, tagKeys := range tagCardinality {
		totalSamples := totalByMetric[metric]
		for key, values := range tagKeys {
			count := len(values)
			pct := 0.0
			if totalSamples > 0 {
				pct = float64(count) / float64(totalSamples)
			}
			result = append(result, CardinalityEntry{
				Metric:         metric + "." + key,
				TagKey:         key,
				UniqueVals:     count,
				TotalSamples:   totalSamples,
				CardinalityPct: pct,
			})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CardinalityPct > result[j].CardinalityPct
	})

	return result
}

func splitTag(tag string) []string {
	tag = strings.TrimSpace(tag)
	parts := strings.SplitN(tag, ":", 2)
	if len(parts) != 2 {
		return nil
	}
	return parts
}
