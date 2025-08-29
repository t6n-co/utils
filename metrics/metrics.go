package metrics

import (
	"sync"

	"github.com/t6n-co/utils/internal"
)

var (
	getOnce   sync.Once
	getClient ClientInterface
)

func GetClient() ClientInterface {
	getOnce.Do(func() {
		getClient = GetOtelClient()
		registry := internal.GetCallbackRegistry()
		registry.Register("logs.error", getClient.incrError)
		registry.Register("logs.warn", getClient.incrWarn)

	})
	return getClient
}
