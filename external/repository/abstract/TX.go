package abstract

// copied from github.com/AdamLuzsi/frameless
import "context"

type OnePhaseCommitProtocol interface {
	// BeginTx creates a context with a transaction.
	// All statements that receive this context should be executed within the given transaction in the context.
	// After a BeginTx command will be executed in a single transaction until an explicit COMMIT or ROLLBACK is given.
	BeginTx(ctx context.Context) (context.Context, error)
	// Commit commits the current transaction.
	// All changes made by the transaction become visible to others and are guaranteed to be durable if a crash occurs.
	CommitTx(context.Context) error
	// Rollback rolls back the current transaction and causes all the updates made by the transaction to be discarded.
	RollbackTx(context.Context) error
}
