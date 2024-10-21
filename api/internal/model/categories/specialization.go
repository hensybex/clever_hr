// model/categories/specialization.go

package category_model

type Specialization struct {
	ID         int      `json:"id"`
	TitleRu    string   `json:"title_ru"`
	TitleEn    string   `json:"title_en"`
	JobGroupID int      `json:"group_id"`
	JobGroup   JobGroup `gorm:"foreignKey:JobGroupID"`
}
