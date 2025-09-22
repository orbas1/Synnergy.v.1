package tests

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func repoPath(parts ...string) string {
	return filepath.Join(append([]string{".."}, parts...)...)
}

func decodeYAMLDocuments(t *testing.T, path string) []map[string]any {
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	dec := yaml.NewDecoder(bytes.NewReader(data))
	var docs []map[string]any
	for {
		var doc map[string]any
		if err := dec.Decode(&doc); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Fatalf("decode %s: %v", path, err)
		}
		if len(doc) == 0 {
			continue
		}
		docs = append(docs, doc)
	}
	return docs
}

func findDocument(docs []map[string]any, kind string) map[string]any {
	for _, doc := range docs {
		if k, _ := doc["kind"].(string); k == kind {
			return doc
		}
	}
	return nil
}

func TestStage77KubernetesNodeManifest(t *testing.T) {
	docs := decodeYAMLDocuments(t, repoPath("deploy", "k8s", "node.yaml"))
	if ns := findDocument(docs, "Namespace"); ns == nil {
		t.Fatalf("expected Namespace document")
	}
	deploy := findDocument(docs, "Deployment")
	if deploy == nil {
		t.Fatalf("expected Deployment document")
	}
	spec, ok := deploy["spec"].(map[string]any)
	if !ok {
		t.Fatalf("deployment spec missing")
	}
	tmpl := spec["template"].(map[string]any)
	tmplSpec := tmpl["spec"].(map[string]any)
	containers := tmplSpec["containers"].([]any)
	var found bool
	for _, c := range containers {
		cm := c.(map[string]any)
		if cm["name"] == "synnergy-node" {
			found = true
			if _, ok := cm["readinessProbe"]; !ok {
				t.Fatalf("synnergy-node container missing readinessProbe")
			}
			if _, ok := cm["livenessProbe"]; !ok {
				t.Fatalf("synnergy-node container missing livenessProbe")
			}
		}
	}
	if !found {
		t.Fatalf("synnergy-node container not found")
	}
	if pdb := findDocument(docs, "PodDisruptionBudget"); pdb == nil {
		t.Fatalf("expected PodDisruptionBudget document")
	}
}

func TestStage77KubernetesWalletManifest(t *testing.T) {
	docs := decodeYAMLDocuments(t, repoPath("deploy", "k8s", "wallet.yaml"))
	if cfg := findDocument(docs, "ConfigMap"); cfg == nil {
		t.Fatalf("expected ConfigMap document")
	}
	if secret := findDocument(docs, "Secret"); secret == nil {
		t.Fatalf("expected Secret document")
	}
	deploy := findDocument(docs, "Deployment")
	if deploy == nil {
		t.Fatalf("expected Deployment document")
	}
	spec := deploy["spec"].(map[string]any)
	tmpl := spec["template"].(map[string]any)
	tmplSpec := tmpl["spec"].(map[string]any)
	containers := tmplSpec["containers"].([]any)
	if len(containers) == 0 || containers[0].(map[string]any)["name"] != "wallet" {
		t.Fatalf("wallet container not found")
	}
	if hpa := findDocument(docs, "HorizontalPodAutoscaler"); hpa == nil {
		t.Fatalf("expected HorizontalPodAutoscaler document")
	}
	if np := findDocument(docs, "NetworkPolicy"); np == nil {
		t.Fatalf("expected NetworkPolicy document")
	}
}

func TestStage77DockerComposeStack(t *testing.T) {
	data, err := os.ReadFile(repoPath("docker", "docker-compose.yml"))
	if err != nil {
		t.Fatalf("read compose: %v", err)
	}
	if !strings.Contains(string(data), "web:") {
		t.Fatalf("expected web service in compose stack")
	}
	if !strings.Contains(string(data), "healthcheck") {
		t.Fatalf("expected health checks in compose stack")
	}
}

func TestStage77TerraformBlueprint(t *testing.T) {
	data, err := os.ReadFile(repoPath("deploy", "terraform", "main.tf"))
	if err != nil {
		t.Fatalf("read terraform: %v", err)
	}
	content := string(data)
	for _, required := range []string{"backend \"s3\"", "aws_lb", "aws_autoscaling_group", "aws_db_instance", "aws_ssm_parameter"} {
		if !strings.Contains(content, required) {
			t.Fatalf("expected %s in terraform blueprint", required)
		}
	}
}
