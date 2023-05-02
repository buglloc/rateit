package contactsys

type Option func(session *Session)

func WithUpstream(upstream string) Option {
	return func(s *Session) {
		s.httpc.SetBaseURL(upstream + "/api/contact/v2")
		s.httpc.SetHeader("Referer", upstream)
	}
}

func WithPartnerID(id string) Option {
	return func(s *Session) {
		s.partnerID = id
	}
}

func WithVerbose(verbose bool) Option {
	return func(s *Session) {
		s.httpc.SetDebug(verbose)
	}
}
