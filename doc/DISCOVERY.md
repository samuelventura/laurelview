# Discovery

## Request

```json
{"name":"nerves","action":"id"}
```

## Response

```json
{
    "data": {
        "hostname":"nerves-4863",
        "ifname":"eth0",
        "macaddr":"E45F01574889",
        "name":"nerves",
        "version":"0.1.2"
    },
    "action":"id",
    "name":"nerves"
}
```

## Sample

```bash
samuel@p3420:~/github/laurelview$ ./discover.sh 
+ MOD=github.com/YeicoLabs/laurelview
+ go install github.com/YeicoLabs/laurelview/cmd/lvndc
+ /home/samuel/go/bin/lvndc
2021/12/03 02:10:43 LocalAddr 0.0.0.0:56876
2021/12/03 02:10:43 > {"name":"nerves","action":"id"}
2021/12/03 02:10:43 < {"data":{"hostname":"nerves-4863","ifname":"eth0","macaddr":"E45F01574889","name":"nerves","version":"0.1.2"},"action":"id","name":"nerves"}
2021/12/03 02:10:43 &{nerves id {nerves-4863 eth0 E45F01574889 nerves 0.1.2}}
2021/12/03 02:10:44 read udp4 0.0.0.0:56876: i/o timeout
```