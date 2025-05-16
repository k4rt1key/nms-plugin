package types

// RequestType defines the type of operation being requested
type RequestType string

const (
	DiscoveryType RequestType = "discovery"
	PollingType   RequestType = "polling"
)

// Credential represents authentication information
type Credential struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// DiscoveryRequest represents the input for discovery operation
type DiscoveryRequest struct {
	Type        RequestType  `json:"type"`
	ID          int          `json:"id"`
	IPs         []string     `json:"ips"`
	Credentials []Credential `json:"credentials"`
	Port        int          `json:"port"`
}

// DiscoveryResult represents the result for a single IP and credential pair
type DiscoveryResult struct {
	Success    bool       `json:"success"`
	IP         string     `json:"ip"`
	Credential Credential `json:"credential"`
	Port       int        `json:"port"`
	Message    string     `json:"message"`
}

// DiscoveryResponse represents the output for discovery operation
type DiscoveryResponse struct {
	Type   RequestType       `json:"type"`
	ID     int               `json:"id"`
	Result []DiscoveryResult `json:"results"`
}

// MetricGroup represents a group of metrics to be polled
type MetricGroup struct {
	MonitorID  int        `json:"monitor_id"`
	Name       string     `json:"name"`
	IP         string     `json:"ip"`
	Port       int        `json:"port"`
	Credential Credential `json:"credential"`
}

// PollingRequest represents the input for polling operation
type PollingRequest struct {
	Type         RequestType   `json:"type"`
	MetricGroups []MetricGroup `json:"metric_groups"`
}

// PollingResult represents the result for a single metric group
type PollingResult struct {
	Success            bool   `json:"success"`
	ProvisionProfileID int    `json:"monitor_id"`
	Name               string `json:"name"`
	Data               string `json:"data"`
	Message            string `json:"message"`
	Time               string `json:"time"`
}

// PollingResponse represents the output for polling operation
type PollingResponse struct {
	Type         RequestType     `json:"type"`
	MetricGroups []PollingResult `json:"metric_groups"`
}
