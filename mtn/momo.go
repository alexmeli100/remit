package mtn

const (
	BaseURL   = "https://sandbox.momodeveloper.mtn.com"
	TargetEnv = "sandbox"
)

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
	if gc != nil {
		if gc.BaseUrl == "" {
			gc.BaseUrl = BaseURL
		}
		if gc.TargetEnv == "" {
			gc.TargetEnv = TargetEnv
		}
	} else {
		gc = &GlobalConfig{BaseUrl: BaseURL, TargetEnv: TargetEnv}
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
