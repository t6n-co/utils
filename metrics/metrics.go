package metrics

func GetClient() ClientInterface {
	return GetOtelClient()
}
