package utils

type GraphNode struct {
	// Component Id in the JFrog standard.
	// For instance, for maven: gav://<groupId>:<artifactId>:<version>
	// For detailed format examples please see:
	// https://www.jfrog.com/confluence/display/JFROG/Xray+REST+API#XrayRESTAPI-ComponentIdentifiers
	Id string `json:"component_id,omitempty"`
	// Sha of the binary representing the component.
	Sha256 string `json:"sha256,omitempty"`
	Sha1   string `json:"sha1,omitempty"`
	// For root file shall be the file name.
	// For internal components shall be the internal path. (Relevant only for binary scan).
	Path string `json:"path,omitempty"`
	// Download url
	DownloadUrl string `json:"-"`
	// List of license names
	Licenses []string `json:"licenses,omitempty"`
	// Component properties
	Properties map[string]string `json:"properties,omitempty"`
	// List of subcomponents.
	Nodes []*GraphNode `json:"nodes,omitempty"`
	// Other component IDs field is populated by the Xray indexer to get a better accuracy in '.deb' files.
	OtherComponentIds []OtherComponentIds `json:"other_component_ids,omitempty"`
	// Node parent (for internal use)
	Parent *GraphNode `json:"-"`
	// Node Can appear in some cases without children. When adding node to flatten graph,
	// we want to process node again if it was processed without children.
	ChildrenExist bool `json:"-"`
}

type OtherComponentIds struct {
	Id     string `json:"component_id,omitempty"`
	Origin int    `json:"origin,omitempty"`
}

func (currNode *GraphNode) NodeHasLoop() bool {
	parent := currNode.Parent
	for parent != nil {
		if currNode.Id == parent.Id {
			return true
		}
		parent = parent.Parent
	}
	return false
}
