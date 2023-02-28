package opensearchtools

// Query wraps all query types in a common interface. Facilitating a common pattern to convert the logical
// struct into the appropriate request JSON.
type Query interface {
	// ToOpenSearchJSON converts the Query struct to the expected OpenSearch JSON
	ToOpenSearchJSON() ([]byte, error)
}

// QueryVersionConverter takes in a domain model Query and makes any modifications or conversions needed for
// a specific version of OpenSearch
type QueryVersionConverter func(Query) (Query, error)

// BoolQueryConverter is a utility support QueryVersionConverter to iterate over all the nested queries in a BoolQuery
func BoolQueryConverter(boolQuery *BoolQuery, converter QueryVersionConverter) (Query, error) {
	must, mErr := convertSubQueries(boolQuery.must, converter)
	if mErr != nil {
		return nil, mErr
	}

	mustNot, mnErr := convertSubQueries(boolQuery.mustNot, converter)
	if mnErr != nil {
		return nil, mnErr
	}

	should, sErr := convertSubQueries(boolQuery.should, converter)
	if sErr != nil {
		return nil, sErr
	}

	filter, fErr := convertSubQueries(boolQuery.filter, converter)
	if fErr != nil {
		return nil, fErr
	}

	return &BoolQuery{
		minimumShouldMatch: boolQuery.minimumShouldMatch,
		must:               must,
		mustNot:            mustNot,
		should:             should,
		filter:             filter,
	}, nil
}

func convertSubQueries(queries []Query, converter QueryVersionConverter) ([]Query, error) {
	var convertedQueries []Query
	for _, q := range queries {
		converted, cErr := converter(q)
		if cErr != nil {
			return nil, cErr
		}

		convertedQueries = append(convertedQueries, converted)
	}

	return convertedQueries, nil
}
