package productdto

type AddProduct struct {
	Nameproduct string `json:"nameproduct" form:"nameproduct" gorm:"type:text"`
	Price       int    `json:"price" form:"price" gorm:"type:int"`
	Qty         int    `json:"qty" form:"qty" gorm:"type:int"`
	Image       string `json:"image" form:"image" gorm:"type:text"`
}

type UpdateProduct struct {
	Nameproduct string `json:"nameproduct" form:"nameproduct" gorm:"type:text"`
	Price       int    `json:"price" form:"price" gorm:"type:int"`
	Qty         int    `json:"qty" form:"qty" gorm:"type:int"`
	Image       string `json:"image" form:"image" gorm:"type:text"`
}
