module github.com/trustbloc/edge-agent/cmd/user-agent-js-worker

go 1.13

require (
	github.com/google/uuid v1.1.1
	github.com/hyperledger/aries-framework-go v0.1.3-0.20200311212058-6f509cae073a
	github.com/mitchellh/mapstructure v1.1.2
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.5.1
	github.com/trustbloc/edge-agent v0.0.0
)

replace github.com/trustbloc/edge-agent => ../../
