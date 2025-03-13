package zzd

import (
	"github.com/davycun/eta/pkg/common/logger"
	"net/http"
)

type Org struct {
	Zzd
}

type GetPropertiesByOrgTypeCodeResp struct {
	Success bool `json:"success"`
	Content struct {
		Data []struct {
			Type           string `json:"type"`
			Code           string `json:"code"`
			CreateRequired string `json:"createRequired"`
			Name           string `json:"name"`
			GroupType      string `json:"groupType"`
		} `json:"data"`
		Success         bool   `json:"success"`
		ResponseMessage string `json:"responseMessage"`
		ResponseCode    string `json:"responseCode"`
		BizErrorCode    string `json:"bizErrorCode"`
	} `json:"content"`
	BizErrorCode string `json:"bizErrorCode"`
}
type PageSearchOrganizationResp struct {
	Success bool `json:"success"`
	Content struct {
		TotalSize int `json:"totalSize"`
		Data      []struct {
			OrganizationName string `json:"organizationName"`
			ParentCode       string `json:"parentCode"`
			OrganizationCode string `json:"organizationCode"`
			Status           string `json:"status"`
		} `json:"data"`
		Success         bool   `json:"success"`
		ResponseMessage string `json:"responseMessage"`
		ResponseCode    string `json:"responseCode"`
		BizErrorCode    string `json:"bizErrorCode"`
	} `json:"content"`
	BizErrorCode string `json:"bizErrorCode"`
}
type GetOrganizationByCodeResp struct {
	Success bool `json:"success"`
	Content struct {
		Data struct {
			DisplayOrder     int    `json:"displayOrder"`
			TypeName         string `json:"typeName"`
			ParentCode       string `json:"parentCode"`
			OrganizationName string `json:"organizationName"`
			Leaf             bool   `json:"leaf"`
			GmtCreate        string `json:"gmtCreate"`
			TypeCode         string `json:"typeCode"`
			ParentName       string `json:"parentName"`
			OrganizationCode string `json:"organizationCode"`
			Status           string `json:"status"`
		} `json:"data"`
		Success         bool   `json:"success"`
		ResponseMessage string `json:"responseMessage"`
		ResponseCode    string `json:"responseCode"`
		BizErrorCode    string `json:"bizErrorCode"`
	} `json:"content"`
	BizErrorCode string `json:"bizErrorCode"`
}
type PageSubOrganizationCodesResp struct {
	Success bool `json:"success"`
	Content struct {
		TotalSize       int      `json:"totalSize"`
		Data            []string `json:"data"`
		Success         bool     `json:"success"`
		RequestId       string   `json:"requestId"`
		PageSize        int      `json:"pageSize"`
		ResponseMessage string   `json:"responseMessage"`
		CurrentPage     int      `json:"currentPage"`
		ResponseCode    string   `json:"responseCode"`
		BizErrorCode    string   `json:"bizErrorCode"`
	} `json:"content"`
}
type ListOrganizationsByCodesResp struct {
	Success bool `json:"success"`
	Content struct {
		Data []struct {
			UnifiedSocialCreditCode  string `json:"unifiedSocialCreditCode"`
			PostalCode               string `json:"postalCode"`
			DisplayOrder             int    `json:"displayOrder"`
			TypeName                 string `json:"typeName"`
			ParentCode               string `json:"parentCode"`
			InstitutionCode          string `json:"institutionCode"`
			ContactNumber            string `json:"contactNumber"`
			BusinessStripCodes       string `json:"businessStripCodes"`
			InstitutionLevelCode     string `json:"institutionLevelCode"`
			Address                  string `json:"address"`
			OrganizationName         string `json:"organizationName"`
			ResponsibleEmployeeCodes string `json:"responsibleEmployeeCodes"`
			Leaf                     bool   `json:"leaf"`
			GmtCreate                string `json:"gmtCreate"`
			TypeCode                 string `json:"typeCode"`
			DivisionCode             string `json:"divisionCode"`
			ParentName               string `json:"parentName"`
			OrganizationCode         string `json:"organizationCode"`
			OtherName                string `json:"otherName"`
			ShortName                string `json:"shortName"`
			Remarks                  string `json:"remarks"`
			Status                   string `json:"status"`
		} `json:"data"`
		Success         bool   `json:"success"`
		ResponseMessage string `json:"responseMessage"`
		ResponseCode    string `json:"responseCode"`
		BizErrorCode    string `json:"bizErrorCode"`
	} `json:"content"`
	BizErrorCode string `json:"bizErrorCode"`
}

/*
GetPropertiesByOrgTypeCode 查询组织类型支持的属性字段

	https://openplatform-portal.dg-work.cn/portal/#/helpdoc?apiType=serverapi&docKey=2967290

参数:
  - tenantId 租户id
  - orgTypeCode 组织类型code
*/
func (o *Org) GetPropertiesByOrgTypeCode(tenantId, orgTypeCode string) *GetPropertiesByOrgTypeCodeResp {
	res := &GetPropertiesByOrgTypeCodeResp{}
	if o.Err != nil {
		return res
	}
	path := "/mozi/fusion/getPropertiesByOrgTypeCode"
	params := buildParam(map[string]interface{}{
		"tenantId":    tenantId,
		"orgTypeCode": orgTypeCode,
	})
	header, query := o.signature(http.MethodPost, path, params)
	resp, err := o.client.R().
		SetHeaders(header).
		SetFormDataFromValues(query).
		SetError(&GetPropertiesByOrgTypeCodeResp{}).
		SetResult(&GetPropertiesByOrgTypeCodeResp{}).
		Post(path)

	if err != nil {
		o.Err = err
		return res
	}
	logger.Debugf("Zzd Org.GetPropertiesByOrgTypeCode resp: %s", resp)
	// {"success":true,"content":{"data":[{"type":"3","code":"govDivision","createRequired":true,"name":"行政区划","groupType":"1"},{"type":"3","code":"govBusinessStrips","createRequired":true,"name":"条线","groupType":"1"},{"type":"0","code":"orgName","createRequired":true,"name":"节点名称","groupType":"1"}],"success":true,"responseMessage":"OK","responseCode":"0","bizErrorCode":"0"},"bizErrorCode":"0"}
	if resp.IsError() {
		return resp.Error().(*GetPropertiesByOrgTypeCodeResp)
	}
	return resp.Result().(*GetPropertiesByOrgTypeCodeResp)
}

/*
PageSearchOrganization 组织信息查询,根据组织名称关键词进行组织信息查询

	https://openplatform-portal.dg-work.cn/portal/#/helpdoc?apiType=serverapi&docKey=2796891

参数:
  - tenant_id: 租户ID
  - nameKeywords: 组织姓名关键字
  - inOrganizationCode: 组织code,在当前组织下查询 优先级高于cascadeOrganizationCode 两个参数至少有一个
  - cascadeOrganizationCode: 组织code,在组织级联下级中查询
  - returnTotalSize: 是否返回查询结果总数 默认不需要
  - pageSize: 每页条数, 默认20, 最大只能100
  - pageNo: 当前页码, 开始页码为1, 小于1认为为1
  - status: A/F（可用/冻结） 默认返回所有
*/
func (o *Org) PageSearchOrganization(tenantId, nameKeywords, inOrganizationCode, cascadeOrganizationCode string,
	returnTotalSize bool, pageSize, pageNo int, status string) *PageSearchOrganizationResp {
	res := &PageSearchOrganizationResp{}
	if o.Err != nil {
		return res
	}
	path := "/mozi/fusion/pageSearchOrganization"
	params := buildParam(map[string]interface{}{
		"tenantId":                tenantId,
		"nameKeywords":            nameKeywords,
		"inOrganizationCode":      inOrganizationCode,
		"cascadeOrganizationCode": cascadeOrganizationCode,
		"returnTotalSize":         returnTotalSize,
		"pageSize":                pageSize,
		"pageNo":                  pageNo,
		"status":                  status,
	})
	header, query := o.signature(http.MethodPost, path, params)
	resp, err := o.client.R().
		SetHeaders(header).
		SetFormDataFromValues(query).
		SetError(&PageSearchOrganizationResp{}).
		SetResult(&PageSearchOrganizationResp{}).
		Post(path)

	if err != nil {
		o.Err = err
		return res
	}
	logger.Debugf("Zzd Org.PageSearchOrganization resp: %s", resp)
	// {"success":true,"content":{"data":[{"organizationName":"信息宣传科","parentCode":"GO_010be2b413884608bec0222264860bc8","organizationCode":"GO_bf5f6add84324d3c8c01c44dd4ec9a94","status":"A"}],"success":true,"responseMessage":"OK","responseCode":"0","bizErrorCode":"0"},"bizErrorCode":"0"}
	if resp.IsError() {
		return resp.Error().(*PageSearchOrganizationResp)
	}
	return resp.Result().(*PageSearchOrganizationResp)
}

/*
GetOrganizationByCode 根据组织 Code 查询详情

	https://openplatform-portal.dg-work.cn/portal/#/helpdoc?apiType=serverapi&docKey=2674856
	https://openplatform-portal.dg-work.cn/portal/#/helpdoc?apiType=serverapi&docKey=2674856

参数:
  - tenantId 租户id
  - organizationCode 组织 code
*/
func (o *Org) GetOrganizationByCode(tenantId, organizationCode string) *GetOrganizationByCodeResp {
	res := &GetOrganizationByCodeResp{}
	if o.Err != nil {
		return res
	}
	path := "/mozi/organization/getOrganizationByCode"
	params := buildParam(map[string]interface{}{
		"tenantId":         tenantId,
		"organizationCode": organizationCode,
	})
	header, query := o.signature(http.MethodPost, path, params)
	resp, err := o.client.R().
		SetHeaders(header).
		SetFormDataFromValues(query).
		SetError(&GetOrganizationByCodeResp{}).
		SetResult(&GetOrganizationByCodeResp{}).
		Post(path)

	if err != nil {
		o.Err = err
		return res
	}
	logger.Debugf("Zzd Org.GetOrganizationByCode resp: %s", resp)
	// {"success":true,"content":{"data":{"displayOrder":3,"typeName":"内设机构","parentCode":"GO_010be2b413884608bec0222264860bc8","organizationName":"信息宣传科","leaf":true,"gmtCreate":"2021-10-25 14:14:40","typeCode":"GOV_INTERNAL_INSTITUTION","parentName":"上城区政协办公室","organizationCode":"GO_bf5f6add84324d3c8c01c44dd4ec9a94","status":"A"},"success":true,"responseMessage":"OK","responseCode":"0","bizErrorCode":"0"},"bizErrorCode":"0"}
	if resp.IsError() {
		return resp.Error().(*GetOrganizationByCodeResp)
	}
	return resp.Result().(*GetOrganizationByCodeResp)
}

/*
PageSubOrganizationCodes 分页获取下⼀级组织 Code 列表

	https://openplatform-portal.dg-work.cn/portal/#/helpdoc?apiType=serverapi&docKey=2674857

参数:
  - tenantId 租户ID
  - organizationCode 组织Code
  - status 查询下一级子组织状态条件
    A - 查询有效的数据
    F - 查询无效的数据
    TOTAL - 查询所有的数据
  - pageSize 每页条数, 默认20, 最大只能100
  - pageNo 当前页码, 开始页码为1, 小于1认为为1
  - returnTotalSize 是否返回查询结果总数默认不需要
*/
func (o *Org) PageSubOrganizationCodes(tenantId, organizationCode, status string,
	pageSize int, pageNo int, returnTotalSize bool) *PageSubOrganizationCodesResp {
	res := &PageSubOrganizationCodesResp{}
	if o.Err != nil {
		return res
	}
	path := "/mozi/organization/pageSubOrganizationCodes"
	params := buildParam(map[string]interface{}{
		"tenantId":         tenantId,
		"organizationCode": organizationCode,
		"status":           status,
		"pageSize":         pageSize,
		"pageNo":           pageNo,
		"returnTotalSize":  returnTotalSize,
	})
	header, query := o.signature(http.MethodPost, path, params)

	resp, err := o.client.R().
		SetHeaders(header).
		SetFormDataFromValues(query).
		SetError(&PageSubOrganizationCodesResp{}).
		SetResult(&PageSubOrganizationCodesResp{}).
		Post(path)

	if err != nil {
		o.Err = err
		return res
	}
	logger.Debugf("Zzd Org.PageSubOrganizationCodes resp: %s", resp)
	// {"success":true,"content":{"success":false,"requestId":"6a5746ec39b64a639f9951fdd7eab243","responseMessage":"scq_ybt没有[GO_3c72b5a3a955461499acfd983ecaeb5f]组织的操作权限","responseCode":"300037","bizErrorCode":"MFH-B003-02-16-0037"},"bizErrorCode":"MFH-B003-02-16-0037"}
	// {"success":true,"content":{"totalSize":4,"data":["GO_3457f9cbc2f8401e8dd4e0108597333d","GO_48885c32539d40c29191f0ead71ae5ef","GO_b40b856faae14995801aefc8d71de500","GO_ee2399d5b3de45b89526af4ec60b3e2a"],"success":true,"requestId":"38fc9a11-60eb-4a94-8eab-aa2c170aad4e","pageSize":100,"responseMessage":"OK","currentPage":1,"responseCode":"0","bizErrorCode":"0"},"bizErrorCode":"0"}
	if resp.IsError() {
		return resp.Error().(*PageSubOrganizationCodesResp)
	}
	return resp.Result().(*PageSubOrganizationCodesResp)
}

/*
ListOrganizationsByCodes 批量根据组织Code查询详情

	https://openplatform-portal.dg-work.cn/portal/#/helpdoc?apiType=serverapi&docKey=2674850

参数:
  - tenantId 租户ID
  - organizationCodes 组织code列表（list最大值100）
*/
func (o *Org) ListOrganizationsByCodes(tenantId string, organizationCodes []string) *ListOrganizationsByCodesResp {
	res := &ListOrganizationsByCodesResp{}
	if o.Err != nil {
		return res
	}
	path := "/mozi/organization/listOrganizationsByCodes"
	params := buildParam(map[string]interface{}{
		"tenantId":          tenantId,
		"organizationCodes": organizationCodes,
	})
	header, query := o.signature(http.MethodPost, path, params)
	resp, err := o.client.R().
		SetHeaders(header).
		SetFormDataFromValues(query).
		SetError(&ListOrganizationsByCodesResp{}).
		SetResult(&ListOrganizationsByCodesResp{}).
		Post(path)

	if err != nil {
		o.Err = err
		return res
	}
	logger.Debugf("Zzd Org.ListOrganizationsByCodes resp: %s", resp)
	// {"success":true,"content":{"data":[{"unifiedSocialCreditCode":"133301020024940758","postalCode":"310000","displayOrder":3,"typeName":"单位","parentCode":"GO_290510a2db9a484d8c4ea3cc42318beb","institutionCode":"133301020024940758","contactNumber":"1900000000","businessStripCodes":"9999","institutionLevelCode":"DAN_WEI_JI_BIE_QUN_ZHONG_ZI_ZHI_TUAN_TI","address":"未知","organizationName":"上城区科学技术协会","responsibleEmployeeCodes":"GE_a2e4803010e5462c819d1ab1517272d8","leaf":false,"gmtCreate":"2020-04-05 07:52:21","typeCode":"GOV_UNIT","divisionCode":"330102000000","parentName":"群众团体","organizationCode":"GO_3c72b5a3a955461499acfd983ecaeb5d","otherName":"区科学技术协会","shortName":"杭州市上城区科学技术协会","remarks":"来自IDaaS推送","status":"A"},{"unifiedSocialCreditCode":"13330102501905350Q","postalCode":"310000","displayOrder":4,"typeName":"单位","parentCode":"GO_290510a2db9a484d8c4ea3cc42318beb","institutionCode":"13330102501905350Q","businessStripCodes":"800","institutionLevelCode":"DAN_WEI_JI_BIE_QUN_ZHONG_ZI_ZHI_TUAN_TI","address":"杭州市凯旋路125号","organizationName":"上城区残疾人联合会","responsibleEmployeeCodes":"GE_e7243381d9d34f00b202e7cbe71a8f50","leaf":false,"gmtCreate":"2020-04-05 07:52:22","typeCode":"GOV_UNIT","divisionCode":"330102000000","parentName":"群众团体","organizationCode":"GO_5ab2c9fd47b04bea8f0c8d06f5eaa4d5","otherName":"上城区残联","shortName":"杭州市上城区残疾人联合会","status":"A"}],"success":true,"responseMessage":"OK","responseCode":"0","bizErrorCode":"0"},"bizErrorCode":"0"}
	if resp.IsError() {
		return resp.Error().(*ListOrganizationsByCodesResp)
	}
	return resp.Result().(*ListOrganizationsByCodesResp)
}
