package utils

// Binary Scan Graph Node
type BinaryGraphNode struct {
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
	// List of license names
	Licenses []string `json:"licenses,omitempty"`
	// Component properties
	Properties map[string]string `json:"properties,omitempty"`
	// List of subcomponents.
	Nodes []*BinaryGraphNode `json:"nodes,omitempty"`
	// Other component IDs field is populated by the Xray indexer to get a better accuracy in '.deb' files.
	OtherComponentIds []OtherComponentIds `json:"other_component_ids,omitempty"`
}

type OtherComponentIds struct {
	Id     string `json:"component_id,omitempty"`
	Origin int    `json:"origin,omitempty"`
}

// Audit Graph Node
type GraphNode struct {
	// Node parent (for internal use)
	Parent *GraphNode `json:"-"`
	// The "classifier" attribute in a Maven pom.xml specifies an additional qualifier for a dependency
	Classifier *string `json:"-"`
	// Node file types (tar, jar, zip, pom)
	Types *[]string `json:"-"`
	Id    string    `json:"component_id,omitempty"`
	// List of subcomponents.
	Nodes []*GraphNode `json:"nodes,omitempty"`
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
