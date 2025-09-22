package whitepaperdetailed

import (
        "errors"
        "fmt"
        "math"
        "sort"
        "strconv"
        "time"
)

// CoinSpecification captures the monetary, governance and operational
// parameters for the Synthron Coin.  It is designed so that external tooling
// (governance dashboards, wallets, treasury automation) can introspect the
// configuration without needing to parse the whitepaper prose.
type CoinSpecification struct {
        Symbol         string
        Name           string
        Precision      uint8
        MaxSupply      uint64
        InitialSupply  uint64
        Treasury       TreasuryPolicy
        FeeModel       FeeModel
        Emissions      []EmissionEpoch
        Governance     GovernancePolicy
        RiskControls   RiskControls
        Metadata       map[string]string
        ComplianceTags []string
}

// TreasuryPolicy describes where minted coins are directed and how reserves are
// protected once in circulation.
type TreasuryPolicy struct {
        TreasuryAccount      string
        InitialReserve       uint64
        OperatingReserveCap  uint64
        EmergencyReserveCap  uint64
        RefillThresholdRatio float64
}

// FeeModel defines how transaction fees are split and whether adaptive pricing
// is enabled.
type FeeModel struct {
        BaseFee        uint64
        DynamicPricing bool
        Allocation     FeeAllocation
}

// FeeAllocation determines the distribution of collected fees. Percentages are
// expressed as fractions that should sum to 1.0.
type FeeAllocation struct {
        Treasury   float64
        Validators float64
        Burn       float64
        Ecosystem  float64
}

// EmissionEpoch details the release of coins at a given block height.
type EmissionEpoch struct {
        Height      uint64
        Release     uint64
        Description string
}

// GovernancePolicy defines the decision-making mechanics for supply or policy
// updates.
type GovernancePolicy struct {
        QuorumRequirement      float64
        VotingPeriod           time.Duration
        ExecutionTimelock      time.Duration
        EmergencyVetoThreshold float64
        DelegateParticipation  float64
}

// RiskControls encode hard guards that must be honoured by minting endpoints
// and treasury automation workflows.
type RiskControls struct {
        CircuitBreakerThreshold float64
        DailyMintLimit          uint64
        MaxValidatorExposure    float64
        SettlementGracePeriod   time.Duration
}

// StakeSnapshot summarises the expected on-chain distribution at a given
// height. It allows operational tooling to reason about inflation in a single
// call without needing to load the entire emission schedule.
type StakeSnapshot struct {
        Height            uint64
        CirculatingSupply uint64
        RemainingEmissions uint64
}

const feeTolerance = 1e-6

// SynthronCoin returns the enterprise specification for the Synthron Coin. The
// specification embeds the emission schedule, treasury policy and governance
// configuration that drive production systems. The function always returns a
// self-consistent configuration; callers can rely on Validate for sanity checks
// when persisting the configuration externally.
func SynthronCoin() CoinSpecification {
        spec := CoinSpecification{
                Symbol:        "SYN",
                Name:          "Synthron Coin",
                Precision:     8,
                MaxSupply:     10_000_000_000,
                InitialSupply: 2_500_000_000,
                Treasury: TreasuryPolicy{
                        TreasuryAccount:      "syn_treasury_vault",
                        InitialReserve:       750_000_000,
                        OperatingReserveCap:  1_500_000_000,
                        EmergencyReserveCap:  500_000_000,
                        RefillThresholdRatio: 0.30,
                },
                FeeModel: FeeModel{
                        BaseFee:        10,
                        DynamicPricing: true,
                        Allocation: FeeAllocation{
                                Treasury:   0.35,
                                Validators: 0.40,
                                Burn:       0.15,
                                Ecosystem:  0.10,
                        },
                },
                Emissions: []EmissionEpoch{
                        {Height: 0, Release: 1_500_000_000, Description: "genesis liquidity provisioning and market makers"},
                        {Height: 210_000, Release: 1_250_000_000, Description: "validator incentives epoch 1"},
                        {Height: 420_000, Release: 1_250_000_000, Description: "ecosystem growth and grants"},
                        {Height: 630_000, Release: 1_000_000_000, Description: "cross-chain incentive alignment"},
                        {Height: 840_000, Release: 1_500_000_000, Description: "long tail developer and community funds"},
                },
                Governance: GovernancePolicy{
                        QuorumRequirement:      0.67,
                        VotingPeriod:           72 * time.Hour,
                        ExecutionTimelock:      24 * time.Hour,
                        EmergencyVetoThreshold: 0.80,
                        DelegateParticipation:  0.55,
                },
                RiskControls: RiskControls{
                        CircuitBreakerThreshold: 0.15,
                        DailyMintLimit:          50_000_000,
                        MaxValidatorExposure:    0.10,
                        SettlementGracePeriod:   48 * time.Hour,
                },
                Metadata: map[string]string{
                        "category":        "utility-governance-hybrid",
                        "supply_model":    "programmatic_emission_with_decay",
                        "audit_reference": "SYN-COIN-2024-Q4",
                },
                ComplianceTags: []string{"aml-compliant", "kyc-required", "multi-jurisdiction"},
        }
        spec.sortEmissions()
        return spec
}

// Validate ensures the specification is internally consistent. It performs a
// comprehensive set of checks so that configuration errors can be surfaced
// during development or CI rather than at runtime.
func (c CoinSpecification) Validate() error {
        if c.Symbol == "" || len(c.Symbol) > 8 {
                return errors.New("symbol must be 1-8 characters")
        }
        if c.Name == "" {
                return errors.New("coin name required")
        }
        if c.InitialSupply == 0 {
                return errors.New("initial supply must be greater than zero")
        }
        if c.InitialSupply > c.MaxSupply {
                return fmt.Errorf("initial supply %d exceeds max supply %d", c.InitialSupply, c.MaxSupply)
        }
        if c.Precision > 18 {
                return errors.New("precision cannot exceed 18 decimal places")
        }
        if c.Treasury.TreasuryAccount == "" {
                return errors.New("treasury account required")
        }
        if c.Treasury.InitialReserve > c.InitialSupply {
                return errors.New("initial reserve cannot exceed initial supply")
        }
        if c.Treasury.OperatingReserveCap+c.Treasury.EmergencyReserveCap > c.MaxSupply {
                return errors.New("aggregate reserve caps exceed max supply")
        }
        if c.Treasury.RefillThresholdRatio <= 0 || c.Treasury.RefillThresholdRatio >= 1 {
                return errors.New("treasury refill threshold must be between 0 and 1")
        }

        alloc := c.FeeModel.Allocation
        sum := alloc.Treasury + alloc.Validators + alloc.Burn + alloc.Ecosystem
        if math.Abs(1-sum) > feeTolerance {
                return fmt.Errorf("fee allocation must sum to 1.0, got %.6f", sum)
        }

        if c.Governance.QuorumRequirement <= 0.5 {
                return errors.New("governance quorum must exceed simple majority")
        }
        if c.Governance.QuorumRequirement > 1 || c.Governance.EmergencyVetoThreshold > 1 {
                return errors.New("governance ratios must be at most 1.0")
        }
        if c.Governance.EmergencyVetoThreshold < c.Governance.QuorumRequirement {
                return errors.New("emergency veto threshold must be >= quorum requirement")
        }
        if c.Governance.VotingPeriod < 24*time.Hour {
                return errors.New("voting period must be at least 24h")
        }
        if c.Governance.ExecutionTimelock < 6*time.Hour {
                return errors.New("execution timelock must allow governance review")
        }

        if c.RiskControls.CircuitBreakerThreshold <= 0 || c.RiskControls.CircuitBreakerThreshold >= 1 {
                return errors.New("circuit breaker threshold must be between 0 and 1")
        }
        if c.RiskControls.MaxValidatorExposure <= 0 || c.RiskControls.MaxValidatorExposure >= 1 {
                return errors.New("max validator exposure must be between 0 and 1")
        }
        if c.RiskControls.SettlementGracePeriod < 24*time.Hour {
                return errors.New("settlement grace period must be >= 24h")
        }

        if len(c.Emissions) == 0 {
                return errors.New("at least one emission epoch required")
        }
        var (
                previousHeight uint64
                cumulative      = c.InitialSupply
        )
        for i, epoch := range c.Emissions {
                if epoch.Release == 0 {
                        return fmt.Errorf("emission epoch %d release must be > 0", i)
                }
                if i > 0 && epoch.Height <= previousHeight {
                        return fmt.Errorf("emission epochs must be strictly increasing (index %d)", i)
                }
                if epoch.Description == "" {
                        return fmt.Errorf("emission epoch %d requires description", i)
                }
                cumulative += epoch.Release
                if cumulative > c.MaxSupply {
                        return fmt.Errorf("emission schedule exceeds max supply at height %d", epoch.Height)
                }
                previousHeight = epoch.Height
        }

        return nil
}

// CirculatingSupply computes the expected circulating supply at the provided
// block height accounting for scheduled emissions.
func (c CoinSpecification) CirculatingSupply(height uint64) uint64 {
        supply := c.InitialSupply
        for _, epoch := range c.Emissions {
                if epoch.Height > height {
                        break
                }
                supply += epoch.Release
        }
        if supply > c.MaxSupply {
                return c.MaxSupply
        }
        return supply
}

// RemainingEmissions returns the total emissions that have yet to unlock after
// the provided height.
func (c CoinSpecification) RemainingEmissions(height uint64) uint64 {
        var remaining uint64
        for _, epoch := range c.Emissions {
                if epoch.Height > height {
                        remaining += epoch.Release
                }
        }
        return remaining
}

// NextEmissionAfter returns the first emission scheduled strictly after the
// provided height.
func (c CoinSpecification) NextEmissionAfter(height uint64) (EmissionEpoch, bool) {
        for _, epoch := range c.Emissions {
                if epoch.Height > height {
                        return epoch, true
                }
        }
        return EmissionEpoch{}, false
}

// Snapshot summarises the emission state for the supplied height, exposing the
// circulating and remaining supply in a single structure.
func (c CoinSpecification) Snapshot(height uint64) StakeSnapshot {
        circulating := c.CirculatingSupply(height)
        return StakeSnapshot{
                Height:             height,
                CirculatingSupply:  circulating,
                RemainingEmissions: c.RemainingEmissions(height),
        }
}

// EmissionRatio returns the share of total supply that has been emitted by the
// supplied height.
func (c CoinSpecification) EmissionRatio(height uint64) float64 {
        if c.MaxSupply == 0 {
                return 0
        }
        return float64(c.CirculatingSupply(height)) / float64(c.MaxSupply)
}

// sortEmissions orders the emission epochs by height.
func (c *CoinSpecification) sortEmissions() {
        sort.SliceStable(c.Emissions, func(i, j int) bool {
                if c.Emissions[i].Height == c.Emissions[j].Height {
                        return c.Emissions[i].Description < c.Emissions[j].Description
                }
                return c.Emissions[i].Height < c.Emissions[j].Height
        })
}

// String renders a human readable identifier for the specification.
func (c CoinSpecification) String() string {
        return fmt.Sprintf("%s (%s) – max supply %s", c.Name, c.Symbol, strconv.FormatUint(c.MaxSupply, 10))
}
