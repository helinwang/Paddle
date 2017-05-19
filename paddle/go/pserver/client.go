package pserver

import "net/rpc"

// Selector selects a client from multiple clients for initializing
// parameters.
type Selector interface {
	Select() bool
}

// ServerInfo is the parameter server information.
type ServerInfo struct {
	Addr  string
	Index int
}

// Lister lists the server info of parameter server
type Lister interface {
	// List will block until all parameter servers are present
	List() []ServerInfo
}

// Client is the client to parameter servers.
type Client struct {
	sel    Selector
	client *rpc.Client
}

// TODO(helin): add TCP re-connect logic

// NewClient creates a new client.
//
// NewClient will block until all parameter servers are ready.
func NewClient(lister Lister, sel Selector) (*Client, error) {
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{sel: sel, client: client}, nil
}

// BeginInitParams begins to initialize parameters on parameter
// servers.
//
// BeginInitParams will be called from multiple trainers, only one
// trainer will be selected to initialize the parameters on parameter
// servers. Other trainers will be blocked until the initialization is
// done, and they need to get the initialized parameters from
// parameter servers using GetParams.
func (c *Client) BeginInitParams(pserverConfigProto []byte) (selected bool, err error) {
	selected = c.sel.Select()
	if !selected {
		return
	}

	var dummy int
	err = c.client.Call("Service.BeginInitParams", pserverConfigProto, &dummy)
	if err != nil {
		return
	}

	return true, nil
}

// InitParam initializes the parameter on parameter servers.
func (c *Client) InitParam(paramWithConfigs ParameterWithConfig) error {
	var dummy int
	return c.client.Call("Service.InitParam", paramWithConfigs, &dummy)
}

// FinishInitParams tells parameter servers client has sent all
// parameters to parameter servers as initialization.
func (c *Client) FinishInitParams() error {
	var dummy int
	return c.client.Call("Service.FinishInitParams", dummy, &dummy)
}

// SendGrads sends gradients to parameter servers for updating
// parameters.
func (c *Client) SendGrads(grads []Gradient) error {
	var dummy int
	return c.client.Call("Service.SendGrads", grads, &dummy)
}

// GetParams gets parameters from parameter servers.
func (c *Client) GetParams(names []string) ([]Parameter, error) {
	var dummy int
	var parameters []Parameter
	err := c.client.Call("Service.GetParams", &parameters, &dummy)
	return parameters, err
}

// SaveModel indicates parameters to save the parameter to the given
// path.
func (c *Client) SaveModel(path string) error {
	var dummy int
	return c.client.Call("Service.SaveModel", path, &dummy)
}
