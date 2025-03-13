package dto

import "gorm.io/gorm/clause"

type ModifyParam struct {
	SingleTransaction bool        `json:"single_transaction,omitempty"`
	Conflict          Conflict    `json:"conflict,omitempty"`
	Data              interface{} `json:"data,omitempty"`
}

type Conflict struct {
	Columns      []string `json:"columns,omitempty"`
	OnConstraint string   `json:"on_constraint,omitempty"`
	DoNothing    bool     `json:"do_nothing,omitempty"`
	UpdateAll    bool     `json:"update_all,omitempty"`
	DoUpdates    []string `json:"do_updates,omitempty"`
}

func ConvertConflict(cf Conflict) clause.OnConflict {
	cols := make([]clause.Column, 0, len(cf.Columns))
	for _, v := range cf.Columns {
		cols = append(cols, clause.Column{Name: v})
	}
	cfl := clause.OnConflict{
		Columns:      cols,
		OnConstraint: cf.OnConstraint,
		DoNothing:    cf.DoNothing,
		UpdateAll:    cf.UpdateAll,
		DoUpdates:    clause.AssignmentColumns(cf.DoUpdates),
	}
	return cfl
}
