package listsort

import (
	"errors"
	"testing"
)

func Test_ParseSortParams(t *testing.T) {
	type args struct {
		rawSort  string
		strType  string
		filter   []string
		expected *SortObj
	}

	tests := []args{
		{
			rawSort: "name_asc,age_desc",
			strType: "_",
			filter:  []string{"name", "age"},
			expected: &SortObj{
				RawSort: "name_asc,age_desc",
				SortParams: SortParams{
					&SortParam{Field: "name", Order: "asc"},
					&SortParam{Field: "age", Order: "desc"},
				},
				Err: nil,
			},
		},
		{
			rawSort: "[{\"field\":\"name\",\"order\":\"asc\"},{\"field\":\"age\",\"order\":\"desc\"}]",
			strType: "json",
			filter:  []string{"name", "age"},
			expected: &SortObj{
				RawSort: "[{\"field\":\"name\",\"order\":\"asc\"},{\"field\":\"age\",\"order\":\"desc\"}]",
				SortParams: SortParams{
					&SortParam{Field: "name", Order: "asc"},
					&SortParam{Field: "age", Order: "desc"},
				},
				Err: nil,
			},
		},
		{
			rawSort: "name_asc1,age_desc",
			strType: "_",
			expected: &SortObj{
				RawSort: "name_asc1,age_desc",
				Err:     errors.New("invalid sort params"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.rawSort, func(t *testing.T) {
			got := ParseSortParams(tt.rawSort, tt.strType, tt.filter...)
			if got.RawSort != tt.expected.RawSort {
				t.Errorf("ParseSortParams() = %v, want %v", got.RawSort, tt.expected.RawSort)
			}
			if tt.expected.Err != nil {
				if got.Err == nil {
					t.Errorf("ParseSortParams() = %v, want %v", got.Err, tt.expected.Err)
				}

				t.SkipNow()
			}

			if got.Err != tt.expected.Err {
				t.Errorf("ParseSortParams() = %v, want %v", got.Err, tt.expected.Err)
			}

			if tt.expected.Err == nil && got.IsError() {
				t.Errorf("ParseSortParams() = %v, want %v, IsError should be false", got.Err, tt.expected.Err)
			}
			if tt.expected.Err != nil && !got.IsError() {
				t.Errorf("ParseSortParams() = %v, want %v, IsError should be true", got.Err, tt.expected.Err)
			}

			if len(got.SortParams) != len(tt.expected.SortParams) {
				t.Errorf("ParseSortParams() = %v, want %v", got.SortParams, tt.expected.SortParams)
			}

			if len(got.SortParams) <= 0 && got.ShouldSort() {
				t.Errorf("SortParams lower than 0 should return false for ShouldSort()")
			}
			if len(got.SortParams) > 0 && !got.ShouldSort() {
				t.Errorf("SortParams greater than 0 should return true for ShouldSort()")
			}

			for i, v := range got.SortParams {
				if v.Field != tt.expected.SortParams[i].Field {
					t.Errorf("ParseSortParams() = %v, want %v", v.Field, tt.expected.SortParams[i].Field)
				} else if !got.Exists(v.Field) {
					t.Errorf("ParseSortParams slice has %v, Exists() should be true", v.Field)
				}
			}
		})
	}

}
