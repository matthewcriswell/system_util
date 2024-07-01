import psutil
import json


def get_system_info():
    info = {}

    # CPU information
    info["cpu_count"] = psutil.cpu_count(logical=False)
    info["cpu_count_logical"] = psutil.cpu_count(logical=True)
    info["cpu_percent"] = psutil.cpu_percent(interval=1)
    info["cpu_times"] = psutil.cpu_times()._asdict()

    # Memory information
    memory = psutil.virtual_memory()
    info["memory_total"] = memory.total
    info["memory_available"] = memory.available
    info["memory_used"] = memory.used
    info["memory_percent"] = memory.percent

    # Swap memory information
    swap = psutil.swap_memory()
    info["swap_total"] = swap.total
    info["swap_used"] = swap.used
    info["swap_free"] = swap.free
    info["swap_percent"] = swap.percent

    # Disk information
    info["disk_partitions"] = [part._asdict() for part in psutil.disk_partitions()]
    info["disk_usage"] = {
        part.mountpoint: psutil.disk_usage(part.mountpoint)._asdict()
        for part in psutil.disk_partitions()
    }

    # Network information
    info["net_io_counters"] = psutil.net_io_counters()._asdict()
    info["net_if_addrs"] = {
        iface: [addr._asdict() for addr in addrs]
        for iface, addrs in psutil.net_if_addrs().items()
    }
    info["net_if_stats"] = {
        iface: stats._asdict() for iface, stats in psutil.net_if_stats().items()
    }

    return info


def write_system_info_to_file(info, filename):
    with open(filename, "w") as file:
        file.write(json.dumps(info, indent=4))


if __name__ == "__main__":
    system_info = get_system_info()
    write_system_info_to_file(system_info, "output.txt")
