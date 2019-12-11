package bundle

import (
	"github.com/cloud66-oss/starter/bundle/templates"
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

func testTemplateJSONDependencyTraversal(t *testing.T, template *templates.Template, initialComponentNames []string, expectedComponentNames []string) {
	requiredComponentNames, err := getDependencyComponents(template, initialComponentNames)
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

func testTemplateJSONDependencyTraversalError(t *testing.T, templateJSON *templates.Template, initialComponentNames []string) {
	_, err := getDependencyComponents(templateJSON, initialComponentNames)
	if err != nil {
		return
	}
	t.Errorf("Expected error when determining dependency tree")
}

func generateTemplateJSON() templates.Template {
	stencilOne := templates.Stencil{
		Name:         "one",
		Dependencies: []string{"stencils/two"},
	}

	stencilTwo := templates.Stencil{
		Name:         "two",
		Dependencies: []string{"stencils/three"},
	}

	stencilThree := templates.Stencil{
		Name:         "three",
		Dependencies: []string{"stencils/one"},
	}

	theTemplates := templates.Templates{
		Stencils: []*templates.Stencil{&stencilOne, &stencilTwo, &stencilThree},
	}

	return templates.Template{
		Templates: &theTemplates,
	}
}

func generateAnotherTemplateJSON() templates.Template {
	stencilOne := templates.Stencil{
		Name:         "one",
		Dependencies: []string{"stencils/two", "stencils/three"},
	}

	stencilTwo := templates.Stencil{
		Name:         "two",
		Dependencies: []string{"stencils/three"},
	}

	stencilThree := templates.Stencil{
		Name:         "three",
		Dependencies: []string{"stencils/one", "stencils/three", "stencils/four"},
	}

	stencilFour := templates.Stencil{
		Name:         "four",
		Dependencies: []string{"stencils/five", "stencils/six"},
	}

	stencilFive := templates.Stencil{
		Name:         "five",
		Dependencies: []string{},
	}

	stencilSix := templates.Stencil{
		Name:         "six",
		Dependencies: []string{},
	}

	stencilSeven := templates.Stencil{
		Name:         "seven",
		Dependencies: []string{"stencils/eight", "stencils/nine"},
	}

	stencilEight := templates.Stencil{
		Name:         "eight",
		Dependencies: []string{"stencils/nine"},
	}

	stencilNine := templates.Stencil{
		Name:         "nine",
		Dependencies: []string{},
	}

	stencilTen := templates.Stencil{
		Name:         "ten",
		Dependencies: []string{},
	}

	stencilEleven := templates.Stencil{
		Name:         "eleven",
		Dependencies: []string{},
	}

	theTemplates := templates.Templates{
		Stencils: []*templates.Stencil{&stencilOne, &stencilTwo, &stencilThree, &stencilFour, &stencilFive, &stencilSix, &stencilSeven, &stencilEight, &stencilNine, &stencilTen, &stencilEleven},
	}

	return templates.Template{
		Templates: &theTemplates,
	}
}

func generateYetAnotherTemplateJSON() templates.Template {
	stencilOne := templates.Stencil{
		Name:         "one",
		Dependencies: []string{"stencils/two"},
	}

	stencilTwo := templates.Stencil{
		Name:         "two",
		Dependencies: []string{"stencils/nonexistent"},
	}

	theTemplates := templates.Templates{
		Stencils: []*templates.Stencil{&stencilOne, &stencilTwo},
	}

	return templates.Template{
		Templates: &theTemplates,
	}
}
