package builder

import (
	"codecommit/builders/cfn-spec-go/spec"
)

type TemplateConfig struct {
	Resources                 map[string]string
	IncludeOptionalProperties bool
}

func NewTemplateConfig() TemplateConfig {
	return TemplateConfig{
		Resources: make(map[string]string),
	}
}

func NewTemplate(config TemplateConfig) map[string]interface{} {
	// Generate resources
	resources := make(map[string]interface{})
	for name, resourceType := range config.Resources {
		resources[name] = newResource(resourceType)
	}

	// Build the template
	return map[string]interface{}{
		"AWSTemplateFormatVersion": "2010-09-09",
		"Description":              "Template generated by cfn-build",
		"Resources":                resources,
		// TODO: "Outputs": outputs,
	}
}

func newResource(resourceType string) map[string]interface{} {
	rSpec, ok := spec.Cfn.ResourceTypes[resourceType]
	if !ok {
		panic("No such resource type: " + resourceType)
	}

	// Generate properties
	properties := make(map[string]interface{})
	for name, pSpec := range rSpec.Properties {
		properties[name] = newProperty(resourceType, pSpec)
	}

	return map[string]interface{}{
		"Type":       resourceType,
		"Properties": properties,
	}
}

func newProperty(resourceType string, pSpec spec.Property) interface{} {
	// Correctly badly-formed entries
	if pSpec.PrimitiveType == "Map" {
		pSpec.PrimitiveType = ""
		pSpec.Type = "Map"
	}

	// Primitive types
	if pSpec.PrimitiveType != "" {
		return newPrimitive(pSpec.PrimitiveType)
	}

	if pSpec.Type == "List" || pSpec.Type == "Map" {
		var value interface{}

		if pSpec.PrimitiveItemType != "" {
			value = newPrimitive(pSpec.PrimitiveItemType)
		} else if pSpec.ItemType != "" {
			value = newPropertyType(resourceType, pSpec.ItemType)
		} else {
			value = "CHANGEME"
		}

		if pSpec.Type == "List" {
			return []interface{}{value}
		}

		return map[string]interface{}{"CHANGEME": value}
	}

	// Fall through to property types
	return newPropertyType(resourceType, pSpec.Type)
}

func newPrimitive(primitiveType string) interface{} {
	switch primitiveType {
	case "String":
		return "CHANGEME"
	case "Integer":
		return 0
	case "Double":
		return 0.0
	case "Long":
		return 0.0
	case "Boolean":
		return false
	case "Timestamp":
		return "1970-01-01 00:00:00"
	case "Json":
		return "{\"JSON\": \"CHANGEME\"}"
	default:
		panic("PRIMITIVE NOT IMPLEMENTED: " + primitiveType)
	}
}

func newPropertyType(resourceType, propertyType string) interface{} {
	var ptSpec spec.PropertyType
	var ok bool

	ptSpec, ok = spec.Cfn.PropertyTypes[propertyType]
	if !ok {
		ptSpec, ok = spec.Cfn.PropertyTypes[resourceType+"."+propertyType]
	}
	if !ok {
		panic("PTYPE NOT IMPLEMENTED: " + resourceType + "." + propertyType)
	}

	// Generate properties
	properties := make(map[string]interface{})
	for name, pSpec := range ptSpec.Properties {
		if pSpec.Type == propertyType || pSpec.ItemType == propertyType {
			properties[name] = make(map[string]interface{})
		} else {
			properties[name] = newProperty(resourceType, pSpec)
		}
	}

	return properties
}