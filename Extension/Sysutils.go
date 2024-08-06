package Extension

import (
	"github.com/elastic/go-sysinfo"
	"github.com/elastic/go-sysinfo/types"
	"time"
)

type sysutils struct {
	host       types.Host
	hostInfo   types.HostInfo
	cputime    types.CPUTimes
	osInfo     *types.OSInfo
	memoryInfo *types.HostMemoryInfo
}

func (sys *sysutils) GetHostName() string {
	return sys.hostInfo.Hostname
}

func (sys *sysutils) GetOsName() string {
	return sys.osInfo.Name
}

func (sys *sysutils) GetOsVersion() string {
	return sys.osInfo.Version
}

func (sys *sysutils) GetPlatform() string {
	return sys.osInfo.Platform
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
	return sys.hostInfo.UniqueID
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
