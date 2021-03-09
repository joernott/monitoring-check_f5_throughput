# check_f5_throughput - Monitoring F5 Loadbalancer Throughput

This check uses SNMP to determine the F5 loadbalancer throughput. It uses the
method described in [K50309321: Viewing BIG-IP system throughput statistics](https://support.f5.com/csp/article/K50309321) to retrieve the
information from the load balancer, stores the information in a file and upon the 2nd call, uses the actual and stored data to calculate
the throughput of the load balancer.

## Usage
```
f5_throughput [flags]

Flags:
  -C, --community string   SNMP community (default public (default "public")
  -f, --config string      config file (default is /etc/icinga2/f5_throughput.yaml) (default "/etc/icinga2/f5_throughput.yaml")
  -c, --critical string    critical range
  -s, --file string        statistics file (default /var/lib/icinga2/f5_throughput_stats.json)
  -h, --help               help for f5_throughput
  -H, --host string        host/ip address of the laod balancer (default 127.0.0.1) (default "127.0.0.1")
  -P, --port int           SNMP port (default 161) (default 161)
  -w, --warning string     warning range
```
