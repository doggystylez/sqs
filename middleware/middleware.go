package middleware

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"time"

	"go.uber.org/zap"
	gotrace "golang.org/x/exp/trace"

	"github.com/labstack/echo/v4"
	"github.com/osmosis-labs/sqs/domain"
	"github.com/osmosis-labs/sqs/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// GoMiddleware represent the data-struct for middleware
type GoMiddleware struct {
	corsConfig         domain.CORSConfig
	flightRecordConfig domain.FlightRecordConfig
	logger             log.Logger
}

var (

	// flight recorder
	recordFlightOnce sync.Once
)

// CORS will handle the CORS middleware
func (m *GoMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", m.corsConfig.AllowedOrigin)
		c.Response().Header().Set("Access-Control-Allow-Headers", m.corsConfig.AllowedHeaders)
		c.Response().Header().Set("Access-Control-Allow-Methods", m.corsConfig.AllowedMethods)
		return next(c)
	}
}

// InitMiddleware initialize the middleware
func InitMiddleware(corsConfig *domain.CORSConfig, flightRecordConfig *domain.FlightRecordConfig, logger log.Logger) *GoMiddleware {
	return &GoMiddleware{
		corsConfig:         *corsConfig,
		flightRecordConfig: *flightRecordConfig,
		logger:             logger,
	}
}

// InstrumentMiddleware will handle the instrumentation middleware
func (m *GoMiddleware) InstrumentMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	// Set up the flight recorder.
	fr := gotrace.NewFlightRecorder()

	if m.flightRecordConfig.Enabled {
		err := fr.Start()
		if err != nil {
			m.logger.Error("failed to start flight recorder", zap.Error(err))
		}
	}

	return func(c echo.Context) error {
		start := time.Now()

		err := next(c)

		duration := time.Since(start)

		// Record outliers to the flight recorder for further analysis
		if m.flightRecordConfig.Enabled && duration > time.Duration(m.flightRecordConfig.TraceThresholdMS)*time.Millisecond {
			recordFlightOnce.Do(func() {
				// Note: we skip error handling since we don't want to interrupt the request handling
				// with tracing errors.

				// Grab the snapshot.
				var b bytes.Buffer
				_, err = fr.WriteTo(&b)
				if err != nil {
					m.logger.Error("failed to write trace to buffer", zap.Error(err))
					return
				}

				// Write it to a file.
				err = os.WriteFile(m.flightRecordConfig.TraceFileName, b.Bytes(), 0o755)
				if err != nil {
					m.logger.Error("failed to write trace to file", zap.Error(err))
					return
				}

				err = fr.Stop()
				if err != nil {
					fmt.Println("failed to stop flight recorder: ", err)
					m.logger.Error("failed to stop fligt recorder", zap.Error(err))
					return
				}
			})
		}

		return err
	}
}

// Middleware to capture request parameters
func (m *GoMiddleware) TraceWithParamsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get span from context (it is initialized by the preceding echo OTEL middleware)
			span := trace.SpanFromContext(c.Request().Context())

			// Iterate through query parameters and add them as attributes to the span
			// Ensure to filter out any sensitive parameters here
			for key, values := range c.QueryParams() {
				// As a simple approach, we're adding only the first value of each parameter
				// Consider handling multiple values differently if necessary
				span.SetAttributes(attribute.String(key, values[0]))
			}

			// Proceed with the request handling
			err := next(c)

			return err
		}
	}
}
