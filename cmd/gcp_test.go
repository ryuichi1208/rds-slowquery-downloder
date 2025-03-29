package cmd

import (
	"errors"
	"testing"
)

// GCPClientMock はGCPClientのモック実装
type GCPClientMock struct {
	InstanceList            []string
	SlowQueryList           []string
	DownloadSlowQueryResult *string
	DownloadError           error
}

func (m GCPClientMock) GetInstanceList() []string {
	return m.InstanceList
}

func (m GCPClientMock) GetSlowQueryList(instance string) ([]string, error) {
	return m.SlowQueryList, nil
}

func (m GCPClientMock) DownloadSlowQueryLog(instance string, logFile string) (*string, error) {
	return m.DownloadSlowQueryResult, m.DownloadError
}

func TestGCPFilterInstance(t *testing.T) {
	testCases := []struct {
		name          string
		instanceList  []string
		targetPrefix  string
		expectedMatch string
	}{
		{
			name:          "完全一致",
			instanceList:  []string{"gcp-instance", "other-instance"},
			targetPrefix:  "gcp-instance",
			expectedMatch: "gcp-instance",
		},
		{
			name:          "前方一致",
			instanceList:  []string{"gcp-instance", "other-instance"},
			targetPrefix:  "gcp",
			expectedMatch: "gcp-instance",
		},
		{
			name:          "一致なし",
			instanceList:  []string{"gcp-instance", "other-instance"},
			targetPrefix:  "non-existent",
			expectedMatch: "",
		},
		{
			name:          "空リスト",
			instanceList:  []string{},
			targetPrefix:  "gcp",
			expectedMatch: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := GCPClientMock{
				InstanceList: tc.instanceList,
			}

			result := FilterInstance(mockClient, tc.targetPrefix)

			if result != tc.expectedMatch {
				t.Errorf("FilterInstance() with GCP client = %v, want %v", result, tc.expectedMatch)
			}
		})
	}
}

func TestGCPGetSlowQueryList(t *testing.T) {
	testCases := []struct {
		name           string
		instanceName   string
		expectedLogs   []string
		expectedLength int
	}{
		{
			name:           "ログファイルあり",
			instanceName:   "gcp-instance",
			expectedLogs:   []string{"slowquery/mysql-slowquery.log.gcp-instance.1", "slowquery/mysql-slowquery.log.gcp-instance.2"},
			expectedLength: 2,
		},
		{
			name:           "ログファイルなし",
			instanceName:   "empty-instance",
			expectedLogs:   []string{},
			expectedLength: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := GCPClientMock{
				SlowQueryList: tc.expectedLogs,
			}

			logs, err := GetSlowQueryList(mockClient, tc.instanceName)

			if err != nil {
				t.Errorf("GetSlowQueryList() with GCP client returned error: %v", err)
			}

			if len(logs) != tc.expectedLength {
				t.Errorf("GetSlowQueryList() with GCP client returned %d logs, want %d", len(logs), tc.expectedLength)
			}

			for i, log := range logs {
				if log != tc.expectedLogs[i] {
					t.Errorf("GetSlowQueryList()[%d] with GCP client = %v, want %v", i, log, tc.expectedLogs[i])
				}
			}
		})
	}
}

func TestGCPDownloadSlowQueryLog(t *testing.T) {
	testCases := []struct {
		name          string
		instance      string
		logFiles      []string
		filter        string
		downloadValue string
		downloadError error
		expectedError bool
	}{
		{
			name:          "正常系",
			instance:      "gcp-instance",
			logFiles:      []string{"slowquery/mysql-slowquery.log.gcp-instance.1"},
			filter:        "slowquery",
			downloadValue: "# Time: 2023-01-01\nSELECT 1",
			downloadError: nil,
			expectedError: false,
		},
		{
			name:          "ダウンロードエラー",
			instance:      "gcp-instance",
			logFiles:      []string{"slowquery/mysql-slowquery.log.gcp-instance.1"},
			filter:        "slowquery",
			downloadValue: "",
			downloadError: errors.New("download error"),
			expectedError: true,
		},
		{
			name:          "フィルター一致なし",
			instance:      "gcp-instance",
			logFiles:      []string{"general/mysql-general.log.gcp-instance.1"},
			filter:        "slowquery",
			downloadValue: "",
			downloadError: nil,
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logValue := tc.downloadValue
			mockClient := GCPClientMock{
				DownloadSlowQueryResult: &logValue,
				DownloadError:           tc.downloadError,
			}

			_, err := DownloadSlowQueryLog(mockClient, tc.instance, tc.filter, tc.logFiles)

			if (err != nil) != tc.expectedError {
				t.Errorf("DownloadSlowQueryLog() with GCP client error = %v, expectedError %v", err, tc.expectedError)
			}
		})
	}
}
