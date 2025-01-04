# Go 如何优雅地读取配置？


{{&lt; admonition type=abstract title=&#34;导语&#34; open=true &gt;}}

{{&lt; /admonition &gt;}}

&lt;!--more--&gt;

{{&lt; figure src=&#34;/posts/go-viper/logo.png&#34; title=&#34;&#34; &gt;}}

**对于一个 Go 应用程序，同城需要解析以下类别的配置：命令行选项、命令行参数、配置文件**，而对于一个非命令行工具的应用程序，不需要考虑读取命令行参数这类场景，其需要的配置内容都可以通过命令行选项或配置文件加载到程序中。

{{&lt; admonition type=Tip title=&#34;Tips&#34; open=true &gt;}}
命令行工具可能会有子命令，例如 `kubectr create` 中的 `create` 就是一个命令行参数
{{&lt; /admonition &gt;}}

## 为何选择配置文件作为配置项的读取方式？

对于一个配置项，既可以通过命令行选项，又能够通过配置文件来读取，而且二者是一个彼此可以取代的，因此，对于非命令行工具的应用程序个人更倾向于通过配置文件完成，原因如下：

- **配置文件更易部署**：可以将应用所需要的所有配置聚合在一个配置文件中。
当需署时，只需要部署、加载这个配置文件即可，不需要配置一大堆命令行选项；
- **配置文件更易维护**：将所有的配置项都保存在配置文件中，加上详细的配置说明，不需要的配置项可以注释掉。
一个具有全量配置项、详细说明的配置文件，更易于理解。并且在修改时，只需要修改配置项的值，而不需要修改配置项名称，更不易出错；
- **配置文件可以实现热加载功能**：应用程序监听配置文件的变更，有变更时，自动重新加载配置文件，实现配置热加载功能；
- **配置层次表达更清晰**：命令行参数无法直接表达&#34;层次&#34;，但配置文件可以。层次化的配置表达，更易于理解，也更易于维护。
- **方便新增配置项**：多数情况下，配置项新增只需在配置文件中新增一行即可，不需要修改源码；

{{&lt; admonition type=Tip title=&#34;总结&#34; open=true &gt;}}
命令行工具可能会有子命令，例如 `kubectr create` 中的 `create` 就是一个命令行参数
总结：当配置项少的时候（比如 5 个以内），可以从命令行选项中读取。
参数较多的时候，建议使用配置文件，配置文件更易部署、维护、热加载、层次表达更清晰。
{{&lt; /admonition &gt;}}

## 为何选择 YAML 作为配置文件的格式？

当打算采用配置文件来读取配置项时，那么就存在多种文件格式，例如：JSON、YAML、TOML、INI 等。
个人推荐使用 YAML，理由如下：

- YAML 语法简单、格式易读、程序易处理；
- YAML 格式可以表达非常丰富、复杂的配置结构；
- YAML 格式普适性高，新人零理解成本；

&gt; 最终配置：使用 YAML 格式的配置文件，并采用 `viper` 来读取配置

---

## 使用 viper 读取配置文件内容

在 [浅析现代化命令行框架 Cobra](https://kydance.github.io/posts/go-cobra/) 中，我们了解到可以通过 `cobra-cli init --viper` 生成一个通过 viper 来配置应用程序的 Demo 应用，那么就可以知道它的应用加载逻辑如下：

```go
/*
Copyright © 2024 Kyden &lt;kytedance@gmail.com&gt;
This file is part of CLI application foo.
*/
package cmd

import (
	&#34;fmt&#34;
	&#34;os&#34;

	&#34;github.com/spf13/cobra&#34;
	&#34;github.com/spf13/viper&#34;
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &amp;cobra.Command{
	Use:   &#34;kydendemo&#34;,
	Short: &#34;A brief description of your application&#34;,
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(
    &amp;cfgFile, &#34;config&#34;, &#34;&#34;, &#34;config file (default is $HOME/.kydendemo.yaml)&#34;)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP(&#34;toggle&#34;, &#34;t&#34;, false, &#34;Help message for toggle&#34;)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != &#34;&#34; {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name &#34;.kydendemo&#34; (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType(&#34;yaml&#34;)
		viper.SetConfigName(&#34;.kydendemo&#34;)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, &#34;Using config file:&#34;, viper.ConfigFileUsed())
	}
}
```

其中，`rootCmd` 设置了命令行选项 `--config`，用于指定配置文件路径，默认值是 `&#34;&#34;`；

通过 `cobra.OnInitialize(initConfig)` 设置了 `kydendemo` 在运行时的回调函数 `initConfig`，
它的执行逻辑主要是：

- 如果指定了 `cfgFile`，则直接读取该配置文件；
- 如果没有指定，则读取 `$HOME/.kydendemo.yaml`，找到则读取；
若 `cfgFile == &#34;&#34;`，且没有找到配置文件，则调用 `viper.ReadInConfig()` 读取配置文件时报错；

## Reference

- [viper](https://github.com/spf13/viper)


---

> : [kyden](https:github.com/kydance)  
> URL: http://kyden.us.kg/posts/go-viper/  

