package resource

import (
	"time"

	"k8s.io/apimachinery/pkg/util/duration"
)

func Age(since time.Time) string {
	age := duration.HumanDuration(time.Since(since))
	if since.IsZero() {
		age = ""
	}

	return age
}
