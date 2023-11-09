package access_control

import (
	"reflect"
	"testing"
)

func TestNewRole(t *testing.T) {
	type args struct {
		name              string
		privilegesStrings []string
	}
	tests := []struct {
		name string
		args args
		want *Role
	}{
		{
			name: "void role",
			args: args{
				name:              "test-role-1",
				privilegesStrings: []string{},
			},
			want: &Role{
				name:       "test-role-1",
				privileges: &Privileges{},
			},
		},
		{
			name: "test role 2",
			args: args{
				name:              "test-role-2",
				privilegesStrings: []string{"content"},
			},
			want: &Role{
				name: "test-role-2",
				privileges: &Privileges{
					"content": &Privilege{
						access: Accessed,
					},
				},
			},
		},
		{
			name: "test role 3",
			args: args{
				name:              "test-role-1",
				privilegesStrings: []string{"content", "content.models", "!content.models.banners", "content.media"},
			},
			want: &Role{
				name: "test-role-1",
				privileges: &Privileges{
					"content": &Privilege{
						access: Accessed,
						privileges: &Privileges{
							"models": &Privilege{
								access: Accessed,
								privileges: &Privileges{
									"banners": &Privilege{
										access: Denied,
									},
								},
							},
							"media": &Privilege{
								access: Accessed,
							},
						},
					},
				},
			},
		},
		{
			name: "",
			args: args{
				name: "",
				privilegesStrings: []string{
					"*",
					"!catalog.models.product.elements.post",
					"!catalog.models.product.elements.*.delete",
					"!catalog.models.client-category.elements.post",
					"!catalog.models.client-category.elements.*.delete",
					"!catalog.models.actarea.elements.post",
					"!catalog.models.actarea.elements.*.delete",
					"!catalog.models.store.elements.post",
					"!catalog.models.store.elements.*.delete",
				},
			},
			want: &Role{
				name: "",
				privileges: &Privileges{
					"*": {access: Accessed},
					"catalog": {privileges: &Privileges{
						"models": {privileges: &Privileges{
							"product": {privileges: &Privileges{
								"elements": {privileges: &Privileges{
									"post": {access: Denied},
									"*": {privileges: &Privileges{
										"delete": {access: Denied},
									}},
								}},
							}},
							"client-category": {privileges: &Privileges{
								"elements": {privileges: &Privileges{
									"post": {access: Denied},
									"*": {privileges: &Privileges{
										"delete": {access: Denied},
									}},
								}},
							}},
							"actarea": {privileges: &Privileges{
								"elements": {privileges: &Privileges{
									"post": {access: Denied},
									"*": {privileges: &Privileges{
										"delete": {access: Denied},
									}},
								}},
							}},
							"store": {privileges: &Privileges{
								"elements": {privileges: &Privileges{
									"post": {access: Denied},
									"*": {privileges: &Privileges{
										"delete": {access: Denied},
									}},
								}},
							}},
						}},
					}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRole(tt.args.name, tt.args.privilegesStrings); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
