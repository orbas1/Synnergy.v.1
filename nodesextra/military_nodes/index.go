package militarynodes

// BaseNode defines minimal functionality expected from a network node.
type BaseNode interface {
	// GetID returns the node identifier.
	GetID() string
}

// WarfareNode extends a base node with military specific operations.
type WarfareNode interface {
	BaseNode
	// SecureCommand executes a privileged command after verifying
	// appropriate authorization. Implementations should ensure commands
	// are transmitted using secure channels.
	SecureCommand(cmd string) error
	// TrackLogistics records movement or status changes for a military asset.
	// assetID uniquely identifies the asset, while location and status capture
	// current logistics information.
	TrackLogistics(assetID, location, status string)
	// ShareTactical distributes tactical information to allied nodes or systems.
	ShareTactical(info string)
}
