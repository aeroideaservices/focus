package access_control

import (
	"testing"
)

func TestPrivileges_accessByAction(t *testing.T) {
	type args struct {
		action string
	}
	tests := []struct {
		name string
		p    Privileges
		args args
		want Access
	}{
		{
			p: Privileges{
				"part1": {
					access: Accessed,
					privileges: &Privileges{
						"part2": {
							access: Nil,
							privileges: &Privileges{
								"part3": {
									access:     Denied,
									privileges: nil,
								},
							},
						},
					},
				},
			},
			args: args{
				action: "part1.part2.part3.post",
			},
			want: Denied,
		},
		{
			p: Privileges{
				"part1": {
					access: Accessed,
					privileges: &Privileges{
						"part2": {
							access: Nil,
							privileges: &Privileges{
								"part3": {
									access:     Denied,
									privileges: nil,
								},
							},
						},
					},
				},
			},
			args: args{
				action: "part1.part2.part3",
			},
			want: Denied,
		},
		{
			p: Privileges{
				"part1": {
					access: Accessed,
					privileges: &Privileges{
						"part2": {
							access: Nil,
							privileges: &Privileges{
								"part3": {
									access:     Denied,
									privileges: nil,
								},
							},
						},
						"*": {
							access: Denied,
						},
					},
				},
			},
			args: args{
				action: "part1.part2",
			},
			want: Denied,
		},
		{
			p: Privileges{
				"part1": {
					access: Accessed,
					privileges: &Privileges{
						"part2": {
							access: Nil,
							privileges: &Privileges{
								"part3": {
									access:     Denied,
									privileges: nil,
								},
							},
						},
						"*": {
							access: Denied,
						},
					},
				},
			},
			args: args{
				action: "part1",
			},
			want: Accessed,
		},
		{
			p: Privileges{
				"part1": {
					access: Accessed,
					privileges: &Privileges{
						"part2": {
							access: Nil,
							privileges: &Privileges{
								"part3": {
									access:     Nil,
									privileges: nil,
								},
							},
						},
						"*": {
							access: Denied,
							privileges: &Privileges{
								"part3": {
									access:     Accessed,
									privileges: nil,
								},
							},
						},
					},
				},
			},
			args: args{
				action: "part1.part2.part3",
			},
			want: Accessed,
		},
		{
			p: Privileges{
				"part1": {
					access: Accessed,
					privileges: &Privileges{
						"part2": {
							access: Nil,
							privileges: &Privileges{
								"part3": {
									access:     Nil,
									privileges: nil,
								},
							},
						},
						"*": {
							access: Denied,
							privileges: &Privileges{
								"part3": {
									access:     Nil,
									privileges: nil,
								},
							},
						},
					},
				},
			},
			args: args{
				action: "part1.part2.part3",
			},
			want: Denied,
		},
		{
			p: Privileges{
				"part1": {
					access: Accessed,
					privileges: &Privileges{
						"part2": {
							access: Nil,
						},
						"*": {
							access: Denied,
							privileges: &Privileges{
								"part3": {
									access:     Nil,
									privileges: nil,
								},
							},
						},
					},
				},
			},
			args: args{
				action: "part1.part2.part3",
			},
			want: Denied,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.accessByAction(tt.args.action); got != tt.want {
				t.Errorf("accessByAction() = %v, want %v", got, tt.want)
			}
		})
	}
}
