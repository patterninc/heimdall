package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/glue"
)

const (
	tableMetadataLocationKey = `metadata_location`
)

var (
	ErrNoTableName                 = fmt.Errorf(`no table and / or namespace name provided`)
	ErrMissingCatalogTableMetadata = fmt.Errorf(`missing table metadata in the glue catalog`)
)

func GetTableMetadata(catalogID, tableName string) ([]byte, error) {

	// split tableName to namespace and table names
	tableNameParts := strings.Split(tableName, `.`)

	if len(tableNameParts) != 2 {
		return nil, fmt.Errorf("unexpected table name: [%s]", tableName)
	}

	// let's get the latest metadata file location
	location, err := getTableMetadataLocation(catalogID, tableNameParts[0], tableNameParts[1])
	if err != nil {
		return nil, err
	}

	// let's pull the file content
	return ReadFromS3(location)

}

// function that calls AWS glue catalog to get the snapshot ID for a given database, table and branch
func getTableMetadataLocation(catalogID, databaseName, tableName string) (string, error) {

	// Return an error if databaseName or tableName is empty
	if databaseName == `` || tableName == `` {
		return ``, ErrNoTableName
	}

	// Create a Glue client using the provided awsConfig
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return ``, err
	}

	// Create an S3 client
	svc := glue.NewFromConfig(awsConfig)

	// Call the GetTable API to retrieve the table metadata
	getObjectOutput, err := svc.GetTable(ctx, &glue.GetTableInput{
		CatalogId:    &catalogID,
		DatabaseName: &databaseName,
		Name:         &tableName,
	})
	if err != nil {
		return ``, err
	}

	if getObjectOutput == nil || getObjectOutput.Table == nil {
		return ``, ErrMissingCatalogTableMetadata
	}

	metdadataLocation, found := getObjectOutput.Table.Parameters[tableMetadataLocationKey]
	if !found || metdadataLocation == `` {
		return ``, ErrMissingCatalogTableMetadata
	}

	return metdadataLocation, nil

}
