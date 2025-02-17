package cmd

import (
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
)

type CmdRootFile struct {
	*generator.BaseGenerator
	Meta      *generator.Meta
	AppName   string
	GlobalDir string
}

const (
	RootFilename  = "root.go"
	VarEnv        = "env"
	VarConfigFile = "configFile"
)

func NewCmdRootFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(meta.Dirs[g.Cmd], RootFilename, false)
	if err != nil {
		return nil, err
	}

	return &CmdRootFile{
		BaseGenerator: baseGenerator,
		Meta:          meta,
		AppName:       meta.RootDir,
		GlobalDir:     meta.Dirs[g.Global],
	}, nil
}

func (f CmdRootFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genImports,
		f.genConst,
		f.genVar,
		f.genInitFunc,
		f.genExecuteFunc,
		f.genLoadConfigFunc,
	}
}

func (f *CmdRootFile) genImports() {
	f.JenFile.ImportName(f.GlobalDir, "g")
}

// genMain generate main function
// const (
//	APP = "app"
// )
func (f CmdRootFile) genConst() {
	found, _ := f.FindConst("APP")
	if found == nil {
		f.Builder.Raw().Const().Defs(
			jen.Id("APP").Op("=").Lit(f.AppName),
		).Line().Line()
	}
}

// var(
//  env        string
//  configFile string
// )
//var rootCmd = &cobra.Command{
//	Use:   APP,
//	Short: "bd server",
//	Long:  `bd server serves for all kinds of API`,
//}
func (f CmdRootFile) genVar() {
	found, _ := f.FindVar(VarEnv)
	if found == nil {
		f.Builder.Raw().Var().Id(VarEnv).String().Line()
	}

	found, _ = f.FindVar(VarConfigFile)
	if found == nil {
		f.Builder.Raw().Var().Id(VarConfigFile).String().Line()
	}

	found, _ = f.FindVar("rootCmd")
	if found == nil {
		f.Builder.Raw().Var().Id("rootCmd").Op("=").Id("&").Qual(g.ImportPaths[g.Cobra], "Command").Values(
			jen.Dict{
				jen.Id("Use"):   jen.Lit(f.AppName),
				jen.Id("Short"): jen.Lit(f.AppName + " short description"),
				jen.Id("Long"):  jen.Lit(f.AppName + " long description"),
			},
		).Line()
	}
}

//func init() {
//	cobra.OnInitialize(loadConfig)
//
//	rootCmd.PersistentFlags().StringP("env", "e", "", "running environment, e,g: [prod, sim, pre, test, dev, local]")
//	rootCmd.PersistentFlags().StringP("config", "c", "", "config file, default: config.toml")
//	rootCmd.AddCommand(runServerCmd)
//}
func (f CmdRootFile) genInitFunc() {
	found, _ := f.FindMethod("init")
	if found == nil {
		body := []jen.Code{
			jen.Qual(g.ImportPaths[g.Cobra], "OnInitialize").Call(jen.Id("loadConfig")),
			jen.Line(),
			jen.Id("rootCmd").Dot("PersistentFlags").Call().Dot("StringVarP").Call(
				jen.Op("&").Id(VarEnv), jen.Lit("env"), jen.Lit("e"), jen.Lit(""), jen.Lit("running environment, e,g: [prod, sim, pre, test, dev, local]"),
			),
			jen.Id("rootCmd").Dot("PersistentFlags").Call().Dot("StringVarP").Call(
				jen.Op("&").Id(VarConfigFile), jen.Lit("config"), jen.Lit("c"), jen.Lit(""), jen.Lit("config file, default: config.toml"),
			),
			jen.Id("rootCmd").Dot("AddCommand").Call(jen.Id("runCmd")),
		}

		f.Builder.AppendFunction(
			"init",
			nil,
			nil,
			nil,
			"",
			body...,
		)
		f.Builder.NewLine()
	}
}

//func Execute() {
//	// 尝试捕获panic并保存到错误中
//	defer func() {
//		if r := recover(); r != nil {
//			hdsdk.RecordErrorStack(APP)
//		}
//	}()
//
//	if err := rootCmd.Execute(); err != nil {
//		hdsdk.Shutdown()
//		os.Exit(1)
//	}
//}
func (f CmdRootFile) genExecuteFunc() {
	found, _ := f.FindMethod("Execute")
	if found == nil {
		body := []jen.Code{
			jen.Defer().Func().Params().Block(
				jen.If(
					jen.Id("r").Op(":=").Id("recover").Call(),
					jen.Id("r").Op("!=").Nil(),
				).Block(
					jen.Qual(g.ImportPaths[g.HdUtils], "RecordErrorStack").Call(jen.Lit(f.AppName)),
				),
			).Call(),
			jen.If(
				jen.Err().Op(":=").Id("rootCmd").Dot("Execute").Call(),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Qual("os", "Exit").Call(jen.Lit(1)),
			),
		}

		f.Builder.AppendFunction(
			"Execute",
			nil,
			nil,
			nil,
			"",
			body...,
		)
		f.Builder.NewLine()
	}
}

//func loadConfig() {
//	// 尝试从各种源加载配置信息
//	v := hdsdk.LoadConfig(APP, env, configFile)
//
//	// 将配置信息转换成对应的数据结构
//	err := v.Unmarshal(&g.Config)
//	if err != nil {
//		utils.Fatal("msg", "unmarshal config", "err", err)
//	}
//
//}
func (f CmdRootFile) genLoadConfigFunc() {
	found, _ := f.FindMethod("loadConfig")
	if found == nil {
		body := []jen.Code{
			jen.Id("v").Op(":=").Qual(g.ImportPaths[g.HdSdk], "LoadConfig").Call(
				jen.Id("APP"), jen.Id(VarEnv), jen.Id(VarConfigFile),
			),
			jen.Line(),
			jen.Id("err").Op(":=").Id("v").Dot("Unmarshal").Call(jen.Op("&").Qual(f.GlobalDir, "Config")),
			jen.If(jen.Id("err").Op("!=").Nil().Block(
				jen.Qual(g.ImportPaths[g.HdUtils], "Fatal").Call(
					jen.Lit("msg"), jen.Lit("unmarshal config"), jen.Lit("err"), jen.Err(),
				),
			)),
		}

		f.Builder.AppendFunction(
			"loadConfig",
			nil,
			nil,
			nil,
			"",
			body...,
		)
		f.Builder.NewLine()
	}
}
