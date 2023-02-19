package request

type ListDeploymentsOptions struct {
	ListOptions
	ChannelUpdate       bool
	DanglingDeployments bool
	OperatorUpdate      bool
}
