package koronapay

type Option func(*Client)

func WithUpstream(upstream string) Option {
	return func(client *Client) {
		client.httpc.SetBaseURL(upstream + "/transfers/online/api")
	}
}

func WithVerbose(verbose bool) Option {
	return func(client *Client) {
		client.httpc.SetDebug(verbose)
	}
}
