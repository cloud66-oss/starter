package bundle

import (
	"reflect"
	"sort"
	"testing"
)

func TestTemplateJSONDependencyTraversal(t *testing.T) {
	templateJSON := generateTemplateJSON()
	testTemplateJSONDependencyTraversal(t, &templateJSON, []string{"stencils/one", "stencils/two", "stencils/three"})
}

func testTemplateJSONDependencyTraversal(t *testing.T, templateJSON *TemplateJSON, expectedComponentNames []string) {
	requiredComponentNames, err := getRequiredComponentNames(templateJSON)
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

func generateTemplateJSON() TemplateJSON {
	stencilOne := StencilTemplate{
		Name:         "one",
		MinUsage:     1,
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
