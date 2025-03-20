package provider

import (
	"testing"
)

func TestGetAWSAccountID(t *testing.T) {
	provider := NewProvider("us-east-1")
	accountID, err := provider.GetAWSAccountID()
	if err != nil {
		t.Fatalf("Failed to get AWS account ID: %v", err)
	}
	t.Logf("AWS account ID: %s", accountID)
}

func TestGetECRPassword(t *testing.T) {
	provider := NewProvider("us-east-1")
	password, err := provider.GetECRPassword()
	if err != nil {
		t.Fatalf("Failed to get ECR password: %v", err)
	}
	t.Logf("ECR password: %s", password)
}
