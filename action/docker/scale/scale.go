// This package does not support edge based actions. Taking action on edges makes sense for alarms
// and other state tracking integrations, scaling operations should continue and repeat periodically
// while a monitor is in an ALARM state. This should ensure that the step increase or decrease 
// will eventually reach some boundary beyond which operations stabilize.
package scale

import (
	"github.com/buildertools/siren"
	"github.com/buildertools/siren/action/common"
	"github.com/buildertools/svctools-go/clients"
	docker "github.com/docker/docker/client"
	dockertypes "github.com/docker/docker/api/types"
	"context"
	"fmt"
	"time"
)

func Alarm(ctx context.Context, m siren.Monitor) {
	scale(ctx, m, siren.STATE_ALARM)
}

func Flapping(ctx context.Context, m siren.Monitor) {
	scale(ctx, m, siren.STATE_FLAPPING)
}

func scale(ctx context.Context, m siren.Monitor, s int) {
	cfg, ok := FromContext(ctx)
	if !ok {
		panic(fmt.Errorf(`No configuration provided for docker-scale`))
	}

	// determine new scale delta
	var scaleDelta int
	switch(s) {
	case siren.STATE_ALARM:
		scaleDelta = cfg.Step
	case siren.STATE_FLAPPING:
		if !cfg.ScaleOnFlap {
			return
		}
		scaleDelta = cfg.Step
	case siren.STATE_CLEAR:
		panic(fmt.Errorf(`Cannot perform scale on CLEAR state`))
		return
	default:
		panic(fmt.Errorf(`Unsupported state passed to docker-scale`))
	}

	// Push desired state
	_, err := clients.RetryLinear(
			func() (interface{}, clients.ClientError) {
				fmt.Println("Running service calls")
				// Connect to Docker endpoint
				dc, err := docker.NewEnvClient()
				if err != nil {
					common.PanicOrLog(cfg.EnablePanic, cfg.EnableLogging, err)
					return nil, clients.RetriableError{E:err}
				}

				// Pull state and scale of service
				service, _, err := dc.ServiceInspectWithRaw(ctx, cfg.ServiceID)
				if err != nil {
					common.PanicOrLog(cfg.EnablePanic, cfg.EnableLogging, err)
					return nil, clients.RetriableError{E:err}
				}

				// Generate updated config
				serviceMode := &service.Spec.Mode
				if serviceMode.Replicated == nil {
					common.PanicOrLog(
						cfg.EnablePanic,
						cfg.EnableLogging,
						fmt.Errorf("scale can only be used with replicated mode"))
					return nil, clients.RetriableError{E:fmt.Errorf(`Service not replicated`)}
				}
				scale := int64(*serviceMode.Replicated.Replicas) + int64(scaleDelta)
				if scale < int64(cfg.Floor) || scale < 0 {
					scale = int64(cfg.Floor)
				}
				if scale > int64(cfg.Ceiling) {
					scale = int64(cfg.Ceiling)
				}
				var tscale uint64 = uint64(scale)
				serviceMode.Replicated.Replicas = &tscale

				_, err = dc.ServiceUpdate(
					ctx, 
					service.ID, 
					service.Version, 
					service.Spec, 
					dockertypes.ServiceUpdateOptions{})
				if err != nil {
					return nil, clients.RetriableError{E: err}
				}
				return nil, nil
			},
			cfg.Timeout,
			cfg.RetryInterval,
			cfg.RetryJitter,
		)
	if err != nil {
		common.PanicOrLog(cfg.EnablePanic, cfg.EnableLogging, err)
		return
	}
}

const ctxKey = `docker-scale`
type DockerScaleConfig struct {
	ServiceID     string
	Floor         int
	Ceiling       int
	Step          int
	ScaleOnFlap   bool
	Timeout       time.Duration
	RetryInterval time.Duration
	RetryJitter   time.Duration
	EnableLogging bool
	EnablePanic   bool
}

func NewContext(ctx context.Context, cfg DockerScaleConfig) context.Context {
	return context.WithValue(ctx, ctxKey, cfg)
}

func FromContext(ctx context.Context) (DockerScaleConfig, bool) {
	v, ok := ctx.Value(ctxKey).(DockerScaleConfig)
	return v, ok
}
