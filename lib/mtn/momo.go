package mtn

type GlobalConfig struct {
	BaseUrl   string
	TargetEnv string
}

type ProductConfig struct {
	PrimaryKey string
	ApiSecret  string
	UserId     string
}

type MomoApp struct {
	gc *GlobalConfig
}

func CreateMomoApp(gc *GlobalConfig) *MomoApp {
	if gc.BaseUrl == "" {
		gc.BaseUrl = BaseURL
	}

	if gc.TargetEnv == "" {
		gc.TargetEnv = TargetEnv
	}

	return &MomoApp{gc}
}

func (m *MomoApp) NewRemittance(pc *ProductConfig) *Remittance {
	conf := &Config{
		targetEnv:  m.gc.TargetEnv,
		baseUrl:    m.gc.BaseUrl,
		primaryKey: pc.PrimaryKey,
		userId:     pc.UserId,
		apiSecret:  pc.ApiSecret,
	}

	return NewRemittance(conf)
}
