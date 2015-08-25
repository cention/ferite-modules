package externalqueue

import (
	"c3/feature"
)

type FeatureInterface interface {
	Int(feature string) int
	Bool(feature string) bool
	Str(feature string) string
	Init()
	SetGlobalContext(string)
	SetDefaultContext(string)
}

//FeatureObject uses the Feature Interface
type FeatureObject struct {
	fo FeatureInterface
}

func (o *FeatureObject) InitFeature(obj FeatureInterface) {
	o.fo = obj
	o.fo.Init()
}

func (o *FeatureObject) Int(tag string) int {
	return o.fo.Int(tag)
}

func (o *FeatureObject) Bool(tag string) bool {
	return o.fo.Bool(tag)
}

func (o *FeatureObject) Str(tag string) string {
	return o.fo.Str(tag)
}

func (o *FeatureObject) SetGlobalContext(context string) {
	o.fo.SetGlobalContext(context)
}

func (o *FeatureObject) SetDefaultContext(context string) {
	o.fo.SetDefaultContext(context)
}

//type FeatureClass implements FeatureInterface
type FeatureClass struct {
	f *feature.Feature
}

func (o *FeatureClass) Init() {
	o.f = feature.New()
}

func (o *FeatureClass) Int(tag string) int {
	return o.f.Int(tag)
}

func (o *FeatureClass) Bool(tag string) bool {
	return o.f.Bool(tag)
}

func (o *FeatureClass) Str(tag string) string {
	return o.f.Str(tag)
}
func (o *FeatureClass) SetGlobalContext(context string) {
	o.f.SetGlobalContext(context)
}

func (o *FeatureClass) SetDefaultContext(context string) {
	o.f.SetDefaultContext(context)
}
