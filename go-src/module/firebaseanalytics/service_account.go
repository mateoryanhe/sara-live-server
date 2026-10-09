package firebaseanalytics

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type serviceAccountCredentials struct {
	Type        string `json:"type"`
	ProjectId   string `json:"project_id"`
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
}

func validateServiceAccountJSON(ctx context.Context, projectId, serviceAccountJson string) error {
	projectId = strings.TrimSpace(projectId)
	serviceAccountJson = strings.TrimSpace(serviceAccountJson)
	if projectId == "" || serviceAccountJson == "" {
		return fmt.Errorf("project id and service account json are required")
	}
	var cred serviceAccountCredentials
	if err := json.Unmarshal([]byte(serviceAccountJson), &cred); err != nil {
		return fmt.Errorf("parse service account json: %w", err)
	}
	if cred.Type != "service_account" ||
		strings.TrimSpace(cred.ProjectId) == "" ||
		strings.TrimSpace(cred.ClientEmail) == "" ||
		strings.TrimSpace(cred.PrivateKey) == "" {
		return fmt.Errorf("service account json is incomplete")
	}
	if strings.TrimSpace(cred.ProjectId) != projectId {
		return fmt.Errorf("project id does not match service account")
	}
	if _, err := newFirebaseAdminApp(ctx, projectId, serviceAccountJson); err != nil {
		return err
	}
	return nil
}
