package dbtag

import (
	"reflect"
	"testing"
	"time"
)

type Test struct {
	AccountID         string    `db:"account_id" selector:"insert"`
	AppID             int64     `db:"appid" selector:"insert"`
	GameName          string    `db:"game_name" selector:"unnecessary,insert"`
	AccountName       string    `db:"account_name" selector:"unnecessary,insert"`
	Balance           int64     `db:"balance"`
	FrozenBalance     int64     `db:"frozen_balance"`
	ReadyDrawBalance  int64     `db:"ready_draw_balance"`
	ReadyTransBalance int64     `db:"ready_trans_balance"`
	Status            int16     `db:"status"`
	Version           int64     `db:"version"`
	CreateTime        time.Time `db:"create_time" selector:"time,unnecessary"`
	UpdateTime        time.Time `db:"update_time" selector:"time,unnecessary"`
}

var testImpl, _ = New(&Test{})

func TestNew(t *testing.T) {
	type args struct {
		sample interface{}
	}

	tests := []struct {
		name string
		args args
		want *instance
	}{
		{
			name: "general test",
			args: args{
				sample: &Test{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := New(tt.args.sample); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_GetColsWithOmit(t *testing.T) {
	type args struct {
		selectors []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "omit noting",
			args: args{
				selectors: nil,
			},
			want: []string{
				"account_id",
				"appid",
				"game_name",
				"account_name",
				"balance",
				"frozen_balance",
				"ready_draw_balance",
				"ready_trans_balance",
				"status",
				"version",
				"create_time",
				"update_time",
			},
		},

		{
			name: "omit time selector",
			args: args{
				selectors: []string{
					"time",
				},
			},
			want: []string{
				"account_id",
				"appid",
				"game_name",
				"account_name",
				"balance",
				"frozen_balance",
				"ready_draw_balance",
				"ready_trans_balance",
				"status",
				"version",
			},
		},
		{
			name: "omit unnecessary selector",
			args: args{
				selectors: []string{
					"unnecessary",
				},
			},
			want: []string{
				"account_id",
				"appid",
				"balance",
				"frozen_balance",
				"ready_draw_balance",
				"ready_trans_balance",
				"status",
				"version",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testImpl.GetColsWithOmit(tt.args.selectors...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("instance.GetColsWithOmit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_GetColsWithSelect(t *testing.T) {
	type args struct {
		selectors []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "select noting",
			args: args{
				selectors: nil,
			},
			want: nil,
		},
		{
			name: "select insert cols",
			args: args{
				selectors: []string{
					"insert",
				},
			},
			want: []string{
				"account_id",
				"appid",
				"game_name",
				"account_name",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testImpl.GetColsWithSelect(tt.args.selectors...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("instance.GetColsWithSelect() = %v, want %v", got, tt.want)
			}
		})
	}
}
