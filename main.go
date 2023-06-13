package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	// Read the keytoken.json file
	keyTokenFile, err := os.ReadFile("keytoken.json")
	if err != nil {
		fmt.Println("Failed to read keytoken.json:", err)
		return
	}

	var keyToken map[string]interface{}
	err = json.Unmarshal(keyTokenFile, &keyToken)
	if err != nil {
		fmt.Println("Failed to parse keytoken.json:", err)
		return
	}
	fmt.Println("keyToken:", keyToken)

	// Read the keys.json file
	keysFile, err := os.ReadFile("keys.json")
	if err != nil {
		fmt.Println("Failed to read keys.json:", err)
		return
	}

	var keys map[string]interface{}
	err = json.Unmarshal(keysFile, &keys)
	if err != nil {
		fmt.Println("Failed to parse keys.json:", err)
		return
	}

	// Process the POST request
	http.HandleFunc("/default", func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var requestBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}

		// Extract the required fields from the request body
		equipmentName, ok := requestBody["equipmentname"].(string)
		if !ok {
			http.Error(w, "Missing or invalid 'equipmentname' field", http.StatusBadRequest)
			return
		}

		namespace, ok := requestBody["namespace"].(string)
		if !ok {
			http.Error(w, "Missing or invalid 'namespace' field", http.StatusBadRequest)
			return
		}

		var transaction map[string]interface{}
		if transactionData, ok := requestBody["transaction"].(map[string]interface{}); ok {
			transaction = transactionData
		}

		attributes, ok := requestBody["attributes"].(map[string]interface{})
		if !ok {
			http.Error(w, "Missing or invalid 'attributes' field", http.StatusBadRequest)
			return
		}
		// Convert the keys in the attributes map to lowercase
		lowercaseAttributes := make(map[string]interface{})
		fmt.Println("----------attributes---start----")
		for key, value := range attributes {
			lowercaseKey := strings.ToLower(key)
			fmt.Printf("lowercaseKey:%s value : %s \n", lowercaseKey, value)
			lowercaseAttributes[lowercaseKey] = value
		}
		fmt.Println("----------attributes----end---")
		fmt.Println("lowercaseAttributes[vendor] :", lowercaseAttributes["vendor"])
		fmt.Println("lowercaseAttributes[software] :", lowercaseAttributes["software"])
		fmt.Println("lowercaseAttributes[model] :", lowercaseAttributes["model"])
		// Create a map to store the calculated key values
		keyValues := make(map[string]interface{})

		// Get the keys list from the keys.json file
		keysList := keys["keys"].(map[string]interface{})["key"].([]interface{})
		fmt.Println("keysList:", keysList)

		//keysList
		for _, key := range keysList {
			fmt.Println("key-top:", key)
			output_keyName := ""
			finalValue := 0
			add_output_result := false
			keyMap, ok := key.(map[string]interface{})
			if !ok {
				fmt.Println("Invalid key format")

			}
			output_keyName, ok = keyMap["value"].(string)
			if ok {
				fmt.Println("output_keyName:", output_keyName)
			} else {
				fmt.Println("Invalid value format")
			}
			ncregexp_Keyvalue, ok := keyMap["ncregexp"].(string)
			if ok {
				fmt.Println("ncregexp_Keyvalue:", ncregexp_Keyvalue)
			} else {
				fmt.Println("Invalid value format ncregexp_Keyvalue")
			}
			// Extract the key-value pair from the map
			for keyName, keyValue := range key.(map[string]interface{}) {
				fmt.Println("===============================================")
				fmt.Println("keyName:", keyName)
				fmt.Println("keyValue:", keyValue)
				if keyName == "value" {
					continue
				}

				keyValueMap_textVal := ""

				keyValueMap_weight := 0
				keyValueStr := ""
				switch keyValue := keyValue.(type) {
				case []interface{}:
					// Check if keyValue is a slice
					if len(keyValue) > 0 {
						// Get the first element of the slice
						firstElem, ok := keyValue[0].(map[string]interface{})
						if !ok {
							fmt.Println("Invalid keyValue format")
							continue
						}
						if text, ok := firstElem["#text"].(string); ok {
							keyValueMap_textVal = text
						}
						if wt, ok := firstElem["weight"].(string); ok {
							keyValueMap_weight, _ = strconv.Atoi(wt)

						}
						fmt.Printf("keyName %s:   keyValueMap_textVal: %s  keyValueMap_weight %d \n", keyName, keyValueMap_textVal, keyValueMap_weight)
					}
				case map[string]interface{}:
					// keyValue is a map
					keyValueMap := keyValue
					fmt.Println("Value is a map:", keyValue)
					fmt.Println("keyValueMap:", keyValueMap)

					if text, ok := keyValueMap["#text"].(string); ok {
						keyValueMap_textVal = text
					}
					if wt, ok := keyValueMap["weight"].(string); ok {
						keyValueMap_weight, _ = strconv.Atoi(wt)

					}
					fmt.Printf("keyName %s:   keyValueMap_textVal: %s  keyValueMap_weight %d \n", keyName, keyValueMap_textVal, keyValueMap_weight)

				case string:
					// keyValue is a string
					keyValueStr = keyValue

					fmt.Printf("keyName  %s :   keyValueStr: %s \n", keyName, keyValueStr)
				default:
					fmt.Println("Invalid keyValue format")
					continue
				}

				fmt.Println("==============================keyToken===================================")
				// Iterate over the fields of the keyToken map and calculate the values
				for keyToken_fieldName, keyToken_fieldValue := range keyToken {
					fmt.Println("keyToken_fieldName:", keyToken_fieldName)
					fmt.Println("keyToken_fieldValue:", keyToken_fieldValue)

					// Check if fieldValue is a map
					keyToken_fieldValueMap, ok := keyToken_fieldValue.(map[string]interface{})
					if !ok {
						fmt.Println("Invalid fieldValue format")

					}
					//ncregexp

					keyToken_fieldValueStr := fmt.Sprint(keyToken_fieldValue)
					fmt.Println("keyToken_fieldValueStr:", keyToken_fieldValueStr)
					matchRegex := false
					if keyName == "ncregexp" && keyName == keyToken_fieldName && len(keyValueStr) > 1 {
						equipmentNameUpper := strings.ToUpper(fmt.Sprint(equipmentName))

						ncregexpstr := fmt.Sprint(keyValueStr)
						fmt.Println("ncregexpstr:", ncregexpstr)
						ncregexp := regexp.MustCompile(ncregexpstr)
						matchRegex = ncregexp.Match([]byte(equipmentNameUpper))
						fmt.Printf("ncregexpstr  %s :   matchRegex: %v \n", ncregexpstr, matchRegex)
					}

					//token
					keyToken_fieldValue_token, ok := keyToken_fieldValueMap["token"].(string)
					if !ok {
						fmt.Println("Missing or invalid 'token' field in fieldValue")

					}

					fmt.Println("keyToken_fieldValue_token:", keyToken_fieldValue_token)

					//source
					keyToken_fieldValue_source, ok := keyToken_fieldValueMap["source"].(string)
					if !ok {
						fmt.Println("Missing or invalid 'token' field in fieldValue")

					}

					fmt.Println("keyToken_fieldValue_source:", keyToken_fieldValue_source)
					//getSource
					keyToken_fieldValue_getSource, ok := keyToken_fieldValueMap["getSource"].(string)
					if !ok {
						fmt.Println("Missing or invalid 'token' field in fieldValue")

					}

					fmt.Println("keyToken_fieldValue_getSource:", keyToken_fieldValue_getSource)

					substr := "attributes/"
					keyToken_fieldValue_getSource_value := ""
					if strings.Contains(keyToken_fieldValue_getSource, substr) {
						fmt.Println("======String contains:", substr)
						split := strings.Split(keyToken_fieldValue_getSource, substr)
						if len(split) > 1 {
							keyToken_fieldValue_getSource_value = split[1]
							fmt.Println("========>>>Extracted value:", keyToken_fieldValue_getSource_value)
						}
					} else {
						keyToken_fieldValue_getSource_value = keyToken_fieldValue_getSource
					}
					/* if strings.Contains(keyToken_fieldValue_getSource, substr) {
						split := strings.Split(keyToken_fieldValue_getSource, substr)
						if len(split) > 1 {
							keyToken_fieldValue_getSource_value = split[1]
							fmt.Println("========>>>Extracted value:", keyToken_fieldValue_getSource_value)
						}
					} else {
						keyToken_fieldValue_getSource_value = keyToken_fieldValue_getSource
					} */

					//default-weight
					keyToken_fieldValue_default_weight := int(keyToken_fieldValueMap["default-weight"].(float64))

					fmt.Println("keyToken_fieldValue_default_weight:", keyToken_fieldValue_default_weight)

					fmt.Println("attributes:", attributes)

					keyToken_fieldValue_getSource_value_Attributes := lowercaseAttributes[strings.ToLower(keyToken_fieldValue_getSource_value)]
					fmt.Printf("keyToken_fieldValue_getSource_value:  %s  keyToken_fieldValue_getSource_value_Attributes:%s \n", keyToken_fieldValue_getSource_value, keyToken_fieldValue_getSource_value_Attributes)
					keyToken_fieldValue_getSource_value_AttributesStr := ""
					if keyToken_fieldValue_getSource_value_Attributes != nil {
						keyToken_fieldValue_getSource_value_AttributesStr = fmt.Sprint(keyToken_fieldValue_getSource_value_Attributes)
					}
					fmt.Println("keyToken_fieldValue_getSource_value_AttributesStr:", keyToken_fieldValue_getSource_value_AttributesStr)

					if output_keyName == "*" || output_keyName == "ALL" {

						add_output_result = true
						fmt.Println("textfieldValueStr:", keyValueMap_textVal)
						fmt.Println("weight:", keyValueMap_weight)

						fmt.Println("keyValueMap_weight:", keyValueMap_weight)
						// Calculate the final value

						finalValue = keyValueMap_weight
						fmt.Println("finalValue:", finalValue)

					} else if (keyName == keyToken_fieldName && len(keyValueStr) > 1 && len(keyToken_fieldValue_getSource_value_AttributesStr) > 1 && keyValueStr == keyToken_fieldValue_getSource_value_AttributesStr) || (matchRegex) {
						add_output_result = true
						fmt.Printf("key  %s :   keyValueStr: %s keyToken_fieldValue_getSource_value_AttributesStr : %s \n", keyName, keyValueStr, keyToken_fieldValue_getSource_value_AttributesStr)

						defaultWeight := keyToken_fieldValue_default_weight
						fmt.Println("defaultWeight:", defaultWeight)
						fmt.Println("keyValueMap_weight:", keyValueMap_weight)
						// Calculate the final value
						finalValue = finalValue + defaultWeight + keyValueMap_weight
						fmt.Println("finalValue:", finalValue)

					}
				} //keyToken
				//}
			} //Extract the key-value pair from the map

			// put finalValue and output_key
			if add_output_result {
				keyValues[output_keyName] = map[string]interface{}{
					"bestkey": 0,
					"value":   finalValue,
				}
			}

		} //keysList
		//

		// Construct the final output
		output := map[string]interface{}{
			"equipmentname": equipmentName,
			"keys":          keyValues,
			"namespace":     namespace,
			"transaction":   transaction,
		}

		// Convert the output to JSON
		outputJSON, err := json.Marshal(output)
		if err != nil {
			http.Error(w, "Failed to serialize output", http.StatusInternalServerError)
			return
		}

		// Set the response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Write the response
		_, err = w.Write(outputJSON)
		if err != nil {
			fmt.Println("Failed to write response:", err)
		}
	})

	// Start the server
	fmt.Println("Server listening on port 8083...")
	err = http.ListenAndServe(":8083", nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
