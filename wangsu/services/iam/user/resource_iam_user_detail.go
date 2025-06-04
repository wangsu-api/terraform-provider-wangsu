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

func ResourceIamUserDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceIamUserDetailRead,
		Schema: map[string]*schema.Schema{
			"login_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy name",
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Response data.",
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
							Description: "The user create time.",
						},
					},
				},
			},
		},
	}
}

func resourceIamUserDetailRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_iam_user_detail.read")
	var diags diag.Diagnostics
	request := &usermanage.QuerySubAccountInfoRequest{}

	var response *usermanage.QuerySubAccountInfoResponse
	if loginName, ok := data.Get("login_name").(string); ok && loginName != "" {
		path := &usermanage.QuerySubAccountInfoPaths{
			LoginName: &loginName,
		}
		var err error
		var requestId string
		err = resource.RetryContext(ctx, time.Duration(2)*time.Minute, func() *resource.RetryError {
			requestId, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseUserManageClient().QueryUser(request, path)
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
		var resultList []interface{}
		var users = map[string]interface{}{
			"login_name":   response.Data.LoginName,
			"display_name": response.Data.DisplayName,
			"status":       response.Data.Status,
			"email":        response.Data.Email,
			"mobile":       response.Data.Mobile,
			"create_time":  response.Data.CreateTime,
		}
		resultList = append(resultList, users)
		err = data.Set("data", resultList)
		if err != nil {
			return nil
		}
		data.SetId(*response.Data.LoginName)
		log.Printf("resource.wangsu_iam_user_detail.read success, requestId: %s", requestId)
		return nil
	} else {
		log.Printf("login_name is required")
		return diags
	}
}
