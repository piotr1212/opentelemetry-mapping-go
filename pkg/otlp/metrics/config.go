// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"fmt"

	"github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/attributes/source"
)

type translatorConfig struct {
	// metrics export behavior
	HistMode                  HistogramMode
	SendHistogramAggregations bool
	Quantiles                 bool
	SendMonotonic             bool
	ResourceAttributesAsTags  bool
	// Deprecated: use InstrumentationScopeMetadataAsTags instead in favor of
	// https://github.com/open-telemetry/opentelemetry-proto/releases/tag/v0.15.0
	// Both must not be enabled at the same time.
	InstrumentationLibraryMetadataAsTags bool
	InstrumentationScopeMetadataAsTags   bool

	// cache configuration
	sweepInterval int64
	deltaTTL      int64

	fallbackSourceProvider source.Provider
}

// TranslatorOption is a translator creation option.
type TranslatorOption func(*translatorConfig) error

// WithDeltaTTL sets the delta TTL for cumulative metrics datapoints.
// By default, 3600 seconds are used.
func WithDeltaTTL(deltaTTL int64) TranslatorOption {
	return func(t *translatorConfig) error {
		if deltaTTL <= 0 {
			return fmt.Errorf("time to live must be positive: %d", deltaTTL)
		}
		t.deltaTTL = deltaTTL
		t.sweepInterval = 1
		if t.deltaTTL > 1 {
			t.sweepInterval = t.deltaTTL / 2
		}
		return nil
	}
}

// WithFallbackSourceProvider sets the fallback source provider.
// By default, an empty hostname is used as a fallback.
func WithFallbackSourceProvider(provider source.Provider) TranslatorOption {
	return func(t *translatorConfig) error {
		t.fallbackSourceProvider = provider
		return nil
	}
}

// WithQuantiles enables quantiles exporting for summary metrics.
func WithQuantiles() TranslatorOption {
	return func(t *translatorConfig) error {
		t.Quantiles = true
		return nil
	}
}

// WithResourceAttributesAsTags sets resource attributes as tags.
func WithResourceAttributesAsTags() TranslatorOption {
	return func(t *translatorConfig) error {
		t.ResourceAttributesAsTags = true
		return nil
	}
}

// WithInstrumentationLibraryMetadataAsTags sets instrumentation library metadata as tags.
func WithInstrumentationLibraryMetadataAsTags() TranslatorOption {
	return func(t *translatorConfig) error {
		t.InstrumentationLibraryMetadataAsTags = true
		return nil
	}
}

// WithInstrumentationScopeMetadataAsTags sets instrumentation scope metadata as tags.
func WithInstrumentationScopeMetadataAsTags() TranslatorOption {
	return func(t *translatorConfig) error {
		t.InstrumentationScopeMetadataAsTags = true
		return nil
	}
}

// HistogramMode is an export mode for OTLP Histogram metrics.
type HistogramMode string

const (
	// HistogramModeNoBuckets disables bucket export.
	HistogramModeNoBuckets HistogramMode = "nobuckets"
	// HistogramModeCounters exports buckets as Datadog counts.
	HistogramModeCounters HistogramMode = "counters"
	// HistogramModeDistributions exports buckets as Datadog distributions.
	HistogramModeDistributions HistogramMode = "distributions"
)

// WithHistogramMode sets the histograms mode.
// The default mode is HistogramModeOff.
func WithHistogramMode(mode HistogramMode) TranslatorOption {
	return func(t *translatorConfig) error {

		switch mode {
		case HistogramModeNoBuckets, HistogramModeCounters, HistogramModeDistributions:
			t.HistMode = mode
		default:
			return fmt.Errorf("unknown histogram mode: %q", mode)
		}
		return nil
	}
}

// WithCountSumMetrics exports .count and .sum histogram metrics.
// Deprecated: Use WithHistogramAggregations instead.
func WithCountSumMetrics() TranslatorOption {
	return WithHistogramAggregations()
}

// WithHistogramAggregations exports .count, .sum, .min and .max histogram metrics when available.
func WithHistogramAggregations() TranslatorOption {
	return func(t *translatorConfig) error {
		t.SendHistogramAggregations = true
		return nil
	}
}

// NumberMode is an export mode for OTLP Number metrics.
type NumberMode string

const (
	// NumberModeCumulativeToDelta calculates delta for
	// cumulative monotonic metrics in the client side and reports
	// them as Datadog counts.
	NumberModeCumulativeToDelta NumberMode = "cumulative_to_delta"

	// NumberModeRawValue reports the raw value for cumulative monotonic
	// metrics as a Datadog gauge.
	NumberModeRawValue NumberMode = "raw_value"
)

// WithNumberMode sets the number mode.
// The default mode is NumberModeCumulativeToDelta.
func WithNumberMode(mode NumberMode) TranslatorOption {
	return func(t *translatorConfig) error {
		switch mode {
		case NumberModeCumulativeToDelta:
			t.SendMonotonic = true
		case NumberModeRawValue:
			t.SendMonotonic = false
		default:
			return fmt.Errorf("unknown number mode: %q", mode)
		}
		return nil
	}
}
