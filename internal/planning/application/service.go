// Package application implements dependency injection container for planning service use cases.
package application

// PlanningApplication is a dependency injection container that aggregates all use cases
// for the planning domain to be injected from main.go.
type PlanningApplication struct {
	GetAllPlans GetAllPlans
	GetPlan     GetPlan
}
