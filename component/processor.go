// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package component

import (
	"context"

	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component/componenterror"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
)

// Processor defines the common functions that must be implemented by TracesProcessor
// and MetricsProcessor.
type Processor interface {
	Component

	// GetCapabilities must return the capabilities of the processor.
	GetCapabilities() ProcessorCapabilities
}

// TracesProcessor is a processor that can consume traces.
type TracesProcessor interface {
	Processor
	consumer.Traces
}

// MetricsProcessor is a processor that can consume metrics.
type MetricsProcessor interface {
	Processor
	consumer.Metrics
}

// LogsProcessor is a processor that can consume logs.
type LogsProcessor interface {
	Processor
	consumer.Logs
}

// ProcessorCapabilities describes the capabilities of a Processor.
type ProcessorCapabilities struct {
	// MutatesConsumedData is set to true if Consume* function of the
	// processor modifies the input TraceData or MetricsData argument.
	// Processors which modify the input data MUST set this flag to true. If the processor
	// does not modify the data it MUST set this flag to false. If the processor creates
	// a copy of the data before modifying then this flag can be safely set to false.
	MutatesConsumedData bool
}

// ProcessorCreateParams is passed to Create* functions in ProcessorFactory.
type ProcessorCreateParams struct {
	// Logger that the factory can use during creation and can pass to the created
	// component to be used later as well.
	Logger *zap.Logger

	// BuildInfo can be used by components for informational purposes
	BuildInfo BuildInfo
}

// ProcessorFactory is factory interface for processors. This is the
// new factory type that can create new style processors.
//
// This interface cannot be directly implemented, implementations need to embed
// the BaseProcessorFactory or use the processorhelper.NewFactory to implement it.
type ProcessorFactory interface {
	Factory

	// CreateDefaultConfig creates the default configuration for the Processor.
	// This method can be called multiple times depending on the pipeline
	// configuration and should not cause side-effects that prevent the creation
	// of multiple instances of the Processor.
	// The object returned by this method needs to pass the checks implemented by
	// 'configcheck.ValidateConfig'. It is recommended to have such check in the
	// tests of any implementation of the Factory interface.
	CreateDefaultConfig() config.Processor

	// CreateTracesProcessor creates a trace processor based on this config.
	// If the processor type does not support tracing or if the config is not valid
	// error will be returned instead.
	CreateTracesProcessor(
		ctx context.Context,
		params ProcessorCreateParams,
		cfg config.Processor,
		nextConsumer consumer.Traces,
	) (TracesProcessor, error)

	// CreateMetricsProcessor creates a metrics processor based on this config.
	// If the processor type does not support metrics or if the config is not valid
	// error will be returned instead.
	CreateMetricsProcessor(
		ctx context.Context,
		params ProcessorCreateParams,
		cfg config.Processor,
		nextConsumer consumer.Metrics,
	) (MetricsProcessor, error)

	// CreateLogsProcessor creates a processor based on the config.
	// If the processor type does not support logs or if the config is not valid
	// error will be returned instead.
	CreateLogsProcessor(
		ctx context.Context,
		params ProcessorCreateParams,
		cfg config.Processor,
		nextConsumer consumer.Logs,
	) (LogsProcessor, error)

	// unexportedProcessor is a dummy method to force this interface to not be implemented.
	unexportedProcessor()
}

// BaseProcessorFactory is the interface that must be embedded by all ProcessorFactory implementations.
type BaseProcessorFactory struct{}

var _ ProcessorFactory = (*BaseProcessorFactory)(nil)

// Type must be override.
func (b BaseProcessorFactory) Type() config.Type {
	panic("implement me")
}

// CreateDefaultConfig must be override.
func (b BaseProcessorFactory) CreateDefaultConfig() config.Processor {
	panic("implement me")
}

// CreateTracesProcessor default implemented as not supported date type.
func (b BaseProcessorFactory) CreateTracesProcessor(context.Context, ProcessorCreateParams, config.Processor, consumer.Traces) (TracesProcessor, error) {
	return nil, componenterror.ErrDataTypeIsNotSupported
}

// CreateMetricsProcessor default implemented as not supported date type.
func (b BaseProcessorFactory) CreateMetricsProcessor(context.Context, ProcessorCreateParams, config.Processor, consumer.Metrics) (MetricsProcessor, error) {
	return nil, componenterror.ErrDataTypeIsNotSupported
}

// CreateLogsProcessor default implemented as not supported date type.
func (b BaseProcessorFactory) CreateLogsProcessor(context.Context, ProcessorCreateParams, config.Processor, consumer.Logs) (LogsProcessor, error) {
	return nil, componenterror.ErrDataTypeIsNotSupported
}

func (b BaseProcessorFactory) unexportedProcessor() {}
