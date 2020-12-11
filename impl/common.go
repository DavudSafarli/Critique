package impl

import (
	"context"
	"github.com/DavudSafarli/Critique/domain/contracts"
)

type common struct {
	txer contracts.OnePhaseCommitProtocol
}

func (c common) commitOrRollback(ctx context.Context, err error) {
	if c.txer == nil {
		panic("common/txer is not initialized")
	}
	if err != nil {
		c.txer.RollbackTx(ctx)
		return
	}
	err = c.txer.CommitTx(ctx)
}
