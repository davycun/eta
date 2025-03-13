package zzd

import (
	"github.com/davycun/eta/pkg/common/logger"
	"net/http"
	"strings"
)

type Emp struct {
	Zzd
}

type ListEmployeeAccountIdsResp struct {
	Success bool `json:"success"`
	Content struct {
		Data []struct {
			AccountId        int    `json:"accountId"`
			AccountCode      string `json:"accountCode"`
			AccountNamespace string `json:"accountNamespace"`
			EmployeeCode     string `json:"employeeCode"`
		} `json:"data"`
		Success         bool   `json:"success"`
		ResponseMessage string `json:"responseMessage"`
		ResponseCode    string `json:"responseCode"`
		BizErrorCode    string `json:"bizErrorCode"`
	} `json:"content"`
	BizErrorCode string `json:"bizErrorCode"`
}
type GetByMobilesResp struct {
	Success bool `json:"success"`
	Content struct {
		Data []struct {
			AccountId    int    `json:"accountId"`
			Mobile       string `json:"mobile"`
			EmployeeCode string `json:"employeeCode"`
			Status       int    `json:"status"`
		} `json:"data"`
		Success         bool   `json:"success"`
		RequestId       string `json:"requestId"`
		ResponseMessage string `json:"responseMessage"`
		ResponseCode    string `json:"responseCode"`
		BizErrorCode    string `json:"bizErrorCode"`
	} `json:"content"`
	BizErrorCode string `json:"bizErrorCode"`
}
type PageOrganizationEmployeePositionsResp struct {
	Success bool `json:"success"`
	Content struct {
		TotalSize int `json:"totalSize"`
		Data      []struct {
			EmployeeName         string `json:"employeeName"`
			GmtCreate            string `json:"gmtCreate"`
			EmpGender            string `json:"empGender"`
			EmployeeCode         string `json:"employeeCode"`
			GovEmpAvatar         string `json:"govEmpAvatar,omitempty"`
			GovEmployeePositions []struct {
				VisibilityIndicatorCode    string `json:"visibilityIndicatorCode"`
				MainJob                    bool   `json:"mainJob"`
				EmpPosInnerInstitutionCode string `json:"empPosInnerInstitutionCode"`
				EmpPosEmployeeRoleCode     string `json:"empPosEmployeeRoleCode"`
				EmployeeCode               string `json:"employeeCode"`
				OrderInOrganization        int    `json:"orderInOrganization"`
				EmpPosUnitCode             string `json:"empPosUnitCode"`
				GmtCreate                  string `json:"gmtCreate"`
				JobAttributesCode          string `json:"jobAttributesCode"`
				OrganizationCode           string `json:"organizationCode"`
				Status                     string `json:"status"`
				GovEmpPosPhoneNo           string `json:"govEmpPosPhoneNo,omitempty"`
				GovEmpPosJob               string `json:"govEmpPosJob,omitempty"`
				PosJobRankCode             string `json:"posJobRankCode,omitempty"`
			} `json:"govEmployeePositions"`
			EmpJobLevelCode        string `json:"empJobLevelCode"`
			EmpBudgetedPostCode    string `json:"empBudgetedPostCode"`
			Status                 string `json:"status"`
			EmpPoliticalStatusCode string `json:"empPoliticalStatusCode,omitempty"`
		} `json:"data"`
		Success         bool   `json:"success"`
		PageSize        int    `json:"pageSize"`
		ResponseMessage string `json:"responseMessage"`
		CurrentPage     int    `json:"currentPage"`
		ResponseCode    string `json:"responseCode"`
		BizErrorCode    string `json:"bizErrorCode"`
	} `json:"content"`
	BizErrorCode string `json:"bizErrorCode"`
}
type PageSearchEmployeeResp struct {
	Success bool `json:"success"`
	Content struct {
		TotalSize int `json:"totalSize"`
		Data      []struct {
			EmployeeName string `json:"employeeName"`
			AccountId    int    `json:"accountId"`
			GovEmpAvatar string `json:"govEmpAvatar"`
			Account      string `json:"account"`
			EmployeeCode string `json:"employeeCode"`
			Status       string `json:"status"`
		} `json:"data"`
		Success         bool   `json:"success"`
		ResponseMessage string `json:"responseMessage"`
		ResponseCode    string `json:"responseCode"`
		BizErrorCode    string `json:"bizErrorCode"`
	} `json:"content"`
	BizErrorCode string `json:"bizErrorCode"`
}

/*
ListEmployeeAccountIds 批量根据员工Code获取员⼯账号ID

https://openplatform-portal.dg-work.cn/portal/#/helpdoc?apiType=serverapi&docKey=2674854
参数:
  - tenantId 租户id
  - employeeCodes 员工code列表（list最大值100）
*/
func (o *Emp) ListEmployeeAccountIds(tenantId string, employeeCodes []string) *ListEmployeeAccountIdsResp {
	res := &ListEmployeeAccountIdsResp{}
	if o.Err != nil {
		return res
	}
	path := "/mozi/employee/listEmployeeAccountIds"
	params := buildParam(map[string]interface{}{
		"tenantId":      tenantId,
		"employeeCodes": employeeCodes,
	})
	header, query := o.signature(http.MethodPost, path, params)
	resp, err := o.client.R().
		SetHeaders(header).
		SetFormDataFromValues(query).
		SetError(&ListEmployeeAccountIdsResp{}).
		SetResult(&ListEmployeeAccountIdsResp{}).
		Post(path)

	if err != nil {
		o.Err = err
		return res
	}
	logger.Debugf("Zzd Emp.ListEmployeeAccountIds resp: %s", resp)
	// {"success":true,"content":{"data":[{"accountId":78698966,"accountCode":"scqfzgghjjxxhj-wzq","accountNamespace":"local","employeeCode":"GE_f843a6af6c41472a8c154d7496550a91"}],"success":true,"responseMessage":"OK","responseCode":"0","bizErrorCode":"0"},"bizErrorCode":"0"}
	if resp.IsError() {
		return resp.Error().(*ListEmployeeAccountIdsResp)
	}
	return resp.Result().(*ListEmployeeAccountIdsResp)
}

/*
GetByMobiles 批量根据手机号获取人员编码

https://openplatform-portal.dg-work.cn/portal/#/helpdoc?apiType=serverapi&docKey=2675012
参数:
  - tenantId 租户ID
  - mobiles 手机号码列表，逗号分隔，最多50个
  - areaCode 手机区号(没有特别说明，固定填写86)
  - namespace 账号类型（没有特别说明,固定填写local）
*/
func (o *Emp) GetByMobiles(tenantId string, mobiles []string, areaCode, namespace string) *GetByMobilesResp {
	res := &GetByMobilesResp{}
	if o.Err != nil {
		return res
	}
	path := "/mozi/employee/get_by_mobiles"
	params := buildParam(map[string]interface{}{
		"tenantId":  tenantId,
		"mobiles":   strings.Join(mobiles, ","),
		"areaCode":  areaCode,
		"namespace": namespace,
	})
	header, query := o.signature(http.MethodPost, path, params)
	resp, err := o.client.R().
		SetHeaders(header).
		SetFormDataFromValues(query).
		SetError(&GetByMobilesResp{}).
		SetResult(&GetByMobilesResp{}).
		Post(path)

	if err != nil {
		o.Err = err
		return res
	}
	logger.Debugf("Zzd Emp.GetByMobiles resp: %s", resp)
	// {"success":true,"content":{"data":[{"accountId":79276621,"mobile":"15279780452","employeeCode":"GE_4b6eb2f4632346e8a386d3206a67d536","status":0},{"accountId":78603252,"mobile":"18818025812","employeeCode":"GE_04815acb3cde41f4a8113b91d0353470","status":0}],"success":true,"requestId":"0aa04cf817037603697985433d0011","responseMessage":"OK","responseCode":"0","bizErrorCode":"0"},"bizErrorCode":"0"}
	if resp.IsError() {
		return resp.Error().(*GetByMobilesResp)
	}
	return resp.Result().(*GetByMobilesResp)
}

/*
PageOrganizationEmployeePositions 查询组织下人员详情

https://openplatform-portal.dg-work.cn/portal/#/helpdoc?apiType=serverapi&docKey=2674970

参数:
  - tenantId 租户id
  - organizationCode 组织code
  - employeeStatus 员工状态，A为有效，F为无效，默认是所有
  - returnTotalSize 是否请求总数，默认是false
  - pageSize 分页大小，默认是20，范围0-100
  - pageNo 请求起始页，默认是1
*/
func (o *Emp) PageOrganizationEmployeePositions(tenantId, organizationCode, employeeStatus string,
	returnTotalSize bool, pageSize, pageNo int) *PageOrganizationEmployeePositionsResp {
	res := &PageOrganizationEmployeePositionsResp{}
	if o.Err != nil {
		return res
	}
	path := "/mozi/organization/pageOrganizationEmployeePositions"
	params := buildParam(map[string]interface{}{
		"tenantId":         tenantId,
		"organizationCode": organizationCode,
		"employeeStatus":   employeeStatus,
		"returnTotalSize":  returnTotalSize,
		"pageSize":         pageSize,
		"pageNo":           pageNo,
	})
	header, query := o.signature(http.MethodPost, path, params)
	resp, err := o.client.R().
		SetHeaders(header).
		SetFormDataFromValues(query).
		SetError(&PageOrganizationEmployeePositionsResp{}).
		SetResult(&PageOrganizationEmployeePositionsResp{}).
		Post(path)

	if err != nil {
		o.Err = err
		return res
	}
	logger.Debugf("Zzd Emp.PageOrganizationEmployeePositions resp: %s", resp)
	// {"success":true,"content":{"totalSize":18,"data":[{"employeeName":"张燕文","gmtCreate":"2023-10-12 11:10:11","empGender":"2","employeeCode":"GE_05922f07dbba4b54860d1c234b941b1a","govEmpAvatar":"$hQHNKaUCo2pwZwMGDAEN2gBAY1BTZEpoL2Q1eU9XTmFwc3QxVnVhcExzaWx1T0lTVENYVmRKclIzQkpici9WTEFOVFY1Q0sxeCtXTWdEQnhjTQ","govEmployeePositions":[{"visibilityIndicatorCode":"XIAN_JI_YI_BAN_GONG_ZUO_REN_YUAN","mainJob":true,"empPosInnerInstitutionCode":"GO_48885c32539d40c29191f0ead71ae5ef","empPosEmployeeRoleCode":"REN_YUAN_JU_SE_QI_TA","employeeCode":"GE_05922f07dbba4b54860d1c234b941b1a","orderInOrganization":14582454,"empPosUnitCode":"GO_3c72b5a3a955461499acfd983ecaeb5d","gmtCreate":"2023-10-12 11:10:11","jobAttributesCode":"1","organizationCode":"GO_48885c32539d40c29191f0ead71ae5ef","status":"A"}],"empJobLevelCode":"ZHI_JI_QI_TA","empBudgetedPostCode":"BIAN_ZHI_QI_TA","status":"A"},{"employeeName":"吕鑫","gmtCreate":"2022-12-08 10:40:19","empGender":"2","employeeCode":"GE_701bc48f835b4ade88bdaec43221d84b","govEmpAvatar":"$hQHNKaUCo2pwZwMGDAEN2gBAUHJPU3ZDSUdieGY0VWJsOGwxUFhSQ3pJMWo5YTZhNzhVaHdESWdEVjdQa2M5UU0xeExsUVFkalBQNFVJakZ0aA","govEmployeePositions":[{"govEmpPosPhoneNo":"0571-89501494","visibilityIndicatorCode":"XIAN_JI_YI_BAN_GONG_ZUO_REN_YUAN","mainJob":true,"empPosInnerInstitutionCode":"GO_48885c32539d40c29191f0ead71ae5ef","empPosEmployeeRoleCode":"REN_YUAN_JU_SE_PU_TONG_YONG_HU","employeeCode":"GE_701bc48f835b4ade88bdaec43221d84b","govEmpPosJob":"科普工作人员","orderInOrganization":3,"empPosUnitCode":"GO_3c72b5a3a955461499acfd983ecaeb5d","gmtCreate":"2022-12-08 10:40:19","jobAttributesCode":"1","organizationCode":"GO_48885c32539d40c29191f0ead71ae5ef","status":"F"}],"empJobLevelCode":"ZHI_JI_QI_TA","empBudgetedPostCode":"BIAN_ZHI_QI_TA","status":"F"},{"employeeName":"丁石峰","gmtCreate":"2023-07-25 16:13:46","empGender":"1","employeeCode":"GE_cb1623ee636e473ab870a5de530b0d20","govEmpAvatar":"$hQHNKaUCo2pwZwMGDAEN2gBASUdKWkg2ZDFMYkVIcHljOStNQUc4RHNHeUxKY0UzMGtxMGg2UnpUNmpaNGM5UU0xeExsUVFkalBQNFVJakZ0aA","govEmployeePositions":[{"visibilityIndicatorCode":"XIAN_JI_YI_BAN_GONG_ZUO_REN_YUAN","mainJob":true,"empPosInnerInstitutionCode":"GO_48885c32539d40c29191f0ead71ae5ef","empPosEmployeeRoleCode":"REN_YUAN_JU_SE_PU_TONG_YONG_HU","employeeCode":"GE_cb1623ee636e473ab870a5de530b0d20","orderInOrganization":14582453,"empPosUnitCode":"GO_3c72b5a3a955461499acfd983ecaeb5d","gmtCreate":"2023-07-25 16:13:46","jobAttributesCode":"1","organizationCode":"GO_48885c32539d40c29191f0ead71ae5ef","status":"A"}],"empJobLevelCode":"ZHI_JI_QI_TA","empBudgetedPostCode":"BIAN_ZHI_QI_TA","status":"A"},{"employeeName":"陈颖力","gmtCreate":"2023-05-19 09:58:27","empGender":"1","employeeCode":"GE_adf70c9394e4463186f4f08cd8a054d5","govEmpAvatar":"$hQHNKaUCo2pwZwMGDAEN2gBATmFEZ0s3TWczSWE2USs3WElKbENGSFl3d01iRUw5dE1ZRnRoRkpLRCtDSC9WTEFOVFY1Q0sxeCtXTWdEQnhjTQ","govEmployeePositions":[{"visibilityIndicatorCode":"OTHER","mainJob":true,"empPosInnerInstitutionCode":"GO_ee2399d5b3de45b89526af4ec60b3e2a","empPosEmployeeRoleCode":"REN_YUAN_JU_SE_QI_TA","employeeCode":"GE_adf70c9394e4463186f4f08cd8a054d5","orderInOrganization":8,"empPosUnitCode":"GO_3c72b5a3a955461499acfd983ecaeb5d","gmtCreate":"2023-05-19 09:58:27","jobAttributesCode":"1","organizationCode":"GO_ee2399d5b3de45b89526af4ec60b3e2a","status":"A"}],"empJobLevelCode":"ZHI_JI_QI_TA","empBudgetedPostCode":"BIAN_ZHI_QI_TA","status":"A"},{"employeeName":"郑瑶","gmtCreate":"2021-03-09 11:24:15","empGender":"2","employeeCode":"GE_9dfaf309d7e44d2cab3aa22665a955dd","govEmpAvatar":"$hQHNKaUCo2pwZwMGDAEN2gBANmErUVlTUG1vdDR0L2lHUnN6ZFpsZmtBREJQVWh2YTFkUmlzd2JodWsxMGM5UU0xeExsUVFkalBQNFVJakZ0aA","govEmployeePositions":[{"govEmpPosPhoneNo":"0571-89501500","visibilityIndicatorCode":"XIAN_JI_YI_BAN_GONG_ZUO_REN_YUAN","mainJob":true,"empPosInnerInstitutionCode":"GO_48885c32539d40c29191f0ead71ae5ef","empPosEmployeeRoleCode":"REN_YUAN_JU_SE_PU_TONG_YONG_HU","employeeCode":"GE_9dfaf309d7e44d2cab3aa22665a955dd","orderInOrganization":2,"empPosUnitCode":"GO_3c72b5a3a955461499acfd983ecaeb5d","gmtCreate":"2021-03-09 11:24:16","jobAttributesCode":"1","organizationCode":"GO_48885c32539d40c29191f0ead71ae5ef","status":"A"}],"empJobLevelCode":"ZHI_JI_QI_TA","empBudgetedPostCode":"BIAN_ZHI_QI_TA","status":"A"},{"employeeName":"施建英","gmtCreate":"2022-04-07 14:55:54","empGender":"9","employeeCode":"GE_443e3ef1a70e4f3cbf2e38308e4ecdbc","govEmpAvatar":"$hQHNKaUCo2pwZwMGDAEN2gBAN2MxRGNaRG5KeUl3QzA5WWdnV3VEWXIyTjFGZ0F5SjB1ZW1yYTFDdXNFMy9WTEFOVFY1Q0sxeCtXTWdEQnhjTQ","govEmployeePositions":[{"govEmpPosPhoneNo":"0571-88258143","visibilityIndicatorCode":"XIAN_JI_YI_BAN_GONG_ZUO_REN_YUAN","mainJob":true,"empPosInnerInstitutionCode":"GO_ee2399d5b3de45b89526af4ec60b3e2a","empPosEmployeeRoleCode":"REN_YUAN_JU_SE_PU_TONG_YONG_HU","employeeCode":"GE_443e3ef1a70e4f3cbf2e38308e4ecdbc","govEmpPosJob":"工作人员","orderInOrganization":5,"empPosUnitCode":"GO_3c72b5a3a955461499acfd983ecaeb5d","gmtCreate":"2022-04-07 14:55:54","jobAttributesCode":"1","organizationCode":"GO_ee2399d5b3de45b89526af4ec60b3e2a","status":"F"}],"empJobLevelCode":"ZHI_JI_QI_TA","empBudgetedPostCode":"BIAN_ZHI_QI_TA","status":"F"},{"employeeName":"田恬","gmtCreate":"2020-04-05 11:44:08","empGender":"2","employeeCode":"GE_db9e813caf65460a9a7341be69745583","govEmpAvatar":"$hQHNKaUCo2pwZwMGDAEN2gAsVDd6cStpdm5DQ2xDdldhd0paTHNCbVEyemtmMFppVDB2eEVseGkwSlJIbz0","govEmployeePositions":[{"posJobRankCode":"CENG_CI_XIANG_KE_JI_ZHENG_ZHI","govEmpPosPhoneNo":"0571-89501498","visibilityIndicatorCode":"XIAN_JI_YI_BAN_GONG_ZUO_REN_YUAN","mainJob":true,"empPosInnerInstitutionCode":"GO_ee2399d5b3de45b89526af4ec60b3e2a","empPosEmployeeRoleCode":"REN_YUAN_JU_SE_PU_TONG_YONG_HU","employeeCode":"GE_db9e813caf65460a9a7341be69745583","govEmpPosJob":"学会部部长","orderInOrganization":4,"empPosUnitCode":"GO_3c72b5a3a955461499acfd983ecaeb5d","gmtCreate":"2021-04-16 11:13:39","jobAttributesCode":"1","organizationCode":"GO_ee2399d5b3de45b89526af4ec60b3e2a","status":"A"}],"empJobLevelCode":"ZHI_JI_ZHENG_XIANG_KE_ZHEN_JI","empBudgetedPostCode":"BIAN_ZHI_CAN_ZHAO_GONG_WU_YUAN_BIAN_ZHI","status":"A"},{"employeeName":"李晖","gmtCreate":"2020-04-05 11:42:45","empGender":"9","employeeCode":"GE_45f7ec74a8ae4f1f8e17284e60d94c8a","govEmpAvatar":"$hQHNKaUCo2pwZwMGDAEN2gBAbUVIREpTRGFmN0ZQbVZ3WUoySk9UV1dvSGxtdzBOcTkrOUlJUmxMRkprWWM5UU0xeExsUVFkalBQNFVJakZ0aA","govEmployeePositions":[{"govEmpPosPhoneNo":"0571-89501496","visibilityIndicatorCode":"XIAN_JI_YI_BAN_GONG_ZUO_REN_YUAN","mainJob":true,"empPosInnerInstitutionCode":"GO_3457f9cbc2f8401e8dd4e0108597333d","empPosEmployeeRoleCode":"REN_YUAN_JU_SE_PU_TONG_YONG_HU","employeeCode":"GE_45f7ec74a8ae4f1f8e17284e60d94c8a","govEmpPosJob":"三调","orderInOrganization":5,"empPosUnitCode":"GO_3c72b5a3a955461499acfd983ecaeb5d","gmtCreate":"2020-04-05 11:42:45","jobAttributesCode":"1","organizationCode":"GO_3457f9cbc2f8401e8dd4e0108597333d","status":"A"}],"empJobLevelCode":"ZHI_JI_QI_TA","empBudgetedPostCode":"BIAN_ZHI_QI_TA","status":"A"},{"employeeName":"陈学军","gmtCreate":"2020-04-05 11:42:44","empGender":"9","employeeCode":"GE_c35f344793c242c88d616f11b0a0aafd","govEmpAvatar":"$hQHNKaUCo2pwZwMGDAEN2gBAbUVIREpTRGFmN0ZQbVZ3WUoySk9UV2dqSHIwZkFWQ2Z1dGVreXVpTmR4WC9WTEFOVFY1Q0sxeCtXTWdEQnhjTQ","govEmployeePositions":[{"govEmpPosPhoneNo":"0571-89501497","visibilityIndicatorCode":"XIAN_JI_YI_BAN_GONG_ZUO_REN_YUAN","mainJob":true,"empPosInnerInstitutionCode":"GO_3457f9cbc2f8401e8dd4e0108597333d","empPosEmployeeRoleCode":"REN_YUAN_JU_SE_PU_TONG_YONG_HU","employeeCode":"GE_c35f344793c242c88d616f11b0a0aafd","govEmpPosJob":"党组成员 四级调研员","orderInOrganization":4,"empPosUnitCode":"GO_3c72b5a3a955461499acfd983ecaeb5d","gmtCreate":"2020-04-05 11:42:45","jobAttributesCode":"1","organizationCode":"GO_3457f9cbc2f8401e8dd4e0108597333d","status":"A"}],"empJobLevelCode":"ZHI_JI_QI_TA","empBudgetedPostCode":"BIAN_ZHI_QI_TA","status":"A"},{"employeeName":"陈伟","gmtCreate":"2020-04-05 11:42:43","empGender":"9","employeeCode":"GE_9b7ca3a5d5454bba90f8499613d75f1b","govEmpAvatar":"$hQHNKaUCo2pwZwMGDAEN2gBAbUVIREpTRGFmN0ZQbVZ3WUoySk9UVHMzWFQrc0ZPQThUVlJ1a3BJMjJTOGM5UU0xeExsUVFkalBQNFVJakZ0aA","govEmployeePositions":[{"govEmpPosPhoneNo":"0571-89501495","visibilityIndicatorCode":"XIAN_JI_YI_BAN_GONG_ZUO_REN_YUAN","mainJob":true,"empPosInnerInstitutionCode":"GO_3457f9cbc2f8401e8dd4e0108597333d","empPosEmployeeRoleCode":"REN_YUAN_JU_SE_PU_TONG_YONG_HU","employeeCode":"GE_9b7ca3a5d5454bba90f8499613d75f1b","govEmpPosJob":"一级调研员","orderInOrganization":3,"empPosUnitCode":"GO_3c72b5a3a955461499acfd983ecaeb5d","gmtCreate":"2020-04-05 11:42:43","jobAttributesCode":"1","organizationCode":"GO_3457f9cbc2f8401e8dd4e0108597333d","status":"A"}],"empJobLevelCode":"ZHI_JI_YI_JI_DIAO_YAN_YUAN","empBudgetedPostCode":"BIAN_ZHI_GONG_WU_YUAN_BIAN_ZHI","status":"A"},{"employeeName":"施潇潇","gmtCreate":"2020-04-05 11:45:06","empGender":"2","employeeCode":"GE_53503ff1b2294aaa8b27da751ebbd4aa","empPoliticalStatusCode":"01","govEmpAvatar":"$hQHNKaUCo2pwZwMGDAEN2gBAM0p1c1VWbkQvWU0yYjRFT0tWQTFndVlnNXlCV0wvUEhEM0dOeTQ0RW9na2M5UU0xeExsUVFkalBQNFVJakZ0aA","govEmployeePositions":[{"govEmpPosPhoneNo":"0571-89501501","visibilityIndicatorCode":"XIAN_JI_YI_BAN_GONG_ZUO_REN_YUAN","mainJob":true,"empPosInnerInstitutionCode":"GO_3457f9cbc2f8401e8dd4e0108597333d","empPosEmployeeRoleCode":"REN_YUAN_JU_SE_PU_TONG_YONG_HU","employeeCode":"GE_53503ff1b2294aaa8b27da751ebbd4aa","govEmpPosJob":"党组成员 副主席","orderInOrganization":2,"empPosUnitCode":"GO_3c72b5a3a955461499acfd983ecaeb5d","gmtCreate":"2022-04-07 16:10:50","jobAttributesCode":"1","organizationCode":"GO_3457f9cbc2f8401e8dd4e0108597333d","status":"A"}],"empJobLevelCode":"ZHI_JI_FU_XIAN_CHU_JI","empBudgetedPostCode":"BIAN_ZHI_GONG_WU_YUAN_BIAN_ZHI","status":"A"},{"employeeName":"张哲持","gmtCreate":"2020-04-05 11:45:54","empGender":"1","employeeCode":"GE_a2e4803010e5462c819d1ab1517272d8","govEmpAvatar":"$hQHNKaUCo2pwZwMGDAEN2gBARS9tU2M4dmxOTEdVb3VpOUI4OWdaN2xPR3NxZXNEVmxTOGJET1lzRzJBbi9WTEFOVFY1Q0sxeCtXTWdEQnhjTQ","govEmployeePositions":[{"posJobRankCode":"CENG_CI_CHU_JI_FU_ZHI","govEmpPosPhoneNo":"0571-89501502","visibilityIndicatorCode":"XIAN_JI_DAN_WEI_FU_ZHI","mainJob":true,"empPosInnerInstitutionCode":"GO_3457f9cbc2f8401e8dd4e0108597333d","empPosEmployeeRoleCode":"REN_YUAN_JU_SE_CHAO_JI_YONG_HU_FU_ZHI","employeeCode":"GE_a2e4803010e5462c819d1ab1517272d8","govEmpPosJob":"党组成员 副主席","orderInOrganization":1,"empPosUnitCode":"GO_3c72b5a3a955461499acfd983ecaeb5d","gmtCreate":"2021-04-20 14:55:33","jobAttributesCode":"1","organizationCode":"GO_3457f9cbc2f8401e8dd4e0108597333d","status":"A"}],"empJobLevelCode":"ZHI_JI_FU_XIAN_CHU_JI","empBudgetedPostCode":"BIAN_ZHI_GONG_WU_YUAN_BIAN_ZHI","status":"A"},{"employeeName":"夏静","gmtCreate":"2020-04-05 11:42:42","empGender":"9","employeeCode":"GE_4e884659152942ea883ee1848e0f94b7","empPoliticalStatusCode":"01","govEmpAvatar":"$hQHNKaUCo2pwZwMGDAEN2gBAbUVIREpTRGFmN0ZQbVZ3WUoySk9UUUM0TUxZZUVnV09icDhYT2d0L3hxVC9WTEFOVFY1Q0sxeCtXTWdEQnhjTQ","govEmployeePositions":[{"posJobRankCode":"CENG_CI_QI_TA","govEmpPosPhoneNo":"0571-89501500","visibilityIndicatorCode":"XIAN_JI_YI_BAN_GONG_ZUO_REN_YUAN","mainJob":true,"empPosInnerInstitutionCode":"GO_48885c32539d40c29191f0ead71ae5ef","empPosEmployeeRoleCode":"REN_YUAN_JU_SE_GAO_JI_YONG_HU","employeeCode":"GE_4e884659152942ea883ee1848e0f94b7","govEmpPosJob":"办公室主任","orderInOrganization":1,"empPosUnitCode":"GO_3c72b5a3a955461499acfd983ecaeb5d","gmtCreate":"2020-04-05 11:42:42","jobAttributesCode":"1","organizationCode":"GO_48885c32539d40c29191f0ead71ae5ef","status":"A"}],"empJobLevelCode":"ZHI_JI_ZHENG_XIANG_KE_ZHEN_JI","empBudgetedPostCode":"BIAN_ZHI_GONG_WU_YUAN_BIAN_ZHI","status":"A"},{"employeeName":"徐天锋","gmtCreate":"2021-05-24 11:17:28","empGender":"2","employeeCode":"GE_40bb2fe7d41748aa8edc7043ad925e1d","govEmpAvatar":"$hQHNKaUCo2pwZwMGDAEN2gBAdFIzOHNvM29Jb3RKNHVIa1FjeVd6b3BkVWxuendBbFhoYVVYQTNWeS9vbi9WTEFOVFY1Q0sxeCtXTWdEQnhjTQ","govEmployeePositions":[{"govEmpPosPhoneNo":"0571-87917109","visibilityIndicatorCode":"XIAN_JI_YI_BAN_GONG_ZUO_REN_YUAN","mainJob":true,"empPosInnerInstitutionCode":"GO_ee2399d5b3de45b89526af4ec60b3e2a","empPosEmployeeRoleCode":"REN_YUAN_JU_SE_QI_TA","employeeCode":"GE_40bb2fe7d41748aa8edc7043ad925e1d","govEmpPosJob":"工作人员","orderInOrganization":1,"empPosUnitCode":"GO_3c72b5a3a955461499acfd983ecaeb5d","gmtCreate":"2021-05-24 11:17:28","jobAttributesCode":"1","organizationCode":"GO_ee2399d5b3de45b89526af4ec60b3e2a","status":"F"}],"empJobLevelCode":"ZHI_JI_QI_TA","empBudgetedPostCode":"BIAN_ZHI_LI_SHI_GONG","status":"F"},{"employeeName":"丁伟韦","gmtCreate":"2021-04-27 12:47:34","empGender":"2","employeeCode":"GE_1830eb0ac28b47e987bb2d9cc79cd65c","govEmpAvatar":"$hQHNKaUCo2pwZwMGDAEN2gBAMkxsLzhmd3JuVW9BWW5VTmZIT1B0K3NXa0x5YUR6akl5UTZZbCtGTDdiNy9WTEFOVFY1Q0sxeCtXTWdEQnhjTQ","govEmployeePositions":[{"govEmpPosPhoneNo":"0571-87910576","visibilityIndicatorCode":"XIAN_JI_YI_BAN_GONG_ZUO_REN_YUAN","mainJob":true,"empPosInnerInstitutionCode":"GO_42c87adcade14f998d6b2d2fbe5c3bf5","empPosEmployeeRoleCode":"REN_YUAN_JU_SE_QI_TA","employeeCode":"GE_1830eb0ac28b47e987bb2d9cc79cd65c","govEmpPosJob":"工作人员","orderInOrganization":1,"empPosUnitCode":"GO_3c72b5a3a955461499acfd983ecaeb5d","gmtCreate":"2021-04-27 12:47:34","jobAttributesCode":"1","organizationCode":"GO_42c87adcade14f998d6b2d2fbe5c3bf5","status":"A"}],"empJobLevelCode":"ZHI_JI_QI_TA","empBudgetedPostCode":"BIAN_ZHI_LI_SHI_GONG","status":"A"},{"employeeName":"陶丽南","gmtCreate":"2020-04-05 11:42:44","empGender":"9","employeeCode":"GE_88fbf7afa5024122b5e82088ecc12986","empPoliticalStatusCode":"","govEmployeePositions":[{"posJobRankCode":"","govEmpPosPhoneNo":"0571-87829060","visibilityIndicatorCode":"XIAN_JI_DAN_WEI_FU_ZHI","mainJob":true,"empPosInnerInstitutionCode":"GO_3457f9cbc2f8401e8dd4e0108597333d","empPosEmployeeRoleCode":"REN_YUAN_JU_SE_PU_TONG_YONG_HU","employeeCode":"GE_88fbf7afa5024122b5e82088ecc12986","govEmpPosJob":"副调研员","orderInOrganization":262166,"empPosUnitCode":"GO_3c72b5a3a955461499acfd983ecaeb5d","gmtCreate":"2020-04-05 11:42:44","jobAttributesCode":"1","organizationCode":"GO_3457f9cbc2f8401e8dd4e0108597333d","status":"F"}],"empJobLevelCode":"ZHI_JI_QI_TA","empBudgetedPostCode":"BIAN_ZHI_QI_TA","status":"F"},{"employeeName":"金燮彪","gmtCreate":"2020-04-05 11:42:44","empGender":"9","employeeCode":"GE_8750346add764fc89712431738110890","govEmployeePositions":[{"posJobRankCode":"","govEmpPosPhoneNo":"0571-87826127","visibilityIndicatorCode":"XIAN_JI_DAN_WEI_FU_ZHI","mainJob":true,"empPosInnerInstitutionCode":"GO_3457f9cbc2f8401e8dd4e0108597333d","empPosEmployeeRoleCode":"REN_YUAN_JU_SE_PU_TONG_YONG_HU","employeeCode":"GE_8750346add764fc89712431738110890","govEmpPosJob":"调研员","orderInOrganization":5,"empPosUnitCode":"GO_3c72b5a3a955461499acfd983ecaeb5d","gmtCreate":"2020-04-05 11:42:44","jobAttributesCode":"1","organizationCode":"GO_3457f9cbc2f8401e8dd4e0108597333d","status":"F"}],"empJobLevelCode":"ZHI_JI_QI_TA","empBudgetedPostCode":"BIAN_ZHI_QI_TA","status":"F"},{"employeeName":"来源","gmtCreate":"2020-04-05 11:42:43","empGender":"9","employeeCode":"GE_6275b11f9023441a9995fdb8c34def4f","empPoliticalStatusCode":"","govEmpAvatar":"$hQHNKaUCo2pwZwMGDAEN2gBAbUVIREpTRGFmN0ZQbVZ3WUoySk9UYnVveTdUdTlSUWp6RHVxZHlQV1BIZi9WTEFOVFY1Q0sxeCtXTWdEQnhjTQ","govEmployeePositions":[{"posJobRankCode":"","govEmpPosPhoneNo":"0571-87816498","visibilityIndicatorCode":"XIAN_JI_YI_BAN_GONG_ZUO_REN_YUAN","mainJob":true,"empPosInnerInstitutionCode":"GO_48885c32539d40c29191f0ead71ae5ef","empPosEmployeeRoleCode":"REN_YUAN_JU_SE_PU_TONG_YONG_HU","employeeCode":"GE_6275b11f9023441a9995fdb8c34def4f","govEmpPosJob":"办公室文秘","orderInOrganization":14582452,"empPosUnitCode":"GO_3c72b5a3a955461499acfd983ecaeb5d","gmtCreate":"2020-04-05 11:42:43","jobAttributesCode":"1","organizationCode":"GO_48885c32539d40c29191f0ead71ae5ef","status":"F"}],"empJobLevelCode":"ZHI_JI_QI_TA","empBudgetedPostCode":"BIAN_ZHI_QI_TA","status":"F"}],"success":true,"pageSize":20,"responseMessage":"OK","currentPage":1,"responseCode":"0","bizErrorCode":"0"},"bizErrorCode":"0"}
	if resp.IsError() {
		return resp.Error().(*PageOrganizationEmployeePositionsResp)
	}
	return resp.Result().(*PageOrganizationEmployeePositionsResp)
}

/*
PageSearchEmployee 人员信息查询

https://openplatform-portal.dg-work.cn/portal/#/helpdoc?apiType=serverapi&docKey=2796890

参数:
  - returnTotalSize 是否返回查询结果总数 默认不需要
  - inOrganizationCode 组织code,在当前组织下查询 优先级高于cascadeOrganizationCode 两个参数至少有一个
  - pageSize 每页条数, 默认20, 最大只能100
  - pageNo 当前页码, 开始页码为1, 小于1认为为1
  - tenantId 租户id
  - cascadeOrganizationCode 组织code,在组织级联下级中查询
  - status A/F （在职/离职）默认返回所有
  - nameKeywords 人员姓名关键字
*/
func (o *Emp) PageSearchEmployee(returnTotalSize bool, inOrganizationCode string, pageSize int, pageNo int,
	tenantId, cascadeOrganizationCode, status, nameKeywords string) *PageSearchEmployeeResp {
	res := &PageSearchEmployeeResp{}
	if o.Err != nil {
		return res
	}
	path := "/mozi/fusion/pageSearchEmployee"
	params := buildParam(map[string]interface{}{
		"returnTotalSize":         returnTotalSize,
		"inOrganizationCode":      inOrganizationCode,
		"pageSize":                pageSize,
		"pageNo":                  pageNo,
		"tenantId":                tenantId,
		"cascadeOrganizationCode": cascadeOrganizationCode,
		"status":                  status,
		"nameKeywords":            nameKeywords,
	})
	header, query := o.signature(http.MethodPost, path, params)
	resp, err := o.client.R().
		SetHeaders(header).
		SetFormDataFromValues(query).
		SetError(&PageSearchEmployeeResp{}).
		SetResult(&PageSearchEmployeeResp{}).
		Post(path)

	if err != nil {
		o.Err = err
		return res
	}
	logger.Debugf("Zzd Emp.PageSearchEmployee resp: %s", resp)

	// {"success":true,"content":{"totalSize":1,"data":[{"employeeName":"王元屏","accountId":32281684,"govEmpAvatar":"$hQHNKaUCo2pwZwMGDAEN2gBAYVNMeVN4R0p0enF6UkJYclVsbWhOYXA2V01MZjUreGVpNHZZMzAvbUtSai9WTEFOVFY1Q0sxeCtXTWdEQnhjTQ","account":"df11939910933","employeeCode":"GE_f2cc16157b9c42e8b617d2d4600747ff","status":"A"}],"success":true,"responseMessage":"OK","responseCode":"0","bizErrorCode":"0"},"bizErrorCode":"0"}
	if resp.IsError() {
		return resp.Error().(*PageSearchEmployeeResp)
	}
	return resp.Result().(*PageSearchEmployeeResp)
}
