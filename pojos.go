package main

func AutoMigrate() {
	DB.AutoMigrate(&PreviewCountry{}, &WaitingUser{}, &User{}, &Country{}, &DiplomaticStance{})
}

type PreviewCountry struct {
	Name           string `json:"name" gorm:"primaryKey"`
	Population     int    `json:"population"`
	Ressources     int    `json:"ressources"`
	Budget         int    `json:"budget"`
	GovernmentType string `json:"government_type"`
	Economy        int    `json:"economy"`
	Happiness      int    `json:"happiness"`
	Visible        bool   `json:"visible"`
}

type WaitingUser struct {
	RealName string `json:"realname"`
	Username string `json:"username" gorm:"primaryKey"`
	Password string `json:"password"`
	Country  string `json:"country"`
}

type User struct {
	UserName     string `json:"username" gorm:"primaryKey"`
	PasswordHash string `json:"password"`
	Country      string `json:"country"`
}

type DiplomaticStance struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	CountryName   string `json:"-"` // Foreign key referencing Country.Name
	TargetCountry string `json:"target_country"`
	StanceType    string `json:"stance_type"` // e.g. "Ally", "Enemy", "Neutral"
}

type Country struct {
	Happiness              int                `json:"happiness"`                                       // Range: 0-100
	Population             int                `json:"population"`                                      // Range: >= 0
	Budget                 int                `json:"budget"`                                          // Range: Any integer (can be negative if bankrupt/in debt)
	Economy                int                `json:"economy"`                                         // Range: 0-100
	Education              int                `json:"education"`                                       // Range: 0-100
	GovernmentType         string             `json:"government_type"`                                 // e.g. "Democracy", "Dictatorship"
	CorruptionLevel        int                `json:"corruption_level"`                                // Range: 0-100
	ActiveTraits           []string           `json:"active_traits" gorm:"serializer:json"`            // Array of current trait IDs or names
	InfrastructureLevel    int                `json:"infrastructure_level"`                            // Range: 0-100
	TechLevel              int                `json:"tech_level"`                                      // Range: 0-100
	GroundMilitaryStrength int                `json:"ground_military_strength"`                        // Range: >= 0
	AirforceStrength       int                `json:"airforce_strength"`                               // Range: >= 0
	NavyStrength           int                `json:"navy_strength"`                                   // Range: >= 0
	SpyPower               int                `json:"spy_power"`                                       // Range: 0-100
	DiplomaticStances      []DiplomaticStance `json:"diplomatic_stance" gorm:"foreignKey:CountryName"` // Relations to other countries
	Secrets                []string           `json:"secrets" gorm:"serializer:json"`                  // Array of strings representing hidden information/secrets
	InternalSecurity       int                `json:"internal_security"`                               // Range: 0-100
	Status                 []string           `json:"status" gorm:"serializer:json"`                   // e.g. "Bankrott", "erobert", "Stabil"
	Name                   string             `json:"name" gorm:"primaryKey"`
}
