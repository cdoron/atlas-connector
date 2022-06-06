/*
 * Data Catalog Service - Asset Details
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	b64 "encoding/base64"

	api "github.com/cdoron/datacatalog-go/api"
	"github.com/go-resty/resty/v2"
)

// DefaultApiService is a service that implements the logic for the DefaultApiServicer
// This service should implement the business logic for every endpoint for the DefaultApi API.
// Include any external packages or services that will be required by this service.
type ApacheApiService struct {
	hostname string
	port     string
	username string
	password string
}

// NewDefaultApiService creates a default api service
func NewApacheApiService(conf map[interface{}]interface{}) AtlasApiServicer {
	return &ApacheApiService{conf["atlas_hostname"].(string),
		strconv.Itoa(conf["atlas_port"].(int)),
		conf["atlas_username"].(string),
		conf["atlas_password"].(string)}
}

func extract_asset_id_from_body(body []byte) (assetId string, err error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Malformed response from Apache Atlas")
			err = r.(error)
		}
	}()

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	assetId = result["mutatedEntities"].(map[string]interface{})["CREATE"].([]interface{})[0].(map[string]interface{})["guid"].(string)

	return assetId, err
}

func extract_metadata_from_body(body []byte) (metadata string, deleted bool, err error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Malformed response from Apache Atlas")
			err = r.(error)
		}
	}()

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	metadata = result["entity"].(map[string]interface{})["customAttributes"].(map[string]interface{})["metadata"].(string)
	deleted = result["entity"].(map[string]interface{})["status"].(string) == "DELETED"

	return metadata, deleted, err
}

// CreateAsset - This REST API writes data asset information to the data catalog configured in fybrik
func (s *ApacheApiService) CreateAsset(ctx context.Context,
	xRequestDatacatalogWriteCred string,
	createAssetRequest api.CreateAssetRequest,
	bodyBytes []byte) (api.ImplResponse, error) {

	//TODO: Uncomment the next line to return response Response(201, CreateAssetResponse{}) or use other options such as http.Ok ...
	//return Response(201, CreateAssetResponse{}), nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	assetName := createAssetRequest.DestinationCatalogID + "/" + createAssetRequest.DestinationAssetID
	metadata := b64.StdEncoding.EncodeToString(bodyBytes)

	body := `
	{
	  "entity": {
		  "typeName": "Asset",
		  "attributes": {
			  "qualifiedName": "` + assetName + `",
			  "name": "` + assetName + `"
		  },
		  "customAttributes": {
			  "metadata": "` + metadata + `"
		  }
	  }
  }
	`

	client := resty.New()

	resp, err := client.R().
		SetBasicAuth(s.username, s.password).
		SetBody(body).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		Post("http://" + s.hostname + ":" + s.port + "/api/atlas/v2/entity")

	if err != nil {
		return api.Response(500, nil), err
	}

	if resp.StatusCode() != 200 {
		return api.Response(resp.StatusCode(), errors.New("Got "+strconv.Itoa(resp.StatusCode())+" from Atlas server")), nil
	}

	assetID, err := extract_asset_id_from_body(resp.Body())
	if err != nil {
		return api.Response(400, nil), err
	}

	return api.Response(200, api.CreateAssetResponse{assetID}), nil
}

// DeleteAsset - This REST API deletes data asset
func (s *ApacheApiService) DeleteAsset(ctx context.Context, xRequestDatacatalogCred string, deleteAssetRequest api.DeleteAssetRequest) (api.ImplResponse, error) {
	assetID := deleteAssetRequest.AssetID

	client := resty.New()

	resp, err := client.R().
		SetBasicAuth(s.username, s.password).
		Delete("http://" + s.hostname + ":" + s.port + "/api/atlas/v2/entity/guid/" + assetID)

	if err != nil {
		return api.Response(500, nil), err
	}

	if resp.StatusCode() != 200 {
		return api.Response(resp.StatusCode(), errors.New("Got "+strconv.Itoa(resp.StatusCode())+" from Atlas server")), nil
	}

	return api.Response(200, api.DeleteAssetResponse{"Deletion Successful"}), nil
}

// GetAssetInfo - This REST API gets data asset information from the data catalog configured in fybrik for the data sets indicated in FybrikApplication yaml
func (s *ApacheApiService) GetAssetInfo(ctx context.Context, xRequestDatacatalogCred string, getAssetRequest api.GetAssetRequest) (api.ImplResponse, error) {
	assetID := getAssetRequest.AssetID

	client := resty.New()
	resp, err := client.R().
		SetBasicAuth(s.username, s.password).
		Get("http://" + s.hostname + ":" + s.port + "/api/atlas/v2/entity/guid/" + assetID)

	if err != nil {
		return api.Response(500, nil), err
	}

	if resp.StatusCode() != 200 {
		return api.Response(resp.StatusCode(), errors.New("Got "+strconv.Itoa(resp.StatusCode())+" from Atlas server")), nil
	}

	metadata, deleted, err := extract_metadata_from_body(resp.Body())
	if err != nil {
		return api.Response(400, nil), err
	}

	if deleted {
		return api.Response(404, "Asset already deleted"), nil
	}

	assetInfo, err := b64.StdEncoding.DecodeString(metadata)
	if err != nil {
		return api.Response(400, nil), err
	}

	var result map[string]interface{}
	json.Unmarshal(assetInfo, &result)

	return api.Response(200, result), nil
}

// UpdateAsset - This REST API updates data asset information in the data catalog configured in fybrik
func (s *ApacheApiService) UpdateAsset(ctx context.Context, xRequestDatacatalogUpdateCred string, updateAssetRequest api.UpdateAssetRequest) (api.ImplResponse, error) {
	// TODO - update UpdateAsset with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, UpdateAssetResponse{}) or use other options such as http.Ok ...
	//return Response(200, UpdateAssetResponse{}), nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	//TODO: Uncomment the next line to return response Response(401, {}) or use other options such as http.Ok ...
	//return Response(401, nil),nil

	return api.Response(http.StatusNotImplemented, nil), errors.New("UpdateAsset method not implemented")
}
