package types

//Message JSON message formwat to return with host list.
type Message struct {
	Status  string
	Message []string
}

//Config file specs for configuration YAML
type Config struct {
	Username  string
	Password  string
	DnsServer string
	Domain    string
}
