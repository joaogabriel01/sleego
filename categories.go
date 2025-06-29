package sleego

// CategoryOperator defines the behavior for managing categories of applications
type CategoryOperator interface {
	GetCategoriesOf(process string) []string
	SetProcessByCategories(categoriesByProcess map[string][]string)
}

var (
	processByCategories map[string][]string
	categoryOperator    *CategoryOperatorImpl
)

type CategoryOperatorImpl struct {
}

func (c *CategoryOperatorImpl) GetCategoriesOf(process string) []string {
	return processByCategories[process]
}

func (c *CategoryOperatorImpl) SetProcessByCategories(categoryByProcesses map[string][]string) {
	if processByCategories == nil {
		processByCategories = make(map[string][]string)
	}
	for categories, processes := range categoryByProcesses {
		for _, process := range processes {
			processByCategories[process] = append(processByCategories[process], categories)
		}
	}

}

func GetCategoryOperator() *CategoryOperatorImpl {
	if categoryOperator == nil {
		categoryOperator = &CategoryOperatorImpl{}
	}
	return categoryOperator
}

var _ CategoryOperator = &CategoryOperatorImpl{}
