package docker_compose

type Healthcheck struct{
	Test []string
	Interval string
	Timeout string
	Retries string
	Disable bool
}
