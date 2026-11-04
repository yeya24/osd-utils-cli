package getcpautoscalingstatus

import (
	"testing"
	"time"
)

func TestApplyFilter(t *testing.T) {
	timestamp := time.Now()

	testClusters := []clusterInfo{
		{
			ClusterID:             "cluster-001",
			ClusterName:           "test-cluster-001",
			Namespace:             "ocm-production-001",
			AutoscalingEnabled:    false,
			HasOverrideAnnotation: true,
			CurrentSize:           "medium",
			RecommendedSize:       "medium",
		},
		{
			ClusterID:             "cluster-002",
			ClusterName:           "test-cluster-002",
			Namespace:             "ocm-production-002",
			AutoscalingEnabled:    true,
			HasOverrideAnnotation: true,
			CurrentSize:           "large",
			RecommendedSize:       "large",
		},
		{
			ClusterID:             "cluster-003",
			ClusterName:           "test-cluster-003",
			Namespace:             "ocm-production-003",
			AutoscalingEnabled:    true,
			HasOverrideAnnotation: false,
			CurrentSize:           "small",
			RecommendedSize:       "small",
		},
		{
			ClusterID:             "cluster-004",
			ClusterName:           "test-cluster-004",
			Namespace:             "ocm-production-004",
			AutoscalingEnabled:    false,
			HasOverrideAnnotation: false,
			CurrentSize:           "medium",
			RecommendedSize:       "N/A",
		},
		{
			ClusterID:             "cluster-005",
			ClusterName:           "test-cluster-005",
			Namespace:             "ocm-production-005",
			AutoscalingEnabled:    true,
			HasOverrideAnnotation: true,
			CurrentSize:           "large",
			RecommendedSize:       "medium",
		},
		{
			ClusterID:             "cluster-006",
			ClusterName:           "test-cluster-006",
			Namespace:             "ocm-production-006",
			AutoscalingEnabled:    true,
			HasOverrideAnnotation: true,
			CurrentSize:           "N/A",
			RecommendedSize:       "N/A",
		},
	}

	testResults := &auditResults{
		Timestamp:         timestamp,
		ManagementCluster: "test-mc",
		TotalClusters:     6,
		Clusters:          testClusters,
	}

	tests := []struct {
		name          string
		showOnly      string
		expectedCount int
		expectedIDs   []string
	}{
		{
			name:          "filter needs-removal",
			showOnly:      "needs-removal",
			expectedCount: 4,
			expectedIDs:   []string{"cluster-001", "cluster-002", "cluster-005", "cluster-006"},
		},
		{
			name:          "filter ready-for-migration",
			showOnly:      "ready-for-migration",
			expectedCount: 2,
			expectedIDs:   []string{"cluster-001", "cluster-004"},
		},
		{
			name:          "filter safe-to-remove-override",
			showOnly:      "safe-to-remove-override",
			expectedCount: 1,
			expectedIDs:   []string{"cluster-002"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &options{
				showOnly: tt.showOnly,
			}

			filtered := opts.applyFilter(testResults)

			if len(filtered.Clusters) != tt.expectedCount {
				t.Errorf("Expected %d clusters, got %d", tt.expectedCount, len(filtered.Clusters))
			}

			if filtered.TotalClusters != tt.expectedCount {
				t.Errorf("Expected TotalClusters=%d, got %d", tt.expectedCount, filtered.TotalClusters)
			}

			foundIDs := make(map[string]bool)
			for _, cluster := range filtered.Clusters {
				foundIDs[cluster.ClusterID] = true
			}

			for _, expectedID := range tt.expectedIDs {
				if !foundIDs[expectedID] {
					t.Errorf("Expected cluster %s in filtered results but not found", expectedID)
				}
			}

			// Verify metadata is preserved
			if filtered.Timestamp != timestamp {
				t.Errorf("Timestamp not preserved in filtered results")
			}
			if filtered.ManagementCluster != testResults.ManagementCluster {
				t.Errorf("ManagementCluster not preserved in filtered results")
			}
		})
	}
}

func TestSafeToRemoveOverrideLogic(t *testing.T) {
	tests := []struct {
		name                  string
		cluster               clusterInfo
		shouldBeInFilteredSet bool
	}{
		{
			name: "safe - autoscaling enabled, has override, sizes match",
			cluster: clusterInfo{
				AutoscalingEnabled:    true,
				HasOverrideAnnotation: true,
				CurrentSize:           "medium",
				RecommendedSize:       "medium",
			},
			shouldBeInFilteredSet: true,
		},
		{
			name: "unsafe - autoscaling disabled",
			cluster: clusterInfo{
				AutoscalingEnabled:    false,
				HasOverrideAnnotation: true,
				CurrentSize:           "medium",
				RecommendedSize:       "medium",
			},
			shouldBeInFilteredSet: false,
		},
		{
			name: "unsafe - no override annotation",
			cluster: clusterInfo{
				AutoscalingEnabled:    true,
				HasOverrideAnnotation: false,
				CurrentSize:           "medium",
				RecommendedSize:       "medium",
			},
			shouldBeInFilteredSet: false,
		},
		{
			name: "unsafe - sizes don't match",
			cluster: clusterInfo{
				AutoscalingEnabled:    true,
				HasOverrideAnnotation: true,
				CurrentSize:           "large",
				RecommendedSize:       "medium",
			},
			shouldBeInFilteredSet: false,
		},
		{
			name: "unsafe - recommended size is N/A",
			cluster: clusterInfo{
				AutoscalingEnabled:    true,
				HasOverrideAnnotation: true,
				CurrentSize:           "medium",
				RecommendedSize:       "N/A",
			},
			shouldBeInFilteredSet: false,
		},
		{
			name: "unsafe - recommended size is empty",
			cluster: clusterInfo{
				AutoscalingEnabled:    true,
				HasOverrideAnnotation: true,
				CurrentSize:           "medium",
				RecommendedSize:       "",
			},
			shouldBeInFilteredSet: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cluster.ClusterID = "test-cluster"
			tt.cluster.ClusterName = "test"
			tt.cluster.Namespace = "ocm-production-test"

			results := &auditResults{
				Timestamp:         time.Now(),
				ManagementCluster: "test-mc",
				TotalClusters:     1,
				Clusters:          []clusterInfo{tt.cluster},
			}

			opts := &options{
				showOnly: "safe-to-remove-override",
			}

			filtered := opts.applyFilter(results)

			if tt.shouldBeInFilteredSet {
				if len(filtered.Clusters) != 1 {
					t.Errorf("Expected cluster to be in safe-to-remove-override set but it was filtered out")
				}
			} else {
				if len(filtered.Clusters) != 0 {
					t.Errorf("Expected cluster to be filtered out but it was included in safe-to-remove-override set")
				}
			}
		})
	}
}

func TestValidateOutput(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid text output",
			output:  "text",
			wantErr: false,
		},
		{
			name:    "valid json output",
			output:  "json",
			wantErr: false,
		},
		{
			name:    "valid yaml output",
			output:  "yaml",
			wantErr: false,
		},
		{
			name:    "valid csv output",
			output:  "csv",
			wantErr: false,
		},
		{
			name:    "invalid output format",
			output:  "xml",
			wantErr: true,
			errMsg:  "invalid output format 'xml'. Valid options: text, json, yaml, csv",
		},
		{
			name:    "invalid output format - uppercase",
			output:  "JSON",
			wantErr: true,
			errMsg:  "invalid output format 'JSON'. Valid options: text, json, yaml, csv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &options{
				output: tt.output,
			}

			validOutputs := map[string]bool{"text": true, "json": true, "yaml": true, "csv": true}
			var err error
			if !validOutputs[opts.output] {
				err = &validationError{message: "invalid output format '" + opts.output + "'. Valid options: text, json, yaml, csv"}
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if err.Error() != tt.errMsg {
					t.Errorf("validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestValidateShowOnly(t *testing.T) {
	tests := []struct {
		name     string
		showOnly string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid needs-removal filter",
			showOnly: "needs-removal",
			wantErr:  false,
		},
		{
			name:     "valid ready-for-migration filter",
			showOnly: "ready-for-migration",
			wantErr:  false,
		},
		{
			name:     "valid safe-to-remove-override filter",
			showOnly: "safe-to-remove-override",
			wantErr:  false,
		},
		{
			name:     "empty filter is valid",
			showOnly: "",
			wantErr:  false,
		},
		{
			name:     "invalid filter",
			showOnly: "invalid-filter",
			wantErr:  true,
			errMsg:   "invalid show-only filter 'invalid-filter'. Valid options: needs-removal, ready-for-migration, safe-to-remove-override",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &options{
				showOnly: tt.showOnly,
			}

			var err error
			if opts.showOnly != "" {
				validFilters := map[string]bool{"needs-removal": true, "ready-for-migration": true, "safe-to-remove-override": true}
				if !validFilters[opts.showOnly] {
					err = &validationError{message: "invalid show-only filter '" + opts.showOnly + "'. Valid options: needs-removal, ready-for-migration, safe-to-remove-override"}
				}
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if err.Error() != tt.errMsg {
					t.Errorf("validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestRequiredFlag(t *testing.T) {
	cmd := newCmdAutoscalingAudit()

	mgmtClusterFlag := cmd.Flags().Lookup("mgmt-cluster-id")
	if mgmtClusterFlag == nil {
		t.Fatal("mgmt-cluster-id flag not found")
	}

	annotations := mgmtClusterFlag.Annotations
	if annotations == nil {
		t.Fatal("mgmt-cluster-id flag should be required")
	}

	requiredAnnotation := annotations["cobra_annotation_bash_completion_one_required_flag"]
	if len(requiredAnnotation) == 0 || requiredAnnotation[0] != "true" {
		t.Error("mgmt-cluster-id flag should be marked as required")
	}
}

type validationError struct {
	message string
}

func (e *validationError) Error() string {
	return e.message
}
