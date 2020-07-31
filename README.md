# merlin-ddns-cloudflare
asuswrt merlin cloudflare ddns using api

## compile

you can compile this tool in any PC or Arm devices with go installed. compile command is quite simple

```
GOOS=linux GOARCH=arm go build .
```

I use NetGear R8000 it is a arm v7 device,so I use `arm` in the command above.if you are using arm 64 bit device,change the `GOARCH` to `arm64`  .

for Non-arm device,you can use the value below

```
➜  ~ go tool dist list
aix/ppc64
android/386
android/amd64
android/arm
android/arm64
darwin/386
darwin/amd64
darwin/arm
darwin/arm64
dragonfly/amd64
freebsd/386
freebsd/amd64
freebsd/arm
freebsd/arm64
illumos/amd64
js/wasm
linux/386
linux/amd64
linux/arm
linux/arm64
linux/mips
linux/mips64
linux/mips64le
linux/mipsle
linux/ppc64
linux/ppc64le
linux/riscv64
linux/s390x
netbsd/386
netbsd/amd64
netbsd/arm
netbsd/arm64
openbsd/386
openbsd/amd64
openbsd/arm
openbsd/arm64
plan9/386
plan9/amd64
plan9/arm
solaris/amd64
windows/386
windows/amd64
windows/arm

```



### config

make a config named `config.toml` and contents like below

```toml
[api]
api_token = "xxxxx"
email = "yyyy"
domain = "zzzzz.com"
sub_domain = "test.zzzzz.com"

[app]
get_ip_from_url= "https://api-ipv4.ip.sb/ip"
```

you should add an `A` record in your domain Control,first

## deploy

1. copy compiled binary into `/jffs/scripts/` and  rename it to `ddns-start`  run

   ```
   chmod +x ddns-start
   ```

   to make it executable

2. login in your router web page

   Go to `Advanced Settings` > `WAN` > `DDNS` Set `Server` to `Custom`  and Click the `Apply` button