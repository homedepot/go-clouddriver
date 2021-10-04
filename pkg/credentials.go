package clouddriver

type Credentials struct {
	AccountType                 string        `json:"accountType"`
	CacheThreads                int           `json:"cacheThreads"`
	ChallengeDestructiveActions bool          `json:"challengeDestructiveActions"`
	CloudProvider               string        `json:"cloudProvider"`
	DockerRegistries            []interface{} `json:"dockerRegistries"`
	Enabled                     bool          `json:"enabled"`
	Environment                 string        `json:"environment"`
	Name                        string        `json:"name"`
	Namespaces                  []string      `json:"namespaces"`
	Permissions                 struct {
		READ  []string `json:"READ"`
		WRITE []string `json:"WRITE"`
	} `json:"permissions"`
	PrimaryAccount          bool              `json:"primaryAccount"`
	ProviderVersion         string            `json:"providerVersion"`
	RequiredGroupMembership []interface{}     `json:"requiredGroupMembership"`
	Skin                    string            `json:"skin"`
	SpinnakerKindMap        map[string]string `json:"spinnakerKindMap"`
	Type                    string            `json:"type"`
}
