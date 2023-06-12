package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

		//
		for _, key := range keysList {
			fmt.Println("key-top:", key)
			output_keyName := ""
			finalValue := 0
			add_output_result := false
			keyMap, ok := key.(map[string]interface{})
			if !ok {
				fmt.Println("Invalid key format")
				//continue
			}
			output_keyName, ok = keyMap["value"].(string)
			if ok {
				fmt.Println("output_keyName:", output_keyName)
			} else {
				fmt.Println("Invalid value format")
			}

			// Extract the key-value pair from the map
			for keyName, keyValue := range key.(map[string]interface{}) {
				fmt.Println("===============================================")
				fmt.Println("keyName:", keyName)
				fmt.Println("keyValue:", keyValue)
				if keyName == "value" {
					continue
				}
				/* // Check if keyValue is a map
				keyValueMap, ok := keyValue.(map[string]interface{})
				if !ok {
					fmt.Println("Invalid keyValue format")
					continue
				} else {
					valueStr, ok := keyValue.(string)
					fmt.Println("Value is a string:", valueStr)
					fmt.Println("Value is a ok:", ok)
				} */
				//var keyValueMap map[string]interface{}
				keyValueMap_textVal := ""
				vendorWeight := 0
				keyValueMap_weight := 0
				keyValueStr := ""
				switch keyValue := keyValue.(type) {
				case map[string]interface{}:
					// keyValue is a map
					keyValueMap := keyValue
					fmt.Println("Value is a map:", keyValue)
					fmt.Println("keyValueMap:", keyValueMap)
					//if value, ok := keyValueMap["value"].(map[string]interface{}); ok {
					if text, ok := keyValueMap["#text"].(string); ok {
						keyValueMap_textVal = text
					}
					if wt, ok := keyValueMap["weight"].(string); ok {
						keyValueMap_weight, _ = strconv.Atoi(wt)
						vendorWeight = keyValueMap_weight
					}
					fmt.Printf("keyName %s:   keyValueMap_textVal: %s  keyValueMap_weight %d \n", keyName, keyValueMap_textVal, keyValueMap_weight)
					//fmt.Println("Value is a string weight:", weight)
					//}
				case string:
					// keyValue is a string
					keyValueStr = keyValue
					//fmt.Println("Value is a string keyValue:", keyValue)
					//fmt.Println("Value is a string keyValueStr:", keyValueStr)
					/* if keyName == "value" {
						output_keyName = keyValue
					} */
					fmt.Printf("keyName  %s :   keyValueStr: %s \n", keyName, keyValueStr)
				default:
					fmt.Println("Invalid keyValue format")
					continue
				}

				//fmt.Println("output_keyName:", output_keyName)
				// Extract the required values from the keyValueMap

				// Extract the values
				/* textfieldValueStr := ""
				weight := 0 */
				/* if value, ok := keyValueMap["value"].(map[string]interface{}); ok {
					if text, ok := value["#text"].(string); ok {
						textfieldValueStr = text
					}
					if wt, ok := value["weight"].(string); ok {
						weight, _ = strconv.Atoi(wt)
					}
				} else {
					if text, ok := keyValueMap["#text"].(string); ok {
						textfieldValueStr = text
					}
					if wt, ok := keyValueMap["weight"].(string); ok {
						weight, _ = strconv.Atoi(wt)
					}
				}

				fmt.Println("textfieldValueStr:", textfieldValueStr)
				fmt.Println("weight:", weight)

				// Extract the token value from the keyValueMap
				fieldValueStr, ok := keyValueMap["value"].(string)
				if !ok {
					fmt.Println("Missing or invalid 'value' field in keyValue")
					continue
				}
				fmt.Println("fieldValueStr:", fieldValueStr) */
				// ... rest of your code for calculating and updating keyValues ...

				fmt.Println("==============================keyToken===================================")
				// Iterate over the fields of the keyToken map and calculate the values
				for keyToken_fieldName, keyToken_fieldValue := range keyToken {
					fmt.Println("keyToken_fieldName:", keyToken_fieldName)
					fmt.Println("keyToken_fieldValue:", keyToken_fieldValue)

					// Check if fieldValue is a map
					keyToken_fieldValueMap, ok := keyToken_fieldValue.(map[string]interface{})
					if !ok {
						fmt.Println("Invalid fieldValue format")
						//continue
					}
					//token
					keyToken_fieldValue_token, ok := keyToken_fieldValueMap["token"].(string)
					if !ok {
						fmt.Println("Missing or invalid 'token' field in fieldValue")
						//continue
					}

					fmt.Println("keyToken_fieldValue_token:", keyToken_fieldValue_token)
					/* if keyName != strings.ToLower(keyToken_fieldValue_token) {
						continue
					} */
					//source
					keyToken_fieldValue_source, ok := keyToken_fieldValueMap["source"].(string)
					if !ok {
						fmt.Println("Missing or invalid 'token' field in fieldValue")
						//continue
					}

					fmt.Println("keyToken_fieldValue_source:", keyToken_fieldValue_source)
					//getSource
					keyToken_fieldValue_getSource, ok := keyToken_fieldValueMap["getSource"].(string)
					if !ok {
						fmt.Println("Missing or invalid 'token' field in fieldValue")
						//continue
					}

					fmt.Println("keyToken_fieldValue_getSource:", keyToken_fieldValue_getSource)

					substr := "attributes/"
					keyToken_fieldValue_getSource_value := ""
					if strings.Contains(keyToken_fieldValue_getSource, substr) {
						fmt.Println("======String contains:", substr)
					} else {
						fmt.Println("String does not contain", substr)
					}
					if strings.Contains(keyToken_fieldValue_getSource, substr) {
						split := strings.Split(keyToken_fieldValue_getSource, substr)
						if len(split) > 1 {
							keyToken_fieldValue_getSource_value = split[1]
							fmt.Println("========>>>Extracted value:", keyToken_fieldValue_getSource_value)
						}
					} else {
						fmt.Println("String does not contain", substr)
					}

					//default-weight
					keyToken_fieldValue_default_weight := int(keyToken_fieldValueMap["default-weight"].(float64))

					fmt.Println("keyToken_fieldValue_default_weight:", keyToken_fieldValue_default_weight)

					fmt.Println("attributes:", attributes)

					keyToken_fieldValue_getSource_value_Attributes := lowercaseAttributes[strings.ToLower(keyToken_fieldValue_getSource_value)]
					fmt.Printf("keyToken_fieldValue_getSource_value:  %s  keyToken_fieldValue_getSource_value_Attributes:%s \n", keyToken_fieldValue_getSource_value, keyToken_fieldValue_getSource_value_Attributes)

					/* if len(keyValueStr) > 1 && keyValueStr != strings.ToLower(keyToken_fieldValue_getSource_value) {
						continue
					} */

					/* value := lowercaseAttributes[strings.ToLower(keyToken_fieldValue_token)]
					fmt.Println("before-value:", value)
					if value == nil {
						continue
						// Handle the case when fieldValue does not exist
						fmt.Println("Field value does not exist")
						continue
					}

					fmt.Println("after-value:", value) */
					// Calculate the default-weight from keytoken.json and vendor.weight from keys.json
					defaultWeight := 0
					//vendorWeight := 0
					//for _, key := range keysList {
					/* fmt.Println("key:", key)
					keyMap, ok := key.(map[string]interface{})
					if !ok {
						continue
					}

					valueStr, ok := keyMap["value"].(string)
					fmt.Println("valueStr:", valueStr)
					if !ok {
						continue
					}
					//textfieldValueStr := ""
					//weight := 0
					_, valueok := keyMap["vendor"].(map[string]interface{})
					_, typeok := keyMap["type"].(map[string]interface{}) */

					//vendorVal := lowercaseAttributes["vendor"]
					//	if (valueStr == "*" || valueStr == "ALL" || keyToken_fieldValue_token == "MODEL" || keyToken_fieldValue_token == "VENDOR") && (valueok || typeok) {
					if output_keyName == "*" || output_keyName == "ALL" {
						/* if value, ok := keyMap["vendor"].(map[string]interface{}); ok {
							if text, ok := value["#text"].(string); ok {
								textfieldValueStr = text
							}
							if wt, ok := value["weight"].(string); ok {
								weight, _ = strconv.Atoi(wt)
								vendorWeight = weight
							}
						} else {
							if text, ok := keyMap["#text"].(string); ok {
								textfieldValueStr = text
							}
							if wt, ok := keyMap["weight"].(string); ok {
								weight, _ = strconv.Atoi(wt)
								vendorWeight = weight
							}
						} */
						add_output_result = true
						fmt.Println("textfieldValueStr:", keyValueMap_textVal)
						fmt.Println("weight:", keyValueMap_weight)
						fmt.Println("defaultWeight:", defaultWeight)
						fmt.Println("vendorWeight:", vendorWeight)
						// Calculate the final value
						finalValue = defaultWeight + vendorWeight
						fmt.Println("finalValue:", finalValue)
						// Update the keyValues map
						/* keyValues[output_keyName] = map[string]interface{}{
							"BestKey": 0,
							"Value":   finalValue,
						} */
						//defaultWeight = int(keyMap["default-weight"].(float64))
						//vendorWeight = int(keyMap["weight"].(float64))
						//break
						//} else if valueStr == "VENDOR" && keyMap["vendor"].(string) == fieldValueStr {
						/*
									else if keyMap["vendor"].(string) == fieldValueStr {
								defaultWeight = int(keyMap["default-weight"].(float64))
								vendorWeight = int(keyMap["weight"].(float64))
								//break
							}
						*/
						//} else if keyMap["vendor"].(string) == fieldValueStr {

					} else if keyMap["vendor"].(string) == keyToken_fieldValue_getSource_value_Attributes {
						add_output_result = true
						fmt.Printf("key  %s :   keyValueStr: %s keyToken_fieldValue_getSource_value_Attributes : %s \n", keyName, keyValueStr, keyToken_fieldValue_getSource_value_Attributes)

						/* if keyName == keyToken_fieldValue_getSource_value_Attributes {
							fmt.Println("keyName == keyToken_fieldValue_getSource_value_Attributes")
						} */
						//defaultWeight = int(keyMap["default-weight"].(float64))
						//vendorWeight = int(keyMap["weight"].(float64))

						//fmt.Println("textfieldValueStr:", textfieldValueStr)
						//fmt.Println("weight:", weight)
						defaultWeight = keyToken_fieldValue_default_weight
						fmt.Println("defaultWeight:", defaultWeight)
						fmt.Println("vendorWeight:", vendorWeight)
						// Calculate the final value
						finalValue = finalValue + defaultWeight + vendorWeight
						fmt.Println("finalValue:", finalValue)
						// Update the keyValues map
						/* keyValues[output_keyName] = map[string]interface{}{
							"BestKey": 0,
							"Value":   finalValue,
						} */

						//break
					} /*  else {
						//textfieldValueStr := ""
						//weight := 0
						if value, ok := keyMap["vendor"].(map[string]interface{}); ok {
							if text, ok := value["#text"].(string); ok {
								textfieldValueStr = text
							}
							if wt, ok := value["weight"].(string); ok {
								weight, _ = strconv.Atoi(wt)
								vendorWeight = weight
							}
						} else {
							if text, ok := keyMap["#text"].(string); ok {
								textfieldValueStr = text
							}
							if wt, ok := keyMap["weight"].(string); ok {
								weight, _ = strconv.Atoi(wt)
								vendorWeight = weight
							}
						}

						fmt.Println("textfieldValueStr:", textfieldValueStr)
						fmt.Println("weight:", weight)
					} */
					//}
					/* fmt.Println("textfieldValueStr:", textfieldValueStr)
					fmt.Println("weight:", weight)
					fmt.Println("defaultWeight:", defaultWeight)
					fmt.Println("vendorWeight:", vendorWeight)
					// Calculate the final value
					finalValue := defaultWeight + vendorWeight
					fmt.Println("finalValue:", finalValue)
					// Update the keyValues map
					keyValues[valueStr] = map[string]interface{}{
						"BestKey": 0,
						"Value":   finalValue,
					} */
				}
			}

			// put finalValue and output_key
			if add_output_result {
				keyValues[output_keyName] = map[string]interface{}{
					"BestKey": 0,
					"Value":   finalValue,
				}
			}

		}
		//

		// Construct the final output
		output := map[string]interface{}{
			"EquipmentName": equipmentName,
			"Keys":          keyValues,
			"Namespace":     namespace,
			"Transaction":   transaction,
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
