package cmd

import (
	"errors"
	"testing"
)

// AWSClientMock はAWSClientのモック実装
type AWSClientMock struct {
	InstanceList            []string
	SlowQueryList           []string
	DownloadSlowQueryResult *string
	DownloadError           error
}

func (m AWSClientMock) GetInstanceList() []string {
	return m.InstanceList
}

func (m AWSClientMock) GetSlowQueryList(instance string) ([]string, error) {
	return m.SlowQueryList, nil
}

func (m AWSClientMock) DownloadSlowQueryLog(instance string, logFile string) (*string, error) {
	return m.DownloadSlowQueryResult, m.DownloadError
}

func TestFilterInstance(t *testing.T) {
	testCases := []struct {
		name          string
		instanceList  []string
		targetPrefix  string
		expectedMatch string
	}{
		{
			name:          "完全一致",
			instanceList:  []string{"test-instance", "other-instance"},
			targetPrefix:  "test-instance",
			expectedMatch: "test-instance",
		},
		{
			name:          "前方一致",
			instanceList:  []string{"test-instance", "other-instance"},
			targetPrefix:  "test",
			expectedMatch: "test-instance",
		},
		{
			name:          "一致なし",
			instanceList:  []string{"test-instance", "other-instance"},
			targetPrefix:  "non-existent",
			expectedMatch: "",
		},
		{
			name:          "空リスト",
			instanceList:  []string{},
			targetPrefix:  "test",
			expectedMatch: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := AWSClientMock{
				InstanceList: tc.instanceList,
			}

			result := FilterInstance(mockClient, tc.targetPrefix)

			if result != tc.expectedMatch {
				t.Errorf("FilterInstance() = %v, want %v", result, tc.expectedMatch)
			}
		})
	}
}

func TestGetSlowQueryList(t *testing.T) {
	testCases := []struct {
		name           string
		instanceName   string
		expectedLogs   []string
		expectedLength int
	}{
		{
			name:           "ログファイルあり",
			instanceName:   "test-instance",
			expectedLogs:   []string{"log1.log", "log2.log"},
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
			mockClient := AWSClientMock{
				SlowQueryList: tc.expectedLogs,
			}

			logs, err := GetSlowQueryList(mockClient, tc.instanceName)

			if err != nil {
				t.Errorf("GetSlowQueryList() returned error: %v", err)
			}

			if len(logs) != tc.expectedLength {
				t.Errorf("GetSlowQueryList() returned %d logs, want %d", len(logs), tc.expectedLength)
			}

			for i, log := range logs {
				if log != tc.expectedLogs[i] {
					t.Errorf("GetSlowQueryList()[%d] = %v, want %v", i, log, tc.expectedLogs[i])
				}
			}
		})
	}
}

func TestDownloadSlowQueryLog(t *testing.T) {
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
			instance:      "test-instance",
			logFiles:      []string{"slowquery/mysql-slowquery.log.1"},
			filter:        "slowquery",
			downloadValue: "# Time: 2023-01-01\nSELECT 1",
			downloadError: nil,
			expectedError: false,
		},
		{
			name:          "ダウンロードエラー",
			instance:      "test-instance",
			logFiles:      []string{"slowquery/mysql-slowquery.log.1"},
			filter:        "slowquery",
			downloadValue: "",
			downloadError: errors.New("download error"),
			expectedError: true,
		},
		{
			name:          "フィルター一致なし",
			instance:      "test-instance",
			logFiles:      []string{"general/mysql-general.log.1"},
			filter:        "slowquery",
			downloadValue: "",
			downloadError: nil,
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logValue := tc.downloadValue
			mockClient := AWSClientMock{
				DownloadSlowQueryResult: &logValue,
				DownloadError:           tc.downloadError,
			}

			_, err := DownloadSlowQueryLog(mockClient, tc.instance, tc.filter, tc.logFiles)

			if (err != nil) != tc.expectedError {
				t.Errorf("DownloadSlowQueryLog() error = %v, expectedError %v", err, tc.expectedError)
			}
		})
	}
}
