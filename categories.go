package sleego

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
