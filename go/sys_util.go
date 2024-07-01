package main

import (
    "encoding/json"
    "fmt"
    "github.com/shirou/gopsutil/v3/cpu"
    "github.com/shirou/gopsutil/v3/disk"
    "github.com/shirou/gopsutil/v3/mem"
    "github.com/shirou/gopsutil/v3/net"
    "github.com/shirou/gopsutil/v3/host"
    "os"
    "time"
)

type SystemInfo struct {
    CPUCount        int                      `json:"cpu_count"`
    CPUCountLogical int                      `json:"cpu_count_logical"`
    CPUPercent      float64                  `json:"cpu_percent"`
    CPUTimes        []cpu.TimesStat          `json:"cpu_times"`
    Memory          *mem.VirtualMemoryStat   `json:"memory"`
    Swap            *mem.SwapMemoryStat      `json:"swap"`
    DiskPartitions  []disk.PartitionStat     `json:"disk_partitions"`
    DiskUsage       map[string]*disk.UsageStat `json:"disk_usage"`
    NetIOCounters   []net.IOCountersStat     `json:"net_io_counters"`
    NetIfAddrs      map[string][]net.InterfaceAddr `json:"net_if_addrs"`
    NetIfStats      []net.InterfaceStat      `json:"net_if_stats"`
    HostInfo        *host.InfoStat           `json:"host_info"`
}

func getSystemInfo() (*SystemInfo, error) {
    info := &SystemInfo{}

    // CPU information
    cpuCount, err := cpu.Counts(false)
    if err != nil {
        return nil, err
    }
    info.CPUCount = cpuCount

    cpuCountLogical, err := cpu.Counts(true)
    if err != nil {
        return nil, err
    }
    info.CPUCountLogical = cpuCountLogical

    cpuPercent, err := cpu.Percent(time.Second, false)
    if err != nil {
        return nil, err
    }
    info.CPUPercent = cpuPercent[0]

    cpuTimes, err := cpu.Times(false)
    if err != nil {
        return nil, err
    }
    info.CPUTimes = cpuTimes

    // Memory information
    memory, err := mem.VirtualMemory()
    if err != nil {
        return nil, err
    }
    info.Memory = memory

    // Swap memory information
    swap, err := mem.SwapMemory()
    if err != nil {
        return nil, err
    }
    info.Swap = swap

    // Disk information
    partitions, err := disk.Partitions(true)
    if err != nil {
        return nil, err
    }
    info.DiskPartitions = partitions

    diskUsage := make(map[string]*disk.UsageStat)
    for _, partition := range partitions {
        usage, err := disk.Usage(partition.Mountpoint)
        if err != nil {
            return nil, err
        }
        diskUsage[partition.Mountpoint] = usage
    }
    info.DiskUsage = diskUsage

    // Network information
    netIOCounters, err := net.IOCounters(true)
    if err != nil {
        return nil, err
    }
    info.NetIOCounters = netIOCounters

    netIfAddrs, err := net.Interfaces()
    if err != nil {
        return nil, err
    }
    netIfAddrsMap := make(map[string][]net.InterfaceAddr)
    for _, iface := range netIfAddrs {
        addrs, err := iface.Addrs()
        if err != nil {
            return nil, err
        }
        netIfAddrsMap[iface.Name] = addrs
    }
    info.NetIfAddrs = netIfAddrsMap

    netIfStats, err := net.Interfaces()
    if err != nil {
        return nil, err
    }
    info.NetIfStats = netIfStats

    // Host information
    hostInfo, err := host.Info()
    if err != nil {
        return nil, err
    }
    info.HostInfo = hostInfo

    return info, nil
}

func writeSystemInfoToFile(info *SystemInfo, filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "    ")

    return encoder.Encode(info)
}

func main() {
    systemInfo, err := getSystemInfo()
    if err != nil {
        fmt.Printf("Error getting system info: %v\n", err)
        return
    }

    err = writeSystemInfoToFile(systemInfo, "output.txt")
    if err != nil {
        fmt.Printf("Error writing system info to file: %v\n", err)
        return
    }

    fmt.Println("System information written to output.txt")
}

