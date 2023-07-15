package redis

import (
	"context"
	"errors"
	"testing"

	"github.com/dnevsky/http-products/models"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestProduct_GetWithOffsetFromJSON(t *testing.T) {
	db, mock := redismock.NewClientMock()

	testCases := []struct {
		name        string
		limit       int
		offset      int
		redisResult []string
		expected    []models.Product
		expectedErr error
	}{
		{
			name:   "Sucefull",
			limit:  2,
			offset: 5,
			redisResult: []string{
				`{"id":"product000005","price":228}`,
				`{"id":"product000006","price":1337}`,
			},
			expected: []models.Product{
				{Id: "product000005", Price: 228},
				{Id: "product000006", Price: 1337},
			},
		},
		{
			name:   "Json decode error",
			limit:  2,
			offset: 5,
			redisResult: []string{
				`invalid json`,
			},
			expectedErr: errors.New("invalid character 'i' looking for beginning of value"),
		},
		{
			name:        "Empty result",
			limit:       2,
			offset:      5,
			redisResult: []string{},
			expected:    []models.Product{},
		},
		{
			name:        "negative offset",
			limit:       2,
			offset:      -5,
			expectedErr: errors.New("invalid offset or limit"),
		},
		{
			name:        "zero limit",
			limit:       0,
			offset:      5,
			expectedErr: errors.New("invalid offset or limit"),
		},
		{
			name:        "negative limit",
			limit:       -2,
			offset:      5,
			expectedErr: errors.New("invalid offset or limit"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			mock.ExpectLRange("products", int64(tc.offset), int64(tc.offset+tc.limit-1)).SetVal(tc.redisResult)

			logger := zap.NewExample().Sugar()

			c := &ProductRedis{logger: logger, client: db}

			products, err := c.GetWithOffsetFromJSON(context.TODO(), tc.offset, tc.limit)

			assert.Equal(t, tc.expected, products)

			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.Nil(t, err)
			}

		})
	}
}

func TestProduct_UpdateCache(t *testing.T) {
	db, mock := redismock.NewClientMock()

	testCase := []struct {
		name        string
		data        []string
		expectedErr error
		setup       func(data []string, expectedErr error)
	}{
		{
			name: "Sucefull",
			data: []string{
				`{"id":"product000005","price":228}`,
				`{"id":"product000006","price":1337}`,
			},
			expectedErr: nil,
			setup: func(data []string, expectedErr error) {
				mock.ExpectTxPipeline()
				mock.ExpectDel("products").SetVal(1)
				mock.ExpectRPush("products", data).SetVal(2)
				mock.ExpectTxPipelineExec()
			},
		},
		{
			name:        "Given zero slice",
			data:        []string{},
			expectedErr: errors.New("UpdateData: there's nothing to stuff"),
			setup: func(data []string, expectedErr error) {
				// nothing to setup
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.data, tc.expectedErr)

			logger := zap.NewExample().Sugar()

			c := &ProductRedis{logger: logger, client: db}

			err := c.UpdateData(context.TODO(), tc.data)

			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}

}
