package consts

const (
	// DefaultFeatureGates is the default feature gates setting that should be
	// provided if none are provided by the user. This generally includes features
	// that are innocuous, or otherwise don't actually get triggered unless the
	// user takes further action.
	DefaultFeatureGates = "GatewayAlpha=true"

	// ConformanceExpressionRoutesTestsFeatureGates is the set of feature gates to be used
	// when running conformance tests with expression routes enabled.
	ConformanceExpressionRoutesTestsFeatureGates = "GatewayAlpha=true,ExpressionRoutes=true"
)
