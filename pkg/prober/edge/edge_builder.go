package edge

// Builder is a builder for application edge.
type Builder struct {
	healthChecks    []CheckFunc
	readinessChecks []CheckFunc
}

// WithHealthChecker adds a health check to application edge.
func (b *Builder) WithHealthChecker(checker HealthChecker) *Builder {
	b.healthChecks = append(
		b.healthChecks,
		func() Status { return checker.CheckHealth() },
	)
	return b
}

// WithReadinessChecker adds a readiness check to application edge.
func (b *Builder) WithReadinessChecker(checker ReadinessChecker) *Builder {
	b.readinessChecks = append(
		b.readinessChecks,
		func() Status { return checker.CheckReadiness() },
	)
	return b
}

// WithK8SChecker adds a k8s checker to application edge.
func (b *Builder) WithK8SChecker(checker K8SChecker) *Builder {
	return b.WithHealthChecker(checker).WithReadinessChecker(checker)
}

// WithHealthCheckFunc adds a health CheckFunc to application edge.
func (b *Builder) WithHealthCheckFunc(fnc CheckFunc) *Builder {
	b.healthChecks = append(b.healthChecks, fnc)
	return b
}

// WithReadinessCheckFunc adds a readiness CheckFunc to application edge.
func (b *Builder) WithReadinessCheckFunc(fnc CheckFunc) *Builder {
	b.readinessChecks = append(b.readinessChecks, fnc)
	return b
}

// Build builds an application edge.
func (b *Builder) Build() Edge {
	return Edge{
		healthChecks:    b.healthChecks,
		readinessChecks: b.readinessChecks,
	}
}
