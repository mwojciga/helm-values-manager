package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

// HelmValues defines a map for storing YAML data.
type HelmValues map[string]interface{}

// Function to read and parse a YAML file.
func readYAMLFile(filename string) (HelmValues, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	values := make(HelmValues)
	err = yaml.Unmarshal(data, &values)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filename, err)
	}

	return values, nil
}

// Merge values, keeping track of where values were overridden and the final value.
func mergeValues(base, override HelmValues, path string, overrideHistory map[string][]string, finalValues HelmValues, filename string) {
	// Iterate over keys in the override map
	for k, v := range override {
		fullKey := path + "." + k

		switch vTyped := v.(type) {
		case map[interface{}]interface{}:
			// Handle nested maps recursively
			baseVal, baseExists := base[k].(map[interface{}]interface{})
			if baseExists {
				if _, ok := finalValues[k].(HelmValues); !ok {
					finalValues[k] = make(HelmValues)
				}
				mergeValues(convertMap(baseVal), convertMap(vTyped), fullKey, overrideHistory, finalValues[k].(HelmValues), filename)
			} else {
				finalValues[k] = convertMap(vTyped)
				overrideHistory[fullKey] = append(overrideHistory[fullKey], fmt.Sprintf("New map structure from %s", filename))
			}
		case []interface{}:
			// Handle lists
			baseVal, baseExists := base[k].([]interface{})
			if baseExists && !reflect.DeepEqual(baseVal, vTyped) {
				overrideHistory[fullKey] = append(overrideHistory[fullKey], fmt.Sprintf("Base list: %v, Override list: %v (%s)", baseVal, vTyped, filename))
			} else if !baseExists {
				overrideHistory[fullKey] = append(overrideHistory[fullKey], fmt.Sprintf("New list: %v (%s)", vTyped, filename))
			}
			finalValues[k] = vTyped // Replace the list with the override list
		default:
			// Log the override and update final value if it's a scalar or leaf value
			finalValues[k] = v
			if baseVal, exists := base[k]; exists && baseVal != v {
				overrideHistory[fullKey] = append(overrideHistory[fullKey], fmt.Sprintf("Base value: %v, Override value: %v (%s)", baseVal, v, filename))
			} else if !exists {
				overrideHistory[fullKey] = append(overrideHistory[fullKey], fmt.Sprintf("New value: %v (%s)", v, filename))
			}
		}
	}
}

// Helper function to convert map[interface{}]interface{} to HelmValues (map[string]interface{}).
func convertMap(input map[interface{}]interface{}) HelmValues {
	output := make(HelmValues)
	for key, value := range input {
		strKey := fmt.Sprintf("%v", key) // Convert all keys to strings
		output[strKey] = value
	}
	return output
}

// Helper function to get final values from nested maps.
func getFinalValue(finalValues HelmValues, keys []string) (interface{}, bool) {
	currentMap := finalValues
	for i, key := range keys {
		if i == len(keys)-1 {
			// Last key, return the value
			val, exists := currentMap[key]
			return val, exists
		}
		// Traverse the nested map
		if nextMap, exists := currentMap[key].(HelmValues); exists {
			currentMap = nextMap
		} else {
			return nil, false
		}
	}
	return nil, false
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Please provide at least a base file and one override file.")
	}

	// Read the base values file (first argument).
	baseFilename := os.Args[1]
	baseValues, err := readYAMLFile(baseFilename)
	if err != nil {
		log.Fatalf("Error reading base values file %s: %v", baseFilename, err)
	}

	// Initialize finalValues with baseValues
	finalValues := make(HelmValues)

	// Initialize tracking for overrides
	overrideHistory := make(map[string][]string)

	// Process override files (subsequent arguments)
	for i := 2; i < len(os.Args); i++ {
		overrideFilename := os.Args[i]
		overrideValues, err := readYAMLFile(overrideFilename)
		if err != nil {
			log.Fatalf("Error reading override values file %s: %v", overrideFilename, err)
		}

		// Merge the current finalValues with the new override values and log overrides
		mergeValues(baseValues, overrideValues, "", overrideHistory, finalValues, overrideFilename)
	}

	// Output the results: base values, overrides, and the final values.
	fmt.Println("--- Override and Final Values Log ---")
	for key, history := range overrideHistory {
		fmt.Printf("Key %s:\n", key)
		for _, entry := range history {
			fmt.Printf("  %s\n", entry)
		}
		// Print the final value for this key
		keys := key[1:]                      // Remove leading dot for correct key lookup
		keyParts := strings.Split(keys, ".") // Split by dot to get the parts
		finalVal, exists := getFinalValue(finalValues, keyParts)
		if exists {
			fmt.Printf("  Final value: %v\n", finalVal)
		} else {
			fmt.Printf("  Final value: <nil>\n")
		}
		fmt.Println("---")
	}
}
