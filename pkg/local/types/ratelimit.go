package types

type LocalRateLimit struct {
	StatPrefix     string                                   `yaml:"stat_prefix,omitempty" json:"stat_prefix,omitempty"`
	TokenBucket    *LocalRateLimit_TokenBucket              `yaml:"token_bucket,omitempty" json:"token_bucket,omitempty"`
	FilterEnabled  *LocalRateLimit_RuntimeFractionalPercent `yaml:"filter_enabled,omitempty" json:"filter_enabled,omitempty"`
	FilterEnforced *LocalRateLimit_RuntimeFractionalPercent `yaml:"filter_enforced,omitempty" json:"filter_enforced,omitempty"`
}

type LocalRateLimit_TokenBucket struct {
	MaxTokens     int    `yaml:"max_tokens,omitempty" json:"max_tokens,omitempty"`
	TokensPerFill int    `yaml:"tokens_per_fill,omitempty" json:"tokens_per_fill,omitempty"`
	FillInterval  string `yaml:"fill_interval,omitempty" json:"fill_interval,omitempty"`
}

type LocalRateLimit_RuntimeFractionalPercent struct {
	DefaultValue *LocalRateLimit_FractionalPercent `yaml:"default_value,omitempty" json:"default_value,omitempty"`
	RuntimeKey   string                            `yaml:"runtime_key,omitempty" json:"runtime_key,omitempty"`
}

type LocalRateLimit_FractionalPercent struct {
	Numerator   int                                              `yaml:"numerator,omitempty" json:"numerator,omitempty"`
	Denominator LocalRateLimit_FractionalPercent_DenominatorType `yaml:"denominator,omitempty" json:"denominator,omitempty"`
}

type LocalRateLimit_FractionalPercent_DenominatorType string

const (
	FractionalPercent_HUNDRED      LocalRateLimit_FractionalPercent_DenominatorType = "HUNDRED"
	FractionalPercent_TEN_THOUSAND LocalRateLimit_FractionalPercent_DenominatorType = "TEN_THOUSAND"
	FractionalPercent_MILLION      LocalRateLimit_FractionalPercent_DenominatorType = "MILLION"
)

