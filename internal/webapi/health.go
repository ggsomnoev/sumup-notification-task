package webapi

type HealthChecker interface {
	Ping() error
}