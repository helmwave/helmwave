package release

import (
	"context"

	"github.com/helmwave/helmwave/pkg/monitor"
	"github.com/invopop/jsonschema"
)

// MonitorFailedAction is a type for enumerating actions for handling failed monitors.
type MonitorFailedAction string

const (
	MonitorActionNone      MonitorFailedAction = ""
	MonitorActionRollback  MonitorFailedAction = "rollback"
	MonitorActionUninstall MonitorFailedAction = "uninstall"
)

func (MonitorFailedAction) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:    "string",
		Default: MonitorActionNone,
		Enum: []any{
			MonitorActionNone,
			MonitorActionRollback,
			MonitorActionUninstall,
		},
	}
}

type MonitorReference struct {
	Name   string              `yaml:"name" json:"name" jsonschema:"required"`
	Action MonitorFailedAction `yaml:"action" json:"action" jsonschema:"title=Action if monitor fails"`
}

func (rel *config) NotifyMonitorsFailed(ctx context.Context, mons ...monitor.Config) {
	action := MonitorActionNone

	for _, mon := range mons {
		allMons := rel.Monitors()
		for i := range allMons {
			monRef := allMons[i]
			if mon.Name() != monRef.Name {
				continue
			}

			if action != monRef.Action {
				if action != MonitorActionNone {
					rel.Logger().Warn("multiple actions to perform found, will use latest one")
				}
				action = monRef.Action
			}
		}
	}

	if action == MonitorActionNone {
		rel.Logger().Info("no actions will be performed for failed monitors")
	} else {
		rel.Logger().WithField("action", action).Info("chose action to perform for failed monitors")
		rel.performMonitorAction(ctx, action)
	}
}

func (rel *config) performMonitorAction(ctx context.Context, action MonitorFailedAction) {
	switch action {
	case MonitorActionRollback:
		err := rel.Rollback(ctx, 0)
		if err != nil {
			rel.Logger().WithError(err).Error("caught error while handling failed monitors")
		}
	case MonitorActionUninstall:
		_, err := rel.Uninstall(ctx)
		if err != nil {
			rel.Logger().WithError(err).Error("caught error while handling failed monitors")
		}
	default:
		rel.Logger().WithField("action", action).Error("unknown action to perform, skipping")
	}
}
