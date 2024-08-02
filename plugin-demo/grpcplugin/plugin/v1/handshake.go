package v1

import "github.com/hashicorp/go-plugin"

func Handshake() plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		// The ProtocolVersion is the version that must match between Velero framework
		// and Velero client plugins. This should be bumped whenever a change happens in
		// one or the other that makes it so that they can't safely communicate.
		ProtocolVersion: 2,

		MagicCookieKey:   "VELERO_PLUGIN",
		MagicCookieValue: "hello",
	}
}
