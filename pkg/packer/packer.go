package packer

import (
	"github/pm/pkg/reconciler"
)

type Packer struct {
	registeredReconcilers map[string]reconciler.Reconciler
}

func (p *Packer) Pack() {
}

func (p *Packer) Unpack() {
}
