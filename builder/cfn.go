package builder

import "codecommit/builders/cfn-spec-go/spec"

type cfnBuilder struct {
	Builder
}

var Cfn = cfnBuilder{}

func init() {
	Cfn.Spec = spec.Cfn
}

func (b cfnBuilder) Template(config map[string]string) map[string]interface{} {
	// Generate resources
	resources := make(map[string]interface{})
	for name, resourceType := range config {
		resources[name] = b.newResource(resourceType)
	}

	// Build the template
	return map[string]interface{}{
		"AWSTemplateFormatVersion": "2010-09-09",
		"Description":              "Template generated by cfn-build",
		"Resources":                resources,
		// TODO: "Outputs": outputs,
	}
}