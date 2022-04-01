package channel

import (
	"atlas-cts/rest/requests"
)

const (
	ServicePrefix string = "/ms/wrg/"
	Service              = requests.BaseRequest + ServicePrefix
	Resource             = Service + "channelServers/"
	ByWorld              = Resource + "?world=%d"
)

func requestChannels() requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](Resource)
}
