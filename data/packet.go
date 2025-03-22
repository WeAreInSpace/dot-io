package data

import (
	"github.com/WeAreInSpace/mlish"
)

type FeildkitGroupData struct {
	Name         string
	Descriptions []string
	FeildkitData *mlish.Model[FeildkitData]
}

type FeildkitData struct {
	Type         string
	Name         string
	Descriptions []string
	Action       string //write, read
}

func NewFieldkitManager() *FieldkitManager {
	fields := mlish.NewModel[FeildkitGroupData]()
	return &FieldkitManager{
		Feilds: fields,
	}
}

type FieldkitManager struct {
	Feilds *mlish.Model[FeildkitGroupData]
}

func (fmgr *FieldkitManager) New(fieldGroupName string, feildGroupDesc ...string) *FeildkitGroup {
	feildModel := mlish.NewModel[FeildkitData]()
	feild := &FeildkitGroup{
		feildModel: feildModel,
	}

	if feildGroupDesc == nil {
		feildGroupDesc = []string{}
	}

	fmgr.Feilds.Add(
		&FeildkitGroupData{Name: fieldGroupName, FeildkitData: feildModel, Descriptions: feildGroupDesc},
	)

	return feild
}

type FeildGroupSchema struct {
	Name         string        `json:"name"`
	Descriptions []string      `json:"descriptions"`
	Feilds       []FeildSchema `json:"feilds"`
}

type FeildSchema struct {
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	Descriptions []string `json:"descriptions"`
	Action       string   `json:"action"` //write, read
}

func (fmgr *FieldkitManager) Export() []FeildGroupSchema {
	var feildGroups []FeildGroupSchema

	fmgr.Feilds.For(
		func(item *mlish.ForParams[FeildkitGroupData]) {
			var feilds []FeildSchema

			item.DataAddr().FeildkitData.For(
				func(item *mlish.ForParams[FeildkitData]) {
					feild := FeildSchema{
						Type:         item.DataAddr().Type,
						Name:         item.DataAddr().Name,
						Descriptions: item.DataAddr().Descriptions,
						Action:       item.DataAddr().Action,
					}
					feilds = append(feilds, feild)
				},
			)

			feildGroup := FeildGroupSchema{
				Name:         item.DataAddr().Name,
				Descriptions: item.DataAddr().Descriptions,
				Feilds:       feilds,
			}

			feildGroups = append(feildGroups, feildGroup)
		},
	)

	return feildGroups
}

type FeildkitGroup struct {
	feildModel *mlish.Model[FeildkitData]
}

func validateFeildkitParams(action string, feildType string, feildName string, feildDesc []string) *FeildkitData {
	feildData := &FeildkitData{}
	if feildName == "" {
		feildData.Name = "feild"
	} else {
		feildData.Name = feildName
	}

	feildData.Action = action
	feildData.Type = feildType

	if feildDesc == nil {
		feildData.Descriptions = []string{}
	} else {
		feildData.Descriptions = feildDesc
	}

	return feildData
}

const (
	T_Int32    = "integer-32bit"
	T_Int64    = "integer-64bit"
	T_Str      = "string"
	T_Json_Str = "json-string"
	T_ByteArr  = "byte-array"
)

//Write

func (f *FeildkitGroup) WriteInt32(feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("write", T_Int32, feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) WriteInt64(feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("write", T_Int64, feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) WriteString(feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("write", T_Str, feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) WriteStreamString(feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("write", T_Str, feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) WriteJson(feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("write", T_Json_Str, feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) WriteBytes(feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("write", T_ByteArr, feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) WriteStreamBytes(feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("write", T_ByteArr, feildName, feildDesc)
	f.feildModel.Add(feildData)
}

//Read

func (f *FeildkitGroup) ReadInt32(feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("read", T_Int32, feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) ReadInt64(feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("read", T_Int64, feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) ReadString(feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("read", T_Str, feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) ReadStreamString(feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("read", T_Str, feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) ReadJson(feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("read", T_Json_Str, feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) ReadBytes(feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("read", T_ByteArr, feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) ReadStreamBytes(feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("read", T_ByteArr, feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func Try(onError func(err error), cb ...error) {
	for _, err := range cb {
		if err != nil {
			onError(err)
		}
	}
}

func TryAndRuturn(onError func(err error) error, cb ...error) error {
	for _, err := range cb {
		if err != nil {
			return onError(err)
		}
	}
	return nil
}

func TryAndRuturnThis(cb ...error) error {
	for _, err := range cb {
		if err != nil {
			return err
		}
	}
	return nil
}
