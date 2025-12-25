package edge

// Status is a result of a check.
type Status = int

const (
	StatusOk = Status(iota)
	StatusBad
)

// CheckFunc is a function used to check.
type CheckFunc = func() Status

// HealthChecker provides health check.
type HealthChecker = interface {
	CheckHealth() Status
}

// ReadinessChecker provides readiness check.
type ReadinessChecker = interface {
	CheckReadiness() Status
}

// K8SChecker is a full standard k8s probes checker.
type K8SChecker = interface {
	HealthChecker
	ReadinessChecker
}
