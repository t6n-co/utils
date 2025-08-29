package metrics

import "github.com/t6n-co/utils/internal"

func GetClient() ClientInterface {
	client := GetOtelClient()

	registry := internal.GetCallbackRegistry()
	registry.Register("logs.error", client.incrError)
	registry.Register("logs.warn", client.incrWarn)

	return client
}
