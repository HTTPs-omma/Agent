package Extension

import (
	"github.com/elastic/go-sysinfo"
	"github.com/elastic/go-sysinfo/types"
	"strings"
	"time"
)

type Sysutils struct {
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

func (sys *Sysutils) GetHostName() string {
	return sys.hostInfo.Hostname
}

func (sys *Sysutils) GetOsName() string {
	return sys.osInfo.Name
}

func (sys *Sysutils) GetOsVersion() string {
	return sys.osInfo.Version
}

/*
*
24.08.07
테스트 진행 안됨.
엄밀한 테스트가 필요함!!
sys.osInfo.Platform 에서 나올 수 있는 모든 경우의 수를 봐야함.
*/
func (sys *Sysutils) GetPlatform() OSPLATFORM {
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

func (sys *Sysutils) GetFamily() string {
	return sys.osInfo.Family
}

func (sys *Sysutils) GetMemoryTotal() uint64 {
	return sys.memoryInfo.Total
}

func (sys *Sysutils) GetMemoryUsed() uint64 {
	return sys.memoryInfo.Used
}

func (sys *Sysutils) GetMemoryFree() uint64 {
	return sys.memoryInfo.Free
}

func (sys *Sysutils) GetArchitecture() string {
	return sys.hostInfo.Architecture
}
func (sys *Sysutils) GetNativeArchitecture() string {
	return sys.hostInfo.NativeArchitecture
}
func (sys *Sysutils) GetKernelVersion() string {
	return sys.hostInfo.KernelVersion
}
func (sys *Sysutils) GetUniqueID() string {
	uuid := strings.Replace(sys.hostInfo.UniqueID, "-", "", -1)
	return uuid
}

func (sys *Sysutils) GetBootTime() time.Time {
	return sys.hostInfo.BootTime.UTC()
}

func (sys *Sysutils) GetIPs() []string {
	return sys.hostInfo.IPs
}
func (sys *Sysutils) GetMACs() []string {
	return sys.hostInfo.MACs
}
func (sys *Sysutils) GetContainerized() *bool {
	return sys.hostInfo.Containerized
}
func (sys *Sysutils) GetTimezoneOffsetSec() int {
	return sys.hostInfo.TimezoneOffsetSec
}

func NewSysutils() (*Sysutils, error) {
	host, err := sysinfo.Host()
	if err != nil {
		return &Sysutils{}, err
	}
	hostInfo := host.Info()
	osInfo := hostInfo.OS

	cputime, err := host.CPUTime()
	if err != nil {
		return &Sysutils{}, err
	}

	memoryInfo, err := host.Memory()
	if err != nil {
		return &Sysutils{}, err
	}

	return &Sysutils{host: host, hostInfo: hostInfo, osInfo: osInfo,
		cputime: cputime, memoryInfo: memoryInfo}, nil
}
