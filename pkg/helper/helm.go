package helper

//func ActionCfg(ns string, settings *helm.EnvSettings) (*action.Configuration, error) {
//	cfg := new(action.Configuration)
//	helmDriver := os.Getenv("HELM_DRIVER")
//	if ns == "" {
//		ns = settings.Namespace()
//	}
//
//	err := cfg.Init(settings.RESTClientGetter(), ns, helmDriver, log.Debugf)
//	return cfg, err
//}
//
//func SetNS(ns string) (*helm.EnvSettings, error) {
//	env := helm.New()
//	fs := &pflag.FlagSet{}
//	env.AddFlags(fs)
//	flag := fs.Lookup("namespace")
//	err := flag.Value.Set(ns)
//	if err != nil {
//		return nil, err
//	}
//
//	return env, nil
//}
