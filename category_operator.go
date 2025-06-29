package sleego

import "sync"

var categoryOperator CategoryOperator

// CategoryOperator defines the behavior for managing categories of applications
type CategoryOperator interface {
	GetCategoriesOf(process string) []string
	SetProcessByCategories(categoriesToProcesses map[string][]string)
}

type CategoryOperatorImpl struct {
	mu                  sync.RWMutex
	processByCategories map[string][]string
}

func newCategoryOperator() CategoryOperator {
	return &CategoryOperatorImpl{
		processByCategories: make(map[string][]string),
	}
}

func (c *CategoryOperatorImpl) GetCategoriesOf(process string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.processByCategories[process]
}

func (c *CategoryOperatorImpl) SetProcessByCategories(categoriesToProcesses map[string][]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.processByCategories = make(map[string][]string)

	for category, processes := range categoriesToProcesses {
		for _, proc := range processes {
			c.processByCategories[proc] = append(c.processByCategories[proc], category)
		}
	}
}

func GetCategoryOperator() CategoryOperator {
	if categoryOperator == nil {
		categoryOperator = newCategoryOperator()
	}
	return categoryOperator
}
