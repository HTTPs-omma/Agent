package Type

type OSPLATFORM string
type PROTOCOL string

const (
	// refer : https://httpsomma.atlassian.net/wiki/spaces/HTTPs/pages/13172747
	WINDOWS       OSPLATFORM = "Windows"
	WINDOWSSERVER OSPLATFORM = "Windows Server"
	UBUNTU        OSPLATFORM = "Ubuntu"
	CENTOS        OSPLATFORM = "Centos"
	MACOS         OSPLATFORM = "MAC OS"
	UNKNOWN       OSPLATFORM = "unknown"
)

const (
	TCP  PROTOCOL = "TCP"
	UDP  PROTOCOL = "UDP"
	HTTP PROTOCOL = "HTTP"
)
