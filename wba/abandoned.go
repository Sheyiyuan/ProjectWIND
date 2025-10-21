package wba

//type APP interface {
//	Get() AppInfo
//	Init(Wba interface{}) error
//	//Init(WspApi WindStandardProtocolAPI, WsdApi WindStandardDataBaseAPI, WstApi WindStandardTools) error
//	//InitWSD(api WindStandardDataBaseAPI) error
//}

//func (ai AppInfo) Get() AppInfo {
//	return ai
//}
//
//func (ai AppInfo) Init(Wba interface{}) error {
//	WBA = Wba.(WindStandardTools)
//	return nil
//}

//func (ai *AppInfo) Init(WspApi WindStandardProtocolAPI, WsdApi WindStandardDataBaseAPI, WstApi WindStandardTools) error {
//	WSP = WspApi
//	WSD = WsdApi
//	WST = WstApi
//	return nil
//}

//func (ai *AppInfo) InitWSD(api WindStandardDataBaseAPI) error {
//	WSD = api
//	return nil
//}

//func WithName(name string) AppInfoOption {
//	return func(ei *AppInfo) {
//		ei.AppKey.Name = name
//	}
//}
//
//func WithVersion(version string) AppInfoOption {
//	return func(ei *AppInfo) {
//		ei.AppKey.Version = version
//	}
//}
//
//func WithAuthor(author string) AppInfoOption {
//	return func(ei *AppInfo) {
//		ei.Author = author
//	}
//}
