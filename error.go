package gothere

import (
	"fmt"
	"strings"
)

const (
	ErrorNoSourceValue = iota
	ErrorGeneratorError
	ErrorConverterError
	ErrorMappingError
)

type RuleValidationError struct {
	Index int
	Src   error
}

func (r RuleValidationError) Error() string {
	return fmt.Sprintf("RuleIndex[%d] : %s", r.Index, r.Src.Error())
}

type CellValidationError struct {
	FromKey string
	ToKey   string
	Value   interface{}
	Type    int
	Detail  error
}

func (c CellValidationError) Error() string {
	return fmt.Sprintf("FromKey[%s],ToKey[%s],Value[%+v],ErrorType[%d] : %s",
		c.FromKey, c.ToKey, c.Value, c.Type, c.Detail.Error())
}

type RowValidationError struct {
	Index      int
	CellErrors []CellValidationError
}

func (r RowValidationError) Error() string {
	builder := strings.Builder{}
	for _, cellError := range r.CellErrors {
		builder.WriteString(fmt.Sprintf("RowIndex[%d],%s\n", r.Index, cellError.Error()))
	}
	return builder.String()
}

type ConvertError struct {
	BlankRowCount int
	CommonErrors  []RuleValidationError
	RowErrors     []RowValidationError
}

func (c *ConvertError) Error() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintln("CommonValidationErrors:"))
	for _, commonError := range c.CommonErrors {
		builder.WriteString(fmt.Sprintln(commonError.Error()))
	}
	builder.WriteString(fmt.Sprintln("RowValidationErrors:"))
	for _, rowError := range c.RowErrors {
		builder.WriteString(fmt.Sprintln(rowError.Error()))
	}
	builder.WriteString(fmt.Sprintf("BlankRowCount : %d\n", c.BlankRowCount))
	return builder.String()
}

func (c *ConvertError) HasCommonError() bool {
	return len(c.CommonErrors) > 0
}

func (c *ConvertError) HasRowError() bool {
	return len(c.RowErrors) > 0
}

func (c *ConvertError) IsNil() bool {
	return !c.HasCommonError() && !c.HasRowError() && c.BlankRowCount == 0
}
