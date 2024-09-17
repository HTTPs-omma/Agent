package Extension

import (
	"github.com/elastic/go-sysinfo"
	"github.com/elastic/go-sysinfo/types"
	"strings"
	"time"
)

type sysutils struct {
	host       types.Host
	hostInfo   types.HostInfo
	cputime    types.CPUTimes
	osInfo     *types.OSInfo
	memoryInfo *types.HostMemoryInfo
}

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

func (sys *sysutils) GetHostName() string {
	return sys.hostInfo.Hostname
}

func (sys *sysutils) GetOsName() string {
	return sys.osInfo.Name
}

func (sys *sysutils) GetOsVersion() string {
	return sys.osInfo.Version
}

/*
*
24.08.07
테스트 진행 안됨.
엄밀한 테스트가 필요함!!
sys.osInfo.Platform 에서 나올 수 있는 모든 경우의 수를 봐야함.
*/
func (sys *sysutils) GetPlatform() OSPLATFORM {
	if strings.Contains(sys.osInfo.Platform, "Windows10") {
		return WINDOWS
	}
	if strings.Contains(sys.osInfo.Platform, "Ubuntu") {
		return UBUNTU
	}
	if strings.Contains(sys.osInfo.Platform, "Centos") {
		return CENTOS
	}
	if strings.Contains(sys.osInfo.Platform, "MAC") {
		return MACOS
	}

	return UNKNOWN
}

func (sys *sysutils) GetFamily() string {
	return sys.osInfo.Family
}

func (sys *sysutils) GetMemoryTotal() uint64 {
	return sys.memoryInfo.Total
}

func (sys *sysutils) GetMemoryUsed() uint64 {
	return sys.memoryInfo.Used
}

func (sys *sysutils) GetMemoryFree() uint64 {
	return sys.memoryInfo.Free
}

func (sys *sysutils) GetArchitecture() string {
	return sys.hostInfo.Architecture
}
func (sys *sysutils) GetNativeArchitecture() string {
	return sys.hostInfo.NativeArchitecture
}
func (sys *sysutils) GetKernelVersion() string {
	return sys.hostInfo.KernelVersion
}
func (sys *sysutils) GetUniqueID() string {
	uuid := strings.Replace(sys.hostInfo.UniqueID, "-", "", -1)
	return uuid
}

func (sys *sysutils) GetBootTime() time.Time {
	return sys.hostInfo.BootTime
}
func (sys *sysutils) GetIPs() []string {
	return sys.hostInfo.IPs
}
func (sys *sysutils) GetMACs() []string {
	return sys.hostInfo.MACs
}
func (sys *sysutils) GetContainerized() *bool {
	return sys.hostInfo.Containerized
}
func (sys *sysutils) GetTimezoneOffsetSec() int {
	return sys.hostInfo.TimezoneOffsetSec
}

func NewSysutils() (*sysutils, error) {
	host, err := sysinfo.Host()
	if err != nil {
		return &sysutils{}, err
	}
	hostInfo := host.Info()
	osInfo := hostInfo.OS

	cputime, err := host.CPUTime()
	if err != nil {
		return &sysutils{}, err
	}

	memoryInfo, err := host.Memory()
	if err != nil {
		return &sysutils{}, err
	}

	return &sysutils{host: host, hostInfo: hostInfo, osInfo: osInfo,
		cputime: cputime, memoryInfo: memoryInfo}, nil
}
