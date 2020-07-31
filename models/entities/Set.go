package entities

type Set struct {
	Code                 string `gorm:"primary_key"`
	Name                 string
	Type                 string
	ReleaseDate          string
	BaseSetSize          int
	Cards                []Card      `gorm:"foreignkey:SetID;PRELOAD:true"`
	BoosterConfiguration BoosterRule `gorm:"foreignkey:ID;PRELOAD:true"`
}
