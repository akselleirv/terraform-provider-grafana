package grafana

import (
	"context"

	goapi "github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/client/users"
	"github.com/grafana/grafana-openapi-client-go/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DatasourceUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: readUsers,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Description: `
* [Official documentation](https://grafana.com/docs/grafana/latest/administration/user-management/server-user-management/)
* [HTTP API](https://grafana.com/docs/grafana/latest/developers/http_api/user/)
		
This data source uses Grafana's admin APIs for reading users which
does not currently work with API Tokens. You must use basic auth.
		`,

		Schema: map[string]*schema.Schema{
			"users": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The Grafana instance's users.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The user ID.",
						},
						"login": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user's login.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user's name.",
						},
						"email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user's email.",
						},
						"is_admin": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the user is admin or not.",
						},
					},
				},
			},
		},
	}
}

func readUsers(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := OAPIGlobalClient(meta) // Users are global/org-agnostic
	allUsers, err := getAllUsers(client)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("grafana_users")
	return diag.FromErr(d.Set("users", flattenUsers(allUsers)))
}

func flattenUsers(items []*models.UserSearchHitDTO) []interface{} {
	userItems := make([]interface{}, 0)
	for _, user := range items {
		f := map[string]interface{}{
			"id":       user.ID,
			"login":    user.Login,
			"name":     user.Name,
			"email":    user.Email,
			"is_admin": user.IsAdmin,
		}
		userItems = append(userItems, f)
	}

	return userItems
}

func getAllUsers(client *goapi.GrafanaHTTPAPI) ([]*models.UserSearchHitDTO, error) {
	allUsers := []*models.UserSearchHitDTO{}
	var page int64 = 1
	params := users.NewSearchUsersParams().WithDefaults()
	for {
		resp, err := client.Users.SearchUsers(params.WithPage(&page), nil)
		if err != nil {
			return nil, err
		}

		allUsers = append(allUsers, resp.Payload...)
		if len(resp.Payload) != int(*params.Perpage) {
			break
		}
		page++
	}
	return allUsers, nil
}
