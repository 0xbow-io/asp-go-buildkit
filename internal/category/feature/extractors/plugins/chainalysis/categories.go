package chainalysis

const (
	ChildAbuseMaterial     = 1
	DarknetMarket          = 2
	SanctionedEntity       = 3
	HighRiskExchange       = 4
	StolenFunds            = 6
	MiningPool             = 7
	Other                  = 9
	EthereumContract       = 10
	HostedWallet           = 11
	Ransomware             = 12
	Mixing                 = 13
	Ico                    = 14
	Erc20Token             = 15
	Gambling               = 16
	MerchantServices       = 17
	Scam                   = 18
	P2pExchange            = 19
	None                   = 20
	Exchange               = 21
	Mining                 = 22
	Terrorism              = 23
	ATM                    = 24
	SanctionedJurisdiction = 25
	LendingContract        = 26
	DecentralizedExchange  = 27
	FraudShop              = 28
	IllicitActorOrg        = 29
	InfrastructureService  = 30
	TokenSmartContract     = 31
	SmartContract          = 32
	ProtocolPrivacy        = 33
	SpecialMeasures        = 34
	Malware                = 35
	OnlinePharmacy         = 36
	Bridge                 = 37
	NFTPlatform            = 38
	SeizedFunds            = 39
	UnnamedService         = 41
	StolenBitcoins         = 42
	StolenEther            = 43
	CustomAddress          = 999
)

var BlockedCategories = []int{
	ChildAbuseMaterial,
	DarknetMarket,
	SanctionedEntity,
	StolenFunds,
	Ransomware,
	Mixing,
	Scam,
	Terrorism,
	FraudShop,
	IllicitActorOrg,
	Malware,
	OnlinePharmacy,
	SeizedFunds,
	StolenBitcoins,
	StolenEther,

	// Not sure if we block these as well
	Gambling,
	HighRiskExchange,
	SanctionedJurisdiction,
	ProtocolPrivacy,
	LendingContract,
	P2pExchange,
}

var BlockedCategoriesMap = map[int]bool{
	ChildAbuseMaterial: true,
	DarknetMarket:      true,
	SanctionedEntity:   true,
	StolenFunds:        true,
	Ransomware:         true,
	Mixing:             true,
	Scam:               true,
	Terrorism:          true,
	FraudShop:          true,
	IllicitActorOrg:    true,
	Malware:            true,
	OnlinePharmacy:     true,
	SeizedFunds:        true,
	StolenBitcoins:     true,
	StolenEther:        true,

	// Not sure if we block these as well
	Gambling:               true,
	HighRiskExchange:       true,
	SanctionedJurisdiction: true,
	ProtocolPrivacy:        true,
	LendingContract:        true,
	P2pExchange:            true,
}

var CatToTname = map[int]string{
	1:   "child abuse material",
	2:   "darknet market",
	3:   "sanctioned entity",
	4:   "high risk exchange",
	6:   "stolen funds",
	7:   "mining pool",
	9:   "other",
	10:  "ethereum contract",
	11:  "hosted wallet",
	12:  "ransomware",
	14:  "ico",
	15:  "erc20 token",
	16:  "gambling",
	17:  "merchant services",
	18:  "scam",
	19:  "p2p exchange",
	20:  "none",
	21:  "exchange",
	22:  "mining",
	23:  "terrorist financing",
	24:  "atm",
	25:  "sanctioned jurisdiction",
	26:  "lending",
	27:  "decentralized exchange",
	28:  "fraud shop",
	30:  "infrastructure as a service",
	31:  "token smart contract",
	32:  "smart contract",
	33:  "protocol privacy",
	34:  "special measures",
	35:  "malware",
	36:  "online pharmacy",
	37:  "bridge",
	38:  "nft platform - collection",
	39:  "seized funds",
	41:  "unnamed service",
	42:  "stolen bitcoins",
	43:  "stolen ether",
	999: "custom address",
}
