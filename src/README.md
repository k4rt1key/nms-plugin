1. Polling 

```json
{
  "type": "polling",
  "metric_groups": [
    {
      "monitor_id": 1,
      "name": "CPUINFO",
      "ip": "172.16.8.104",
      "port": 5985,
      "credential": {
        "id": 1,
        "name": "one",
        "protocol": "winrm",
        "credential": {
          "password": "Mind@123",
          "username": "Administrator"
        }
      }
    }
  ]
}
```

2. Discover 

```json
{
  "type": "discovery",
  "id": 1,
  "ips": [
    "172.16.8.113",
    "172.16.8.57",
    "172.16.8.128",
    "172.16.8.86",
    "172.16.8.104",
    "172.16.8.98"
  ],
  "port": 5985,
  "credentials": [
    {
      "id": 1,
      "name": "one",
      "protocol": "winrm",
      "credential": {
        "password": "Mind@123",
        "username": "Administrator"
      }
    },
    {
      "id": 2,
      "name": "two",
      "protocol": "winrm",
      "credential": {
        "password": "Mindarray@8",
        "username": "Administrator"
      }
    },
    {
      "id": 3,
      "name": "three",
      "protocol": "winrm",
      "credential": {
        "password": "Mindarray@8",
        "username": "admin"
      }
    },
    {
      "id": 4,
      "name": "four",
      "protocol": "winrm",
      "credential": {
        "password": "Mind@123",
        "username": "admin"
      }
    }
  ]
}
```