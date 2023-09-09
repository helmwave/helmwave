package release

import (
	"context"

	"github.com/helmwave/helmwave/pkg/monitor"
)

const (
	ActionNone      = ""
	ActionRollback  = "rollback"
	ActionUninstall = "uninstall"
)

type MonitorReference struct {
	Name   string `yaml:"name" json:"name" jsonschema:"required"`
	Action string `yaml:"action" json:"action" jsonschema:"enum=,enum=rollback,enum=uninstall,title=Action if monitor fails"`
}

func (rel *config) NotifyMonitorsFailed(ctx context.Context, mons ...monitor.Config) {
	action := ActionNone

	for _, mon := range mons {
		allMons := rel.Monitors()
		for i := range allMons {
			monRef := allMons[i]
			if mon.Name() != monRef.Name {
				continue
			}

			if action != monRef.Action {
				if action != ActionNone {
					rel.Logger().Warn("multiple actions to perform found, will use latest one")
				}
				action = monRef.Action
			}
		}
	}

	if action == ActionNone {
		rel.Logger().Info("no actions will be performed for failed monitors")
	} else {
		rel.Logger().WithField("action", action).Info("chose action to perform for failed monitors")
		rel.performMonitorAction(ctx, action)
	}
}

func (rel *config) performMonitorAction(ctx context.Context, action string) {
	switch action {
	case ActionRollback:
		err := rel.Rollback(ctx, 0)
		if err != nil {
			rel.Logger().WithError(err).Error("caught error while handling failed monitors")
		}
	case ActionUninstall:
		_, err := rel.Uninstall(ctx)
		if err != nil {
			rel.Logger().WithError(err).Error("caught error while handling failed monitors")
		}
	default:
		rel.Logger().WithField("action", action).Error("unknown action to perform, skipping")
	}
}
