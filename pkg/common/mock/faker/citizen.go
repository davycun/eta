package faker

// CitizenLeaveReason 市民-搬离原因
func CitizenLeaveReason() string {
	return getRandValue([]string{"citizen", "leave_reason"})
}

// CitizenBasic 市民-居住事由
func CitizenBasic() string { return getRandValue([]string{"citizen", "basic"}) }

// CitizenLiveType 市民-居住类型
func CitizenLiveType() string { return getRandValue([]string{"citizen", "live_type"}) }

// CitizenOwnerRelation 市民-与房主关系
func CitizenOwnerRelation() string { return getRandValue([]string{"citizen", "owner_relation"}) }

// CitizenPep2pepType 市民-人与人的关系类型
func CitizenPep2pepType() string { return getRandValue([]string{"citizen", "pep2pep_type"}) }

// CitizenJob 市民-岗位
func CitizenJob() string { return getRandValue([]string{"citizen", "job"}) }

// CitizenJobType 市民-商业职务
func CitizenJobType() string { return getRandValue([]string{"citizen", "job_type"}) }

// CitizenFloorType 市民-楼层属性
func CitizenFloorType() string { return getRandValue([]string{"citizen", "floor_type"}) }

// CitizenIndustryType 市民-行业分类
func CitizenIndustryType() string { return getRandValue([]string{"citizen", "industry_type"}) }

// CitizenEnterpriseStatus 市民-企业状态
func CitizenEnterpriseStatus() string { return getRandValue([]string{"citizen", "enterprise_status"}) }

// CitizenEnterpriseType 市民-企业类型
func CitizenEnterpriseType() string { return getRandValue([]string{"citizen", "enterprise_type"}) }

// CitizenPoliticsStatus 市民-政治面貌
func CitizenPoliticsStatus() string { return getRandValue([]string{"citizen", "politics_status"}) }

// CitizenEthnicity 市民-民族
func CitizenEthnicity() string { return getRandValue([]string{"citizen", "politics_status"}) }

// CitizenEducationLevel 市民-受教育程度
func CitizenEducationLevel() string { return getRandValue([]string{"citizen", "education_level"}) }

// CitizenReligiousBelief 市民-宗教信仰
func CitizenReligiousBelief() string { return getRandValue([]string{"citizen", "religious_belief"}) }
