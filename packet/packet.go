package packet

import (
	"io"

	"github.com/WeAreInSpace/mlish"
)

type FeildGroupData struct {
	Name      string
	FeildData *mlish.Model[FeildData]
}

type FeildData struct {
	Type         string
	Name         string
	Descriptions []string
	Action       string //write, read
}

func NewFieldManager() *FieldManager {
	fields := mlish.NewModel[FeildGroupData]()
	return &FieldManager{
		Feilds: fields,
	}
}

type FieldManager struct {
	Feilds *mlish.Model[FeildGroupData]
}

func (fmgr *FieldManager) New(fieldGroupName string) *FeildGroup {
	feildModel := mlish.NewModel[FeildData]()
	feild := &FeildGroup{
		feildModel: feildModel,
	}

	fmgr.Feilds.Add(
		&FeildGroupData{Name: fieldGroupName, FeildData: feildModel},
	)

	return feild
}

type FeildGroupSchema struct {
	Name   string        `json:"name"`
	Feilds []FeildSchema `json:"feilds"`
}

type FeildSchema struct {
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	Descriptions []string `json:"descriptions"`
	Action       string   `json:"action"` //write, read
}

func (fmgr *FieldManager) Export() []FeildGroupSchema {
	var feildGroups []FeildGroupSchema

	fmgr.Feilds.For(
		func(item *mlish.ForParams[FeildGroupData]) {
			var feilds []FeildSchema

			item.DataAddr().FeildData.For(
				func(item *mlish.ForParams[FeildData]) {
					feild := FeildSchema{
						Type:         item.DataAddr().Type,
						Name:         item.DataAddr().Name,
						Descriptions: item.DataAddr().Descriptions,
						Action:       item.DataAddr().Action,
					}
					feilds = append(feilds, feild)
				},
			)

			feildGroup := &FeildGroupSchema{
				Name:   item.DataAddr().Name,
				Feilds: feilds,
			}

			feildGroups = append(feildGroups, *feildGroup)
		},
	)

	return feildGroups
}

type FeildGroup struct {
	feildModel *mlish.Model[FeildData]
}

func validateFeildParams(action string, feildType string, feildName string, feildDesc []string) *FeildData {
	feildData := &FeildData{}
	if feildName == "" {
		feildData.Name = "feild"
	} else {
		feildData.Name = feildName
	}

	feildData.Action = action
	feildData.Type = feildType
	feildData.Descriptions = feildDesc

	return feildData
}

//Write

func (f *FeildGroup) WriteInt32(data int32, feildName string, feildDesc ...string) {
	feildData := validateFeildParams("write", "integer-32bit", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildGroup) WriteInt64(data int64, feildName string, feildDesc ...string) {
	feildData := validateFeildParams("write", "integer-64bit", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildGroup) WriteString(data string, feildName string, feildDesc ...string) {
	feildData := validateFeildParams("write", "string", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildGroup) WriteStreamString(len int64, data io.Reader, feildName string, feildDesc ...string) {
	feildData := validateFeildParams("write", "string", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildGroup) WriteJson(data any, feildName string, feildDesc ...string) {
	feildData := validateFeildParams("write", "json-string", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildGroup) WriteBytes(data []byte, feildName string, feildDesc ...string) {
	feildData := validateFeildParams("write", "byte-array", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildGroup) WriteStreamBytes(len int64, data io.Reader, feildName string, feildDesc ...string) {
	feildData := validateFeildParams("write", "byte-array", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

//Read

func (f *FeildGroup) ReadInt32(feildName string, feildDesc ...string) {
	feildData := validateFeildParams("read", "integer-32bit", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildGroup) ReadInt64(feildName string, feildDesc ...string) {
	feildData := validateFeildParams("read", "integer-64bit", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildGroup) ReadString(feildName string, feildDesc ...string) {
	feildData := validateFeildParams("read", "string", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildGroup) ReadStreamString(feildName string, feildDesc ...string) {
	feildData := validateFeildParams("read", "string", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildGroup) ReadJson(val any, feildName string, feildDesc ...string) {
	feildData := validateFeildParams("read", "json-string", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildGroup) ReadBytes(feildName string, feildDesc ...string) {
	feildData := validateFeildParams("read", "byte-array", feildName, feildDesc)
	f.feildModel.Add(feildData)
}

func (f *FeildGroup) ReadStreamBytes(feildName string, feildDesc ...string) {
	feildData := validateFeildParams("read", "byte-array", feildName, feildDesc)
	f.feildModel.Add(feildData)
}
