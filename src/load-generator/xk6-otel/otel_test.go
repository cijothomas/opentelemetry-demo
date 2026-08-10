// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package xk6otel

import (
	"context"
	"testing"
	"time"

	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
)

type recordingLogger struct {
	embedded.Logger
	records []otellog.Record
}

func (l *recordingLogger) Enabled(context.Context, otellog.EnabledParameters) bool { return true }

func (l *recordingLogger) Emit(_ context.Context, record otellog.Record) {
	l.records = append(l.records, record.Clone())
}

func TestEmitPipelineProbe(t *testing.T) {
	recorder := &recordingLogger{}
	previousLogger := globalLogger
	globalLogger = recorder
	t.Cleanup(func() { globalLogger = previousLogger })

	before := time.Now()
	emitPipelineProbe()
	after := time.Now()

	if len(recorder.records) != 1 {
		t.Fatalf("got %d probe records, want 1", len(recorder.records))
	}

	record := recorder.records[0]
	if record.EventName() != "demo.telemetry.pipeline.probe" {
		t.Errorf("event name = %q, want %q", record.EventName(), "demo.telemetry.pipeline.probe")
	}
	if body := record.Body().AsString(); body != "Telemetry pipeline latency probe" {
		t.Errorf("body = %q, want %q", body, "Telemetry pipeline latency probe")
	}
	if record.Timestamp().Before(before) || record.Timestamp().After(after) {
		t.Errorf("timestamp %v is outside emission interval [%v, %v]", record.Timestamp(), before, after)
	}
	if !record.ObservedTimestamp().Equal(record.Timestamp()) {
		t.Errorf("observed timestamp %v does not match timestamp %v", record.ObservedTimestamp(), record.Timestamp())
	}
}
