package opensearchtools

import (
	"fmt"
)

func validateMGetRequestForV2(m MGetRequest) ValidationResults {
	validationResults := make(ValidationResults, 0)

	// ensure Index is either set at the top level or set in each of the Docs
	// ensure that ID() is non-empty for each Doc
	topLevelIndexIsEmpty := m.Index == ""
	for _, d := range m.Docs {
		if topLevelIndexIsEmpty && d.Index() == "" {
			validationResults = append(validationResults, NewValidationResult(fmt.Sprintf("Index not set at the MGetRequest level nor in the Doc with ID %s", d.ID()), true))
		}
		if d.ID() == "" {
			validationResults = append(validationResults, NewValidationResult("Doc ID is empty", true))
		}
	}

	return validationResults
}
