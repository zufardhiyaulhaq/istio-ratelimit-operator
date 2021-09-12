package types

type ClusterPatchValues struct {
	Name                 string               `yaml:"name,omitempty"`
	Type                 string               `yaml:"type,omitempty"`
	ConnectTimeout       string               `yaml:"connect_timeout,omitempty"`
	LbPolicy             string               `yaml:"lb_policy,omitempty"`
	HTTP2ProtocolOptions HTTP2ProtocolOptions `yaml:"http2_protocol_options"`
	LoadAssignment       LoadAssignment       `yaml:"load_assignment,omitempty"`
}

type HTTP2ProtocolOptions struct {
}

type LoadAssignment struct {
	ClusterName string      `yaml:"cluster_name,omitempty"`
	Endpoints   []Endpoints `yaml:"endpoints,omitempty"`
}

type Endpoints struct {
	LbEndpoints []LbEndpoints `yaml:"lb_endpoints,omitempty"`
}

type LbEndpoints struct {
	Endpoint Endpoint `yaml:"endpoint,omitempty"`
}

type Endpoint struct {
	Address Address `yaml:"address,omitempty"`
}

type Address struct {
	SocketAddress SocketAddress `yaml:"socket_address,omitempty"`
}

type SocketAddress struct {
	Address   string `yaml:"address,omitempty"`
	PortValue int    `yaml:"port_value,omitempty"`
}
