package user

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	"github.com/wangsu-api/wangsu-sdk-go/wangsu/usermanage"
	"golang.org/x/net/context"
)

func ResourceIamUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceIamUsersRead,
		Schema: map[string]*schema.Schema{
			"page_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Policy name",
			},
			"page_number": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Policy name",
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Detailed data on the results of the request",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"login_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User login name.",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User display name.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "User status. Options: 1-active, 0-inactive.",
						},
						"email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User email address.",
						},
						"mobile": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User mobile number.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether console access is enabled. Options: 1-enabled, 0-disabled.",
						},
					},
				},
			},
		},
	}
}

func resourceIamUsersRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_iam_users.user_list")
	var diags diag.Diagnostics

	request := &usermanage.GetSubAccountListRequest{}
	if pageSize, ok := data.Get("page_size").(int); ok {
		request.PageSize = &pageSize
	}
	if pageIndex, ok := data.Get("page_number").(int); ok {
		request.PageIndex = &pageIndex
	}
	var response *usermanage.GetSubAccountListResponse
	var err error
	var requestId string
	err = resource.RetryContext(ctx, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseUserManageClient().ListUsers(request)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	if response == nil || response.Data == nil {
		data.SetId("")
		return nil
	}

	users := make([]map[string]interface{}, 0)
	for _, user := range response.Data.Rows {
		users = append(users, map[string]interface{}{
			"login_name":   user.LoginName,
			"display_name": user.DisplayName,
			"status":       user.Status,
			"email":        user.Email,
			"mobile":       user.Mobile,
			"create_time":  user.CreateTime,
		})
	}

	data.SetId("user_list")
	err = data.Set("data", users)
	if err != nil {
		return nil
	}
	log.Printf("resource.wangsu_iam_users.user_list success, requestId: %s", requestId)
	return diags
}
