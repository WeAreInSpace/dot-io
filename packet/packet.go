package packet

import (
	"io"

	"github.com/WeAreInSpace/mlish"
)

type FeildkitGroupData struct {
	Name         string
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

func (fmgr *FieldkitManager) New(fieldGroupName string) *FeildkitGroup {
	feildModel := mlish.NewModel[FeildkitData]()
	feild := &FeildkitGroup{
		feildModel: feildModel,
	}

	fmgr.Feilds.Add(
		&FeildkitGroupData{Name: fieldGroupName, FeildkitData: feildModel},
	)

	return feild
}

type FeildkitGroupSchema struct {
	Name   string           `json:"name"`
	Feilds []FeildkitSchema `json:"feilds"`
}

type FeildkitSchema struct {
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	Descriptions []string `json:"descriptions"`
	Action       string   `json:"action"` //write, read
}

func (fmgr *FieldkitManager) Export() []FeildkitGroupSchema {
	var feildGroups []FeildkitGroupSchema

	fmgr.Feilds.For(
		func(item *mlish.ForParams[FeildkitGroupData]) {
			var feilds []FeildkitSchema

			item.DataAddr().FeildkitData.For(
				func(item *mlish.ForParams[FeildkitData]) {
					feild := FeildkitSchema{
						Type:         item.DataAddr().Type,
						Name:         item.DataAddr().Name,
						Descriptions: item.DataAddr().Descriptions,
						Action:       item.DataAddr().Action,
					}
					feilds = append(feilds, feild)
				},
			)

			feildGroup := &FeildkitGroupSchema{
				Name:   item.DataAddr().Name,
				Feilds: feilds,
			}

			feildGroups = append(feildGroups, *feildGroup)
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

//Write

func (f *FeildkitGroup) WriteInt32(data int32, feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("write", "integer-32bit", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) WriteInt64(data int64, feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("write", "integer-64bit", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) WriteString(data string, feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("write", "string", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) WriteStreamString(len int64, data io.Reader, feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("write", "string", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) WriteJson(data any, feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("write", "json-string", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) WriteBytes(data []byte, feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("write", "byte-array", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) WriteStreamBytes(len int64, data io.Reader, feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("write", "byte-array", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

//Read

func (f *FeildkitGroup) ReadInt32(feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("read", "integer-32bit", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) ReadInt64(feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("read", "integer-64bit", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) ReadString(feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("read", "string", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) ReadStreamString(feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("read", "string", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) ReadJson(val any, feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("read", "json-string", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) ReadBytes(feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("read", "byte-array", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildkitGroup) ReadStreamBytes(feildName string, feildDesc ...string) {
	feildData := validateFeildkitParams("read", "byte-array", feildName, feildDesc)
	f.feildModel.Add(feildData)
}
