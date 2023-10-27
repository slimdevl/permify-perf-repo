package permify

import (
	"context"
	"encoding/json"
	"net/http"
)

func (c *client) SaveModelSchema(ctx context.Context, schema *SaveSchemaRequest) (*SaveSchemaResponse, error) {
	url := c.constructURL(SchemaWriteAPIPath)

	body, err := c.sendRequest(ctx, http.MethodPost, url, schema)
	if err != nil {
		return nil, err
	}

	var response SaveSchemaResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, ErrBodyDecodeFailure
	}

	if response.ErrorResponse != nil && response.ErrorResponse.Code != 0 {
		return nil, ErrUnableToWriteSchema
	}

	return &response, nil
}
