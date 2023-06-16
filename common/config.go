package common

type Config struct {
	Demo []Demo `mapstructure:"demo"`
}

type Demo struct {
	Controller []Controller `mapstructure:"controller"`
	CoPilot    []CoPilot    `mapstructure:"copilot"`
	Azure      []Azure      `mapstructure:"azure"`
}

type Controller struct {
	Cloud      string `mapstructure:"cloud"`
	Region     string `mapstructure:"region"`
	Name       string `mapstructure:"name"`
	Version    string `mapstructure:"version"`
	Email      string `mapstructure:"email"`
	Password   string `mapstructure:"password"`
	CustomerId string `mapstructure:"customerId"`
	Keypair    string `mapstructure:"keypair"`
	AccID      string `mapstructure:"accId"`
	VpcCidr    string `mapstructure:"vpcCidr"`
}

type CoPilot struct {
	Cluster bool `mapstructure:"cluster"`
}

type Azure struct {
	SubscriptionId          string `mapstructure:"subscriptionid"`
	ApplicationEndpoint     string `mapstructure:"directoryid"`
	ApplicationClientId     string `mapstructure:"applicationid"`
	ApplicationClientSecret string `mapstructure:"applicationclientsecret"`
}
