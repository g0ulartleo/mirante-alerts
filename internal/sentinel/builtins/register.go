package builtins

import "github.com/g0ulartleo/mirante-alerts/internal/sentinel"

func Register(f *sentinel.SentinelFactory) {
	f.Register("endpoint-checker", NewEndpointCheckerSentinel)
	f.Register("mysql-query-checker", NewMySQLQueryCheckerSentinel)
}
