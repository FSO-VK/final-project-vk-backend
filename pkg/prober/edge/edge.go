// Package edge implements multiple probe checks.
package edge

func checkAll(checks []CheckFunc) Status {
	if checks == nil {
		return StatusOk
	}
	for _, hc := range checks {
		partialStatus := hc()
		if partialStatus != StatusOk {
			return partialStatus
		}
	}
	return StatusOk
}

// Edge is an application edge.
type Edge struct {
	healthChecks    []CheckFunc
	readinessChecks []CheckFunc
}

// CheckHealth is a HealthChecker implementation for Edge.
// It is healthy if all of the healthChecks passed.
func (e *Edge) CheckHealth() Status {
	return checkAll(e.healthChecks)
}

// CheckReadiness is a ReadinessChecker implementation for Edge.
// It is ready if all of the readyChecks passed.
func (e *Edge) CheckReadiness() Status {
	return checkAll(e.readinessChecks)
}
