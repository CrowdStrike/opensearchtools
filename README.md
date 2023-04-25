# opensearchtools

[![GoDoc](https://pkg.go.dev/badge/github.com/CrowdStrike/opensearchtools.svg)](https://pkg.go.dev/github.com/CrowdStrike/opensearchtools)

opensearchtools is a Go library for working with OpenSearch v2 and beyond. It is designed to assist users who are familiar with Elasticsearch and the [olivere/elastic library](https://github.com/olivere/elastic), which is no longer maintained. The library aims to provide common use-cases for OpenSearch, while remaining mindful that OpenSearch's surface area is quite large and it may not be feasible to implement every possible feature.

## Features

The current set of features that opensearchtools provides is:

- [Bulk](https://opensearch.org/docs/latest/api-reference/document-apis/bulk/)
  - BulkIndex: Index actions create a document if it doesn’t yet exist and replace the document if it already exists.
  - BulkCreate: Creates a document if it doesn’t already exist and returns an error otherwise. Note that this will return an error if you attempt to create a document with an ID that already exists.
  - BulkUpdate: This action updates existing documents and returns an error if the document doesn’t exist. 
  - BulkDelete: Deletes a document if it exists
- [Multi-Get](https://opensearch.org/docs/latest/api-reference/document-apis/multi-get/)
- [Search](https://opensearch.org/docs/latest/api-reference/search/)
  - Query
  - Aggregations

## Concepts

- executor
- validation

## Detailed Walk-through

Here is a detailed walk-through that shows the overal concepts of using opensearchtools.

opensearchtools provides an agnostic to OpenSearch version. However, it supports only v2 (not v1). We intend to add support for future OpenSearch versions as they arise.

This version-agnostic API is implemented using an executor type pattern: you create your get/bulk/search related objects using the version agnostic domain objects, and then you create an instance of the version-specific executor (currently there is only an executor for v2 since that is the only version supported) and pass your domain-specific objects to that. 

To create the OSv2 executor:

```
import (
    "github.com/opensearch-project/opensearch-go/v2" // the official OpenSearch v2 lib
    "github.com/CrowdStrike/opensearchtools"
    "github.com/CrowdStrike/opensearchtools/osv2"
)

...

// create your opensearch.Client
// opensearchV2Client := ...

// now pass your opensearch.Client to the executor constructor:
osv2Executor := osv2.NewExecutor(opensearchV2Client)
```

Say you have a slice of objects called `objectsToWrite` you would like to index. Each object must implement the [RoutableDoc](document.go#L7) interface and therefore have `ID()` and `Index()` methods defined.

You can than create a bulk request to index all these objects as follows:

```
bulkRequest := opensearchtools.
    NewBulkRequest().
    WithIndex("my_index")

for _, o := range objectsToWrite {
    bulkRequest.Add(opensearchtools.NewCreateBulkAction(o))
}
```

And then you can executute your bulk request using the osv2 executor:

```
resp, rErr := r.osv2Executor.Bulk(ctx, bulkRequest)
```

The `resp` is of type `opensearchtools.OpenSearchResponse[opensearchtools.BulkResponse]`. The generic `opensearchtools.OpenSearchResponse[T any]` encapsulates the actual desired result type (`opensearchtools.BulkResponse` in this case) as well as a `ValidationResult` which provides validation information about the request.

The `resp` should be validated by checking for fatal validation errors:

```
if resp.ValidationResults.IsFatal() {
    // create error from the validation results:
    return nil, opensearchtools.NewValidationError(resp.ValidationResults)
}
```

If there are no fatal validation results, one should still look to see if there are any non-fatal validation results.




## Examples

These examples are relatively short and to the point and skip over the creation of the executor and the validation in order to
emphasize the specics for each type of functionality. For a more hand-holding approach, see the Detailed Walkthrough section.

### Bulk

#### Bulk Index

```
// TODO
```

#### Bulk Create

```
bulkRequest := opensearchtools.
    NewBulkRequest().
    WithIndex("my_index")

for _, o := range objectsToWrite {
    bulkRequest.Add(opensearchtools.NewCreateBulkAction(o))
}
```

#### Bulk Update

```
// TODO
```

#### Bulk Delete

```
// TODO
```

### Search 

```
// TODO
```


### Multi-Get

```
ids := []id{"100", "101", "102"}

req := opensearchtools.
    NewMGetRequest().
    WithIndex("my_index")


type MyDocType struct {
    id string
    index string
    // your other fields
}

func (d *MyDocType) ID() string {
    return d.id
}

func (d *MyDocType) Index() string {
    return d.index
}

for _, id := range ids {
    nextObj := &MyDocType{
        id: id,
    }
    req.AddDocs(nextObj)
}

resp, rErr := r.osv2Executor.MGet(ctx, req)
```

## License

opensearchtools is licensed under the [MIT License](https://opensource.org/licenses/MIT).
