package externalqueue

//FeatureStub implements FeatureInteface
//shouldonly be used or testing
type FeatureStub struct {
}

var RetInt int
var RetBool bool
var RetString string

func (o *FeatureStub) Int(tag string) int {
	return RetInt
}

func (o *FeatureStub) Bool(tag string) bool {
	return RetBool
}

func (o *FeatureStub) Str(tag string) string {
	return RetString
}

func (o *FeatureStub) SetGlobalContext(context string) {
}

func (o *FeatureStub) SetDefaultContext(context string) {
}

func (o *FeatureStub) Init() {
}
