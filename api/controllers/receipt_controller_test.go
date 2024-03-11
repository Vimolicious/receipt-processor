package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/vimolicious/receipt-processor/data/models"
	"github.com/vimolicious/receipt-processor/data/repositories/inmemory"
)

type receiptTestCase struct {
	ExpectedPoints int            `json:"expectedPoints"`
	Receipt        models.Receipt `json:"receipt"`
}

func makeReceiptController() *ReceiptController {
	receiptRepo := inmemory.NewInMemoryReceiptRepository()
	receiptController := NewReceiptController(receiptRepo)

	return receiptController
}

func loadTestCaseBytes(name string) ([]byte, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("Error getting caller info")
	}

	testFilename := fmt.Sprintf("%s.json", name)
	testFilepath := filepath.Join(
		filepath.Dir(currentFile), "..", "..", "test", "receipts", testFilename,
	)

	contents, err := os.ReadFile(testFilepath)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func loadReceiptTestCase(name string) (*receiptTestCase, error) {
	contents, err := loadTestCaseBytes(name)
	if err != nil {
		return nil, err
	}

	var testCase receiptTestCase
	if err := json.Unmarshal(contents, &testCase); err != nil {
		return nil, err
	}

	return &testCase, nil
}

func assertStatusCode(t *testing.T, rr *httptest.ResponseRecorder, status int) {
	if rr.Code != status {
		t.Fatalf(
			"Wrong status code: '%d' expected '%d'; body '%s'",
			rr.Code, status, strings.TrimSpace(rr.Body.String()),
		)
	}
}

/*
 * Process Receipt Tests
 */

func TestProcessReceiptHandler(t *testing.T) {
	/* Good Cases */
	passingCases := []string{"pass1", "pass2"}

	receiptController := makeReceiptController()

	for _, pc := range passingCases {
		testCase, err := loadReceiptTestCase(pc)
		if err != nil {
			t.Fatal(err)
		}

		assertOkProcessResponse(t, receiptController, testCase)
	}

	/* Bad Cases */
	failingCases := []string{
		"failBadSyntax", "failMissingFields", "failMissingItemFields",
		"failMalformattedFields", "failWrongTotal", "failInvalidValues",
		"failPriceTooLarge",
	}

	for _, pc := range failingCases {
		receiptBytes, err := loadTestCaseBytes(pc)
		if err != nil {
			t.Fatalf(
				"Couldn't load raw test case '%s'; error: %s",
				receiptBytes, err.Error(),
			)
		}

		assertBadRequestProcessResponse(t, receiptController, receiptBytes)
	}

	emptyReceipt := make([]byte, 0)
	assertBadRequestProcessResponse(t, receiptController, emptyReceipt)
}

func assertBadRequestProcessResponse(
	t *testing.T, rc *ReceiptController, b []byte,
) *httptest.ResponseRecorder {
	res := callProcessReceiptHandler(t, rc, b)

	assertStatusCode(t, res, http.StatusBadRequest)

	return res
}

func assertOkProcessResponse(
	t *testing.T, rc *ReceiptController, tc *receiptTestCase,
) *httptest.ResponseRecorder {
	testBody, err := json.Marshal(tc.Receipt)
	if err != nil {
		t.Fatal(err)
	}

	res := callProcessReceiptHandler(t, rc, testBody)

	assertStatusCode(t, res, http.StatusOK)

	return res
}

func callProcessReceiptHandler(
	t *testing.T, rc *ReceiptController, b []byte,
) *httptest.ResponseRecorder {
	req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(rc.processReceiptHandler)
	handler.ServeHTTP(rr, req)

	return rr
}

/*
 * Get Points Tests
 */

func TestGetPointsHandler(t *testing.T) {
	/* Good Cases */
	passingCases := []string{"pass1", "pass2"}

	receiptController := makeReceiptController()

	for _, pc := range passingCases {
		testCase, err := loadReceiptTestCase(pc)
		if err != nil {
			t.Fatal(err)
		}

		assertCorrectPoints(t, receiptController, testCase)
	}

	/* Bad Cases */
	invalidID := "invalid ID"
	res := callGetPointsHandler(t, receiptController, invalidID)

	assertStatusCode(t, res, http.StatusBadRequest)

	nonExistantID := "00000000-0000-0000-0000-000000000000"
	res = callGetPointsHandler(t, receiptController, nonExistantID)

	assertStatusCode(t, res, http.StatusNotFound)
}

func assertCorrectPoints(
	t *testing.T, rc *ReceiptController, tc *receiptTestCase,
) *httptest.ResponseRecorder {
	res := assertOkProcessResponse(t, rc, tc)

	var processResponse processReceiptResponse
	if err := json.Unmarshal(res.Body.Bytes(), &processResponse); err != nil {
		t.Fatalf("Couldn't unmarshal process receipt response: '%s'", err.Error())
	}

	res = callGetPointsHandler(t, rc, processResponse.Id)

	assertStatusCode(t, res, http.StatusOK)

	var pointsResponse getPointsResponse
	if err := json.Unmarshal(res.Body.Bytes(), &pointsResponse); err != nil {
		t.Fatalf("Couldn't unmarshal process receipt response: '%s'", err.Error())
	}

	if tc.ExpectedPoints != pointsResponse.Points {
		t.Fatalf(
			"Wrong number of points: '%d' expected '%d'",
			pointsResponse.Points, tc.ExpectedPoints,
		)
	}

	return res
}

func callGetPointsHandler(
	t *testing.T, rc *ReceiptController, id string,
) *httptest.ResponseRecorder {
	req, err := http.NewRequest(
		"GET", fmt.Sprintf("/receipts/%s/points", id), nil,
	)
	req.SetPathValue("id", id)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(rc.getPointsHandler)
	handler.ServeHTTP(rr, req)

	return rr
}
