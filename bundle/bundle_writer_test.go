package bundle

import (
	"reflect"
	"sort"
	"testing"
)

func TestTemplateJSONDependencyTraversal(t *testing.T) {
	templateJSON := generateTemplateJSON()
	testTemplateJSONDependencyTraversal(t, &templateJSON, []string{"stencils/one"}, []string{"stencils/one", "stencils/two", "stencils/three"})
	testTemplateJSONDependencyTraversal(t, &templateJSON, []string{"stencils/two"}, []string{"stencils/one", "stencils/two", "stencils/three"})
	testTemplateJSONDependencyTraversal(t, &templateJSON, []string{"stencils/three"}, []string{"stencils/one", "stencils/two", "stencils/three"})
	testTemplateJSONDependencyTraversal(t, &templateJSON, []string{"stencils/one", "stencils/two", "stencils/three"}, []string{"stencils/one", "stencils/two", "stencils/three"})
	testTemplateJSONDependencyTraversal(t, &templateJSON, []string{}, []string{})

	anotherTemplateJSON := generateAnotherTemplateJSON()
	testTemplateJSONDependencyTraversal(t, &anotherTemplateJSON, []string{}, []string{})
	testTemplateJSONDependencyTraversal(t, &anotherTemplateJSON, []string{"stencils/one"}, []string{"stencils/one", "stencils/two", "stencils/three", "stencils/four", "stencils/five", "stencils/six"})
	testTemplateJSONDependencyTraversal(t, &anotherTemplateJSON, []string{"stencils/two"}, []string{"stencils/one", "stencils/two", "stencils/three", "stencils/four", "stencils/five", "stencils/six"})
	testTemplateJSONDependencyTraversal(t, &anotherTemplateJSON, []string{"stencils/three"}, []string{"stencils/one", "stencils/two", "stencils/three", "stencils/four", "stencils/five", "stencils/six"})
	testTemplateJSONDependencyTraversal(t, &anotherTemplateJSON, []string{"stencils/four"}, []string{"stencils/four", "stencils/five", "stencils/six"})
	testTemplateJSONDependencyTraversal(t, &anotherTemplateJSON, []string{"stencils/five"}, []string{"stencils/five"})
	testTemplateJSONDependencyTraversal(t, &anotherTemplateJSON, []string{"stencils/six"}, []string{"stencils/six"})
	testTemplateJSONDependencyTraversal(t, &anotherTemplateJSON, []string{"stencils/seven"}, []string{"stencils/seven", "stencils/eight", "stencils/nine"})
	testTemplateJSONDependencyTraversal(t, &anotherTemplateJSON, []string{"stencils/eight"}, []string{"stencils/eight", "stencils/nine"})
	testTemplateJSONDependencyTraversal(t, &anotherTemplateJSON, []string{"stencils/nine"}, []string{"stencils/nine"})
	testTemplateJSONDependencyTraversal(t, &anotherTemplateJSON, []string{"stencils/ten"}, []string{"stencils/ten"})
	testTemplateJSONDependencyTraversal(t, &anotherTemplateJSON, []string{"stencils/eleven"}, []string{"stencils/eleven"})
	testTemplateJSONDependencyTraversal(t, &anotherTemplateJSON, []string{"stencils/four", "stencils/seven", "stencils/eleven"}, []string{"stencils/four", "stencils/five", "stencils/six", "stencils/seven", "stencils/eight", "stencils/nine", "stencils/eleven"})

	yetAnotherTemplateJSON := generateYetAnotherTemplateJSON()
	testTemplateJSONDependencyTraversalError(t, &yetAnotherTemplateJSON, []string{"stencils/one"})

	testTemplateJSONDependencyTraversalError(t, &anotherTemplateJSON, []string{"stencils/nonexistent"})
}

func testTemplateJSONDependencyTraversal(t *testing.T, templateJSON *TemplateJSON, initialComponentNames []string, expectedComponentNames []string) {
	requiredComponentNames, err := getRequiredComponentNames(templateJSON, initialComponentNames)
	if err != nil {
		t.Errorf("Obtained error when determining dependency tree: %s\n", err)
		return
	}

	sort.Strings(requiredComponentNames)
	sort.Strings(expectedComponentNames)
	if !reflect.DeepEqual(requiredComponentNames, expectedComponentNames) {
		t.Errorf("Expected dependency tree to yield %v, but got %v\n", expectedComponentNames, requiredComponentNames)
	}
}

func testTemplateJSONDependencyTraversalError(t *testing.T, templateJSON *TemplateJSON, initialComponentNames []string) {
	_, err := getRequiredComponentNames(templateJSON, initialComponentNames)
	if err != nil {
		return
	}
	t.Errorf("Expected error when determining dependency tree")
}

func generateTemplateJSON() TemplateJSON {
	stencilOne := StencilTemplate{
		Name:         "one",
		Dependencies: []string{"stencils/two"},
	}

	stencilTwo := StencilTemplate{
		Name:         "two",
		Dependencies: []string{"stencils/three"},
	}

	stencilThree := StencilTemplate{
		Name:         "three",
		Dependencies: []string{"stencils/one"},
	}

	templateStruct := TemplatesStruct{
		Stencils: []*StencilTemplate{&stencilOne, &stencilTwo, &stencilThree},
	}

	return TemplateJSON{
		Templates: &templateStruct,
	}
}

func generateAnotherTemplateJSON() TemplateJSON {
	stencilOne := StencilTemplate{
		Name:         "one",
		Dependencies: []string{"stencils/two", "stencils/three"},
	}

	stencilTwo := StencilTemplate{
		Name:         "two",
		Dependencies: []string{"stencils/three"},
	}

	stencilThree := StencilTemplate{
		Name:         "three",
		Dependencies: []string{"stencils/one", "stencils/three", "stencils/four"},
	}

	stencilFour := StencilTemplate{
		Name:         "four",
		Dependencies: []string{"stencils/five", "stencils/six"},
	}

	stencilFive := StencilTemplate{
		Name:         "five",
		Dependencies: []string{},
	}

	stencilSix := StencilTemplate{
		Name:         "six",
		Dependencies: []string{},
	}

	stencilSeven := StencilTemplate{
		Name:         "seven",
		Dependencies: []string{"stencils/eight", "stencils/nine"},
	}

	stencilEight := StencilTemplate{
		Name:         "eight",
		Dependencies: []string{"stencils/nine"},
	}

	stencilNine := StencilTemplate{
		Name:         "nine",
		Dependencies: []string{},
	}

	stencilTen := StencilTemplate{
		Name:         "ten",
		Dependencies: []string{},
	}

	stencilEleven := StencilTemplate{
		Name:         "eleven",
		Dependencies: []string{},
	}

	templateStruct := TemplatesStruct{
		Stencils: []*StencilTemplate{&stencilOne, &stencilTwo, &stencilThree, &stencilFour, &stencilFive, &stencilSix, &stencilSeven, &stencilEight, &stencilNine, &stencilTen, &stencilEleven},
	}

	return TemplateJSON{
		Templates: &templateStruct,
	}
}

func generateYetAnotherTemplateJSON() TemplateJSON {
	stencilOne := StencilTemplate{
		Name:         "one",
		Dependencies: []string{"stencils/two"},
	}

	stencilTwo := StencilTemplate{
		Name:         "two",
		Dependencies: []string{"stencils/nonexistent"},
	}

	templateStruct := TemplatesStruct{
		Stencils: []*StencilTemplate{&stencilOne, &stencilTwo},
	}

	return TemplateJSON{
		Templates: &templateStruct,
	}
}
