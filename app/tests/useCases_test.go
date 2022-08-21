package tests

import (
	"errors"
	"github.com/AlfonsoCebollero/asyncReply-RequestAPI/app/server/entities"
	"github.com/AlfonsoCebollero/asyncReply-RequestAPI/app/useCases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type CadenceManagerMock struct {
	mock.Mock
}

func (m *CadenceManagerMock) StartWorkflow(name string, list string, duration int) (string, string, error) {
	args := m.Called(name, list, duration)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *CadenceManagerMock) RetrieveWFStatus(id string) (int, bool, error) {
	args := m.Called(id)
	return args.Int(0), args.Bool(1), args.Error(2)
}

func TestCreateWorkflowSuccess(t *testing.T) {
	mockClient := new(CadenceManagerMock)
	interval := "35"
	mockClient.
		On(
			"StartWorkflow",
			"WaitingWorkflow", "asyncAPI", 35).
		Return("inventedID",
			"inventedRunID", nil)

	id, runId, err := useCases.CreateWorkflow(mockClient, "", interval, "waitingwf")

	if err != nil {
		t.Fail()
		return
	}

	mockClient.AssertExpectations(t)
	mockClient.AssertNumberOfCalls(t, "StartWorkflow", 1)

	assert.Equal(t, "inventedID", id)
	assert.Equal(t, "inventedRunID", runId)

}

func TestBadWorkflow(t *testing.T) {
	mockClient := new(CadenceManagerMock)
	interval := "35"

	id, runId, err := useCases.CreateWorkflow(mockClient, "", interval, "invalidWF")

	if err == nil {
		t.Fail()
	}
	if err != entities.BadWorkflow {
		t.Fail()
	}

	mockClient.AssertExpectations(t)
	mockClient.AssertNumberOfCalls(t, "StartWorkflow", 0)

	assert.Equal(t, "", id)
	assert.Equal(t, "", runId)
}

func TestValidationError(t *testing.T) {
	mockClient := new(CadenceManagerMock)
	interval := "22"

	id, runId, err := useCases.CreateWorkflow(mockClient, "", interval, "invalidWF")

	if err == nil {
		t.Fail()
	}
	if err != entities.ValidationError {
		t.Fail()
	}

	mockClient.AssertExpectations(t)
	mockClient.AssertNumberOfCalls(t, "StartWorkflow", 0)

	assert.Equal(t, "", id)
	assert.Equal(t, "", runId)
}

func TestRetrieveWorkflowStatusCompleted(t *testing.T) {
	mockClient := new(CadenceManagerMock)
	mockClient.On("RetrieveWFStatus", "inventedID").Return(0, false, nil)
	status, err := useCases.RetrieveWorkflowStatus(mockClient, "inventedID")

	if err != nil {
		t.Fail()
	}

	mockClient.AssertExpectations(t)
	assert.Equal(t, "COMPLETED", status)
}

func TestRetrieveWorkflowStatusPending(t *testing.T) {
	mockClient := new(CadenceManagerMock)
	mockClient.On("RetrieveWFStatus", "inventedID").Return(-1, true, nil)
	status, err := useCases.RetrieveWorkflowStatus(mockClient, "inventedID")

	if err != nil {
		t.Fail()
	}

	mockClient.AssertExpectations(t)
	assert.Equal(t, "PENDING", status)
}

func TestRetrieveWorkflowStatusRunning(t *testing.T) {
	mockClient := new(CadenceManagerMock)
	mockClient.On("RetrieveWFStatus", "inventedID").Return(-1, false, nil)
	status, err := useCases.RetrieveWorkflowStatus(mockClient, "inventedID")

	if err != nil {
		t.Fail()
	}

	mockClient.AssertExpectations(t)
	assert.Equal(t, "RUNNING", status)
}

func TestRetrieveWorkflowStatusError(t *testing.T) {
	mockClient := new(CadenceManagerMock)
	mockClient.On("RetrieveWFStatus", "inventedID").Return(0, false, errors.New(""))
	status, err := useCases.RetrieveWorkflowStatus(mockClient, "inventedID")

	if err == nil {
		t.Fail()
	}

	assert.Equal(t, "", status)
	mockClient.AssertExpectations(t)
}
