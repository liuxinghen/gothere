package gothere

import "fmt"

func CheckRule(rules []Rule) []RuleValidationError {
	var commonErrors []RuleValidationError
	// rules check : duplicated ToKeys in rules are not allowed
	toKeyRecord := map[string]int{}
	for index, rule := range rules {
		if len(rule.ToKey) == 0 {
			// validation: ToKey is required
			commonErrors = append(commonErrors,
				RuleValidationError{
					Index: index,
					Src:   fmt.Errorf("no ToKey is specified"),
				})
		} else if previousIndex, exists := toKeyRecord[rule.ToKey]; exists {
			// validation: rule ToKey should not be duplicated
			commonErrors = append(commonErrors,
				RuleValidationError{
					Index: index,
					Src:   fmt.Errorf("duplicated ToKey[%s] with rule index[%d]", rule.ToKey, previousIndex),
				})
		} else {
			toKeyRecord[rule.ToKey] = index
		}
		if rule.Generator == nil && len(rule.FromKey) == 0 {
			// validation: Either Generator or FromKey should be specified
			commonErrors = append(commonErrors,
				RuleValidationError{
					Index: index,
					Src:   fmt.Errorf("at least one of Generator and FromKey should be specified"),
				})
		}
	}
	return commonErrors
}

// convert source data map items to another data map items according to the convert rule
func Convert(source []map[string]interface{}, rules []Rule) ([]map[string]interface{}, *ConvertError) {
	if len(source) == 0 {
		// nothing to convert
		return source, nil
	}
	if len(rules) == 0 {
		// if no rule is specified, return source directly
		return source, nil
	}
	convertError := &ConvertError{}
	convertError.CommonErrors = CheckRule(rules)
	if convertError.HasCommonError() {
		// convert process will be skipped when common error happens
		return nil, convertError
	}
	var rowErrors []RowValidationError
	blankRowCount := 0
	result := make([]map[string]interface{}, 0, len(source))
	for rowIndex, row := range source {
		rowResult := make(map[string]interface{})
		var cellErrors []CellValidationError
		blankCellCount := 0
		for _, rule := range rules {
			if rule.Generator != nil {
				if generatedValue, err := rule.Generator(row); err == nil {
					rowResult[rule.ToKey] = generatedValue
				} else {
					// generate value failed
					cellErrors = append(cellErrors, CellValidationError{
						ToKey:  rule.ToKey,
						Type:   ErrorGeneratorError,
						Detail: err,
					})
				}
			} else if sourceValue, exists := row[rule.FromKey]; !exists && rule.Required {
				// sourceValue doesn't exist but required
				cellErrors = append(cellErrors, CellValidationError{
					FromKey: rule.FromKey,
					ToKey:   rule.ToKey,
					Value:   nil,
					Type:    ErrorNoSourceValue,
					Detail:  fmt.Errorf("required but not exists"),
				})
				blankCellCount++
			} else if !exists && !rule.Required {
				// sourceValue doesn't exist and not required, use default
				rowResult[rule.ToKey] = rule.Default()
				blankCellCount++
			} else if rule.Converter != nil {
				if convertedValue, err := rule.Converter(sourceValue); err != nil {
					// convert value failed
					cellErrors = append(cellErrors, CellValidationError{
						FromKey: rule.FromKey,
						ToKey:   rule.ToKey,
						Value:   sourceValue,
						Type:    ErrorConverterError,
						Detail:  err,
					})
				} else {
					rowResult[rule.ToKey] = convertedValue
				}
			} else if rule.Mapping != nil {
				if mappedValue, exists := rule.Mapping[sourceValue]; !exists {
					// map value failed
					cellErrors = append(cellErrors, CellValidationError{
						FromKey: rule.FromKey,
						ToKey:   rule.ToKey,
						Value:   sourceValue,
						Type:    ErrorMappingError,
						Detail:  fmt.Errorf("not exists in mapping"),
					})
				} else {
					rowResult[rule.ToKey] = mappedValue
				}
			} else {
				// directly delivery source value to target
				rowResult[rule.ToKey] = sourceValue
			}
		}
		if blankCellCount == len(rules) {
			blankRowCount++
		}
		if len(rowResult) == len(rules) {
			// successfully convert
			result = append(result, rowResult)
		}
		if len(cellErrors) > 0 {
			// converting failed with cell errors
			rowErrors = append(rowErrors, RowValidationError{
				Index:      rowIndex,
				CellErrors: cellErrors,
			})
		}
	}
	convertError.BlankRowCount = blankRowCount
	convertError.RowErrors = rowErrors
	return result, convertError
}
