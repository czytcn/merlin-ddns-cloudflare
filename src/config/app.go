package config

type AppSetting struct {
	UpdateIpv4DNSRecord         bool `toml:"update_ipv4_dns_record"`
	UpdateIpv6DNSRecord         bool `toml:"update_ipv6_dns_record"`
	EnableCloudflareProxied     bool `toml:"enable_cloudflare_proxied"`
	TTL                         int  `toml:"ttl"`
	SendDdnsCustomUpdatedNotify bool `toml:"send_ddns_custom_updated_notify"`
}
