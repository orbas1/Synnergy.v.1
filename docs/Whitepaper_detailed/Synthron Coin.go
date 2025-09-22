package whitepaperdetailed

import (
	"context"

	"synnergy/treasury"
)

// NewSynthronCoinTreasury constructs a fresh treasury instance for whitepaper
// examples. The helper mirrors the runtime behaviour exposed through the CLI
// and web dashboards while keeping documentation snippets concise.
func NewSynthronCoinTreasury(ctx context.Context) (*treasury.SynthronTreasury, error) {
	return treasury.NewSynthronTreasury(ctx)
}

// DefaultSynthronCoinTreasury returns the shared treasury singleton used across
// documentation-driven tests. It aligns with the CLI behaviour where Stage 80
// telemetry is initialised lazily on demand.
func DefaultSynthronCoinTreasury(ctx context.Context) (*treasury.SynthronTreasury, error) {
	return treasury.DefaultSynthronTreasury(ctx)
}

// TreasuryTelemetry provides a convenience wrapper so guides can reference the
// structured diagnostics emitted by the treasury without repeating boilerplate.
func TreasuryTelemetry(ctx context.Context) (treasury.SynthronTreasuryTelemetry, error) {
	treas, err := DefaultSynthronCoinTreasury(ctx)
	if err != nil {
		return treasury.SynthronTreasuryTelemetry{}, err
	}
	return treas.Diagnostics(ctx), nil
}

// TreasurySummary renders the short form diagnostics used throughout the
// whitepaper to describe Stage 80 orchestration outcomes.
func TreasurySummary(ctx context.Context) (string, error) {
	diag, err := TreasuryTelemetry(ctx)
	if err != nil {
		return "", err
	}
	return treasury.SynthronTreasurySummary(diag), nil
}

// TreasuryAuditTrail exposes the signed event log so documentation snippets can
// demonstrate tamper-evident provenance for Stage 80 monetary actions.
func TreasuryAuditTrail(ctx context.Context) ([]treasury.SynthronCoinEvent, error) {
	treas, err := DefaultSynthronCoinTreasury(ctx)
	if err != nil {
		return nil, err
	}
	return treas.AuditTrail(), nil
}

// TreasuryAuthorizeOperator promotes an address to operator status for Stage 80
// narratives.
func TreasuryAuthorizeOperator(ctx context.Context, addr string) error {
	treas, err := DefaultSynthronCoinTreasury(ctx)
	if err != nil {
		return err
	}
	return treas.AuthorizeOperator(ctx, addr)
}

// TreasuryRevokeOperator demotes an operator, preserving the guardian wallet.
func TreasuryRevokeOperator(ctx context.Context, addr string) error {
	treas, err := DefaultSynthronCoinTreasury(ctx)
	if err != nil {
		return err
	}
	return treas.RevokeOperator(ctx, addr)
}

// TreasuryWithOperator annotates documentation contexts with operator identity.
func TreasuryWithOperator(ctx context.Context, operator string) context.Context {
	return treasury.WithOperator(ctx, operator)
}
